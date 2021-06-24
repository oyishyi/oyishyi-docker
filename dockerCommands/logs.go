package dockerCommands

import (
	"fmt"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func LogContainer(containerName string) {
	// find the dir path
	dirPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	filePath := dirPath + container.ContainerLogName

	// open the log file
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		logrus.Errorf("open file %s fails: %v", filePath, err)
		return
	}

	// read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read file %s fails: %v", filePath, err)
		return
	}

	// print log to stdout
	fmt.Fprint(os.Stdout, string(content))
}