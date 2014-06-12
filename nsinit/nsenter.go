package main

import (
	"fmt"

	"github.com/docker/libcontainer/namespaces"
)

func nsenterAction(nspid int, label, json string, cmd []string) error {
	container, err := loadContainerFromJson(json)
	if err != nil {
		return fmt.Errorf("unable to load container: %s", err)
	}

	if nspid <= 0 {
		return fmt.Errorf("cannot enter into namespaces without valid pid: %q", nspid)
	}

	if err := namespaces.NsEnter(container, label, nspid, cmd); err != nil {
		return fmt.Errorf("failed to nsenter: %s", err)
	}

	return nil
}
