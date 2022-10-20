[Back](/cherri/language/)

# Includes

Use the following syntax to include other Cherri files in a Cherri file.

```cherri
#include "path/to/include.cherri"
```

The file must exist and be a `.cherri` file.

Includes are checked just before parsing your file. If an include is found, the file at the path will be checked, if
valid it will be replaced with the contents of the file at that path.

You can include other Cherri files at any point in your Cherri file .

Be careful of conflicts between the included code and the current file.
