package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	database_url = "mongodb://root:example@localhost:27017"
)

func main() {
	orderIDStr := "6823dc5bcb80d8ea88f9b32b"
	orderID, err := bson.ObjectIDFromHex(orderIDStr)
	if err != nil {
		fmt.Printf("Invalid ObjectID string: %v", err)
	}
	fmt.Printf("orderID: %v", orderID)

	orderIDStr = orderID.Hex()
	fmt.Printf("orderIDStr: %v, type: %T", orderIDStr, orderIDStr)

	orderIDStr2 := orderID.String()
	fmt.Printf("orderIDStr2: %v, type: %T", orderIDStr2, orderIDStr2)
}
