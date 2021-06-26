package dockerCommands

import (
	"fmt"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func CommitContainer(containerName string, imageName string) {
	mntURL := fmt.Sprintf(container.MntURL, containerName) + "/"
	storeURL := container.ImagesURL + "/" + imageName + ".tar"
	logrus.Infof("stored path: %v", storeURL)
	cmd := exec.Command("tar", "-czf", storeURL, "-C", mntURL, ".")
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Tar folder %s fails %v", mntURL, err)
	}
}