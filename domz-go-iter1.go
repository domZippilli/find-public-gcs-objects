/*
This package is an executable.

It provides efficient scanning of a GCS bucket's objects for those with ACLs which permit
"allUsers" to read the object.
*/
package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Main performs setup with GCS, and dispatches the root goroutine to list /
func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if len(os.Args) < 2 {
		panic("Please give bucket name as the first and only arg.")
	}

	bh := client.Bucket(os.Args[1])
	if _, err = bh.Attrs(ctx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go ListAndCheck(ctx, &wg, "/", "", bh)
	wg.Wait()
}

// List implements a traversal of the bucket's prefixes, starting from the prefix
// with which it is invoked. Objects with prefix values (a synthetic "directory") are
// dispatched into a new List goroutine, while objects without are dispatched into a
// Check goroutine.
func ListAndCheck(ctx context.Context, wg *sync.WaitGroup, delimiter string, prefix string, bh *storage.BucketHandle) {
	//fmt.Println("LISTING: " + prefix)
	dirQuery := storage.Query{Delimiter: delimiter, Prefix: prefix}
	entries := bh.Objects(ctx, &dirQuery)
	for {
		attrs, err := entries.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}

		if attrs.Prefix == "" {
			//fmt.Println("CHECKING: " + objAttrs.Name)
			for _, acl := range attrs.ACL {
				if acl.Entity == "allUsers" && acl.Role == "READER" {
					fmt.Println(attrs.Name)
				}
			}
		} else {
			wg.Add(1)
			go ListAndCheck(ctx, wg, delimiter, attrs.Prefix, bh)
		}
	}
	wg.Done()
}
