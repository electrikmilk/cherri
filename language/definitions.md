<style>
.color {
    display: inline-block;
    height: 15px;
    width: 15px;
    border-radius: 50%;
    margin-bottom: -3.5px
}
</style>

[Back](/cherri/language/)

# Definitions

Define aspects of your Shortcut, such as the color and glyph of the icon, how it should respond
to no input, what it accepts as input, etc.

## Icon

Define the look of your Shortcut using one of the supported colors or glyphs.

```cherri
#define color red
#define glyph apple
```

Most glyphs are supported, but not all the newest are yet. [List of supported glyphs](glyphs.md).

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

These values should be separated by commas.

Inputs will default to all types. Outputs will default to no types. This is done to be consistent with the Shortcuts file format.

```cherri
#define inputs image, text
#define outputs app, file
```

## NoInput

Define how your shortcut should respond to no input.

Stop and give a specific response:

```cherri
#define noinput stopwith "Response"
```

Get the contents of the clipboard:

```cherri
#define noinput getclipboard
```

Ask for a [content item type](content-item-types.md):

```cherri
#define noinput askfor text
```

## From (Workflows)

This defines where your Shortcut should show up, `quickactions`, `sleepmode`, etc.

These values should be separated by commas.

```cherri
#define from menubar, sleepmode, onscreen
```

### Workflows

- menubar - Menubar
- quickactions - Quick Actions
- sharesheet - Share Sheet
- notifications - Notifications Center Widget
- sleepmode - Sleep Mode
- watch - Apple Watch
- onscreen - Receive On-Screen Content

## Name

This one is not widely supported. Defines the name of your Shortcut.

```cherri
#define name Test
```

## Version

Defines the version of Shortcuts your Shortcut supports.

```cherri
#define version 16
```