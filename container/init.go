package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

// already in container
// initiate the container
func InitProcess() error {

	// read command from pipe, will plug if write side is not ready
	containerCmd := readCommand()
	if containerCmd == nil || len(containerCmd) == 0 {
		return fmt.Errorf("Init process fails, containerCmd is nil")
	}

	// setup all mount commands
	if err:= setupMount(); err != nil {
		logrus.Errorf("setup mount fails: %v", err)
		return err
	}

	// look for the path of container command
	// so we don't need to type "/bin/ls", but "ls"
	commandPath, err := exec.LookPath(containerCmd[0])
	if err != nil {
		logrus.Errorf("initProcess look path fails: %v", err)
		return err
	}

	// log commandPath info
	// if you type "ls", it will be "/bin/ls"
	logrus.Infof("Find commandPath: %v", commandPath)
	if err := syscall.Exec(commandPath, containerCmd, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}

func readCommand() []string {
	// 3 is the index of readPipe
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("read pipe fails: %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}


// integration of all mount commands
func setupMount() error {

	// ensure that container mount and parent mount has no shared propagation
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		logrus.Errorf("mount / fails: %v", err)
		return err
	}

	// get current directory
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	logrus.Infof("current location is: %v", pwd)
	// use current directory as the root
	if  err:= pivotRoot(pwd); err != nil {
		logrus.Errorf("pivot root fails: %v", err)
		return err
	}

	// mount proc filesystem
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("mount /proc fails: %v", err)
		return err
	}

	// mount tmpfs
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755"); err != nil {
		logrus.Errorf("mount /dev fails: %v", err)
		return err
	}

	return nil
}

// change the container rootfs to image rootfs
func pivotRoot(root string) error {

	// what it does?
	// remember root is just a parameter now,it's not the rootfs, it's what we want to create
	// this command ensure that root is a mount point which bind itself
	if err:= syscall.Mount(root, root, "bind", syscall.MS_BIND | syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("remount root fails: %v", err)
	}

	// create the putOld directory to store old rootfs
	putOld := path.Join(root, ".put_old")
	if err:= os.Mkdir(putOld, 0777); err != nil {
		return fmt.Errorf("create putOld directory fails: %v", err)
	}

	// pivot old root mount to putOld
	// and mount the first parameter as the new root mount
	// which means, '/.put_old/' is exactly the old rootfs
	// the first parameter must be a mount point, that's why we remount root itself at the beginning
	if err := syscall.PivotRoot(root, putOld); err != nil {
		return fmt.Errorf("pivot_root fails: %v", err)
	}

	// chdir do exactly the same as cd. chdir is a syscall, cd is a program
	// change to root directory
	if err:= syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir fails: %v", err)
	}

	// after the previous process, the current filesystem is the new root
	// the old filesystem is .put_old
	// finally, we need to unmount the old root mount before remove it
	// change the putOld dir, because we are in the new rootfs now
	// the root became "/"
	putOld = path.Join("/", ".put_old")
	if err:= syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount fails: %v", err)
	}

	// remove the old mount point
	return os.Remove(putOld)
}