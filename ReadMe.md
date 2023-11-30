![Cherri](https://github.com/electrikmilk/cherri/blob/main/assets/cherri_icon.png?raw=true)

# Cherri

[![Build & Test](https://github.com/electrikmilk/cherri/actions/workflows/go.yml/badge.svg)](https://github.com/electrikmilk/cherri/actions/workflows/go.yml)
[![macOS Build](https://github.com/electrikmilk/cherri/actions/workflows/go-macos.yml/badge.svg)](https://github.com/electrikmilk/cherri/actions/workflows/go-macos.yml)
[![Releases](https://img.shields.io/github/v/release/electrikmilk/cherri?include_prereleases)](https://github.com/electrikmilk/cherri/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/go.mod)
[![License](https://img.shields.io/github/license/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/LICENSE)
![Platform](https://img.shields.io/badge/platform-macOS-red)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/electrikmilk/cherri?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrikmilk/cherri)](https://goreportcard.com/report/github.com/electrikmilk/cherri)

**Cherri** (pronounced cherry) is a [iOS Siri Shortcuts](https://apps.apple.com/us/app/shortcuts/id915249334)
programming language, that compiles directly to a valid runnable Shortcut.

[![Hello World Example](https://github.com/electrikmilk/cherri/assets/4368524/4205308b-4358-47c9-a1b4-66930c8db1c2)](https://playground.cherrilang.org)

### 🌟 Top Features

- 🎓 Easy to learn and syntax similar to other languages
- 🐞 1-1 translation to Shortcut actions as much as possible to make debugging easier
- 🪄 No magic variables syntax, they're constants instead
- 🪶 Optimized to create small as possible Shortcuts and reduces memory usage at runtime
- #️⃣ Include files within each other for large Shortcut projects
- 🔧 Define your own custom actions to paste actions
- 🥩 Enter actions raw with custom identifier and parameters
- ❓ Define import questions
- 📇 Generate VCards for menus
- 🔢 Strongly statically typed with type inference

### Learn more

- 🛝 [Playground](https://playground.cherrilang.org/)
- 📄 [Documentation](https://cherrilang.org/language/)
- 🧑‍💻 [Code Tour](https://youtu.be/gU8TsI96uww)
- 🗺️ [_Idealistic_ roadmap](https://github.com/electrikmilk/cherri/wiki/Project-Roadmap)

### **📣 WIP 📣**

This project has not yet reached a stable version. It is under heavy development and backward
incompatible changes may be made.

## Usage

```bash
cherri file.cherri
```

Run `cherri` without any arguments to see all options and usage. For development, use the `--debug` (or `-d`) option to print
stack traces, debug information, and output a .plist file.

## Why macOS only?

Generating valid Shortcuts is only possible on macOS. However, I am hoping to add a signing server to
the [web editor](https://playground.cherrilang.org) that will turn out valid Shortcuts on any platform with a web
browser.

### Development on other platforms

As it stands, I don't want someone to get confused and think Shortcuts compiled using Cherri on other platforms will run
on their Mac or iOS device. However, you can build the compiler for your platform, it will just skip signing the
compiled Shortcut, so it will not run on iOS 15+ or macOS 12+. Also, note that the compiler is primarily developed and
tested on Unix-like systems.

[Read my full thoughts on this](https://github.com/electrikmilk/cherri/wiki/Why-macOS-only%3F)

## Why another Shortcuts language?

Because it's fun :)

Some languages have been abandoned, don't work very well, or no longer work. I don't want Shortcuts languages to die.
There should be more, not less.

Some stability comes with the project being on macOS and not iOS.

## Community

- [VS Code Syntax Highlighting](https://marketplace.visualstudio.com/items?itemName=erenyenigul.cherri) ([Github](https://github.com/erenyenigul/cherri-vscode-highlight))

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

_The original Workflow app assigned a code name to each release. Cherri is named after the second to last
update "Cherries" (also cherry is one of my favorite flavors)._

This project started on Oct 5, 2022.
