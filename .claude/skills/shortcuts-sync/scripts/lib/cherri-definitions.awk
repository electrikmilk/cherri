# Turns the raw output of `cherri --action= --debug --no-ansi` into one JSON line per action:
# {"identifier":"<full identifier>","data":<parameters JSON>}
#
# Each action prints its parameters line first and then its definition, whose quoted identifier is the
# full Shortcuts identifier, so the parameters are held until the definition line arrives.

/^parameters: / {
	parameters = substr($0, 13)
}

/^(#builtin)? action (default )?(mac |!mac )?(v[0-9]+> )?(v[0-9]+< )?'/ {
	match($0, /'[^']+'/)
	identifier = substr($0, RSTART + 1, RLENGTH - 2)
	printf "{\"identifier\":\"%s\",\"data\":%s}\n", identifier, parameters
}
