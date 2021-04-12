package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

// Define salt size
const saltSize = 16

type BCrypt interface {
	generateRandomSalt(saltSize int)
	HashPassword(password string) string
	DoPasswordsMatch(hashedPassword string, currPassword string) bool
}

type bCrypt struct {
	Salt []byte
}

func NewBCrypt(saltSize int) BCrypt {
	bc := bCrypt{}
	bc.generateRandomSalt(saltSize)
	return &bCrypt{}
}

// Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func (bc *bCrypt) generateRandomSalt(saltSize int) {
	salt := make([]byte, saltSize)

	_, err := rand.Read(salt[:])
	if err != nil {
		panic(err)
	}

	bc.Salt = salt
}

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a base64 encoded string
func (bc *bCrypt) HashPassword(password string) string {
	// Convert password string to byte slice
	passwordBytes := []byte(password)

	// Create sha-512 hasher
	sha512Hasher := sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, bc.Salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	hashedPasswordBytes := sha512Hasher.Sum(nil)

	// Convert the hashed password to a base64 encoded string
	base64EncodedPasswordHash := base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

// Check if two passwords match
func (bc *bCrypt) DoPasswordsMatch(hashedPassword string, currPassword string) bool {
	currPasswordHash := bc.HashPassword(currPassword)

	return hashedPassword == currPasswordHash
}
