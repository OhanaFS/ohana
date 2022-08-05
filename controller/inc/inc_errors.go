package inc

import "errors"

var (
	JobCurrentlyRunning          = errors.New("job is currently running")
	JobCurrentlyRunningWarning   = errors.New("job is currently running. warning")
	ErrServerFailed              = errors.New("server failed")
	ErrServerTimeout             = errors.New("server timeout")
	ErrJobFailed                 = errors.New("job failed")
	ErrOrphanedShardsFound       = errors.New("orphaned shards found")
	ErrMissingShardsFound        = errors.New("missing shards found")
	ErrCannotFindNetworkAdaptor  = errors.New("cannot find network adaptor that binds to the IP address in config")
	ErrCannotFindDriveDeviceName = errors.New("cannot find drive device name that bounded to shards location")
	ErrServerNotFound            = errors.New("server not found")
)
