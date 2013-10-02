// nvram show
// nvram get key
// nvram set key=value
// nvram set key=
package main

import (
	"bufio"
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
	lock    *sync.RWMutex
	path    string
	keyList []string
	data    map[string]string
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
	//fmt.Println(os.Stderr, "Save file!")
	os.Exit(2)
	return false
}

func (this *nvram) Show() string {
	for _, k := range this.keyList {
		value := this.data[k]
		fmt.Printf("%s=%s\n", k, value)
	}
	return "ERROR"
}

func (this *nvram) Get(key string) string {
	if _, ok := this.data[key]; ok {
		return this.data[key]
	}

	return fmt.Sprintf("key[%s] is not exist", key)
}

func (this *nvram) Set(key, value string) bool {
	//fmt.Printf("%s=%s", key, value)
	var isExist bool = false
	this.data[key] = value
	for _, k := range this.keyList {
		if k == key {
			isExist = true
			break
		}
	}
	if isExist == false {
		this.keyList = append(this.keyList, key)
	}
	return true
}

func init() {

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
		if len(args) < 3 {
			return
		}
		nvram.Set(args[1], args[2])
	case "get":
		if len(args) < 2 {
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

var usageTemplate = "nvram show\nnvram set key value\nnvram get key\n"

func usage() {
	fmt.Fprintf(os.Stderr, usageTemplate)
	os.Exit(2)
}
