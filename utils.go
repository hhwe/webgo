package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

var (
	// todo: add this to config file or os environment
	salt = "123456"
)

// HashSHA256 hash string to string, add salt to hash
func HashSHA256(p string) (d string) {
	h := md5.New()
	io.WriteString(h, p)
	d = string(h.Sum(nil))
	io.WriteString(h, salt)
	io.WriteString(h, d)
	d = string(h.Sum(nil))
	logger.Println(d)
	return
}

// Recovery recovers from any panics and writes a 500 if there was one.
func Recovery(w http.ResponseWriter) {
	if e := recover(); e != nil {
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "Error: %v\n", e)

		for i := 1; ; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		}

		var content = buf.String()
		logger.Println(content)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
