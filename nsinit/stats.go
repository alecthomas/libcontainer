package main

import (
	"encoding/json"
	"fmt"

	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/cgroups/fs"
)

func statsAction() error {
	container, err := loadContainer()
	if err != nil {
		return err
	}

	stats, err := getContainerStats(container)
	if err != nil {
		return fmt.Errorf("failed to get stats - %v\n", err)
	}

	fmt.Printf("Stats:\n%v\n", stats)
	return nil
}

// returns the container stats in json format.
func getContainerStats(container *libcontainer.Container) (string, error) {
	stats, err := fs.GetStats(container.Cgroups)
	if err != nil {
		return "", err
	}

	out, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		return "", err
	}

	return string(out), nil
}
