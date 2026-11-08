# Compares the parameters and enum values Cherri defines with a Shortcuts ToolKit database.
# Run by actions-sync-params.sh with $label, $cherri, $dbparams and $dbenums.

def sample:
  if length > 8
  then (.[:8] | join(", ")) + ", ... (+\(length - 8) more)"
  else join(", ")
  end;

# One entry per identifier, merging actions that share an identifier (variants of one action).
def cherri_by_identifier:
  reduce $cherri[] as $entry ({};
    .[$entry.identifier] = (
      (.[$entry.identifier] // {keys: [], fixed: [], enums: {}, dynamic: false})
      | .dynamic = (.dynamic or $entry.data.dynamic)
      | .keys = ((.keys + ($entry.data.parameters | map(.key))) | unique)
      | .fixed = ((.fixed + ($entry.data.fixed // [])) | unique)
      | .enums = (reduce ($entry.data.parameters[] | select(.enum != null)) as $param (.enums;
          .[$param.key] = (((.[$param.key] // []) + $param.enum) | unique)))
    ));

# A tool with no parameters still gets an empty list so it is compared.
def database_keys:
  (($dbparams[0] // []) | reduce .[] as $row ({};
    .[$row.identifier] += (if $row.key == null then [] else [$row.key] end)));

def database_enums:
  (($dbenums[0] // []) | reduce .[] as $row ({}; .[$row.identifier][$row.key] += [$row.title]));

def title($text; $count): "\n### \($text) (\($count))";

cherri_by_identifier as $cherri_actions
| database_keys as $keys_by_identifier
| database_enums as $enums_by_identifier
| [
    $keys_by_identifier | to_entries[]
    | select($cherri_actions[.key] != null)
    | .key as $identifier
    | .value as $database_keys
    | $cherri_actions[$identifier] as $action
    | {
        identifier: $identifier,
        dynamic: $action.dynamic,
        unexposed: ($database_keys - $action.keys - $action.fixed),
        unknown: ($action.keys - $database_keys),
        enums: [
          ($enums_by_identifier[$identifier] // {}) | to_entries[]
          | .key as $key
          | select($action.keys | any(. == $key))
          | (.value | unique) as $database_values
          | $action.enums[$key] as $cherri_values
          | if $cherri_values == null
            then {key: $key, kind: "free", database_count: ($database_values | length)}
            else {key: $key, kind: "differs", only_database: ($database_values - $cherri_values), only_cherri: ($cherri_values - $database_values)}
            end
        ]
      }
  ] as $rows
| ($rows | map(select(.unexposed | length > 0))) as $unexposed
| ($rows | map(select(.unknown | length > 0))) as $unknown
| ($rows | map(. as $row | .enums[] | select(.kind == "differs" and ((.only_database | length) > 0 or (.only_cherri | length) > 0)) | {identifier: $row.identifier} + .)) as $differs
| ($rows | map(. as $row | .enums[] | select(.kind == "free") | {identifier: $row.identifier} + .)) as $free
| "## \($label) database: \($rows | length) actions defined in both Shortcuts and Cherri",

  title("Shortcuts parameters Cherri does not expose"; $unexposed | length),
  ($unexposed[] | "\(.identifier)\(if .dynamic then " (keys also built in Go, so this may be incomplete)" else "" end): \(.unexposed | join(", "))"),

  title("Cherri parameters Shortcuts does not list"; $unknown | length),
  ($unknown[] | "\(.identifier): \(.unknown | join(", "))"),

  title("Enum values that differ"; $differs | length),
  ($differs[] | "\(.identifier) \(.key)",
    (if (.only_database | length) > 0 then "  only in Shortcuts (\(.only_database | length)): \(.only_database | sample)" else empty end),
    (if (.only_cherri | length) > 0 then "  only in Cherri (\(.only_cherri | length)): \(.only_cherri | sample)" else empty end)),

  title("Free text in Cherri where Shortcuts lists values"; $free | length),
  ($free[] | "\(.identifier) \(.key) (\(.database_count) values)")
