package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

// FnContext represents the context provided by the function invocation,
// including the request path from the API Gateway.
type FnContext struct {
	Path string `json:"path"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// --- 1. Extract the jobID from the request path ---
	var fnCtx FnContext
	// The request body from the API Gateway contains context, including the path
	if err := json.NewDecoder(in).Decode(&fnCtx); err != nil {
		log.Printf("Error decoding function context: %v", err)
		http.Error(out.(http.ResponseWriter), "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jobID := extractJobID(fnCtx.Path)
	if jobID == "" {
		log.Printf("Could not extract jobID from path: %s", fnCtx.Path)
		http.Error(out.(http.ResponseWriter), "Bad Request: Missing jobID", http.StatusBadRequest)
		return
	}

	log.Printf("Received download request for job: %s", jobID)

	// --- 2. Generate the Pre-Authenticated Request (PAR) URL ---
	parURL, err := createPAR(ctx, jobID)
	if err != nil {
		log.Fatalf("Failed to create PAR for job %s: %v", jobID, err)
		http.Error(out.(http.ResponseWriter), "Failed to generate download link", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully generated PAR for job %s", jobID)

	// --- 3. Perform the Redirect ---
	w, ok := out.(http.ResponseWriter)
	if !ok {
		log.Fatal("Output is not an http.ResponseWriter, cannot perform redirect")
		return
	}

	w.Header().Set("Location", *parURL)
	w.WriteHeader(http.StatusFound)
}

// extractJobID helper function to extract the job ID from the request path.
// IMPORTANT: Adjust the logic here if your API Gateway route path is different.
func extractJobID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Example path: /api/v1/jobs/some-job-id/download
	// After splitting, "some-job-id" would be at index 3.
	if len(parts) >= 4 && parts[2] == "jobs" && parts[4] == "download" {
		return parts[3]
	}
	return ""
}

// createPAR generates a Pre-Authenticated Request for the specified object.
func createPAR(ctx context.Context, jobID string) (*string, error) {
	// Use instance principal for authentication within the OCI Function.
	provider, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		return nil, err
	}

	osClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// Get required details from the function's environment variables.
	namespace := os.Getenv("OBJECT_STORAGE_NAMESPACE")
	bucketName := os.Getenv("OUTPUT_BUCKET_NAME")
	region := os.Getenv("OCI_REGION")
	objectName := jobID + ".zip" // The converted file should be named after the job ID.

	// Set the expiration time for the PAR.
	expirationTime := time.Now().Add(15 * time.Minute)

	// Create the request to generate the PAR.
	req := objectstorage.CreatePreauthenticatedRequestRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:       common.String("par-for-" + jobID),
			ObjectName: common.String(objectName),
			AccessType: objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectread,
			// FIX: Create an SDKTime struct and pass a pointer to it.
			TimeExpires: &common.SDKTime{Time: expirationTime},
		},
	}

	resp, err := osClient.CreatePreauthenticatedRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// Construct the full, absolute URL for the download.
	fullURL := "https://objectstorage." + region + ".oraclecloud.com" + *resp.AccessUri
	return &fullURL, nil
}
