package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/kleimak/mini-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("	add -b BLOCK_DATA - add a block to the chain")
	fmt.Println("	print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// Exist the application by shutting down the go routine
		// Badger needs to properly garbage collect the values and keys before its shuts down
		// So then its avoid corrupting data
		runtime.Goexit()
	}
}

// Add new block into the blockchain
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block added successfully")
}

// Print data inside the DB if there is one
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Block => data : %x \n", string(block.Data))
		fmt.Printf("Block => previous hash : %x \n", block.PrevHash)
		fmt.Printf("Block => hash : %x \n", block.Hash)
		proof := blockchain.NewProof(block)
		fmt.Printf("Proof => %v\n", strconv.FormatBool(proof.Validate()))
		fmt.Printf("------------------------ \n")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// CLI that help to add block manualy or vizualise the blocks inside the blockchain
func (cli *CommandLine) run() {
	cli.validateArgs()

	// set all flags
	addBlockCommand := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCommand := flag.NewFlagSet("print", flag.ExitOnError)
	// add subset
	addBlockData := addBlockCommand.String("b", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCommand.Parse(os.Args[2:])
		blockchain.Catch(err)
	case "print":
		err := printChainCommand.Parse(os.Args[2:])
		blockchain.Catch(err)
	default:
		cli.printUsage()
		runtime.Goexit()
		break
	}

	if addBlockCommand.Parsed() {
		if *addBlockData == "" {
			addBlockCommand.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCommand.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}

	cli.run()
}
