package cli

import (
	"fmt"
	"os"
	"strings"
)

type Runable interface {
	Run([]string) int
	Describe() string
	Match([]string) bool
}

type Cmd[T any] struct {
	Name        string
	Description string
	Examples    []Example
	Usage       string
	ParseArgs   func([]string) (T, error)
	Fn          func(T) error
	SubCmds     []Runable
}

func (c *Cmd[T]) Run(args []string) int {
	if shouldPrintHelp(args) {
		fmt.Println(c.Describe())
		fmt.Println("USAGE: ", c.Usage)
		if len(c.Examples) > 0 {
			fmt.Println("EXAMPLES:")
		}
		for _, example := range c.Examples {
			fmt.Println(example.render(c.Name) + "\n")
		}
		if len(c.SubCmds) > 0 {
			fmt.Println("SUB COMMANDS:")
		}
		for _, sub := range c.SubCmds {
			fmt.Println(" - ", sub.Describe())
		}
		return 0
	}

	for _, sub := range c.SubCmds {
		if sub.Match(args) {
			return sub.Run(args[1:])
		}
	}

	if c.ParseArgs == nil {
		fmt.Fprintln(os.Stderr, "failed to parse args")
		fmt.Fprintln(os.Stderr, "USAGE: ", c.Usage)
		return 1
	}

	t, err := c.ParseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid args: ", err)
		fmt.Fprintln(os.Stderr, "USAGE: ", c.Usage)
		return 1
	}

	err = c.Fn(t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "command failed: ", err)
		return 1
	}
	return 0
}

func shouldPrintHelp(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "help":
		return true
	case "-h":
		return true
	case "--help":
		return true
	default:
		return false
	}
}

func (c *Cmd[T]) Describe() string {
	return c.Name + ": " + c.Description
}

func (c *Cmd[T]) Match(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return c.Name == args[0]
}

type Example struct {
	Description string
	Args        []string
	Output      string
	Error       error
}

func (e Example) render(name string) string {
	switch {
	case e.Output != "":
		return e.Description + "\n" +
			"$ " + name + " " + strings.Join(e.Args, " ") + "\n" +
			"\t" + strings.Join(strings.Split(e.Output, "\n"), "\n\t")
	case e.Error != nil:
		return e.Description + "\n" +
			"$ " + name + " " + strings.Join(e.Args, " ") + "\n" +
			"\tcommand failed: " + e.Error.Error()
	default:
		return e.Description + "\n" +
			"$ " + name + " " + strings.Join(e.Args, " ")
	}
}
