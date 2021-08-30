package main

import (
	"fmt"
	"log"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	bip39 "github.com/tyler-smith/go-bip39"
)

func main() {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Entropy:", entropy)

	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("Mnemonic:", mnemonic)

	seed := bip39.NewSeed(mnemonic, "")

	fmt.Println("Seed:", seed)

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wallet:", wallet)

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Addr 1: ", account.Address.Hex())

	path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	account, err = wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Address 2: ", account.Address.Hex())
}
