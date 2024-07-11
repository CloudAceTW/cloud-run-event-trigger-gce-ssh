package model

import (
	"context"

	"cloud.google.com/go/compute/apiv1/computepb"
)

type GceInstance struct {
	Project       string
	Zone          string
	Instance      string
	RestartStatus bool
}

func NewGceInstance(project, zone, instance string) *GceInstance {
	return &GceInstance{
		Project:       project,
		Zone:          zone,
		Instance:      instance,
		RestartStatus: false,
	}
}

func (gceInstance *GceInstance) RestartVM() error {
	client, err := NewComputeEngineClient()
	if err != nil {
		return err
	}
	req := &computepb.ResetInstanceRequest{
		Project:  gceInstance.Project,
		Zone:     gceInstance.Zone,
		Instance: gceInstance.Instance,
	}
	_, err = client.Reset(context.Background(), req)
	if err != nil {
		return err
	}
	gceInstance.RestartStatus = true
	return nil
}
