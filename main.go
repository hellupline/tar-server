package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/hellupline/leto"
)

func main() {
	var fname string
	flag.StringVar(&fname, "filename", "", "tar filename")

	var portNumber int
	flag.IntVar(&portNumber, "port", 8001, "http port")

	flag.Parse()

	if err := loadFiles(fname); err != nil {
		log.Fatal(err)
	}

	fmt.Println("files:")
	for k, _ := range leto.Default().Files {
		fmt.Println(k)
	}
	fmt.Println()

	http.Handle("/", http.FileServer(leto.Default()))

	log.Printf("Listening on %d...\n", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil))
}

func loadFiles(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return errors.Wrapf(err, "error opening %s", fname)
	}
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return errors.Wrap(err, "error reading tar archive")
		}

		data, err := ioutil.ReadAll(tr)
		if err != nil {
			return errors.Wrapf(err, "error reading %s from tar", hdr.Name)
		}

		leto.Register("/"+path.Clean(hdr.Name), data)
	}
	return nil
}
