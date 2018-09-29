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

package asset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/zipper-project/z0/crypto"
)

//AccountAssetInfo .
type AccountAssetInfo struct {
	Name     string
	Symbol   string
	Total    *big.Int
	Decimals uint64
	Owner    common.Address
}

func issueAccountAsset(db StateDB, accountAddr common.Address, nonce uint64, desc string) (common.Address, error) {
	var info AccountAssetInfo
	err := json.Unmarshal([]byte(desc), &info)
	if err != nil {
		return common.Address{}, err
	}

	//save asset info
	assetAddr := crypto.CreateAssetAddress(accountAddr, nonce, info.Name)
	save(db, &info, assetAddr, assetAddr.String())
	//save base type
	assetTypeKey := assetAddr.String() + string(assetType)

	save(db, uint64(AccountModel), assetAddr, assetTypeKey)
	//issue balance to owner
	setAccountList(AccountModel, db, info.Owner, assetAddr)
	key := info.Owner.String() + assetAddr.String()
	save(db, info.Total, info.Owner, key)
	return assetAddr, nil
}

func getAccountAssetInfo(db StateDB, assetAddr common.Address) (AccountAssetInfo, error) {
	var info AccountAssetInfo
	asset := db.GetAccount(assetAddr, assetAddr.String())
	if !bytes.Equal(asset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(asset), &info)
		if err != nil {
			return info, err
		}
	} else {
		return info, fmt.Errorf("Asset not exit")
	}
	return info, nil
}

func setAccountNewOwner(db StateDB, oldOwner common.Address, assetAddr common.Address, newOwner common.Address) (bool, error) {
	var info AccountAssetInfo
	asset := db.GetAccount(assetAddr, assetAddr.String())
	if !bytes.Equal(asset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(asset), &info)
		if err != nil {
			return false, err
		}
	} else {
		return false, fmt.Errorf("Asset not exit")
	}
	if strings.Compare(info.Owner.String(), oldOwner.String()) != 0 {
		return false, nil
	}
	info.Owner = newOwner

	//save asset info
	save(db, &info, assetAddr, assetAddr.String())
	return true, nil
}

func increaseAccountAsset(db StateDB, ownerAddr common.Address, assetAddr common.Address, value *big.Int) error {
	//get asset balance
	var info AccountAssetInfo
	asset := db.GetAccount(assetAddr, assetAddr.String())
	if !bytes.Equal(asset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(asset), &info)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Asset not exit")
	}
	if strings.Compare(info.Owner.String(), ownerAddr.String()) != 0 {
		return fmt.Errorf("Owner error ")
	}

	//save asset info
	info.Total = new(big.Int).Add(info.Total, value)
	save(db, &info, assetAddr, assetAddr.String())
	//save owner balance
	key := ownerAddr.String() + assetAddr.String()
	selfAsset := db.GetAccount(ownerAddr, key)
	var balance *big.Int

	if !bytes.Equal(selfAsset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(selfAsset), &balance)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Owner asset not exit")
	}
	balance = new(big.Int).Add(balance, value)
	save(db, balance, ownerAddr, key)
	return nil
}

func subAccountBalance(db StateDB, targetAddr common.Address, assetAddr common.Address, value *big.Int) error {
	key := targetAddr.String() + assetAddr.String()
	selfAsset := db.GetAccount(targetAddr, key)
	var balance *big.Int

	if !bytes.Equal(selfAsset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(selfAsset), &balance)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Asset not exit")
	}

	if balance.Cmp(value) < 0 {
		return fmt.Errorf("Asset not enough")
	}
	balance = new(big.Int).Sub(balance, value)
	save(db, balance, targetAddr, key)
	return nil
}

func addAccountlBalance(db StateDB, targetAddr common.Address, assetAddr common.Address, value *big.Int) error {
	key := targetAddr.String() + assetAddr.String()
	selfAsset := db.GetAccount(targetAddr, key)
	var balance *big.Int

	if !bytes.Equal(selfAsset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(selfAsset), &balance)
		if err != nil {
			return err
		}
	} else {
		balance = big.NewInt(0)
		setAccountList(AccountModel, db, targetAddr, assetAddr)
	}
	balance = new(big.Int).Add(balance, value)
	save(db, balance, targetAddr, key)
	return nil
}

func enoughAccountBalance(db StateDB, targetAddr common.Address, assetAddr common.Address, value *big.Int) (bool, error) {
	key := targetAddr.String() + assetAddr.String()
	selfAsset := db.GetAccount(targetAddr, key)
	var balance *big.Int

	if !bytes.Equal(selfAsset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(selfAsset), &balance)
		if err != nil {
			return false, err
		}
	} else {
		return false, fmt.Errorf("Asset not exit")
	}

	if balance.Cmp(value) < 0 {
		return false, nil
	}

	return true, nil
}

func getAccountBalance(db StateDB, targetAddr common.Address, assetAddr common.Address) (*big.Int, error) {
	key := targetAddr.String() + assetAddr.String()
	selfAsset := db.GetAccount(targetAddr, key)
	var balance *big.Int

	if !bytes.Equal(selfAsset, []byte{}) {
		err := rlp.Decode(bytes.NewReader(selfAsset), &balance)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Asset not exit")
	}

	return balance, nil
}
