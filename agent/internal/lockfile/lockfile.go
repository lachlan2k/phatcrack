package lockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Lockfile struct {
	path string

	id          string
	mu          sync.Mutex
	stopWriting context.CancelFunc
	created     time.Time
}

const (
	writeInterval = time.Second

	// If a lockfile hasn't been updated in this long, consider the owner dead
	staleAge = 10 * time.Second

	// Time to wait after we delete a stale lockfile
	// This prevents the race condition where we have 2 writers who both delete the lockfile, and writer B deletes writer A's new lockfile
	afterStaleDelay = 3 * time.Second
)

type lockdata struct {
	Created int64
	Updated int64
	ID      string
}

func (data lockdata) isStale() bool {
	updatedT := time.Unix(data.Updated, 0)

	return time.Since(updatedT) > staleAge
}

func New(path string) Lockfile {
	return Lockfile{
		path: path,
	}
}

func (l *Lockfile) MyID() string {
	return l.id
}

func (l *Lockfile) Acquire(ctx context.Context) error {
	for {
		err := l.tryAcquire()
		if err == nil {
			// success
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}

func (l *Lockfile) tryAcquire() error {
	// Check if there's already a lockfile we can read
	data, err := l.readData()
	if err == nil {
		// Remove if stale, or exit if its live
		if data.isStale() {
			os.Remove(l.path)
			// Sleep to prevent the race condition of 2 writers deleting simultaneously, then write B deletes writer A's lockfile
			time.Sleep(afterStaleDelay)
		} else {
			return errors.New("lock is held by another process")
		}
	}
	// If we fail to read it that probably means the lockfile doesn't exist, in which case let's proceed

	l.mu.Lock()
	defer l.mu.Unlock()

	l.id = uuid.NewString()
	l.created = time.Now()

	// Try to claim the lockfile
	// Write, read it back to make sure we have the lock
	err = l.write(true)
	if err != nil {
		l.id = ""
		return err
	}
	err = l.ensureHeld()
	if err != nil {
		l.id = ""
		return err
	}

	// Lockfile obtained
	l.startWriting()
	return nil
}

func (l *Lockfile) readData() (*lockdata, error) {
	buff, err := os.ReadFile(l.path)
	if err != nil {
		return nil, err
	}

	data := &lockdata{
		Created: 0,
		Updated: 0,
		ID:      "",
	}
	// If it fails, that's fine (such as corrupt content of lockfile)
	json.Unmarshal(buff, &data)
	return data, nil
}

func (l *Lockfile) ensureHeld() error {
	// If someone else is writing, their write will appear in N seconds, so we check after N+1
	time.Sleep(writeInterval + time.Second)
	data, err := l.readData()
	if err != nil {
		return err
	}
	if l.id != data.ID {
		return errors.New("ID in lockfile doesn't match our ID")
	}
	return nil
}

func (l *Lockfile) write(create bool) error {
	flags := os.O_WRONLY
	if create {
		flags |= os.O_CREATE | os.O_EXCL
	}

	f, err := os.OpenFile(l.path, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	d := lockdata{
		Created: l.created.Unix(),
		Updated: time.Now().Unix(),
		ID:      l.id,
	}

	// buff, err := json.Marshal(d)
	err = json.NewEncoder(f).Encode(d)
	if err != nil {
		return err
	}

	return f.Sync()
}

func (l *Lockfile) startWriting() {
	ctx, cancel := context.WithCancel(context.Background())
	l.stopWriting = cancel
	go l.writeLoop(ctx)
}

func (l *Lockfile) writeLoop(ctx context.Context) {
	for {
		l.mu.Lock()
		err := l.write(false)
		if err != nil {
			logrus.WithError(err).Warn("Unexpected problem when writing lockfile")
		}
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return
		case <-time.After(writeInterval):
		}
	}
}

func (l *Lockfile) Unlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.stopWriting()
	l.delete()

	l.stopWriting = nil
}

func (l *Lockfile) delete() error {
	return os.Remove(l.path)
}
