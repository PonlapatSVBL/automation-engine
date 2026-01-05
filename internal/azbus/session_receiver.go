package azbus

import (
	"automation-engine/internal/httpclient"
	"automation-engine/internal/service"
	"automation-engine/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Message struct {
	AutomationID string `json:"automation_id"`
	Url          string `json:"url"`
	Version      int32  `json:"version"`
}

type SessionReceiverOptions struct {
	SessionPool int
	BatchSize   int
	ProcessPool int
	RetryDelay  int
}

type Option func(*SessionReceiverOptions)

type SessionReceiver struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	client     *azservicebus.Client
	queueName  string
	runService service.RunService
	options    SessionReceiverOptions
}

func NewSessionReceiver(ctx context.Context, wg *sync.WaitGroup, client *azservicebus.Client, queueName string, runService service.RunService, opts *SessionReceiverOptions) *SessionReceiver {
	// 1. กำหนดค่า Default
	defaultOpts := SessionReceiverOptions{
		SessionPool: 20,
		BatchSize:   5,
		ProcessPool: 1,
		RetryDelay:  5,
	}

	// 2. Apply options ที่กำหนดมา หากมี
	if opts != nil {
		if opts.SessionPool > 0 {
			defaultOpts.SessionPool = opts.SessionPool
		}
		if opts.BatchSize > 0 {
			defaultOpts.BatchSize = opts.BatchSize
		}
		if opts.ProcessPool > 0 {
			defaultOpts.ProcessPool = opts.ProcessPool
		}
		if opts.RetryDelay > 0 {
			defaultOpts.RetryDelay = opts.RetryDelay
		}
	}

	// 3. สร้าง Struct โดยใช้ opts ที่ได้มา
	return &SessionReceiver{
		ctx:        ctx,
		wg:         wg,
		client:     client,
		queueName:  queueName,
		runService: runService,
		options:    defaultOpts,
	}
}

// RunDispatcher continuously accepts sessions from the queue and dispatches workers to process them
func (sr *SessionReceiver) RunDispatcher() {
	// Spawn worker (Initialization)
	workerCH := make(chan int, sr.options.SessionPool)
	for i := range make([]int, sr.options.SessionPool) {
		workerCH <- i + 1
	}

	for {
		select {
		case <-sr.ctx.Done():
			log.Printf("[%s] Shutting down Session Receiver. Waiting for active sessions to complete...", sr.queueName)
			sr.wg.Wait()
			log.Printf("[%s] All active sessions completed. Session Receiver stopped.", sr.queueName)
			return
		default:
		}

		select {
		case workerNo := <-workerCH:
			acceptSessionCtx, acceptSessionCancel := context.WithTimeout(sr.ctx, 5*time.Second)

			// sessionReceiver, err := sr.client.AcceptSessionForQueue(acceptSessionCtx, sr.queueName, "202205072D33E7BA048F", nil)
			sessionReceiver, err := sr.client.AcceptNextSessionForQueue(acceptSessionCtx, sr.queueName, nil)
			if err != nil {
				acceptSessionCancel()

				workerCH <- workerNo
				if errors.Is(err, context.DeadlineExceeded) {
					log.Printf("[%s] Worker<%d>: Session accept timed out. Retrying in %d seconds...", sr.queueName, workerNo, sr.options.RetryDelay)
				}
				time.Sleep(time.Duration(sr.options.RetryDelay) * time.Second)
				continue
			}

			sr.wg.Add(1)
			go sr.runSessionWorker(sessionReceiver, acceptSessionCancel, workerCH, workerNo)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// runSessionWorker processes messages from a single session
func (sr *SessionReceiver) runSessionWorker(sessionReceiver *azservicebus.SessionReceiver, acceptSessionCancel context.CancelFunc, workerCH chan int, workerNo int) {
	defer sr.wg.Done()
	defer func() { workerCH <- workerNo }()
	defer acceptSessionCancel()
	defer sessionReceiver.Close(sr.ctx)
	// defer log.Printf("[%s] Worker<%d>: %s [/]\n", sr.queueName, workerNo, sessionReceiver.SessionID())

	log.Printf("[%s] Worker<%d>: %s [ ]\n", sr.queueName, workerNo, sessionReceiver.SessionID())

	recvCtx, cancel := context.WithTimeout(sr.ctx, 10*time.Second)

	msgs, err := sessionReceiver.ReceiveMessages(recvCtx, sr.options.BatchSize, nil)
	cancel()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("[%s] Worker<%d>: No messages, polling again... [/]\n", sr.queueName, workerNo)
			return
		}

		log.Printf("[%s] Receive error: %v", sr.queueName, err)
		return
	}

	defer log.Printf("[%s] Worker<%d>: %s | Messages: %d [/]\n", sr.queueName, workerNo, sessionReceiver.SessionID(), len(msgs))

	if len(msgs) == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, sr.options.ProcessPool)

	for _, msg := range msgs {
		sem <- struct{}{}
		wg.Add(1)
		go sr.runMessageWorker(sem, &wg, sessionReceiver, msg)
	}

	wg.Wait()
}

// runMessageWorker handles a single message, applies business logic and updates Redis
func (sr *SessionReceiver) runMessageWorker(sem chan struct{}, wg *sync.WaitGroup, sessionReceiver *azservicebus.SessionReceiver, msg *azservicebus.ReceivedMessage) {
	defer wg.Done()
	defer func() { <-sem }()

	err := sr.handleMessage(msg)
	if err != nil {
		// Log: Failed/Abandon
		fmt.Println(err)
		sessionReceiver.AbandonMessage(sr.ctx, msg, nil)
		return
	}

	// Log: Success/Complete
	sessionReceiver.CompleteMessage(sr.ctx, msg, nil)
}

// handleMessage contains the actual business logic for processing a single message
func (sr *SessionReceiver) handleMessage(msg *azservicebus.ReceivedMessage) error {
	// Check for nil message early
	if msg == nil {
		return fmt.Errorf("received nil message")
	}

	var body Message
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("invalid message json: %w", err)
	}

	// Check for required fields in the message body
	if body.Url == "" {
		return fmt.Errorf("url is empty in message")
	}

	// Fetch automation details using RunService
	automation, err := sr.runService.GetAutomationByID(sr.ctx, body.AutomationID)
	if err != nil {
		return fmt.Errorf("failed to get automation by ID: %w", err)
	}

	if automation.Version != body.Version {
		// return fmt.Errorf("automation version mismatch: message version %d, automation version %d", body.Version, automation.Version)
		return nil
	}
	if automation.IsActive != "Y" {
		// return fmt.Errorf("automation is not active")
		return nil
	}
	if automation.NextRun.After(time.Now()) {
		// return fmt.Errorf("automation is not due to run yet")
		return nil
	}

	switch automation.IntervalType {
	case "daily":
		next, err := utils.CalculateDailyNextRun(
			time.Now(),
			automation.Time,
			time.Local,
		)
		if err != nil {
			return fmt.Errorf("failed to calculate next daily run: %w", err)
		}
		log.Printf("Next daily run calculated: %s", next)
	case "weekly":
		next, err := utils.CalculateWeeklyNextRun(
			time.Now(),
			automation.Time,
			automation.DayOfWeek,
			time.Local,
		)
		if err != nil {
			return fmt.Errorf("failed to calculate next weekly run: %w", err)
		}
		log.Printf("Next weekly run calculated: %s", next)
	case "monthly":
		next, err := utils.CalculateMonthlyNextRun(
			time.Now(),
			automation.Time,
			int(automation.DayOfMonth),
			time.Local,
		)
		if err != nil {
			return fmt.Errorf("failed to calculate next monthly run: %w", err)
		}
		log.Printf("Next monthly run calculated: %s", next)
	default:
		return fmt.Errorf("unknown interval type: %s", automation.IntervalType)
	}

	statusCode, _, err := httpclient.PostRequest(body.Url, msg.Body)

	// Handle non-200 status code or network/client errors
	if err != nil || statusCode != 200 {
		return fmt.Errorf("post request failed")
	}

	// Successfully processed the message
	return nil
}
