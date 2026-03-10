# Book Introductions

> Each book gets a short introductory chapter that appears before Part 1.
> It flows through the pipeline with some stages auto-skipped.

---

## Overview

Each of the three books has a single Introduction — 1.5 to 2 pages (~400–500
words) of prose that opens the book. It appears before the first section
divider. It is the reader's first encounter with the book's voice and territory.

The introduction hints at the journey ahead without revealing the structure.
It does not summarize the parts or list the essays. A careful reader may sense
the arc; a casual reader simply feels welcomed.

---

## File Naming

```
cChapter - 2026 - I.00.00 Introduction.docx
cChapter - 2026 - II.00.00 Introduction.docx
cChapter - 2026 - III.00.00 Introduction.docx
```

Uses `cChapter` (not `cSection`) — this is real prose. The `.00.00` sorts
before `cSection - 2026 - I.01.00 ...` in the collection.

---

## Metadata

```yaml
slug: intro-i
title: "Introduction"
type: introduction
book: I
part: 0
part_title: "Introduction"
order: 0
status: pending
model: claude-sonnet-4-20250514
created: 2026-03-10
```

The `type: introduction` value joins `essay` and `section` as the three
pipeline item types.

---

## Idea File

**`ideas/intro-i.md`:**
```markdown
# Introduction

**Slug:** intro-i

## Book

Everything Is a Rate of Change

## Book Character

Warm, surprising, grounded. The mathematics hidden in what you can touch,
hear, and see. Starts with the feeling in an elevator. Ends with a glimpse
upward. The reader finishes with new eyes.

## Parts

1. Everything Is a Rate of Change
2. Your Body Already Knows
3. What Your Kitchen Knows
4. What Nature Computed
5. Hidden in Plain Sight
6. The Mathematics of Beautiful Things
7. The Mathematics of Living Together
8. A Glimpse Upward

## Placement

- **Book:** I
- **Part:** 0
- **Order:** 0
```

The idea file carries the book title, character description, and part names
— all drawn from the Plan for the Book document. This gives later stages
everything they need.

---

## Pipeline Flow

| Stage | Behavior for introductions |
|-------|---------------------------|
| **ideas** | Real content. Scaffolded from the Plan for the Book. |
| **research** | Auto-skip. Copy ideas .md forward. |
| **outline** | Real Claude call. Produces a 4–5 beat outline for ~450 words. |
| **draft** | Real Claude call. Writes the introduction from the outline. |
| **factcheck** | Auto-skip. Copy draft .md forward. |
| **illustrate** | Auto-skip. Copy forward. |
| **draft2** | Real Claude call. Polish pass — tighten, match book voice. |
| **export** | Normal md2docx and upgradeDocx. No imagerender or imageswap (no images). |

**Skipped stages:** research, factcheck, illustrate.

**Active stages:** ideas, outline, draft, draft2, export.

The auto-skip mechanism is identical to sections: copy previous stage's .md
forward, mark meta.yaml complete, zero tokens, zero cost.

---

## Skip List by Type

All three item types use the same auto-skip mechanism. The only difference
is which stages appear in each type's skip list:

| Stage | essay | section | introduction |
|-------|-------|---------|--------------|
| ideas | run | run | run |
| research | run | skip | skip |
| outline | run | skip | run |
| draft | run | skip | run |
| factcheck | run | skip | skip |
| illustrate | run | skip | skip |
| draft2 | run | run | run |
| export | run | run | run |

This can be represented as a simple lookup: given a `type` and a `stage`,
is it in the skip list? One map or switch statement in `processEssay`.

---

## Prompt Guidance

### Outline stage

The outline prompt for introductions receives the idea file (book title,
character, part names). It produces a short outline (~4–5 beats) for a
400–500 word piece. The outline should:

- Open with a concrete, sensory image (not an abstraction)
- Build toward the book's governing idea without naming it outright
- End with an invitation, not a table of contents

### Draft stage

The draft prompt writes from the outline. Requirements:

- 400–500 words maximum (1.5–2 pages)
- Match the book's character (warm for I, intellectual for II, philosophical for III)
- No images, no equations, no section headers
- Hint at the journey without mapping it
- The reader should feel curiosity, not obligation

### Draft2 stage

Polish pass. Tighten language, cut anything that sounds like a sales pitch
or a course syllabus. The introduction should read like the opening of a
conversation, not the preface of a textbook.

---

## Arc and Ending

Introductions do not have narrative arc or ending type assignments.
The arc/ending system applies only to items with `type: essay`. The
pipeline's arc assignment logic ignores `type: introduction` and
`type: section`.

---

## Export

At export, introductions follow the same md2docx and upgradeDocx path
as essays, but skip imagerender and imageswap since there are no images.

The export logic can check: if no `[[IMG:...]]` or `[[R:...]]` tags are
present in the markdown, skip the image processing steps. This is a
general-purpose check that works for introductions without special-casing
the type.

---

## Count

| Book | Introductions |
|------|---------------|
| I | 1 |
| II | 1 |
| III | 1 |
| **Total** | **3** |

These 3 introductions join 24 sections and ~105 essays, bringing the
total pipeline item count to ~132.
