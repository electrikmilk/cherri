---
title: Documents
layout: default
grand_parent: Documentation
parent: Actions
nav_order: 3
---

# Documents Actions

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Archives

### `extractArchive(file)`

Extract files from archive `file`.

---

### `makeArchive(name,format,files)`

Create an archive of `files` named `name` in `format`.

## Books

### `addToBooks(input)`

Add `input` to books. `input` is expected to be a PDF.

## Editing

### `markup(input)`

Opens `input` in a markup editor.

## File Storage

### `createFolder(path)`

Create a folder at `path` in the users Shortcuts folder in their iCloud Drive.

---

### `deleteFiles(files,immediately)`

Delete `files`. `immediately` is a boolean, default is `false`.

---

### `getFolderContents(folder,recursive)`

Get contents of `folder`. `recursive` is a boolean, default is `true`.

---

### `getFile(path,errorIfNotFound)`

Get `path` in the users Shortcuts folder. `errorIfNotFound` is a boolean, default is `true`.

---

### `getFileFromFolder(folder,path,errorIfNotFound)`

Get `path` in `folder`. `errorIfNotFound` is a boolean, default is `true`.

---

### `getFileLink(file)`

Get a link to `file`.

---

### `getSelectedFiles()`

Get selected files in Finder.

---

### `rename(file,newName)`

Rename `file` to `newName`.

---

### `reveal(files)`

Reveal `files` in Finder.

---

### `saveFile(file,path,overwrite)`

Save `file` to `path`. `overwrite` is a boolean, default is `false`.

---

### `saveFilePrompt(file,overwrite)`

Prompt to save `file`. `overwrite` is a boolean, default is `false`.

---

### `selectFile(multiple)`

Prompt to select a file.

`multiple` is a boolean, default is `false`.

## Files

### `getFileDetail(file,detail)`

Get `detail` of `file`.

---

### `getParentDirectory(input)`

Get parent directory of `input`.

## Network

### `connectToServer(url)`

Connect to server at `url`.

## Notes

### `appendNote(note,input)`

Append `input` to `note`.

---

### `showNote(note)`

Show `note`.

## Previewing

### `quicklook(input)`

Preview `input`.

---

### `show(input)`

Show `input` in a dialog.

## Printing

### `print(input)`

Print `input`.

---

### `splitPDF(pdf)`

Split `pdf` into pages.


## QR Codes

### `makeQRCode(input,errorCorrection)`

Generate a QR code using `input` with error correction level `errorCorrection`.

#### Error Correction Levels

- Low
- Medium
- Quartile
- High

## Rich Text

### `makeHTML(input,makeFullDocument)`

Convert `richText` into HTML.

`makeFullDocument` is a boolean, default is `false`.

---

### `makeMarkdown(richText)`

Convert `richText` into Markdown.

---

### `getRichTextFromHTML(html)`

Get rich text from `html`.

---

### `getRichTextFromMarkdown(markdown)`

Get rich text from `markdown`.

## Text

### `getTextFromImage(image)`

Get text from `image`.

---

### `getEmojiName(emoji)`

Get the emoji name of `emoji`.

---

### `getText(input)`

Get text from `input`.

---

### `define(word)`

Show the definition of `word`.

---

### `matchedTextGroupIndex(matches,index)`

Get group at `index` in `matches`.

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
