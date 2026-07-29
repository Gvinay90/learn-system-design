package prototype

import "testing"

func newTestDocument() *Document {
	return &Document{
		Title:    "Q3 Report",
		Meta:     &Metadata{Author: "Vinay", Tags: []string{"finance", "draft"}},
		Sections: []string{"Intro", "Numbers"},
		Props:    map[string]string{"status": "draft"},
	}
}

func TestCloneProducesEqualButDistinctValues(t *testing.T) {
	original := newTestDocument()
	clone := original.Clone()

	if clone.Title != original.Title {
		t.Fatalf("expected equal titles, got %s vs %s", clone.Title, original.Title)
	}
	if clone.Meta.Author != original.Meta.Author {
		t.Fatalf("expected equal authors")
	}
	if clone == original {
		t.Fatalf("clone should be a distinct pointer from original")
	}
	if clone.Meta == original.Meta {
		t.Fatalf("clone.Meta should be a distinct pointer from original.Meta")
	}
}

func TestMutatingCloneTagsDoesNotAffectOriginal(t *testing.T) {
	original := newTestDocument()
	clone := original.Clone()

	clone.Meta.Tags[0] = "mutated"
	clone.Meta.Tags = append(clone.Meta.Tags, "extra")

	if original.Meta.Tags[0] != "finance" {
		t.Fatalf("mutating clone's tags leaked into original: %v", original.Meta.Tags)
	}
	if len(original.Meta.Tags) != 2 {
		t.Fatalf("appending to clone's tags leaked into original: %v", original.Meta.Tags)
	}
}

func TestMutatingCloneSectionsDoesNotAffectOriginal(t *testing.T) {
	original := newTestDocument()
	clone := original.Clone()

	clone.Sections[0] = "Rewritten Intro"
	clone.Sections = append(clone.Sections, "Appendix")

	if original.Sections[0] != "Intro" {
		t.Fatalf("mutating clone's sections leaked into original: %v", original.Sections)
	}
	if len(original.Sections) != 2 {
		t.Fatalf("appending to clone's sections leaked into original: %v", original.Sections)
	}
}

func TestMutatingClonePropsDoesNotAffectOriginal(t *testing.T) {
	original := newTestDocument()
	clone := original.Clone()

	clone.Props["status"] = "final"
	clone.Props["reviewer"] = "Asha"

	if original.Props["status"] != "draft" {
		t.Fatalf("mutating clone's props leaked into original: %v", original.Props)
	}
	if _, ok := original.Props["reviewer"]; ok {
		t.Fatalf("adding a key to clone's props leaked into original: %v", original.Props)
	}
}

func TestMutatingCloneAuthorDoesNotAffectOriginal(t *testing.T) {
	original := newTestDocument()
	clone := original.Clone()

	clone.Meta.Author = "Someone Else"

	if original.Meta.Author != "Vinay" {
		t.Fatalf("mutating clone's author leaked into original: %s", original.Meta.Author)
	}
}
