package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func evaluatePolynomial(coefficients []*big.Int, x *big.Int, prime *big.Int) *big.Int {
	result := new(big.Int).Set(coefficients[0])
	xPower := new(big.Int).SetInt64(1)
	for i := 1; i < len(coefficients); i++ {
		xPower.Mul(xPower, x)
		xPower.Mod(xPower, prime)
		term := new(big.Int).Mul(coefficients[i], xPower)
		result.Add(result, term)
		result.Mod(result, prime)
	}
	return result
}

func generateShares(coefficients []*big.Int, totalShares int, prime *big.Int) [][2]*big.Int {
	shares := make([][2]*big.Int, totalShares)
	for i := 1; i <= totalShares; i++ {
		x := big.NewInt(int64(i))
		y := evaluatePolynomial(coefficients, x, prime)
		shares[i-1] = [2]*big.Int{x, y}
	}
	return shares
}

// Encrypt a file using AES
func encryptFile(filename string, key []byte) ([]byte, []byte, error) {
	// plaintext is the content of the file
	plaintext, _ := ioutil.Re-+--adFile(filename)
	// cipherblock
	block, _ := aes.NewCipher(key)
	// nonce same saize as block size
	iv := make([]byte, aes.BlockSize)
	// stream encrypts with CFB - iv & block must be same size
	stream := cipher.NewCFBEncrypter(block, iv)
	// ciphertext is created by the length of plaintext
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, iv, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Generate a random 16-byte AES key
	key := make([]byte, 16)
	rand.Read(key)
	// Print the AES key
	fmt.Println("Generated AES key (hex):", hex.EncodeToString(key))
	// Convert the key to a big.Int
	secret := new(big.Int).SetBytes(key)

	// Enter the total number of shares
	fmt.Print("Enter the total number of shares: ")
	totalSharesString, _ := reader.ReadString('\n')
	totalSharesString = strings.TrimSpace(totalSharesString)
	totalShares, _ := strconv.Atoi(totalSharesString)

	// Enter the prime number
	fmt.Print("Enter the prime number: ")
	primeStr, _ := reader.ReadString('\n')
	primeStr = strings.TrimSpace(primeStr)
	prime, _ := new(big.Int).SetString(primeStr, 10)

	// Enter the t
	fmt.Print("Enter the t: ")
	tStr, _ := reader.ReadString('\n')
	tStr = strings.TrimSpace(tStr)
	t, _ := strconv.Atoi(tStr)

	// Define the polynomial coefficients with the secret as the constant term
	coefficients := []*big.Int{secret}

	// Prompt user for (t) values in the coefficients
	for i := 1; i <= t; i++ {
		fmt.Printf("Enter coefficient %d: ", i)
		coeffStr, _ := reader.ReadString('\n')
		coeffStr = strings.TrimSpace(coeffStr)
		coeff, _ := new(big.Int).SetString(coeffStr, 10)
		coefficients = append(coefficients, coeff)
	}

	// Generate shares
	shares := generateShares(coefficients, totalShares, prime)
	fmt.Println("Shares:")
	for _, share := range shares {
		fmt.Printf("x: %s, y: %s\n", share[0], share[1])
	}

	// Ensure the shares directory exists
	sharesDir := "shares"
	if _, err := os.Stat(sharesDir); os.IsNotExist(err) {
		err := os.Mkdir(sharesDir, 0755)
		if err != nil {
			fmt.Println("Error creating shares directory:", err)
			return
		}
	}

	// Save each share to its own file in the shares directory
	for _, share := range shares {
		shareFilename := fmt.Sprintf("%s/share%s.txt", sharesDir, share[0])
		shareFile, _ := os.Create(shareFilename)
		defer shareFile.Close()
		shareFile.WriteString(fmt.Sprintf("%s %s\n", share[0], share[1]))
	}

	// Encrypt the file using the generated AES key
	filename := "secret.txt"
	ciphertext, iv, _ := encryptFile(filename, key)

	// Save the ciphertext and IV to files
	ioutil.WriteFile("encrypted.dat", ciphertext, 0644)

	ioutil.WriteFile("iv.dat", iv, 0644)

	fmt.Println("File encryption complete. Shares and encrypted file saved.")
}
