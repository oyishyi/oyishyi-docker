package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CPUAmountSubsystem struct {
}

// return the name of the subsystem
func (c *CPUAmountSubsystem) Name() string {
	return "cpuset"
}

// set the cpu amount of a cgroup
func (c *CPUAmountSubsystem) Set(cgroupPath string, res *ResourceConfig) error  {
	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true); err != nil {
		return err
	} else {
		targetFilePath := path.Join(subsystemCgroupPath, "cpuset.cpus")
		if res.CPUAmount != "" {
			fmt.Println("no way")
			if err := ioutil.WriteFile(targetFilePath, []byte(res.CPUAmount), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu amount fail: %v", err)
			}
		}
	}
	return nil
}

func (css *CPUAmountSubsystem) AddProcess(cgroupPath string, pid int) error {

	if subsystemCgroupPath, err := GetCgroupPath(css.Name(), cgroupPath, false); err != nil {
		return err
	} else {
		targetFilePath := path.Join(subsystemCgroupPath, "tasks")
		if err := ioutil.WriteFile(targetFilePath, []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("cgroup add process fail: %v", err)
		}
	}
	return nil
}

func (css *CPUAmountSubsystem) RemoveCgroup(cgroupPath string) error {
	if SubsystemCgroupPath, err := GetCgroupPath(css.Name(), cgroupPath, false); err != nil {
		return err
	} else {
		return os.Remove(SubsystemCgroupPath)
	}
}