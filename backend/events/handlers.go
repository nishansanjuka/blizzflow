package events

import "context"

type EventHandler interface {
	OnDomReady(ctx context.Context)
	OnBeforeClose(ctx context.Context) bool
	OnShutdown(ctx context.Context)
}

type HandlerRegistry struct {
	handlers []EventHandler
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make([]EventHandler, 0),
	}
}

func (r *HandlerRegistry) Register(handler EventHandler) {
	r.handlers = append(r.handlers, handler)
}

func (r *HandlerRegistry) TriggerDomReady(ctx context.Context) {
	for _, h := range r.handlers {
		h.OnDomReady(ctx)
	}
}

func (r *HandlerRegistry) TriggerBeforeClose(ctx context.Context) bool {
	for _, h := range r.handlers {
		if h.OnBeforeClose(ctx) {
			return true
		}
	}
	return false
}

func (r *HandlerRegistry) TriggerShutdown(ctx context.Context) {
	for _, h := range r.handlers {
		h.OnShutdown(ctx)
	}
}
