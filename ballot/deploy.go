package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func deploy(conn *ethclient.Client) (*Ballot, error) {
	privatekey, err := crypto.HexToECDSA("7c75aa52379eda9d434d493a58aa1e14a2d210bff2adaf257ef3df290ecc47cc")
	if err != nil {
		return nil, err
	}

	publickey := privatekey.Public()
	publickeyECDSA, ok := publickey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	address := crypto.PubkeyToAddress(*publickeyECDSA)

	nonce, err := conn.PendingNonceAt(context.TODO(), address)
	if err != nil {
		return nil, err
	}

	chainID, err := conn.NetworkID(context.TODO())
	if err != nil {
		return nil, err
	}

	gasPrice, err := conn.SuggestGasPrice(context.TODO())
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

	proposals := []string{"Ezekiel", "Javier"}

	// Deploy the contract, get the bytecode
	addr, tx, instance, err := DeployBallot(auth, conn, proposals)
	if err != nil {
		return nil, err
	}

	// Print information
	fmt.Println("Addr: ", addr.Hex())
	fmt.Println("Hash: ", tx.Hash().Hex())
	fmt.Println("Gas: ", tx.Gas())

	return instance, err
}
