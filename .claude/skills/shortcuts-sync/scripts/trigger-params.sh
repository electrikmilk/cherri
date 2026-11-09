#!/usr/bin/env bash
# Prints every variant, parameter key and enum value the Shortcuts ToolKit database lists for one
# trigger family, so the permutations of a trigger can be read without building each one by hand.
#
# Usage: trigger-params.sh <family, e.g. WFWifiTrigger> [--db PATH]
#
# Defaults to the Mac database. For the simulator database, pass
# --db "$(xcrun simctl getenv booted HOME)/Library/Shortcuts/ToolKit/Tools-active".
# A parameter with no enum values is free-form (a number, boolean or picker). A family whose
# variants differ only by trigger id, such as WFWifiTrigger, does not list the values that id
# implies, so read those from a reference plist with trigger-plist.sh.
set -euo pipefail

db="$HOME/Library/Shortcuts/ToolKit/Tools-active"
family=""

while [ $# -gt 0 ]; do
	case "$1" in
	--db)
		db="$2"
		shift 2
		;;
	*)
		if [ -n "$family" ]; then
			echo "Unexpected argument: $1" >&2
			exit 2
		fi
		family="$1"
		shift
		;;
	esac
done

if [ -z "$family" ]; then
	echo "Usage: trigger-params.sh <family> [--db PATH]" >&2
	exit 2
fi

if [ ! -r "$db" ]; then
	echo "ToolKit database not found at $db. Pass --db to point at another." >&2
	exit 1
fi

if [[ ! "$family" =~ ^[A-Za-z]+$ ]]; then
	echo "Trigger family must be letters only, like WFWifiTrigger. Got '$family'." >&2
	exit 2
fi

variant_count="$(sqlite3 -readonly "$db" "select count(*) from Triggers where id like 'com.apple.shortcuts.$family.%'")"
if [ "$variant_count" -eq 0 ]; then
	echo "No trigger family '$family' in $db." >&2
	exit 1
fi

sqlite3 -readonly -separator $'\t' "$db" \
	"select t.id, coalesce(l.name, ''), coalesce(p.key, ''), coalesce(group_concat(e.id || '=' || coalesce(e.title, ''), ' | '), '')
	from Triggers t
	left join TriggerLocalizations l on l.triggerId = t.rowId and l.locale = 'en'
	left join TriggerParameters p on p.triggerId = t.rowId
	left join EnumerationCases e on e.typeId = p.typeId and e.locale = 'en'
	where t.id like 'com.apple.shortcuts.$family.%'
	group by t.id, p.key
	order by t.id, p.sortOrder"
