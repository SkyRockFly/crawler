package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "No flags provided, using value from config")
	}
	if err := godotenv.Load("config/settings.env"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env file:\n %v \n", err)
		os.Exit(1)
	}

	var parser Parser
	if err := parser.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create parser from config file:\n %v \n", err)
		os.Exit(1)
	}

	author := flag.String("author", parser.author, "Override author")
	file := flag.String("log", "log.txt", "Log file")
	start := flag.Int("start", parser.url.startID, "Override start ID")
	end := flag.Int("end", parser.url.endID, "Override end ID")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n🐻 Bear Hunter v1.0\n")
		fmt.Fprintf(os.Stderr, "Scan posts by author and ID range.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  parser.exe [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *author != "" {
		parser.author = *author
	}
	if *start > *end {
		fmt.Fprintln(os.Stderr, "Start is more than end")
		flag.Usage()
		os.Exit(1)
	}
	if *start < 0 || *end < 0 {
		fmt.Fprintln(os.Stderr, "Start and/or end is negative")
		flag.Usage()
		os.Exit(1)
	}

	parser.url.startID = *start
	parser.url.endID = *end

	if err := parser.createLog(*file); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log file:\n %v \n", err)
		os.Exit(1)
	}

	printBanner(*author, *start, *end, *file)

	for i := parser.url.startID; i <= parser.url.endID; i++ {
		resp, err := parser.makeRequest(i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get response:\n %v \n", err)
			os.Exit(1)
		}
		if resp.StatusCode != http.StatusOK {
			parser.logger.Info().Int("Id", i).Int("Status", resp.StatusCode).Msg("Failed to get OK status")
			continue
		}

		answ, err := parser.findAuthor(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find author:\n %v \n", err)
			os.Exit(1)
		}
		if !answ {
			parser.logger.Info().Int("Id", i).Str("Author", parser.author).
				Msg("No such")
		} else {
			parser.logger.Info().Int("Id", i).Str("Author", parser.author).
				Msg("Found ")
		}

		parser.sleepParam.sleep()
		parser.sleepParam.randomPause()
	}
	parser.logger.Info().Msg("Done searching")
}

func printBanner(author string, start, end int, logfile string) {
	fmt.Println(Green + "╔════════════════════════════╗" + Reset)
	fmt.Println(Green + "║   🐻  Bear Hunter v1.0     ║" + Reset)
	fmt.Println(Green + "╚════════════════════════════╝" + Reset)

	fmt.Printf(Yellow+"Author   : %s\n"+Reset, author)
	fmt.Printf(Yellow+"Range    : %d – %d\n"+Reset, start, end)
	fmt.Printf(Yellow+"Output   : %s\n\n"+Reset, logfile)

	colors := []string{Red, Yellow, Green}
	for i := 3; i > 0; i-- {
		fmt.Printf(colors[3-i]+"Scanning will start in %d seconds\n"+Reset, i)
		time.Sleep(1 * time.Second)
	}
}
