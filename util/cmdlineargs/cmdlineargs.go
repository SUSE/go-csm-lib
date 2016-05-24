package cmdlineargs

import "strings"
import "fmt"

//Param - defines the properties of a parameter to the command line
type Param struct {
	Shortform   string
	Longform    string
	Description string
	Value       *string
}

//ShowHelp returns the help string based on the parameters you require
func ShowHelp(programName string, params []*Param) string {
	rez := fmt.Sprintf("Usage of %s:\n", programName)
	//fmt.Printf("Usage of %s:\n", programName)
	for _, p := range params {
		rez = fmt.Sprintf("%s%s, %s string\n\t%s\n", rez, p.Shortform, p.Longform, p.Description)
	}
	rez = fmt.Sprintf("%s-h, --help\n\tShow this text\n", rez)
	return rez
}

func identify(pos int, args []string, shortform string, longform string) (*string, bool) {
	var value string
	token := args[pos]
	if token == shortform || token == longform {
		if pos < (len(args) - 1) {
			value = args[pos+1]
			return &value, true
		}
		return nil, false
	} else if strings.HasPrefix(token, shortform) {
		value = strings.TrimPrefix(token, shortform)
		return &value, true
	}
	return nil, false
}

//ParseParamsHasHelp parses the parameters received by the program and
//puts the values in the parameters defined by you
//returns true if user asked for help, false otherwise
func ParseParamsHasHelp(args []string, params []*Param) bool {
	var (
		help bool
	)
	argpos := 0
	help = false
	for argpos < len(args) {
		if args[argpos] == "-h" || args[argpos] == "--help" {
			argpos++
			help = true
			break
		}
		for _, p := range params {
			v, found := identify(argpos, args, p.Shortform, p.Longform)
			if found {
				p.Value = v
			}
		}

		argpos++
	}
	return help
}
