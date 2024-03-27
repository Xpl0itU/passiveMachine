package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func stringIsEmpty(s string) bool {
	return s == ""
}

func generateEarnAppUUID() string {
	return fmt.Sprintf("sdk-node-%s", strings.ReplaceAll(uuid.NewString(), "-", ""))
}

func saveToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
