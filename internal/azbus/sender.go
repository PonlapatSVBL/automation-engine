package azbus

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Sender struct {
	sender *azservicebus.Sender
	client *azservicebus.Client
	queue  string
}

func NewSender(ctx context.Context, client *azservicebus.Client, queueName string) (*Sender, error) {
	s, err := client.NewSender(queueName, nil)
	if err != nil {
		return nil, err
	}

	return &Sender{
		sender: s,
		client: client,
		queue:  queueName,
	}, nil
}

// SendMessage sends a normal message (no scheduling)
func (s *Sender) SendMessage(ctx context.Context, body []byte) error {
	msg := &azservicebus.Message{
		Body: body,
	}

	return s.sender.SendMessage(ctx, msg, nil)
}

// ScheduleMessage sends a message scheduled for a future time
func (s *Sender) ScheduleMessage(ctx context.Context, sessionID string, body []byte, runAt time.Time) error {
	msg := &azservicebus.Message{
		SessionID:            &sessionID,
		Body:                 body,
		ScheduledEnqueueTime: &runAt,
	}

	return s.sender.SendMessage(ctx, msg, nil)
}

// Close cleans up the sender
func (s *Sender) Close(ctx context.Context) error {
	return s.sender.Close(ctx)
}
