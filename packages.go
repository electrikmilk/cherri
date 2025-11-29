/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
	"github.com/go-git/go-git/v5"
	"howett.net/plist"
)

/*
Package management through remote Git repositories
*/

var currentPkg *cherriPackage
var visitedPackages []string
var pkgSignatureRegex = regexp.MustCompile(`^(https://(.*?)/)?@?([A-Za-z0-9_.-]+)/(?:cherri-)?([A-Za-z0-9_.-]+)$`)

type cherriPackage struct {
	Name     string
	User     string
	Uri      string
	Archived bool
	Packages []cherriPackage
}

func (pkg *cherriPackage) installed() (installed bool) {
	if _, statErr := os.Stat(pkg.path()); os.IsNotExist(statErr) {
		return false
	}
	return true
}

func (pkg *cherriPackage) trusted() bool {
	loadTrustedPackages()
	return slices.Contains(trusted.Packages, pkg.signature())
}

func (pkg *cherriPackage) install() (installed bool) {
	fmt.Println(fmt.Sprintf("Installing %s from %s...", pkg.signature(), pkg.url()))

	var packagePath = pkg.path()
	var _, cloneErr = git.PlainClone(packagePath, false, &git.CloneOptions{
		URL:          pkg.url(),
		SingleBranch: true,
		Depth:        5,
	})
	if cloneErr != nil {
		pkg.failed(cloneErr)
		return
	}

	var entryPath = fmt.Sprintf("%s/main.cherri", packagePath)
	if _, statErr := os.Stat(entryPath); os.IsNotExist(statErr) {
		pkg.failed(statErr)
		return
	}

	var infoPlist = fmt.Sprintf("%s/info.plist", packagePath)
	if _, infoStatErr := os.Stat(infoPlist); os.IsNotExist(infoStatErr) {
		pkg.failed(infoStatErr)
		return
	}

	fmt.Println(ansi(fmt.Sprintf("[+] %s installed: %s\n", pkg.signature(), pkg.path()), green))

	return true
}

func (pkg *cherriPackage) removeFiles() {
	var packagePath = pkg.path()
	if _, pkgStatErr := os.Stat(packagePath); os.IsNotExist(pkgStatErr) {
		return
	}
	var removeErr = os.RemoveAll(packagePath)
	handle(removeErr)
}

func (pkg *cherriPackage) removeFromManifest() {
	for i, dep := range currentPkg.Packages {
		if dep.signature() != pkg.signature() {
			continue
		}
		currentPkg.Packages = append(currentPkg.Packages[:i], currentPkg.Packages[i+1:]...)
		break
	}
}

func (pkg *cherriPackage) loadDependencies(reinstall bool) {
	if len(pkg.Packages) == 0 {
		return
	}
	if pkg, found := loadPackage(fmt.Sprintf("%s/info.plist", pkg.path())); found {
		installPackages(pkg.Packages, reinstall)
	}
}

// Reports that the package failed to install.
func (pkg *cherriPackage) failed(err error) {
	fmt.Println(ansi(fmt.Sprintf("[x] %s - unable to install: %s\n", pkg.signature(), err), yellow))
}

// url returns the repository git URL of the package.
func (pkg *cherriPackage) url() string {
	if pkg.Uri == "" {
		exit("Package has no URL host (e.g. github.com)!")
	}
	return fmt.Sprintf("https://%s/%s/cherri-%s.git", pkg.Uri, pkg.User, pkg.Name)
}

// signature returns a formatted string of the author and name of the package.
func (pkg *cherriPackage) signature() string {
	return fmt.Sprintf("@%s/%s", pkg.User, pkg.Name)
}

// path returns the absolute path of the package.
func (pkg *cherriPackage) path() string {
	return fmt.Sprintf("./packages/%s", pkg.signature())
}

// listPackage shows the current package info.
func listPackage() {
	if pkg, found := loadPackage("info.plist"); found {
		printPackage(pkg)
		if len(pkg.Packages) != 0 {
			fmt.Println(ansi("\nDependent packages:", green, underline))
			for _, pkg := range pkg.Packages {
				printPackage(&pkg)
			}
		}
	} else {
		initPackageError()
	}
}

// listPackages shows the current installed packages info.
func listPackages() {
	fmt.Println(ansi("Installed packages:\n", green))
	if pkg, found := loadPackage("info.plist"); found {
		if len(pkg.Packages) != 0 {
			for _, pkg := range pkg.Packages {
				printPackage(&pkg)
			}
		}
	} else {
		initPackageError()
	}
}

// Add host remote Git repository URI to current package.
func addUri() {
	if pkg, found := loadPackage("info.plist"); found {
		currentPkg = pkg
		currentPkg.Uri = args.Value("add-uri")
		writePackage()
		fmt.Println(ansi(fmt.Sprintf("Updated URI for package %s to %s", currentPkg.signature(), currentPkg.Uri), green))
	} else {
		initPackageError()
	}
}

func printPackage(pkg *cherriPackage) {
	var isArchived string
	if pkg.Archived {
		isArchived = " (archived)"
	}
	fmt.Println("-", ansi(pkg.signature(), blue))
	fmt.Println("\tName:", pkg.Name, isArchived)
	fmt.Println("\tUser:", pkg.User)
	if pkg.Uri == "" {
		fmt.Println("\tHosted: Private")
	} else {
		fmt.Println("\tHosted At:", pkg.Uri)
	}
	fmt.Println("\tDep. packages:", len(pkg.Packages))
	fmt.Println("\tInstall path:", pkg.path())
}

// initPackage initializes a package in the current directory using an info.plist file based on cherriPackage.
func initPackage() {
	if _, statErr := os.Stat("info.plist"); !os.IsNotExist(statErr) {
		exit("info.plist already exists. Delete it to create new package.")
	}
	var pkgSig = args.Value("init")
	var newPkg = newPackage(pkgSig)
	currentPkg = &newPkg
	writePackage()

	fmt.Println(ansi("Initialized Cherri package", green))
}

var trustedPackagesPlistPath = os.ExpandEnv("$HOME/.cherri/trusted.plist")

type trustedPackages struct {
	Packages []string
}

var trusted trustedPackages

func loadTrustedPackages() {
	createInternalDir()

	if _, statErr := os.Stat(trustedPackagesPlistPath); os.IsNotExist(statErr) {
		return
	}
	var trustPlist, readErr = os.ReadFile(trustedPackagesPlistPath)
	handle(readErr)

	var _, plistErr = plist.Unmarshal(trustPlist, &trusted)
	handle(plistErr)
}

func writeTrustedPackages() {
	if _, statErr := os.Stat(internalDirectoryPath); os.IsNotExist(statErr) {
		var intDirErr = os.Mkdir(internalDirectoryPath, 0777)
		handle(intDirErr)

		var trustedPlistFile, createErr = os.Create(trustedPackagesPlistPath)
		handle(createErr)

		defer func(f *os.File) {
			var fileCloseErr = f.Close()
			handle(fileCloseErr)
		}(trustedPlistFile)
	}

	var marshaledPlist, plistErr = plist.MarshalIndent(trusted, plist.XMLFormat, "\t")
	handle(plistErr)

	var writeErr = os.WriteFile(trustedPackagesPlistPath, marshaledPlist, 0600)
	handle(writeErr)
}

// loadPackage loads the package in the current directory.
func loadPackage(path string) (pkg *cherriPackage, found bool) {
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return nil, false
	}
	var pkgPlist, pkgPlistError = os.ReadFile(path)
	handle(pkgPlistError)

	var info cherriPackage
	var _, plistErr = plist.Unmarshal(pkgPlist, &info)
	handle(plistErr)

	return &info, true
}

// writePackage writes the current package to the info.plist file in the current directory.
func writePackage() {
	var writeDebugOutput = args.Using("debug")
	if writeDebugOutput {
		fmt.Printf("Writing to info.plist...")
	}

	var info, infoErr = os.OpenFile("info.plist", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	handle(infoErr)
	defer info.Close()

	var plistEncoder = plist.NewEncoder(info)
	if args.Using("debug") {
		plistEncoder.Indent("\t")
	}

	var encodeErr = plistEncoder.Encode(currentPkg)
	handle(encodeErr)

	if writeDebugOutput {
		fmt.Println(ansi("Done.", green))
	}
}

// tidyPackage re-installs all packages in the package in the current directory.
func tidyPackage() {
	if pkg, found := loadPackage("info.plist"); found {
		currentPkg = pkg
		installPackages(currentPkg.Packages, true)
	} else {
		initPackageError()
	}
}

// newPackage creates a cherriPackage type from a string matching pkgSignatureRegex.
func newPackage(name string) cherriPackage {
	var matches = pkgSignatureRegex.FindAllStringSubmatch(name, -1)
	if len(matches) == 0 {
		exit(fmt.Sprintf("Invalid package signature: %s\nPackage signature must follow one of these patterns:\n\nPrivate: @{author}/{package_name}\n\nor\n\nPublic Remote Git repository: https://{domain}/{username}/cherri-{package_name}", name))
	}
	var uri = matches[0][2]
	var user = matches[0][3]
	var pkg = matches[0][4]

	if uri == user {
		uri = ""
	}

	return cherriPackage{
		Uri:  uri,
		Name: pkg,
		User: user,
	}
}

// addPackage adds a package to the dependencies for the package in the current directory and triggers lazy installation.
func addPackage() {
	var pkg, found = loadPackage("info.plist")
	if !found {
		initPackageError()
	}

	currentPkg = pkg
	var name = args.Value("install")
	var newPkg = newPackage(name)
	if newPkg.installed() || addedPackage(&newPkg) {
		exit(fmt.Sprintf("Package %s already installed.", newPkg.signature()))
	}

	fmt.Println(ansi(fmt.Sprintf("Install package %s\n", newPkg.signature()), green))

	checkTrustedPackages(&newPkg)

	currentPkg.Packages = append(currentPkg.Packages, newPkg)
	installPackages(currentPkg.Packages, false)
	writePackage()
}

func checkTrustedPackages(newPkg *cherriPackage) {
	if !newPkg.trusted() {
		var packagePrompt = fmt.Sprintf("Do you trust this package?\n\nThis will download this remote git repository and automatically include it in this project:\n%s", newPkg.url())
		fmt.Println(ansi(packagePrompt, yellow))
		if !yesNo() {
			return
		}
		fmt.Print("\n")

		trusted.Packages = append(trusted.Packages, newPkg.signature())
		writeTrustedPackages()
	}
}

func addedPackage(pkg *cherriPackage) (added bool) {
	for _, p := range currentPkg.Packages {
		if p.signature() == pkg.signature() {
			return true
		}
	}
	return
}

// installPackages installs the given dependencies.
func installPackages(packages []cherriPackage, tidy bool) {
	if len(packages) == 0 {
		return
	}

	for _, pkg := range packages {
		if pkg.Archived {
			fmt.Println(ansi(fmt.Sprintf("[!] Archived package: %s", pkg.Name), yellow))
		}
		if slices.Contains(visitedPackages, pkg.path()) {
			continue
		}
		if _, statErr := os.Stat(pkg.path()); os.IsNotExist(statErr) || tidy {
			visitedPackages = append(visitedPackages, pkg.path())
			if tidy && pkg.installed() {
				pkg.removeFiles()
			}
			if pkg.install() {
				pkg.loadDependencies(tidy)
			}
		}
	}
}

// includePackages adds lines to include all the packages for the current directory.
// Sorts packages in deterministic order. Only includes files in info.plist.
func includePackages() {
	if len(currentPkg.Packages) == 0 {
		return
	}

	var sortedPackages = make([]cherriPackage, len(currentPkg.Packages))
	copy(sortedPackages, currentPkg.Packages)
	slices.SortFunc(sortedPackages, func(a, b cherriPackage) int {
		return strings.Compare(a.signature(), b.signature())
	})

	for _, pkg := range sortedPackages {
		var packageInclude = fmt.Sprintf("#include './packages/%s/main.cherri'\n", pkg.signature())
		lines = append([]string{packageInclude}, lines...)
	}

	resetParse()
}

// removePackage uninstalls a package from the current directory.
func removePackage() {
	if pkg, found := loadPackage("info.plist"); found {
		currentPkg = pkg
		var name = args.Value("remove")
		var targetPkg = newPackage(name)
		if !targetPkg.installed() {
			exit(fmt.Sprintf("Package %s is not installed.", targetPkg.signature()))
		}

		targetPkg.removeFiles()
		targetPkg.removeFromManifest()
		writePackage()
		fmt.Println(ansi(fmt.Sprintf("[-] %s removed: %s", targetPkg.signature(), targetPkg.path()), red))
	} else {
		initPackageError()
	}
}

func initPackageError() {
	exit("install: info.plist does not exist. Use --init argument to create a package.")
}
