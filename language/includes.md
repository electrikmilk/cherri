---
title: Includes
layout: default
parent: Documentation
nav_order: 6
---

# Includes

Use the following syntax to include other Cherri files in a Cherri file.

```c
#include "path/to/include.cherri"
```

The file must exist and be a `.cherri` file.

Includes are checked just before parsing your file. If an include is found, the file at the path will be checked, if
valid it will be replaced with the contents of the file at that path.

The path of your file will be added to the beginning of every include path unless it contains `..`.

You can include other Cherri files at any point in your Cherri file.

You can include Cherri files in files that you include creating an include chain, but a single file can only be included once ever.

Be careful of conflicts between the included code and the current file.
