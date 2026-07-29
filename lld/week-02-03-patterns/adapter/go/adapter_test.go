package adapter

import "testing"

func TestAdapterTranslatesLegacyCall(t *testing.T) {
	legacy := NewLegacyXMLDataProvider(map[string]string{"u1": "42"})
	adapted := NewXMLToJSONAdapter(legacy)

	rec, err := adapted.FetchJSON("u1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rec["id"] != "u1" || rec["value"] != "42" {
		t.Fatalf("expected {id:u1 value:42}, got %v", rec)
	}
}

func TestAdapterPropagatesNotFound(t *testing.T) {
	legacy := NewLegacyXMLDataProvider(map[string]string{})
	adapted := NewXMLToJSONAdapter(legacy)

	if _, err := adapted.FetchJSON("missing"); err == nil {
		t.Fatal("expected error for missing record")
	}
}

func TestClientCodeIsProviderAgnostic(t *testing.T) {
	legacy := NewXMLToJSONAdapter(NewLegacyXMLDataProvider(map[string]string{"a": "10", "b": "20"}))
	modern := NewModernDataProvider(map[string]string{"a": "10", "b": "20"})

	legacySum, err := FetchAndSum(legacy, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	modernSum, err := FetchAndSum(modern, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if legacySum != 30 || modernSum != 30 {
		t.Fatalf("expected both sums to be 30, got legacy=%d modern=%d", legacySum, modernSum)
	}
}

func TestLegacyProviderRawXML(t *testing.T) {
	legacy := NewLegacyXMLDataProvider(map[string]string{"u1": "42"})
	raw, err := legacy.FetchXML("u1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if raw == "" {
		t.Fatal("expected non-empty raw XML")
	}
}

func TestDemo(t *testing.T) {
	legacy := NewLegacyXMLDataProvider(map[string]string{"u1": "100"})
	var provider DataProvider = NewXMLToJSONAdapter(legacy)
	rec, err := provider.FetchJSON("u1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	t.Logf("adapted record: %v", rec)
}
