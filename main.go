package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"aead.dev/minisign"
)

const (
	boldYellow = "\033[1;33m"
	gray       = "\033[37m"
	reset      = "\033[0m"
)

type result struct {
	pk    minisign.PublicKey
	sk    minisign.PrivateKey
	count uint64
}

func main() {
	overwrite := flag.Bool("overwrite", false, "overwrite existing ./minisign.key and ./minisign.pub files")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-overwrite] <regexp> [<regexp> ...]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	patterns := args
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid regexp: %v\n", err)
			os.Exit(1)
		}
		regexps[i] = re
	}

	if !*overwrite {
		if _, err := os.Stat("minisign.key"); err == nil {
			fmt.Fprintf(os.Stderr, "error: ./minisign.key already exists, use -overwrite to replace it\n")
			os.Exit(1)
		}
		if _, err := os.Stat("minisign.pub"); err == nil {
			fmt.Fprintf(os.Stderr, "error: ./minisign.pub already exists, use -overwrite to replace it\n")
			os.Exit(1)
		}
	}

	start := time.Now()
	var totalCount uint64

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan result, 1)

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go worker(ctx, &wg, regexps, resultChan, &totalCount)
	}

	go progressReporter(ctx, &totalCount, start)

	res := <-resultChan
	cancel()
	wg.Wait()

	pkStr := res.pk.String()

	fmt.Printf("Found match after %d keys in %v\n", atomic.LoadUint64(&totalCount), time.Since(start).Truncate(time.Second))
	fmt.Printf("Public key: %s\n", pkStr)

	pkText, err := res.pk.MarshalText()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling public key: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("minisign.pub", []byte(pkText), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing minisign.pub: %v\n", err)
		os.Exit(1)
	}

	skText, err := res.sk.MarshalText()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling secret key: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("minisign.key", skText, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "error writing minisign.key: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Keys saved to minisign.pub and minisign.key")
	fmt.Println()
	fmt.Printf("%sWARNING: The secret key is unencrypted!%s\n", boldYellow, reset)
	fmt.Println("To protect it with a password, run:")
	fmt.Println("  minisign -C -s minisign.key")
}

func worker(ctx context.Context, wg *sync.WaitGroup, regexps []*regexp.Regexp, resultChan chan result, totalCount *uint64) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		pk, sk, err := minisign.GenerateKey(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating key: %v\n", err)
			os.Exit(1)
		}

		atomic.AddUint64(totalCount, 1)

		pkStr := pk.String()

		matchAll := true
		for _, re := range regexps {
			if !re.MatchString(pkStr) {
				matchAll = false
				break
			}
		}

		if matchAll {
			select {
			case resultChan <- result{pk, sk, atomic.LoadUint64(totalCount)}:
			case <-ctx.Done():
				return
			}
			return
		}
	}
}

func progressReporter(ctx context.Context, totalCount *uint64, start time.Time) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			count := atomic.LoadUint64(totalCount)
			fmt.Printf("%sGenerated %d keys in %v...%s\n", gray, count, time.Since(start).Truncate(time.Second), reset)
		}
	}
}
