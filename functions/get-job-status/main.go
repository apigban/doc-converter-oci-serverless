package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	jobID := "jobId-from-request" // Extract jobID from request path

	// Replace with OCI Object Storage SDK to check for file existence
	zipFileExists := true // Placeholder
	log.Printf("Checking status for job %s", jobID)

	var status string
	if zipFileExists {
		status = "completed"
	} else {
		status = "pending"
	}

	json.NewEncoder(out).Encode(map[string]string{"status": status})
}
