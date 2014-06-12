package main

import (
	"encoding/json"
	"fmt"

	"github.com/docker/libcontainer"
)

func specAction() error {
	container, err := loadContainer()
	if err != nil {
		return err
	}

	spec, err := getContainerSpec(container)
	if err != nil {
		return fmt.Errorf("Failed to get spec - %v\n", err)
	}

	fmt.Printf("Spec:\n%v\n", spec)
	return nil
}

// returns the container spec in json format.
func getContainerSpec(container *libcontainer.Container) (string, error) {
	spec, err := json.MarshalIndent(container, "", "\t")
	if err != nil {
		return "", err
	}

	return string(spec), nil
}
