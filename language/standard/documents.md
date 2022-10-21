[Back](../actions.md)

# Documents Actions

## Archives

### `makeArchive(name,format,files)`

Create an archive of `files` named `name` in `format`.

---

### `extractArchive(file)`

Extract files from archive `file`.

## Previewing

### `quicklook(input)`

Preview `input`.

---

### `show(input)`

Show `input` in a dialog.

## Text

### `getText(input)`

Get text from `input`.

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

## Translation

### `translate(text,to)`

Translate `text` from the detected language of `text` to `to`.

---

### `translateFrom(text,from,to)`

Translate `text` from `from` to `to`.

---

### `detectLanguage(text)`

Detect the language of `text`.