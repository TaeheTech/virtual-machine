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

package runtime

import (
	"math"
	"math/big"
	"time"
	"fmt"
	"encoding/json"
	"github.com/vm-project/vm/params"
	"github.com/vm-project/vm"
	"github.com/vm-project/common"
	"github.com/vm-project/dep/asset"
	"github.com/vm-project/dep/crypto"
	"github.com/vm-project/dep/statedb"
	"github.com/vm-project/dep/memdb"
)

// Config is a basic type specifying certain configuration flags for running
// the EVM.
type Config struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config
	asset       *asset.Asset
	State     *statedb.StateDB
	GetHashFn func(n uint64) common.Hash
}

// sets defaults on the config
func setDefaults(cfg *Config) {
	if cfg.ChainConfig == nil {
		//cfg.ChainConfig = &params.ChainConfig{
		//	ChainID:        big.NewInt(1),
		//	HomesteadBlock: new(big.Int),
		//	DAOForkBlock:   new(big.Int),
		//	DAOForkSupport: false,
		//	EIP150Block:    new(big.Int),
		//	EIP155Block:    new(big.Int),
		//	EIP158Block:    new(big.Int),
		//}
	}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	cfg.Origin = common.BytesToAddress([]byte("sender"))
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
}
//create a new evm env
func NewEnv(cfg *Config) *vm.EVM {
	fmt.Println("in NewEnv ...")
	context := vm.Context{
		CanTransfer: vm.CanTransfer,
		Transfer:    vm.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} },

		From:      cfg.Origin,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		GasPrice:    cfg.GasPrice,
	}

	return vm.NewEVM(context, cfg.asset ,*cfg.State,  cfg.EVMConfig)
}

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Executes sets up a in memory, temporarily, environment for the execution of
// the given code. It makes sure that it's restored to it's original state afterwards.
func Execute(code []byte, input []byte, cfg *Config) ([]byte, *statedb.StateDB, error) {
	fmt.Println("in runtime.execute ...")
	if cfg == nil {
		cfg = new(Config)

	}
	setDefaults(cfg)

	if cfg.asset == nil {
		cfg.State, _ = statedb.New(common.Hash{})

		//cfg.State, _ = state.New(parent.Root(), state.NewDatabase())
		//cfg.EvmDB =
		cfg.asset = asset.NewAsset(cfg.State)
	}

	//if cfg.State == nil {
	//	cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(zdb.NewMemDatabase()))
	//
	//	//cfg.State, _ = state.New(parent.Root(), state.NewDatabase())
	//	//cfg.EvmDB =
	//	cfg.EvmDB = asset.NewAsset(cfg.State)
	//}
	fmt.Println("in runtime.execute 2...")
	var (
		ToAddress = common.BytesToAddress([]byte("contractTest"))
		vmenv   = NewEnv(cfg)
		//sender  = vm.AccountRef(cfg.Origin)
		//assetAddr = common.BytesToAddress([]byte("asset"))
		//assetAddr = types.ZipAssetID
	)
	//fmt.Println("in runtime.execute 31...")
	cfg.asset.CreateAccount(ToAddress)
	account1 := common.Address{10}
	info := &asset.AccountAssetInfo{
		Name:     "test",
		Symbol:   "BTC",
		Total:    big.NewInt(2100),
		Decimals: 8,
		Owner:  cfg.Origin  }

	b, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Unexpected error : %v", err)
	}
    //注册并发行
	assetAddress, err := cfg.asset.IssueAsset(asset.AccountModel, account1, string(b))
	if err != nil {
		fmt.Println("Unexpected error : %v", err)
	}
	//增发20
	err = cfg.asset.IncreaseAsset(cfg.Origin , assetAddress, big.NewInt(20))
	if err != nil {
		fmt.Println("Unexpected error : %v", err)
	}

	v := cfg.asset.GetBalance(cfg.Origin , assetAddress)
	value := v.(*big.Int)
	if value.Cmp(big.NewInt(2120)) != 0 {
		fmt.Println("Unexpected error : %v", err)
	}
	fmt.Println("asset value=",value)
	// set the receiver's (the executing contract) code for execution.
	cfg.State.SetCode(ToAddress, code)
	// Call the code with the given configuration.
	ret, _, err := vmenv.Call(
		vm.AccountRef(cfg.Origin),
		ToAddress,
		assetAddress,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	fmt.Println("out runtime.execute ...")
	return ret, cfg.State, err
}

// Create executes the code using the EVM create method
func Create(input []byte, cfg *Config) ([]byte, common.Address, uint64, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = statedb.New(common.Hash{}, statedb.NewDatabase(zdb.NewMemDatabase()))
	}
	var (
		vmenv  = NewEnv(cfg)
		sender = vm.AccountRef(cfg.Origin)
		assetAddr = common.BytesToAddress([]byte("asset"))
	)

	// Call the code with the given configuration.
	code, address, leftOverGas, err := vmenv.Create(
		sender,
		assetAddr,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return code, address, leftOverGas, err
}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
//
// Call, unlike Execute, requires a config and also requires the State field to
// be set.
func Call(address common.Address, input []byte, cfg *Config) ([]byte, uint64, error) {
	setDefaults(cfg)

	vmenv := NewEnv(cfg)
	var assetAddr = common.BytesToAddress([]byte("asset"))
	sender := cfg.State.GetOrNewStateObject(cfg.Origin)
	// Call the code with the given configuration.
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		assetAddr,
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, leftOverGas, err
}
