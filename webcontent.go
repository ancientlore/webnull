package main

import (
	"github.com/golang/snappy"
	"encoding/base64"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

const (
	cMedia_NullPng = "iQRMiVBORw0KGgoAAAANSUhEUgAAACQBBPBSCAMAAADW3miqAAAABGdBTUEAALGPC_xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAANlBMVEUAAAAVFDyyAwDsAAAAyZVz2QAAABB0Uk5TAMAwINCQYBBQgKBw8EDgsCRsCvcAAAABYktHRACIBR1IAAAACXBIWXMAAC4jAQT0MAEBeKU_dgAAARRJREFUOMuFU1sChCAIVDczzJT7n3Z94KPEmq_SEYYBhOAhEZX4wA9x--IIjbiv7swBVkp9KofOLHJYhw2W5-wDBVGykXa6lLKQN4Z1pQtf1NrCmkl-eEx5Ty5ZLSh-25D-uUDVmeS2cThZpYaHxe1oJ4Y76RycKW6rOV8YgpPbKd-9x1s8ubpsUlaPCKniQXY1CxakNiTwQoIq7o0UDTNfpKNZ8UIKrSQ9kTzVG2X7XuXDAklmQu-9n8wEakuTnWI-23KVoy47d_PRYFEmo8u-zU5FqsWpLvvkhu6XQrkmO28FPEn5KZJsA4tFoBXxAbTMX_wK23E32bVLSQaKA8EjqtIQ4gZLOMQK3e01rtlhVvcqyx_mHhaYGJ81fQAAAABJRU5ErkJggg=="
)

var staticFiles = map[string]string{
	"/media/null.png": cMedia_NullPng,
}

func Lookup(path string) []byte {
	s, ok := staticFiles[path]
	if !ok {
		return nil
	} else {
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
}

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
		w.Write(b)
	} else {
		http.NotFound(w, r)
	}
}
