package main

import (
	"fmt"
	"log"
	"os"
)

func testFunc() {
	fmt.Println("Hello, world!")
	panic("oops")     // want "unexpected panic outside of main package"
	log.Fatal("test") // want "unexpected log.Fatal outside of main package"
	os.Exit(1)        // want "unexpected os.Exit outside of main package"
}
