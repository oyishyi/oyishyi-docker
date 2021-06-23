package main

import (
	"errors"
	"fmt"
	"github.com/oyishyi/docker/cgroups/subsystems"
	"github.com/oyishyi/docker/container"
	"github.com/oyishyi/docker/dockerCommands"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// docker init, but cannot be used by user
var initCommand = cli.Command{
	Name:  "init",
	Usage: "init a container",
	Action: func(context *cli.Context) error {
		logrus.Infof("Start initiating...")
		return container.InitProcess()
	},
}

// docker run
var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container",
	Flags: []cli.Flag{
		// integrate -i and -t for convenience
		&cli.BoolFlag{
			Name:  "it",
			Usage: "open an interactive tty(pseudo terminal)",
		},
		// detach
		&cli.BoolFlag{
			Name: "d",
			Usage: "detach the container process",
		},
		// resource limit config
		&cli.StringFlag{
			Name: "m",
			Usage: "limit the memory",
		},&cli.StringFlag{
			Name: "cpu",
			Usage: "limit the cpu amount",
		},&cli.StringFlag{
			Name: "cpushare",
			Usage: "limit the cpu share",
		},
		// volume
		&cli.StringFlag{
			Name: "v",
			Usage: "generate volume",
		},
		// name
		&cli.StringFlag{
			Name: "name",
			Usage: "container name",
		},
	},
	Action: func(context *cli.Context) error {
		args := context.Args()
		if args.Len() == 0 {
			return errors.New("Run what?")
		}

		// transfer from cli.Args to []string
		containerCmd := make([]string, args.Len())        // command
		for index, cmd := range args.Slice() {
			containerCmd[index] = cmd
		}

		// check whether open a pseudo terminal
		tty := context.Bool("it") // presudo terminal
		// check whether detach the process
		detach := context.Bool("d")

		if tty && detach {
			return fmt.Errorf("-it & -d cannot appear together")
		}

		// get the resource config
		resourceConfig := subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CPUShare:    context.String("cpushare"),
			CPUAmount:   context.String("cpu"),
		}
		// get the volume config
		volume := context.String("v")
		// get the container name
		name := context.String("name")
		dockerCommands.Run(tty, containerCmd, &resourceConfig, volume, name)

		return nil
	},
}

var commitCommand = cli.Command{
	Name: "commit",
	Usage: "commit the container into image",
	Action: func(context *cli.Context) error {
		args := context.Args()
		if args.Len() == 0 {
			return errors.New("Commit what?")
		}
		imageName := args.Get(0)
		dockerCommands.CommitContainer(imageName)
		return nil
	},
}

var psCommand = cli.Command{
	Name: "ps",
	Usage: "list all containers",
	Action: func(context *cli.Context) error {
		dockerCommands.ListContainers()
		return nil
	},
}