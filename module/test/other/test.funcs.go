package main

import (
	"context"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	adminAddr = "0QANsjLvOX2MERlT4oyv2bSPEVc9lunSPIs5a1kPthCXydUX"
	csAddr    = "kQCKW5X6AqcHY5if5QQvChBOwdvqUz_zODy2-BxHzvAtriiJ"
)

type JettonEntry struct {
	Address *address.Address
	Amount  uint64
}

func CreateJettonsDictionary(entries []JettonEntry) *cell.Dictionary {
	dict := cell.NewDict(267)
	for _, entry := range entries {
		valueCell := cell.BeginCell().MustStoreCoins(entry.Amount).EndCell()
		if err := dict.Set(cell.BeginCell().MustStoreAddr(entry.Address).EndCell(), valueCell); err != nil {
			panic(err)
		}
	}

	return dict
}

func main() {
	client := liteclient.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(context.Background(), "https://ton-blockchain.github.io/testnet-global.config.json")
	if err != nil {
		panic(err)
	}
	entries := []JettonEntry{
		{
			Address: address.MustParseAddr("0QAYRw04JzUo1IEK6TL6vKfos66gsdN6vUFfJeA3OOjOfDPG"),
			Amount:  tlb.MustFromTON("0.2").Nano().Uint64(),
		},
		{
			Address: address.MustParseAddr("0QD-q5a1Z3kYfDBgYUcUX_MigynA5FuiNx0i5ySt37rfrFeP"),
			Amount:  tlb.MustFromTON("0.3").Nano().Uint64(),
		},
	}

	// Создаем словарь
	dictionary := CreateJettonsDictionary(entries)

	api := ton.NewAPIClient(client)
	body := cell.BeginCell().MustStoreUInt(0xf8a7ea5, 32).
		MustStoreUInt(uint64(time.Now().Unix()), 32).
		MustStoreCoins(tlb.MustFromDecimal("0.5", 9).Nano().Uint64()). // колличество жетонов
		MustStoreAddr(address.MustParseAddr(csAddr)).
		MustStoreUInt(0, 2).
		MustStoreUInt(0, 1).
		MustStoreCoins(tlb.MustFromTON("0.1").Nano().Uint64()). // а тут тон
		MustStoreBoolBit(true).
		MustStoreRef(cell.BeginCell().MustStoreDict(dictionary).EndCell()).
		EndCell()

}
