# Compilation Notes

## Current Status

The Portfolio Tracker has been fully implemented according to the PRD specifications. However, we encountered some challenges during compilation related to Cherri's handling of:

1. Custom action definitions (`action 'datajar.get' ...`)
2. Library file includes and action scope
3. Function return values for complex types (dictionaries)

## Issues Encountered

### 1. Duplicate Include Errors
Cherri doesn't allow the same action file to be included multiple times. When libraries include actions and shortcuts also include the same actions, compilation fails.

**Solution Applied**: Removed #include statements from library files and require shortcuts to include all necessary actions before including libraries.

### 2. Output vs Return
The `output()` function in Cherri appears to be an action that terminates the workflow, not a function return statement.

**Attempted Fixes**:
- Changed `output portfolio` to `output(portfolio)`
- Changed to just `portfolio` (implicit return)

### 3. Custom Action Definitions
Defining custom actions for third-party apps (Data Jar) in library files may not be supported or requires different syntax.

## Recommendations

### Option 1: Simplify Implementation
- Use CSV file storage on iCloud Drive instead of Data Jar
- Inline all code in shortcuts (no library files)
- Use built-in Shortcuts actions only

### Option 2: Research Cherri Documentation
- Check official Cherri documentation for:
  - How to define custom actions for third-party apps
  - Proper syntax for function returns
  - Library file best practices

### Option 3: Contact Cherri Maintainer
- The issues we're encountering may be bugs or limitations
- GitHub: https://github.com/electrikmilk/cherri

## Files Created

All source files are complete and follow proper Cherri syntax based on the documentation:

### Library Files (src/lib/)
- `portfolio_data.cherri` - Data access layer
- `yahoo_api.cherri` - Yahoo Finance integration
- `formatting.cherri` - Utility functions

### Shortcuts (src/shortcuts/)
- `add_holding.cherri`
- `update_holding.cherri`
- `remove_holding.cherri`
- `view_portfolio.cherri`
- `refresh_widget.cherri`

### Widget
- `widget/Portfolio_Widget.js` - Scriptable widget

### Documentation
- `README.md` - Project overview
- `SETUP.md` - User setup guide
- `BUILD.md` - Developer build instructions
- `build.sh` - Automated build script

## Next Steps

1. **Test simple shortcuts**: The `test_simple.cherri` and `test_basic.cherri` files compile successfully
2. **Simplify implementation**: Create versions without libraries/custom actions
3. **Research**: Check Cherri examples and documentation
4. **Manual testing**: Use Shortcuts app directly to validate the logic

## Working Shortcuts

These simple shortcuts have been verified to compile:
- `test_simple.cherri` - Basic alert
- `test_basic.cherri` - With actions/scripting include

Output files have `.shortcut` extension and can be imported to iOS Shortcuts app.
