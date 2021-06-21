package dockerCommands

import (
	"github.com/sirupsen/logrus"
	"os/exec"
)

func CommitContainer(imageName string) {
	mntURL := "./mnt"
	storeURL := "./images/" + imageName + ".tar"
	logrus.Infof("stored path: %v", storeURL)
	cmd := exec.Command("tar", "-czf", storeURL, "-C", mntURL, ".")
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Tar folder %s fails %v", mntURL, err)
	}
}