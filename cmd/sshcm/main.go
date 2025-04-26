/*
sshcm is a simple SSH connection manager written in Go.
This is a re-write of a tool originally written in Tcl.

Usage:

	sshcm [command]

Available Commands:

	add         Add a connection
	completion  Generate the autocompletion script for the specified shell
	connect     Start a connection
	def         Set program default settings
	defaults    List program defaults
	export      Export all connections
	get         Print existing connection details
	help        Help about any command
	import      Import connections
	list        list all connections
	remove      Remove connection
	search      Search for connections
	set         Alter an existing connection
	version     Print program version

Flags:

	    --db string   Path to connection DB file (ssh-cm.connections).
	-h, --help        help for sshcm
	-t, --toggle      Help message for toggle
	-v, --verbose     Verbose output

Use "sshcm [command] --help" for more information about a command.
*/
package main

import "github.com/cannable/sshcm/cmd/sshcm/cmd"

func main() {
	cmd.Execute()
}
