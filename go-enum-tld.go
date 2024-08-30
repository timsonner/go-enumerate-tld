package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run go-enum-tld.go <domain or IP>")
		return
	}

	domain := os.Args[1]

	file, err := os.Open("tld-list.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tld := scanner.Text()
		fullURL := fmt.Sprintf("http://%s.%s", domain, tld)

		// Increment the WaitGroup counter
		wg.Add(1)

		// Launch a goroutine for each tld check
		go func(url string) {
			// Decrement the counter when the goroutine completes
			defer wg.Done()

			response, err := http.Get(url)
			if err != nil {
				// Uncomment next line for error output
				// fmt.Println("Error connecting to:", url, "-", err)
				return
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Valid tld:", url)
			} else {
				// Uncomment next line if looking for other status codes... Also, good to see how query string is constructed.
				// fmt.Println("Invalid tld or other status code:", url, "Status Code:", response.StatusCode)
			}
		}(fullURL) // Pass fullURL as an argument to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
