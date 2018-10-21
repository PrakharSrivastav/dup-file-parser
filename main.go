package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"encoding/json"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/PrakharSrivastav/dup-file-parser/model"
	"google.golang.org/api/iterator"
)

var mu sync.Mutex

func main() {
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, "digital-utility-playground")
	if err != nil {
		fmt.Println(err.Error())
	}
	received := 0
	sub := pubsubClient.Subscription("processFileUpload")
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		var b model.Bucket
		err = json.Unmarshal([]byte(msg.Data), &b)
		if err != nil {
			fmt.Println(err.Error())
		}
		// fmt.Printf("Got message: %#v\n", &b)
		readFile(ctx, b)
		mu.Lock()
		defer mu.Unlock()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func readFile(ctx context.Context, b model.Bucket) error {
	fmt.Println("Reading from storage")
	fmt.Printf("Bucket Details %#v \n", b)
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer storageClient.Close()

	bucket := storageClient.Bucket(b.Bucket)
	query := &storage.Query{}
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Println(attrs.Name)
	}
	return nil
}
