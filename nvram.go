// nvram show
// nvram get key
// nvram set key=value
// nvram set key=
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	DEFAULT_PATH = "/tmp/nvram.conf"
)

type nvram struct {
	lock *sync.RWMutex
	path string
	data map[string]string
}

func newNvram() *nvram {
	n := new(nvram)
	n.lock = new(sync.RWMutex)
	n.data = make(map[string]string)
	return n
}

func loadFile(path string) (*nvram, bool) {
	var f *os.File
	var err error
	if f, err = os.Open(path); err != nil {
		return nil, false
	}

	defer f.Close()

	n := newNvram()
	n.path = path

	// Create buffer reader.
	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil {
			// Unexpected error
			if err != io.EOF {
				return nil, false
			}

			// Reached end of file, if nothing to read then break,
			// otherwise handle the last line.
			if len(line) == 0 {
				break
			}
		}

		// switch written for readability (not performance)
		switch {
		case len(line) == 0: // Empty line
			continue
		case line[0] == '#' || line[0] == ';': // Comment
			continue
		default: // Other alternatives
			i := strings.IndexAny(line, "=:")
			if i > 0 {
				key := strings.TrimSpace(line[0:i])
				value := strings.TrimSpace(line[i+1:])
				n.Set(key, value)

			}
		}

		// Reached end of file
		if err == io.EOF {
			break
		}
	}

	return n, true
}

func saveFile(n *nvram) bool {
	// Write configuration file by filename
	var f *os.File
	var err error
	if f, err = os.Create(n.path); err != nil {
		return false
	}

	// Data buffer
	buf := bytes.NewBuffer(nil)
	// Write sections
	for k, v := range n.data {
		s := fmt.Sprintf("%s=%s\n", k, v)

		if v == "" {
			continue
		}

		if _, err = buf.WriteString(s); err != nil {
			return false
		}
	}
	// Write to file
	buf.WriteTo(f)
	f.Sync()
	f.Close()
	return false
}

func (this *nvram) Show() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	for k, v := range this.data {
		fmt.Printf("%s=%s\n", k, v)
	}
	return "ERROR"
}

func (this *nvram) Get(key string) string {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.data[key]; ok {
		return this.data[key]
	}

	return fmt.Sprintf("key[%s] is not exist", key)
}

func (this *nvram) Set(key, value string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data[key] = value
	return true
}

func init() {

}

var usageTemplate = "nvram show\nnvram set key=[value]\nnvram get key\n"

func usage() {
	fmt.Fprintf(os.Stderr, usageTemplate)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
		return
	}

	if args[0] == "help" {
		usage()
		return
	}

	nvram, ok := loadFile(DEFAULT_PATH)
	if !ok {
		fmt.Println(os.Stderr, "Load", DEFAULT_PATH, "error")
	}
	defer saveFile(nvram)

	switch args[0] {
	case "set":
		if len(args) != 2 {
			usage()
			return
		}

		line := args[1]
		i := strings.IndexAny(line, "=:")
		if i > 0 {
			key := strings.TrimSpace(line[0:i])
			value := strings.TrimSpace(line[i+1:])
			nvram.Set(key, value)
		}
	case "get":
		if len(args) != 2 {
			usage()
			return
		}
		value := nvram.Get(args[1])
		fmt.Printf("%s\n", value)
	case "show":
		nvram.Show()
	default:
		usage()
	}
	return
}
