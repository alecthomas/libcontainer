package main

import (
	"fmt"
	"os"

	"github.com/docker/libcontainer/namespaces"
)

func initAction(command []string, console string, pipeFd int) error {
	container, err := loadContainer()
	if err != nil {
		return err
	}

	rootfs, err := os.Getwd()
	if err != nil {
		return err
	}

	syncPipe, err := namespaces.NewSyncPipeFromFd(0, uintptr(pipeFd))
	if err != nil {
		return fmt.Errorf("unable to create sync pipe: %s", err)
	}

	if err := namespaces.Init(container, rootfs, console, syncPipe, command); err != nil {
		return fmt.Errorf("unable to initialize for container: %s", err)
	}

	return nil
}
