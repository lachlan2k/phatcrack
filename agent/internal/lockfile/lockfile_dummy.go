package lockfile

import "time"

type LockfileDummy struct{}

func (l *LockfileDummy) AcquireWithTimeout(timeout time.Duration) error {
	return nil
}

func (l *LockfileDummy) Unlock() {
	// nop
}

func NewDummy() *LockfileDummy {
	return &LockfileDummy{}
}
