package notificationsystem

import (
	"testing"
	"time"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]string
		want     string
	}{
		{
			name:     "single placeholder",
			template: "Hello {name}",
			data:     map[string]string{"name": "Alice"},
			want:     "Hello Alice",
		},
		{
			name:     "multiple placeholders",
			template: "Hello {name}, your order {orderId} shipped",
			data:     map[string]string{"name": "Bob", "orderId": "42"},
			want:     "Hello Bob, your order 42 shipped",
		},
		{
			name:     "missing key left untouched",
			template: "Hi {name}, code {code}",
			data:     map[string]string{"name": "Carl"},
			want:     "Hi Carl, code {code}",
		},
		{
			name:     "no placeholders",
			template: "plain message",
			data:     map[string]string{"unused": "x"},
			want:     "plain message",
		},
		{
			name:     "unterminated brace",
			template: "Hello {name",
			data:     map[string]string{"name": "Dee"},
			want:     "Hello {name",
		},
		{
			name:     "empty template",
			template: "",
			data:     map[string]string{},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderTemplate(tt.template, tt.data)
			if got != tt.want {
				t.Fatalf("RenderTemplate(%q) = %q, want %q", tt.template, got, tt.want)
			}
		})
	}
}

func TestNotifyDispatchesPerPreference(t *testing.T) {
	email := NewEmailChannel()
	sms := NewSMSChannel()
	push := NewPushChannel()

	svc := NewNotificationService(RetryPolicy{MaxAttempts: 1})
	svc.RegisterChannel(email)
	svc.RegisterChannel(sms)
	svc.RegisterChannel(push)

	svc.SetPreferences("u1", Email, SMS)
	svc.SetPreferences("u2", Push)

	results, err := svc.Notify("u1", "u1@example.com", "Hello {name}", map[string]string{"name": "Ann"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results for u1, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("unexpected send error on channel %s: %v", r.Channel, r.Err)
		}
	}
	if len(email.Sent()) != 1 || email.Sent()[0].Message != "Hello Ann" {
		t.Fatalf("expected email to receive rendered message, got %+v", email.Sent())
	}
	if len(sms.Sent()) != 1 {
		t.Fatalf("expected sms to receive one message, got %+v", sms.Sent())
	}
	if len(push.Sent()) != 0 {
		t.Fatalf("expected push to receive nothing for u1, got %+v", push.Sent())
	}

	if _, err := svc.Notify("u2", "device-token", "Ping {name}", map[string]string{"name": "Bo"}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(push.Sent()) != 1 || push.Sent()[0].Message != "Ping Bo" {
		t.Fatalf("expected push to receive rendered message for u2, got %+v", push.Sent())
	}
}

func TestNotifyUnknownUserReturnsError(t *testing.T) {
	svc := NewNotificationService(DefaultRetryPolicy())
	svc.RegisterChannel(NewEmailChannel())

	if _, err := svc.Notify("ghost", "x", "hi", nil); err == nil {
		t.Fatalf("expected error for user with no preferences")
	}
}

func TestRetrySucceedsAfterNFailures(t *testing.T) {
	underlying := NewEmailChannel()
	flaky := NewFlakyChannel(underlying, 2) // fails twice, succeeds on 3rd attempt

	svc := NewNotificationService(RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond})
	svc.RegisterChannel(flaky)
	svc.SetPreferences("u1", Email)

	results, err := svc.Notify("u1", "u1@example.com", "Order {orderId} shipped", map[string]string{"orderId": "7"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Err != nil {
		t.Fatalf("expected eventual success, got err: %v", r.Err)
	}
	if r.Attempts != 3 {
		t.Fatalf("expected 3 attempts (2 failures + 1 success), got %d", r.Attempts)
	}
	if flaky.Attempts() != 3 {
		t.Fatalf("expected channel to have been called 3 times, got %d", flaky.Attempts())
	}
	if len(underlying.Sent()) != 1 || underlying.Sent()[0].Message != "Order 7 shipped" {
		t.Fatalf("expected underlying channel to record exactly one successful send, got %+v", underlying.Sent())
	}
}

func TestRetryGivesUpAfterMaxAttempts(t *testing.T) {
	failing := NewAlwaysFailChannel(SMS)

	svc := NewNotificationService(RetryPolicy{MaxAttempts: 4, Delay: time.Millisecond})
	svc.RegisterChannel(failing)
	svc.SetPreferences("u1", SMS)

	results, err := svc.Notify("u1", "5551234", "hi {name}", map[string]string{"name": "X"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Err == nil {
		t.Fatalf("expected send to ultimately fail")
	}
	if r.Attempts != 4 {
		t.Fatalf("expected exactly MaxAttempts=4 attempts, got %d", r.Attempts)
	}
	if failing.Attempts() != 4 {
		t.Fatalf("expected channel to have been called 4 times, got %d", failing.Attempts())
	}
}

func TestNotifyMultipleChannelsIndependentRetry(t *testing.T) {
	email := NewEmailChannel()
	flakySMS := NewFlakyChannel(NewSMSChannel(), 1)

	svc := NewNotificationService(RetryPolicy{MaxAttempts: 2, Delay: time.Millisecond})
	svc.RegisterChannel(email)
	svc.RegisterChannel(flakySMS)
	svc.SetPreferences("u1", Email, SMS)

	results, err := svc.Notify("u1", "recipient", "msg", nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("expected both channels to eventually succeed, channel %s failed: %v", r.Channel, r.Err)
		}
	}
	if results[0].Attempts != 1 {
		t.Fatalf("expected email to succeed on first attempt, got %d", results[0].Attempts)
	}
	if results[1].Attempts != 2 {
		t.Fatalf("expected sms to succeed on second attempt, got %d", results[1].Attempts)
	}
}
