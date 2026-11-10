#!/usr/bin/env bash
# Compares the parameters and enum values Cherri defines for each action with what the Shortcuts
# ToolKit databases (macOS and the booted iOS simulator) list for the same identifier.
#
# Usage: actions-sync-params.sh [--cherri PATH] [--mac-db PATH] [--sim-db PATH]
set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/lib/toolkit.sh"

parse_toolkit_arguments "$@"
prepare_toolkit_inputs

"$cherri_bin" --action= --debug --no-ansi >"$work_dir/cherri.raw"

awk -f "$(dirname "${BASH_SOURCE[0]}")/lib/cherri-definitions.awk" "$work_dir/cherri.raw" >"$work_dir/cherri.jsonl"

compare_database() {
	local label="$1"
	local db="$2"

	sqlite3 -readonly -json "$db" \
		"select t.id as identifier, p.key as key
		from Tools t
		left join Parameters p on p.toolId = t.rowId
		where $toolkit_scope" >"$work_dir/$label.params.json"

	sqlite3 -readonly -json "$db" \
		"select t.id as identifier, p.key as key, e.title as title
		from Tools t
		join Parameters p on p.toolId = t.rowId
		join EnumerationCases e on e.typeId = 'com.apple.shortcuts.' || t.id || '.' || p.key and e.locale = 'en'
		where $toolkit_scope" >"$work_dir/$label.enums.json"

	jq -n -r --arg label "$label" \
		--slurpfile cherri "$work_dir/cherri.jsonl" \
		--slurpfile dbparams "$work_dir/$label.params.json" \
		--slurpfile dbenums "$work_dir/$label.enums.json" \
		-f "$(dirname "${BASH_SOURCE[0]}")/lib/compare-params.jq"
}

compare_database "Mac" "$mac_db"
echo
compare_database "Simulator" "$sim_db"
