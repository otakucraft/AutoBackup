package rsync

// Manager threaded rsync instance manager.
// Controls how many instances there are running and processes them.
type Manager struct {
	instances map[*Instance]bool
	addInstance chan *Instance
	removeInstance chan *Instance
}

// Start manager main loop, receives instances and starts them, adds or removes them from the instance array.
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

// InitInstance start a rsync instance and wait for it to finish.
func (rManager *Manager) InitInstance(instance *Instance) {
	rManager.addInstance <- instance
	instance.Run()
	rManager.removeInstance <- instance
}

// GetInstances instance getter if needed.
func (rManager *Manager) GetInstances() map[*Instance]bool {
	return rManager.instances
}
