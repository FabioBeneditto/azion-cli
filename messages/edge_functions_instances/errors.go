package edge_functions_instances

import "errors"

var (
	ErrorMissingArgumentsDelete = errors.New("Required flags are missing. You must supply application-id and instance-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDeleteFuncInst   = errors.New("Failed to delete the Edge Function Instance: %s. Check your settings and try again. If the error persists, contact Azion support")
)