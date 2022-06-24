package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 10

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// Produce a proof of work
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// generate left shifted result from target starting point (1)
	// to the number of bytes inside of one of our hashes - the difficulty
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// Generate a hash trough the manipulated block data (prevHash, data, nonce, actual difficulty)
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// Algorithm that check if there is the right hash and if the block wasn't already signed
// It will return a nonce as int and byte as hash
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		// Create block hash from nonce
		data := pow.InitData(nonce)
		// Transform the [32]byte returned by initdata to a real hash
		hash = sha256.Sum256(data)
		// Convert the hash into big int
		intHash.SetBytes(hash[:])

		fmt.Printf("\r%x", hash)
		// If : the result is -1, it will mean that this hash
		// is actually less than the target it's looking for
		// this means that the block has already been signed
		// Else : increment the nonce in order to find the correct one
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Check if the pow is correct and return a bool
func (pow *ProofOfWork) Validate() bool {
	// From pow nonce create some byte[32]
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	// convert it to a hash
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	// check if it is valid
	return intHash.Cmp(pow.Target) == -1
}

// Take a int64 number and decode it into bytes
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	// bigEndian signifies how the bytes gonna be organised
	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
