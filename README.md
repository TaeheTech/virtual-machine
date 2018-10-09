# Build

go build -o vm.exe ./vm/cmd/

# usage

C:\gopath\src\github.com\vm-project>vm.exe

NAME:
   vm.exe - the vm command line interface

USAGE:
   vm.exe [global options] command [command options] [arguments...]

VERSION:
   1.0

AUTHOR:
   albertFan

COMMANDS:
     compile  compiles easm source to evm binary
     disasm   disassembles evm binary
     run      run arbitrary vm binary
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --create            indicates the action should be create rather than call
   --debug             output full trace logs
   --verbosity value   sets the verbosity level (default: 0)
   --code value        VM code
   --codefile value    File containing VM code. If '-' is specified, code is read from stdin
   --gas value         gas limit for the vm (default: 10000000000)
   --price "0"         price set for the evm
   --value "0"         value set for the evm
   --dump              dumps the state after the run
   --input value       input for the VM
   --memprofile value  creates a memory profile at the given path
   --cpuprofile value  creates a CPU profile at the given path
   --statdump          displays stack and heap memory information
   --prestate value    JSON file with prestate (genesis) config
   --json              output trace logs in machine readable format (json)
   --sender value      The transaction origin
   --receiver value    The transaction receiver (execution context)
   --nomemory          disable memory output
   --nostack           disable stack output
   --help, -h          show help
   --version, -v       print the version
#example
   
vm -input=0xa9059cbb000000000000000000000000d90ed3659b1953fddde628183dacb316b9a3cc89000000000000000000000000000000000000000000013da329b6336471800000 

-dump -sender=0x915d7915f2b469bb654a7d903a5d4417cb8ea7df 

-receiver=0xa9d2927d3a04309e008b6af6e2e282ae2952e7fd 

-price=18 -codefile=e:/sol/zip.bin run


