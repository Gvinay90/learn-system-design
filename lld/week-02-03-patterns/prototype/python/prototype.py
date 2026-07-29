"""Prototype pattern — cloning a Document with nested Metadata instead of
rebuilding it. See ../README.md for the design writeup.
"""
from __future__ import annotations

import copy
from dataclasses import dataclass, field
from typing import Dict, List


@dataclass
class Metadata:
    author: str
    tags: List[str] = field(default_factory=list)


@dataclass
class Document:
    title: str
    meta: Metadata
    sections: List[str] = field(default_factory=list)
    props: Dict[str, str] = field(default_factory=dict)

    def clone(self) -> "Document":
        # copy.deepcopy walks nested Metadata, the sections list, and the
        # props dict, so the clone shares no mutable state with the original.
        return copy.deepcopy(self)


def _demo() -> None:
    original = Document(
        title="Q3 Report",
        meta=Metadata(author="Vinay", tags=["finance", "draft"]),
        sections=["Intro", "Numbers"],
        props={"status": "draft"},
    )

    clone = original.clone()
    clone.meta.author = "Someone Else"
    clone.meta.tags.append("reviewed")
    clone.sections.append("Appendix")
    clone.props["status"] = "final"

    print("Original author:", original.meta.author)
    print("Clone author:", clone.meta.author)
    print("Original tags:", original.meta.tags)
    print("Clone tags:", clone.meta.tags)
    print("Original sections:", original.sections)
    print("Clone sections:", clone.sections)
    print("Original props:", original.props)
    print("Clone props:", clone.props)


if __name__ == "__main__":
    _demo()
