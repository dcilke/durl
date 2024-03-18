package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dcilke/durl/pkg/urlax"
	flags "github.com/jessevdk/go-flags"
)

type Cmd struct {
	Password bool `short:"p" long:"password" description:"Extract password from URL"`
	Encode   bool `short:"e" long:"encode" description:"Encode URL"`
	Decode   bool `short:"d" long:"decode" description:"Decode URL"`
}

func main() {
	// parse command line flags
	var cmd Cmd
	parser := flags.NewParser(&cmd, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "[URL]"
	url, err := parser.Parse()
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		parser.WriteHelp(os.Stdout)
		return
	}
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("unable to parse arguments: %w", err))
	}

	if len(url) > 0 {
		for _, arg := range url {
			cmd.process(arg)
		}
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		cmd.process(line)
	}
}

func (c *Cmd) process(arg string) {
	u, err := urlax.Parse(arg)
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("unable to parse url %q: %w", arg, err))
		return
	}

	switch {
	case c.Password:
		password, _ := u.User.Password()
		fmt.Println(password)
	case c.Encode:
		fmt.Println(u.String())
	case c.Decode:
		fmt.Println(urlax.Decode(u))
	}
}
