package asset
import (

	"github.com/vm-project/vm"

)

// Database wraps all database operations. All methods are safe for concurrent use.
type asset interface {
	GetCodeSize(vm.Address)
	SetCode(vm.Address, []byte)
	CreateAccount(vm.Address)

}