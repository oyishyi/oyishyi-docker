package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ContainerInfo struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"created_time"`
	Status      string `json:"status"`
}

var (
	RUNNING string = "running"
	STOP string = "stopped"
	Exit string = "exited"
	DefaultInfoLocation = "/var/run/oyishyi-docker/%s/"
	ConfigName = "config.json"
	ContainerLogName = "container.log"
)

func GetContainerInfoByName(containerName string) (*ContainerInfo, error) {
	dirPath := fmt.Sprintf(DefaultInfoLocation, containerName)
	configPath := dirPath + ConfigName

	contentBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ContainerInfo
	if err := json.Unmarshal(contentBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}