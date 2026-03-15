package noop

import (
	"context"

	"github.com/fundacaobeta/base-canalgov-monorepo/internal/conversation/models"
)

type Noop struct {
	id      int
	channel string
	from    string
}

type Opts struct {
	ID      int
	Channel string
	From    string
}

func New(opts Opts) *Noop {
	return &Noop{
		id:      opts.ID,
		channel: opts.Channel,
		from:    opts.From,
	}
}

func (n *Noop) Close() error {
	return nil
}

func (n *Noop) Identifier() int {
	return n.id
}

func (n *Noop) Receive(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (n *Noop) Send(models.Message) error {
	return nil
}

func (n *Noop) FromAddress() string {
	return n.from
}

func (n *Noop) Channel() string {
	return n.channel
}
