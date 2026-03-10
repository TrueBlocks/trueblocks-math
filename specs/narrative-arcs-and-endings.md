# Narrative Arcs and Endings

> How individual essays are structured and how they land. This spec replaces the
> random-assignment system with an editorial framework where arc and ending are
> chosen based on the essay's topic, position in its book, and role in the series.

---

## The 15 Narrative Arcs

Each arc defines the structural beats of an essay — the shape of the reader's
journey from first sentence to last. An essay is ~1,200–1,600 words (4–6 minutes),
so these beats are rhythmic markers, not rigid sections.

Arcs and endings apply only to items with `type: essay`. Section dividers
(`type: section`) and book introductions (`type: introduction`) do not have
arc or ending assignments.

### Original 8

| # | Arc | Beats | Best for |
|---|-----|-------|----------|
| 1 | **The Slow Build** | Quiet scene → gradual layers → unified click → return transformed | Everyday topics where the math sneaks up on you |
| 2 | **The Cold Open Mystery** | Mystery first → investigate → eliminate suspects → solution | Unsolved problems, surprising results |
| 3 | **The Paradox** | Impossible statement → build tools → resolve → intuition updated | Counterintuitive results (Banach-Tarski, birthday paradox) |
| 4 | **The Demolition** | Textbook answer → systematic debunking → build correct → why myth stuck | Misconceptions (ice slipperiness, shower curtain) |
| 5 | **The Dual Timeline** | Interleave history + modern math → converge | Topics with rich human stories (Noether, Euler, Ramanujan) |
| 6 | **The Zoom** | Camera through scales, same math at each level | Scale-invariant topics (fractals, constructal law, power laws) |
| 7 | **The Argument** | Present competing explanations → evidence → verdict or honest ambiguity | Debated topics (Mpemba effect, fine-tuning) |
| 8 | **The Catalog of Wonders** | 4–5 escalating instances of same principle, most jaw-dropping last | One principle with many manifestations (derivative chains, Fibonacci) |

### New 7

| # | Arc | Beats | Best for |
|---|-----|-------|----------|
| 9 | **The Confession** | "I used to think X. Here's why. Then this changed my mind." | Topics where the common explanation is wrong |
| 10 | **The Letter** | Written to a specific person — intimate, direct, second-person | Biographical essays, philosophical closers |
| 11 | **The Walk** | Narrated as a journey through a place, noticing math along the way | Sensory/spatial topics (sound, light, architecture) |
| 12 | **The Countdown** | Start with the answer, peel back layers of "but why?" to bedrock | Deep topics where the real story is the foundation |
| 13 | **The Recipe** | Step-by-step instructions for *doing* something, math revealed in each step | Procedural topics (shuffling, cake-cutting, origami) |
| 14 | **The Debate** | Two figures argue across centuries, reader as judge | Historical conflicts (Cantor vs. Kronecker, EPR vs. Bohr) |
| 15 | **The Inheritance** | A concept passed from mind to mind across generations, each transforming it | Topics with long intellectual lineages (calculus of variations, topology) |

---

## The 6 Ending Types

Not every essay should land on the same emotional note. The ending type is chosen
based on the topic's nature and the essay's role in its book.

| Ending | Emotional note | When to use |
|--------|----------------|-------------|
| **Resolution** | "Now you see." Satisfying closure. | Early in a book or part. Topics with clean answers. Builds reader confidence. |
| **Provocation** | "But wait — what about...?" Opens a door. | Mid-book. Creates forward momentum. Best when the next essay picks up the thread. |
| **Awe** | "And it keeps going." Vertigo of scale or beauty. | The "zoom" and "catalog" essays. Cosmic topics. When the math is genuinely breathtaking. |
| **Honesty** | "We don't know." Unsolved, open, frank. | Late in a book. Frontier topics. When pretending to resolve would be dishonest. |
| **Communion** | "You already knew this." Returns to the body, to intuition. | Book closers. Brings the reader home. Connects abstract back to the felt. |
| **Handoff** | The essay's closing image becomes the next essay's opening. | Within a part, to link essays that share a thread. Subtle — the reader may not notice. |

---

## How Arc and Ending Interact

The arc shapes the journey. The ending shapes the landing. They are independent
choices, though some combinations are more natural:

| Arc | Natural endings | Unusual but effective |
|-----|-----------------|---------------------|
| Slow Build | Resolution, Communion | Awe (if the "click" is enormous) |
| Cold Open Mystery | Resolution, Honesty | Provocation (if the mystery deepens) |
| Paradox | Resolution, Awe | Honesty (if the paradox is truly unresolved) |
| Demolition | Resolution | Provocation (if the correct answer raises new questions) |
| Dual Timeline | Communion, Awe | Handoff (historical figure's question becomes next essay's topic) |
| Zoom | Awe | Honesty (if the smallest or largest scale is unknown) |
| Argument | Honesty, Resolution | Provocation (if the jury is still out) |
| Catalog of Wonders | Awe | Communion (if the last wonder is personal) |
| Confession | Communion | Resolution (if the new understanding is clean) |
| Letter | Communion, Provocation | Honesty (if the letter admits what isn't known) |
| Walk | Communion | Awe (if the walk ends at a vista) |
| Countdown | Awe, Honesty | Resolution (rare — the bedrock is usually mysterious) |
| Recipe | Resolution | Communion (the reader can now *do* this thing) |
| Debate | Honesty, Resolution | Provocation (if the debate is ongoing) |
| Inheritance | Awe, Communion | Handoff (the inheritance continues to the next essay) |

---

## Position-Sensitive Defaults

Arc and ending assignments should follow a loose pattern based on position:

### Within a part (typically 3–6 essays)

- **First essay:** Accessible arc (Slow Build, Walk, Recipe, Confession). Resolution or Communion ending.
- **Middle essays:** Any arc. Provocation or Handoff endings to maintain momentum.
- **Last essay:** Ambitious arc (Zoom, Catalog, Countdown, Debate). Awe or Honesty ending.

### Within a book

- **Early parts:** More Resolution and Communion endings. Build the reader's confidence.
- **Middle parts:** Mix of everything. Provocation endings create forward drive.
- **Late parts:** More Honesty and Awe endings. The book grows in ambition.
- **Final part/coda:** Communion. Bring the reader home.

### Across the series

- **Book I:** Weighted toward Resolution and Communion. The reader is learning to see.
- **Book II:** Balanced. Resolution where patterns are clear, Honesty where they aren't.
- **Book III:** Weighted toward Awe and Honesty. The reader is comfortable with open questions.

These are tendencies, not rules. A Book I essay about an unsolved problem should end
with Honesty. A Book III essay about something beautiful and solved can end with Resolution.
The pattern should be felt, not announced.

---

## Assignment

Arc and ending are recorded in each essay's `.meta.yaml`:

```yaml
arc: slow-build
ending: resolution
```

Assignment happens during the editorial/planning phase — when essays are organized
into books and parts. The outline and draft prompts use the arc and ending to shape
the AI's output.

The `RandomArc()` function is replaced by explicit editorial assignment. When an arc
is not yet assigned, the pipeline should skip the essay rather than guess.

---

## Biographical Enrichment

Some arcs naturally accommodate biographical material:

- **Dual Timeline** — interleave a mathematician's story with the modern understanding
- **Letter** — write to or from a historical figure
- **Inheritance** — trace a concept through generations of thinkers
- **Debate** — dramatize a real historical disagreement

Rather than creating standalone "biography essays," the preferred approach is to
use biographical material as enrichment *within* technical essays. The Lives in
Mathematics origin document provides source material for this:

| Technical essay | Biographical enrichment | Source |
|----------------|------------------------|--------|
| noethers-theorem | Emmy Noether's life, Hilbert's defense | Lives #2 |
| eulers-bridges | Euler's blindness and productivity | Lives #1 |
| hear-shape-drum | Ramanujan's intuition about eigenvalues | Lives #3 |
| arrows-impossibility | Arrow's Nobel, social choice as mathematics | — |
| navier-stokes | Generations of failed attempts | — |
| banach-tarski | The axiom of choice controversy | Lives #6 (Cantor) |
| self-organized-criticality | Per Bak's combative personality | — |
| algorithms-blind-spots | Turing and Gödel | Lives #8 |

A few essays are biographical enough to stand on their own:
- **ramanujan-notebooks** — the letter, the notebooks, what "discovery" means
- **galois-last-night** — group theory invented the night before a duel

These are interspersed among technical essays, not grouped into a "biography" section.
