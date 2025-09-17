// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tradingstate

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/XinFinOrg/XDPoSChain/common"
	"github.com/XinFinOrg/XDPoSChain/rlp"
	"github.com/XinFinOrg/XDPoSChain/trie"
)

type DumpOrderList struct {
	Volume *big.Int
	Orders map[*big.Int]*big.Int
}
type DumpLendingBook struct {
	Volume       *big.Int
	LendingBooks map[common.Hash]DumpOrderList
}

type DumpOrderBookInfo struct {
	LastPrice              *big.Int
	LendingCount           *big.Int
	MediumPrice            *big.Int
	MediumPriceBeforeEpoch *big.Int
	Nonce                  uint64
	TotalQuantity          *big.Int
	BestAsk                *big.Int
	BestBid                *big.Int
	LowestLiquidationPrice *big.Int
}

func (t *TradingStateDB) DumpAskTrie(orderBook common.Hash) (map[*big.Int]DumpOrderList, error) {
	exhangeObject := t.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("not found orderBook: %v", orderBook.Hex())
	}
	mapResult := map[*big.Int]DumpOrderList{}
	it := trie.NewIterator(exhangeObject.getAsksTrie(t.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(it.Key)
		if priceHash.IsZero() {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.stateAskObjects[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("fail when decode order iist orderBook: %v , price :%v", orderBook.Hex(), price)
			}
			stateOrderList := newStateOrderList(t, Ask, orderBook, priceHash, data, nil)
			mapResult[price] = stateOrderList.DumpOrderList(t.db)
		}
	}
	for priceHash, stateOrderList := range exhangeObject.stateAskObjects {
		if stateOrderList.Volume().Sign() > 0 {
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = stateOrderList.DumpOrderList(t.db)
		}
	}
	listPrice := []*big.Int{}
	for price := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]DumpOrderList{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return result, nil
}

func (t *TradingStateDB) DumpBidTrie(orderBook common.Hash) (map[*big.Int]DumpOrderList, error) {
	exhangeObject := t.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("not found orderBook: %v", orderBook.Hex())
	}
	mapResult := map[*big.Int]DumpOrderList{}
	it := trie.NewIterator(exhangeObject.getBidsTrie(t.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(it.Key)
		if priceHash.IsZero() {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.stateBidObjects[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("fail when decode order iist orderBook: %v , price :%v", orderBook.Hex(), price)
			}
			stateOrderList := newStateOrderList(t, Bid, orderBook, priceHash, data, nil)
			mapResult[price] = stateOrderList.DumpOrderList(t.db)
		}
	}
	for priceHash, stateOrderList := range exhangeObject.stateBidObjects {
		if stateOrderList.Volume().Sign() > 0 {
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = stateOrderList.DumpOrderList(t.db)
		}
	}
	listPrice := []*big.Int{}
	for price := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]DumpOrderList{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return mapResult, nil
}

func (t *TradingStateDB) GetBids(orderBook common.Hash) (map[*big.Int]*big.Int, error) {
	exhangeObject := t.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("not found orderBook: %v", orderBook.Hex())
	}
	mapResult := map[*big.Int]*big.Int{}
	it := trie.NewIterator(exhangeObject.getBidsTrie(t.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(it.Key)
		if priceHash.IsZero() {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.stateBidObjects[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("fail when decode order iist orderBook: %v , price :%v", orderBook.Hex(), price)
			}
			stateOrderList := newStateOrderList(t, Bid, orderBook, priceHash, data, nil)
			mapResult[price] = stateOrderList.data.Volume
		}
	}
	for priceHash, stateOrderList := range exhangeObject.stateBidObjects {
		if stateOrderList.Volume().Sign() > 0 {
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = stateOrderList.data.Volume
		}
	}
	listPrice := []*big.Int{}
	for price := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]*big.Int{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return mapResult, nil
}

func (t *TradingStateDB) GetAsks(orderBook common.Hash) (map[*big.Int]*big.Int, error) {
	exhangeObject := t.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("not found orderBook: %v", orderBook.Hex())
	}
	mapResult := map[*big.Int]*big.Int{}
	it := trie.NewIterator(exhangeObject.getAsksTrie(t.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(it.Key)
		if priceHash.IsZero() {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.stateAskObjects[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("fail when decode order iist orderBook: %v , price : %v", orderBook.Hex(), price)
			}
			stateOrderList := newStateOrderList(t, Ask, orderBook, priceHash, data, nil)
			mapResult[price] = stateOrderList.data.Volume
		}
	}
	for priceHash, stateOrderList := range exhangeObject.stateAskObjects {
		if stateOrderList.Volume().Sign() > 0 {
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = stateOrderList.data.Volume
		}
	}
	listPrice := []*big.Int{}
	for price := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]*big.Int{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return result, nil
}

func (s *stateOrderList) DumpOrderList(db Database) DumpOrderList {
	mapResult := DumpOrderList{Volume: s.Volume(), Orders: map[*big.Int]*big.Int{}}
	orderListIt := trie.NewIterator(s.getTrie(db).NodeIterator(nil))
	for orderListIt.Next() {
		keyHash := common.BytesToHash(orderListIt.Key)
		if keyHash.IsZero() {
			continue
		}
		if _, exist := s.cachedStorage[keyHash]; exist {
			continue
		} else {
			_, content, _, _ := rlp.Split(orderListIt.Value)
			mapResult.Orders[new(big.Int).SetBytes(keyHash.Bytes())] = new(big.Int).SetBytes(content)
		}
	}
	for key, value := range s.cachedStorage {
		if !value.IsZero() {
			mapResult.Orders[new(big.Int).SetBytes(key.Bytes())] = new(big.Int).SetBytes(value.Bytes())
		}
	}
	listIds := []*big.Int{}
	for id := range mapResult.Orders {
		listIds = append(listIds, id)
	}
	sort.Slice(listIds, func(i, j int) bool {
		return listIds[i].Cmp(listIds[j]) < 0
	})
	result := DumpOrderList{Volume: s.Volume(), Orders: map[*big.Int]*big.Int{}}
	for _, id := range listIds {
		result.Orders[id] = mapResult.Orders[id]
	}
	return mapResult
}

func (t *TradingStateDB) DumpOrderBookInfo(orderBook common.Hash) (*DumpOrderBookInfo, error) {
	exhangeObject := t.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("not found orderBook: %v", orderBook.Hex())
	}
	result := &DumpOrderBookInfo{}
	result.LastPrice = exhangeObject.data.LastPrice
	result.LendingCount = exhangeObject.data.LendingCount
	result.MediumPrice = exhangeObject.data.MediumPrice
	result.MediumPriceBeforeEpoch = exhangeObject.data.MediumPriceBeforeEpoch
	result.Nonce = exhangeObject.data.Nonce
	result.TotalQuantity = exhangeObject.data.TotalQuantity
	result.BestAsk = new(big.Int).SetBytes(exhangeObject.getBestPriceAsksTrie(t.db).Bytes())
	result.BestBid = new(big.Int).SetBytes(exhangeObject.getBestBidsTrie(t.db).Bytes())
	lowestPrice, _ := exhangeObject.getLowestLiquidationPrice(t.db)
	result.LowestLiquidationPrice = new(big.Int).SetBytes(lowestPrice.Bytes())
	return result, nil
}

func (s *stateLendingBook) DumpOrderList(db Database) DumpOrderList {
	mapResult := DumpOrderList{Volume: s.Volume(), Orders: map[*big.Int]*big.Int{}}
	orderListIt := trie.NewIterator(s.getTrie(db).NodeIterator(nil))
	for orderListIt.Next() {
		keyHash := common.BytesToHash(orderListIt.Key)
		if keyHash.IsZero() {
			continue
		}
		if _, exist := s.cachedStorage[keyHash]; exist {
			continue
		} else {
			_, content, _, _ := rlp.Split(orderListIt.Value)
			mapResult.Orders[new(big.Int).SetBytes(keyHash.Bytes())] = new(big.Int).SetBytes(content)
		}
	}
	for key, value := range s.cachedStorage {
		if !value.IsZero() {
			mapResult.Orders[new(big.Int).SetBytes(key.Bytes())] = new(big.Int).SetBytes(value.Bytes())
		}
	}
	listIds := []*big.Int{}
	for id := range mapResult.Orders {
		listIds = append(listIds, id)
	}
	sort.Slice(listIds, func(i, j int) bool {
		return listIds[i].Cmp(listIds[j]) < 0
	})
	result := DumpOrderList{Volume: s.Volume(), Orders: map[*big.Int]*big.Int{}}
	for _, id := range listIds {
		result.Orders[id] = mapResult.Orders[id]
	}
	return mapResult
}

func (s *liquidationPriceState) DumpLendingBook(db Database) (DumpLendingBook, error) {
	result := DumpLendingBook{Volume: s.Volume(), LendingBooks: map[common.Hash]DumpOrderList{}}
	it := trie.NewIterator(s.getTrie(db).NodeIterator(nil))
	for it.Next() {
		lendingBook := common.BytesToHash(it.Key)
		if lendingBook.IsZero() {
			continue
		}
		if _, exist := s.stateLendingBooks[lendingBook]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return result, fmt.Errorf("failed to decode state lending book orderbook: %s, liquidation price: %s , lendingBook: %s , err: %v", s.orderBook, s.liquidationPrice, lendingBook, err)
			}
			stateLendingBook := newStateLendingBook(s.orderBook, s.liquidationPrice, lendingBook, data, nil)
			result.LendingBooks[lendingBook] = stateLendingBook.DumpOrderList(db)
		}
	}
	for lendingBook, stateLendingBook := range s.stateLendingBooks {
		if !lendingBook.IsZero() {
			result.LendingBooks[lendingBook] = stateLendingBook.DumpOrderList(db)
		}
	}
	return result, nil
}

func (t *TradingStateDB) DumpLiquidationPriceTrie(orderBook common.Hash) (map[*big.Int]DumpLendingBook, error) {
	exhangeObject := t.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("not found orderBook: %v", orderBook.Hex())
	}
	mapResult := map[*big.Int]DumpLendingBook{}
	it := trie.NewIterator(exhangeObject.getLiquidationPriceTrie(t.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(it.Key)
		if priceHash.IsZero() {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.liquidationPriceStates[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("fail when decode order iist orderBook: %v , price : %v", orderBook.Hex(), price)
			}
			liquidationPriceState := newLiquidationPriceState(t, orderBook, priceHash, data, nil)
			dumpLendingBook, err := liquidationPriceState.DumpLendingBook(t.db)
			if err != nil {
				return nil, err
			}
			mapResult[price] = dumpLendingBook
		}
	}
	for priceHash, liquidationPriceState := range exhangeObject.liquidationPriceStates {
		if liquidationPriceState.Volume().Sign() > 0 {
			dumpLendingBook, err := liquidationPriceState.DumpLendingBook(t.db)
			if err != nil {
				return nil, err
			}
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = dumpLendingBook
		}
	}
	listPrice := []*big.Int{}
	for price := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]DumpLendingBook{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return mapResult, nil
}
