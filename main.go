package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// Generate a hash for a block through the previous hash and the block data
// and then it will update the manipulated block hash
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// Create block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// Insert the block in the chain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

// Generate the genesis block
func Genesis() *Block {
	return CreateBlock("Genesis block", []byte{})
}

// Init blockchain with the genesis first block
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("Init the first block")
	chain.AddBlock("Init the secondary block")
	chain.AddBlock("Init the third block")

	for i, block := range chain.Blocks {
		fmt.Printf("Block[%v] => data : %v \n", i, string(block.Data))
		fmt.Printf("Block[%v] => previous hash : %v \n", i, block.PrevHash)
		fmt.Printf("Block[%v] => hash : %v \n", i, block.Hash)
		fmt.Printf("------------------------ \n")
	}
}
