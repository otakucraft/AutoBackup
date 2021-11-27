package rsync

type Manager struct {
	instances map[*Instance]bool
	addInstance chan *Instance
	removeInstance chan *Instance
}

func (rManager *Manager) Start() {
	for {
		select {
		case i := <- rManager.addInstance:
			rManager.instances[i] = true
		case i := <- rManager.removeInstance:
			delete(rManager.instances, i)
		}
	}
}

func (rManager *Manager) InitInstance(instance *Instance) {
	rManager.addInstance <- instance
	instance.Run()
	rManager.removeInstance <- instance
}

func (rManager *Manager) GetInstances() map[*Instance]bool {
	return rManager.instances
}
