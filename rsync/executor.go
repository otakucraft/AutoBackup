package rsync

var RsyncExecutor *Executor

// Executor structs that will contain the rsync instance manager.
type Executor struct {
	rManager *Manager
}

func NewExecutor() *Executor {
	return &Executor{}
}

// Start init the rsync manager variables and start it.
func (executor *Executor) Start() {
	executor.rManager = &Manager{
		instances:      make(map[*Instance]bool),
		addInstance:    make(chan *Instance),
		removeInstance: make(chan *Instance),
	}
	go executor.rManager.Start()
}

// StartInstance given a rsync instance, start it on a different thread but still inside the manager.
func (executor *Executor) StartInstance(instance *Instance) {
	go executor.rManager.InitInstance(instance)
}

// StopInstance not implemented yet but to kill the instance.
func (executor *Executor) StopInstance(instance *Instance) {
	instance.Stop()
}

// GetManager manager getter if needed.
func (executor *Executor) GetManager() *Manager {
	return executor.rManager
}
