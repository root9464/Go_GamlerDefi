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
	logger               *logger.Logger
	smartContractAddress string
}

func NewReferralHelper(logger *logger.Logger, smartContractAddress string) IReferralHelper {
	return &ReferralHelper{
		logger:               logger,
		smartContractAddress: smartContractAddress,
	}
}
