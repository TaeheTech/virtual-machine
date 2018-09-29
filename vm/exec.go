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
	"errors"
	"math"
	"math/big"

	"github.com/vm-project/vm/params"
	"github.com/vm-project/common"

	"github.com/zipper-project/z0/consensus"
	"github.com/zipper-project/z0/types"
	"github.com/zipper-project/z0/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

const (
	//new issue an asset
	NEWASSET  =  iota
	//add asset amount
	ADDASSET
	//asset operation
	OPASSET
	//change asset owner
    OPOWNER
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}


func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool) Message {
	return Message{
		From:       from,
		To:         to,
		Nonce:      nonce,
		//amount:     amount,
		//gasLimit:   gasLimit,
		GasPrice:   gasPrice,
		Data:       data,
		//checkNonce: checkNonce,
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	var cache map[uint64]common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		if cache == nil {
			cache = map[uint64]common.Hash{
				ref.Number.Uint64() - 1: ref.ParentHash,
			}
		}
		// Try to fulfill the request from the cache
		if hash, ok := cache[n]; ok {
			return hash
		}
		// Not cached, iterate the blocks and cache the hashes
		for header := chain.GetHeader(ref.ParentHash, ref.Number.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
			cache[header.Number.Uint64()-1] = header.ParentHash
			if n == header.Number.Uint64()-1 {
				return header.ParentHash
			}
		}
		return common.Hash{}
	}
}
// NewEVMContext creates a new context for use in the EVM.
func NewEVMContext(msg Message, header *types.Header, chain ChainContext, author *common.Address) vm.Context {
	// If we don't have an explicit author (i.e. not mining), extract from the header
	var beneficiary common.Address
	if author == nil {
		//todo
		//beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}
	return vm.Context{
		CanTransfer: vm.CanTransfer,
		Transfer:    vm.Transfer,
		GetHash:     GetHashFn(header, chain),
		From:        msg.From,
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).Set(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		GasLimit:    header.GasLimit,
		GasPrice:    new(big.Int).Set(msg.GasPrice),
	}
}


// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation bool, homestead bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation && homestead {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/params.TxDataNonZeroGas < nz {
			return 0, ErrOutOfGas
		}
		gas += nz * params.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, ErrOutOfGas
		}
		gas += z * params.TxDataZeroGas
	}
	return gas, nil
}


// to returns the recipient of the message.
func To(cfg *execConfig) common.Address {
	if cfg.msg.To == nil /* contract creation */ {
		return common.Address{}
	}
	return *cfg.msg.To
}


// TransactionExec will execute the transaction
//return receipts, gas, error
func TransactionExec(exeCfg *execConfig,msg Message) (ret []byte, usedGas uint64, failed bool, err error){

	if exeCfg != nil {
		return nil, 0, true,nil
	}
	if exeCfg.gp == nil {
		//cfg = new(Config)
		return nil, 0, true,nil
	}
	if preCheck(exeCfg) != nil {
		return nil, 0, true,nil
	}
	//
	if exeCfg.statedb == nil {
		//exeCfg.statedb, _ = state.New(common.Hash{}, state.NewDatabase(zdb.NewMemDatabase()))
		return nil, 0, true,nil
	}

	//homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	contractCreation := msg.To == nil

	// Pay intrinsic gas
	gas, err := IntrinsicGas(exeCfg.data, contractCreation, false)
	if err != nil {
		return nil, 0,true, err
	}
	if err = useGas(exeCfg,gas); err != nil {
		return nil, 0,true, err
	}

	var (
		vm = exeCfg.vm
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
	)
	//to address is nil
	if contractCreation {
		//
		ret, _, exeCfg.gas, vmerr = vm.Create(msg.From, msg.AssetId, exeCfg.data, msg.Gas, msg.Value)
	} else {
		// Increment the nonce for the next transaction
		//vm.Asset.SetNonce(msg.From, exeCfg.asset.GetNonce( msg.From ) + 1 )
		ret, exeCfg.gas, vmerr = vm.Call(msg.From, To(exeCfg), msg.AssetId,msg.Data, msg.Gas,msg.Value)
	}

	if vmerr != nil {
		log.Debug("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == ErrInsufficientBalance {
			return nil, 0,true, vmerr
		}
	}
	refundGas(exeCfg)
	exeCfg.usedGas += gasUsed(exeCfg)

	return ret, gas, false,err
}

func refundGas(cfg *execConfig) {
	// Apply refund counter, capped to half of the used gas.
	refund := gasUsed(cfg) / 2
	if refund > cfg.statedb.GetRefund() {
		refund = cfg.statedb.GetRefund()
	}
	cfg.gas += refund

	// Return for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(cfg.gas), cfg.gasPrice)
	cfg.asset.AddBalance(cfg.msg.From,*cfg.msg.AssetId ,remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	cfg.gp.AddGas(cfg.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func gasUsed(cfg *execConfig) uint64 {
	return cfg.initialGas - cfg.gas
}


func  useGas(cfg *execConfig,amount uint64) error {
	if cfg.gas < amount {
		return vm.ErrOutOfGas
	}
	cfg.gas -= amount

	return nil
}

func buyGas(cfg *execConfig) error {
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(cfg.msg.Gas), cfg.gasPrice)
	value := cfg.asset.GetBalance(cfg.msg.From,*cfg.msg.AssetId)

	if value.(big.Int).Cmp(mgval) < 0 {
		return errors.New("insufficient balance to pay for gas")
	}
	if err := cfg.gp.SubGas(cfg.msg.Gas); err != nil {
		return err
	}
	cfg.gas += cfg.msg.Gas

	cfg.initialGas = cfg.msg.Gas
	cfg.asset.SubBalance(cfg.msg.From,*cfg.msg.AssetId ,mgval)
	return nil
}

func preCheck(cfg *execConfig) error {
	// Make sure this transaction's nonce is correct.
	if cfg.msg.CheckNonce {
		nonce := cfg.asset.GetNonce(cfg.msg.From)

		if nonce < cfg.msg.Nonce {
			return errors.New("nonce too high")
		} else if nonce > cfg.msg.Nonce {
			return errors.New("nonce too low")
		}
	}
	return buyGas(cfg)
}