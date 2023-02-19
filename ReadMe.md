<img src="https://github.com/electrikmilk/cherri/blob/main/assets/cherri_icon.png?raw=true" alt="cherri"/>

# Cherri

[![Build](https://github.com/electrikmilk/cherri/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/electrikmilk/cherri/actions/workflows/go.yml) [![Releases](https://img.shields.io/github/v/release/electrikmilk/cherri?include_prereleases)](https://github.com/electrikmilk/cherri/releases) [![Go](https://img.shields.io/github/go-mod/go-version/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/go.mod) [![License](https://img.shields.io/github/license/electrikmilk/cherri)](https://github.com/electrikmilk/cherri/blob/main/LICENSE) ![Platform](https://img.shields.io/badge/platform-macOS-red)
<a href="https://pkg.go.dev/github.com/electrikmilk/cherri?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
<a href="https://goreportcard.com/report/github.com/electrikmilk/cherri"><img src="https://goreportcard.com/badge/github.com/electrikmilk/cherri"/></a>

**Cherri** (pronounced cherry) is a [iOS Siri Shortcuts](https://apps.apple.com/us/app/shortcuts/id915249334)
programming language, that compiles directly to a valid runnable Shortcut.

[Playground](https://playground.cherrilang.org/) • [Documentation](https://cherrilang.org/language/) • [Code Tour](https://youtu.be/gU8TsI96uww)

[![Hello World Example](https://cherrilang.org/assets/example.png)](examples/hello-world.cherri)

This project is in the early stages of development, but it does have an [_idealistic_ roadmap](https://github.com/electrikmilk/cherri/wiki/Project-Roadmap).

### Usage

```bash
cherri file.cherri
```

Run `cherri` without any arguments to see all options and usage. For development, use the `--debug` (or `-d`) option for
stack traces and output a plist file.

## Why macOS only?

Generating valid Shortcuts is only possible on macOS. However, I am hoping to add a signing server to the [web editor](https://playground.cherrilang.org) that will turn out valid Shortcuts on any platform with a web browser.

### Development on other platforms

As it stands, I don't want someone to get confused and think Shortcuts compiled using Cherri on other platforms will run
on their Mac or iOS device. However, you can build the compiler for your platform and use the `--unsigned` (or `-u` for
short) to skip signing the compiled Shortcut, but the compiled Shortcut will not run on iOS or macOS, obviously. Also,
the compiler is primarily developed and tested on Unix-like systems.

[Read my full thoughts on this](https://github.com/electrikmilk/cherri/wiki/Why-macOS-only%3F)

## Why another Shortcuts language?

Because it's fun :)

Some languages have been abandoned, don't work very well, or no longer work. I don't want Shortcuts languages to die. There should be more, not less.

Some stability that comes with the project being on macOS and not iOS. I am not aware of any project [other than one](https://github.com/zachary7829/Buttermilk) that compiles a
   Shortcut in a way that is meant for a desktop OS.

## Credits

### Reference

- [zachary7829](https://github.com/zachary7829)'
  s [Shortcut File Format Reference](https://zachary7829.github.io/blog/shortcuts/fileformat)
- [sebj](https://github.com/sebj)'s [Shortcut File Format Reference](https://github.com/sebj/iOS-Shortcuts-Reference)

### Inspiration

- Go syntax
- Ruby syntax
- [ScPL](https://github.com/pfgithub/scpl)
- [Buttermilk](https://github.com/zachary7829/Buttermilk)
- [Jelly](https://jellycuts.com)

---

_The original Workflow app assigned a code name to each release. Cherri is named after the second to last
update "Cherries" (also cherry is one of my favorite flavors)._
