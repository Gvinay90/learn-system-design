package factory

import "testing"

func TestCreateEmailNotification(t *testing.T) {
	n, err := CreateNotification(Email)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	got := n.Send("alice@example.com", "hello")
	want := "Email to alice@example.com: hello"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCreateSMSNotification(t *testing.T) {
	n, err := CreateNotification(SMS)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	got := n.Send("+1-555-0100", "hello")
	want := "SMS to +1-555-0100: hello"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCreatePushNotification(t *testing.T) {
	n, err := CreateNotification(Push)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	got := n.Send("device-123", "hello")
	want := "Push to device-123: hello"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCreateUnknownNotificationType(t *testing.T) {
	_, err := CreateNotification(NotificationType(99))
	if err != ErrUnknownNotificationType {
		t.Fatalf("expected ErrUnknownNotificationType, got %v", err)
	}
}

func TestDemo(t *testing.T) {
	for _, nt := range []NotificationType{Email, SMS, Push} {
		n, err := CreateNotification(nt)
		if err != nil {
			t.Fatalf("unexpected err for type %v: %v", nt, err)
		}
		if n.Send("user", "test") == "" {
			t.Fatalf("expected non-empty send result")
		}
	}
}
