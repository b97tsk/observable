package rx

import (
	"context"
	"sync"
	"time"
)

// A ConnectableObservable is an Observable that only subscribes to the source
// Observable by calling its Connect method. Calling its Subscribe method will
// not subscribe the source, instead, it subscribes to a local Subject, which
// means that it can be called many times with different Observers.
type ConnectableObservable struct {
	*connectableObservable
}

type connectableObservable struct {
	Observable
	mux            sync.Mutex
	source         Observable
	subjectFactory func() Subject
	connection     context.Context
	disconnect     context.CancelFunc
	subject        Subject
	refCount       int
}

func newConnectableObservable(source Observable, subjectFactory func() Subject) *connectableObservable {
	connectable := connectableObservable{
		source:         source,
		subjectFactory: subjectFactory,
	}
	connectable.Observable = Observable(
		func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return connectable.getSubject().Subscribe(ctx, sink)
		},
	)
	return &connectable
}

func (obs *connectableObservable) getSubjectLocked() Subject {
	if obs.subject.Observer == nil {
		obs.subject = obs.subjectFactory()
	}
	return obs.subject
}

func (obs *connectableObservable) getSubject() Subject {
	obs.mux.Lock()
	defer obs.mux.Unlock()
	return obs.getSubjectLocked()
}

func (obs *connectableObservable) connect(addRef bool) (context.Context, context.CancelFunc) {
	obs.mux.Lock()
	defer obs.mux.Unlock()

	connection := obs.connection

	if connection == nil {
		type X struct{}
		cx := make(chan X, 1)
		cx <- X{}

		subject := obs.getSubjectLocked()

		ctx, cancel := obs.source.Subscribe(context.Background(), func(t Notification) {
			if t.HasValue {
				t.Observe(subject.Observer)
				return
			}

			x, cxLocked := <-cx
			if !cxLocked {
				obs.mux.Lock()
			}

			if connection == obs.connection {
				obs.connection = nil
				obs.disconnect = nil
				obs.subject = Subject{}
				obs.refCount = 0
			}

			if cxLocked {
				cx <- x
			} else {
				obs.mux.Unlock()
			}

			t.Observe(subject.Observer)
		})

		<-cx
		defer close(cx)

		if ctx.Err() != nil {
			return Done()
		}

		connection = ctx
		obs.connection = ctx
		obs.disconnect = cancel
	}

	if addRef {
		obs.refCount++

		return connection, func() {
			obs.mux.Lock()
			defer obs.mux.Unlock()

			if connection != obs.connection {
				return
			}
			if obs.refCount == 0 {
				return
			}

			obs.refCount--

			if obs.refCount == 0 {
				obs.disconnect()
				obs.connection = nil
				obs.disconnect = nil
				obs.subject = Subject{}
			}
		}
	}

	return connection, func() {
		obs.mux.Lock()
		defer obs.mux.Unlock()

		if connection != obs.connection {
			return
		}

		obs.disconnect()
		obs.connection = nil
		obs.disconnect = nil
		obs.subject = Subject{}
		obs.refCount = 0
	}
}

func (obs *connectableObservable) connectAddRef() (context.Context, context.CancelFunc) {
	return obs.connect(true)
}

// Connect invokes an execution of an ConnectableObservable.
func (obs ConnectableObservable) Connect() (context.Context, context.CancelFunc) {
	return obs.connect(false)
}

type refCountOperator struct {
	Connectable ConnectableObservable
}

func (op refCountOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := op.Connectable.Subscribe(ctx, sink)
	_, releaseRef := op.Connectable.connectAddRef()

	go func() {
		<-ctx.Done()
		releaseRef()
	}()

	return ctx, cancel
}

// RefCount creates an Observable that keeps track of how many subscribers
// it has. When the number of subscribers increases from 0 to 1, it will call
// Connect() for us, which starts the shared execution. Only when the number
// of subscribers decreases from 1 to 0 will it be fully unsubscribed, stopping
// further execution.
func (obs ConnectableObservable) RefCount() Observable {
	op := refCountOperator{obs}
	return Empty().Lift(op.Call)
}

// Multicast returns a ConnectableObservable, which is a variety of Observable
// that waits until its Connect method is called before it begins emitting
// items to those Observers that have subscribed to it.
func (obs Observable) Multicast(subjectFactory func() Subject) ConnectableObservable {
	return ConnectableObservable{newConnectableObservable(obs, subjectFactory)}
}

// Publish is like Multicast, but it uses only one subject.
func (obs Observable) Publish() ConnectableObservable {
	subject := NewSubject()
	return obs.Multicast(func() Subject { return subject })
}

// PublishBehavior is like Publish, but it uses a BehaviorSubject instead.
func (obs Observable) PublishBehavior(val interface{}) ConnectableObservable {
	subject := NewBehaviorSubject(val)
	return obs.Multicast(func() Subject { return subject.Subject })
}

// PublishReplay is like Publish, but it uses a ReplaySubject instead.
func (obs Observable) PublishReplay(bufferSize int, windowTime time.Duration) ConnectableObservable {
	subject := NewReplaySubject(bufferSize, windowTime)
	return obs.Multicast(func() Subject { return subject.Subject })
}

// Share returns a new Observable that multicasts (shares) the original
// Observable. When subscribed multiple times, it guarantees that only one
// subscription is made to the source Observable at the same time. When all
// subscribers have unsubscribed it will unsubscribe from the source Observable.
func (Operators) Share() OperatorFunc {
	return func(source Observable) Observable {
		return source.Multicast(NewSubject).RefCount()
	}
}

// ShareReplay is like Share, but it uses a ReplaySubject instead.
func (Operators) ShareReplay(bufferSize int, windowTime time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		newSubject := func() Subject {
			return NewReplaySubject(bufferSize, windowTime).Subject
		}
		return source.Multicast(newSubject).RefCount()
	}
}
