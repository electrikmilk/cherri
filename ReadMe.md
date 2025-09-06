<img src="https://github.com/electrikmilk/cherri/blob/main/assets/cherri_icon.png" width="200"/>

# Cherri

[![Build & Test](https://github.com/electrikmilk/cherri/actions/workflows/go.yml/badge.svg)](https://github.com/electrikmilk/cherri/actions/workflows/go.yml)
[![Releases](https://img.shields.io/github/v/release/electrikmilk/cherri?include_prereleases)](https://github.com/electrikmilk/cherri/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/go.mod)
[![License](https://img.shields.io/github/license/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/LICENSE)
![Platform](https://img.shields.io/badge/platform-macOS-red)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/electrikmilk/cherri?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrikmilk/cherri)](https://goreportcard.com/report/github.com/electrikmilk/cherri)

**Cherri** (pronounced cherry) is a [Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752) programming language that compiles directly to a valid runnable Shortcut.

The primary goal is to make it practical to create large Shortcut projects (within the limitations of Shortcuts) and maintain them long term.

[![Hello World Example](https://github.com/electrikmilk/cherri/blob/main/assets/hello_world.png)](https://playground.cherrilang.org)

### 🌟 Top Features

- 🖥️ Laptop/Desktop-based development (CLI, [VSCode extension](https://marketplace.visualstudio.com/items?itemName=electrikmilk.cherri-vscode-extension), macOS app)
- 🎓 Easy to learn and syntax similar to other languages
- 🐞 1-1 translation to Shortcut actions as much as possible to make debugging easier
- 🪄 No magic variables syntax, they're constants instead
- 🪶 Optimized to create as small as possible Shortcuts and reduce memory usage at runtime
- #️⃣ Include files within others for large Shortcut projects
- 🔧 Define custom actions with type checking, enums, optionals, default values, raw identifiers, and raw keys.
- 📋 Copy-paste actions automatically
- 🥩 Enter actions raw with a custom identifier and parameters.
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
