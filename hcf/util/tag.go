package util

import (
	"github.com/df-mc/atomic"
	"time"
)

// TagFunc is a function that is called when a tag is set.
type TagFunc func(t *Tag)

// UntagFunc is a function that is called when a tag expires.
type UntagFunc = TagFunc

// Tag is a tag that can be used to cancel a task.
type Tag struct {
	expiration atomic.Value[time.Time]

	tag   TagFunc
	untag UntagFunc

	c chan struct{}
}

// NewTag returns a new tag.
func NewTag(t TagFunc, f UntagFunc) *Tag {
	return &Tag{
		tag:   t,
		untag: f,

		c: make(chan struct{}),
	}
}

// Active returns true if the tag is active.
func (t *Tag) Active() bool {
	return t.expiration.Load().After(time.Now())
}

// Remaining returns the remaining time of the tag.
func (t *Tag) Remaining() time.Duration {
	return time.Until(t.expiration.Load())
}

// Set adds a duration to the tag.
func (t *Tag) Set(d time.Duration) {
	if t.Active() {
		t.Cancel()
	}
	t.c = make(chan struct{})

	if t.tag != nil {
		t.tag(t)
	}
	go func() {
		select {
		case <-time.After(d):
			if t.untag != nil {
				t.untag(t)
			}
		case <-t.c:
			return
		}
	}()
	t.expiration.Store(time.Now().Add(d))
}

// Reset resets the tag.
func (t *Tag) Reset() {
	if t.Active() {
		t.Cancel()
	}
	if t.untag != nil {
		t.untag(t)
	}
	t.expiration.Store(time.Time{})
}

// C returns the channel of the tag.
func (t *Tag) C() <-chan struct{} {
	return t.c
}

// Cancel cancels the tag.
func (t *Tag) Cancel() {
	t.c <- struct{}{}
}
