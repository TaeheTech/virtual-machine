package statedb
import (

	"github.com/vm-project/common"

)

// Database wraps all database operations. All methods are safe for concurrent use.
type db interface {
	CreateAccount(common.Address)

	SubBalance(  common.Address,  common.Address,  interface{})
	AddBalance(  common.Address,  common.Address,  interface{})
	GetBalance( common.Address,  common.Address) interface{}
	//GetBalance(common.Address, int, *big.Int)

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	GetRefund() uint64

	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	//AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)
}