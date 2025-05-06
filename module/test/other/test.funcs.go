package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/address"
)

func main() {
	const addr = "kQCKW5X6AqcHY5if5QQvChBOwdvqUz_zODy2-BxHzvAtriiJ"
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
		fmt.Printf("Balance: %s %s\n", b.Balance, b.WalletAddress.Address)
	}

	v := "0:5ac4333e06ef077965bf0933b9405b3980175b6ac135966f89eca3fda784a5b2"

	rawAddr, err := address.ParseRawAddr(v)
	if err != nil {
		fmt.Printf("Ошибка при парсинге адреса: %v\n", err)
		return
	}

	// Преобразуем в user-friendly формат (bounceable, mainnet)
	userFriendlyAddr := rawAddr.String()
	fmt.Printf("User-friendly адрес: %s\n", userFriendlyAddr)
}
