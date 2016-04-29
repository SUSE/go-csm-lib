package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strings"
)

func NormalizeGuid(id string) string {
	name := "d" + strings.Replace(id, "-", "", -1)
	name = strings.Replace(name, "`", "", -1)
	name = strings.Replace(name, ";", "", -1)
	if len(name) > 64 {
		name = name[:64]
	}

	return name
}
func GetMD5Hash(text string, size int) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")
	generated = reg.ReplaceAllString(generated, "")
	if len(generated) > size {
		generated = generated[:size]
	}
	return generated, nil
}

func SecureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(rb), nil
}
