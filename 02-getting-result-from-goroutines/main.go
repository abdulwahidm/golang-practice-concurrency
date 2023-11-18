package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

// sha1Sig returns SHA1 signature in the format "35aabcd5a32e01d18a5ef688111624f3c547e13b"
func sha1Sig(data []byte) (string, error) {
	w := sha1.New()
	r := bytes.NewReader(data)
	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}

	sig := fmt.Sprintf("%x", w.Sum(nil))
	return sig, nil
}

type File struct {
	Name      string
	Content   []byte
	Signature string
}

// ValidateSigs returns slices of OK files and mismatched files concurrently
func ValidateSigs(files []File) ([]string, []string, error) {
	var okFiles []string
	var badFiles []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Channel to communicate results from goroutines
	resultCh := make(chan struct {
		name string
		ok   bool
	})

	// Process each file concurrently
	for _, file := range files {
		wg.Add(1)
		go func(f File) {
			defer wg.Done()

			sig, err := sha1Sig(f.Content)
			if err != nil {
				// Handle error if needed
				return
			}

			mu.Lock()
			defer mu.Unlock()

			if sig != f.Signature {
				badFiles = append(badFiles, f.Name)
				resultCh <- struct {
					name string
					ok   bool
				}{name: f.Name, ok: false}
			} else {
				okFiles = append(okFiles, f.Name)
				resultCh <- struct {
					name string
					ok   bool
				}{name: f.Name, ok: true}
			}
		}(file)
	}

	// Close the result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Read from the result channel
	for result := range resultCh {
		log.Printf("File %s: %s", result.name, getStatus(result.ok))
	}

	return okFiles, badFiles, nil
}

func getStatus(ok bool) string {
	if ok {
		return "OK"
	}
	return "Mismatch"
}

func main() {
	files := []File{
		{Name: "file1.txt", Content: []byte("content1"), Signature: "signature1"},
		{Name: "file2.txt", Content: []byte("content1"), Signature: "signature1"},
		// Add more files as neededs
	}

	start := time.Now()

	ok, bad := newFunction(files)

	duration := time.Since(start)
	log.Printf("info: %d files in %v\n", len(ok)+len(bad), duration)
	log.Printf("ok: %v", ok)
	log.Printf("bad: %v", bad)
}

func newFunction(files []File) ([]string, []string) {
	ok, bad, _ := ValidateSigs(files)
	return ok, bad
}
