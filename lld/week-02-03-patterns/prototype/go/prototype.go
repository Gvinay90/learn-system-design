// Package prototype implements the Prototype pattern via a Document with
// nested Metadata and section content, cloned via a deep-copy Clone method
// instead of being rebuilt field by field.
package prototype

type Metadata struct {
	Author string
	Tags   []string
}

// clone returns a Metadata with its own backing array for Tags, so mutating
// the clone's tag slice never affects the original.
func (m *Metadata) clone() *Metadata {
	tags := make([]string, len(m.Tags))
	copy(tags, m.Tags)
	return &Metadata{Author: m.Author, Tags: tags}
}

type Document struct {
	Title    string
	Meta     *Metadata
	Sections []string
	Props    map[string]string
}

// Clone performs a deep copy: the nested Metadata, the Sections slice, and
// the Props map are all duplicated, not shared, with the original.
func (d *Document) Clone() *Document {
	sections := make([]string, len(d.Sections))
	copy(sections, d.Sections)

	props := make(map[string]string, len(d.Props))
	for k, v := range d.Props {
		props[k] = v
	}

	return &Document{
		Title:    d.Title,
		Meta:     d.Meta.clone(),
		Sections: sections,
		Props:    props,
	}
}
