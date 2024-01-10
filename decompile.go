/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	plists "howett.net/plist"
)

func decompile(b []byte) {
	var data map[string]interface{}
	var _, marshalErr = plists.Unmarshal(b, &data)
	handle(marshalErr)

	fmt.Println(data)
}
