[Back](/cherri/language/)

# VCard Menus

You can easily create a VCard menu using the built-in action `makeVCard()`.

```cherri
@vcard = makeVCard("Title","Subtitle","path/to/image.jpg")
```

This generates a VCard using your parameters at compile time, inserting the title as the name, the subtitle as the org/company, and base 64 encodes the image at the path you specified and inserts it as the photo.

In the example above, `vcard` will contain the generated VCard (a text block containing the generated VCard). So, we can then use it to make a VCard menu like so:

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