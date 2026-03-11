package pipeline

type StructureHint struct {
	Name        string
	OutlineHint string
	DraftHint   string
}

type EntryHint struct {
	Name      string
	DraftHint string
}

type RegisterHint struct {
	Name      string
	DraftHint string
}

type MathVisHint struct {
	Name        string
	OutlineHint string
	DraftHint   string
}

var structureHints = map[string]StructureHint{
	"narrative": {
		Name: "narrative",
		OutlineHint: `ESSAY STRUCTURE: Narrative
Build this essay around a story with characters, a setting, and events that unfold.
The mathematics emerges from the narrative — it should feel discovered, not inserted.`,
		DraftHint: `STRUCTURE: Narrative. Write this as a story. The reader follows characters or events,
and the mathematics grows out of what happens. Don't pause the story to explain — weave
the math into the telling.`,
	},
	"in-medias-res": {
		Name: "in-medias-res",
		OutlineHint: `ESSAY STRUCTURE: In Medias Res
Open mid-action — the reader is dropped into a moment already underway. Fill in context
only as needed, through the action itself. The mathematics is part of the unfolding scene.`,
		DraftHint: `STRUCTURE: In Medias Res. Start in the middle. No setup paragraph, no throat-clearing.
The reader should be disoriented for exactly two sentences, then oriented by the third.
Backstory arrives through action, not exposition.`,
	},
	"frame": {
		Name: "frame",
		OutlineHint: `ESSAY STRUCTURE: Frame Narrative
Build a story-within-a-story. An outer frame (a conversation, a memory, a letter discovered)
contains the real mathematical content. The frame provides emotional context; the inner
story carries the math.`,
		DraftHint: `STRUCTURE: Frame Narrative. Establish the outer frame quickly — a grandmother's story, a
found notebook, a conversation overheard. Let the inner story carry the weight. Return to
the frame at the end to land the emotional note.`,
	},
	"braided": {
		Name: "braided",
		OutlineHint: `ESSAY STRUCTURE: Braided
Two seemingly unrelated threads — different topics, different eras, different scales —
that converge. Plan the convergence point carefully. The reader should not see the
connection coming until it arrives.`,
		DraftHint: `STRUCTURE: Braided. Alternate between two threads that seem to have nothing in common.
Each thread should be compelling on its own. When they converge, the reader should feel
the connection physically — like two melodies resolving into harmony.`,
	},
	"reverse-chronology": {
		Name: "reverse-chronology",
		OutlineHint: `ESSAY STRUCTURE: Reverse Chronology
Open with the conclusion — the surprising result, the final state. Then work backward
through the chain of causes. Each section peels back one layer of "how did we get here?"`,
		DraftHint: `STRUCTURE: Reverse Chronology. Start with the punchline. The reader knows the ending
from sentence one. The pleasure is in the unwinding — each section reveals an earlier
cause, an older insight, a deeper layer. End at the true beginning.`,
	},
	"compare-contrast": {
		Name: "compare-contrast",
		OutlineHint: `ESSAY STRUCTURE: Compare and Contrast
Two things that seem different but share hidden math — or two things that seem similar
but are mathematically alien. Build both sides fairly before revealing the connection
(or the gulf).`,
		DraftHint: `STRUCTURE: Compare and Contrast. Present both subjects with equal care and specificity.
Don't telegraph which direction the comparison will go. The reveal — "these are the same"
or "these are nothing alike" — should land with genuine surprise.`,
	},
	"catalogue": {
		Name: "catalogue",
		OutlineHint: `ESSAY STRUCTURE: Catalogue
A curated list of instances where the same mathematical principle appears. Start familiar,
escalate to surprising. The power is in accumulation — each instance adds evidence until
the pattern becomes undeniable.`,
		DraftHint: `STRUCTURE: Catalogue. Each entry should be punchy — a paragraph or two, no more. Let the
list build momentum. Resist the urge to explain the connecting principle too early. The
most shocking example goes last. By the end, the reader should be seeing the pattern
everywhere.`,
	},
	"exposition": {
		Name: "exposition",
		OutlineHint: `ESSAY STRUCTURE: Exposition
No story, no characters, no narrative arc. A clean, direct walk through an idea. This
structure earns its place through clarity and elegance — every sentence advances understanding.`,
		DraftHint: `STRUCTURE: Exposition. No story needed. Write with the directness of a well-organized
mind explaining something it finds beautiful. Every paragraph should move the reader
forward. Cut anything that doesn't teach or illuminate.`,
	},
	"argument": {
		Name: "argument",
		OutlineHint: `ESSAY STRUCTURE: Argument
Take a position and defend it. State the claim early. Anticipate objections. Marshal
evidence. The mathematics is your proof — not in the formal sense, but in the rhetorical
sense. The reader should end convinced.`,
		DraftHint: `STRUCTURE: Argument. Be bold. State your thesis in the first paragraph. Don't hedge.
Address the strongest counterargument, not the weakest. Let the math do the persuading.
End with the implications of being right.`,
	},
	"letter": {
		Name: "letter",
		OutlineHint: `ESSAY STRUCTURE: Letter
Addressed to a specific person — a barista, a child, a historical figure, a skeptic.
The mathematics is a gift being unwrapped for the recipient. The letter format creates
intimacy and allows a natural, conversational unfolding.`,
		DraftHint: `STRUCTURE: Letter. Use "you" naturally. Write as if the recipient is sitting across
from you. The tone is generous — you're sharing something you love. Let digressions
happen (they would in a real letter). Sign off with warmth.`,
	},
	"field-guide": {
		Name: "field-guide",
		OutlineHint: `ESSAY STRUCTURE: Field Guide / How-To
Instructional voice. The reader is given steps to follow or observations to make.
The mathematics emerges from doing — the reader has agency. Structure as actual
instructions that someone could follow.`,
		DraftHint: `STRUCTURE: Field Guide. Write in imperative voice when giving instructions. "Take a
piece of paper. Fold it in half." The reader should feel they could do this right now.
Name the mathematical principle only after the reader has felt it through action.`,
	},
	"meditation": {
		Name: "meditation",
		OutlineHint: `ESSAY STRUCTURE: Meditation
Contemplative, essayistic in the Montaigne sense. No clear arc, no plot. The mind
wanders through the territory of an idea, making connections, doubling back, arriving
nowhere in particular — but the reader ends up somewhere they didn't expect.`,
		DraftHint: `STRUCTURE: Meditation. Let the mind wander. This essay thinks on the page. There is no
thesis to defend, no mystery to solve. Just a mind turning an idea over, noticing its
facets. The rhythm should be unhurried. The ending arrives, not concludes.`,
	},
	"dialogue": {
		Name: "dialogue",
		OutlineHint: `ESSAY STRUCTURE: Dialogue
Two voices — a student and a teacher, two colleagues, a parent and child. The mathematics
emerges through conversation. Disagreement is productive. Neither voice is a strawman.`,
		DraftHint: `STRUCTURE: Dialogue. Write as a genuine conversation between two people who think
differently. Both voices should sound like real humans — interruptions, incomplete
thoughts, moments of confusion. The math arrives through the exchange, not despite it.`,
	},
	"q-and-a": {
		Name: "q-and-a",
		OutlineHint: `ESSAY STRUCTURE: Q&A
Question-driven. Each section opens with a question the reader is likely asking. Answers
lead to new questions. The escalation from simple to profound should feel inevitable.`,
		DraftHint: `STRUCTURE: Q&A. Let the questions drive. Each answer should satisfy and simultaneously
open the next question. Start with the question a child would ask. End with the question
a mathematician can't answer. The reader's curiosity should accelerate, not wind down.`,
	},
	"vignette-collage": {
		Name: "vignette-collage",
		OutlineHint: `ESSAY STRUCTURE: Vignette Collage
Three to four tiny scenes — different places, different people, different moments — with
no connective tissue between them. The mathematics is the invisible thread. The reader
must do the connecting.`,
		DraftHint: `STRUCTURE: Vignette Collage. Write each scene as a self-contained miniature. No
transitions between them — just white space. Each vignette should be vivid and specific.
The math is never stated explicitly; it's the pattern the reader discovers by reading
all the pieces together.`,
	},
	"chronicle": {
		Name: "chronicle",
		OutlineHint: `ESSAY STRUCTURE: Chronicle
Recount how something was discovered or evolved. Follow the timeline faithfully. The
pleasure is in the human detail — the wrong turns, the accidents, the people who
almost got there first.`,
		DraftHint: `STRUCTURE: Chronicle. Follow the history. Be faithful to the timeline but not enslaved
by it — skip the boring parts, linger on the turning points. People should feel alive:
give them names, locations, motivations. The math is what they were chasing.`,
	},
}

var entryHints = map[string]EntryHint{
	"cold-open": {
		Name:      "cold-open",
		DraftHint: `OPENING: Cold open. First sentence drops the reader into action or a striking fact. No setup, no context, no "Imagine if..." — just the thing itself. Orient the reader through the telling, not before it.`,
	},
	"on-the-math": {
		Name:      "on-the-math",
		DraftHint: `OPENING: Lead with the math. First sentence states the equation, the theorem, or the number. Then humanize it. The reader meets the abstraction first and discovers it has a body.`,
	},
	"question": {
		Name:      "question",
		DraftHint: `OPENING: Start with a question. Not rhetorical — genuine. The kind of question that once you hear it, you can't unhear it. The rest of the essay is the answer.`,
	},
	"contradiction": {
		Name:      "contradiction",
		DraftHint: `OPENING: State something that should be impossible. "This shouldn't work, but it does." The reader's disbelief is the engine — the essay is the proof.`,
	},
	"setting": {
		Name:      "setting",
		DraftHint: `OPENING: Paint a scene. Sensory details — light, sound, temperature, texture. The reader is somewhere specific before they know what the essay is about. The math grows out of the place.`,
	},
	"history": {
		Name:      "history",
		DraftHint: `OPENING: Start in the past. A date, a name, a place. "In 1854, a physician named John Snow removed the handle from a water pump." The reader enters through a door in history.`,
	},
	"at-the-end": {
		Name:      "at-the-end",
		DraftHint: `OPENING: Begin at the conclusion. The result, the consequence, the aftermath. "By the time they finished, the answer was three inches too short." Then unwind backward to show how we got there.`,
	},
	"voice": {
		Name:      "voice",
		DraftHint: `OPENING: The narrator's attitude is the hook. The first sentence establishes a personality — sardonic, wonder-struck, weary, delighted — before it establishes a topic. The reader stays for the company.`,
	},
	"data": {
		Name:      "data",
		DraftHint: `OPENING: A number. Cold, unadorned, unexplained. "0.0000000000000000000016." Or "23." Let the number sit for a moment. Then explain why it matters. The starkness is the hook.`,
	},
}

var registerHints = map[string]RegisterHint{
	"casual": {
		Name:      "casual",
		DraftHint: `REGISTER: Casual. Write like you're explaining this to a friend at a bar. Contractions, sentence fragments, the occasional "look." The math is serious; the delivery isn't.`,
	},
	"academic": {
		Name:      "academic",
		DraftHint: `REGISTER: Academic. Precise language, careful qualifications, proper attribution. Not dry — engaged and passionate about accuracy. The reader should feel they're learning from someone who has spent years with this material.`,
	},
	"whimsical": {
		Name:      "whimsical",
		DraftHint: `REGISTER: Whimsical. Light, playful, willing to make an unexpected comparison or a bad pun if it illuminates. The math is a game — not trivial, but genuinely fun. Delight is the engine.`,
	},
	"lyrical": {
		Name:      "lyrical",
		DraftHint: `REGISTER: Lyrical. Pay attention to the sound of your sentences. Rhythm matters. Metaphors should be precise, not decorative. This register earns the right to be beautiful by being exact.`,
	},
	"dry": {
		Name:      "dry",
		DraftHint: `REGISTER: Dry. Humor through understatement and deadpan delivery. State extraordinary things as if they were ordinary. The contrast between the wild math and the flat delivery is the comedy.`,
	},
	"urgent": {
		Name:      "urgent",
		DraftHint: `REGISTER: Urgent. Short sentences. Active voice. This matters right now. Cut every word that doesn't earn its place. The reader should feel the essay is pulling them forward.`,
	},
	"warm": {
		Name:      "warm",
		DraftHint: `REGISTER: Warm. Personal memory, sensory detail, genuine affection for the subject. The reader should feel invited, not lectured. Share the wonder without performing it.`,
	},
	"clinical": {
		Name:      "clinical",
		DraftHint: `REGISTER: Clinical. Observational, almost like scientific field notes. Report what you see without editorializing. Let the facts carry the weight. The restraint makes the extraordinary details hit harder.`,
	},
	"sardonic": {
		Name:      "sardonic",
		DraftHint: `REGISTER: Sardonic. Slightly cynical, world-weary, but secretly delighted. Mock the obvious, praise the unexpected. The reader should feel like they're in on a joke that most people miss.`,
	},
	"wonder-struck": {
		Name:      "wonder-struck",
		DraftHint: `REGISTER: Wonder-struck. Genuine amazement with no irony. This essay is written by someone who cannot believe how beautiful this is. Not performed awe — earned awe. The math itself is the spectacle.`,
	},
}

var mathVisHints = map[string]MathVisHint{
	"front-and-center": {
		Name: "front-and-center",
		OutlineHint: `MATH VISIBILITY: Front and Center
Equations, notation, and worked examples are the essay's spine. The reader should see
the math clearly — inline notation, display equations, step-by-step reasoning. Don't
hide it. The beauty IS the formalism.`,
		DraftHint: `MATH VISIBILITY: Front and Center. Show your work. Use inline math ($...$) freely and
display equations ($$...$$) when a result deserves emphasis. Walk through calculations
step by step. The reader came for the math — give it to them.`,
	},
	"woven-in": {
		Name: "woven-in",
		OutlineHint: `MATH VISIBILITY: Woven In
Mathematics is present but integrated into the narrative. Equations appear when they
illuminate but don't dominate. The reader absorbs the math through the story without
feeling like they stopped to study.`,
		DraftHint: `MATH VISIBILITY: Woven In. Use the math where it helps. An equation here, a ratio there.
But keep the prose flowing around it. The reader who skips every equation should still
follow the argument. The reader who reads them should get a bonus.`,
	},
	"buried": {
		Name: "buried",
		OutlineHint: `MATH VISIBILITY: Buried
The mathematics is present but invisible. No equations, no notation, no "let x be..."
The reader discovers they've been doing math only at the end — or never. The essay
teaches through story, analogy, and experience.`,
		DraftHint: `MATH VISIBILITY: Buried. No equations. No notation. No formal math at all. Use
analogies, physical intuition, sensory description. The reader should feel the
mathematics without ever seeing a symbol. If you must reference a number, make it
concrete — "twice as fast" not "velocity doubled."`,
	},
}

func StructureByName(name string) (StructureHint, bool) {
	h, ok := structureHints[name]
	return h, ok
}

func EntryByName(name string) (EntryHint, bool) {
	h, ok := entryHints[name]
	return h, ok
}

func RegisterByName(name string) (RegisterHint, bool) {
	h, ok := registerHints[name]
	return h, ok
}

func MathVisByName(name string) (MathVisHint, bool) {
	h, ok := mathVisHints[name]
	return h, ok
}
