[Back](../index.md)

# Actions

Actions in Cherri are intended to be easier to use, as in some cases, single actions have been split up into multiple
actions to reduce the number of arguments and complexity.

Some arguments are optional and some are required.

Some actions may not be implemented due to difficulty implementing them practically into the language.

## Standard Actions

These are the standard Shortcuts actions currently supported. Currently, more are being added all the time so this list
may be inaccurate. Not all standard actions in each category are implemented yet, Scripting is the most complete.

- [ ] (3/23) [Calendar](standard/calendar.md)
- [ ] (6/15) [Contacts](standard/contacts.md)
- [ ] (40/64) [Documents](standard/documents.md)
- [ ] (9/17) [Location](standard/location.md)
- [ ] (7/68) [Media](standard/media.md)
- [ ] (54/74) [Scripting](standard/scripting.md)
- [ ] (4/10) [Sharing](standard/sharing.md)
- [ ] (19/26) [Web](standard/web.md)

[Please report incomplete or non-working actions](https://github.com/electrikmilk/cherri/issues)

## Non-standard actions

Once standard actions are finished being implemented, or a non-standard action is contributed, it is planned to implement an `import` syntax to use a library of non-standard actions so that they do not take up memory and compile time if they don't have to.

## Can I contribute actions, even non-standard actions?

Yes.

[Learn more about contributing actions](../compiler/actions.md)
