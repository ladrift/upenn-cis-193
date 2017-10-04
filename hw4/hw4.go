// Homework 4: Concurrency
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	fmt.Println("========== Problem 1a: File processing ==========")
	FileSum("file_sum.txt", "sum.txt") // It's implementation is IOSum() from Problem 1b.

	fmt.Println("========== Problem 1b: IO processing with interfaces ==========")

	fmt.Println("========== Problem 2: Concurrent map access ==========")
	d := PennDirectory{
		directory: make(map[int]string),
	}
	total := 10000000
	var wg sync.WaitGroup
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			switch i % 3 {
			case 0:
				d.Add(1, "Ziyi Yan")
			case 1:
				d.Add(2, "Xuefei Chen")
			case 2:
				d.Remove(2)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("d.directory = %v\n", d.directory)
}

// Problem 1a: File processing
// You will be provided an input file consisting of integers, one on each line.
// Your task is to read the input file, sum all the integers, and write the
// result to a separate file.

// FileSum sums the integers in input and writes them to an output file.
// The two parameters, input and output, are the filenames of those files.
// You should expect your input to end with a newline, and the output should
// have a newline after the result.
func FileSum(input, output string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() error: %v\n", err)
	}
	infile, err := os.Open(filepath.Join(dir, input))
	if err != nil {
		log.Fatalf("os.Open() error: %v\n", err)
	}
	defer infile.Close()
	outfile, err := os.OpenFile(filepath.Join(dir, output), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("os.Open() error: %v\n", err)
	}
	defer outfile.Close()

	IOSum(infile, outfile)
}

// Problem 1b: IO processing with interfaces
// You must do the exact same task as above, but instead of being passed 2
// filenames, you are passed 2 interfaces: io.Reader and io.Writer.
// See https://golang.org/pkg/io/ for information about these two interfaces.
// Note that os.Open returns an io.Reader, and os.Create returns an io.Writer.

// IOSum sums the integers in input and writes them to output
// The two parameters, input and output, are interfaces for io.Reader and
// io.Writer. The type signatures for these interfaces is in the Go
// documentation.
// You should expect your input to end with a newline, and the output should
// have a newline after the result.
func IOSum(input io.Reader, output io.Writer) {
	in := bufio.NewScanner(input)
	sum := 0
	for in.Scan() {
		n, err := strconv.Atoi(in.Text())
		if err != nil {
			log.Fatalf("strconv.Atoi() error: %v\n", err)
		}
		sum += n
	}

	out := bufio.NewWriter(output)
	out.WriteString(strconv.Itoa(sum) + "\n")
	out.Flush()
}

// Problem 2: Concurrent map access
// Maps in Go [are not safe for concurrent use](https://golang.org/doc/faq#atomic_maps).
// For this assignment, you will be building a custom map type that allows for
// concurrent access to the map using mutexes.
// The map is expected to have concurrent readers but only 1 writer can have
// access to the map.

// PennDirectory is a mapping from PennID number to PennKey (12345678 -> adelq).
// You may only add *private* fields to this struct.
// Hint: Use an embedded sync.RWMutex, see lecture 2 for a review on embedding
type PennDirectory struct {
	mu        sync.RWMutex
	directory map[int]string
}

// Add inserts a new student to the Penn Directory.
// Add should obtain a write lock, and should not allow any concurrent reads or
// writes to the map.
// You may NOT write over existing data - simply raise a warning.
func (d *PennDirectory) Add(id int, name string) {
	d.mu.Lock()
	d.directory[id] = name
	d.mu.Unlock()
}

// Get fetches a student from the Penn Directory by their PennID.
// Get should obtain a read lock, and should allow concurrent read access but
// not write access.
func (d *PennDirectory) Get(id int) string {
	d.mu.RLock()
	name := d.directory[id]
	d.mu.RUnlock()
	return name
}

// Remove deletes a student to the Penn Directory.
// Remove should obtain a write lock, and should not allow any concurrent reads
// or writes to the map.
func (d *PennDirectory) Remove(id int) {
	d.mu.Lock()
	delete(d.directory, id)
	d.mu.Unlock()
}
