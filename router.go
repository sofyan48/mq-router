package mq

import (
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
)

// MessageAttributeNameRoute is a MessageAttribute name used as a routing key by
// the Router.
const MessageAttributeNameRoute = "path"
const MessageAttributeNameMethod = "method"

// Router will route a message based on MessageAttributes to other registered Handlers.
type Router struct {
	sync.Mutex

	// Resolver maps a Message to a string identifier used to match to a registered Handler. The
	// default implementation returns a MessageAttribute named "route".
	Resolver func(*Message) (string, bool)

	// A map of handlers to route to. The return value of Resolver should match a key in this map.
	handlers map[string]Handler
	routes   string

	handle Handler
}

// NewRouter returns a new Router.
func NewRouter() *Router {
	return &Router{
		Resolver: defaultResolver,
		handlers: map[string]Handler{},
	}
}

// Handle registers a Handler under a route key.
func (r *Router) Handle(route string, h Handler) *Router {
	r.routes = route
	r.handle = h
	return r
}

func (r *Router) Method(method string) {
	r.Lock()
	defer r.Unlock()

	r.handlers[r.routes] = r.handle
	r.handlers[method] = r.handle
}

// HandleMessage satisfies the Handler interface.
func (r *Router) HandleMessage(m *Message) error {
	key, ok := r.Resolver(m)
	if !ok {
		return errors.New("no routing key for message")
	}

	if h, ok := r.handlers[key]; ok {
		return h.HandleMessage(m)
	}

	return fmt.Errorf("no handler matched for routing key: %s", key)
}

func defaultResolver(m *Message) (string, bool) {
	r := ""
	v, ok := m.SQSMessage.MessageAttributes[MessageAttributeNameRoute]
	if ok {
		r = aws.StringValue(v.StringValue)
	}
	d, ok := m.SQSMessage.MessageAttributes[MessageAttributeNameMethod]
	if ok {
		r = aws.StringValue(d.StringValue)
	}
	return r, ok
}
