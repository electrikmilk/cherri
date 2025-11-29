<img width="200" height="200" alt="Cherri Icon" src="https://github.com/user-attachments/assets/a9c23532-a1df-41ec-bd5b-6621f54064c8" />

# Cherri

[![Build & Test](https://github.com/electrikmilk/cherri/actions/workflows/build-test.yml/badge.svg)](https://github.com/electrikmilk/cherri/actions/workflows/build-test.yml)
[![Releases](https://img.shields.io/github/v/release/electrikmilk/cherri?include_prereleases)](https://github.com/electrikmilk/cherri/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/go.mod)
[![License](https://img.shields.io/github/license/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/LICENSE)
![Platform](https://img.shields.io/badge/platform-macOS-red)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/electrikmilk/cherri?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrikmilk/cherri)](https://goreportcard.com/report/github.com/electrikmilk/cherri)

**Cherri** (pronounced cherry) is a [Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752) programming language that compiles directly to a valid runnable Shortcut.

The primary goal is to make it practical to create large Shortcut projects (within the limitations of Shortcuts) and maintain them long term.

[![Hello World Example](https://github.com/user-attachments/assets/dc2ce82e-f85e-44ab-9f43-50f9518ddcda)](https://playground.cherrilang.org)

### 🌟 Top Features

- 🖥️ Laptop/Desktop-based development (CLI, [VSCode extension](https://marketplace.visualstudio.com/items?itemName=electrikmilk.cherri-vscode-extension), macOS app)
- 🎓 Easy to learn and syntax similar to other languages
- 🐞 1-1 translation to Shortcut actions as much as possible to make debugging easier
- 🥾 Half-bootstrapped: Most actions and types are written in the language
- 💻 Import actions on your Mac
- 📦 Package manager: Remote Git repo-based package manager built in, allowing for automatic inclusion and updates.
- 🪄 No magic variables syntax, they're constants instead
- 🪶 Optimized to create as small as possible Shortcuts and reduce memory usage at runtime
- #️⃣ Include files within others for large Shortcut projects
- 🔧 Define actions with type checking, enums, optionals, default values, raw identifiers, and raw keys.
- 🔄 Define functions to run within their own scope at the top of your Shortcut to reduce duplicate actions.
- 📋 Copy-paste actions automatically
- 🥩 Enter action identifier and parameters manually using Raw Actions.
- ❓ Define import questions
- 📇 Generate VCards for menus
- 📄 Embed files in base64
- 🔀 Convert Shortcuts from an iCloud link with the `--import=` option
- 🔢 Type system and type inference
- 🔏 Signs using macOS, falls back on [HubSign](https://routinehub.co/membership) or another server that uses [scaxyz/shortcut-signing-server](https://github.com/scaxyz/shortcut-signing-server).

### Resources

- 🍒 [Cherri VSCode Extension](https://marketplace.visualstudio.com/items?itemName=electrikmilk.cherri-vscode-extension)
- 🛝 [Playground](https://playground.cherrilang.org/) - Try out Cherri on any platform, preview the result, and export signed Shortcuts
- 🖥️ [macOS IDE](https://github.com/electrikmilk/cherri-macos-app) - Defines Cherri file type, write and build Shortcuts on Mac with a GUI
- 📄 [Documentation](https://cherrilang.org/language/) - Learn Cherri or how to contribute
- 🔍 [Glyph Search](https://glyphs.cherrilang.org/) - Search glyphs you can use in Cherri!
- ❓ [FAQ](https://cherrilang.org/faq)

## Installation

You can install Cherri by downloading the latest release or via the Homebrew package manager:

**Add Tap:**

```console
brew tap electrikmilk/cherri
```

**Install:**

```console
brew install electrikmilk/cherri/cherri
```

## Usage

```bash
cherri file.cherri
```

Run `cherri` without any arguments to see all options and usage. For development, use the `--debug` (or `-d`) option to print
stack traces, debug information, and output a `.plist` file.

## Why another Shortcuts language?

Because it's fun :)

Some languages have been abandoned, don't work well, or no longer work. I don't want Shortcuts languages to die.
There should be more, not less.

Plus, some stability comes with this project being on macOS and not iOS, and I'm not aware of another Shortcuts language with macOS as its platform other than [Buttermilk](https://github.com/zachary7829/Buttermilk).

## Community

- [VS Code Syntax Highlighting](https://marketplace.visualstudio.com/items?itemName=electrikmilk.cherri-vscode-extension) ([repo](https://github.com/electrikmilk/cherri-vscode))
- [Zed Editor](https://github.com/videah/zed-cherri)

## Credits

### Reference

- [zachary7829/Shortcuts File Format Documentation](https://zachary7829.github.io/blog/shortcuts/fileformat)
- [sebj/iOS-Shortcuts-Reference](https://github.com/sebj/iOS-Shortcuts-Reference)
- [[Tip] Reducing memory usage of repeat loops](https://www.reddit.com/r/shortcuts/comments/taceg7/tip_reducing_memory_usage_of_repeat_loops/)

### Inspiration

- Go syntax
- Ruby syntax
- [ScPL](https://github.com/pfgithub/scpl)
- [Buttermilk](https://github.com/zachary7829/Buttermilk)
- [Jelly](https://jellycuts.com)

---

_The original Workflow app assigned a code name to each release. Cherri is named after the second-to-last
update "Cherries" (also cherry is one of my favorite flavors)._

This project started on October 5, 2022.
