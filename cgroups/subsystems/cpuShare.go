package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CPUShareSubsystem struct {
}

// return the name of the subsystem
func (cs *CPUShareSubsystem) Name() string {
	return "cpu"
}

// set the cpu share of a cgroup
func (cs *CPUShareSubsystem) Set(cgroupPath string, res *ResourceConfig) error  {
	if subsystemCgroupPath, err := GetCgroupPath(cs.Name(), cgroupPath, true); err != nil {
		return err
	} else {
		targetFilePath := path.Join(subsystemCgroupPath, "cpu.shares")
		if res.CPUShare != "" {
			if err := ioutil.WriteFile(targetFilePath, []byte(res.CPUShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail: %v", err)
			}
		}
	}
	return nil
}

func (cs *CPUShareSubsystem) AddProcess(cgroupPath string, pid int) error {
	if SubsystemCgroupPath, err := GetCgroupPath(cs.Name(), cgroupPath, false); err != nil {
		return err
	} else {
		targetFilePath := path.Join(SubsystemCgroupPath, "tasks")
		if err := ioutil.WriteFile(targetFilePath, []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("cgroup add process fail: %v", err)
		}
	}
	return nil
}

func (cs *CPUShareSubsystem) RemoveCgroup(cgroupPath string) error {
	if SubsystemCgroupPath, err := GetCgroupPath(cs.Name(), cgroupPath, false); err != nil {
		return err
	} else {
		return os.Remove(SubsystemCgroupPath)
	}
}