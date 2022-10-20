[Back](/cherri/language/)

# Actions

These are the standard Shortcuts actions currently built-in and more are being added all the time.

[Learn about contributing actions here](/cherri/compiler/actions.md).

Actions in Cherri are intended to be easier to use, as in some cases, single actions have been split up into multiple
actions to reduce the number of arguments.

All arguments are currently required. `""` is accepted as an empty string value.

## Scripting

### `wait(number)`

Wait `number` of seconds.

---

### `waitToReturn()`

Wait for the user to return to the Shortcut.

---

### `alert(alert,title,cancelButton)`

Show an alert with `title` as the title and `alert` as the body. `cancelButton` is a boolean.

---

### `askForInput(type,prompt,default)`

Ask for input of `type` with `prompt`. Optionally specify `default` or just `""`

---

### `notification(body,title,playSound)`

Trigger a custom notification.

---

### `chooseFromList(list)`

Prompt the user to choose an item from `list`. Returns the item chosen.

---

### `chooseMultipleFromList(list)`

Prompt the user to choose multiple items from `list`. Returns all of the items chosen.

---

### `stop()`

Stop the Shortcut. Equivalent to `exit` or `die` in other languages.

---

### `nothing()`

Used as a way to prevent the output of the previous action being used as input in the next action or as the Shortcuts
output.

---

### `setVolume(number)`

Set device volume to `number`.

---

### `setBrightness(number)`

Set display brightness to `number`.

---

### `countItems(input)`

Returns the number of items in `input`.

---

### `countChars(input)`

Returns the number of characters in `input`.

---

### `countWords(input)`

Returns the words of characters in `input`.

---

### `countSentences(input)`

Returns the number of sentences in `input`.

---

### `countLines(input)`

Returns the number of lines in `input`.

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

### `getName(input)`

Get the name of `input`.

---

### `setName(input,name,includeFileExtension)`

Set the name of `input`. `includeFileExtension` is a boolean.

---

### `getType(input)`

Get the type of `input`.

---

### `openURL(url)`

Open `url` in default browser.

---

### `quicklook(input)`

Preview `input`.

---

### `show(input)`

Show `input` in a dialog.

---

### `getBatteryLevel()`

Get the current battery level.

---

### `getCurrentLocation()`

Get the users current location.

## Translation

### `translate(text,to)`

Translate `text` from the detected language of `text` to `to`.

---

### `translateFrom(text,from,to)`

Translate `text` from `from` to `to`.

---

### `detectLanguage(text)`

Detect the language of `text`.

## Social

### `share(input)`

Share `input` via the Share Sheet.

---

### `copyToClipboard(input,local,expire)`

Copy `input` to the clipboard. `local` is a boolean. `expire` is a date as a string (e.g. Today at 3pm).

---

### `getClipboard()`

Get the current contents of the clipboard.

## Files

### `makeArchive(name,format,files)`

Create an archive of `files` named `name` in `format`.

---

### `extractArchive(file)`

Extract files from archive `file`.

## Dictionaries

### `getKeys(dictionary)`

Get the keys of `dictionary`.

---

### `getValues(dictionary)`

Get the values of `dictionary`.

---

### `getValue(dictionary,key)`

Get the value of `key` in `dictionary`.

## Values

### `url(...url)`

Store URL(s) in a variable.

---

### `list(...item)`

Store a list in a variable.

## Encoding & Hashes

### `hash(type,input)`

Generate a hash of `type` using `input`.

---

### `base64Encode(input)`

Base 64 encode `input`.

---

### `base64Decode(input)`

Base 64 decode `input`.

## Numbers

### `formatNumber(number,decimalPlaces)`

Format `number` with `decimalPlaces` number of decimal places.

---

### `randomNumber(min,max)`

Generate a random number with `min` as the minimum value and `max` as the maximum value.

## Shortcuts

### `getShortcuts()`

Get all the users shortcuts.

---

### `open(name)`

Open Shortcut with name `name`.

---

### `run(name,input,isSelf)`

Run Shortcut with name `name` giving `input`. `isSelf` is a boolean.

## Text Editing

### `splitText(text,separator)`

Split `text` by `separator`.

---

### `combineText(text,glue)`

Combine `text` with `glue`.

---

### `replaceText(find,replacement,subject)`

Find `find` in `subject` and replace with `replacement`.

---

### `iReplaceText(find,replacement,subject)`

Case-insensitive find `find` in `subject` and replace with `replacement`.

---

### `regReplaceText(expression,replacement,subject)`

Use a regular-expression to find `find` in `subject` and replace with `replacement`.

---

### `iRegReplaceText(expression,replacement,subject)`

Use a case-insensitive regular-expression to find `find` in `subject` and replace with `replacement`.

---

### `uppercase(text)`

Change case of `text` to UPPERCASE.

---

### `lowercase(text)`

Change case of `text` to lowercase.

---

### `capitalize(text)`

Capitalize the first word in `text`.

---

### `capitalizeAll(text)`

Capitalize all the words in `text`.

---

### `titleCase(text)`

Change case of `text` to Title Case.

---

### `alternateCase(text)`

Change case of `text` to aLtErNaTiNg cAsE.

---

### `correctSpelling(text)`

Correct the spelling of `text`.