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

package vm

import (
	"math/big"


	"github.com/vm-project/common"
	"github.com/vm-project/dep/statedb"
	"github.com/vm-project/vm/params"
	"github.com/vm-project/dep/asset"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/state"
)

// execConfig provide basic config data for running
//all parameter transaction needed
type execConfig struct {
	//total gas of block
	gp         *GasPool
    //
	vm         *vm.EVM
	//global config
	config     params.ChainConfig
	//
	tx         *types.Transaction
	//
	bc          ChainContext
	//
	header     *types.Header
	//all Transactions in block total used
	usedGas     uint64
	//from
	author     *common.Address
	//
	statedb    *statedb.StateDB
	//asset func
	asset     *asset.Asset
	//
	cfg         vm.Config
	//gas used
	gas        uint64
	//
	gasPrice   *big.Int
	//Gas already get via buy gas
	initialGas uint64
	//contract data,input data
	data       []byte
}

// Message represents a message sent to a contract.
//one msg only process one operation
type Message struct {
	From common.Address
	//FromFrontier() (common.Address, error)
	To  *common.Address
	//one type of asset
	AssetId common.Address
	//
	GasPrice *big.Int
	//
	Gas uint64
	//
	Value *big.Int
	//opertype:  new asset, add issue, transfer
	OpType int
	//assettype:  account , utxo , ...
	AssetType uint64
	//
	Nonce uint64
	CheckNonce  bool
	//contract data,input data
	Data []byte
}

// Receipt represents the results of a transaction.
type Receipt struct {
	// Consensus fields
	PostState         []byte        `json:"root"`
	Status            uint64        `json:"status"`
	Internal          []*InternalTx `json:"internal" `
	CumulativeGasUsed uint64        `json:"cumulativeGasUsed"`
	Bloom             Bloom         `json:"logsBloom"        `
	Logs              []*Log        `json:"logs"             `

	// Implementation fields (don't reorder!)
	TxHash          common.Hash    `json:"transactionHash"`
	ContractAddress common.Address `json:"contractAddress"`
	GasUsed         uint64         `json:"gasUsed"`
}
//zipper vm interface for external use
type zvm interface{

	NewExecConfig()
	//
	NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool) Message
	//NewMessage( Address,  *Address,  uint64,  *big.Int,  uint64,  *big.Int,  []byte,  bool) Message
	//
	TransactionExec( *execConfig, *Message) (*Receipt, uint64, error)
}

