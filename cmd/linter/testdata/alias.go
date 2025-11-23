package main

import (
	"fmt"
	logalias "log"
	osalias "os"
)

func testAlias() {
	fmt.Println("Hello, world!")
	logalias.Fatal("test") // want "unexpected log.Fatal outside of main package"
	osalias.Exit(1)        // want "unexpected os.Exit outside of main package"
}
