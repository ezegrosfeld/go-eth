package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://192.168.0.233:7545")
	if err != nil {
		log.Fatalf("Couldn't connect to client: %v", err)
	}

	conn, err := deploy(client)
	if err != nil {
		log.Fatalf("Couldn't deploy: %v", err)
	}

	owner, err := getAuth(client, "7c75aa52379eda9d434d493a58aa1e14a2d210bff2adaf257ef3df290ecc47cc")
	if err != nil {
		log.Fatalf("Couldn't get owner: %v", err)
	}

	voter1, err := getAuth(client, "33a39f7a73ad8fdaf97ea341992ab35bcefc12be6b8cd00fdbb4d4a1efb583e7")
	if err != nil {
		log.Fatalf("Couldn't get voter1: %v", err)
	}

	voter2, err := getAuth(client, "7253885c5f42db65ff701630f908d3f7d95c63ab1f7711d26cca65161240ae50")
	if err != nil {
		log.Fatalf("Couldn't get voter2: %v", err)
	}

	voter3, err := getAuth(client, "ab7492e03bc108ff8dc1956b1a05efae73113eb086e201997bbe5ac6d8b52bbc")
	if err != nil {
		log.Fatalf("Couldn't get voter3: %v", err)
	}

	_, err = conn.Vote(owner, common.Big0)
	if err != nil {
		log.Fatalf("Couldn't give right to vote: %v", err)
	}

	_, err = conn.Vote(voter1, common.Big1)
	if err != nil {
		log.Fatalf("Couldn't give right to vote: %v", err)
	}

	_, err = conn.Vote(voter2, common.Big1)
	if err != nil {
		log.Fatalf("Couldn't give right to vote: %v", err)
	}

	_, err = conn.Vote(voter3, common.Big1)
	if err != nil {
		log.Fatalf("Couldn't give right to vote: %v", err)
	}

	p, err := conn.WinnerName(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Couldn't get winner name: %v", err)
	}

	fmt.Println("El ganador es: ", p)
}

func getAuth(client *ethclient.Client, pk string) (*bind.TransactOpts, error) {
	privatekey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}

	publickey := privatekey.Public()
	publickeyECDSA, ok := publickey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	address := crypto.PubkeyToAddress(*publickeyECDSA)

	nonce, err := client.PendingNonceAt(context.TODO(), address)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.TODO())
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privatekey, chainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice        // in wei

	return auth, nil
}
