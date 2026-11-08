#!/usr/bin/env bash
# Compares the automation trigger families in the Shortcuts ToolKit databases (macOS and the booted
# iOS simulator) against the trigger families Cherri defines in automations.go.
#
# Usage: triggers-sync.sh [--mac-db PATH] [--sim-db PATH]
set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/lib/toolkit.sh"

parse_toolkit_arguments "$@"
prepare_toolkit_databases

# Trigger ids look like com.apple.shortcuts.<family>.<variant>. The family is what a shortcut stores
# as WFTriggerIdentifier, and a variant is chosen by which parameters are set.
export_db() {
	local db="$1"
	local out="$2"
	sqlite3 -readonly -separator $'\t' "$db" \
		"select substr(t.id, 21, instr(substr(t.id, 21), '.') - 1) as family, coalesce(min(l.name), '')
		from Triggers t
		left join TriggerLocalizations l on l.triggerId = t.rowId and l.locale = 'en'
		group by family
		order by family" >"$out"
}

export_db "$mac_db" "$work_dir/mac.tsv"
export_db "$sim_db" "$work_dir/sim.tsv"

cut -f1 "$work_dir/mac.tsv" >"$work_dir/mac.ids"
cut -f1 "$work_dir/sim.tsv" >"$work_dir/sim.ids"

repo_root="$(git rev-parse --show-toplevel)"
sed -n 's/^[[:space:]]*"[A-Za-z]*":[[:space:]]*"\(WF[A-Za-z]*Trigger\)",$/\1/p' "$repo_root/automations.go" |
	sort -u >"$work_dir/cherri.ids"

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

echo "Mac database: $mac_db ($(wc -l <"$work_dir/mac.ids" | tr -d ' ') trigger families)"
echo "Simulator database: $sim_db ($(wc -l <"$work_dir/sim.ids" | tr -d ' ') trigger families)"
echo "Cherri: $(wc -l <"$work_dir/cherri.ids" | tr -d ' ') trigger families"

print_section "In simulator only and missing from Cherri" "$work_dir/sim-only-missing.ids" "$work_dir/sim.tsv"
print_section "In Mac database but missing from Cherri" "$work_dir/mac-missing.ids" "$work_dir/mac.tsv"
print_section "In simulator database but missing from Cherri" "$work_dir/sim-missing.ids" "$work_dir/sim.tsv"
print_section "In Cherri but not in Mac database" "$work_dir/cherri-not-mac.ids"
print_section "In Cherri but not in simulator database" "$work_dir/cherri-not-sim.ids"
