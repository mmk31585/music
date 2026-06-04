package eventbus

import "context"

type Handler func(ctx context.Context, event Event) error

type Subscriber interface {
	Subscribe(eventName string, handler Handler)
}
