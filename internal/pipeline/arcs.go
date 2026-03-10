package pipeline

import randv2 "math/rand/v2"

type NarrativeArc struct {
	Name        string
	Label       string
	Weight      int
	OutlineHint string
	DraftHint   string
}

var narrativeArcs = []NarrativeArc{
	{
		Name:   "slow-build",
		Label:  "The Slow Build",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Slow Build"

Structure this essay as a gradual revelation. Open with a quiet everyday scene — something the reader has
experienced a thousand times without thinking about it. Let curiosity build naturally. Introduce the
mathematics one layer at a time, each layer deepening the previous. The "aha" moment arrives late, as the
accumulated layers suddenly click into a unified picture. End by returning to the opening scene, now
transformed.

Key structural beats:
1. A calm, specific everyday observation (not a question — a scene)
2. "Here's something odd about that..." — a wrinkle the reader didn't notice
3. First mathematical layer — accessible, concrete
4. Second layer — deeper, connects to something unexpected
5. The reveal — everything clicks together
6. Return to the opening, seen with new eyes`,
		DraftHint: `This essay uses "The Slow Build" arc. Let the reader settle into the everyday scene before
introducing any mathematics. Build layers gradually — each section should deepen the previous one,
not jump to a new topic. The mathematical reveal should feel earned, not announced. Resist the urge
to tip your hand early. The reader should feel the pieces clicking together, not be told they're about to.`,
	},
	{
		Name:   "cold-open-mystery",
		Label:  "The Cold Open Mystery",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Cold Open Mystery"

Structure this essay as a detective story. Open cold — drop a baffling fact, an unsolved puzzle, or a
question that nobody can satisfactorily answer. No preamble, no everyday scene. Just the mystery.
Then investigate: examine the obvious suspects, show why each fails, eliminate wrong answers one by one.
Build tension as the remaining possibilities narrow. The mathematical explanation is the solution to
the case — the detective's final reveal.

Key structural beats:
1. The mystery — stated bluntly and provocatively in the first paragraph
2. The obvious suspect — the explanation "everyone knows" (and why it fails)
3. Following the evidence — what experiments or observations reveal
4. A false lead or red herring — something that almost works but doesn't
5. The real answer — the mathematical mechanism, presented as discovery
6. Case closed — but leave a thread dangling (what's still unknown)`,
		DraftHint: `This essay uses "The Cold Open Mystery" arc. Open with the mystery immediately — no warm-up,
no scene-setting. The reader should feel the puzzle in the first paragraph. Write like a detective
following leads: examine explanations, eliminate them with evidence, build toward the real answer.
The mathematics should feel like cracking the case, not a textbook aside. Maintain tension throughout.`,
	},
	{
		Name:   "paradox",
		Label:  "The Paradox",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Paradox"

Structure this essay around a central paradox or impossibility. Open by stating something that sounds
flatly impossible — a result that contradicts common sense or intuition. Let the reader sit with the
discomfort. Don't rush to resolve it. Instead, build the mathematical machinery piece by piece,
showing how each tool gets you closer to understanding why the impossible is true. The resolution
should transform the reader's intuition, not just satisfy their curiosity.

Key structural beats:
1. The impossibility — stated plainly, with no softening ("This can't be true, but it is")
2. Why it feels wrong — articulate the reader's intuition and take it seriously
3. The first tool — a mathematical concept that starts to crack the paradox
4. Going deeper — a second concept that shifts the ground further
5. The resolution — the moment intuition updates; the reader sees why it MUST be true
6. The new normal — show how accepting this changes how you see other things`,
		DraftHint: `This essay uses "The Paradox" arc. State the impossible-sounding result early and bluntly.
Respect the reader's disbelief — don't dismiss it, engage with it. Build mathematical tools one
at a time, each one making the paradox slightly less impossible. The resolution should feel like
an intellectual level-up, not just "oh, I see." The reader's intuition should be permanently updated.`,
	},
	{
		Name:   "demolition",
		Label:  "The Demolition",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Demolition"

Structure this essay as a controlled demolition of a widely-believed explanation. Open by confidently
presenting the "textbook answer" — the explanation everyone learned and nobody questions. Make it
sound plausible. Then systematically dismantle it with evidence, calculations, and experiments.
Each section should remove another support beam until the whole structure collapses. Then build the
correct explanation from the rubble — more interesting than what it replaced.

Key structural beats:
1. The confident myth — present the standard explanation as if it's obviously true
2. The first crack — a calculation or observation that doesn't quite fit
3. More cracks — additional evidence piling up against the myth
4. Collapse — the moment it's clear the standard explanation cannot be right
5. From the rubble — build the real explanation, which is more surprising and elegant
6. Why the myth persisted — a brief, human note on why wrong ideas are sticky`,
		DraftHint: `This essay uses "The Demolition" arc. Start by selling the wrong explanation convincingly —
the reader should nod along. Then undermine it methodically. Each piece of counter-evidence should
feel like removing a Jenga block. When the old explanation collapses, build the new one from scratch.
The correct answer should feel more satisfying than the myth it replaced.`,
	},
	{
		Name:   "dual-timeline",
		Label:  "The Dual Timeline",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Dual Timeline"

Structure this essay by interleaving two threads that converge. Thread A is the historical discovery
story — the humans, their arguments, their breakthroughs, their mistakes. Thread B is the modern
mathematical understanding — clean, elegant, powerful. Alternate between threads, section by section.
The historical thread gives narrative momentum and emotional stakes. The mathematical thread gives
clarity and depth. They converge at the end when the reader sees how the messy human journey produced
the beautiful mathematics.

Key structural beats:
1. Thread A: The problem as it first appeared (historical context, the humans involved)
2. Thread B: What we now know — the clean mathematical statement
3. Thread A: Early attempts, wrong turns, arguments between researchers
4. Thread B: The key insight — why this particular approach works
5. Thread A: The breakthrough moment — who cracked it and how
6. Convergence: The human story and the mathematics meet; show the reader both at once`,
		DraftHint: `This essay uses "The Dual Timeline" arc. Alternate between historical narrative and modern
mathematical understanding. Use section breaks or clear transitions to signal shifts between threads.
The historical sections should read like storytelling — names, dates, motivations, conflicts.
The mathematical sections should be crisp and illuminating. The two threads should feel like they're
racing toward each other and finally meeting at the end.`,
	},
	{
		Name:   "zoom",
		Label:  "The Zoom",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Zoom"

Structure this essay as a journey through scales — either zooming in (from cosmos to atom) or zooming
out (from a single grain of sand to the universe). Each scale level reveals new mathematics at work.
The essay is a continuous camera move, pausing at each level to explore what's happening there.
The unifying insight is that the same mathematical principle operates at every scale, but looks
different at each one.

Key structural beats:
1. Start at one extreme — the smallest or largest scale, something concrete and vivid
2. First zoom — move one level, show what changes and what stays the same
3. Second zoom — another level, the pattern deepening
4. Third zoom — by now the reader sees the pattern and starts anticipating
5. The furthest zoom — the most surprising scale where this math still applies
6. Pull back — show all the scales at once; the reader sees the unity`,
		DraftHint: `This essay uses "The Zoom" arc. Structure the essay as a continuous journey through scales.
Each section should feel like moving the camera — zooming in or out to a new level. At each scale,
show something concrete and specific, then reveal the mathematics at work. The reader should feel
the thrill of the same principle appearing again and again at wildly different scales. Use vivid,
sensory details at each level to keep the journey grounded.`,
	},
	{
		Name:   "argument",
		Label:  "The Argument",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Argument"

Structure this essay as a genuine intellectual debate. Present two (or more) competing explanations
for the same phenomenon. Give each side its best case — don't set up a straw man. Let the reader
weigh the evidence like a juror. The essay should feel balanced and suspenseful until the resolution.
When the verdict comes, it should feel earned, not imposed. If the debate is genuinely unresolved,
say so — honest uncertainty is more compelling than false resolution.

Key structural beats:
1. The phenomenon — what needs explaining, presented as a puzzle
2. Team A's case — the first explanation, presented persuasively with evidence
3. Team B's case — the competing explanation, equally persuasive
4. The evidence that breaks the tie — what distinguishes the explanations
5. The verdict — which explanation wins (or why neither fully does)
6. What the debate reveals — what we learn from the argument itself, beyond who won`,
		DraftHint: `This essay uses "The Argument" arc. Present competing explanations as a genuine debate.
Be fair to both sides — the reader should find each plausible before seeing the distinguishing
evidence. Write the competing sections with equal conviction. The resolution should feel like
a jury reaching a verdict after careful deliberation, not a teacher revealing the "right answer."`,
	},
	{
		Name:   "catalog-of-wonders",
		Label:  "The Catalog of Wonders",
		Weight: 2,
		OutlineHint: `NARRATIVE ARC: "The Catalog of Wonders"

Structure this essay as an escalating series of variations on a single mathematical theme. No single
narrative through-line — instead, present 4-5 instances where the same hidden principle appears,
each more surprising than the last. The essay should feel like a magician pulling increasingly
impossible objects from the same hat. The unifying mathematics is the thread connecting them all.
End with the most jaw-dropping example.

Key structural beats:
1. The first instance — familiar, accessible, the reader nods along
2. The second instance — different domain, same math; reader's eyebrows rise
3. The principle — name and explain the underlying mathematics
4. The third instance — now the reader starts looking for it themselves
5. The fourth instance — genuinely shocking; the reader didn't see this one coming
6. The pattern — step back and show why this principle appears everywhere`,
		DraftHint: `This essay uses "The Catalog of Wonders" arc. Present each instance as its own mini-story,
but keep them punchy — the power is in the accumulation, not the depth of any single example.
Each instance should be more surprising than the last. The mathematical explanation threads between
them, growing clearer with each example. The reader should feel the escalation — by the end,
they should be actively looking for the next occurrence themselves.`,
	},
	{
		Name:   "letter-to-a-friend",
		Label:  "The Letter to a Friend",
		Weight: 1,
		OutlineHint: `NARRATIVE ARC: "The Letter to a Friend"

Structure this essay as if writing to a specific curious person who asked you a question. Open with
"You asked me once why..." or a similar direct address. The entire essay is a personal explanation —
intimate, patient, building from what the friend already knows. This creates natural permission to
explain from scratch while feeling warm rather than pedagogical. The mathematics arrives as something
you're sharing, not teaching.

Key structural beats:
1. The question — what the friend asked, in their words
2. "Here's the thing..." — your first attempt to explain, starting simple
3. Going deeper — you realize the simple answer isn't enough, and dig into the math
4. The connection — show why this matters beyond the original question
5. The surprise — something you discovered while preparing this answer
6. The sign-off — return to the friend, leaving them with a new way to see it`,
		DraftHint: `This essay uses "The Letter to a Friend" arc. Write as if addressing a specific curious person.
Use "you" naturally and often. The tone should be intimate and generous — you're sharing something
you find fascinating, not proving you're smart. Let the explanation build patiently from what the
friend already knows. The mathematics should feel like a gift you're unwrapping together.`,
	},
	{
		Name:   "time-lapse",
		Label:  "The Time-Lapse",
		Weight: 1,
		OutlineHint: `NARRATIVE ARC: "The Time-Lapse"

Structure this essay as a compressed journey through time. Start at T=0 — the Big Bang, the first
cell division, the moment a raindrop hits a puddle, the instant a string is plucked — and fast-forward
through time, pausing at key moments when the mathematics becomes visible. Each "pause" is a section.
The essay should feel like a nature documentary's time-lapse: vast spans compressed into vivid
snapshots, with the mathematics as the narrator explaining what's happening at each stop.

Key structural beats:
1. T=0 — the initiating event, described vividly and specifically
2. First pause — seconds/minutes later, the first mathematical pattern emerges
3. Second pause — hours/days/years later, the pattern has evolved or scaled
4. Third pause — a much longer timescale, the same mathematics in a new form
5. The long view — the final timescale, showing where this all leads
6. Now — return to the present moment; the reader sees the whole arc at once`,
		DraftHint: `This essay uses "The Time-Lapse" arc. Write as a journey through time, with each section
anchored to a specific moment. Use concrete timestamps or time references to ground each pause.
The reader should feel time accelerating between sections. At each stop, show something vivid
and specific, then reveal the mathematics at work. The cumulative effect should be awe at how
much unfolds from a single starting moment.`,
	},
	{
		Name:   "wrong-turn",
		Label:  "The Wrong Turn",
		Weight: 1,
		OutlineHint: `NARRATIVE ARC: "The Wrong Turn"

Structure this essay so the reader follows a plausible but incorrect line of reasoning — and discovers
the error themselves. Unlike "The Demolition" (which debunks a myth others believe), this arc has
the READER doing the reasoning, building confidence in a wrong answer, and then hitting a wall.
The essay then backtracks and follows the correct path, which the reader now appreciates far more
because they've felt the failure firsthand. The wrong turn must be genuinely tempting, not a straw man.

Key structural beats:
1. The setup — present a question that invites a natural first approach
2. Following the wrong path — walk through the reasoning step by step, building confidence
3. The wall — show where this approach breaks, with a specific counterexample or contradiction
4. "Where did we go wrong?" — identify the exact moment the reasoning went astray
5. The right path — now follow the correct approach, which the reader appreciates deeply
6. The lesson — what this teaches about mathematical thinking itself`,
		DraftHint: `This essay uses "The Wrong Turn" arc. Lead the reader down the wrong path deliberately —
make it feel natural and logical. Build their confidence before pulling the rug out. The moment
of failure should be a genuine "wait, that can't be right" realization, not a lecture about what's
wrong. When you backtrack to the correct approach, the reader should feel relief and deeper
understanding precisely because they experienced the failure themselves.`,
	},
	{
		Name:   "expedition",
		Label:  "The Expedition",
		Weight: 1,
		OutlineHint: `NARRATIVE ARC: "The Expedition"

Structure this essay as a physical journey through a real place. Walk through a forest, ride a train,
cross a bridge, wander a city. The mathematics emerges from what you encounter along the way — each
location reveals a new mathematical layer. The reader is literally moving through the essay. Sensory
details (sounds, textures, light) ground the mathematics in physical experience. The journey has a
destination, and arriving there brings the mathematical theme into focus.

Key structural beats:
1. Setting out — describe where you are and what you see; the journey begins
2. First encounter — something you pass reveals the first mathematical principle
3. Deeper in — further along, a second encounter deepens the theme
4. The unexpected — something surprising along the way connects to a different branch of math
5. Arrival — reaching the destination, where the full mathematical picture comes together
6. Looking back — from the destination, the whole journey makes new sense`,
		DraftHint: `This essay uses "The Expedition" arc. Write as a journey through a real physical place.
Use sensory details — what you see, hear, feel — to ground every mathematical idea in a location.
The reader should feel like they're walking alongside you. Each stop along the way should reveal
mathematics naturally, not as a digression but as something the place itself is showing you.
The journey should have momentum — the reader should want to see what's around the next corner.`,
	},
}

func RandomArc() NarrativeArc {
	totalWeight := 0
	for _, arc := range narrativeArcs {
		totalWeight += arc.Weight
	}
	pick := randv2.IntN(totalWeight)
	for _, arc := range narrativeArcs {
		pick -= arc.Weight
		if pick < 0 {
			return arc
		}
	}
	return narrativeArcs[0]
}

func ArcByName(name string) (NarrativeArc, bool) {
	for _, arc := range narrativeArcs {
		if arc.Name == name {
			return arc, true
		}
	}
	return NarrativeArc{}, false
}
