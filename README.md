# installgo

***installgo*** will check <https://go.dev> if updates are available for your installed version of go.  If found you can optionally install the updated version of GO.  You can also reinstall the current version if you installed version is the latest one.

## Usage

```text
  installgo [flags]
  installgo [command]
```

## Available Commands

```text
  completion  Generate the autocompletion script for the specified shell
  get         get the value associated with the given keys from the config file.
  help        Help about any command
  status      status will check for a newer version of GO.
  update      update will install the latest version of GO if not already installed.
  version     Display version information
```

## Flags

```text
      --config string       the config file to use
  -h, --help                help for installgo
  -d, --installdir string   the target directory where Go is installed.
  -v, --version             version for installgo
```

Use `installgo [command] --help` for more information about a command.

# completion

Generate the autocompletion script for installgo for the specified shell.
See each sub-command's help for details on how to use the generated script.

## Usage

```text
  installgo completion [command]
```

## Available Commands

```text
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh
```

## Flags

```text
  -h, --help   help for completion 
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```

Use `installgo completion [command] --help` for more information about a command.

# edit

***edit*** starts an editor to edit the `installgo` configuration file.

## Usage

```text
  installgo list [flags]
```

## Flags

```text
  -h, --help   help for list
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```

# help  

***help*** provides help for any command in the application.  Simply type `installgo help [path to command]` for full details.

## Usage

```text
  installgo help [command] [flags]
```

## Flags

```text
  -h, --help   help for help
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```

# list

***list*** the contents of the config file.

## Usage

```text
  installgo list [flags]
```

## Flags

```text
  -h, --help   help for list
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```

# status

***status*** will check <https://go.dev> for the latest version of GO and optionally install it if the `--autoinstall` option is given.

## Usage

```text
  installgo status [flags]
```

## Flags

```text
  -a, --autoupdate                 install the latest version automatically.
  -h, --help                       help for status
  -m, --maxcachetime float[=0.0]   time (in hours) that the cache is valid for. (default 6)
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```

# update  

***update*** will check <https://go.dev> for the latest version of GO and optionally install it.  You can also have update reinstall the latest version if it is already installed on your system.

## Usage

```text
  installgo update [flags]
```

## Flags

```text
  -a, --autoupdate           install the latest version without asking.
  -h, --help                 help for update
  -m, --maxcachetime float   time (in hours) that the cache is valid for. (default 6)
  -r, --reinstall            reinstall the latest version if already installed.
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```

# version

Display version and detailed build information for installgo.

## Usage

```text
  installgo version [flags]
```

## Flags

```text
  -h, --help   help for version
```

## Global Flags

```text
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
