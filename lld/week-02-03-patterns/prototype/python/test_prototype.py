from prototype import Document, Metadata


def new_test_document() -> Document:
    return Document(
        title="Q3 Report",
        meta=Metadata(author="Vinay", tags=["finance", "draft"]),
        sections=["Intro", "Numbers"],
        props={"status": "draft"},
    )


def test_clone_produces_equal_but_distinct_values():
    original = new_test_document()
    clone = original.clone()

    assert clone.title == original.title
    assert clone.meta.author == original.meta.author
    assert clone is not original
    assert clone.meta is not original.meta


def test_mutating_clone_tags_does_not_affect_original():
    original = new_test_document()
    clone = original.clone()

    clone.meta.tags[0] = "mutated"
    clone.meta.tags.append("extra")

    assert original.meta.tags[0] == "finance"
    assert len(original.meta.tags) == 2


def test_mutating_clone_sections_does_not_affect_original():
    original = new_test_document()
    clone = original.clone()

    clone.sections[0] = "Rewritten Intro"
    clone.sections.append("Appendix")

    assert original.sections[0] == "Intro"
    assert len(original.sections) == 2


def test_mutating_clone_props_does_not_affect_original():
    original = new_test_document()
    clone = original.clone()

    clone.props["status"] = "final"
    clone.props["reviewer"] = "Asha"

    assert original.props["status"] == "draft"
    assert "reviewer" not in original.props


def test_mutating_clone_author_does_not_affect_original():
    original = new_test_document()
    clone = original.clone()

    clone.meta.author = "Someone Else"

    assert original.meta.author == "Vinay"


def test_shallow_copy_would_have_shared_nested_state():
    import copy

    original = new_test_document()
    shallow = copy.copy(original)

    shallow.meta.tags.append("leaks-into-original")

    assert "leaks-into-original" in original.meta.tags
