package dockerCommands

import (
	"encoding/json"
	"fmt"
	"github.com/oyishyi/docker/cgroups"
	"github.com/oyishyi/docker/cgroups/subsystems"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// This is the function what `docker run` will call
func Run(tty bool, containerCmd []string, res *subsystems.ResourceConfig, volume string, containerName string, imageName string) {

	// this is "docker init <containerCmd>"
	initProcess, writePipe := container.NewProcess(tty, volume, containerName, imageName)
	if initProcess == nil {
		logrus.Error("create init process fails")
	}

	// start the init process
	if err := initProcess.Start(); err != nil {
		logrus.Error(err)
	}

	// record container info
	containerName, err := recordContainerInfo(initProcess.Process.Pid, containerCmd, containerName, volume)
	if err != nil {
		logrus.Errorf("record container %s fails: %v", containerName, err)
		return
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

	if tty {
		// if pseudo terminal
		if err := initProcess.Wait(); err != nil {
			logrus.Error(err)
		}
		deleteContainerInfo(containerName)
		container.DeleteWorkspace(volume, containerName)
	}

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

func getContainerID(n int) string {
	candidate := "abcdefghijklmnopqrstuvwxyz1234567890"
	rand.Seed(time.Now().UnixNano())

	res := make([]byte, n)
	for i := range res {
		res[i] = candidate[rand.Intn(len(candidate))]
	}
	return string(res)
}

func recordContainerInfo(containerPID int, containerCmd []string, containerName string, volume string) (string, error) {
	pid := containerPID
	id := getContainerID(10)
	name := containerName
	cmd := strings.Join(containerCmd, " ")
	createdTime := time.Now().Format("2006-01-02 15:04:05")
	if name == "" {
		name = id
	}

	containerInfo := &container.ContainerInfo{
		Pid:         strconv.Itoa(pid),
		Id:          id,
		Name:        name,
		Command:     cmd,
		CreatedTime: createdTime,
		Status:      container.RUNNING,
		Volume:      volume,
	}

	// create json string
	jsonBytes, err := json.Marshal(containerInfo)
	jsonStr := string(jsonBytes)
	if err != nil {
		return "", err
	}
	// create path
	dirPath := fmt.Sprintf(container.DefaultInfoLocation, name)
	if err := os.MkdirAll(dirPath, 0622); err != nil {
		logrus.Errorf("mkdir %s fails %v", dirPath, err)
		return "", err
	}
	// create file
	filePath := dirPath + "/" + container.ConfigName
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		logrus.Errorf("create file %s fails %v", filePath, err)
		return "", err
	}
	// write json string
	if _, err := file.WriteString(jsonStr); err != nil {
		logrus.Errorf("write file %s fails: %v", filePath, err)
		return "", err
	}
	return name, nil
}

func deleteContainerInfo(containerName string) {
	dirPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirPath); err != nil {
		logrus.Errorf("delete container config %s fails %v", dirPath, err)
	}
}
