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
// along with the z0 library. If not, see <http://www.gnu.org/licenses/>.s

package statedb

import (
	"sync"
    "github.com/vm-project/common"
	"github.com/vm-project/dep/crypto"
)
//type commAddress [AddressLength]byte
type revision struct {
	id           int
	journalIndex int
}

var (
	emptyState = crypto.Keccak256Hash(nil)

	emptyCode = crypto.Keccak256Hash(nil)
)

type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash"`
	// index of the log in the receipt
	Index uint `json:"logIndex" gencodec:"required"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}
type StateDB struct {
	//db   Database
	//trie Trie
    //
	stateObjects      map[common.Address]*stateObject

	stateObjectsDirty map[common.Address]struct{}
	//

    //
	dbErr  error
	refund uint64

	//thash, bhash common.Hash
	//txIndex      int
	//addlog
	logs    map[common.Hash][]*Log
	logSize uint
    // ?
	preimages map[common.Hash][]byte

	journal        *journal
	//validRevisions []revision
	//nextRevisionId int

	lock sync.Mutex
}

func New(root common.Hash ) (*StateDB, error) {
	//tr, err := db.OpenTrie(root)
	//if err != nil {
	//	return nil, err
	//}
	return &StateDB{
		//db:                db,
		//trie:              tr,
		stateObjects:      make(map[common.Address]*stateObject),
		//stateObjectsDirty: make(map[common.Address]struct{}),
		logs:              make(map[common.Hash][]*Log),
		preimages:         make(map[common.Hash][]byte),
		journal:           newJournal(),
	}, nil
}

func (self *StateDB) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *StateDB) Error() error {
	return self.dbErr
}

func (self *StateDB) Reset(root common.Hash) error {
	//tr, err := self.db.OpenTrie(root)
	//if err != nil {
	//	return err
	//}
	//self.trie = tr
	self.stateObjects = make(map[common.Address]*stateObject)
	self.stateObjectsDirty = make(map[common.Address]struct{})
	//self.thash = common.Hash{}
	//self.bhash = common.Hash{}
	//self.txIndex = 0
	self.logs = make(map[common.Hash][]*Log)
	self.logSize = 0
	self.preimages = make(map[common.Hash][]byte)
	//self.clearJournalAndRefund()
	return nil
}

func (self *StateDB) AddLog(log *Log) {
	//self.journal.append(addLogChange{txhash: self.thash})
	//
	//log.TxHash = self.thash
	//log.BlockHash = self.bhash
	//log.TxIndex = uint(self.txIndex)
	//log.Index = self.logSize
	//self.logs[self.thash] = append(self.logs[self.thash], log)
	self.logSize++
}

func (self *StateDB) GetLogs(hash common.Hash) []*Log {
	return self.logs[hash]
}


func (self *StateDB) Logs() []*Log {
	var logs []*Log
	for _, lgs := range self.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

func (self *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := self.preimages[hash]; !ok {
		self.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		self.preimages[hash] = pi
	}
}

func (self *StateDB) Preimages() map[common.Hash][]byte {
	return self.preimages
}

func (self *StateDB) AddRefund(gas uint64) {
	self.journal.append(refundChange{prev: self.refund})
	self.refund += gas
}

func (self *StateDB) Exist(addr common.Address) bool {
	return self.getStateObject(addr) != nil
}

func (self *StateDB) Empty(addr common.Address) bool {
	so := self.getStateObject(addr)
	return so == nil || so.empty()
}

func (self *StateDB) GetCode(addr common.Address) []byte {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code()
	}
	return nil
}
func (self *StateDB) GetCodeSize(addr common.Address) int {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return 0
	}
	if stateObject.code != nil {
		return len(stateObject.code)
	}
	//size, err := self.db.ContractCodeSize(stateObject.addrHash, common.BytesToHash(stateObject.CodeHash()))
	//if err != nil {
	//	self.setError(err)
	//}
	return 0
}
func (self *StateDB) SetCode(addr common.Address, code []byte) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetCode(crypto.Keccak256Hash(code), code)
	}
}

func (self *StateDB) GetCodeHash(addr common.Address) common.Hash {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return common.Hash{}
	}
	return common.BytesToHash(stateObject.CodeHash())
}
//get hash
func (self *StateDB) GetState(addr common.Address, key common.Hash) common.Hash {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetState( key)
	}
	return common.Hash{}
}

func (self *StateDB) SetState(addr common.Address, key, value common.Hash) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetState( key, value)
	}
}

func (self *StateDB) GetAccount(addr common.Address, key string) []byte {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetAccount(key)
	}
	return []byte{}
}

func (self *StateDB) SetAccount(addr common.Address, key string, value []byte) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetAccount(key, value)
	}
}

func (self *StateDB) DeleteAccount(addr common.Address, key string) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.DeleteAccount(key)
	}
}








func (self *StateDB) getStateObject(addr common.Address) (stateObject *stateObject) {
	// Prefer 'live' objects.
	if obj := self.stateObjects[addr]; obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}

	//enc, err := self.trie.TryGet(addr[:])
	//if len(enc) == 0 {
	//	self.setError(err)
	//	return nil
	//}
	var data Account
	//if err := rlp.DecodeBytes(enc, &data); err != nil {
	//	log.Error("Failed to decode state object", "addr", addr, "err", err)
	//	return nil
	//}
	obj := newObject(self, addr, data)
	self.setStateObject(obj)
	return obj
}

func (self *StateDB) setStateObject(object *stateObject) {
	self.stateObjects[object.Address()] = object
}

// Retrieve a state object or create a new state object if nil.
func (self *StateDB) GetOrNewStateObject(addr common.Address) *stateObject {
	stateObject := self.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject, _ = self.createObject(addr)
	}
	return stateObject
}

func (self *StateDB) createObject(addr common.Address) (newobj, prev *stateObject) {
	prev = self.getStateObject(addr)
	newobj = newObject(self, addr, Account{})
	if prev == nil {
		self.journal.append(createObjectChange{account: &addr})
	} else {
		self.journal.append(resetObjectChange{prev: prev})
	}
	self.setStateObject(newobj)
	return newobj, prev
}

func (self *StateDB) CreateAccount(addr common.Address) {
	//new, prev := self.createObject(addr)
	//if prev != nil {
	//	new.setBalance(prev.data.Balance)
	//}
	self.createObject(addr)
}

//func (db *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) {
//	so := db.getStateObject(addr)
//	if so == nil {
//		return
//	}
//
//	// When iterating over the storage check the cache first
//	for h, value := range so.cachedStorage {
//		cb(h, value)
//	}
//
//	it := trie.NewIterator(so.getTrie(db.db, so.data.StRoot, STROOTFlAG).NodeIterator(nil))
//	for it.Next() {
//		key := common.BytesToHash(db.trie.GetKey(it.Key))
//		if _, ok := so.cachedStorage[key]; !ok {
//			cb(key, common.BytesToHash(it.Value))
//		}
//	}
//}

//func (self *StateDB) Copy() *StateDB {
//	self.lock.Lock()
//	defer self.lock.Unlock()
//
//	state := &StateDB{
//		db:                self.db,
//		trie:              self.db.CopyTrie(self.trie),
//		stateObjects:      make(map[common.Address]*stateObject, len(self.journal.dirties)),
//		stateObjectsDirty: make(map[common.Address]struct{}, len(self.journal.dirties)),
//		refund:            self.refund,
//		logs:              make(map[common.Hash][]*types.Log, len(self.logs)),
//		logSize:           self.logSize,
//		preimages:         make(map[common.Hash][]byte),
//		journal:           newJournal(),
//	}
//	for addr := range self.journal.dirties {
//		if object, exist := self.stateObjects[addr]; exist {
//			state.stateObjects[addr] = object.deepCopy(state)
//			state.stateObjectsDirty[addr] = struct{}{}
//		}
//	}
//	for addr := range self.stateObjectsDirty {
//		if _, exist := state.stateObjects[addr]; !exist {
//			state.stateObjects[addr] = self.stateObjects[addr].deepCopy(state)
//			state.stateObjectsDirty[addr] = struct{}{}
//		}
//	}
//
//	for hash, logs := range self.logs {
//		state.logs[hash] = make([]*types.Log, len(logs))
//		copy(state.logs[hash], logs)
//	}
//	for hash, preimage := range self.preimages {
//		state.preimages[hash] = preimage
//	}
//	return state
//}

func (self *StateDB) Snapshot() int {
	//id := self.nextRevisionId
	//self.nextRevisionId++
	//self.validRevisions = append(self.validRevisions, revision{id, self.journal.length()})
	var id int
	return id
}

func (self *StateDB) RevertToSnapshot(revid int) {
	//idx := sort.Search(len(self.validRevisions), func(i int) bool {
	//	return self.validRevisions[i].id >= revid
	//})
	//if idx == len(self.validRevisions) || self.validRevisions[idx].id != revid {
	//	panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	//}
	//snapshot := self.validRevisions[idx].journalIndex
	//
	//self.journal.revert(self, snapshot)
	//self.validRevisions = self.validRevisions[:idx]
}

func (self *StateDB) GetRefund() uint64 {
	return self.refund
}


//func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
//	s.Finalise(deleteEmptyObjects)
//	return s.trie.Hash()
//}
//
//func (self *StateDB) Prepare(thash, bhash common.Hash, ti int) {
//	self.thash = thash
//	self.bhash = bhash
//	self.txIndex = ti
//}
//
//func (s *StateDB) clearJournalAndRefund() {
//	s.journal = newJournal()
//	s.validRevisions = s.validRevisions[:0]
//	s.refund = 0
//}


