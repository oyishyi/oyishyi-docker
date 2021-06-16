package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// as the function name shows, find the root path of hierarchy
func FindHierarchyMountRootPath(subsystemName string) string  {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		// find whether "subsystemName" appear in the last field
		// if so, then the fifth field is the path
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystemName {
				return fields[4]
			}
		}
	}
	return ""
}

// get the absolute path of a cgroup
func GetCgroupPath(subsystemName string, cgroupPath string, autoCreate bool) (string, error)  {
	cgroupRootPath := FindHierarchyMountRootPath(subsystemName)
	expectedPath := path.Join(cgroupRootPath, cgroupPath)
	// find the cgroup or create a new cgroup
	if _, err := os.Stat(expectedPath); err == nil  || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(expectedPath, 0755); err != nil {
				return "", fmt.Errorf("error when create cgroup: %v", err)
			}
		}
		return expectedPath, nil
	} else {
		return "", fmt.Errorf("cgroup path error: %v", err)
	}
}