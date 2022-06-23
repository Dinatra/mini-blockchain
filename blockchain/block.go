package blockchain

type BlockChain struct {
	Blocks []*Block
}

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
