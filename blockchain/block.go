package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/kleimak/mini-blockchain/logger"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

// Need to hash the transactions to be able to use the proof of work algorithm
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Create block
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	//
	block := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Generate the genesis block
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize the block because badgerDB only accept slice or arrays of bytes
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	logger.Catch(err)

	return result.Bytes()
}

// Deserialize the created slice/array of bytes and transform it to a block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	logger.Catch(err)

	return &block
}
