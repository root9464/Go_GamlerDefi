package referral_helper

import (
	"strconv"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type JettonEntry struct {
	Address *address.Address
	Amount  uint64
}

func (h *ReferralHelper) createJettonsDictionary(entries []JettonEntry) (*cell.Dictionary, error) {
	dict := cell.NewDict(267)
	for _, entry := range entries {
		valueCell := cell.BeginCell().MustStoreCoins(entry.Amount).EndCell()
		if err := dict.Set(cell.BeginCell().MustStoreAddr(entry.Address).EndCell(), valueCell); err != nil {
			return nil, err
		}
	}

	return dict, nil
}

func (h *ReferralHelper) CellTransferJettonsFromLeader(dict []JettonEntry, amountJettons float64) (*cell.Cell, error) {
	h.logger.Infof("create cell transfer jettons from leader")
	h.logger.Infof("create jettons dictionary: %v", dict)

	dictionary, err := h.createJettonsDictionary(dict)
	if err != nil {
		h.logger.Errorf("create jettons dictionary error: %s", err)
		return cell.BeginCell().EndCell(), err
	}
	h.logger.Infof("create jettons dictionary successful: %s\n", dictionary)

	payload := cell.BeginCell().
		MustStoreUInt(0xf8a7ea5, 32).
		MustStoreUInt(uint64(time.Now().Unix()), 64).
		MustStoreCoins(tlb.MustFromDecimal(strconv.FormatFloat(amountJettons, 'f', -1, 64), 9).Nano().Uint64()).
		MustStoreAddr(address.MustParseAddr(h.smartContractAddress)).
		MustStoreUInt(0, 2).
		MustStoreUInt(0, 1).
		MustStoreCoins(tlb.MustFromTON("0.1").Nano().Uint64()).
		MustStoreBoolBit(true).
		MustStoreRef(cell.BeginCell().MustStoreDict(dictionary).EndCell()).
		EndCell()

	return payload, nil
}

func (h *ReferralHelper) CellTransferJettonsFromPlatform(dict []JettonEntry) (*cell.Cell, error) {
	h.logger.Infof("create cell transfer jettons from platform")
	h.logger.Infof("create jettons dictionary: %v", dict)

	dictionary, err := h.createJettonsDictionary(dict)
	if err != nil {
		h.logger.Errorf("create jettons dictionary error: %s", err)
		return cell.BeginCell().EndCell(), err
	}
	h.logger.Infof("create jettons dictionary successful: %s\n", dictionary)

	payload := cell.BeginCell().
		MustStoreUInt(0xfba77a9, 32).
		MustStoreUInt(uint64(time.Now().Unix()), 64).
		MustStoreDict(dictionary).
		EndCell()

	return payload, nil
}
