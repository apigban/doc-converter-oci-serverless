package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/fnproject/fdk-go"
	// This would be replaced with OCI Queue SDK
)

type ConversionRequest struct {
	URLs     []string `json:"urls"`
	Selector string   `json:"selector"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// OCI Function Handler Logic
	var req ConversionRequest
	json.NewDecoder(in).Decode(&req)

	jobID := "generated-job-id" // Generate a unique job ID

	// Replace with OCI Queue SDK to publish the job
	log.Printf("Job %s queued successfully", jobID)

	out.Write([]byte("{\"jobId\": \"" + jobID + "\"}"))
}
