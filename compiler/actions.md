[Back](/cherri/compiler/)

# Actions

Defining actions is easy but might be a little complicated to understand at first. Actions are defined in one
place, [`actions.go`](https://github.com/electrikmilk/cherri/blob/main/actions.go).

Actions are added to a map that accepts a key of type `string` and a value of type `actionDefinition`.

An action definition consists of the identifier, the argument definitions, an optional `check` function of custom
type `actionCheck` that is called when the action is checked at parsing, and a call of the custom type `actionCall`.

```go
type actionDefinition struct {
    ident string
    args  []argumentDefinition
    check actionCheck
    call  actionCall
}
```

All types and helper functions associated with constructing actions and checking parsed argument inputs, etc. are
in [`action.go`](https://github.com/electrikmilk/cherri/blob/main/action.go).

Here is an example of a simple action definition:

```go
actions["getType"] = actionDefinition{
    ident: "getitemtype",
    args: []argumentDefinition{
        {
            field:     "input",
            validType: STRING,
        },
    },
    call: func (args []actionArgument) []plistData {
        return []plistData{
            argumentValue("WFInput", args, 0),
        }
    },
}
```

### `ident`

The identifier is the identifier at the end of `is.workflows.actions.**identifier**`.

The identifier is optional if the key of the action matches the action identifier, even in camelCase, as if no `ident`
is specified, not only is the key used instead, but its case will be changed to lowercase.

### `args`

Arguments are defined using a `argumentDefinition` for each argument. It has two fields, one that defines the field
name, this will be used in error messages. The other defines the valid type for the input, this is compared against the
value type of the argument received in parsing. This is also used to know the minimum number of arguments. Both of these
checks happen during parsing right after parsing the arguments for the action.

### `check`

This field takes an `actionCheck` which is a function that accepts a slice of `actionArguments`.

### `call`

This field takes an `actionCall` which is a function that accepts a slice of `actionArguments` and must return a slice
of `plistData`. These usually contain the `argumentValue(key,args,argsIndex)` function which handles the argument value
based on its definition.

Otherwise, you might use a `inputValue(key,name,uuid)` function if the argument only accepts variables.

You can also obviously directly add a `plistData` value to this slice.

This slice will be used as the value of `WFWorkflowActionParameters` dictionary for the action.

---

When contributing actions, if an action has a complex number of arguments, try your best to split the action into
multiple actions to reduce the number of arguments and complexity.
