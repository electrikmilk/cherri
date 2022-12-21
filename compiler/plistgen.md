---
title: Property List Generation
layout: default
parent: Contributing
nav_order: 4
---

# Property List (.plist) Generation

Property list generation is done using various functions and custom data types.

The file [`plistgen.go`](https://github.com/electrikmilk/cherri/blob/main/plistgen.go) generates the file,
while [`plist.go`](https://github.com/electrikmilk/cherri/blob/main/plist.go) contains definitions of the data types and
helper functions used in generating the plist.

A `plistData` value consists of a key, type and value. This will almost directly translate to an XML string.

```go
type plistData struct {
    key      string
    dataType plistDataType
    value    any
}
```

### `plistDataType`

Here are the constants for `plistDataType`:

- `Text`
- `Number`
- `Dictionary`
- `Array`
- `Boolean`

It is made easier by the fact that you can have a slice of `plistData` as the `value` if the `dataType` is `Dictionary`.
This is because the function that generates the plist syntax for dictionaries will recurse. This allows you to in some
cases almost completely forget about the resulting plist and mainly construct it via abstraction.

At the end of the day, these functions and types come together to build a string that is saved as a `.shortcut` file.

To save the output as a separate plist file from the resulting signed Shortcut, use the `--debug` (or `-d`) option when
running the compiler.
