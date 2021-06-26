package dockerCommands

import (
	"github.com/oyishyi/docker/container"
	_ "github.com/oyishyi/docker/setns"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func ExecContainer(containerName string, containerCmd []string) {
	containerInfo, err := container.GetContainerInfoByName(containerName)
	if err != nil {
		logrus.Errorf("get container info of %s fails: %v", containerName, err)
		return
	}
	pid := containerInfo.Pid
	containerCmdStr := strings.Join(containerCmd, " ")
	logrus.Infof("container pid: %s", pid)
	logrus.Infof("container cmd: %s", containerCmdStr)

	// run docker exec again
	// with env this time, so that bunch of cgo codes will execute
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//cmd.Env = []string{fmt.Sprintf("%s=%s", ENV_EXEC_PID, pid), fmt.Sprintf("%s=%s", ENV_EXEC_CMD, containerCmdStr)}

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, containerCmdStr)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("run the second docker exec with env fails: %v", err)
		return
	}

}
