package cmdlineargs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParamsHasHelpLongFormParams(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for no help asked, long form parameters")
	commandLineArgs := []string{"--one", "firstparam", "--two", "secondparam"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(false, help, "No help parameter was passed, but one was incorrectly detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Equal("secondparam", *param2.Value, "Incorrect value received, should be secondparam")
}
func TestParseParamsHasHelpShortFormParamsUnSeparated(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for no help asked, short form parameters unseparated")
	commandLineArgs := []string{"-ofirstparam", "-tsecondparam"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(false, help, "No help parameter was passed, but one was incorrectly detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Equal("secondparam", *param2.Value, "Incorrect value received, should be secondparam")

}
func TestParseParamsHasHelpShortFormParamsSeparated(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	//for a succesful check
	t.Log("Cheching for no help asked, short form parameters separated")
	commandLineArgs := []string{"-o", "firstparam", "-t", "secondparam"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(false, help, "No help parameter was passed, but one was incorrectly detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Equal("secondparam", *param2.Value, "Incorrect value received, should be secondparam")
}

func TestParseParamsHasHelpLongFormParamsOneUnSeparated(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for no help asked, one incorect parameter, long form")
	commandLineArgs := []string{"--one", "firstparam", "--twos", "secondparam"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(false, help, "No help parameter was passed, but one was incorrectly detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Nil(param2.Value, "Incorrect value received, should be nil")

}
func TestParseParamsHasHelpShortFormParamsOneWrong(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for no help asked, one incorect parameter, shortform separated")
	commandLineArgs := []string{"-o", "firstparam", "-x", "secondparam"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(false, help, "No help parameter was passed, but one was incorrectly detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Nil(param2.Value, "Incorrect value received, should be nil")

}

func TestParseParamsHasHelpShortFormParamsWithHelp(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for help asked in shortform, both parameters supplied")
	commandLineArgs := []string{"-o", "firstparam", "-t", "secondparam", "-h"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(true, help, "Help parameter was passed, and was not detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Equal("secondparam", *param2.Value, "Incorrect value received, should be secondparam")
}

func TestParseParamsHasHelpLongFormParamsWithHelp(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for help asked in longform, both parameters supplied")
	commandLineArgs := []string{"-o", "firstparam", "-t", "secondparam separated", "--help"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(true, help, "Help parameter was passed, and was not detected")
	assert.Equal("firstparam", *param1.Value, "Incorrect value received, should be firstparam")
	assert.Equal("secondparam separated", *param2.Value, "Incorrect value received, should be secondparam")
}

func TestParseParamsHasHelpWithHelpLongForm(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for help asked in longform, no other parameter supplied")
	commandLineArgs := []string{"--help"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(true, help, "Help parameter was passed, and was not detected")
}
func TestParseParamsHasHelpWithHelpShortForm(t *testing.T) {
	assert := assert.New(t)

	param1 := Param{"-o", "--one", "The first parameter", nil}
	param2 := Param{"-t", "--two", "The second parameter", nil}
	arguments := []*Param{&param1, &param2}

	t.Log("Cheching for help asked in shortform, no other parameter supplied")
	commandLineArgs := []string{"-h"}
	help := ParseParamsHasHelp(commandLineArgs, arguments)
	assert.Equal(true, help, "Help parameter was passed, and was not detected")
}
