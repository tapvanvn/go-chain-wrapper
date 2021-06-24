package utility

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"log"
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var randArray string = "abcdefghijklmnopqrstuvwxyz0123456789"
var randArrayCase string = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//GenVerifyCode generate a verify code
func GenVerifyCode(length int) string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return code
}

//GenOTPCode generate otp code
func GenOTPCode(length int) string {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return strings.ToUpper(code)
}

//GenerateSalt generate a verify code
func GenerateSalt(length int) string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return code
}

//GenSecretKey generate secret key
func GenSecretKey(length int) string {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArrayCase)
	for i := 0; i < length; i++ {
		code += string(randArrayCase[rand.Intn(arrayLen)])
	}

	return code
}

//HashPassword Generate hash from password and salt
func HashPassword(password string, salt string) string {
	pwd := []byte(salt + password + salt)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

//VerifyPasswords verify password
func VerifyPasswords(passwordShadow string, password string, salt string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(passwordShadow)
	plainPwd := []byte(salt + password + salt)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func GenJWTSecret() (string, error) {
	b := make([]byte, 265)
	_, err := crypto_rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {

		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
