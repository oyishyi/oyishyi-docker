package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// mount related variable
var (
	ImagesURL = "./runtime"
	ReadOnlyLayerURL = "./runtime/readOnlyLayer/%s"
	MntURL = "./runtime/mnt/%s"
	WriteLayerURL = "./runtime/writeLayer/%s"
)


func NewProcess(tty bool, volume string, containerName string, imageName string) (*exec.Cmd, *os.File) {

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
		// if tty
		// attach to stdio
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// if detach
		// generate log file
		dirPath := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err:= os.MkdirAll(dirPath, 0622); err != nil {
			logrus.Errorf("mkdir %s fails %v", dirPath, err)
			return nil, nil
		}
		filePath := dirPath + ContainerLogName
		file, err := os.Create(filePath)
		if err != nil {
			logrus.Errorf("create file %s fails: %v", filePath, err)
			return nil, nil
		}
		cmd.Stdout = file
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	//imagesRootURL := "./images/"
	//mntURL := "./mnt/"
	newWorkspace(volume, containerName, imageName)
	cmd.Dir = fmt.Sprintf(MntURL, containerName)

	return cmd, writePipe
}

func newWorkspace(volume string, containerName string, imageName string) {
	createReadOnlyLayer(imageName)
	createWriteLayer(containerName)
	createMountPoint(containerName, imageName)
	createVolume(containerName, volume)
}

func createReadOnlyLayer(imageName string) {
	imageTarURL := ImagesURL + "/" + imageName + ".tar"
	readOnlyLayerURL := fmt.Sprintf(ReadOnlyLayerURL, imageName) + "/"
	isExist, err := pathExist(readOnlyLayerURL)
	if err != nil {
		logrus.Infof("fail to judge whether path exist: %v", err)
	}
	if isExist == false {
		if err := os.MkdirAll(readOnlyLayerURL, 0777); err != nil {
			logrus.Errorf("fail to create dir %s: %v", readOnlyLayerURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", imageTarURL, "-C", readOnlyLayerURL).CombinedOutput(); err != nil {
			logrus.Errorf("fail to untar %s: %v", imageTarURL, err)
		}
	}
}

func createWriteLayer(containerName string) {
	writeLayerURL := fmt.Sprintf(WriteLayerURL, containerName)
	fmt.Println(writeLayerURL)
	if err := os.MkdirAll(writeLayerURL, 0777); err != nil {
		logrus.Errorf("fail to create dir %s: %v", writeLayerURL, err)
	}
}

func createMountPoint(containerName string, imageName string) {
	mntURL := fmt.Sprintf(MntURL, containerName)
	if err := os.MkdirAll(mntURL, 0777); err != nil {
		logrus.Errorf("fail to create dir %s: %v", mntURL, err)
	}

	writeLayerURL := fmt.Sprintf(WriteLayerURL, containerName)
	imageTmpLocationURL := fmt.Sprintf(ReadOnlyLayerURL, imageName)
	// mount the readOnly layer and writeLayer on the mntURL
	dirs := "dirs=" + writeLayerURL + ":" + imageTmpLocationURL


	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount dirs %s fails: %v", dirs, err)
		return
	}
	logrus.Infof("mount dirs: %s", dirs)
}

func createVolume(containerName string, volume string) {
	if volume == "" {
		return
	}
	// extract url from volume input
	volumeURls := strings.Split(volume, ":")
	if len(volumeURls) == 2 && volumeURls[0] != "" && volumeURls[1] != "" {
		mountVolume(containerName, volumeURls)
		logrus.Infof("volume created: %q", volumeURls)
	} else {
		logrus.Warn("volume parameter not valid")
	}

}

func mountVolume(containerName string, volumeURls []string) {
	hostVolumeURL := volumeURls[0]
	containerVolumeURL := fmt.Sprintf(MntURL, containerName) + "/" + volumeURls[1]
	// create host dir, which store the real data
	if err:= os.MkdirAll(hostVolumeURL, 0777); err != nil {
		logrus.Errorf("create host volume directory %s fails: %v", containerVolumeURL, err)
	}

	// create container dir, which is a mount point
	// currently not in container(container root is not "/"), so mntURL prefix is needed
	if err:= os.MkdirAll(containerVolumeURL, 0777); err != nil {
		logrus.Errorf("create container volume %s fails: %v", containerVolumeURL, err)
	}

	dirs :="dirs=" + hostVolumeURL
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	if err:= cmd.Run(); err != nil {
		logrus.Errorf("mount container volume %s fails: %v", containerVolumeURL, err)
	}
	logrus.Infof("volume dirs: %s", dirs)
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
func DeleteWorkspace(volume string, containerName string) {
	deleteMountWithVolume(volume, containerName)
	deleteWriteLayer(containerName)
}


func deleteMountWithVolume(volume string, containerName string) {

	mntURL := fmt.Sprintf(MntURL, containerName)

	// umount volume mount
	volumeURls := strings.Split(volume, ":")
	if len(volumeURls) == 2 && volumeURls[0] != "" && volumeURls[1] != "" {

		volumeMountURL := mntURL + "/" + volumeURls[1]
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

func deleteWriteLayer(containerName string) {
	writeLayerURL := fmt.Sprintf(WriteLayerURL, containerName)
	if err:= os.RemoveAll(writeLayerURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", writeLayerURL, err)
	}
}