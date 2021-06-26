package dockerCommands

import (
	"encoding/json"
	"fmt"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"syscall"
)

func StopContainer(containerName string) {
	containerInfo, err := container.GetContainerInfoByName(containerName)
	if err != nil {
		logrus.Errorf("get containerInfo of %s fails: %v", containerName, err)
		return
	}

	pid, err := strconv.Atoi(containerInfo.Pid)
	if err != nil {
		logrus.Errorf("convert Pid to int fails: %v", err)
		return
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		logrus.Errorf("stopping container fais: %v", err)
		return
	}

	containerInfo.Status = container.STOP
	containerInfo.Pid = ""
	containerInfoBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("encode containerInfo to UTF-8(bytes) fails: %v", err)
		return
	}
	dirPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirPath + container.ConfigName
	if err := ioutil.WriteFile(configFilePath, containerInfoBytes, 0622); err != nil {
		logrus.Errorf("change config file %s fails: %v", configFilePath, err)
		return
	}
	return
}