package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"strconv"
)

// Lagrange interpolation to reconstruct the secret
func lagrangeInterpolation(x, y []*big.Int, prime *big.Int) *big.Int {
	secret := big.NewInt(0)
	for j := 0; j < len(y); j++ {
		num := big.NewInt(1)
		denom := big.NewInt(1)
		for m := 0; m < len(x); m++ {
			if m != j {
				num.Mul(num, new(big.Int).Neg(x[m]))
				num.Mod(num, prime)
				denom.Mul(denom, new(big.Int).Sub(x[j], x[m]))
				denom.Mod(denom, prime)
			}
		}
		denomInv := new(big.Int).ModInverse(denom, prime)
		term := new(big.Int).Mul(y[j], num)
		term.Mul(term, denomInv)
		term.Mod(term, prime)
		secret.Add(secret, term)
		secret.Mod(secret, prime)
	}
	return secret
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for t
	fmt.Print("Enter the t: ")
	tStr, _ := reader.ReadString('\n')
	tStr = strings.TrimSpace(tStr)
	t, _ := strconv.Atoi(tStr)

	// Prompt user for t number of points
	var points [][2]*big.Int
	for i := 0; i <= t; i++ {
		fmt.Printf("Enter point %d (x y): ", i+1)
		pointStr, _ := reader.ReadString('\n')
		pointStr = strings.TrimSpace(pointStr)
		pointParts := strings.Split(pointStr, " ")
		x, _ := new(big.Int).SetString(pointParts[0], 10)
		y, _ := new(big.Int).SetString(pointParts[1], 10)
		points = append(points, [2]*big.Int{x, y})
	}

	// Prompt user for the prime number
	fmt.Print("Enter the prime number: ")
	primeStr, _ := reader.ReadString('\n')
	primeStr = strings.TrimSpace(primeStr)
	prime, _ := new(big.Int).SetString(primeStr, 10)

	// Reconstruct the secret using the selected points
	xValues := make([]*big.Int, len(points))
	yValues := make([]*big.Int, len(points))
	for i, point := range points {
		xValues[i] = point[0]
		yValues[i] = point[1]
	}
	reconstructedSecret := lagrangeInterpolation(xValues, yValues, prime)

	// Ensure the reconstructed secret is 16 bytes long
	key := reconstructedSecret.Bytes()
	if len(key) > 16 {
		key = key[len(key)-16:]
	} else if len(key) < 16 {
		key = append(make([]byte, 16-len(key)), key...)
	}

	// Print the reconstructed AES key
	fmt.Printf("Reconstructed AES key (hex): %s\n", hex.EncodeToString(key))

	// Read the ciphertext and IV from files
	ciphertext, _ := ioutil.ReadFile("encrypted.dat")
	iv, _ := ioutil.ReadFile("iv.dat")
	// create a cipher based on the reconstructedSecret 
	block, _ := aes.NewCipher(key)
	// decrypt
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	// Write the plaintext to a file
	ioutil.WriteFile("decrypted.txt", plaintext, 0644)
	fmt.Println("File decryption complete. Decrypted content saved to decrypted.txt.")
}
