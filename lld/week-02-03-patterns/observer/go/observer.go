// Package observer demonstrates the Observer pattern: a WeatherStation
// subject notifies multiple registered Observers whenever its readings change.
package observer

import "fmt"

// Observer receives updates whenever the subject's state changes.
type Observer interface {
	Update(tempC float64, humidity float64)
}

// WeatherStation is the subject: it holds the current readings and notifies
// all subscribed observers whenever they change.
type WeatherStation struct {
	tempC     float64
	humidity  float64
	observers []Observer
}

func NewWeatherStation() *WeatherStation {
	return &WeatherStation{}
}

func (w *WeatherStation) Subscribe(o Observer) {
	w.observers = append(w.observers, o)
}

// Unsubscribe removes the first matching observer, if present.
func (w *WeatherStation) Unsubscribe(o Observer) {
	for i, existing := range w.observers {
		if existing == o {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
			return
		}
	}
}

// SetMeasurements updates the readings and notifies every subscribed observer.
func (w *WeatherStation) SetMeasurements(tempC, humidity float64) {
	w.tempC = tempC
	w.humidity = humidity
	for _, o := range w.observers {
		o.Update(tempC, humidity)
	}
}

// PhoneDisplay is a concrete observer that remembers only the latest reading.
type PhoneDisplay struct {
	LastTempC    float64
	LastHumidity float64
	UpdateCount  int
}

func (d *PhoneDisplay) Update(tempC, humidity float64) {
	d.LastTempC = tempC
	d.LastHumidity = humidity
	d.UpdateCount++
}

// LoggingDisplay is a concrete observer that keeps a full history.
type LoggingDisplay struct {
	History []string
}

func (d *LoggingDisplay) Update(tempC, humidity float64) {
	d.History = append(d.History, fmt.Sprintf("temp=%.1f humidity=%.1f", tempC, humidity))
}
