---
title: Definitions
layout: default
parent: Documentation
nav_order: 8
---

# Definitions

Define aspects of your Shortcut, such as the color and glyph of the icon, how it responds
to no input, what it accepts as input, etc.

## Icon

Define the look of your Shortcut using one of the supported colors or glyphs.

```bash
#define color red
#define glyph apple
```

Most glyphs are supported, but not all the newest are yet. [Click here for the full list of supported glyphs](/language/glyphs.html)
.

### Color

- <span class="color" style="background-color: #FF4351"></span> `red`
- <span class="color" style="background-color: #FD6631"></span> `darkOrange`
- <span class="color" style="background-color: #FE9949"></span> `orange`
- <span class="color" style="background-color: #FEC418"></span> `yellow`
- <span class="color" style="background-color: #FFD426"></span> `green`
- <span class="color" style="background-color: #19BD03"></span> `teal`
- <span class="color" style="background-color: #55DAE1"></span> `lightBlue`
- <span class="color" style="background-color: #1B9AF7"></span> `blue`
- <span class="color" style="background-color: #3871DE"></span> `darkBlue`
- <span class="color" style="background-color: #7B72E9"></span> `violet`
- <span class="color" style="background-color: #DB49D8"></span> `purple`
- <span class="color" style="background-color: #ED4694"></span> `pink`
- <span class="color" style="background-color: #B4B2A9"></span> `taupe`
- <span class="color" style="background-color: #A9A9A9"></span> `gray`
- <span class="color" style="background-color: #555555"></span> `darkGray`

## Inputs & Outputs

Inputs and outputs accept [content item types](content-item-types.md).

Inputs will default to all types. Outputs will default to no types. This is done to be consistent with the Shortcuts
file format.

```bash
#define inputs image, text
#define outputs app, file
```

These values must be separated by commas.

## NoInput

Define how your shortcut responds to no input.

Stop and give a specific response:

```bash
#define noinput stopwith "Response"
```

Get the contents of the clipboard:

```bash
#define noinput getclipboard
```

Ask for a [content item type](/language/content-item-types.html):

```bash
#define noinput askfor text
```

## From (Workflows)

This defines where your Shortcut shows up, `quickactions`, `sleepmode`, etc.

```bash
#define from menubar, sleepmode, onscreen
```

These values must be separated by commas.

### Workflows

- `menubar` - Menubar
- `quickactions` - Quick Actions
- `sharesheet` - Share Sheet
- `notifications` - Notifications Center Widget
- `sleepmode` - Sleep Mode
- `watch` - Apple Watch
- `onscreen` - Receive On-Screen Content

## Name

This definition is not widely supported. Defines the name of your Shortcut alternative to the file name.

```bash
#define name Test
```

## Version

Defines the minimum version of iOS your Shortcut supports. Warnings will be printed if you use actions that are not supported in the targeted version.

```bash
#define version 16.2
```
