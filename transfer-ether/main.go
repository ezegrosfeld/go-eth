package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to blockchain node
	conn, err := ethclient.Dial("http://192.168.0.233:7545")
	if err != nil {
		log.Fatalf("Couldn't connect to client: %v", err)
	}

	// Get private key
	pk, err := crypto.HexToECDSA("568ec0233599ade0f6fdab22b4a034476861ab8a4e22c343f9d090b74c00460d")
	if err != nil {
		log.Fatalf("Couldn't get private key: %v", err)
	}

	// Get public key and cast it
	pubKey := pk.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("error casting public key to ECDSA")
	}

	// Get address from public key
	fromAddress := crypto.PubkeyToAddress(*pubKeyECDSA)

	// Get nonce
	nonce, err := conn.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Couldn't get nonce: %v", err)
	}

	// Create transaction
	// We are sending 1 ether
	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000) // in units
	gasPrice, err := conn.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Couldn't get gas price: %v", err)
	}

	toAddress := common.HexToAddress("0xf865f3026d1b71aA35633B6D8014e1915CF92bfd")
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, pk)
	if err != nil {
		log.Fatalf("Couldn't sign transaction: %v", err)
	}

	// Send transaction
	err = conn.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Couldn't send transaction: %v", err)
	}
}
