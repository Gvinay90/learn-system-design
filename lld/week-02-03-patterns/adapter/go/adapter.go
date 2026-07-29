// Package adapter demonstrates the Adapter structural pattern: a legacy
// XML-based data provider is wrapped so it satisfies a modern
// JSON-style DataProvider interface, without modifying the legacy code.
package adapter

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

// DataProvider is the modern interface the rest of the system depends on.
type DataProvider interface {
	FetchJSON(id string) (map[string]string, error)
}

// LegacyXMLDataProvider is the pre-existing incompatible interface;
// it can't be changed (owned by another team / third-party).
type LegacyXMLDataProvider struct {
	store map[string]string
}

func NewLegacyXMLDataProvider(store map[string]string) *LegacyXMLDataProvider {
	return &LegacyXMLDataProvider{store: store}
}

type xmlRecord struct {
	XMLName xml.Name `xml:"record"`
	ID      string   `xml:"id,attr"`
	Value   string   `xml:",chardata"`
}

// FetchXML returns the raw legacy XML payload for id.
func (l *LegacyXMLDataProvider) FetchXML(id string) (string, error) {
	v, ok := l.store[id]
	if !ok {
		return "", errors.New("record not found: " + id)
	}
	rec := xmlRecord{ID: id, Value: v}
	out, err := xml.Marshal(rec)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// XMLToJSONAdapter adapts a LegacyXMLDataProvider to the DataProvider interface,
// translating XML responses into the map-based shape callers expect.
type XMLToJSONAdapter struct {
	Legacy *LegacyXMLDataProvider
}

func NewXMLToJSONAdapter(legacy *LegacyXMLDataProvider) *XMLToJSONAdapter {
	return &XMLToJSONAdapter{Legacy: legacy}
}

func (a *XMLToJSONAdapter) FetchJSON(id string) (map[string]string, error) {
	raw, err := a.Legacy.FetchXML(id)
	if err != nil {
		return nil, err
	}
	var rec xmlRecord
	if err := xml.Unmarshal([]byte(raw), &rec); err != nil {
		return nil, err
	}
	return map[string]string{"id": rec.ID, "value": strings.TrimSpace(rec.Value)}, nil
}

// ModernDataProvider is a native implementation of the target interface,
// included to show the adapter is a drop-in replacement wherever DataProvider is used.
type ModernDataProvider struct {
	store map[string]string
}

func NewModernDataProvider(store map[string]string) *ModernDataProvider {
	return &ModernDataProvider{store: store}
}

func (m *ModernDataProvider) FetchJSON(id string) (map[string]string, error) {
	v, ok := m.store[id]
	if !ok {
		return nil, errors.New("record not found: " + id)
	}
	return map[string]string{"id": id, "value": v}, nil
}

// FetchAndSum is a client function that only knows about DataProvider,
// oblivious to whether the underlying source is legacy XML or native JSON.
func FetchAndSum(p DataProvider, ids []string) (int, error) {
	total := 0
	for _, id := range ids {
		rec, err := p.FetchJSON(id)
		if err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(rec["value"])
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}
