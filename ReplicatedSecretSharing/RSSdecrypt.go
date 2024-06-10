package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

// Combine shares to reconstruct the key
func combineShares(shares1, shares2, shares3 *big.Int, p *big.Int) *big.Int {
	// reconstruction formula from rapport: 
	//x = r1 + r2 + r3 mod p
	x := new(big.Int).Add(shares1, shares2)
	x.Add(x, shares3)
	x.Mod(x, p)
	return x
}

// Decrypt a file using AES
func decryptFile(ciphertext, iv, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for their respective shares
	fmt.Println("Enter r1/share1:")
	share1Str, _ := reader.ReadString('\n')
	share1Str = strings.TrimSpace(share1Str)
	share1 := new(big.Int)
	share1.SetString(share1Str, 10)

	fmt.Println("Enter r2/share2:")
	share2Str, _ := reader.ReadString('\n')
	share2Str = strings.TrimSpace(share2Str)
	share2 := new(big.Int)
	share2.SetString(share2Str, 10)

	fmt.Println("Enter your r3/share3:")
	share3Str, _ := reader.ReadString('\n')
	share3Str = strings.TrimSpace(share3Str)
	share3 := new(big.Int)
	share3.SetString(share3Str, 10)

	// Prompt user for the prime number p
	fmt.Println("Enter the prime p:")
	pStr, _ := reader.ReadString('\n')
	pStr = strings.TrimSpace(pStr)
	p := new(big.Int)
	p.SetString(pStr, 10)

	// Read the IV from the file
	iv, err := ioutil.ReadFile("iv.dat")
	if err != nil {
		fmt.Println("Error reading IV file:", err)
		return
	}
	// Read the ciphertext from the encrypted file
	ciphertext, _ := ioutil.ReadFile("encrypted.dat")

	// Reconstruct the key using shares
	x := combineShares(share1, share2, share3, p)
	key := x.Bytes()

	// Print the reconstructed key in hexadecimal format
	fmt.Printf("Reconstructed AES Key (hex): %x\n", key)

	// Decrypt the file using the reconstructed key
	plaintext, _ := decryptFile(ciphertext, iv, key)

	// Write the decrypted plaintext to a new file
	outputFile := "decrypted_output.txt"
	ioutil.WriteFile(outputFile, plaintext, 0644)

	fmt.Println("File decryption complete. Decrypted file saved as", outputFile)
}
