---
title: Actions
layout: default
parent: Contributing
nav_order: 1
---

# Defining Actions

Defining actions is easy but might be a little complicated to understand at first. Standard actions are defined in one
place, [`actions_std.go`](https://github.com/electrikmilk/cherri/blob/main/actions_std.go).

Non-standard actions should be added to a new or existing file for those actions.

Actions are added to a map that accepts a key of type `string` and a value of type `actionDefinition`.

An action definition consists of the identifier, the parameter definitions, an optional `check` function of custom
type `paramCheck` that is called when the action is checked at parsing, and `make` of the custom type `makeParams`
that returns the parameters for the action.

```go
type actionDefinition struct {
	identifier    string
	appIdentifier string
	parameters    []parameterDefinition
	check         paramCheck
	make          makeParams
}
```

All types and helper functions associated with constructing actions and checking parsed argument inputs, etc. are
in [`action.go`](https://github.com/electrikmilk/cherri/blob/main/action.go).

Here is an example of a simple action definition:

```go
actions["takePhoto"] = actionDefinition{
     parameters: []parameterDefinition{
		{
			field:     "showPreview",
			validType: Bool,
			defaultValue: actionArgument{
				valueType: Bool,
				value:     true,
			},
		},
	},
	make: func(args []actionArgument) []plistData {
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

### `identifier`

The identifier is the identifier at the end of `is.workflows.actions.**identifier**`.

The identifier is optional if the key of the action matches the action identifier, even in camelCase, as if no `identifier`
is defined, not only is the key used instead, but its case will be changed to lowercase.

So there is no need to do the below example, remove the ident property and the key will be used instead.

```go
actions["takePhoto"] = actionDefinition{
     identifier: "takephoto",
     // ...
}
```

`identifier` exists so that you can do this:

```go
actions["takeMorePhotos"] = actionDefinition{
     identifier: "takephoto",
     // ...
}
```

### `parameters`

Parameters are defined using a `parameterDefinition` for each parameter. It has two main fields, one that defines the argument
`name` for the action, this will be used in error messages. The other defines the valid type for the argument, this is compared against the
value type of the argument received in parsing. This is also used to know the minimum number of arguments for the action. Both of these
checks happen during parsing, right after parsing the arguments for the action.

Parameters correlate almost directly to arguments.

```go
type parameterDefinition struct {
	name         string
	validType    tokenType
	key          string
	defaultValue actionArgument
	optional     bool
	noMax        bool
}
```

The `defaultValue` takes an `actionArgument`, you must give a value type and value. This defines a default value for this argument, this is used to compare the value given to the default value for this action paramter, we then print a warning that this argument value could be removed. This is mainly for booleans and enums.

`optional` tells Cherri that this parameter is completely optional for this action. this defaults to false. It will not write a key value pair for this parameter if it is optional and no argument value is given.

The Shortcuts app does this with many parameters that it has a default value for, if an actions parameter is not specified, it fills in the gap and goes with the default value for that parameter for that action as defined in the Shortcuts app itself for that action. This should be done when possible as it makes for a much smaller Shortcut file on average.

The `key` is used if there is no need to process arguments. If you use `key` to specify the key of this action parameter, do not add a `make` to the action definition, otherwise this is pointless and will be ignored.

### `check`

This field takes an `paramCheck` which is a function that accepts a slice of `actionArguments`.

### `make`

This field takes an `makeParams` which is a function that accepts a slice of `actionArguments` and must return a slice
of `plistData`. These usually contain the `argumentValue(key,args,argsIndex)` function which handles the argument value
based on its definition.

For a variable only argument, use the `variableInput(key,value)` function. `variableInput(key,value)` is used when a parameter uses an input value, usually with the key `WFInput`, these are parameters whose values must be inserted as a variable value.

You can also obviously directly add a `plistData` value to this slice. This slice will be used as the value of
`WFWorkflowActionParameters` dictionary for the action.

If the action has mutliple arguments without a variable only argument, it's best to return the output of `argumentValues()`
instead. This function takes a reference to the `args` and an unlimited strings argument of the keys for each parameter.

If it is not necessary to process arguments before they are used as parameter values, simply add the key of the argument to the parameter definition like this:

```go
actions["takePhoto"] = actionDefinition{
     parameters: []parameterDefinition{
		{
			name:     "showPreview",
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

If you do this, the parameter value for that argument will be done for you, just make sure to not add the `make` property and add the `key` property to each of your parameter definitions. This is done to help make defining actions simpler and faster.

---

When contributing actions, if an action has a complex number of arguments, try your best to split the action into
multiple actions to reduce the number of arguments and complexity.
