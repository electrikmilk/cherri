#!/usr/bin/env bash
# Prints the properties the Shortcuts ToolKit database lists for each content item type Cherri
# knows, which are the values a key read (getAs) on a coerced reference can take.
#
# Usage: content-item-properties.sh [content item class, e.g. WFArticleContentItem ...] [--db PATH]
#
# Without a class it covers every class in the contentItems map in shortcut.go. Defaults to the Mac
# database. For the simulator database, pass
# --db "$(xcrun simctl getenv booted HOME)/Library/Shortcuts/ToolKit/Tools-active".
# The database only describes content items that use App Intents metadata. A type with no rows,
# such as WFContactContentItem, keeps its properties in WorkflowKit code, so a class listed as
# missing needs its properties read from a reference plist or the Shortcuts app instead.
set -euo pipefail

db="$HOME/Library/Shortcuts/ToolKit/Tools-active"
classes=()

while [ $# -gt 0 ]; do
	case "$1" in
	--db)
		db="$2"
		shift 2
		;;
	*)
		classes+=("$1")
		shift
		;;
	esac
done

if [ ! -r "$db" ]; then
	echo "ToolKit database not found at $db. Pass --db to point at another." >&2
	exit 1
fi

repo_root="$(git rev-parse --show-toplevel)"

if [ ${#classes[@]} -eq 0 ]; then
	while IFS= read -r class; do
		classes+=("$class")
	done < <(sed -n '/^var contentItems/,/^}/p' "$repo_root/shortcut.go" | grep -o '"WF[A-Za-z]*ContentItem"' | tr -d '"' | sort -u)
fi

for class in "${classes[@]}"; do
	if [[ ! "$class" =~ ^[A-Za-z]+$ ]]; then
		echo "Content item class must be letters only, like WFArticleContentItem. Got '$class'." >&2
		exit 2
	fi

	echo "## $class"
	properties="$(sqlite3 -readonly "$db" \
		"select p.id from EntityProperties p
		where p.typeId = 'com.apple.shortcuts.$class'
		order by p.id")"
	if [ -z "$properties" ]; then
		echo "(no properties in database)"
	else
		echo "$properties"
	fi
	echo
done
