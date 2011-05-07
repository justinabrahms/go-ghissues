package main

import (
	"github.com/justinlilly/go-ghissues"
	"fmt"
)

func main() {
	c := ghissues.NewClient("", "") // name/token unnecessary for list call.
	list, err := c.List("technomancy", "emacs-starter-kit", "open")
	if (err != nil) {
		fmt.Println("Error: " + err.String())
	}
	for i := range list {
		item := list[i]
		fmt.Printf("%v: %s\n", item.Number, item.Title)
	}
}
