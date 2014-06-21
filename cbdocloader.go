package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/couchbaselabs/go-couchbase"
)

func readFile(f zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	return string(data), err
}

func storeDesignDoc(bucket couchbase.Bucket, doc string, f zip.File) error {
	value, err := readFile(f)
	if err != nil {
		return err
	}
	return bucket.PutDDoc(doc, value)
}

func storeDocument(bucket couchbase.Bucket, doc string, f zip.File) error {
	key := strings.Split(doc, ".")
	if len(key[0]) == 0 {
		return errors.New("invaid key name for " + f.Name)
	}

	value, err := readFile(f)
	if err != nil {
		return err
	}

	return bucket.Set(key[0], 0, value)
}

func main() {
	_ = flag.String("u", "Ignored", "This parameter is ignored")
	password := flag.String("p", "", "Password for bucket")
	bucketname := flag.String("b", "mybucket", "The name of the bucket to create")
	url := flag.String("n", "127.0.0.1:8091", "Node address")
	_ = flag.Int("s", 100, "This parameter is ignored")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: cbdocloader [arguments] zipfile")
		os.Exit(1)
	}

	r, err := zip.OpenReader(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	uri := "http://" + *bucketname + ":" + *password + "@" + *url
	bucket, err := couchbase.GetBucket(uri, "default", *bucketname)
	if err != nil {
		log.Fatalf("Failed to connect to cluster %s: %v", *url, err)
	}

	for _, f := range r.File {
		doc := strings.Split(f.Name, "/docs/")
		if len(doc) == 2 && len(doc[1]) > 0 {
			err := storeDocument(*bucket, doc[1], *f)
			if err != nil {
				log.Fatalf("Failed to store %s: %v", doc[1], err)
			}
		} else {
			doc := strings.Split(f.Name, "/design_docs/")
			if len(doc) == 2 && len(doc[1]) > 0 {
				err := storeDesignDoc(*bucket, doc[1], *f)
				if err != nil {
					log.Fatalf("Failed to store %s: %v", doc[1], err)
				}
			}
		}
	}
}
