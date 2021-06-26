package dockerCommands

import (
	"fmt"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"os"
)

func RemoveContainer(containerName string) {
	containerInfo, err := container.GetContainerInfoByName(containerName)
	if err != nil {
		logrus.Errorf("get container info of %s fails: %v", containerName, err)
		return
	}

	if containerInfo.Status != container.STOP {
		logrus.Errorf("couldn't remove not stopped container")
		return
	}

	dirPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err:= os.RemoveAll(dirPath); err != nil {
		logrus.Errorf("remove dir %s fails: %v", dirPath, err)
		return
	}
}
