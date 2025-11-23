package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	fmt.Println("Hello, world!")
	panic("oops") // want "unexpected panic outside of main package"
	log.Fatal("test")
	os.Exit(1)

	var wg sync.WaitGroup
	wg.Go(func() {
		fmt.Println("Hello, goroutine!")
		panic("oops goroutine")     // want "unexpected panic outside of main package"
		log.Fatal("test goroutine") // want "unexpected log.Fatal outside of main package"
		os.Exit(1)                  // want "unexpected os.Exit outside of main package"
	})
	wg.Wait()

	// никакой пощады к IIFE
	func() {
		fmt.Println("Hello, IIFE!")
		panic("oops IIFE")     // want "unexpected panic outside of main package"
		log.Fatal("test IIFE") // want "unexpected log.Fatal outside of main package"
		os.Exit(1)             // want "unexpected os.Exit outside of main package"
	}()

	callback(func() {
		fmt.Println("Hello, callback!")
		panic("oops callback")     // want "unexpected panic outside of main package"
		log.Fatal("test callback") // want "unexpected log.Fatal outside of main package"
		os.Exit(1)                 // want "unexpected os.Exit outside of main package"
	})
}

func callback(f func()) {
	f()
}
