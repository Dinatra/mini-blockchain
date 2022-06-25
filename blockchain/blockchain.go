package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First transaction from the genesis block"
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

func DBexists() bool {
	// Checks if a file named MANIFEST has been created from badger
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Catch(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Catch(err)

		item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	chain := BlockChain{lastHash, db}

	return &chain
}

// Init blockchain with the genesis first block
func InitBlockChain(address string) *BlockChain {
	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Catch(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis block created successfully")
		err := txn.Set(genesis.Hash, genesis.Serialize())
		Catch(err)

		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})
	Catch(err)
	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

// Insert the block in the chain
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil { // check if its inside the map
					for _, spentOut := range spentTXOs[txID] {
						// check if the output is already spent
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					// determine how much tokens the user have in the wallet
					unspentTxs = append(unspentTxs, *tx)
				}

				if tx.IsCoinbase() == false { // avoid genesis block
					// find others outputs that are spent by this transaction
					for _, in := range tx.Inputs {
						if in.CanUnlock(address) { // check if we can unlock those outputs
							inTxID := hex.EncodeToString(in.ID)
							spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out) // add the output to the map with the index where its appears
						}
					}
				}
			}
		}

		// break if we reach the genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs // all this unspent transactions are the ones that the user have in the wallet
}

// find the non consumed outputs of the user
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput

	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// Enable us to know if a user can enable a transaction in depend of the tokens of its balance
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int) // unspentOuts is the map with the index of the output that we can use to spend
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			// Check if the accumulated is less than the amount that he want to send and if he can unlock the output
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value                              // increment the value that he want to spend
				unspentOuts[txID] = append(unspentOuts[txID], outIdx) // add it to the unspentOutputs

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
