# installgo
***installgo*** will check https://go.dev if updates are available for your installed version of go.  If found you can optionally install the updated version of GO.  You can also reinstall the current version if you installed version is the latest one.
## Usage:
```
  installgo [flags]
  installgo [command]
```
## Available Commands:
```
  completion  Generate the autocompletion script for the specified shell
  get         get the value associated with the given keys from the config file.
  help        Help about any command
  status      status will check for a newer version of GO.
  update      update will install the latest version of GO if not already installed.
  version     Display version information
```
## Flags:
```
      --config string       the config file to use
  -h, --help                help for installgo
  -d, --installdir string   the target directory where Go is installed.
  -v, --version             version for installgo
```
Use `installgo [command] --help` for more information about a command.
# completion
Generate the autocompletion script for installgo for the specified shell.
See each sub-command's help for details on how to use the generated script.

## Usage:
```
  installgo completion [command]
```
## Available Commands:
```
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh
```
## Flags:
```
  -h, --help   help for completion
```
## Global Flags:
```
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
Use `installgo completion [command] --help` for more information about a command.
# get
***get*** the values associated with the given keys from the config file.  If no
keys are given, display all the key/value pairs in the config file.

## Usage:
```
  installgo get [flags]
```
## Flags:
```
  -h, --help   help for get
```
## Global Flags:
```
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
# help  
***help*** provides help for any command in the application.  Simply type `installgo help [path to command]` for full details.

## Usage:
```
  installgo help [command] [flags]
```
## Flags:
```
  -h, --help   help for help
```
## Global Flags:
```
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
# status
***status*** will check https://go.dev for the latest version of GO and optionally install it if the `--autoinstall` option is given.
## Usage:
```
  installgo status [flags]
```
## Flags:
```
  -a, --autoupdate                 install the latest version automatically.
  -h, --help                       help for status
  -m, --maxcachetime float[=0.0]   time (in hours) that the cache is valid for. (default 6)
```
## Global Flags:
```
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
# update  
***update*** will check https://go.dev for the latest version of GO and optionally install it.  You can also have update reinstall the latest version if it is already installed on your system.
## Usage:
```
  installgo update [flags]
```
## Flags:
```
  -a, --autoupdate           install the latest version without asking.
  -h, --help                 help for update
  -m, --maxcachetime float   time (in hours) that the cache is valid for. (default 6)
  -r, --reinstall            reinstall the latest version if already installed.
```
## Global Flags:
```
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
# version 
Display version and detailed build information for installgo.
## Usage:
```
  installgo version [flags]
```
## Flags:
```
  -h, --help   help for version
```
## Global Flags:
```
      --config string       the config file to use
  -d, --installdir string   the target directory where Go is installed.
```
