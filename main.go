package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"github.com/ericchiang/pup/selector"
	"io"
	"os"
	"strconv"
	"strings"
)

const VERSION string = "0.1.0"

var (
	// Flags
	inputStream   io.ReadCloser = os.Stdin
	indentString  string        = " "
	maxPrintLevel int           = -1
	printNumber   bool          = false
	printColor    bool          = false
)

// Print to stderr and exit
func Fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

// Print help to stderr and quit
func PrintHelp() {
	helpString := `Usage

    pup [list of css selectors]

Version

    %s

Flags

    -c --color         print result with color
    -f --file          file to read from
    -h --help          display this help
    -i --indent        number of spaces to use for indent or character
    -n --number        print number of elements selected
    -l --limit         restrict number of levels printed
    --version          display version`
	Fatal(helpString, VERSION)
}

// Process command arguments and return all non-flags.
func ProcessFlags(cmds []string) []string {
	var i int
	var err error
	defer func() {
		if r := recover(); r != nil {
			Fatal("Option '%s' requires an argument", cmds[i])
		}
	}()
	nonFlagCmds := make([]string, len(cmds))
	n := 0
	for i = 0; i < len(cmds); i++ {
		cmd := cmds[i]
		switch cmd {
		case "-c", "--color":
			printColor = true
		case "-f", "--file":
			filename := cmds[i+1]
			inputStream, err = os.Open(filename)
			if err != nil {
				Fatal(err.Error())
			}
			i++
		case "-h", "--help":
			PrintHelp()
			os.Exit(1)
		case "-i", "--indent":
			indentLevel, err := strconv.Atoi(cmds[i+1])
			if err == nil {
				indentString = strings.Repeat(" ", indentLevel)
			} else {
				indentString = cmds[i+1]
			}
			i++
		case "-n", "--number":
			printNumber = true
		case "-l", "--limit":
			maxPrintLevel, err = strconv.Atoi(cmds[i+1])
			if err != nil {
				Fatal("Argument for '%s' must be numeric",
					cmds)
			}
			i++
		case "--version":
			Fatal(VERSION)
		default:
			if cmd[0] == '-' {
				Fatal("Unrecognized flag '%s'", cmd)
			}
			nonFlagCmds[n] = cmds[i]
			n++
		}
	}
	return nonFlagCmds[:n]
}

// pup
func main() {
	cmds := ProcessFlags(os.Args[1:])
	root, err := html.Parse(inputStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
	inputStream.Close()
	if len(cmds) == 0 {
		PrintNode(root, 0)
		os.Exit(0)
	}
	selectors := make([]*selector.Selector, len(cmds))
	for i, cmd := range cmds {
		selectors[i], err = selector.NewSelector(cmd)
		if err != nil {
			Fatal("Selector parse error: %s", err)
		}
	}
	currNodes := []*html.Node{root}
	var selected []*html.Node
	for _, selector := range selectors {
		selected = []*html.Node{}
		for _, node := range currNodes {
			selected = append(selected,
				selector.FindAllChildren(node)...)
		}
		currNodes = selected
	}
	if printNumber {
		fmt.Println(len(currNodes))
	} else {
		for _, s := range currNodes {
			// defined in `printing.go`
			PrintNode(s, 0)
		}
	}
}
