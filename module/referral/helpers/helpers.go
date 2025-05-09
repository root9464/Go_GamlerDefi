package referral_helper

import (
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type IReferralHelper interface {
	CellTransferJettonsFromPlatform(dict []JettonEntry) (*cell.Cell, error)
	CellTransferJettonsFromLeader(dict []JettonEntry, amountJettons float64) (*cell.Cell, error)
}

type ReferralHelper struct {
	logger                 *logger.Logger
	smart_contract_address string
}

func NewReferralHelper(logger *logger.Logger, smart_contract_address string) IReferralHelper {
	return &ReferralHelper{
		logger:                 logger,
		smart_contract_address: smart_contract_address,
	}
}
