# Cherri Language Guide for LLMs

## Overview
Cherri is a programming language that compiles to Apple Shortcuts. When writing Cherri code, you're essentially creating iOS/macOS Shortcuts using a more traditional programming syntax.

## Key Concepts

### 1. Basic Syntax
- **Comments**: Use `/* */` for block comments, `//` for line comments
- **Variables**: Use `@varName = value` to declare variables
- **Constants**: Use `const varName = value` for constants
- **Strings**: Can use double quotes `"text"` or template literals with variables `"{varName}"`

### 2. Actions
Actions are the core building blocks - they map directly to Shortcut actions:

```cherri
// Display an alert
alert("Message", "Title")

// Show output
show("Text to display")

// Get text input
@userInput = askText("Enter text:", "Default value")
```

### 3. Control Flow

#### If Statements
```cherri
if condition {
    // code
} else {
    // code
}
```

#### Loops
```cherri
// For each item in list (CORRECT syntax)
for item in list {
    alert(item)
}

// Repeat n times
repeat 5 {
    alert("{RepeatIndex}")
}

// Repeat with index variable
repeat i for 5 {
    alert("Index: {i}")
}
```

### 4. Menus

#### Correct Menu Syntax (IMPORTANT!)
```cherri
menu "Choose an option:" {
    item "Option 1":
        alert("You chose option 1")
    item "Option 2":
        alert("You chose option 2")
}
```

#### Alternative: List-Based Menus
```cherri
#include 'actions/scripting'

@list = list("Option 1", "Option 2", "Option 3")
@chosenItem = chooseFromList(list, "Choose an item")
```

#### Key Menu Rules:
- Use **colon `:` after item name**, NOT curly braces `{}`
- Menu prompts can be variables or string literals
- Item labels can be variables or string literals
- Each item can have multiple action lines
- Variables assigned in menu items persist after menu execution

### 5. Common Patterns

#### Variables and Text Manipulation
```cherri
@name = "World"
@greeting = "Hello, {name}!"
alert(greeting)

// Text operations
@upper = changeCase(greeting, "uppercase")
@replaced = replaceText("World", "Cherri", greeting)
```

#### Lists/Arrays
```cherri
@items = ["apple", "banana", "cherry"]
repeatEach fruit in items {
    alert(fruit)
}
```

#### Getting User Input
```cherri
@name = askText("What's your name?")
@age = askNumber("What's your age?")
alert("Hello {name}, you are {age} years old")
```

### 6. File Organization

#### Includes
Use `#include` to split code across files:
```cherri
#include "other_file.cherri"
```

#### Defines (Preprocessor)
```cherri
#define color red
#define glyph star
#define name "My Shortcut"
```

### 7. Custom Actions
Define reusable actions:
```cherri
action greet(text name) {
    alert("Hello, {name}!")
}

// Use it
greet("World")
```

## Best Practices for LLMs

### When Writing Cherri Code:

1. **Start Simple**: Begin with basic alert() or show() to test compilation
2. **Use Examples**: Reference tests/*.cherri files for syntax patterns
3. **Check Action Definitions**: Actions are defined in actions/*.cherri files
4. **Variable Syntax**: Always use @ for variable assignment, not declaration
5. **String Interpolation**: Use `"{varName}"` to include variables in strings

### Common Mistakes to Avoid:

1. **Wrong Alert Syntax**: `alert(message, title)` not `alert(title, message)`
2. **Missing @**: Variables need @ when assigned: `@var = value`
3. **Case Sensitivity**: Cherri is case-sensitive
4. **Action Names**: Many actions have specific names (e.g., `askText` not `prompt`)

### Compilation Tips:

1. **Build Command**: `./cherri filename.cherri`
2. **Fast Iteration**: Use `--skip-sign` flag during development to skip signing (much faster)
   - Development: `./cherri filename.cherri --skip-sign`
   - Production: `./cherri filename.cherri` (includes signing)
3. **Debug Mode**: Use `--debug` or `-d` flag for detailed output
4. **Output**: Creates `.shortcut` file with same base name
5. **Signing**: Will attempt macOS signing, falls back to HubSign service (unless --skip-sign is used)

### Understanding Errors:

1. **Parse Errors**: Usually syntax issues - check brackets, quotes
2. **Unknown Action**: Action doesn't exist - check actions/*.cherri
3. **Type Errors**: Wrong parameter types - check action definitions
4. **Variable Errors**: Undefined variables - ensure @ is used correctly

## Quick Reference

### Essential Actions:
- `alert(message, ?title)` - Show alert dialog
- `show(text)` - Display output
- `askText(prompt, ?default)` - Get text input
- `askNumber(prompt, ?default)` - Get number input
- `nothing()` - Clear output
- `stop()` or `exit()` - Stop shortcut
- `comment(text)` - Add comment

### Variables:
- `@var = value` - Assign variable
- `const name = value` - Declare constant
- `{varName}` - Use in strings
- Global variables: `Ask`, `Clipboard`, `CurrentDate`, `ShortcutInput`, `RepeatIndex`, `RepeatItem`

### Types:
- `text` - String type
- `number` - Numeric type
- `boolean` - True/false
- `array` - List/collection
- `dictionary` - Key-value pairs

## Testing Workflow

### Development (Fast Iteration)
1. Write `.cherri` file
2. Compile with skip-sign: `./cherri file.cherri --skip-sign`
3. Check for errors in output
4. Fix any issues and repeat
5. Test the unsigned `.shortcut` file locally

### Production (Final Version)
1. Ensure code compiles without errors
2. Compile with signing: `./cherri file.cherri`
3. Wait for signing process (macOS or HubSign)
4. Distribute the signed `.shortcut` file

## Example Programs

### Hello World
```cherri
alert("Hello World", "Greeting")
```

### User Input
```cherri
@name = askText("What's your name?")
alert("Hello, {name}!", "Personalized Greeting")
```

### Loop Example
```cherri
repeat 3 {
    alert("Count: {RepeatIndex}", "Loop Demo")
}
```

### Menu Example
```cherri
menu "Pick a color:" {
    item "Red":
        alert("You chose red!")
    item "Blue":
        alert("You chose blue!")
}
```

This guide should help any LLM understand and write Cherri code effectively. Always refer to test files and action definitions for the most accurate syntax.
