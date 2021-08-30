package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to blockchain node
	conn, err := ethclient.Dial("http://192.168.0.233:7545")
	if err != nil {
		log.Fatalf("Couldn't connect to client: %v", err)
	}

	ctx := context.Background()

	// Get the balance of the address
	tx, err := conn.BalanceAt(ctx, common.HexToAddress("0xf865f3026d1b71aA35633B6D8014e1915CF92bfd"), common.Big0)
	if err != nil {
		log.Fatalf("Couldn't get balance: %v", err)
	}
	log.Printf("Balance in Wei: %v", tx)
	log.Printf("Balance in Ether: %v", tx.Div(tx, big.NewInt(1000000000000000000)))
}
