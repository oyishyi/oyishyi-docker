package dockerCommands

import (
	"github.com/oyishyi/oyishyi-docker/cgroups"
	"github.com/oyishyi/oyishyi-docker/cgroups/subsystems"
	"github.com/oyishyi/oyishyi-docker/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

// This is the function what `docker run` will call
func Run(tty bool, containerCmd []string, res *subsystems.ResourceConfig) {

	// this is "docker init <containerCmd>"
	initProcess, writePipe := container.NewProcess(tty)

	// start the init process
	if err := initProcess.Start(); err != nil{
		logrus.Error(err)
	}

	// create container manager to control resource config on all hierarchies
	// this is the cgroupPath
	cm := cgroups.NewCgroupManager("oyishyi-docker-first-cgroup")
	defer cm.Remove()

	if err := cm.Set(res); err != nil {
		logrus.Error(err)
	}
	if err := cm.AddProcess(initProcess.Process.Pid); err != nil {
		logrus.Error(err)
	}

	// send command to write side
	// will close the plug
	sendInitCommand(containerCmd, writePipe)

	if err := initProcess.Wait(); err != nil {
		logrus.Error(err)
	}

	imagesRootURL := "./images/"
	mntURL := "./mnt/"
	container.DeleteWorkspace(imagesRootURL, mntURL)

	os.Exit(-1)
}

func sendInitCommand(containerCmd []string, writePipe *os.File) {
	cmdString := strings.Join(containerCmd, " ")
	logrus.Infof("whole init command is: %v", cmdString)
	if _, err := writePipe.WriteString(cmdString); err != nil {
		logrus.Error(err)
	}
	if err := writePipe.Close(); err != nil {
		logrus.Error(err)
	}
}