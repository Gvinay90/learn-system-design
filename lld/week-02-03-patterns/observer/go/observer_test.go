package observer

import "testing"

func TestAllObserversNotifiedOnChange(t *testing.T) {
	station := NewWeatherStation()
	phone := &PhoneDisplay{}
	logger := &LoggingDisplay{}
	station.Subscribe(phone)
	station.Subscribe(logger)

	station.SetMeasurements(21.5, 60)

	if phone.LastTempC != 21.5 || phone.LastHumidity != 60 {
		t.Fatalf("expected phone display updated, got temp=%v humidity=%v", phone.LastTempC, phone.LastHumidity)
	}
	if len(logger.History) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(logger.History))
	}
}

func TestMultipleUpdatesAccumulate(t *testing.T) {
	station := NewWeatherStation()
	phone := &PhoneDisplay{}
	station.Subscribe(phone)

	station.SetMeasurements(20, 50)
	station.SetMeasurements(22, 55)

	if phone.UpdateCount != 2 {
		t.Fatalf("expected 2 updates, got %d", phone.UpdateCount)
	}
	if phone.LastTempC != 22 {
		t.Fatalf("expected latest temp 22, got %v", phone.LastTempC)
	}
}

func TestUnsubscribeStopsNotifications(t *testing.T) {
	station := NewWeatherStation()
	phone := &PhoneDisplay{}
	station.Subscribe(phone)
	station.SetMeasurements(20, 50)

	station.Unsubscribe(phone)
	station.SetMeasurements(99, 99)

	if phone.LastTempC != 20 {
		t.Fatalf("expected phone unaffected after unsubscribe, got %v", phone.LastTempC)
	}
	if phone.UpdateCount != 1 {
		t.Fatalf("expected only 1 update recorded, got %d", phone.UpdateCount)
	}
}

func TestNoObserversDoesNotPanic(t *testing.T) {
	station := NewWeatherStation()
	station.SetMeasurements(10, 10)
}
