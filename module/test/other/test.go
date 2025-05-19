package main

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/address"
)

func IsTransactionValid(tx *tonapi.Trace) bool {
	if tx.Transaction.ComputePhase.IsSet() &&
		!tx.Transaction.ComputePhase.Value.Skipped &&
		!tx.Transaction.ComputePhase.Value.Success.Value {
		return false
	}

	if tx.Transaction.ActionPhase.IsSet() &&
		(!tx.Transaction.ActionPhase.Value.Success || tx.Transaction.ActionPhase.Value.ResultCode != 0) {
		return false
	}

	return lo.EveryBy(tx.Children, func(child tonapi.Trace) bool {
		return IsTransactionValid(&child)
	})
}

func valid() {
	client, err := tonapi.NewClient(tonapi.TestnetTonApiURL, &tonapi.Security{})
	if err != nil {
		fmt.Printf("Failed to create ton api client: %v", err)
	}

	txTrace, err := client.GetTrace(context.Background(), tonapi.GetTraceParams{
		TraceID: "105f7620bf78d534941ebcf97dda0dbe8e79c134a8ab346843787c71fe3308d5",
	})
	if err != nil {
		fmt.Printf("Failed to get trace: %v", err)
	}

	fmt.Printf("%+v", txTrace.Transaction.InMsg.Value.Hash)
	isSuccess := IsTransactionValid(txTrace)
	fmt.Print(isSuccess)
}

func main() {
	const addr = "UQANsjLvOX2MERlT4oyv2bSPEVc9lunSPIs5a1kPthCXyW6d"
	rawAddr := address.MustParseAddr(addr)
	fmt.Printf("%+v", rawAddr.StringRaw())
}
