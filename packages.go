/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/electrikmilk/args-parser"
	"github.com/go-git/go-git/v5"
	"howett.net/plist"
)

/*
Package management through GitHub repositories
*/

var currentPkg *cherriPackage
var pkgRegex = regexp.MustCompile(`^@([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)$`)

type cherriPackage struct {
	Name     string
	User     string
	Archived bool
	Packages []cherriPackage
}

func (pkg *cherriPackage) installed() (installed bool) {
	if _, statErr := os.Stat(pkg.path()); os.IsNotExist(statErr) {
		return false
	}
	return true
}

func (pkg *cherriPackage) install() (installed bool) {
	var packagePath = pkg.path()
	fmt.Println(fmt.Sprintf("Installing %s from %s...", pkg.signature(), pkg.url()))
	var _, cloneErr = git.PlainClone(packagePath, false, &git.CloneOptions{
		URL: pkg.url(),
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

func (pkg *cherriPackage) uninstall() {
	for i, dep := range currentPkg.Packages {
		if dep.signature() == pkg.signature() {
			var packagePath = pkg.path()
			if _, pkgStatErr := os.Stat(packagePath); os.IsNotExist(pkgStatErr) {
				break
			}
			var gitDirRemoveErr = os.RemoveAll(packagePath)
			handle(gitDirRemoveErr)
			currentPkg.Packages = append(currentPkg.Packages[:i], currentPkg.Packages[i+1:]...)
		}
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

// url returns the GitHub repository git URL of the package.
func (pkg *cherriPackage) url() string {
	return fmt.Sprintf("https://github.com/%s/cherri-%s.git", pkg.User, pkg.Name)
}

// signature returns a formatted string of the author and name of the package.
func (pkg *cherriPackage) signature() string {
	return fmt.Sprintf("@%s/%s", pkg.User, pkg.Name)
}

// path returns the absolute path of the package.
func (pkg *cherriPackage) path() string {
	return fmt.Sprintf("./packages/%s", pkg.signature())
}

// initPackage initializes a package in the current directory using an info.plist file based on cherriPackage.
func initPackage() {
	if _, statErr := os.Stat("info.plist"); !os.IsNotExist(statErr) {
		exit("info.plist already exists. Delete it to create new package.")
	}
	var pkgSig = args.Value("init")
	var newPkg = createPackage(pkgSig)
	currentPkg = &newPkg
	writePackage()

	fmt.Println(ansi("Initialized Cherri package", green))
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
	var marshaledPlist, plistErr = plist.MarshalIndent(currentPkg, plist.XMLFormat, "\t")
	handle(plistErr)

	compiled = string(marshaledPlist)
	writeFile("info.plist", "info.plist")
}

// tidyPackage re-installs all packages in the package in the current directory.
func tidyPackage() {
	if pkg, found := loadPackage("info.plist"); found {
		currentPkg = pkg
		installPackages(currentPkg.Packages, true)
		return
	}

	exit("info.plist does not exist. Use --init argument to initialize a package.")
}

// addPackage adds a package to the dependencies for the package in the current directory and triggers lazy installation.
func addPackage() {
	if pkg, found := loadPackage("info.plist"); found {
		currentPkg = pkg
		var name = args.Value("install")
		var newPkg = createPackage(name)
		if newPkg.installed() || addedPackage(&newPkg) {
			exit(fmt.Sprintf("Package %s already installed.", newPkg.signature()))
		}

		fmt.Println(ansi(fmt.Sprintf("Install package %s\n", newPkg.signature()), green))

		var packagePrompt = fmt.Sprintf("Do you trust this package?\n\nThis will download this GitHub repository and automatically include it in this project:\n%s", newPkg.url())
		fmt.Println(ansi(packagePrompt, yellow))
		if !yesNo() {
			return
		}
		fmt.Println("")

		currentPkg.Packages = append(currentPkg.Packages, newPkg)
		installPackages(currentPkg.Packages, false)
		writePackage()
	} else {
		exit("install: info.plist does not exist. Use --init argument to initialize a package.")
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

// createPackage creates a cherriPackage type from a string matching pkgRegex.
func createPackage(name string) cherriPackage {
	var matches = pkgRegex.FindAllStringSubmatch(name, -1)
	if len(matches) == 0 {
		exit(fmt.Sprintf("Package must follow pattern: {github_username}/{repo_package_name}, got: %s", name))
	}
	var user = matches[0][1]
	var pkg = matches[0][2]

	return cherriPackage{
		Name: pkg,
		User: user,
	}
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
		if _, statErr := os.Stat(pkg.path()); os.IsNotExist(statErr) || tidy {
			if tidy && pkg.installed() {
				pkg.uninstall()
			}
			if pkg.install() {
				pkg.loadDependencies(tidy)
			}
		}
	}
}

// includePackages adds lines to include all the packages for the current directory.
func includePackages(path string) {
	var packages, dirErr = os.ReadDir(path)
	handle(dirErr)
	for _, user := range packages {
		includeUserPackages(fmt.Sprintf("%s/%s", path, user.Name()))
	}
	resetParse()
}

func includeUserPackages(path string) {
	var packages, dirErr = os.ReadDir(path)
	handle(dirErr)
	for _, pkg := range packages {
		var packageInclude = fmt.Sprintf("#include '%s/%s/main.cherri'\n", path, pkg.Name())
		lines = append([]string{packageInclude}, lines...)
	}
}

// removePackage uninstalls a package from the current directory.
func removePackage() {
	if pkg, found := loadPackage("info.plist"); found {
		currentPkg = pkg
		var name = args.Value("remove")
		var targetPkg = createPackage(name)
		if !targetPkg.installed() {
			exit(fmt.Sprintf("Package %s is not installed.", targetPkg.signature()))
		}

		targetPkg.uninstall()
		writePackage()
		fmt.Println(ansi(fmt.Sprintf("[-] %s removed: %s", targetPkg.signature(), targetPkg.path()), red))
	} else {
		exit("install: info.plist does not exist. Use --init argument to create a package.")
	}
}
