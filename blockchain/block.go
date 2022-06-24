package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// Create block
func CreateBlock(data string, prevHash []byte) *Block {
	//
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Generate the genesis block
func Genesis() *Block {
	return CreateBlock("Genesis block", []byte{})
}

// Serialize the block because badgerDB only accept slice or arrays of bytes
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	Catch(err)

	return result.Bytes()
}

// Deserialize the created slice/array of bytes and transform it to a block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Catch(err)

	return &block
}

func Catch(err error) {
	if err != nil {
		log.Panic(err)
	}
}
