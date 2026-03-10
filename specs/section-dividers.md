# Section Dividers

> How book parts get their own .docx title pages using the existing pipeline
> with minimal code changes.

---

## Overview

Each book is divided into parts (typically 8 per book). Each part needs a
**section divider** — a short .docx file that acts as a title page when the
book is assembled into a collection. Sections are treated as pipeline items
alongside essays, sharing the same file structure, state management, and
export path. Most pipeline stages are auto-skipped; only draft2 and export
do real work.

---

## File Naming

Section exports follow this pattern:

```
cSection - 2026 - {Book}.{Part}.00 {Part Name}.docx
```

Examples:
```
cSection - 2026 - I.01.00 Everything Is a Rate of Change.docx
cSection - 2026 - I.03.00 What Your Kitchen Knows.docx
cSection - 2026 - II.04.00 Signals and Information.docx
cSection - 2026 - III.08.00 Coda — The Nature of Mathematics.docx
```

Essays in the same part start at `.01`:
```
cChapter - 2026 - I.03.01 The Cheerios Effect.docx
cChapter - 2026 - I.03.02 The Coffee Ring Effect.docx
```

When the collection is assembled (outside this pipeline), the `.00` section
sorts before the `.01` first essay.

---

## Metadata

Each section has a `.meta.yaml` identical in structure to essay metadata,
with one addition — a `type` field:

```yaml
slug: section-i-03
title: "What Your Kitchen Knows"
type: section
book: I
part: 3
part_title: "What Your Kitchen Knows"
order: 0
status: pending
model: claude-sonnet-4-20250514
created: 2026-03-10
```

The `type` field distinguishes sections from essays:

| Value | Meaning |
|-------|---------|
| `essay` | Default. Normal pipeline processing. |
| `section` | Section divider. Most stages auto-skipped. |

The `type` field is displayed in the dashboard/UI table alongside status,
book, part, and order.

---

## Idea Files

Each section gets an idea file scaffolded from the Plan for its book.

**`ideas/section-i-03.md`:**
```markdown
# What Your Kitchen Knows

**Slug:** section-i-03

## Section

This is a section divider for Book I, Part 3.

## The Part

What Your Kitchen Knows — essays about the mathematics hiding in domestic
mysteries: why Cheerios clump, why coffee rings form, why the shower curtain
misbehaves.

## Placement

- **Book:** I
- **Part 3:** What Your Kitchen Knows
- **Order:** 0
```

The idea file carries enough context for the draft2 stage to write the
introductory paragraph. The "Part" section summarizes what the reader is
about to encounter, drawn from the Plan for the Book document.

---

## Pipeline Flow

Sections traverse all 8 stages. At stages that are auto-skipped, the
pipeline copies the previous stage's `.md` verbatim into the next stage
directory and marks the meta.yaml as complete. No Claude API call is made
for skipped stages. No tokens are consumed. No cost is incurred.

| Stage | Behavior for sections |
|-------|----------------------|
| **ideas** | Real content. Scaffolded from the Plan for the Book. |
| **research** | Auto-skip. Copy ideas .md forward. |
| **outline** | Auto-skip. Copy forward. |
| **draft** | Auto-skip. Copy forward. |
| **factcheck** | Auto-skip. Copy forward. |
| **illustrate** | Auto-skip. Copy forward. |
| **draft2** | Real Claude call. Writes the 2–3 sentence introductory paragraph + image tag. |
| **export** | Normal. md2docx → imagerender → imageswap → upgradeDocx. |

### Auto-skip logic

At the start of `processEssay`, before building the prompt:

1. Read the essay's `type` from meta.yaml.
2. If `type` is `section` and the current stage is in the skip list
   (research, outline, draft, factcheck, illustrate):
   - Read the .md from the previous stage directory.
   - Write it verbatim to the current stage directory.
   - Write a meta.yaml marking this stage as complete (status: final,
     tokens: 0, cost: 0).
   - Return — no API call.
3. Otherwise, proceed normally.

This requires no changes to state management, cycle selection, dashboard,
or any code outside `processEssay`.

---

## Draft2 Stage (the only AI stage for sections)

The draft2 prompt for sections receives:

- The idea file (part name, book, part description)
- The Plan for the Book document (book arc, all parts, all essay assignments)

It produces:

- A short introductory paragraph (2–3 sentences maximum). The paragraph
  sets up what this part of the book explores, in a voice consistent with
  the book's character. It does not summarize — it invites.
- An image tag referencing the placeholder image: `[[IMG:section-placeholder.png|Insert cartoon here]]`

The prompt should emphasize brevity and a light touch. The paragraph should
feel like a breath between chapters, not an announcement.

**Example output for Book I, Part 3:**

```markdown
# What Your Kitchen Knows

Your kitchen is a physics laboratory with terrible safety protocols.
Every morning, your cereal, your coffee, and your shower conspire to
demonstrate graduate-level mathematics — without your permission.

[[IMG:section-placeholder.png|Insert cartoon here]]
```

---

## Placeholder Image

A single static PNG image used by all sections:

- **Dimensions:** 4.5" × 3" at 300 DPI (1350 × 900 pixels)
- **Content:** The text "Insert cartoon here" centered on a white background
- **Font:** Simple sans-serif, medium gray, centered horizontally and vertically
- **Location:** `images/shared/section-placeholder.png`

Each section references this file via the `[[IMG:section-placeholder.png|...]]`
tag. The imageswap tool embeds a copy of the PNG bytes into each section's
.docx — the sections do not share a single file in the output; each .docx
contains its own embedded copy.

These placeholder images are intended to be replaced later with commissioned
or hand-drawn cartoons specific to each section.

---

## Source of Truth

The **Plan for Book** documents are the authoritative source for:

- Part names (used in section titles and filenames)
- Part descriptions (used in idea files as context for draft2)
- Number of parts per book (determines how many sections exist)
- Essay assignments per part (available as context for the draft2 prompt)

The current "Table of Contents" in `design/` is renamed to "Plan for Book I."
Plans for Book II and Book III are written from the straw man before sections
are scaffolded.

---

## Scaffolding

The `scaffold` tool is extended to create section idea files alongside essay
idea files. For each part in each book, it creates:

- `ideas/section-{book}-{part}.md`
- `ideas/section-{book}-{part}.meta.yaml`

Slug convention: `section-i-03`, `section-ii-01`, `section-iii-08`, etc.

The scaffold reads part names and descriptions from the Plan for the Book
documents (or from the hardcoded essay list, which mirrors the plans).

---

## Count

Based on the three-book straw man:

| Book | Parts | Sections |
|------|-------|----------|
| I | 8 | 8 |
| II | 8 | 8 |
| III | 8 | 8 |
| **Total** | **24** | **24** |

These 24 sections join the ~105 essays in the pipeline, bringing the total
item count to ~129.
