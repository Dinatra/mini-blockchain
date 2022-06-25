package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/kleimak/mini-blockchain/logger"
)

type Transaction struct {
	ID      []byte // it will be a hash
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int    // Value in Tokens locked here
	PubKey string // Users address (Public key) needed to unlock the tokens
}

// Reference of output e.g [Txn] -> hash : xxx | index : 10
type TxInput struct {
	ID  []byte // references the transaction that the output is inside of
	Out int    // Index where the output appears
	Sig string // Script who provides the used data in the output pubKey
}

func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	// input has empty slice & negative index because because it is the first transaction
	txIn := TxInput{[]byte{}, -1, data}
	// First args is the reward (100 tokens) and second is the address of the receiver
	txOut := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetID()

	return &tx
}

// Create ID Based on the txn data bytes
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	logger.Catch(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Determine if a txn is a coinbase txn
func (tx *Transaction) IsCoinbase() bool {
	// check if the length of the inputs is 1 because the coinbase only has one input
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// If true its means its own the data inside the output which is referenced by the input
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// If true its means the account which has the "data" owns the information
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func NewTransaction(from string, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput   // contain all the inputs for the transaction
	var outputs []TxOutput // contain all the outputs for the transaction

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		fmt.Printf("User %v has only %v tokens\n", from, acc) // account doesnt have enough tokens
		log.Panic("ERROR: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		logger.Catch(err)

		for _, out := range outs {
			input := TxInput{txID, out, from} // reference the outputs into the inputs (creating input for all the unspent output of the sender)
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to}) // create the output for the receiver

	if acc > amount { // check if the amount that he want to send is not greater than the amount he has
		outputs = append(outputs, TxOutput{acc - amount, from}) // create the output for the sender
	}
	// Finalize the transaction
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
