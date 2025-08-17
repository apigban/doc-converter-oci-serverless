package main

import (
	"context"
	"doc-converter-oci-serverless/pkg/queue"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/fnproject/fdk-go"
)

type ConversionRequest struct {
	URLs     []string `json:"urls"`
	Selector string   `json:"selector"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var req ConversionRequest
	json.NewDecoder(in).Decode(&req)

	jobID := "generated-job-id"

	// 1. Get Queue OCID from environment variable (set in function config)
	queueID := os.Getenv("QUEUE_OCID")
	if queueID == "" {
		log.Fatal("QUEUE_OCID environment variable not set")
	}

	// 2. Create the OCI Queue client
	queueClient, err := queue.NewOCIQueueClient(queueID)
	if err != nil {
		log.Fatalf("Failed to create OCI Queue client: %v", err)
	}

	// 3. Create and publish the job
	job := &queue.ConversionJob{
		URLs:       req.URLs,
		Selector:   req.Selector,
		DownloadID: jobID,
	}

	err = queueClient.PutMessage(job)
	if err != nil {
		log.Fatalf("Failed to publish job to queue: %v", err)
	}

	log.Printf("Job %s queued successfully", jobID)

	out.Write([]byte("{\"jobId\": \"" + jobID + "\"}"))
}
