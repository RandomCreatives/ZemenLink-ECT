package kernel

import (
	"context"
)

type Message struct {
	ID      string
	Content string
	Metadata map[string]interface{}
}

type Interceptor func(ctx context.Context, msg *Message) error

type Pipeline struct {
	interceptors []Interceptor
}

func NewPipeline(interceptors ...Interceptor) *Pipeline {
	return &Pipeline{interceptors: interceptors}
}

func (p *Pipeline) Execute(ctx context.Context, msg *Message) error {
	for _, interceptor := range p.interceptors {
		if err := interceptor(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
