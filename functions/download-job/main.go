package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	jobID := "jobId-from-request" // Extract jobID from request path

	// Replace with OCI Object Storage SDK to generate a PAR
	parURL := "https://objectstorage.us-ashburn-1.oraclecloud.com/p/..." // Placeholder
	log.Printf("Generating PAR for job %s", jobID)

	// Redirect the user to the PAR URL
	w := out.(http.ResponseWriter)
	http.Redirect(w, nil, parURL, http.StatusFound)
}
