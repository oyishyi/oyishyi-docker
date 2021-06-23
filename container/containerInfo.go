package container

type ContainerInfo struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"created_time"`
	Status      string `json:"status"`
}

var (
	RUNNING string = "running"
	STOP string = "stopped"
	Exit string = "exited"
	DefaultInfoLocation = "/var/run/oyishyi-docker/%s/"
	ConfigName = "config.json"
)