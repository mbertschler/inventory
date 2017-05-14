package main

import (
	"fmt"
	"inventory/db"
	"inventory/parts"
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
	case "part":
		addPart(c)
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

	rmType := c.Args[0]
	rmName := strings.Join(c.Args[1:], " ")

	switch rmType {
	case "type":
		types.Rm(rmName)
	case "part":
		parts.Rm(rmName)
	default:
		fmt.Printf("can't remove %s: unknown type", rmType)
		fmt.Println()
		return
	}

	fmt.Println("removed", rmType, rmName)
}

func list(c *ishell.Context) {
	if len(c.Args) < 1 {
		listTypes(c)
		listParts(c)
		return
	}

	switch c.Args[0] {
	case "types":
		listTypes(c)
	case "parts":
		listParts(c)

	default:
		log.Fatalf("can't list %s: unknown type", os.Args[2])
	}
}

func addPart(c *ishell.Context) {
	// disable the '>>>' for cleaner same line input.
	c.ShowPrompt(false)
	defer c.ShowPrompt(true) // yes, revert after login.

	// get name
	partName := c.Args[1]
	c.Println("Adding part ", partName)
	// get type
	c.Print("Type: ")
	partType := c.ReadLine()
	// get value
	c.Print("Value: ")
	partValue := c.ReadLine()
	// get location
	c.Print("Location: ")
	partLocation := c.ReadLine()
	// get datasheet
	c.Print("Datasheet URL: ")
	partDatasheet := c.ReadLine()

	parts.Add(partName, partType, partValue, partLocation, partDatasheet)
}

func listTypes(c *ishell.Context) {
	c.Println("Types:")
	for _, e := range types.List() {
		c.Println("   ", e)
	}
	c.Println()
}

func listParts(c *ishell.Context) {
	c.Println("Parts:")
	for _, e := range parts.List() {
		c.Println("   ", e)
	}
	c.Println()
}
