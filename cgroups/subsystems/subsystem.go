package subsystems

type ResourceConfig struct {
	MemoryLimit string // memory limit
	CPUShare    string // cpu share
	CPUAmount   string // cpu amount
}

type Subsystem interface {
	// return the name of which type of subsystem
	Name() string
	// set a resource limit on a cgroup
	Set(cgroupPath string, res *ResourceConfig) error
	// add a processs with the pid to a group
	AddProcess(cgroupPath string, pid int) error
	// remove a cgroup
	RemoveCgroup(cgroupPath string) error
}

// instance of a subsystems
var SubsystemsInstance = []Subsystem{
	&CPUShareSubsystem{},
	&CPUAmountSubsystem{},
	&MemorySubsystem{},
}
