package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// GetUUID ...
func GetUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return ""
	}

	return u.String()
}