package dockerCommands

import (
	"encoding/json"
	"fmt"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

func ListContainers() {
	dirPath := strings.Split(container.DefaultInfoLocation, "%s")[0]
	dirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logrus.Errorf("read dir %s fails: %v", dirPath, err)
		return
	}

	// stored in a slice
	var containerInfos []*container.ContainerInfo
	for _, dir := range dirs {
		tempContainer, err := getContainerInfo(dir)
		if err != nil {
			logrus.Error(err)
		}
		containerInfos = append(containerInfos, tempContainer)
	}

	// output
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	// store output in the writer
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, containerInfo := range containerInfos {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			containerInfo.Id,
			containerInfo.Name,
			containerInfo.Pid,
			containerInfo.Status,
			containerInfo.Command,
			containerInfo.CreatedTime,
		)
	}
	// print the output stored onto the os.stdout
	if err := w.Flush(); err != nil {
		logrus.Errorf("flush error: %v", err)
		return
	}
}

func getContainerInfo(dir os.FileInfo)(*container.ContainerInfo, error) {
	containerName := dir.Name()
	configPath := fmt.Sprintf(container.DefaultInfoLocation, containerName) + container.ConfigName
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read container config %s fails: %v", configPath, err)
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		return nil, fmt.Errorf("unmarshal json file fails: %v", err)
	}

	return &containerInfo, nil
}
