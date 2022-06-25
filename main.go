package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/kleimak/mini-blockchain/blockchain"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("	getbalance -address ADDRESS : get the balance for an address")
	fmt.Println("	createblockchain -address ADDRESS : create a blockchain and send genesis block reward to address")
	fmt.Println("	send -from FROM -to TO -amount AMOUNT : send amount of coins from FROM address to TO address")
	fmt.Println("	printchain : Prints the blocks in the chain")
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

// Print data inside the DB if there is one
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

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

func (cli *CommandLine) createblockchain(address string) {
	chain := blockchain.InitBlockChain(address) // address of the person who mines the genesis block
	chain.Database.Close()
	fmt.Println("Blockchain created successfully")
}

func (cli *CommandLine) getbalance(address string) {
	chain := blockchain.ContinueBlockChain(address) // open blockchain on the address
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address) // get all the unspend transactions outputs of the address

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s' is %d\n", address, balance)
}

func (cli *CommandLine) send(from string, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Transaction sended successfully !")
}

// CLI that help to add block manualy or vizualise the blocks inside the blockchain
func (cli *CommandLine) run() {
	cli.validateArgs()

	// set all flags
	createBlockchainCommand := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCommand := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCommand := flag.NewFlagSet("printchain", flag.ExitOnError)

	// set subset flags
	createBlockchainAddress := createBlockchainCommand.String("address", "", "The address to get balance for")
	sendFrom := sendCommand.String("from", "", "Origin wallet address")
	sendTo := sendCommand.String("to", "", "Receiver wallet address")
	sendAmount := sendCommand.Int("amount", 0, "Amount to send")
	getBalanceAddress := getBalanceCommand.String("address", "", "The address to get balance for")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCommand.Parse(os.Args[2:])
		panicLog(err)
	case "send":
		err := sendCommand.Parse(os.Args[2:])
		panicLog(err)
	case "getbalance":
		err := getBalanceCommand.Parse(os.Args[2:])
		panicLog(err)
	case "printchain":
		err := printChainCommand.Parse(os.Args[2:])
		panicLog(err)
	default:
		cli.printUsage()
		runtime.Goexit()
		break
	}
	// create blockchain
	if createBlockchainCommand.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCommand.Usage()
			runtime.Goexit()
		}
		cli.createblockchain(*createBlockchainAddress)
	}
	// send command
	if sendCommand.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCommand.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	// get balance of the address
	if getBalanceCommand.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCommand.Usage()
			runtime.Goexit()
		}
		cli.getbalance(*getBalanceAddress)
	}
	// print the blockchain
	if printChainCommand.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)

	cli := CommandLine{}

	cli.run()
}

func panicLog(err error) {
	if err != nil {
		log.Panic(err)
	}
}
