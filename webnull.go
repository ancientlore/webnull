package main

import (
	"embed"
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

	"github.com/NYTimes/gziphandler"
	"github.com/ancientlore/kubismus"
)

var (
	addr string = ":8080"
	help bool
)

//go:embed media/*.png
var media embed.FS

func init() {
	// http service/status address
	flag.StringVar(&addr, "addr", addr, "HTTP service address for monitoring.")

	// help
	flag.BoolVar(&help, "help", false, "Show help.")
}

func showHelp() {
	fmt.Fprintln(os.Stderr, `webnull

/dev/null for http

Usage:
  webnull [options]

Example:
  webnull -addr :8080

Options:`)
	flag.PrintDefaults()
}

func main() {
	// Parse flags from command-line
	flag.Parse()

	if help {
		showHelp()
		return
	}

	name, _ := os.Hostname()

	http.Handle("/status/", gziphandler.GzipHandler(http.StripPrefix("/status", cspHandler(http.HandlerFunc(kubismus.ServeHTTP)))))
	http.Handle("/", kubismus.HttpRequestMetric("Requests", cspHandler(handleRequest())))
	http.Handle("/http", kubismus.HttpRequestMetric("Requests", cspHandler(handleRequestStatus())))
	http.Handle("/http/", kubismus.HttpRequestMetric("Requests", cspHandler(handleRequestStatus())))
	http.Handle("/delay", kubismus.HttpRequestMetric("Requests", cspHandler(handleRequestDelayMs())))
	http.Handle("/delay/", kubismus.HttpRequestMetric("Requests", cspHandler(handleRequestDelayMs())))
	http.Handle("/media/", cspHandler(http.FileServer(http.FS(media))))

	kubismus.Setup("/web/null", "/media/null.png")
	kubismus.Note("Host Name", name)
	kubismus.Note("CPUs", fmt.Sprintf("%d", runtime.NumCPU()))
	kubismus.Define("Requests", kubismus.COUNT, "Requests/second")
	kubismus.Define("Requests", kubismus.SUM, "Bytes/second")

	go calcMetrics()

	log.Fatal(http.ListenAndServe(addr, nil))
}

const csp = "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' fonts.googleapis.com; font-src 'self' fonts.googleapis.com fonts.gstatic.com"

func cspHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", csp)
		h.ServeHTTP(w, r)
	})
}

func calcMetrics() {
	tck := time.NewTicker(time.Duration(10) * time.Second)
	for range tck.C {
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

func handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body = "OK"
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		io.WriteString(w, body)
	}
}

func handleRequestStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s := strings.TrimPrefix(req.URL.Path, "/http/")
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

func handleRequestDelayMs() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s := strings.TrimPrefix(req.URL.Path, "/delay/")
		ms, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			ms = 0
		}
		delay := time.Duration(ms) * time.Millisecond
		status := http.StatusOK
		if delay < 0 || delay > 30*time.Second {
			status = http.StatusBadRequest
		} else {
			time.Sleep(delay)
		}
		body := http.StatusText(status)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(status)
		io.WriteString(w, body)
	}
}
