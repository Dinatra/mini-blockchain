package blockchain

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// This struct help only to iterate through the blockchain data
type BlockChainIterator struct {
	currentHash []byte
	Database    *badger.DB
}

// Init blockchain with the genesis first block
func InitBlockChain() *BlockChain {
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Catch(err)

	err = db.Update(func(txn *badger.Txn) error {
		// At first, check if there is a blockchain stored in the db through "lh" key referencing the genesis block
		if _, err = txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No database has been found")

			genesis := Genesis()
			fmt.Println("Genesis proved")
			// The key of the genesis block is its hash and the value is the output of the serialized data of the block
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Catch(err)
			// Set the hash of genesis block as the last hash because he is at this moment the only in the db
			err = txn.Set([]byte("lh"), genesis.Hash)

			// set the genesis block hash as the lasthash variable and put it to the memory storage
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Catch(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
	})
	Catch(err)
	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

// Insert the block in the chain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		// Get the lasthash from the lastblock
		item, err := txn.Get([]byte("lh"))
		Catch(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	Catch(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		// Assign the newblock hash to the lasthash key
		// it help to easily get it out of database and derive a new block into it
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Catch(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Catch(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	// Set the chain last registered hash into LH
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

// Iterate backward Until the genesis block
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.currentHash)

		item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})

		return err
	})
	Catch(err)

	iter.currentHash = block.PrevHash

	return block
}
