# Essay Variation Spec

## The Problem

After reading ~100 generated essays, nearly every one follows the same shape:

> "I was doing some normal shit. I noticed a pattern. I was surprised.
> Here's the math. Isn't it interesting. Next time you're doing some shit,
> remember how cool math is."

The existing narrative arc system (15 arcs) varies the *plot shape* but leaves the
*texture* identical. A "Cold Open Mystery" and a "Slow Build" with the same register,
same entry style, and same setting feel much more similar than they should.

The goal: shatter sameness across 103 essays while maintaining a coherent authorial
voice and serving each book's larger arc.

---

## The Seven Axes of Variation

Each essay carries seven curated attributes. Together they form a **recipe** that makes
the essay feel unlike its neighbors while belonging to the same collection.

### 1. Arc (existing — 15 types)

The emotional/structural shape of the essay's narrative.

| Arc | Core Shape |
|-----|-----------|
| Slow Build | Calm scene → odd detail → math layers → reveal → return transformed |
| Cold Open Mystery | Baffling fact → obvious suspect → evidence → false lead → solution |
| Paradox | Impossible statement → why it feels wrong → tools → resolution → new normal |
| Demolition | Confident myth → first crack → collapse → rebuild → why myth persisted |
| Dual Timeline | Historical (A) / modern (B) threads racing toward convergence |
| Zoom | Continuous camera move through scales — same principle at each level |
| Argument | Phenomenon → Team A → Team B → evidence → verdict → what debate reveals |
| Catalog of Wonders | Instance after instance, power through accumulation, most shocking last |
| Confession | "I used to think X" → why it made sense → cracks → new understanding |
| Letter | Addressed to someone specific, math as a gift being unwrapped |
| Walk | Journey through a real place, sensory details ground math in location |
| Countdown | Answer given upfront, then drilling down through layers of "but why?" |
| Recipe | Actual instructions, imperative voice, math emerges from doing |
| Debate | Two historical figures, equal conviction, reader weighs both sides |
| Inheritance | Chain of handoffs across generations, idea evolving through each mind |

### 2. Ending (existing — 6 types)

How the essay resolves (or doesn't).

| Ending | Feel |
|--------|------|
| Resolution | The puzzle clicks shut. Satisfying, complete. |
| Awe | The reader is left staring at something vast. |
| Communion | The reader feels connected — to the author, to the world, to the math. |
| Provocation | The essay leaves a splinter. Something to argue with. |
| Honesty | "We don't know." The frontier is real and that's okay. |
| Handoff | The essay opens a door and invites the reader to walk through it alone. |

### 3. Structure (new — 16 types)

The rhetorical form — how the essay is *built*, independent of its narrative arc.

| Structure | Description |
|-----------|------------|
| narrative | A story with characters, setting, and events. The default — use sparingly. |
| in-medias-res | Opens mid-action, fills in context later. |
| frame | A story within a story. "My grandmother used to say..." |
| braided | Two seemingly unrelated threads that converge. |
| reverse-chronology | Start with the punchline, unpack backward. |
| compare-contrast | Two things that seem different but share hidden math (or vice versa). |
| catalogue | "Seven places π shows up where you'd never expect." |
| exposition | No story at all — a clean walk through an idea. |
| argument | Take a position. "Averages lie. Here's why you should care." |
| letter | Addressed to someone specific. "Dear barista, about that foam..." |
| field-guide | Instructional tone. "How to lose money slowly: a guide to compound fees." |
| meditation | Contemplative, essayistic in the Montaigne sense. No clear arc. |
| dialogue | Two voices, Socratic or conversational. |
| q-and-a | Question-driven structure. |
| vignette-collage | 3–4 tiny scenes, no connective tissue, the math is the thread. |
| chronicle | Recount how something was discovered or evolved over time. |

### 4. Entry (new — 9 types)

Where the reader lands on the first sentence.

| Entry | Opens with... |
|-------|--------------|
| cold-open | Dropped into action, no setup. |
| on-the-math | Lead with the equation or concept, then humanize it. |
| question | "Have you ever wondered why..." |
| contradiction | "This should be impossible, but..." |
| setting | Paint a scene, let the math emerge. |
| history | "In 1854, a cholera outbreak..." |
| at-the-end | "By the time she finished, the bridge was 3 inches too short." |
| voice | The narrator's attitude is the hook, not the content. |
| data | A number, a stat, a measurement — cold and unadorned. |

### 5. Register (new — 10 types)

The diction, tone, and level of formality. Overlaps with voice but doesn't destroy it —
a single author can write in all these registers. That's range, not inconsistency.

| Register | Feel |
|----------|------|
| casual | Explaining to a friend at a bar. |
| academic | Formal, careful, citation-flavored. |
| whimsical | Light, punny, irreverent. |
| lyrical | Rhythmic prose, metaphor-heavy. |
| dry | Humor through understatement. |
| urgent | Short sentences, punchy. "Here's what matters." |
| warm | Personal memory, sensory detail. |
| clinical | Observational, almost scientific field notes. |
| sardonic | Slightly cynical, world-weary humor. |
| wonder-struck | Genuine amazement, no irony. |

### 6. Setting (new — compound)

A compressed location + time tag. Not a controlled vocabulary — each essay gets a
bespoke setting string that the prompt interprets. Examples:

**Time dimension:**
- Historical deep past (ancient, medieval)
- Historical near past (1800s, early 1900s)
- Mid-20th century
- Contemporary / current day
- Near future (plausible tech, same society)
- Far future (solarpunk, post-scarcity, or dystopian)

**Place dimension:**
- Domestic (kitchen, garage, garden)
- Urban (city street, subway, office)
- Rural (farm, forest, open road)
- Institutional (school, hospital, courtroom, lab)
- Commercial (market, restaurant, factory floor)
- Exotic / travel (Istanbul, Kyoto, Buenos Aires)
- Digital / virtual (inside a game, an algorithm, a network)

**Geography / culture:**
- North America, Europe, East Asia, South Asia, Middle East, Africa, Latin America
- Each brings different measurement systems, cultural relationships to number

**Scale:**
- Microscopic (cells, atoms, grains of sand)
- Human-scale (rooms, streets, bodies)
- Architectural (buildings, bridges, cities)
- Planetary (orbits, tides, weather)
- Cosmic (stars, galaxies, the observable universe)

Setting is a free-form string in the metadata, e.g.:
- `"kitchen-morning"`, `"1854-London"`, `"near-future-Tokyo"`, `"medieval-Baghdad"`
- `"highway-afternoon"`, `"microscope-lab"`, `"digital-network"`

### 7. Math Visibility (new — 3 types)

How prominently the mathematics appears in the text.

| MathVis | Description |
|---------|------------|
| front-and-center | Equations, notation, worked examples are the essay's spine. |
| woven-in | Math is present but integrated into the narrative — the reader absorbs it without stopping. |
| buried | The reader doesn't realize they've learned mathematics until the last paragraph (or never). |

---

## Flow Principles

Hard rules ("no more than 2 catalogues per part") produce sterile regularity — exactly
the sameness we're trying to escape. Instead, attribute assignments follow *flow
principles* drawn from the same mathematics the books explore:

### Repetition with Variation

Like a musical theme that returns in different keys. If three consecutive essays use
warm register, the fourth in sardonic hits harder *because* of the warmth before it.
The pattern isn't "never repeat" — it's "repeat with enough transformation that the
echo adds meaning."

### The Clothoid Principle

Curvature changes gradually, not in steps. Book I doesn't jump from all-casual to
all-academic. The register drifts, like an Euler spiral, from warm into more challenging
territory as the reader's confidence grows.

### Benford's Law for Distribution

Some attributes should be common (warm register, narrative structure, setting-as-entry)
because they're the baseline that makes the rare ones feel rare. If every essay were a
dialogue or a letter, none would feel surprising. Uncommon forms get their power from
scarcity.

### Autocorrelation with Decay

Adjacent essays share *some* texture — they're in the same Part. But the correlation
decays with distance. Two essays in the same Part might share a register but differ in
structure and entry. Two essays in different Parts should feel like different rooms.

### Material Affinity

Each topic pulls toward a natural form. "The Cheerios Effect" (tiny domestic mystery)
wants warm register, slow-build, kitchen-morning. "Galois's Last Night" (Book III)
wants lyrical register, reverse-chronology, 1832-Paris. The attribute system's job is
to *follow that pull* instead of flattening everything to the same mold — and to notice
when two adjacent essays are pulling the same direction and nudge one.

---

## Metadata Schema

Each essay's `.meta.yaml` gains five new fields alongside the existing `arc` and `ending`:

```yaml
slug: snap-crackle-pop
title: "Snap, Crackle, and Pop"
type: essay
book: I
part: 1
part_title: "Everything Is a Rate of Change"
order: 1
status: pending
model: claude-sonnet-4-20250514

# Variation attributes (all 7 axes)
arc: slow-build
ending: resolution
structure: narrative
entry: setting
register: warm
setting: "elevator-morning"
math_visibility: woven-in

created: 2026-03-10
```

Sections and introductions carry empty strings for all variation attributes (they have
their own dedicated prompts).

---

## How Attributes Enter Prompts

Each attribute is injected into the pipeline stage where it has the most leverage:

| Stage | Attributes Used | Why |
|-------|----------------|-----|
| Research | setting | Guides examples, historical context, cultural references |
| Outline | arc, structure, entry, math_visibility | These shape the skeleton — *what goes where* |
| Draft | All seven | Full recipe — the draft is where texture lives |
| Factcheck | None | Fact-checking is attribute-agnostic |
| Illustrate | setting, math_visibility | Guides image style and density |
| Draft2 | register, arc | Final polish enforces diction; arc ensures emotional shape holds |

---

## Plan File Format

The Plan for Math Book files remain markdown (human-readable) with an enhanced table:

```markdown
| # | Type | Slug | Title | Hook | Hidden Math | Arc | Ending | Structure | Entry | Register | Setting | MathVis |
```

A companion YAML block or sidecar file provides machine-parseable data for the pipeline:

```yaml
# Plan for Math Book I — attributes.yaml
essays:
  - slug: snap-crackle-pop
    arc: slow-build
    ending: resolution
    structure: narrative
    entry: setting
    register: warm
    setting: "elevator-morning"
    math_visibility: woven-in
  - slug: derivative-chains
    arc: catalog-of-wonders
    ending: awe
    structure: catalogue
    entry: on-the-math
    register: wonder-struck
    setting: "various-domains"
    math_visibility: front-and-center
  # ...
```

The pipeline reads attributes from the YAML. The markdown table is the authoring surface.

---

## Book-Level Character Guidelines

Each book has a gravitational center — a default register and emotional territory that
most of its essays orbit, with deliberate departures for contrast.

### Book I — *Everything Is a Rate of Change*

**Center of gravity:** Warm, grounded, wonder-struck. Domestic and human-scale.
Most essays resolve. The reader finishes with new eyes.

**Register distribution:** Heavy warm and casual, with whimsical and wonder-struck
accents. Academic appears rarely (once or twice) as contrast. Sardonic appears once
at most.

**Setting tendency:** Kitchen, body, neighborhood, nature. Contemporary. Close-in.

**Math visibility:** Mostly woven-in. A few front-and-center (derivative chains,
Fourier). A few buried (the reader discovers the math was there all along).

### Book II — *The Hidden Architecture*

**Center of gravity:** More intellectual, more "wait, really?" Structures and systems
made visible. Some essays resolve, some leave the door open.

**Register distribution:** Shifts toward dry, academic, sardonic. Still warm in places
but earned — the reader is more sophisticated now. Urgent appears (information theory,
strategy). Clinical appears (signals, data).

**Setting tendency:** Wider geography. Cities, institutions, historical periods.
More variety in scale (microscopic to architectural).

**Math visibility:** More front-and-center. The reader has been trained by Book I
to handle equations. Still some buried essays for texture.

### Book III — *The Edge of Knowing*

**Center of gravity:** Philosophical, frontier-facing. Many topics don't resolve.
Heavy on awe and honesty endings.

**Register distribution:** Full range. Lyrical for the beautiful topics (Ramanujan,
Galois). Clinical for the unsolved problems. Wonder-struck for the cosmological.
Sardonic for the paradoxes. This book earns the right to use every register because
the reader has traveled with you through two books.

**Setting tendency:** Historical deep past to far future. Planetary and cosmic scale.
Global geography. Some essays are placeless — pure abstraction.

**Math visibility:** Oscillates between front-and-center (this is the frontier) and
buried (some of the deepest math hides in the simplest questions).

---

## Regeneration Strategy

Since all 103 essays regenerate from the ideas stage:

1. **Rewrite the three Plan files** — add 5 new attribute columns per essay, curated
   to serve book flow. Generate companion YAML.
2. **Regenerate idea files** — same slugs and titles, but idea content now includes
   all 7 attributes as seeds.
3. **Update EssayMeta** — add `structure`, `entry`, `register`, `setting`,
   `math_visibility` fields.
4. **Update prompts** — each stage incorporates relevant attributes (see table above).
5. **Wipe stages research through export** — delete all downstream files.
6. **Pipeline regenerates everything** — with the new attribute-enriched prompts.

Slugs, titles, book/part/order remain unchanged. Content changes radically.

---

## Dashboard Enhancement

Replace the single **Arc** column with an **Attributes** column:

**Default display:** Compact top-3 attributes, e.g.:
> `braided · sardonic · Tokyo`

**Popover on hover:** Full attribute card:
```
Structure:    braided
Entry:        cold-open
Register:     sardonic
Setting:      near-future, urban Tokyo
Arc:          wonder-driven
Ending:       provocation
Math:         buried
```

---

## Deferred: Guest Essayist Personas

Tracked in [issue #124](https://github.com/TrueBlocks/trueblocks-art/issues/124).
A `persona` field would assign some essays to named guest voices (Marguerite Torque,
Peppe Ratio, Irene Null, etc.). Implement after the core attribute system is proven.

---

## Risk Acceptance

Some structures are risky. Dialogue might not land. Q&A might feel gimmicky. A letter
to a barista about foam mathematics might be brilliant or terrible. We accept this.
Risky is human. The factcheck and draft2 stages exist to catch essays that don't land,
and a failed experiment is more interesting than a safe repetition.
