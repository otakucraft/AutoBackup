package rsync

var RsyncExecutor *Executor

type Executor struct {
	rManager *Manager
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (executor *Executor) Start() {
	executor.rManager = &Manager{
		instances:      make(map[*Instance]bool),
		addInstance:    make(chan *Instance),
		removeInstance: make(chan *Instance),
	}
	go executor.rManager.Start()
}

func (executor *Executor) StartInstance(instance *Instance) {
	go executor.rManager.InitInstance(instance)
}

func (executor *Executor) StopInstance(instance *Instance) {
	instance.Stop()
}

func (executor *Executor) GetManager() *Manager {
	return executor.rManager
}
