package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// MultiURLTime calls URLTime for every URL in URLs.
func MultiURLTime(urls []string) {
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()
			URLTime(u)
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// URLTime checks how much time it takes url to respond.
func URLTime(url string) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error: %q - %s", url, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("error: %q - bad status - %s", url, resp.Status)
		return
	}
	// Read body
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Printf("error: %q - %s", url, err)
		return
	}

	duration := time.Since(start)
	log.Printf("info: %q - %s", url, formatDuration(duration))
}

func formatDuration(d time.Duration) string {
	// Format duration as seconds with 4 decimal places
	return fmt.Sprintf("%.4f seconds", d.Seconds())
}

func server() {
	http.HandleFunc("/50", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is the response for /50")
	})

	http.HandleFunc("/100", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is the response for /100")
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is the response for /200")
	})

	// Start the server on port 8080
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting the server:", err)
		}
	}()

	fmt.Println("server running on port 8080...")
}

func main() {
	go server()

	start := time.Now()

	urls := []string{
		"http://localhost:8080/200",
		"http://localhost:8080/100",
		"http://localhost:8080/50",
	}

	MultiURLTime(urls)

	duration := time.Since(start)
	log.Printf("%d URLs in %v", len(urls), duration)
}
