package main

import (
	"github.com/golang/snappy"
	"encoding/base64"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	cMedia_NullPng = "iQRMiVBORw0KGgoAAAANSUhEUgAAACQBBPBSCAMAAADW3miqAAAABGdBTUEAALGPC_xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAANlBMVEUAAAAVFDyyAwD0cAEAAADJlXPZAAAAEHRSTlMAwDAg0JBgEFCAoHDwQOCwJGwK9wAAAAFiS0dEAIgFHUgAAAAJcEhZcwAALiMAAC4jAXilP3YAAAEUSURBVDjLhVNbAoQgCFQ3M8yU-592feCjxJqv0hGGAYTgIRGV-MAPcfviCI24r-7MAVZKfSqHzixyWIcNlufsAwVRspF2upSykDeGdaULX9TawppJfnhMeU8uWS0oftuQ_rlA1ZnktnE4WaWGh8XtaCeGO-kcnCluqzlfGIKT2ynfvcdbPLm6bFJWjwip4kF2NQsWpDYk8EKCKu6NFA0zX6SjWfFCCq0kPZE81Rtl-17lwwJJZkLvvZ_MBGpLk51iPttylaMuO3fz0WBRJqPLvs1ORarFqS775Ibul0K5JjtvBTxJ-SmSbAOLRaAV8QG0zF_8CttxN9m1S0kGigPBI6rSEOIGSzjECt3tNa7ZYVb3Kssf5h4WmBifNX0AAAAASUVORK5CYII="
)

var staticFiles = map[string]string{
	"/media/null.png": cMedia_NullPng,
}

// Lookup returns the bytes associated with the given path, or nil if the path was not found.
func Lookup(path string) []byte {
	s, ok := staticFiles[path]
	if ok {
		d, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			log.Print("main.Lookup: ", err)
			return nil
		}
		r, err := snappy.Decode(nil, d)
		if err != nil {
			log.Print("main.Lookup: ", err)
			return nil
		}
		return r
	}
	return nil
}

// ServeHTTP serves the stored file data over HTTP.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if strings.HasSuffix(p, "/") {
		p += "index.html"
	}
	b := Lookup(p)
	if b != nil {
		mt := mime.TypeByExtension(path.Ext(p))
		if mt != "" {
			w.Header().Set("Content-Type", mt)
		}
		w.Header().Set("Expires", time.Now().AddDate(0, 0, 1).Format(time.RFC1123))
		w.Write(b)
	} else {
		http.NotFound(w, r)
	}
}
