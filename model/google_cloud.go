package model

import (
	"context"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
)

func NewSecretClient() (*secretmanager.Client, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Printf("failed to setup client: %v", err)
		return nil, err
	}
	return client, nil
}

func NewComputeEngineClient() (*compute.InstancesClient, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Printf("failed to setup client: %v", err)
		return nil, err
	}
	return instancesClient, nil
}
