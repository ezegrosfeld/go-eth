package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ezegrosfeld/go-eth/event-listener/token"
)

type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func main() {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/58d336155af0427b8810f159bfafdfc6")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("Block Number: %d\n", vLog.BlockNumber)
			fmt.Printf("Index: %d\n", vLog.Index)
			fmt.Printf("Tx Hash: %s\n", vLog.TxHash.Hex())

			contractAbi, err := abi.JSON(strings.NewReader(token.TokenABI))
			if err != nil {
				log.Fatal(err)
			}

			logTransferSig := []byte("Transfer(address,address,uint256)")
			logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				fmt.Printf("Log Name: Transfer\n")

				var transferEvent LogTransfer

				err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}

				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

				valueInt := transferEvent.Value.Int64() // 1 USD = 1,000,000 wei
				valueInUSD := float64(valueInt) / 1000000.0

				fmt.Printf("From: %s\n", transferEvent.From.Hex())
				fmt.Printf("To: %s\n", transferEvent.To.Hex())
				fmt.Printf("Tokens: %v\n", transferEvent.Value)
				fmt.Printf("USD: %.6f\n", valueInUSD)

				fmt.Printf("\n\n")
			}
		}
	}
}
