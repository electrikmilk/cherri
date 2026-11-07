---
name: actions-sync
description: Use when developing, debugging or verifying Cherri actions. Finds standard Shortcuts and Apple app actions Cherri is missing or has lost, new actions after an OS release, and gets the plist Shortcuts itself stores for an installed shortcut so it can be compared with the plist Cherri compiles.
---

# Actions sync

Three tools that keep Cherri's action definitions matched to what Shortcuts actually does. They are plain bash scripts in this skill's `scripts/` folder, so any agent or person can run them directly. Run them from the repo root.

| Need | Tool |
|---|---|
| Which actions exist in Shortcuts but not Cherri, or the reverse | `scripts/actions-sync.sh` |
| Which parameters and enum values differ for actions both define | `scripts/actions-sync-params.sh` |
| What plist does Shortcuts produce for this action | `scripts/shortcut-plist.sh` |

The paths in the table are relative to `.claude/skills/actions-sync/`. The commands below use the full path from the repo root.

## Find missing or lost actions

```bash
bash .claude/skills/actions-sync/scripts/actions-sync.sh
```

Compares the Shortcuts ToolKit databases (macOS and the booted iOS simulator) with the actions Cherri defines, by full Shortcuts identifier such as `is.workflow.actions.alert`. It builds Cherri from the working tree first. The brew binary can lag the repo by months, so do not compare against it. To reuse a binary you already built, pass `--cherri <path>`.

Prerequisites:

- Open the Shortcuts app on the Mac once, so the Mac database exists. Override with `--mac-db <path>`.
- Boot an iOS simulator and open Shortcuts in it once, so the simulator database exists. The script finds it with `xcrun simctl getenv booted HOME`. Override with `--sim-db <path>`.

The databases are opened read-only because Shortcuts writes to them while running.

### Scope

Standard Shortcuts actions (`is.workflow.actions.*`) and Apple app actions (`com.apple.*`). Third-party app actions are excluded unless asked for. To include them, change the `scope` in the script.

Within that, only `action`, `appIntent` and `siriIntent` tools are compared. Identifiers ending in `Entity` are skipped because they are data types, not actions. That is a naming heuristic, so check a skipped-looking result before trusting it. `visibilityFlags` is not used as a filter because flag `2` hides entities and internal tools but also real intents such as `com.apple.MobileSMS.DraftMessageIntent`.

### Reading the report

| Section | Meaning |
|---|---|
| In simulator only and missing from Cherri | Actions Shortcuts has on iOS but not on macOS, and Cherri lacks. Start here for iOS-only gaps. |
| In Mac database but missing from Cherri | Candidates to implement. |
| In simulator database but missing from Cherri | Same for iOS. |
| In Cherri but not in Mac database | Mac-missing is expected for iOS-only actions. Otherwise the action was renamed, removed or deprecated. |
| In Cherri but not in simulator database | Mac-only actions are expected. Otherwise renamed, removed or deprecated. |

Each missing entry shows the identifier and its English name from the database.

### Limits

- Identifier level only. An action present in both can still differ in parameters, so run `.claude/skills/actions-sync/scripts/actions-sync-params.sh` for those.
- Cherri needs `--action= --debug` to print identifiers. Without `--debug` the listing hides actions that cannot be defined in the Cherri DSL.
- The databases are versioned per OS release, and the Mac and the simulator can be on different releases. Mention the versions when reporting results.

## Compare parameters and enums

```bash
bash .claude/skills/actions-sync/scripts/actions-sync-params.sh
```

For every action both Shortcuts and Cherri define, it compares the plist parameter keys and the enum values. It takes the same options and prerequisites as `actions-sync.sh`, and reports separately for the Mac and the simulator database.

| Section | Meaning |
|---|---|
| Shortcuts parameters Cherri does not expose | Keys the database lists that no Cherri parameter or fixed value sets. |
| Cherri parameters Shortcuts does not list | Keys Cherri sets that the database does not list for that action. |
| Enum values that differ | Values only one side has. Check these first, since a wrong value can compile to a plist Shortcuts rejects. |
| Free text in Cherri where Shortcuts lists values | A Cherri parameter with no enum where the database has a fixed set of values. |

Treat the report as a list of leads, not defects:

- Some keys are left out on purpose, such as `ShowWhenRun`, `IntentAppDefinition`, `OpenWhenRun` and the filter keys (`WFContentItemFilter`, `WFCompoundType`). Check whether the action ever needs them before adding one.
- Actions marked "keys also built in Go" set some keys from their arguments in Go code, which the script cannot see, so an unexposed key may already be set.
- Enum values are compared exactly. A value that differs only because Cherri maps a friendlier name to the real one still shows up.
- Only enums the database stores under the tool's own parameter are compared. A parameter whose values live in a shared type, such as the Health sample types, is not matched automatically. Compare it with the shared type's `EnumerationCases` by hand.

`cherri --action=<name> --debug` prints the same keys and enum values for one action, plus the fixed values, so use it to see what a single flagged action actually defines.

## App intent actions

Two sources answer different questions about an app intent:

| Source | Answers |
|---|---|
| ToolKit database (`Tools` table) | What Shortcuts offers as actions on that OS. The Health actions in the simulator's action list match its Health rows exactly. |
| `Metadata.appintents/extract.actionsdata` | What an app or framework defines as App Intents: intent names, parameter names and types, titles. |

An app can define an intent that Shortcuts does not offer. In the iOS 27.2 simulator, `HealthIntents.framework` defines `OpenViewIntent`, `OpenDataTypeIntent`, `OpenSearchIntent` and `OpenSleepScheduleIntent`, but the database has no rows for them and the action list does not show them. So use the metadata to check an intent's name, bundle and parameters, and use the database or a reference plist from the Shortcuts UI to confirm it is a usable action.

The plist identifier of an app intent is its `attributionBundleIdentifier` plus `.` plus its `identifier`, for example `com.apple.ShortcutsActions` plus `GetPhysicalActivity`.

To find the metadata for an app in the simulator:

1. `xcrun simctl listapps booted` lists each app's `Bundle` path. The file is `<bundle>/Metadata.appintents/extract.actionsdata`.
2. If its `actions` object is empty, `extract.packagedata` lists the frameworks that define the intents under `includes`. The Swift-mangled name starts with the module name's length, so `13HealthIntents...` is `HealthIntents`. That framework's file is at `<runtime root>/System/Library/PrivateFrameworks/HealthIntents.framework/Metadata.appintents/extract.actionsdata`. The runtime root is the part of the app's bundle path before `/Applications/`.
3. `jq -r '.actions | keys[]' <file>` lists the intents, and `jq '.actions.<Name>' <file>` shows their parameters and types.

## Check the plist Cherri compiles

Shortcuts is the source of truth for how an action and its parameters are structured. Instead of sharing through iCloud and downloading, read the plist straight from the Shortcuts database.

1. Build the action in the Shortcuts app with the parameter values to cover. Shortcuts saves it to its database on its own.
2. Print the plist Shortcuts stored for it:

   ```bash
   bash .claude/skills/actions-sync/scripts/shortcut-plist.sh "My Test Shortcut" > reference.plist
   ```

   This reads the Mac library. For the simulator library, pass `--db "$(xcrun simctl getenv booted HOME)/Library/Shortcuts/Shortcuts.sqlite"`. Use that for actions that do not exist on the Mac.

3. Compile the equivalent Cherri with `--debug` and keep only its actions array, so it has the same shape as the reference:

   ```bash
   go run . file.cherri --debug --output=/tmp/out.shortcut
   plutil -extract WFWorkflowActions xml1 -o cherri.plist "<shortcut name>.plist"
   ```

   `--debug` writes `<shortcut name>.plist` and a `_processed.cherri` next to the source file, named after the shortcut's `#define name`. Delete them afterwards.

4. Diff the two. `UUID`, `OutputUUID` and `GroupingIdentifier` values differ on every compile, so ignore those lines. Any other difference is a bug, per the verification rules in `AGENTS.md`.

Both plists are an array of action dicts, so a plain `diff reference.plist cherri.plist` lines up action by action.

The shortcut name must be unique in the library. The script stops with an error if it is not, or if it is missing.

The check covers structure, not behavior. A matching plist does not prove the action runs correctly, so the signed `go test -run '^TestCherri$'` and the Shortcuts app still apply.

## After finding a gap

Follow the "Action Definition Workflow" in `AGENTS.md`: build the action in Shortcuts, get its plist as above, write the definition, and add a test file in `tests/`.
