package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ancientlore/flagcfg"
	"github.com/ancientlore/kubismus"
	"github.com/facebookgo/flagenv"
)

// github.com/ancientlore/binder is used to package the web files into the executable.
//go:generate binder -package main -o webcontent.go media/*.png

const (
	VERSION = "0.2"
)

var (
	addr string = ":8080"
	cpus int    = 1
	help bool
)

func init() {
	// http service/status address
	flag.StringVar(&addr, "addr", addr, "HTTP service address for monitoring.")

	// runtime
	flag.IntVar(&cpus, "cpu", cpus, "Number of CPUs to use.")

	// help
	flag.BoolVar(&help, "help", false, "Show help.")
}

func showHelp() {
	fmt.Println(`webnull

/dev/null for http

Usage:
  webnull [options]

Example:
  webnull -addr :8080

Options:`)
	flag.PrintDefaults()
	fmt.Println(`
All of the options can be set via environment variables prefixed with "WEBNULL_".

Options can also be specified in a TOML configuration file named "webnull.config". The location
of the file can be overridden with the WEBNULL_CONFIG environment variable.`)
}

func main() {
	// Parse flags from command-line
	flag.Parse()

	// Parser flags from config
	flagcfg.AddDefaults()
	flagcfg.Parse()

	// Parse flags from environment (using github.com/facebookgo/flagenv)
	flagenv.Prefix = "WEBNULL_"
	flagenv.Parse()

	if help {
		showHelp()
		return
	}

	// setup number of CPUs
	runtime.GOMAXPROCS(cpus)

	name, _ := os.Hostname()

	http.Handle("/status/", http.StripPrefix("/status", http.HandlerFunc(kubismus.ServeHTTP)))
	http.Handle("/", kubismus.HttpRequestMetric("Requests", handleRequest()))
	http.Handle("/http", kubismus.HttpRequestMetric("Requests", handleRequest2()))
	http.Handle("/http/", kubismus.HttpRequestMetric("Requests", handleRequest2()))
	http.HandleFunc("/media/", ServeHTTP)

	kubismus.Setup("/web/null", "/media/null.png")
	kubismus.Note("Host Name", name)
	kubismus.Note("CPUs", fmt.Sprintf("%d of %d", runtime.GOMAXPROCS(0), runtime.NumCPU()))
	kubismus.Define("Requests", kubismus.COUNT, "Requests/second")
	kubismus.Define("Requests", kubismus.SUM, "Bytes/second")

	go calcMetrics()

	log.Fatal(http.ListenAndServe(addr, nil))
}

func calcMetrics() {
	tck := time.NewTicker(time.Duration(10) * time.Second)
	for {
		select {
		case <-tck.C:
			kubismus.Note("Goroutines", fmt.Sprintf("%d", runtime.NumGoroutine()))
			go func() {
				v := kubismus.GetMetrics("Requests", kubismus.SUM)
				defer kubismus.ReleaseMetrics(v)
				c := kubismus.GetMetrics("Requests", kubismus.COUNT)
				defer kubismus.ReleaseMetrics(c)
				sz := len(c)
				T := 0.0
				C := 0.0
				for i := sz - 60; i < sz; i++ {
					C += c[i]
					T += v[i]
				}
				A := 0.0
				if C > 0.0 {
					A = T / C
				}
				kubismus.Note("Last One Minute", fmt.Sprintf("%.0f Requests, %.0f Average Size, %0.f Bytes", C, A, T))
			}()
		}
	}
}

func handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body = "OK"
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		io.WriteString(w, body)
	}
}

func handleRequest2() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s := strings.TrimPrefix(req.URL.Path, "/http/")
		log.Print(s)
		id, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			id = 0
		}
		body := http.StatusText(int(id))
		if body == "" {
			body = http.StatusText(http.StatusOK)
			id = http.StatusOK
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(int(id))
		io.WriteString(w, body)
	}
}
