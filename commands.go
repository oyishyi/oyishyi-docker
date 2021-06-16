package main

import (
	"errors"
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

		// check whether type `-it`
		tty := context.Bool("it") // presudo terminal

		// get the resource config
		resourceConfig := subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CPUShare:    context.String("cpushare"),
			CPUAmount:   context.String("cpu"),
		}
		dockerCommands.Run(tty, containerCmd, &resourceConfig)

		return nil
	},
}


