package referral_service

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
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
