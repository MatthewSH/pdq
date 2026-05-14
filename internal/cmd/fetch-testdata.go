package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const baseURL = "https://raw.githubusercontent.com/facebook/ThreatExchange/main/pdq/data"

var fixtures = []string{
	"misc-images/c.png",
	"misc-images/small.jpg",
	"misc-images/wee.jpg",
	"reg-test-input/labelme-subset/q0003.jpg",
	"reg-test-input/labelme-subset/q0004.jpg",
	"reg-test-input/labelme-subset/q0122.jpg",
	"reg-test-input/labelme-subset/q0291.jpg",
	"reg-test-input/labelme-subset/q0746.jpg",
	"reg-test-input/labelme-subset/q1050.jpg",
	"reg-test-input/labelme-subset/q2821.jpg",
	"reg-test-input/dih/bridge-1-original.jpg",
	"reg-test-input/dih/bridge-2-rotate-90.jpg",
	"reg-test-input/dih/bridge-3-rotate-180.jpg",
	"reg-test-input/dih/bridge-4-rotate-270.jpg",
	"reg-test-input/dih/bridge-5-flipx.jpg",
	"reg-test-input/dih/bridge-6-flipy.jpg",
	"reg-test-input/dih/bridge-7-flip-plus-1.jpg",
	"reg-test-input/dih/bridge-8-flip-minus-1.jpg",
}

func main() {
	dir := flag.String("dir", "testdata", "directory to write testdata to")
	force := flag.Bool("force", false, "re-download testdata even if it exists")
	flag.Parse()

	var errs int

	for _, rel := range fixtures {
		dest := filepath.Join(*dir, filepath.FromSlash(rel))

		if !*force {
			if _, err := os.Stat(dest); err == nil {
				fmt.Printf("skip %s\n", rel)
				continue
			}
		}

		if err := fetch(baseURL+"/"+rel, dest); err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", rel, err)
			errs++
		} else {
			fmt.Printf("ok    %s\n", rel)
		}
	}

	if errs > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "%d files failed to download\n", errs)
		os.Exit(1)
	}

	fmt.Printf("\nAll testdata files ready in %s\n", *dir)
}

func fetch(url, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}

	f, err := os.CreateTemp(filepath.Dir(dst), ".fetch-*")
	if err != nil {
		return fmt.Errorf("temp file: %w", err)
	}
	tmpPath := f.Name()

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close: %w", err)
	}
	if err := os.Rename(tmpPath, dst); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}
