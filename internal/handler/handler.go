package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/satyamsoni2211/go_webapp_encoded_files/internal"
)

type FileHandler struct {
}

var (
	Dir string
	err error
)

func init() {
	Dir, err = ioutil.TempDir("", "webapp")
	if err != nil {
		fmt.Printf("Could not create temp dir %s", err)
		os.Exit(1)
	}
}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// Since we have encoded all the static files
// We are using custom handler to serve the static files from the
// generated encoded files
func (f *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	content, hasEncodedFile := internal.Data[path]
	if !hasEncodedFile {
		w.Write([]byte(strings.Join([]string{"No Such file ", path}, " ")))
		w.WriteHeader(404)
		return
	}
	file, err := ioutil.TempFile(Dir, path)
	if err != nil {
		w.Write([]byte("Error occurred while creating temp file"))
		w.WriteHeader(500)
		return
	}
	defer file.Close()
	n, err := file.Write(content)
	if err != nil {
		w.Write([]byte("Error occurred while writing temp file"))
		w.WriteHeader(500)
		return
	}
	stat, _ := file.Stat()
	fmt.Printf("Written %d bytes \n", n)
	http.ServeContent(w, r, path, stat.ModTime(), file)
}
