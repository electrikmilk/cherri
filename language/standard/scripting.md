---
title: Scripting
layout: default
grand_parent: Documentation
parent: Actions
nav_order: 6
---

# Scripting Actions

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Apps

### `openApp(appId)`

Open app with `appId`. For example `com.apple.AppStore`.

_Contributed by [JosephShenton](https://github.com/JosephShenton)_.

#### App ID Shorthands

- `appstore`
- `files`
- `shortcuts`
- `safari`
- `facetime`
- `notes`
- `phone`
- `reminders`
- `mail`
- `music`
- `calendar`
- `maps`
- `contacts`
- `health`
- `photos`

---

### `hideApp(appId)`

Hide app with `appId`. For example `com.apple.AppStore`.

---

### `quitApp(appId)`

Quit app with `appId`. For example `com.apple.DocumentsApp`.

---

### `killApp(appId)`

Quit app with `appId` without asking to save changes. For example `com.apple.facetime`.

## Content

### `getOnScreenContent()`

Get content on device screen.

## Control Flow

Some of these actions are abstracted into statements:

- [Choose From Menu](../menus.md#1-choose-from-menu)
- [If](../statements.md#ifelse)
- [Repeat](../statements.md#repeat)
- [Repeat With Each](../statements.md#repeat-with-each)

---

### `output(output)`

Stop and output `output`. Do nothing if there is nowhere to output.

---

### `mustOutput(output,response)`

Stop and output `output`. Respond with `response` if there is nowhere to output.

---

### `outputOrClipboard(output)`

Stop and output `output`. Copy to the clipboard if there is nowhere to output.

---

### `stop()`

Stop the Shortcut. Equivalent to `exit` or `die` in other languages.

---

### `wait(number)`

Wait `number` of seconds.

---

### `waitToReturn()`

Wait for the user to return to the Shortcut.

## Device

### `getBatteryLevel()`

Get the current battery level of device.

---

### `isCharging()`

Check if the device is charging.

_**Minimum iOS version:** 16.2_

---

### `connectedToCharger()`

Check if device is connected to charger.

_**Minimum iOS version:** 16.2_

---

### `getDeviceDetail(detail)`

Get `detail` of current device.

---

### `toggleApperance()`

Toggles system appearance from light to dark, or dark to light.

---

### `lightMode()`

Change system appearance to light.

---

### `darkMode()`

Change system appearance to dark.

---

### `toggleBluetooth()`

Toggle device bluetooth.

---

### `setBluetooth(status)`

Set device bluetooth on or off.

---

### `setBrightness(number)`

Set display brightness to `number`.

---

### `setVolume(number)`

Set device volume to `number`.

--- 

### `startScreensaver()`

Start screen saver on Mac.

## Dictionaries

### [Dictionary](../variables-globals.md#dictionaries)

---

### `getDictionary(input)`

Get dictionary from `input`.

---

### `getValues(dictionary)`

Get values of `dictionary`.

---

### `getKeys(dictionary)`

Get keys of `dictionary`.

---

### `getValue(dictionary,key)`

Get the value of `key` in `dictionary`.

---

### `setValue(key,value,dictionary)`

Set the value of `key` to `value` in `dictionary`.

## Files

### `base64Encode(input)`

Base 64 encode `input`.

---

### `base64Decode(input)`

Base 64 decode `input`.

---

### `hash(input,type)`

Generate a hash of `type` using `input`.

#### Hash Types

- MD5 (default)
- SHA1
- SHA256
- SHA512

## Items

### `countItems(input)`

Returns the number of items in `input`.

---

### `countChars(input)`

Returns the number of characters in `input`.

---

### `countWords(input)`

Returns the number of words in `input`.

---

### `countSentences(input)`

Returns the number of sentences in `input`.

---

### `countLines(input)`

Returns the number of lines in `input`.

---

### `getName(input)`

Get the name of `input`.

---

### `getType(input)`

Get the type of `input`.

---

### `setName(input,name,includeFileExtension)`

Set the name of `input`. `includeFileExtension` is a boolean.

## Lists

### `chooseFromList(list)`

Prompt the user to choose an item from `list`. Returns the item chosen.

---

### `chooseMultipleFromList(list)`

Prompt the user to choose multiple items from `list`. Returns all of the items chosen.

---

### `firstListItem(list)`

Get the first item from `list`.

---

### `lastListItem(list)`

Get the last item from `list`.

---

### `randomListItem(list)`

Get a random item from `list`.

---

### `getListItem(list,index)`

Get item at `index` from `list`. Counting in Shortcuts starts at `1`.

---

### `getListItems(list,start,end)`

Get items in range of `start` to `end`.

---

### `list(...item)`

Store a list of `item` in a variable. No limit on `item` arguments.

## Math

### Calculate Expression

To do this you make a variable and set the value as an expression:

```cherri
@expression = 25 * 6 + (5 / 6)
```

---

### Rounding

- `round(number,roundTo)` - Normal
- `roundUp(number,roundTo)` - Always Round Up
- `roundDown(number,roundTo)` - Always Round Down

Round `number` to `roundTo`.

Shorthands for `roundTo` value:

- `1` - Ones Place
- `10` - Tens Place
- `100` - Hundreds Place
- `1000` - Thousands
- `10000` - Ten Thousands
- `100000` - Hundred Thousands
- `1000000` - Millions

## Network

### `isOnline()`

This is an alias of the standard **Get IP Address** action.

### `getLocalIP(type)`

Get the local IP of the user of `type`.

`type` is optional.

### Types

- IPv4
- IPv6

---

### `getExternalIP(type)`

Get the external IP of the user of `type`.

`type` is optional.

### Types

- IPv4
- IPv6

---

### `setCellularData(bool)`

Turn cellular data to on or off.

---

### `setCellularVoice(bool)`

Turn cellular voice and data to on or off.

---

### `setWifi(bool)`

Turn cellular voice and data to on or off.

## No-ops (noonce)

### Comments

```js
// Single line comment
```

```js
/*
Multiline
comment
*/
```

---

### `nothing()`

Do nothing and output nothing.

## Notification

### `askForInput(type,prompt,default)`

Ask for input of `type` with `prompt`. `default` is optional.

#### Input Types

- Text
- Number
- URL
- Date
- Time
- Date and Time

---

### `playSound(input)`

Play sound `input`.

---

### `alert(alert,title,cancelButton)`

Show an alert with `title` as the title and `alert` as the body. `cancelButton` is a boolean, default is `true`.

---

### `notification(body,title,playSound)`

Trigger a custom notification. `playSound` is a boolean, default is `true`.

## Numbers

### `fileSize(file,format)`

Format the size of `file` to `format`.

---

### `formatNumber(number,decimalPlaces)`

Format `number` with `decimalPlaces` number of decimal places.

---

### `getNumber(input)`

Get numbers from `input`.

---

### `randomNumber(min,max)`

Generate a random number with `min` as the minimum value and `max` as the maximum value.

## Shell

### `runShellScript(script,input,shell,inputMode)`

Run `script` with `input` as `inputMode` in `shell`.

`shell` and `inputMode` are not required. Default shell is `/bin/zsh` and input mode is `to stdin`.

#### Input Modes

- `to stdin`
- `as arguments`

## Shortcuts

### `getShortcuts()`

Get all the users shortcuts.

---

### `open(name)`

Open Shortcut with name `name`.

---

### `run(name,input,isSelf)`

Run Shortcut with name `name` giving `input`. `isSelf` is a boolean, default is `false`.

## System

### `dismissSiri()`

Dismiss Siri and continue.

---

### `setWallpaper(image)`

Set device wallpaper to `image`.

---


### `getWallpaper()`

Get device wallpaper.

_**Minimum iOS version:** 16.2_

## Variables

See [Variables & Globals](../variables-globals.md).

## X-Callback

### `openXCallbackURL(url)`

Open X-Callback URL `url`.

---

### `openXCustomCallbackURL(url,successKey,cancelKey,errorKey,successURL)`

Open X-Callback URL `url`, with Success Key `successKey`, Cancel Key `cancelKey`, and Error Key `errorKey`, and custom
X-Success URL `successURL`.
