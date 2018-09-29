// Copyright 2018 The zipper team Authors
// This file is part of the z0 library.
//
// The z0 library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The z0 library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the z0 library. If not, see <http://www.gnu.org/licenses/>.

// vm executes VM bytecode.
package main

import (
	"os"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
	"github.com/vm-project/vm/log"
	"github.com/vm-project/utils"
	"fmt"
	"math/big"
)



// NewApp creates an app with sane defaults.
func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = "albertFan"
	//app.Authors = nil
	app.Email = ""
	app.Version = "1.0"
	app.Usage = usage
	return app
}
// BigFlag is a command line flag that accepts 256 bit big integers in decimal or
// hexadecimal syntax.
type BigFlag struct {
	Name  string
	Value *big.Int
	Usage string
}

var (
	app = NewApp("the vm command line interface")

	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "output full trace logs",
	}
	MemProfileFlag = cli.StringFlag{
		Name:  "memprofile",
		Usage: "creates a memory profile at the given path",
	}
	CPUProfileFlag = cli.StringFlag{
		Name:  "cpuprofile",
		Usage: "creates a CPU profile at the given path",
	}
	StatDumpFlag = cli.BoolFlag{
		Name:  "statdump",
		Usage: "displays stack and heap memory information",
	}
	CodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "VM code",
	}
	CodeFileFlag = cli.StringFlag{
		Name:  "codefile",
		Usage: "File containing VM code. If '-' is specified, code is read from stdin ",
	}
	GasFlag = cli.Uint64Flag{
		Name:  "gas",
		Usage: "gas limit for the vm",
		Value: 10000000000,
	}
	PriceFlag = utils.BigFlag{
		Name:  "price",
		Usage: "price set for the evm",
		Value: new(big.Int),
	}
	ValueFlag = utils.BigFlag{
		Name:  "value",
		Usage: "value set for the evm",
		Value: new(big.Int),
	}
	DumpFlag = cli.BoolFlag{
		Name:  "dump",
		Usage: "dumps the state after the run",
	}
	InputFlag = cli.StringFlag{
		Name:  "input",
		Usage: "input for the VM",
	}
	VerbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "sets the verbosity level",
	}
	CreateFlag = cli.BoolFlag{
		Name:  "create",
		Usage: "indicates the action should be create rather than call",
	}
	GenesisFlag = cli.StringFlag{
		Name:  "prestate",
		Usage: "JSON file with prestate (genesis) config",
	}
	MachineFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output trace logs in machine readable format (json)",
	}
	SenderFlag = cli.StringFlag{
		Name:  "sender",
		Usage: "The transaction origin",
	}
	ReceiverFlag = cli.StringFlag{
		Name:  "receiver",
		Usage: "The transaction receiver (execution context)",
	}
	DisableMemoryFlag = cli.BoolFlag{
		Name:  "nomemory",
		Usage: "disable memory output",
	}
	DisableStackFlag = cli.BoolFlag{
		Name:  "nostack",
		Usage: "disable stack output",
	}
)

func init() {

	app.Flags = []cli.Flag{
		CreateFlag,
		DebugFlag,
		VerbosityFlag,
		CodeFlag,
		CodeFileFlag,
		GasFlag,
		PriceFlag,
		ValueFlag,
		DumpFlag,
		InputFlag,
		MemProfileFlag,
		CPUProfileFlag,
		StatDumpFlag,
		GenesisFlag,
		MachineFlag,
		SenderFlag,
		ReceiverFlag,
		DisableMemoryFlag,
		DisableStackFlag,
	}
	app.Commands = []cli.Command{
		//compileCommand,
		//disasmCommand,
		runCommand,
		//stateTestCommand,
	}
}

func main() {
	//set log level
	vmlog.Config("VM",vmlog.DebugLevel);
	vmlog.DebugPrint("vm start ...\n")
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	vmlog.InfoPrint("vm end .\n")
}
