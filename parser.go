package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Parser struct {
	cookie     authCookies
	url        urlBuilder
	author     string
	sleepParam SleepParam
	client     *http.Client
	logger     zerolog.Logger
}

func (p *Parser) LoadConfig() error {
	p.cookie.cookieFromEnv()
	p.author = os.Getenv("AUTHOR")

	if err := p.url.urlFromEnv(); err != nil {
		return fmt.Errorf("func LoadConfig"+
			":failed to load urlBuilder: %w", err)
	}

	if err := p.sleepParam.sleepFromEnv(); err != nil {
		return fmt.Errorf("func LoadConfig"+
			":Failed to load SleepParam: %w", err)
	}

	p.client = &http.Client{
		Timeout: 5 * time.Second,
	}

	return nil
}

func (p *Parser) createLog(file string) error {
	if !strings.HasSuffix(file, ".txt") {
		return fmt.Errorf("func createLog: unsupported extension (need .txt): %v", file)
	}
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("func createLog: Cannot create log file, %w", err)
	}

	logOutput := zerolog.ConsoleWriter{
		Out:        logFile,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}
	consoleOutput := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	multi := zerolog.MultiLevelWriter(logOutput, consoleOutput)

	p.logger = zerolog.New(multi).With().Timestamp().Logger()

	return nil
}

func (p *Parser) makeRequest(id int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, p.url.buildURL(id), nil)
	if err != nil {
		return nil, fmt.Errorf("func makeRequest: Can't create request : %w", err)
	}
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.AddCookie(&http.Cookie{Name: "a", Value: p.cookie.cookieA})
	req.AddCookie(&http.Cookie{Name: "b", Value: p.cookie.cookieB})
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("func MakeRequest: failed to get response from : %w", err)
	}

	return resp, nil
}

func (p *Parser) findAuthor(resp *http.Response) (bool, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("func findAuthor: failed to read html body: %w", err)
	}

	if strings.Contains(strings.ToLower(string(body)), strings.ToLower(p.author)) {
		return true, nil
	}

	return false, nil
}
