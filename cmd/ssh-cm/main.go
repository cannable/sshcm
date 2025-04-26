/*
ssh-cm is a simple SSH connection manager written in Go.
This is a re-write of a tool originally written in Tcl.

Usage:

	ssh-cm [command]

Available Commands:

	add         Add a connection
	completion  Generate the autocompletion script for the specified shell
	connect     Start a connection
	def         Set program default settings
	defaults    List program defaults
	export      Export all connections
	help        Help about any command
	import      Import connections
	list        list all connections
	remove      Remove connection
	search      Search for connections
	set         Alter an existing connection

Flags:

	    --db string   Path to connection DB file (ssh-cm.connections).
	-h, --help        help for ssh-cm
	-t, --toggle      Help message for toggle
	-v, --verbose     Verbose output

Use "ssh-cm [command] --help" for more information about a command.
*/
package main

import "github.com/cannable/ssh-cm-go/cmd/ssh-cm/cmd"

func main() {
	cmd.Execute()
}
