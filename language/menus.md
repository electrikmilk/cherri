[Back](/cherri/language/)

# Menus

You can create menus in 3 ways.

## 1. VCard Menu

You can easily create a VCard menu using the built-in action `makeVCard()`.

```cherri
@vcard = makeVCard("Title","Subtitle","path/to/image.jpg")
```

This generates a VCard using your parameters at compile time, inserting the title as the name, the subtitle as the org/company, and base 64 encodes the image at the path you specified and inserts it as the photo.

In the example above, `vcard` will contain the generated VCard. So, we can then use it to make a VCard menu like so:

```cherri
// Generate VCard menu
@items
repeat 3 {
    @items += makeVCard("{menuPrefix} Title {RepeatIndex}","{menuPrefix} Subtitle","assets/cherri_icon.png")
}
@menuItems = "{items}"
@vcf = setName(menuItems,"menu.vcf",false)

// Coerce type to contact
@contact = vcf.contact

// Use chooseFromList to prompt the user with our menu
@chosenItem = chooseFromList(contact,"Prompt")

// chosenItem contains the title of the chosen item
alert(chosenItem,"You chose:",false)
```

## 2. Choose From Menu

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

## 3. Choose From List

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
