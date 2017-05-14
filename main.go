package main

import (
	"fmt"
	"inventory/db"
	"inventory/types"
	"log"
	"os"

	"strings"

	"github.com/abiosoft/ishell"
)

func main() {
	db.Open()
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("Inventory Manager v0.1")

	shell.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "add something to the db",
		Func: func(c *ishell.Context) {
			add(c)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "rm",
		Help: "remove something from the db",
		Func: func(c *ishell.Context) {
			rm(c)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "ls",
		Help: "list db contents",
		Func: func(c *ishell.Context) {
			list(c)
		},
	})

	// start shell
	shell.Start()

}

func add(c *ishell.Context) {
	if len(c.Args) < 2 {
		fmt.Println("usage: add type name")
		return
	}

	addType := c.Args[0]
	addName := strings.Join(c.Args[1:], " ")

	switch addType {
	case "type":
		types.Add(addName)
	default:
		fmt.Printf("can't add %s: unknown type", addType)
		fmt.Println()
		return
	}

	fmt.Println("added", addType, addName)
}

func rm(c *ishell.Context) {
	if len(c.Args) < 2 {
		fmt.Println("usage: rm type name")
		return
	}

	addType := c.Args[0]
	addName := strings.Join(c.Args[1:], " ")

	switch addType {
	case "type":
		types.Rm(addName)
	default:
		fmt.Printf("can't remove %s: unknown type", addType)
		fmt.Println()
		return
	}

	fmt.Println("removed", addType, addName)
}

func list(c *ishell.Context) {
	if len(c.Args) < 1 {
		listTypes()
		listParts()
		return
	}

	switch c.Args[0] {
	case "types":
		listTypes()
	case "parts":
		listParts()

	default:
		log.Fatalf("can't list %s: unknown type", os.Args[2])
	}
}

func listTypes() {
	fmt.Println("Types:")
	for _, e := range types.List() {
		fmt.Println("   ", e)
	}
	fmt.Println()
}

func listParts() {
	fmt.Println("parts:")
	// for _, e := range parts.Parts {
	// 	fmt.Println("   ", e)
	// }
	fmt.Println()
}
