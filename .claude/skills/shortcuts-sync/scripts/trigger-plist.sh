#!/usr/bin/env bash
# Prints the automations of an installed shortcut as a JSON array, exactly as the Shortcuts app
# stored them. Each entry has WFTriggerIdentifier, WFTriggerSerializedParameters and WFTriggerUUID,
# the same shape as an entry of WFWorkflowTriggers in a shortcut file.
#
# Usage: trigger-plist.sh <shortcut name> [--db PATH]
#
# The database is opened read-only because Shortcuts keeps it open and writes to it.
set -euo pipefail

db="$HOME/Library/Shortcuts/Shortcuts.sqlite"
shortcut_name=""

while [ $# -gt 0 ]; do
	case "$1" in
	--db)
		db="$2"
		shift 2
		;;
	*)
		if [ -n "$shortcut_name" ]; then
			echo "Unexpected argument: $1" >&2
			exit 2
		fi
		shortcut_name="$1"
		shift
		;;
	esac
done

if [ -z "$shortcut_name" ]; then
	echo "Usage: trigger-plist.sh <shortcut name> [--db PATH]" >&2
	exit 2
fi

if [ ! -r "$db" ]; then
	echo "Shortcuts database not found at $db. Pass --db to point at another." >&2
	exit 1
fi

work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT

escaped_name="${shortcut_name//\'/\'\'}"

match_count="$(sqlite3 -readonly "$db" "select count(*) from ZSHORTCUT where ZNAME = '$escaped_name'")"
if [ "$match_count" -eq 0 ]; then
	echo "No shortcut named '$shortcut_name' in $db." >&2
	exit 1
fi
if [ "$match_count" -gt 1 ]; then
	echo "$match_count shortcuts are named '$shortcut_name' in $db. Rename one so the name is unique." >&2
	exit 1
fi

trigger_keys="$(sqlite3 -readonly "$db" \
	"select Z_PK from ZUNIFIEDTRIGGER
	where ZSHORTCUT = (select Z_PK from ZSHORTCUT where ZNAME = '$escaped_name')
	order by Z_PK")"
if [ -z "$trigger_keys" ]; then
	echo "'$shortcut_name' has no automations in $db." >&2
	exit 1
fi

json_files=()
for trigger_key in $trigger_keys; do
	sqlite3 -readonly "$db" \
		"select writefile('$work_dir/$trigger_key.plist', ZDATA) from ZUNIFIEDTRIGGER where Z_PK = $trigger_key" >/dev/null
	plutil -convert json -o "$work_dir/$trigger_key.json" "$work_dir/$trigger_key.plist"
	json_files+=("$work_dir/$trigger_key.json")
done

jq -s '.' "${json_files[@]}"
