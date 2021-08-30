#!bin/sh

solcjs --bin --abi --optimize --overwrite -o ./build/ ./contracts/ballot.sol

abigen --abi ./build/__contracts_ballot_sol_Ballot.abi --bin ./build/__contracts_ballot_sol_Ballot.bin --pkg main --type Ballot --out ./ballot.go