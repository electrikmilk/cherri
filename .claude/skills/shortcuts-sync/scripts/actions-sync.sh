#!/usr/bin/env bash
# Compares standard and Apple actions in the Shortcuts ToolKit databases (macOS and the booted iOS
# simulator) against the actions Cherri defines, matching on the full Shortcuts identifier.
#
# Usage: actions-sync.sh [--cherri PATH] [--mac-db PATH] [--sim-db PATH]
set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/lib/toolkit.sh"

parse_toolkit_arguments "$@"
prepare_toolkit_inputs

export_db() {
	local db="$1"
	local out="$2"
	sqlite3 -readonly -separator $'\t' "$db" \
		"select t.id, coalesce(min(l.name), '')
		from Tools t
		left join ToolLocalizations l on l.toolId = t.rowId and l.locale = 'en'
		where $toolkit_scope
		group by t.id
		order by t.id" >"$out"
}

export_db "$mac_db" "$work_dir/mac.tsv"
export_db "$sim_db" "$work_dir/sim.tsv"

cut -f1 "$work_dir/mac.tsv" >"$work_dir/mac.ids"
cut -f1 "$work_dir/sim.tsv" >"$work_dir/sim.ids"

"$cherri_bin" --action= --debug --no-ansi >"$work_dir/cherri.raw"
awk -f "$(dirname "${BASH_SOURCE[0]}")/lib/cherri-definitions.awk" "$work_dir/cherri.raw" |
	jq -r '.identifier' |
	grep -E '^(is\.workflow\.actions|com\.apple)\.' |
	sort -u >"$work_dir/cherri.ids"

# Both lists must be sorted for comm and join.
comm -13 "$work_dir/cherri.ids" "$work_dir/mac.ids" >"$work_dir/mac-missing.ids"
comm -13 "$work_dir/cherri.ids" "$work_dir/sim.ids" >"$work_dir/sim-missing.ids"
comm -13 "$work_dir/mac.ids" "$work_dir/sim.ids" >"$work_dir/sim-only.ids"
comm -12 "$work_dir/sim-only.ids" "$work_dir/sim-missing.ids" >"$work_dir/sim-only-missing.ids"
comm -23 "$work_dir/cherri.ids" "$work_dir/mac.ids" >"$work_dir/cherri-not-mac.ids"
comm -23 "$work_dir/cherri.ids" "$work_dir/sim.ids" >"$work_dir/cherri-not-sim.ids"

print_section() {
	local title="$1"
	local ids_file="$2"
	local names_file="${3:-}"

	echo
	echo "## $title ($(wc -l <"$ids_file" | tr -d ' '))"
	if [ -n "$names_file" ]; then
		join -t $'\t' "$ids_file" "$names_file"
		return
	fi
	cat "$ids_file"
}

echo "Mac database: $mac_db ($(wc -l <"$work_dir/mac.ids" | tr -d ' ') in scope)"
echo "Simulator database: $sim_db ($(wc -l <"$work_dir/sim.ids" | tr -d ' ') in scope)"
echo "Cherri: $(wc -l <"$work_dir/cherri.ids" | tr -d ' ') distinct identifiers in scope"

print_section "In simulator only and missing from Cherri" "$work_dir/sim-only-missing.ids" "$work_dir/sim.tsv"
print_section "In Mac database but missing from Cherri" "$work_dir/mac-missing.ids" "$work_dir/mac.tsv"
print_section "In simulator database but missing from Cherri" "$work_dir/sim-missing.ids" "$work_dir/sim.tsv"
print_section "In Cherri but not in Mac database" "$work_dir/cherri-not-mac.ids"
print_section "In Cherri but not in simulator database" "$work_dir/cherri-not-sim.ids"
