package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tonkeeper/tonapi-go"
)

func main() {
	const addr = "0QANsjLvOX2MERlT4oyv2bSPEVc9lunSPIs5a1kPthCXydUX"
	client, err := tonapi.NewClient(tonapi.TestnetTonApiURL, &tonapi.Security{})
	if err != nil {
		log.Fatal(err)
	}

	// Get account information
	balance, err := client.GetAccountJettonsBalances(context.Background(), tonapi.GetAccountJettonsBalancesParams{
		AccountID: addr,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range balance.Balances {
		fmt.Printf("Balance: %s %s\n", b.Balance, b.Jetton.Name)
	}
}
