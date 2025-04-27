package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type urlBuilder struct {
	baseURL string
	startID int
	endID   int
}

func (u *urlBuilder) urlFromEnv() error {
	u.baseURL = os.Getenv("BASE_URL")
	startID, err := strconv.Atoi(os.Getenv("START_ID"))
	if err != nil {
		return fmt.Errorf("func urlFromEnv: failed to load config: %w", err)
	}
	endID, err := strconv.Atoi(os.Getenv("END_ID"))
	if err != nil {
		return fmt.Errorf("func urlFromEnv: failed to load config: %w", err)
	}
	u.startID = startID
	u.endID = endID

	return nil
}

func (u *urlBuilder) buildURL(id int) string {
	var url strings.Builder
	stringID := strconv.Itoa(id)

	url.WriteString(u.baseURL)
	url.WriteString(stringID)

	return url.String()
}
