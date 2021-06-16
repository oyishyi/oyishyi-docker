package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func NewProcess(tty bool) (*exec.Cmd, *os.File) {

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
	newWorkspace(imagesRootURL, mntURL)
	cmd.Dir = mntURL

	return cmd, writePipe
}

func newWorkspace(imagesRootURL string, mntURL string) {
	createReadOnlyLayer(imagesRootURL)
	createWriteLayer(imagesRootURL)
	createMountPoint(imagesRootURL, mntURL)
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
func DeleteWorkspace(imagesRootURL string, mntURL string) {
	deleteMount(mntURL)
	deleteWriteLayer(imagesRootURL)
}

func deleteMount(mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("umount fails: %v", err)
	}
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
