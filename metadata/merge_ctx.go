package metadata

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// mergeCtx is a context that merges two contexts.
type mergeCtx struct {
	doneOne    sync.Once
	cancelOnce sync.Once

	ctx1, ctx2 context.Context

	doneErr  error
	doneMark atomic.Bool

	done     chan struct{}
	cancelCh chan struct{}
}

// Merge merges two contexts into one. The returned context is done when either of the input contexts is done or when the returned cancel function is called.
func Merge(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	mc := &mergeCtx{
		ctx1:     ctx1,
		ctx2:     ctx2,
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	select {
	case <-ctx1.Done():
		_ = mc.finish(ctx1.Err())
	case <-ctx2.Done():
		_ = mc.finish(ctx2.Err())
	default:
		go mc.wait()
	}
	return mc, mc.cancel
}

// finish sets the doneErr and closes the done channel.
func (mc *mergeCtx) finish(err error) error {
	mc.doneOne.Do(func() {
		mc.doneErr = err
		mc.doneMark.Store(true)
		close(mc.done)
	})
	return mc.doneErr
}

// wait waits for either of the two contexts to be done or for cancellation.
func (mc *mergeCtx) wait() {
	var err error
	select {
	case <-mc.ctx1.Done():
		err = mc.ctx1.Err()
	case <-mc.ctx2.Done():
		err = mc.ctx2.Err()
	case <-mc.cancelCh:
		err = context.Canceled
	}
	_ = mc.finish(err)
}

// cancel cancels the mergeCtx.
func (mc *mergeCtx) cancel() {
	mc.cancelOnce.Do(func() {
		close(mc.cancelCh)
	})
}

// Deadline returns the time when work done on behalf of this context
// should be canceled.  Deadline returns ok==false when no deadline is
// set.  Successive calls to Deadline return the same results.
func (mc *mergeCtx) Deadline() (deadline time.Time, ok bool) {
	d1, ok1 := mc.ctx1.Deadline()
	d2, ok2 := mc.ctx2.Deadline()
	if !ok1 && !ok2 {
		return time.Time{}, false
	}
	if !ok1 {
		return d2, true
	}
	if !ok2 {
		return d1, true
	}
	if d1.Before(d2) {
		return d1, true
	}
	return d2, true
}

// Done returns a channel that's closed when work done on behalf of this context should be canceled.
func (mc *mergeCtx) Done() <-chan struct{} {
	return mc.done
}

// Err returns a non-nil error value after Done is closed.
func (mc *mergeCtx) Err() error {
	if mc.doneMark.Load() {
		return mc.doneErr
	}
	var err error
	select {
	case <-mc.ctx1.Done():
		err = mc.ctx1.Err()
	case <-mc.ctx2.Done():
		err = mc.ctx2.Err()
	case <-mc.cancelCh:
		err = context.Canceled
	default:
		return nil
	}
	return mc.finish(err)
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.  Successive calls to Value with
// the same key returns the same result.
func (mc *mergeCtx) Value(key any) any {
	if v := mc.ctx1.Value(key); v != nil {
		return v
	}
	return mc.ctx2.Value(key)
}
