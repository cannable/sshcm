# SSH Connection Manager

A simple SSH connection manager written in Go.

This is a work-in-progress rewrite of a project in Tcl.

# Background

This started as a "let's learn Go" thing. I chose to rewrite the Tcl script because
it's a reasonably complicated small project and there are some features I wanted
to add to the tool.


# Status

It's not done. This version of the tool is still missing one feature from the
original. Also the documentation could be better. Things on the to-do list:

- Schema upgrades
- Actual documentation

# Installation

```
go install github.com/cannable/sshcm/cmd/sshcm@latest
```

## Switching from the Tcl script

You should take a backup of your connection DB file before running sshcm.

The first official release (v1.2.0) of this tool uses DB schema 1.1. The tool is
able to use connection DBs from the last version of the Tcl script as-is.
Additionally, the Tcl script should continue to work until an sshcm release
upgrades the schema.

The Tcl script is not able to use connection DBs created by sshcm. Export/import
from sshcm to the Tcl tool should work, with the proviso that connection IDs may
change. If you need to avoid this, you can try opening the connection DB file in
a sqlite tool and flipping the version setting from "v1.1" to "1.1".

# Usage

```
This section will be expanded in the future. For now, it contains output from the tool's help (thank you cobra!).

A simple SSH manager, written in Go, that uses a Sqlite DB.

Usage:
  sshcm [command]

Available Commands:
  add         Add a connection
  completion  Generate the autocompletion script for the specified shell
  connect     Start a connection
  def         Set program default settings
  defaults    List program defaults
  get         Print existing connection settings
  help        Help about any command
  list        List all connections
  remove      Remove a connection
  set         Change connection settings
  version     Print program version

Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -h, --help        help for sshcm
  -v, --verbose     Verbose output

Use "sshcm [command] --help" for more information about a command.
```

## Connections

### Add a connection

Add a new connection.

All connection settings are expected to be passed via flags. Most are optional,
but a nickname and host are required. The nickname must be unique.

```
Usage:
  sshcm add [flags]

Aliases:
  add, a

Examples:

sshcm add --nickname something --user me --host 127.0.0.1

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

Start a connection.

A connection ID or nickname must be specified as the only positional argument.

Some connection settings (ex. command) can be overridden at runtime by passing flags.

```
Usage:
  sshcm connect { id | nickname } [flags]

Aliases:
  connect, c

Examples:

sshcm connect something
sshcm c 22
sshcm c something --user=someone


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

Print connection settings.

A valid connection ID or nickname must be specified.

```
Usage:
  sshcm get { id | nickname } [flags]

Aliases:
  get, g

Examples:

sshcm get asdf
sshcm g 42


Flags:
  -h, --help   help for get

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Search for connections

Search for connections. Search is case-insensitive.

Search references the following properties:

- nickname
- host
- user
- description

```
Usage:
  sshcm search [flags]

Aliases:
  search, f

Flags:
  -a, --all    List all connection details (wide output).
  -h, --help   help for search

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output

```

### List all connections

List all connections.

```
Usage:
  sshcm list [flags]

Aliases:
  list, l

Examples:

sshcm list

Flags:
  -a, --all    List all connection details (wide output).
  -h, --help   help for list

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Change connection settings

Change connection settings.
A valid ID or nickname must be specified.

A connection can be renamed by passing  `--nickname="new_nickname"`.

```
Usage:
  sshcm set { id | nickname } [flags]

Aliases:
  set, s

Examples:

sshcm set 42 --user="blarg"
sshcm s asdf --nickname fdsa

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

Remove a connection.

A valid connection ID or nickname must be specified.

```
Usage:
  sshcm remove { id | nickname } [flags]

Aliases:
  remove, rm, delete, del

Examples:

sshcm rm asdf
sshcm delete 42

Flags:
  -h, --help   help for remove

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```


## Program Defaults

### Change program default settings

Set program default settings.

```
Usage:
  sshcm def setting value [flags]

Examples:

sshcm def user asdf


Flags:
  -h, --help   help for def

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### List program defaults

List program defaults.

```
Usage:
  sshcm defaults [flags]

Examples:

sshcm defaults

Flags:
  -h, --help   help for defaults

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

## Import/Export

### Import connections

Import connections from standard input (default) or a file.

The import process will update existing connections and append new ones.

The default format is CSV. To use json, pass `--format json`.

```
Usage:
  sshcm import [flags]

Flags:
      --format string   Export format. Valid formats: csv or json. (default "csv")
  -h, --help            help for import
  -f, --path string     Export destination path.

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```

### Export connections

Export connections to standard output (default) or a file.

The export process will update existing connections and append new ones.

The default format is CSV. To use json, pass `--format json`.

```
Usage:
  sshcm export [flags]

Flags:
      --format string   Export format. Valid formats: csv or json. (default "csv")
  -h, --help            help for export
  -f, --path string     Export destination path.

Global Flags:
      --db string   Path to connection DB file (ssh-cm.connections).
  -v, --verbose     Verbose output
```
