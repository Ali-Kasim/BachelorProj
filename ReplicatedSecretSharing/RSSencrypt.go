package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
)


// Generate shares 
func generateShares(x *big.Int, p *big.Int) ([]*big.Int, []*big.Int, []*big.Int) {
	// r3 = x - r1 - r2 mod p
	r1 := generateRandom(64)
	r2 := generateRandom(64) 
	r3 := new(big.Int).Sub(x, new(big.Int).Add(r1, r2))
	r3.Mod(r3, p)
	return []*big.Int{r1}, []*big.Int{r2}, []*big.Int{r3}
}

// Generate a random prime number with 128 as bitsize
func generatePrime() *big.Int {
	prime, _ := rand.Prime(rand.Reader, 128)
	return prime
}


// Generate a random number with bit size
func generateRandom(bitSize int) *big.Int {
    // Generate a random number of given bitSize
    num, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bitSize)), nil))
    return num
}


// Encrypt a file using AES
func encryptFile(filename string, key []byte) ([]byte, []byte, error) {
	plaintext, _ := ioutil.ReadFile(filename)

	block, _ := aes.NewCipher(key)

	iv := make([]byte, aes.BlockSize)
	rand.Read(iv);

	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, iv, nil
}

func main() {
	// Generate a random 16-byte AES key
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		fmt.Println("Error generating AES key:", err)
		return
	}

	// Print the AES key
	fmt.Println("Generated AES key (hex):", hex.EncodeToString(key))

	// Define the prime p
	p := generatePrime()
	fmt.Println("prime: ", p)

	// Generate shares according to the protocol
	x := new(big.Int).SetBytes(key)
	shares1, shares2, shares3 := generateShares(x, p)

	// Print shares for each participant
	fmt.Printf("Shares for P1: r1: %s, r2: %s\n", shares1[0].String(), shares2[0].String())
	fmt.Printf("Shares for P2: r2: %s, r3: %s\n", shares2[0].String(), shares3[0].String())
	fmt.Printf("Shares for P3: r3: %s, r1: %s\n", shares3[0].String(), shares1[0].String())

	// Encrypt the file using the generated AES key
	filename := "secret.txt"
	ciphertext, iv, _ := encryptFile(filename, key)

	// Save the ciphertext and IV to files
	ioutil.WriteFile("encrypted.dat", ciphertext, 0644)
	ioutil.WriteFile("iv.dat", iv, 0644)
	fmt.Println("File encryption complete. Shares, encrypted file, and IV saved.")
}

