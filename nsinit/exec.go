package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/namespaces"
)

func execAction(command []string) error {
	var nspid, exitCode int

	container, err := loadContainer()
	if err != nil {
		log.Fatal(err)
	}

	if nspid, err = readPid(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to read pid: %s", err)
	}

	if nspid > 0 {
		err = namespaces.ExecIn(container, nspid, command)
	} else {
		term := namespaces.NewTerminal(os.Stdin, os.Stdout, os.Stderr, container.Tty)
		exitCode, err = startContainer(container, term, dataPath, command)
	}

	if err != nil {
		return fmt.Errorf("failed to exec: %s", err)
	}

	os.Exit(exitCode)
	return nil
}

// startContainer starts the container. Returns the exit status or -1 and an
// error.
//
// Signals sent to the current process will be forwarded to container.
func startContainer(container *libcontainer.Container, term namespaces.Terminal, dataPath string, args []string) (int, error) {
	var (
		cmd  *exec.Cmd
		sigc = make(chan os.Signal, 10)
	)

	signal.Notify(sigc)

	createCommand := func(container *libcontainer.Container, console, rootfs, dataPath, init string, pipe *os.File, args []string) *exec.Cmd {
		cmd = namespaces.DefaultCreateCommand(container, console, rootfs, dataPath, init, pipe, args)
		if logPath != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("log=%s", logPath))
		}
		return cmd
	}

	startCallback := func() {
		go func() {
			for sig := range sigc {
				cmd.Process.Signal(sig)
			}
		}()
	}

	return namespaces.Exec(container, term, "", dataPath, args, createCommand, startCallback)
}
