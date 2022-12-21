---
title: Variables & Globals
layout: default
parent: Documentation
nav_order: 2
---

# Variables & Globals

Variables are initialized using the following syntax:

```cherri
@text = "value"
```

Insert variables into a string value:

```cherri
@inserted = "Value: {text}"

// Get as...
@inserted = "Value: {text[Get]}"

// Type coersion...
@inserted = "Value: {text(type)}"
```

### Variable as a value

```cherri
@inserted = variable

// Get as...
@inserted = variable[Get]

// Type coersion...
@inserted = variable.type
```

### Globals

All globals are implemented. Globals are case-sensitive.

```cherri
@input = ShortcutInput
@date = CurrentDate
@clipboard = Clipboard
@device = Device

alert(Ask, "", false)

// But you can also just inline a global in a string like other variables
@shortcutInput = "{ShortcutInput}"
```

### Text

```cherri
@text = "text"
@multi = "multi
line
text"
```

### Numbers

```cherri
@number = 42
@expression = 54 * 4 + (6 * 7)
```

### Action Variables

```cherri
@urls = url("https://apple.com","https://google.com")
@list = list("Item 1","Item 2","Item 3")
@email = emailAddress("test@test.org")
@phone = phoneNumber("(555) 555-5555")
@date = date("October 5, 2022")
@location = getCurrentLocation()
```

### Dictionaries

You can declare a dictionary using a valid JSON object.

```cherri
@dictionary = {
    "key1": "value",
    "key2": 5,
    "key3": true,
    "key4": [
        "item1",
        "item 2",
        "item3"
    ]
}
```

### Booleans

Booleans translate to a number value of 1 for true, and 0 for false. Boolean variables can be used in conditionals.

```cherri
@boolVarTrue = true
@boolVarFalse = false

if boolVarTrue == true {
    // ...
}
```

### Misc

You can declare a variable without a value:

```cherri
@emptyVar
```

Add to a variable using the standard `+=` syntax:

```cherri
@stringVar += "test"
```
