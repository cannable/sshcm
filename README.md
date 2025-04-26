# SSH Connection Manager

A simple SSH connection manager written in Go.

This is a work-in-progress rewrite of a project in Tcl.

# Background

This started as a "let's learn Go" thing. I chose to rewrite the Tcl script because
it's a reasonably complicated small project and there are some features I wanted
to add to the tool.


# Status

It's not done. This version of the tool is still missing some features from the
original. Also the documentation is a mess. Things on the to-do list:

- Connection searching
- Import/export (csv and json)
- Actual documentation

# Installation

```
go install github.com/cannable/sshcm/cmd/sshcm
```

# Usage

This section will be expanded in the future. For now, it contains output from the tool's help (thank you cobra!).

```
A simple SSH manager, written in Go, that uses a Sqlite DB.

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
  list        List all connections
  remove      Remove connection
  search      Search for connections
  set         Alter an existing connection
  version     Print program version

Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -h, --help        help for sshcm
  -v, --verbose     Verbose output

Use "sshcm [command] --help" for more information about a command.
```

## Connections

### Add a connection

```
Usage:
  sshcm add [flags]

Aliases:
  add, a

Flags:
  -a, --args string          Arguments to pass to SSH command
  -c, --command string       SSH command to run
  -d, --description string   Short description of the connection
  -h, --help                 help for add
      --host string          Connection hostname (or IP address)
      --identity string      SSH identity to use for connection (a la '-i')
  -n, --nickname string      Nickname for connection
  -u, --user string          User name for connection

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Start a connection

```
Usage:
  sshcm connect [flags]

Aliases:
  connect, c

Flags:
  -a, --args string       Arguments to pass to SSH command
  -c, --command string    SSH command to run
  -h, --help              help for connect
      --identity string   SSH identity to use for connection (a la '-i')
  -u, --user string       User name for connection

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Get connection settings

```
Usage:
  sshcm get [flags]

Aliases:
  get, g

Flags:
  -h, --help   help for get

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### List all connections

```
Usage:
  sshcm list [flags]

Aliases:
  list, l

Flags:
  -a, --all    List all connection details (wide output).
  -h, --help   help for list

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Change connection settings

```
Usage:
  sshcm set [flags]

Aliases:
  set, s

Flags:
  -a, --args string          Arguments to pass to SSH command
  -c, --command string       SSH command to run
  -d, --description string   Short description of the connection
  -h, --help                 help for set
      --host string          Connection hostname (or IP address)
      --identity string      SSH identity to use for connection (a la '-i')
  -n, --nickname string      Nickname for connection
  -u, --user string          User name for connection

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Remove a connection

```
Usage:
  sshcm remove [flags]

Aliases:
  remove, rm, delete, del

Flags:
  -h, --help   help for remove

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```


## Program Defaults

### Change program default settings

```
Usage:
  sshcm def [flags]

Flags:
  -h, --help   help for def

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### List program defaults

```
Usage:
  sshcm defaults [flags]

Flags:
  -h, --help   help for defaults

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```
