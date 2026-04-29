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

**Sequential test isolation:** The test functions are not designed to run sequentially in the same process. The global `actions` map and related state accumulate across test functions, so running `go test` (all tests together) may produce failures that do not occur in CI. Always run tests individually with `-run`, matching how the GitLab pipeline executes them.

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
