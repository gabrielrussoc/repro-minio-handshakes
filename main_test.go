package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"runtime/pprof"
	"sync"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func get(ctx context.Context, minioCore *minio.Core, bucket string, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, _, _, err := minioCore.GetObject(
		ctx,
		bucket,
		key,
		minio.GetObjectOptions{},
	)
	if err != nil {
		fmt.Println("failed to get object", err)
	}
}

func TestHandshake(t *testing.T) {
	opts := &minio.Options{
		Creds:  credentials.NewEnvAWS(),
		Region: "us-west-2",
		Secure: true,
	}
	minioCore, err := minio.NewCore("s3.us-west-2.amazonaws.com", opts)
	if err != nil {
		log.Fatalln("failed to create minio client", err)
		return
	}
	bucket := "my_bucket"
	key := "my_key"
	var wg sync.WaitGroup
	calls := 100
	wg.Add(calls)
	file, err := ioutil.TempFile("/tmp", "cpu-*.pprof")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	pprof.StartCPUProfile(file)
	for i := 0; i < calls; i++ {
		go get(ctx, minioCore, bucket, key, &wg)
	}
	wg.Wait()
	pprof.StopCPUProfile()
	fmt.Printf("Profile available at %s\n", file.Name())
}
