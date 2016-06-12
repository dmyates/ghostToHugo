package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dmyates/ghostToHugo/lib/ghost"
	"github.com/dmyates/ghostToHugo/lib/hugo"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] <Ghost Export>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	var c hugo.Config
	var l string

	flag.Usage = usage

	flag.StringVar(&c.Path, "hugo", ".", "Path to hugo project")
	flag.StringVar(&l, "location", "",
		"Location to use for time conversions (default: local)")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if err := hugo.Init(c); err != nil {
		log.Fatalf("Error initializing Hugo Config (%v)", err)
	}

	if l != "" {
		location, err := time.LoadLocation(l)
		if err != nil {
			log.Fatalf("Error loading location %s: %v", l, err)
		}
		ghost.SetLocation(location)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("Error opening export: %v", err)
	}
	defer file.Close()

	reader := ghost.ExportReader{file}

	entries, err := ghost.Process(reader)
	if err != nil {
		log.Fatalf("Error processing Ghost export: %v", err)
	}

	var wg sync.WaitGroup
	for _, entry := range entries {
		wg.Add(1)
		go func(data ghost.ExportData) {
			defer wg.Done()
			hugo.ExportGhost(&data)
		}(entry.Data)
	}

	wg.Wait()
}
