# AGENTS.md

This file provides guidance to agents when working with code in this repository.

You must read `CONTRIBUTING.md` at the start of every session.

## What This Is

Language documentation: check `~/cherrilang.org` for a local copy first, otherwise fall back to https://cherrilang.org/language/

Cherri is a compiled DSL that targets Apple Shortcuts plist XML. A `.cherri` source file compiles to a signed `.shortcut` file that Shortcuts can run directly. The compiler is a single Go binary.

## Commands

```bash
# Build
go build ./...

# Run compiler on a file
go run . file.cherri
go run . --debug file.cherri      # also writes .plist and preprocessed file

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

Code style: use `var x = ...` (not `:=`) except in `for`/`if`/`switch` init statements. Types are declared without values (`var x string`). Comments explain *why*, not *what*.

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

For function call argument serialization specifically, use `functionCallArgValue(arg, paramType)` in `functions.go` — this produces values suitable for the function call dict (strings with `{ref}` for variables, which `makeDictionaryItem` later resolves to proper attachments).

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

Each file in `tests/*.cherri` (except `decomp-expected.cherri` and `decomp_me.cherri`) is compiled by `TestCherri`/`TestCherriNoSign`. When adding or changing a language feature, add or update a test file in `tests/`.

The test runner calls `compile()` which calls `main()`, so `os.Args[1]` is set to the test file path before each run. Global state is fully reset via `resetParser()` after each test.

`TestDecomp` is a diff test: it decompiles `tests/decomp-me.plist` and expects byte-for-byte equality with `tests/decomp-expected.cherri`. When decompiler output changes intentionally, regenerate the expected file.

**Format correctness caveat:** `TestCherriNoSign` only verifies that compilation does not panic — it does not validate plist format. Shortcuts signing (`TestCherri` on macOS) will fail if the plist is structurally invalid, making it a stronger format check. Even a successful sign is not sufficient on its own: the resulting Shortcut must be manually opened and run in Shortcuts to confirm it behaves correctly. Automated tests cannot substitute for this manual verification step.
