package csm

import (
	"errors"

	"github.com/hpcloud/go-csm-lib/util/cmdlineargs"
)

type CSMRequest struct {
	WorkspaceID  string
	ConnectionID string
	OutputPath   string
}

//GetCSMRequest returns a request based on the received parameters
//The received parameters are supossed to be passed in a
//POSIX compatible way:
//-o, --output The file where the response will be written
//-w, --workspace The workspace on which the action will take place
//-c, --connection The connection on which the action will take place
func GetCSMRequest(args []string) (*CSMRequest, error) {

	outputpath := cmdlineargs.Param{"-o", "--output", "The file where the response will be written", nil}
	workspace := cmdlineargs.Param{"-w", "--workspace", "The workspace on which the action will take place", nil}
	connection := cmdlineargs.Param{"-c", "--connection", "The connection on which the action will take place", nil}

	arguments := []*cmdlineargs.Param{&outputpath, &workspace, &connection}

	help := cmdlineargs.ParseParamsHasHelp(args, arguments)

	//if help was asked we return it
	if help {
		return nil, errors.New(cmdlineargs.ShowHelp(args[0], arguments))
	}

	if outputpath.Value == nil {
		return nil, errors.New("No output path was passed")
	}

	if workspace.Value == nil {
		return nil, errors.New("No workspace was passed")
	}

	request := CSMRequest{}

	request.OutputPath = *outputpath.Value
	request.WorkspaceID = *workspace.Value

	if connection.Value != nil {
		request.ConnectionID = *connection.Value
	}

	return &request, nil
}
