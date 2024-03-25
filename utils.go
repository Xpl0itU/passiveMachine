package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func stringIsEmpty(s string) bool {
	return s == ""
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateEarnAppUUID() string {
	return fmt.Sprintf("sdk-node-%x", md5.Sum([]byte(randomString(32))))
}
