# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

You must read `CONTRIBUTING.md` at the start of every session.

## What This Is

Language documentation: check `~/cherrilang.org` for a local copy first, otherwise fall back to https://cherrilang.org/language/

Cherri is a compiled DSL that targets Apple Shortcuts plist XML. A `.cherri` source file compiles to a signed `.shortcut` file that Shortcuts can run directly. The compiler is a single Go binary.

## Commands

```bash
# Build
go build ./...

# Run compiler on a file — filename MUST be the first argument
go run . file.cherri
go run . file.cherri --debug      # also writes .plist and preprocessed file

# Tests
go test -run TestCherriNoSign     # compile all tests/*.cherri files, skip signing (use on non-macOS)
go test -run TestCherri           # same but with signing (macOS only)
go test -run TestDecomp           # decompile tests/decomp-me.plist, compare to tests/decomp-expected.cherri
go test -run TestActionIdentifiers # verify specific action identifier output

# Search actions and glyphs (requires build)
cherri --action=actionName
cherri --glyph=glyphName
```

The `--debug` flag is essential during development: it prints parse state, outputs a `.plist` for inspection, and writes a `_processed.cherri` showing the source after pre-processing (includes, function header injection).

**Debugging caveat:** To test a file individually with `--debug` and inspect its plist, you must temporarily add any required action `#include` statements directly to that test file. The standard library includes are only injected automatically during the full test suite run, not when compiling a single file directly.

Code style: use `var x = ...` (not `:=`) except in `for`/`if`/`switch` init statements. Types are declared without values (`var x string`). Comments explain *why*, not *what*.

## Language Syntax Basics

> This is a brief reference for contributors. The [official docs](https://cherrilang.org/language/) (or `~/cherrilang.org/language/` locally) are the authoritative source.

### Comments

```ruby
// single-line
/* multi-line */
comment('explicit comment action — always included in the Shortcut')
```

### Definitions

```ruby
#define name My Shortcut
#define color blue
#define glyph apple
#define inputs image, text
#define outputs file
#define noinput stopwith "No input provided"
#define from menubar, sharesheet
#define mac true
#define version 18.4
```

### Includes

```ruby
#include 'actions/scripting'   // standard library category
#include 'path/to/file.cherri' // arbitrary file
```

Includes are resolved in `preParse()` before any tokenization.

### Variables, Constants, Globals

```ruby
// Mutable variable — compiles to a "Set Variable" action
@name = "value"
@count = 0
@count += 1    // also -=, *=, /=

// Constant (magic variable) — references the action output directly, smaller Shortcut
const result = someAction()

// Type declaration with no initial value
@builder: text
@items: array

// Globals (case-sensitive)
@input = ShortcutInput
@now   = CurrentDate
@clip  = Clipboard
@dev   = Device

// Ask the user at runtime
wait(Ask)
wait(Ask: 'How many seconds?')
```

Mutable variables require `@` prefix when referenced inside inline strings: `"{@name}"`.

### Types

| Syntax | Type | Notes |
|---|---|---|
| `"hello {@var}"` | text | interpolates `{@var}` and escape sequences |
| `'raw text'` | rawtext | no interpolation; not allowed in dicts/arrays |
| `42` | number | |
| `0.5` | float | |
| `true` / `false` | bool | compiles to `1`/`0` |
| `{"k": "v"}` | dictionary | valid JSON syntax |
| `["a", "b"]` | array | valid JSON syntax |
| `nil` | empty | skips optional arguments; faster than `""`, `[]`, `{}` |

**Type declaration** (no initial value): `@x: text`, `@x: number`, `@x: array`, `@x: dictionary`, `@x: bool`, `@x: float`, `@x: variable`

**Type coercion**: `@var.number`, `@var.text`, `"{@var.number}"`, or coercion actions (`number()`, `getDictionary()`)

**Expressions** (arithmetic): `@result = 5 + (2 * @n)` — two-operand expressions compile to a Math action.

**Enumerations**:
```ruby
enum Color { 'Red', 'Green', 'Blue' }
```

### Control Flow

```ruby
// If / else
if @x > 0 {
    // ...
} else {
    // ...
}

// Operators: == != > >= < <= contains !contains beginsWith endsWith <> (between)
// Logical: && (all) || (any)
// Has value: if @x { } / if !@x { }
// Between: if @x <> 5 10 { }

// Repeat N times (i is the index variable)
repeat i for 6 {
    @items += "Item {@i}"
}

// For-each
@items = ["a","b","c"]
for item in @items {
    alert("@item")
}

// Control flow output (assign result to a constant)
const result = if @x == "iPhone" {
    getCellularDetail("Carrier Name")
} else {
    getWifiDetail("Network Name")
}
```

### Functions

```ruby
// Definition
function add(number op1, number op2): number {
    const s = @op1 + @op2
    output("{s}")
}

// Call
const sum = add(2, 3)

// Argument modifiers:
//   text? message   — optional
//   text! message   — literal value required (no variable)
//   text message = "default"  — default value
```

Functions compile to a `runSelf` call with a dictionary; a generated header at the top of the Shortcut dispatches to the correct body. See **Functions Abstraction** below.

### Actions

```ruby
alert("Hello!")                  // no include needed for basic actions
#include 'actions/network'
@data = getContentsOfURL("https://example.com")
```

Use `cherri --action=actionName` to look up a specific action's signature.

## Compilation Pipeline

```
.cherri source
    ↓ handleFile()       — read file, split into lines/chars
    ↓ initParse()        — initialize global state
    ↓ preParse()         — sequential pre-processing passes:
         handleIncludes()      — inline #include files, resetParse()
         handleFunctions()     — collect function defs, inject header, resetParse()
         handleActionDefinitions() — collect user-defined action blocks
    ↓ parse()            — main pass: walk chars, emit tokens[]
    ↓ generateShortcut() — walk tokens[], build shortcut.WFWorkflowActions[]
    ↓ createShortcut()   — serialize to plist, sign, write .shortcut
```

**Critical:** The compiler uses heavy global state (see `resetParser()` in `cherri_test.go` for the full list). `resetParse()` (lowercase) rebuilds `contents`/`chars`/`lines` after any source modification. `resetParser()` (in test file) resets all global state between test runs.

## Key Files

| File | Purpose |
|---|---|
| `token.go` | `tokenType` (string constants) and `token` struct |
| `shortcut.go` | All plist structs: `Shortcut`, `ShortcutAction`, `Value`, `WFTextTokenAttachment`, `WFTextTokenString`, `WFDictionaryFieldValue`, `WFArrayValue`, `WFBoolValue`, etc. |
| `action.go` | `actionDefinition`, `actionArgument`, `parameterDefinition`; action validation, parameter construction |
| `actions_std.go` | The large `actions` map of all built-in Go-defined actions |
| `shortcutgen.go` | Plist generation: `paramValue`, `variableValue`, `attachmentValues`, `makeDictionaryItem` |
| `variables.go` | `varValue` struct, variable/global storage and lookup |
| `functions.go` | Function definitions, header generation, `makeFunctionCall` |
| `parser.go` | Source parsing cursor, token emission, `initParse`/`preParse`/`parse` |
| `decompile.go` | Reverse: `.plist` → `.cherri` source |
| `actions/` | Standard library actions written in Cherri (embedded via `go:embed`) |
| `stdlib.cherri` | Standard library Cherri code (also embedded) |

## Type System

`tokenType` is `type tokenType string`. Key types:

| Cherri keyword | tokenType constant | Notes |
|---|---|---|
| `text` | `String` | |
| `number` | `Integer` | |
| `float` | `Float` | |
| `bool` | `Bool` | |
| `dictionary` | `Dict` | |
| `array` | `Arr` | |
| `variable` | `Variable` | pass-through, no coercion |
| `rawtext` | `RawString` | no interpolation |

## Plist Value Generation

The central value-generation path is in `shortcutgen.go`:

- **`paramValue(arg actionArgument, handleAs tokenType) any`** — canonical router for converting a collected argument to its plist representation. All action parameter value generation flows through here.
- **`variableValue(v varValue) any`** → **`variableValueWithSerialization`** — produces `WFTextTokenAttachment` for variable references, handling aggrandizements (`.getAs`/`.coerce`), globals, constants vs. mutable variables.
- **`attachmentValues(str string) any`** — processes strings containing `{varName}` into `WFTextTokenString` with `attachmentsByRange`.
- **`makeDictionaryItem(key, value) WFDictionaryFieldValueItem`** — converts a Go value into a typed plist dict item using a type-switch. Returns typed structs (`WFTextTokenString`, `WFArrayValue`, `WFBoolValue`, `WFDictionaryFieldValue`) rather than `map[string]any`.
- **`argumentValue(args, idx) any`** — wraps `paramValue`, using the current action's parameter definition to determine the expected type.

For function call argument serialization specifically, use `makeFunctionArgValue(arg, paramType)` in `functions.go` — this produces values suitable for the function call dict (`{ref}` strings for variables, which `makeDictionaryItem` later resolves to proper attachments).

## Action Definitions

Actions are defined two ways:

1. **Go (`actions_std.go`)**: entries in the `actions map[string]*actionDefinition`. The `makeParams` field fully overrides parameter construction; `appendParams`/`appendParamsFunc` adds extra parameters without overriding automatic handling.

2. **Cherri DSL (`actions/*.cherri`)**: user-defined action blocks parsed at compile time into the same `actions` map. These are embedded and loaded via `loadStandardActions()`.

`actionDefinition.parameters []parameterDefinition` drives automatic parameter construction in `makeActionParams` and type-checking in `checkArg`/`typeCheck`.

## Functions Abstraction

Functions are compiled to a Cherri-language header injected at the top of the Shortcut. Calling `myFunc(arg1, arg2)` generates:
1. A dictionary variable: `{cherri_functions: 1, function: "myFunc", arguments: [arg1, arg2]}`
2. A `runSelf(dict)` action

The header intercepts `ShortcutInput`, validates it as a Cherri function call, matches the function name, coerces each argument back to its declared type, and runs the function body. Type coercion per parameter type:
- `number`/`float`/`bool` → `number(argRef)`
- `text` → `"{argRef}"`
- `dictionary` → `getDictionary(argRef)`
- `array` → `getDictionary` + `getValue(dict,"array")` + `for` loop
- `variable` → direct assignment (no coercion)

## Testing

Each file in `tests/*.cherri` (except `decomp-expected.cherri` and `decomp-me.cherri`) is compiled by `TestCherri`/`TestCherriNoSign`. When adding or changing a language feature, add or update a test file in `tests/`. Group tests by domain (e.g. `web.cherri`, `shortcuts.cherri`) rather than one file per action. Test files do not need explicit `#include` statements — all standard library action categories are auto-injected during suite runs; includes are only needed when running a file individually with `go run .`.

The test runner calls `compile()` which calls `main()`, so `os.Args[1]` is set to the test file path before each run. Global state is fully reset via `resetParser()` after each test.

`TestDecomp` is a diff test: it decompiles `tests/decomp-me.plist` and expects byte-for-byte equality with `tests/decomp-expected.cherri`. When decompiler output changes intentionally, regenerate the expected file.

**Verification hierarchy:** On macOS, always run `go test -run TestCherri` as the final automated check — it is never sufficient to stop at `TestCherriNoSign`. `TestCherriNoSign` only confirms compilation does not panic; it does not validate plist structure. The signing step in `TestCherri` is the only automated confirmation that the generated plist is structurally accepted by the Shortcuts runtime. Use `TestCherriNoSign` only for rapid iteration or on non-macOS hosts. After `TestCherri` passes, the user performs the final verification: open the compiled `.shortcut` in QuickLook and import it into the Shortcuts app to confirm it runs correctly — this step cannot be automated.

**Sequential test isolation:** The test functions are not designed to run sequentially in the same process. The global `actions` map and related state accumulate across test functions, so running `go test` (all tests together) may produce failures that do not occur in CI. Always run tests individually with `-run`, matching how the GitLab pipeline executes them.

**Dual-purpose test files:** Test files are designed to both compile clean in CI *and* run in the Shortcuts app to verify runtime behavior. Each assertable test file follows this pattern:

```ruby
// run assertion inline — no helper function needed
const result = someAction("input")
if result != "expected" {
    alert("❌ FAIL: someAction — got {result}, expected 'expected'")
}

show("✅ All tests passed")  // only executes if no alert() fired
```

Files that cannot be behaviorally asserted (interactive UI, hardware-dependent, non-deterministic) begin with `// compile-only: <reason>` and still end with `show("✅ All tests passed")`.

**Comparison type constraints:** The `!=` / `==` operators (and other conditionals) require the left-side value to have a declared type in `{text, number, bool, action, date}`. A `const` bound to an action with no declared return type gets type `''` and cannot be used directly in a comparison — add `: type` to the action definition, or assign to a mutable variable via interpolation: `@s = "{const}"`. Constants and globals *can* be used directly on the left side when their type is known.

**Conditional left-side must be a variable or const, not a literal:** `if 5 == 5` causes a compile panic — always put a variable or const on the left: `@n = 5; if @n == 5`.

**`contains` only works for text and arrays:** The `contains` / `!contains` operators are not valid for dictionaries. To check if a key exists in a dict, use `getValue(dict, key)` and check `if !result`.

**Control flow output blocks require action calls, not bare literals:** Inside `const result = if ... { }` blocks, every branch must call an action. Use `text("literal value")` (the `gettext` action aliased in `actions/basic.cherri`) to output a plain string: `const r = if @x > 3 { text("yes") } else { text("no") }`. The result const has type `''` regardless — assign to a mutable variable via interpolation before comparing: `@s = "{r}"; if @s != "yes"`.

## Feature Verification Against Shortcuts Plist

A signed Shortcut is necessary but not sufficient to confirm correct behavior. The authoritative source of truth for how any action or parameter should be structured is the plist XML produced by the Shortcuts app itself.

**Workflow:** Build the equivalent Shortcut in the Shortcuts app, share it via iCloud, then retrieve its canonical plist to compare against Cherri's output. Use `--debug` (`-d`) to make Cherri emit its own plist for the comparison, and `--output=` (`-o=`) to control where the compiled `.shortcut` file is written. Run `go run . --help` for the full list of CLI flags.

**Fetching a Shortcuts iCloud plist:**
```bash
# Given a share URL like https://www.icloud.com/shortcuts/{identifier}
# Replace /shortcuts/ with /shortcuts/api/records/ to get the metadata JSON
curl "https://www.icloud.com/shortcuts/api/records/{identifier}" | jq '.fields.shortcut.value.downloadURL'

# Download the binary plist from the returned URL
curl -L "{downloadURL}" -o reference.shortcut

# Convert binary plist to readable XML (macOS)
plutil -convert xml1 reference.shortcut -o reference.plist
```

When implementing or debugging a feature, obtain the reference plist for that action type and align Cherri's generated plist structure to match. Any structural difference is a bug — the Shortcuts app's output defines correct behavior, not assumptions or prior output.

## Action Definition Workflow

The full loop for adding a new action:

1. Build the action in the Shortcuts app.
2. Share via iCloud, download, and convert to XML (see above).
3. Read the plist to extract the identifier and parameter keys.
4. Write the definition — DSL for simple cases, Go for complex ones.
5. Write a test `.cherri` file in `tests/`, compile with `--debug`, diff the output plist against the reference plist.

### Reading a plist to write a definition

A Shortcuts action in plist looks like:

```xml
<dict>
    <key>WFWorkflowActionIdentifier</key>
    <string>is.workflow.actions.gettext</string>
    <key>WFWorkflowActionParameters</key>
    <dict>
        <key>WFTextActionText</key>
        <string>hello</string>
    </dict>
</dict>
```

Map fields to the definition:

| Plist field | Definition field | Rule |
|---|---|---|
| `is.workflow.actions.foo` | `identifier: "foo"` | Strip the standard prefix; it is auto-prepended |
| `com.apple.App.Bar` | `appIdentifier: "com.apple.App"` + `identifier: "Bar"` | Non-standard prefix replaces `is.workflow.actions` |
| Fully custom identifier | `overrideIdentifier: "..."` | Used verbatim, no prefix logic applies |
| Each key in `WFWorkflowActionParameters` | `key` in `parameterDefinition` | The Cherri `name` can differ — `key` is what goes in the plist |
| A parameter always present with a fixed value | `appendParams` map | Static keys go here, not as user-facing parameters |

Parameters that vary by argument value but can't be expressed as simple pass-throughs need `makeParams` (Go only).

### DSL definition (simple actions)

Use the DSL in `actions/*.cherri` when every parameter is a direct value pass-through:

```cherri
// [Doc]: Get Text: Returns `text`.
action 'gettext' getText(text text: 'WFTextActionText'): text
```

Syntax breakdown:

- `'gettext'` — short identifier (appended after `is.workflow.actions.`); omit if the Cherri name lowercased already matches
- `getText` — the Cherri call name; this becomes the key in the `actions` map
- `text text: 'WFTextActionText'` — `type name: 'plistKey'`; omit `: 'plistKey'` when the name and plist key are identical
- `: text` — output type; omit when the action produces no output

Static parameters go in an inline dict block:

```cherri
action 'output' outputOrClipboard(text output: 'WFOutput') {
    "WFNoOutputSurfaceBehavior": "Copy to Clipboard"
}
```

### Go definition (complex actions)

Use Go in `actions_std.go` when you need any of:

- Custom parameter construction (`makeParams`) — the function receives `[]actionArgument` and returns `map[string]any`; it **fully replaces** automatic handling
- Argument validation beyond type-checking (`check`)
- Dynamic extra parameters without disabling automatic handling (`appendParamsFunc`) — use this to inject derived plist keys (e.g. boolean enable-flags like `WFXCallbackCustomCallbackEnabled`) computed from argument values. Every `args[i]` access inside this func must be guarded by a `len(args)` check that covers index `i` (e.g. `if len(args) >= 5 { args[4]... }`).
- Decompiler support (`decomp`)

Minimal Go definition:

```go
"getText": {
    parameters: []parameterDefinition{
        {name: "text", validType: String, key: "WFTextActionText"},
    },
    outputType: String,
},
```

Use `argumentValue(args, i)` or `paramValue(arg, type)` inside `makeParams` to convert collected arguments to plist values. Use `appendParamsFunc` instead of `makeParams` when you want automatic parameter handling to still apply but need to inject extra keys based on argument values.
