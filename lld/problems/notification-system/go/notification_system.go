// Package notificationsystem implements a multi-channel notification LLD
// problem: pluggable delivery channels via the Strategy pattern, per-user
// channel preferences, retry-with-fixed-delay on send failure, and basic
// "{placeholder}" template rendering.
package notificationsystem

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ErrSendFailed is returned by a channel implementation to simulate a
// delivery failure (e.g. a flaky network call).
var ErrSendFailed = errors.New("notification: send failed")

// Channel identifies a delivery mechanism.
type Channel string

const (
	Email Channel = "EMAIL"
	SMS   Channel = "SMS"
	Push  Channel = "PUSH"
)

// Notifier is the Strategy interface every delivery channel implements.
type Notifier interface {
	Channel() Channel
	Send(recipient, message string) error
}

// sentMessage records one successful delivery, used by the in-memory
// channels below so tests can assert on what was actually sent.
type sentMessage struct {
	Recipient string
	Message   string
}

// EmailChannel simulates sending email by recording messages in memory.
type EmailChannel struct {
	mu   sync.Mutex
	sent []sentMessage
}

func NewEmailChannel() *EmailChannel { return &EmailChannel{} }

func (c *EmailChannel) Channel() Channel { return Email }

func (c *EmailChannel) Send(recipient, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, sentMessage{Recipient: recipient, Message: message})
	return nil
}

func (c *EmailChannel) Sent() []sentMessage {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]sentMessage, len(c.sent))
	copy(out, c.sent)
	return out
}

// SMSChannel simulates sending SMS by recording messages in memory.
type SMSChannel struct {
	mu   sync.Mutex
	sent []sentMessage
}

func NewSMSChannel() *SMSChannel { return &SMSChannel{} }

func (c *SMSChannel) Channel() Channel { return SMS }

func (c *SMSChannel) Send(recipient, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, sentMessage{Recipient: recipient, Message: message})
	return nil
}

func (c *SMSChannel) Sent() []sentMessage {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]sentMessage, len(c.sent))
	copy(out, c.sent)
	return out
}

// PushChannel simulates sending a push notification by recording messages
// in memory.
type PushChannel struct {
	mu   sync.Mutex
	sent []sentMessage
}

func NewPushChannel() *PushChannel { return &PushChannel{} }

func (c *PushChannel) Channel() Channel { return Push }

func (c *PushChannel) Send(recipient, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, sentMessage{Recipient: recipient, Message: message})
	return nil
}

func (c *PushChannel) Sent() []sentMessage {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]sentMessage, len(c.sent))
	copy(out, c.sent)
	return out
}

// FlakyChannel wraps a delegate Notifier and fails the first N sends before
// delegating for real. It exists purely so tests can exercise the retry
// mechanism deterministically.
type FlakyChannel struct {
	mu         sync.Mutex
	delegate   Notifier
	failures   int // remaining forced failures
	attempts   int
	AttemptLog []string
}

// NewFlakyChannel returns a channel that fails the first failCount sends
// (returning ErrSendFailed) and then delegates to delegate.
func NewFlakyChannel(delegate Notifier, failCount int) *FlakyChannel {
	return &FlakyChannel{delegate: delegate, failures: failCount}
}

func (c *FlakyChannel) Channel() Channel { return c.delegate.Channel() }

func (c *FlakyChannel) Send(recipient, message string) error {
	c.mu.Lock()
	c.attempts++
	if c.failures > 0 {
		c.failures--
		c.mu.Unlock()
		return ErrSendFailed
	}
	c.mu.Unlock()
	return c.delegate.Send(recipient, message)
}

func (c *FlakyChannel) Attempts() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.attempts
}

// AlwaysFailChannel always fails, useful for exercising "retries exhausted".
type AlwaysFailChannel struct {
	channel  Channel
	mu       sync.Mutex
	attempts int
}

func NewAlwaysFailChannel(channel Channel) *AlwaysFailChannel {
	return &AlwaysFailChannel{channel: channel}
}

func (c *AlwaysFailChannel) Channel() Channel { return c.channel }

func (c *AlwaysFailChannel) Send(string, string) error {
	c.mu.Lock()
	c.attempts++
	c.mu.Unlock()
	return ErrSendFailed
}

func (c *AlwaysFailChannel) Attempts() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.attempts
}

// RenderTemplate substitutes "{key}" tokens in template with the string
// form of data[key]. Tokens with no matching key are left untouched.
func RenderTemplate(template string, data map[string]string) string {
	var b strings.Builder
	i := 0
	for i < len(template) {
		open := strings.IndexByte(template[i:], '{')
		if open == -1 {
			b.WriteString(template[i:])
			break
		}
		open += i
		close := strings.IndexByte(template[open:], '}')
		if close == -1 {
			b.WriteString(template[i:])
			break
		}
		close += open

		b.WriteString(template[i:open])
		key := template[open+1 : close]
		if val, ok := data[key]; ok {
			b.WriteString(val)
		} else {
			b.WriteString(template[open : close+1])
		}
		i = close + 1
	}
	return b.String()
}

// RetryPolicy configures the retry-on-failure wrapper around a channel send.
type RetryPolicy struct {
	MaxAttempts int           // total attempts, including the first; must be >= 1
	Delay       time.Duration // fixed delay between attempts
}

// DefaultRetryPolicy retries up to 3 times total with a 10ms delay, suitable
// for tests; production callers should tune this.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{MaxAttempts: 3, Delay: 10 * time.Millisecond}
}

// SendResult captures the outcome of dispatching to a single channel.
type SendResult struct {
	Channel  Channel
	Attempts int
	Err      error // nil on success
}

// NotificationService dispatches rendered notifications to a user's
// preferred channels, retrying each channel independently on failure.
type NotificationService struct {
	mu          sync.RWMutex
	channels    map[Channel]Notifier
	preferences map[string][]Channel
	retry       RetryPolicy
}

func NewNotificationService(retry RetryPolicy) *NotificationService {
	return &NotificationService{
		channels:    make(map[Channel]Notifier),
		preferences: make(map[string][]Channel),
		retry:       retry,
	}
}

// RegisterChannel makes a Notifier available for dispatch under its own
// Channel() identity.
func (s *NotificationService) RegisterChannel(n Notifier) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.channels[n.Channel()] = n
}

// SetPreferences records which channels a given user wants to receive
// notifications on, in the order they should be attempted.
func (s *NotificationService) SetPreferences(userID string, channels ...Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.preferences[userID] = append([]Channel(nil), channels...)
}

// Notify renders template with data and sends the result to every channel
// preferred by userID, retrying each channel per the service's RetryPolicy.
// It returns one SendResult per attempted channel and does not stop early
// if one channel ultimately fails.
func (s *NotificationService) Notify(userID, recipient, template string, data map[string]string) ([]SendResult, error) {
	s.mu.RLock()
	prefs, ok := s.preferences[userID]
	if !ok {
		s.mu.RUnlock()
		return nil, fmt.Errorf("notification: no channel preferences registered for user %q", userID)
	}
	prefsCopy := append([]Channel(nil), prefs...)
	s.mu.RUnlock()

	message := RenderTemplate(template, data)

	results := make([]SendResult, 0, len(prefsCopy))
	for _, ch := range prefsCopy {
		s.mu.RLock()
		notifier, ok := s.channels[ch]
		s.mu.RUnlock()
		if !ok {
			results = append(results, SendResult{Channel: ch, Err: fmt.Errorf("notification: no channel registered for %q", ch)})
			continue
		}
		results = append(results, s.sendWithRetry(notifier, recipient, message))
	}
	return results, nil
}

// sendWithRetry attempts notifier.Send up to s.retry.MaxAttempts times,
// sleeping s.retry.Delay between attempts, and returns the outcome.
func (s *NotificationService) sendWithRetry(notifier Notifier, recipient, message string) SendResult {
	maxAttempts := s.retry.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lastErr = notifier.Send(recipient, message)
		if lastErr == nil {
			return SendResult{Channel: notifier.Channel(), Attempts: attempt, Err: nil}
		}
		if attempt < maxAttempts {
			time.Sleep(s.retry.Delay)
		}
	}
	return SendResult{Channel: notifier.Channel(), Attempts: maxAttempts, Err: lastErr}
}
