---
title: Non-Standard Actions
layout: default
parent: Contributing
nav_order: 2
---

# Defining Non-Standard Actions

## Define a library

Libraries are defined in the `makeLibraries()` function in [`action.go`](https://github.com/electrikmilk/cherri/blob/main/action.go).

Your library name must be lowercase or camelCase.

```go
// func makeLibraries() {
libraries["app"] = libraryDefinition{
	identifier: "com.company.app",
	make: func(identifier string) {
		appActions(identifier)
	},
}
// ...
```

Create a file called `actions_APP.go`, replace `APP` with a unique name for the app, etc. that you are adding actions for.

Create a function in that file:

```go
func appActions(identifier string) {
// ...
```

Again, replace `app` with the unique name you gave the file or a variant of it.

## Add to an Existing Library

An existing library will have it's own file and a library definition in `makeLibraries()` in [`action.go`](https://github.com/electrikmilk/cherri/blob/main/action.go).

Go to the file for the library (e.g. `actions_APP.go`) and define actions in the same way as explained on this page, but use the `appIdentifier` field instead of the `identifier` field. Unlike standard actions, you must specify a `appIdentifier` field even if it matches the key.

Use the `identifier` provided to the make actions function.

```go
actions["doThing"] = actionDefinition{
     appIdentifier: identifier + "dothing",
     // ...
}
```

Libraries are made available using the `#import {library name}` syntax.

---

When contributing actions, if an action has a complex number of arguments, try your best to split the action into
multiple actions to reduce the number of arguments and complexity.
