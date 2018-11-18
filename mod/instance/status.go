package instance

// Status - get current status
type Status = string

const (
	statusReady   Status = "STATUS_READY"
	statusRunning Status = "STATUS_RUNNING"
	statusStopped Status = "STATUS_STOPPED"
)

func setStatus(inst *Instance, status Status) {
	inst.Lock()
	defer inst.Unlock()

	inst.status = status
}

func getStatus(inst *Instance) Status {
	inst.RLock()
	defer inst.RUnlock()

	return inst.status
}
