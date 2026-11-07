#!/usr/bin/env bash
# Shared setup for the actions-sync scripts. Source it, do not run it.
#
# Call parse_toolkit_arguments "$@" and then prepare_toolkit_inputs. Afterwards these are set:
# cherri_bin, mac_db, sim_db, work_dir and toolkit_scope.
#
# The databases must only be opened read-only because Shortcuts keeps them open and writes to them.

export LC_ALL=C

# Entities are data types, not actions. flowTool is Siri internal plumbing.
toolkit_scope="(t.id like 'is.workflow.actions.%' or t.id like 'com.apple.%')
	and t.toolType in ('action', 'appIntent', 'siriIntent')
	and t.id not like '%Entity'"

cherri_bin=""
mac_db="$HOME/Library/Shortcuts/ToolKit/Tools-active"
sim_db=""
work_dir=""

parse_toolkit_arguments() {
	while [ $# -gt 0 ]; do
		case "$1" in
		--cherri)
			cherri_bin="$2"
			shift 2
			;;
		--mac-db)
			mac_db="$2"
			shift 2
			;;
		--sim-db)
			sim_db="$2"
			shift 2
			;;
		*)
			echo "Unknown argument: $1" >&2
			exit 2
			;;
		esac
	done
}

prepare_toolkit_inputs() {
	local repo_root
	repo_root="$(git rev-parse --show-toplevel)"
	work_dir="$(mktemp -d)"
	trap 'rm -rf "$work_dir"' EXIT

	if [ -z "$cherri_bin" ]; then
		# The brew binary can lag the repo by months, so always compare against the working tree.
		cherri_bin="$work_dir/cherri"
		(cd "$repo_root" && go build -o "$cherri_bin" .)
	fi

	if [ ! -r "$mac_db" ]; then
		echo "Mac ToolKit database not found at $mac_db. Open the Shortcuts app once, or pass --mac-db." >&2
		exit 1
	fi

	if [ -z "$sim_db" ]; then
		local sim_home
		sim_home="$(xcrun simctl getenv booted HOME)"
		sim_db="$sim_home/Library/Shortcuts/ToolKit/Tools-active"
	fi

	if [ ! -r "$sim_db" ]; then
		echo "Simulator ToolKit database not found at $sim_db. Boot a simulator and open Shortcuts once, or pass --sim-db." >&2
		exit 1
	fi
}
