package main

import (
	"fmt"
	"strconv"

	"github.com/kleimak/mini-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("Init the first block")
	chain.AddBlock("Init the secondary block")
	chain.AddBlock("Init the third block")
	chain.AddBlock("Init the forth block")
	chain.AddBlock("Init the fifth block")

	for i, block := range chain.Blocks {
		fmt.Printf("Block[%x] => data : %x \n", i, string(block.Data))
		fmt.Printf("Block[%x] => previous hash : %x \n", i, block.PrevHash)
		fmt.Printf("Block[%x] => hash : %x \n", i, block.Hash)
		proof := blockchain.NewProof(block)
		fmt.Printf("Proof[%v] => %v\n", i, strconv.FormatBool(proof.Validate()))
		fmt.Printf("------------------------ \n")
	}
}
