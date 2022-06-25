# Prototype of Mini blockchain with golang

This is a prototype project whose goal is to implement the basic mechanisms of a blockchain using go language

## Usage

```
go run main.go

Usage:
        getbalance -address ADDRESS : get the balance for an address
        createblockchain -address ADDRESS : create a blockchain and send genesis block reward to address
        send -from FROM -to TO -amount AMOUNT : send amount of coins from FROM address to TO address
        printchain : Prints the blocks in the chain
```