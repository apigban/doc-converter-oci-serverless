package main

import (
	"context"
	"doc-converter-oci-serverless/pkg/converter"
	"doc-converter-oci-serverless/pkg/queue"
	"encoding/json"
	"io"
	"log"

	"github.com/fnproject/fdk-go"
)

// OCIQueueEvent represents the structure of the event from an OCI Queue trigger
type OCIQueueEvent struct {
	Messages []struct {
		Content string `json:"content"`
	} `json:"messages"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// 1. Decode the incoming event from the OCI Queue trigger
	var event OCIQueueEvent
	json.NewDecoder(in).Decode(&event)

	// A single trigger can contain multiple messages
	for _, message := range event.Messages {
		var job queue.ConversionJob
		// 2. Unmarshal the message content into your job struct
		if err := json.Unmarshal([]byte(message.Content), &job); err != nil {
			log.Printf("ERROR: Failed to unmarshal job from queue message: %v", err)
			continue // Move to the next message
		}

		log.Printf("Processing job %s", job.DownloadID)

		c, err := converter.NewConverterForJob(job.DownloadID)
		if err != nil {
			log.Printf("ERROR: Failed to create new converter for job %s: %v", job.DownloadID, err)
			return
		}

		resultsChan, summaryChan := c.Convert(job.URLs, job.Selector)

		for range resultsChan {
			// Drain results
		}

		summary := <-summaryChan
		log.Printf("INFO: Conversion finished for job %s. Successful: %d, Failed: %d",
			job.DownloadID, summary.Successful, summary.Failed)
	}
}
