package asset
import (

	"github.com/vm-project/common"
)

// Database wraps all database operations. All methods are safe for concurrent use.
type StateDB interface {
	GetAccount(addr common.Address, key string) []byte
	SetAccount(addr common.Address, key string, value []byte)
	GetRefund() uint64
}