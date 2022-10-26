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
actions["takePhoto"] = actionDefinition{
     args: []argumentDefinition{
		{
			field:     "showPreview",
			validType: Bool,
			defaultValue: actionArgument{
				valueType: Bool,
				value:     true,
			},
		},
	},
	call: func(args []actionArgument) []plistData {
		return []plistData{
			argumentValue("WFCameraCaptureShowPreview", args, 0),
		}
	},
}
```

All of these options are optional, as long as the key matches the identifier you could have an action definition as simple as:

```go
actions["identifier"] = actionDefinition{}
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

The `defaultValue` takes an `actionArgument`, you must give a value type and value. This defines a default value if no value
is specified, to require a value, simply don't define this field. Default values should only be given after any required values
otherwise they are pointless.

### `check`

This field takes an `actionCheck` which is a function that accepts a slice of `actionArguments`.

### `call`

This field takes an `actionCall` which is a function that accepts a slice of `actionArguments` and must return a slice
of `plistData`. These usually contain the `argumentValue(key,args,argsIndex)` function which handles the argument value
based on its definition.

For a variable only argument, use the `variableInput(key,value)` function.

You can also obviously directly add a `plistData` value to this slice. This slice will be used as the value of
`WFWorkflowActionParameters` dictionary for the action.

If the action has mutliple arguments without a variable only argument, it's best to return the output of `argumentValues()`
instead. This function takes a reference to the `args` and a `[]paramsMap` slice.

```go
type paramMap struct {
	idx int
	key string
}
```

If it is not necessary to process arguments before they are used as values, simply add the key of the argument to the argument definition like this:

```go
actions["takePhoto"] = actionDefinition{
     args: []argumentDefinition{
		{
			field:     "showPreview",
			validType: Bool,
			key: "WFCameraCaptureShowPreview"
			defaultValue: actionArgument{
				valueType: Bool,
				value:     true,
			},
		},
	},
}
```

If you do this, the argument value for that argument will be done for you, just make sure to not add the call property and add the key property to each of your argument definitions. This is done to help make defining actions simpler and faster.

---

When contributing actions, if an action has a complex number of arguments, try your best to split the action into
multiple actions to reduce the number of arguments and complexity.
