#!/usr/bin/env bash
# Prints the actions of an installed shortcut as XML plist, exactly as the Shortcuts app stored them.
# Use it as the reference when checking the plist Cherri compiles for an action.
#
# Usage: shortcut-plist.sh <shortcut name> [--db PATH]
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
	echo "Usage: shortcut-plist.sh <shortcut name> [--db PATH]" >&2
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

binary_plist="$work_dir/actions.plist"
sqlite3 -readonly "$db" \
	"select writefile('$binary_plist', ZDATA) from ZSHORTCUTACTIONS
	where ZSHORTCUT = (select Z_PK from ZSHORTCUT where ZNAME = '$escaped_name')" >/dev/null

if [ ! -s "$binary_plist" ]; then
	echo "'$shortcut_name' has no stored actions in $db." >&2
	exit 1
fi

plutil -convert xml1 -o - "$binary_plist"
