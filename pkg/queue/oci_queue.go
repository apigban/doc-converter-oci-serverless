package queue

import (
	"context"
	"encoding/json"

	// Import the auth package
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth" // <-- Add this import
	"github.com/oracle/oci-go-sdk/v65/queue"
)

// OCIQueueClient holds the client and queue ID
type OCIQueueClient struct {
	client  queue.QueueClient
	queueID string
}

// ConversionJob defines the structure for a job message.
// This struct remains the same as before.
type ConversionJob struct {
	URLs       []string `json:"urls"`
	Selector   string   `json:"selector"`
	DownloadID string   `json:"downloadId"`
}

// NewOCIQueueClient creates a new client to interact with OCI Queues.
func NewOCIQueueClient(queueID string) (*OCIQueueClient, error) {
	// FIX: Use the InstancePrincipalConfigurationProvider from the auth package
	provider, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		return nil, err
	}

	client, err := queue.NewQueueClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	return &OCIQueueClient{
		client:  client,
		queueID: queueID,
	}, nil
}

// PutMessage publishes a new job to the OCI Queue.
// This function remains the same as before.
func (c *OCIQueueClient) PutMessage(job *ConversionJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	req := queue.PutMessagesRequest{
		QueueId: &c.queueID,
		PutMessagesDetails: queue.PutMessagesDetails{
			Messages: []queue.PutMessagesDetailsEntry{
				{
					Content: common.String(string(body)),
				},
			},
		},
	}

	_, err = c.client.PutMessages(context.Background(), req)
	return err
}
