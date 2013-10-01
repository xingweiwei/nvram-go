// nvram show
// nvram get key
// nvram set key=value
// nvram set key=
package main

import (
	"flag"
	"fmt"
	//"io"
	"log"
	"os"
	//"strings"
	"sync"
)

const (
	DEFAULT_PATH = "/tmp/nvram.conf"
)

type nvram struct {
	lock *sync.RWMutex
	path string
	data map[string]map[string]string
}

func newNvram() *nvram {
	n := new(nvram)
	n.lock = new(sync.RWMutex)
	n.data = make(map[string]map[string]string)
	return n
}

func loadFile(path string) (*nvram, bool) {
	n := newNvram()
	return n, true
}

func saveFile(n *nvram) bool {
	fmt.Println(os.Stderr, "Save file!")
	os.Exit(2)
	return false
}

func (this *nvram) Show() string {
	for _, n := range this.data {
		fmt.Println(os.Stderr, n)
	}
	return "ERROR"
}

func (this *nvram) Get(key string) string {
	if _, ok := this.data[key]; ok {
		return this.data[key][key]
	}
	return "ERROR"
}

func (this *nvram) Set(key, value string) bool {
	fmt.Println(os.Stderr, "key=", key, "value=", value)
	n := make(map[string]string)
	n[key] = value
	this.data[key] = n
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
		nvram.Get(args[1])
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
