[Back](index.md)

# Menus

## 1. Choose From Menu

The syntax for Menus is similar to a `switch` statement in other languages. Use the following syntax to create a menu:

```cherri
menu "Prompt" {
    case "Item 1":
        // do something...
    case "Item 2":
        // do something else...
}
```

The menu prompt can be a variable, so can each case label, they also support inserted variables.

## 2. Choose From List

Create a variable with a `list()` action as its value.

Just like in Shortcuts each item must be a string, but you can still insert variables.

```cherri
@list = list("Item 1", "Item 2", "Item 3")
```

Then simply use the `chooseFromList()` action with the list and a prompt.

```cherri
@chosenItem = chooseFromList(list,"Choose a item")
```

`chosenItem` will hold the item chosen from your list by the user.
