package main

import (
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	logPath  string
	dataPath string
)

func initLogging() error {
	if logPath != "" {
		if err := openLog(logPath); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	app := kingpin.New("nsinit", "libcontainer reference application").Version("0.1")
	app.Flag("log_path", "Path to log file.").OverrideDefaultFromEnvar("log").Required().Dispatch(initLogging).StringVar(&logPath)
	app.Flag("data_path", "Path to data directory.").OverrideDefaultFromEnvar("data_path").Required().ExistingDirVar(&dataPath)

	execCommand := app.Command("exec", "Execute a new command inside a container.")
	execCommandCommand := execCommand.Arg("command", "Command to execute.").Required().Strings()
	execCommand.Dispatch(func() error { return execAction(*execCommandCommand) })

	initCommand := app.Command("init", "Runs the init process inside the namespace.")
	initCommandConsole := initCommand.Flag("console", "Path to console file.").OverrideDefaultFromEnvar("console").Required().String()
	initCommandRawPipeFd := initCommand.Flag("pipe", "FD to pipe output to.").OverrideDefaultFromEnvar("pipe").Required().Int()
	initCommandCommand := initCommand.Arg("command", "Command to execute.").Required().Strings()
	initCommand.Dispatch(func() error { return initAction(*initCommandCommand, *initCommandConsole, *initCommandRawPipeFd) })

	statsCommand := app.Command("stats", "Display statistics for the container.")
	statsCommand.Dispatch(statsAction)

	specCommand := app.Command("spec", "Display the container specification.")
	specCommand.Dispatch(specAction)

	nsenterCommand := app.Command("nsenter", "Init process for entering an existing namespace.")
	nsenterCommandPid := nsenterCommand.Arg("pid", "PID of existing container.").Required().Int()
	nsenterCommandProcessLabel := nsenterCommand.Arg("label", "Process label.").Required().String()
	nsenterCommandContainerJSON := nsenterCommand.Arg("container-json", "Container JSON.").Required().String()
	nsenterCommandCmd := nsenterCommand.Arg("cmd", "Command to execute in existing namespace.").Required().Strings()
	nsenterCommand.Dispatch(func() error {
		return nsenterAction(*nsenterCommandPid, *nsenterCommandProcessLabel, *nsenterCommandContainerJSON, *nsenterCommandCmd)
	})

	if kingpin.MustParse(app.Parse(os.Args[1:])) == "" {
		kingpin.Fatalf("command not provided")
	}
}
