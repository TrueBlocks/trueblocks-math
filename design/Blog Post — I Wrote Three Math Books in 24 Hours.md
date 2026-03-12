# I Wrote a Three-Part Series of Mathematics Books in 24 Hours

**141,000 words. 132 essays. Three books. Illustrated. Fact-checked. Formatted for print. Under $90.**

Let that sink in for a moment. Then let me tell you how.

---

## The Lily Pad Problem

There's a classic brain teaser in mathematics: A lily pad doubles in size every day. On day 30, it covers the entire pond. On what day did it cover half the pond?

Day 29.

This is one of the essays in my book series, actually — *Exponential Growth and Your Intuition*. The thesis is simple: humans are wired to think linearly. We see a pond that's 1% covered on day 23 and think we have ages to go. We are catastrophically wrong.

I bring this up because the publishing industry is standing on day 23 right now, staring at a mostly empty pond, telling itself everything is fine.

---

## What I Actually Did

I built a pipeline. Not a metaphorical one — a real software system, written in Go, that orchestrates the transformation of rough essay ideas into publication-ready Word documents with embedded illustrations. Here's what it does, stage by stage:

**Ideas** — I started with 132 two-sentence essay seeds. A hook and a hidden mathematical punchline. "The pond is half-empty one day before it's full." "Humans linearize growth; exponentials ambush us every time."

**Research** — The pipeline calls Claude's API and generates a 2,000-3,000 word research brief for each essay: mathematical foundations, historical context, worked examples, surprising connections, citations.

**Outline** — Each essay gets structured according to one of 15 narrative arcs I designed. Not random templates — editorial choices. "The Slow Build" opens with a quiet scene and accumulates layers. "The Cold Open Mystery" drops a baffling fact and investigates. "The Confession" starts with what everyone gets wrong and walks the reader through changing their mind. Each arc shapes the outline differently.

**Draft** — A complete first draft, written in an Asimov-meets-Mary-Roach voice. Warm, clear, specific. Real numbers, not variables. Every equation preceded by intuition and followed by interpretation.

**Fact-Check** — The draft gets checked against the research. Factual errors, misleading simplifications, wrong attributions, tone problems. Each issue is quoted and corrected. The essay gets a verdict: PASS, REVISE, or FAIL.

**Illustrate** — The system plans 1-2 images per essay and generates the actual source code: Mermaid diagrams, R plots, or DALL-E prompts. These get rendered to PNG and embedded in the document.

**Final Draft** — All fact-check corrections are incorporated, image references are woven into the prose, and the narrative arc is tightened.

**Export** — Markdown gets converted to Word, images get embedded, styles get applied. Out comes a print-ready .docx formatted for a 6x9 trade paperback.

The whole thing runs concurrently — three essays processing simultaneously, cycling every 15 seconds. I watched it on a web dashboard while drinking coffee.

---

## The Numbers That Should Terrify Publishers

Let me be precise, because precision matters when you're making a point about exponential change:

- **132 items** through the pipeline (94 essays, 24 section dividers, 3 book introductions, plus 11 new essays)
- **141,062 words** of final prose
- **~12.4 million tokens** processed
- **Total API cost: $86.48**
- **Average cost per essay: $0.66**

Sixty-six cents per illustrated, fact-checked, narratively structured essay.

A traditional publisher would spend 18-24 months and six figures to produce what I produced in a day for the cost of a decent dinner.

Now, I should be honest. It wasn't *just* 24 hours. I spent weeks before that designing the pipeline, choosing the essay topics, defining the narrative arcs, writing the prompts, and building the software. The creative and editorial decisions — which topics to include, how to structure each book, what voice to use, what narrative arc fits each essay — those were mine. The 24 hours is the production time. The button-push-to-books time.

But that's exactly the point.

---

## What the Publishing Industry Gets Wrong

The conversation about AI and publishing has been stuck in a stupid loop: "Will AI replace authors?" That's the wrong question. It's like asking in 1995 whether the internet would replace newspapers. The answer was more nuanced and more devastating: it replaced the *business model*.

Here's what individual writers should understand: **your job is safe, and it might actually get better.** The things that make a writer valuable — taste, voice, editorial judgment, the ability to know what's interesting and why — those are exactly the things AI can't do. I chose every one of those 132 essay topics. I designed the 15 narrative arcs. I decided that "The Cheerios Effect" should use "The Slow Build" arc and "The Monty Hall Problem" should use "The Paradox" arc. I wrote the prompts that encode my aesthetic preferences: no false enthusiasm, no hedging, no jargon without explanation, never the word "actually."

The AI did the production work. The grinding, scaling, formatting, fact-checking, illustrating, converting. The work that publishers charge six figures and two years to do.

**That's** what should keep publishing executives up at night. Not that authors will be replaced — that authors won't need publishers.

---

## The Pond Is Filling Faster Than You Think

My pipeline today uses Claude Sonnet — a mid-tier model. Not the most powerful available. The essays are good. Some are genuinely delightful. A few need human editing. The fact-checker catches real errors.

But here's the lily pad math: if the models improve by even 30% per year (they're improving faster), and costs drop by 50% per year (they're dropping faster), then within two years, a system like mine produces essays that need zero human editing at a cost approaching zero.

That's not speculation. That's the curve we're on. And like every exponential, it will feel like nothing is happening right up until the moment everything has happened.

Consider what I have that a traditional publisher doesn't:

- **Speed**: 24 hours vs. 18-24 months
- **Cost**: $86 vs. $100,000+
- **Iteration**: I can regenerate any essay with different parameters in minutes
- **Consistency**: Every essay gets the same editorial framework, the same fact-checking rigor, the same illustration quality
- **Customization**: I can tune the voice, the reading level, the target length, the narrative structure — per essay, in real time, by editing a config file

And I'm one person. With a laptop. And a $86 API bill.

---

## What This Means (and What It Doesn't)

Let me be clear about what I'm not saying.

I'm not saying AI-generated books are as good as the best human-written books. They're not. Mary Roach is a better writer than any language model, and she will be for the foreseeable future.

I'm not saying editorial judgment doesn't matter. It matters more than ever. The difference between my pipeline and someone dumping "write me a book about math" into ChatGPT is the same difference between a film director and someone who owns a camera. The tooling amplifies taste; it doesn't replace it.

I'm not saying individual authors should be worried. If you're a writer with voice, perspective, and editorial instincts, you just got a superpower. You can now produce at a scale that was previously only available to institutions.

What I *am* saying is this: the publishing industry's value proposition — "we turn manuscripts into books and get them to readers" — just got mass-automated. The editorial curation, the production pipeline, the formatting, the illustration, the quality control. A single person with the right software can do all of it, in a day, for pocket change.

The big five publishers are standing at the edge of a pond that's 1% covered in lily pads. They're looking at the open water and seeing decades of runway.

They have seven days.

---

## The Books

The series is called *The Hidden Mathematics Series*. Three volumes:

**Book I: *Everything Is a Rate of Change*** — From the Cheerios floating in your cereal bowl to the GPS satellites correcting for relativity, the derivative chain connects everything. A journey from the intimate to the cosmic.

**Book II: *The Edge of Knowing*** — The frontiers where mathematics breaks down, where probability defies intuition, and where the answers are still unknown. Benford's Law, the Monty Hall Problem, the Birthday Paradox, the Traveling Salesman.

**Book III: *The Hidden Architecture*** — The deep structural patterns that underlie nature, art, music, and thought. Fractals, hyperbolic crochet, the mathematics of braiding, the arrow of time.

132 essays. Each one starts with something you've seen a thousand times — soap bubbles, piano tuning, a shower curtain blowing inward — and reveals the mathematics hiding inside it.

Written by a human. Produced by a machine. In 24 hours. For $86.

The pond is filling.

---

*The pipeline that produced these books is open source and written in Go. The total codebase is ~4,300 lines with a single external dependency. It runs on a laptop.*
