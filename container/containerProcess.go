package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

func NewProcess(tty bool, volume string) (*exec.Cmd, *os.File) {

	readPipe, writePipe, err := os.Pipe()

	if err != nil {
		logrus.Errorf("New Pipe Error: %v", err)
		return nil, nil
	}
	// create a new command which run itself
	// the first arguments is `init` which is in the "container/init.go" file
	// so, the <cmd> will be interpret as "docker init <containerCmd>"
	cmd := exec.Command("/proc/self/exe", "init")

	// new namespaces, thanks to Linux
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	// this is what presudo terminal means
	// link the container's stdio to os
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.ExtraFiles = []*os.File{readPipe}

	imagesRootURL := "./images/"
	mntURL := "./mnt/"
	newWorkspace(imagesRootURL, mntURL, volume)
	cmd.Dir = mntURL

	return cmd, writePipe
}

func newWorkspace(imagesRootURL string, mntURL string, volume string) {
	createReadOnlyLayer(imagesRootURL)
	createWriteLayer(imagesRootURL)
	createMountPoint(imagesRootURL, mntURL)
	createVolume(imagesRootURL, mntURL, volume)
}

func createReadOnlyLayer(imagesRootURL string) {
	readOnlyLayerURL := imagesRootURL + "busybox/"
	imageTarURL := imagesRootURL + "busybox.tar"
	isExist, err := pathExist(readOnlyLayerURL)
	if err != nil {
		logrus.Infof("fail to judge whether path exist: %v", err)
	}
	if isExist == false {
		if err := os.Mkdir(readOnlyLayerURL, 0777); err != nil {
			logrus.Errorf("fail to create dir %s: %v", readOnlyLayerURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", imageTarURL, "-C", readOnlyLayerURL).CombinedOutput(); err != nil {
			logrus.Errorf("fail to untar %s: %v", imageTarURL, err)
		}
	}
}

func createWriteLayer(imagesRootURL string) {
	writeLayerURL := imagesRootURL + "writeLayer/"
	if err := os.Mkdir(writeLayerURL, 0777); err != nil {
		logrus.Errorf("fail to create dir %s: %v", writeLayerURL, err)
	}
}

func createMountPoint(imagesRootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		logrus.Errorf("fail to create dir %s: %v", mntURL, err)
	}

	// mount the readOnly layer and writeLayer on the mntURL
	dirs := "dirs=" + imagesRootURL + "writeLayer:" + imagesRootURL + "busybox"
	fmt.Println(dirs)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Error(err)
	}
}

func createVolume(imagesRootURL string, mntURL string, volume string) {
	if volume == "" {
		return
	}
	// extract url from volume input
	volumeURls := strings.Split(volume, ":")
	if len(volumeURls) == 2 && volumeURls[0] != "" && volumeURls[1] != "" {
		mountVolume(imagesRootURL, mntURL, volumeURls)
		logrus.Infof("volume created: %q", volumeURls)
	} else {
		logrus.Warn("volume parameter not valid")
	}

}

func mountVolume(imagesRootURL string, mntURL string, volumeURls []string) {
	hostVolumeURL := volumeURls[0]
	containerVolumeURL := path.Join(mntURL, volumeURls[1])
	// create host dir, which store the real data
	if err:= os.Mkdir(hostVolumeURL, 0777); err != nil {
		logrus.Errorf("create host volume directory %s fails: %v", containerVolumeURL, err)
	}

	// create container dir, which is a mount point
	// currently not in container(container root is not "/"), so mntURL prefix is needed
	if err:= os.Mkdir(containerVolumeURL, 0777); err != nil {
		logrus.Errorf("create container volume %s fails: %v", containerVolumeURL, err)
	}

	dirs :="dirs=" + hostVolumeURL
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err:= cmd.Run(); err != nil {
		logrus.Errorf("mount container volume %s fails: %v", containerVolumeURL, err)
	}


}

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// delete AUFS (delete write layer)
func DeleteWorkspace(imagesRootURL string, mntURL string, volume string) {
	deleteMountWithVolume(mntURL, volume)
	deleteWriteLayer(imagesRootURL)
}



func deleteMountWithVolume(mntURL string, volume string) {

	// umount volume mount
	volumeURls := strings.Split(volume, ":")
	if len(volumeURls) == 2 && volumeURls[0] != "" && volumeURls[1] != "" {
		volumeMountURL := path.Join(mntURL, volumeURls[1])
		cmd := exec.Command("umount", volumeMountURL)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			logrus.Errorf("umount volume mount %s fails: %v", volumeMountURL, err)
		}
	}

	// umount container mount
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("umount container mount %s fails: %v", mntURL, err)
	}

	// delete whole container mount
	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", mntURL, err)
	}
}

func deleteWriteLayer(imagesRootURL string) {
	writeLayerURL := imagesRootURL + "writeLayer/"
	if err:= os.RemoveAll(writeLayerURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", writeLayerURL, err)
	}
}