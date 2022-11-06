package util

import "time"

type BackoffEntry struct {
	AfterTime time.Duration
	TimeApart time.Duration
}

type Backoff struct {
	Entries     []BackoffEntry
	activeEntry int
	startTime   time.Time
	lastTime    time.Time
}

func (b *Backoff) Start() {
	b.startTime = time.Now()
	b.activeEntry = 0
	b.lastTime = b.startTime
}

func (b *Backoff) Ready() bool {
	now := time.Now()
	sinceStart := now.Sub(b.startTime)
	sinceLast := now.Sub(b.lastTime)

	if b.activeEntry+1 < len(b.Entries) {
		nextEntry := b.Entries[b.activeEntry+1]
		if sinceStart >= nextEntry.AfterTime {
			b.activeEntry++
			b.lastTime = now
			return true
		}
	}

	return sinceLast >= b.Entries[b.activeEntry].TimeApart
}
