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

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var job queue.ConversionJob
	json.NewDecoder(in).Decode(&job)

	log.Printf("Processing job %s", job.DownloadID)

	// The converter logic is modified to use the OCI SDK for object storage
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

	// After conversion, create a ZIP in-memory and upload to Object Storage
}
