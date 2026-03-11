package pipeline

import "fmt"

func ResearchPrompt(title, hook, hiddenMath string) string {
	return fmt.Sprintf(`You are a research assistant preparing material for an Asimov-style popular mathematics essay.

ESSAY: "%s"

EVERYDAY HOOK: %s

HIDDEN MATHEMATICS: %s

Your task is to produce a thorough research brief that the author will use to write the essay. Include:

1. **Mathematical Foundation** (2-3 paragraphs)
   - State the core mathematical concept precisely
   - Include the key equations or theorems, explained clearly
   - Note any common misconceptions (especially important: flag anything "everyone knows" that is actually wrong)

2. **History and Discovery** (1-2 paragraphs)
   - Who first described this? When? What was the context?
   - Any surprising or delightful historical details?
   - Priority disputes, overlooked contributors, or amusing anecdotes

3. **The Everyday Connection** (1-2 paragraphs)
   - Explain exactly how the everyday experience connects to the math
   - Be specific: what would you measure, observe, or calculate?
   - Include concrete numbers where possible (dimensions, rates, magnitudes)

4. **Worked Examples** (2-3 examples)
   - Provide specific numerical examples the author could use
   - Show the math step by step
   - Choose examples that create "aha" moments

5. **Surprising Extensions** (3-5 bullet points)
   - Connections to other fields the reader wouldn't expect
   - Modern applications or active research
   - Anything that made YOU say "I didn't know that"

6. **Sources and Citations**
   - List key papers, books, or references (author, title, year)
   - Prefer original sources and authoritative textbooks
   - Flag any claims that are contested or uncertain

FORMAT: Use markdown. Be precise but not dry. The author writes in an engaging, conversational style — give them material they can work with. Target 2000-3000 words.`, title, hook, hiddenMath)
}

func OutlinePrompt(title, researchContent string, targetWords int, arc NarrativeArc) string {
	return fmt.Sprintf(`You are structuring an Asimov-style popular mathematics essay.

ESSAY: "%s"

Below is the research brief. Using this material, create a detailed outline for an essay of approximately %d words (~%d minute read).

RESEARCH:
%s

%s

Create an outline that follows the narrative arc above. For each structural beat, provide:
- A section title (informal, curiosity-driven)
- 2-3 sentence summary of what this section covers
- Key facts or examples to include
- The "turn" — the moment of surprise or insight in that section

Also include:

1. **The Big Reveal**
   - Where does the mathematical punchline land?
   - What should the reader feel? (wonder, surprise, satisfaction)

2. **Tone Notes**
   - Conversational, warm, occasionally wry
   - Like explaining something fascinating to a bright friend over coffee
   - NEVER condescending. Assume the reader is intelligent but not specialized.

FORMAT: Use markdown with clear section headers. Be specific — don't say "discuss the math," say WHICH math and HOW.`, title, targetWords, targetWords/265, researchContent, arc.OutlineHint)
}

func DraftPrompt(title, outlineContent, researchContent string, targetWords int, arc NarrativeArc) string {
	arcDirective := ""
	if arc.DraftHint != "" {
		arcDirective = fmt.Sprintf("\nNARRATIVE ARC: %s\n%s\n", arc.Label, arc.DraftHint)
	}
	return fmt.Sprintf(`You are drafting an Asimov-style popular mathematics essay.

ESSAY: "%s"

Below are the outline and research materials. Write a complete first draft.

TARGET LENGTH: Approximately %d words (~%d minute read). This should read like a focused blog post, not an exhaustive treatment. Be concise and engaging.
%s
OUTLINE:
%s

RESEARCH:
%s

WRITING GUIDELINES:

1. **Voice**: Warm, clear, conversational. Think Isaac Asimov meets Mary Roach. You are a knowledgeable friend, not a lecturer.

2. **Structure**: Follow the outline closely but let the prose breathe. Transitions should feel natural, not mechanical. Respect the narrative arc — the outline was designed with a specific storytelling structure in mind.

3. **Mathematics**: Present equations when they illuminate, but always in service of understanding. Every equation must be preceded by intuition and followed by interpretation. If a reader skips the equation, they should still follow the argument.

4. **Examples**: Use concrete, specific examples with real numbers. "A 2-meter chain hanging 30 cm" not "a chain of length L."

5. **Surprise**: Every essay needs at least one moment where the reader thinks "Wait, really?" Build toward it. Don't give it away too early.

6. **Length**: Target approximately %d words (~%d minute read). This is a focused blog post — be selective about which examples and details to include. Every paragraph must earn its place.

7. **Avoid**: Jargon without explanation, false enthusiasm ("Amazing!"), hedging ("It's kind of like..."), and the word "actually."

FORMAT: Markdown with section headers. Use **bold** for emphasis, *italics* for book/paper titles and foreign words. Use $...$ for inline math and $$...$$ for display math.`, title, targetWords, targetWords/265, arcDirective, outlineContent, researchContent, targetWords, targetWords/265)
}

func FactcheckPrompt(title, draftContent, researchContent string) string {
	return fmt.Sprintf(`You are fact-checking a popular mathematics essay before publication.

ESSAY: "%s"

DRAFT:
%s

ORIGINAL RESEARCH:
%s

Check the draft against the research and your own knowledge. For each issue found, report:

1. **Factual Errors**
   - Mathematical statements that are wrong or imprecise
   - Historical claims that are incorrect (wrong dates, wrong attribution)
   - Numbers that don't check out
   - For each: quote the problematic text, explain the error, provide the correction

2. **Misleading Simplifications**
   - Places where simplification crosses into inaccuracy
   - Missing caveats that a knowledgeable reader would notice
   - Claims presented as settled that are actually debated

3. **Missing Attribution**
   - Results attributed to the wrong person
   - Important contributors not mentioned
   - Theorems stated without naming them

4. **Tone Issues**
   - Anything condescending or overly simplified
   - Passages that might bore or lose the reader
   - Moments where the "surprise" falls flat

5. **Verification Summary**
   - Overall assessment: PASS (minor issues only), REVISE (significant issues), or FAIL (fundamental errors)
   - List of specific changes needed, in order of importance

FORMAT: Markdown. Be specific and quote the problematic passages directly. If everything checks out, say so — don't invent problems.`, title, draftContent, researchContent)
}

func IllustratePrompt(title, draftContent, factcheckContent, slug string) string {
	return fmt.Sprintf(`You are an illustration specialist for an Asimov-style popular mathematics essay collection.

ESSAY: "%s"
SLUG: %s

DRAFT:
%s

FACT-CHECK REPORT:
%s

Your task is to INSERT image tags into the essay text and PRODUCE the source code for each image. You are working with the DRAFT (before final revision). The next stage will revise the text to incorporate both the fact-check corrections AND natural references to these images.

## IMAGE REQUIREMENTS

1. **Image 1 (MANDATORY)**: The "hero image" — captures the article's core concept at a glance. This image must work standalone (e.g., as a social media preview). Choose the single most visually compelling aspect of the essay.

2. **Image 2 (only if there is a clear win)**: Only include if there is a CLEAR visual win — a concept that becomes dramatically easier to understand with a diagram, plot, or visual. Ask: "Would a reader genuinely struggle here without a visual?" If the answer is no, stop at one image.

3. **Not every essay needs two images.** One excellent, well-placed image is far better than two mediocre ones.

## IMAGE SIZING (6×9 TRADE PAPERBACK)

These images will appear in a 6×9 inch trade paperback. After margins, the live area is roughly 4.5 inches wide × 7 inches tall. Follow these rules:

- **Landscape or wide images**: max width 4.5 inches, let the height be whatever is natural. Do NOT stretch to fill the page.
- **Portrait or tall images**: max height 4.5 inches, let the width be whatever is natural. A tall, narrow diagram is fine — it should NOT consume the entire page.
- **No single image should dominate a full page.** The reader should always see essay text above and below the image on the same page.
- Mermaid font sizes and node counts should be calibrated so text remains legible at 4.5 inches wide.
- R plots should use readable axis labels and annotations at the printed size.

## IMAGE STYLE AND CREATIVITY

- **Vary the image types.** Do NOT default to flowcharts for everything. Use scatter plots, annotated mathematical diagrams, concept maps, state diagrams, timeline charts, function plots, phase diagrams, comparison tables, or whatever best fits the essay.
- **Vary the placement.** Images do not always go in the same structural position. Place each image exactly where it has the most explanatory impact — sometimes early, sometimes mid-essay, sometimes near the end.
- **Keep Mermaid diagrams simple.** Fewer boxes, fewer arrows, clearer labels. A 4–6 node diagram that communicates one idea crisply is far better than a 12-node diagram that tries to capture everything. If you need more than ~8 nodes, rethink the diagram — split it or simplify.
- **No captions.** Omit the caption portion of the image tag — leave it empty or use a very short alt-text only. Good axis labels, node labels, and clear visual design make captions redundant. Captions compete for visual space and cause clutter.
- **Be creative, not boring.** Surprise the reader. An unexpected visualization choice that illuminates a concept is worth far more than a predictable flowchart.
- **Every image must explain something.** No decorative images. If it does not deepen understanding, cut it.

## IMAGE TAG FORMAT

Insert tags directly into the essay text at the most impactful location:

%s

The tag goes on its own line, exactly where the image should appear in the final document. The caption should be descriptive and meaningful.

## CHOOSING RENDERING METHOD

For each image, choose the BEST rendering method — and mix it up across essays:

**Mermaid** (method: mermaid) — Use for:
- Concept maps showing relationships between ideas
- Simple flowcharts (keep to 4–6 nodes max)
- State diagrams showing transitions
- Sequence diagrams showing step-by-step processes
- Keep diagrams SIMPLE. Fewer nodes with clear labels. If it looks crowded at 4.5 inches wide, you have too many nodes.

**R script** (method: r) — Use for:
- Mathematical function plots (curves, surfaces)
- Statistical or data visualizations
- Annotated mathematical diagrams with coordinates
- Any visualization involving axes or numerical data
- Ensure labels and annotations remain legible at print size

**AI image generation** (method: ai) — Use when a visual metaphor or real-world scene would land better than a code-generated diagram:
- Visual metaphors or analogies (e.g., Cheerios clumping in a bowl to illustrate capillary forces)
- Real-world scenes that illustrate mathematical principles (e.g., a soap bubble cluster showing minimal surfaces)
- Artistic/conceptual illustrations that require photorealism or painterly style
- Physical phenomena that are best shown as realistic renderings
- Do not shy away from AI images — use them whenever they are the most compelling choice

Prefer Mermaid or R when the concept is inherently diagrammatic or quantitative. Use AI when the concept is best captured by a scene, metaphor, or physical illustration. The goal is the clearest, most engaging result — not the most technically reproducible one.

## OUTPUT FORMAT

Your response must contain TWO sections:

### SECTION 1: MODIFIED ESSAY

Output the COMPLETE essay text with image tags inserted at the chosen locations. Preserve ALL original content — do not edit, rewrite, or shorten the essay. Only ADD the image tags.

### SECTION 2: IMAGE SOURCES

After the essay, output a separator line:

%s

Then for EACH image, output:

%s

For Mermaid sources, write valid Mermaid diagram syntax. Use the following theme colors:
- Primary nodes: fill #2563EB (blue), stroke #1D4ED8, color white
- Secondary nodes: fill #7C3AED (purple), stroke #6D28D9, color white
- Accent nodes: fill #059669 (emerald), stroke #047857, color white
- Highlight nodes: fill #DC2626 (red), stroke #B91C1C, color white
- Background: transparent
- Edge labels: #1E293B (dark slate)
- Edges: #64748B (gray)

For R sources, write a complete, self-contained R script that:
- Sources common.R using the environment variable at the top: source(Sys.getenv("COMMON_R"))
- Uses ggplot2 for plotting (and optionally dplyr, scales — already installed)
- Uses the color palette from common.R (math_blue, math_purple, math_emerald, math_red, math_amber, math_sky, math_rose)
- Uses math_theme() for consistent plot styling
- Calls save_chart(p) at the end to save the output — NEVER use ggsave() directly
- Produces a clear, well-labeled visualization
- Is fully self-contained except for the common.R dependency
- NEVER use library(plotly), library(rgl), library(gridExtra), library(patchwork), library(cowplot), or library(magick) — these are NOT installed
- For line geoms (geom_line, geom_path, geom_step, geom_segment, geom_curve, geom_vline, geom_hline, geom_abline, element_line), use linewidth= NOT size= (size= is deprecated in ggplot2 3.4+)
- For annotate(), valid geoms are: "text", "segment", "rect", "point", "label" — NEVER use annotate("arrow", ...). For arrows, use annotate("segment", ..., arrow = arrow(length = unit(0.1, "inches")))
- In element_text(), use face= NOT fontface= (fontface= is for annotate("text") only)

For AI sources, write a detailed image generation prompt that:
- Describes the exact scene, composition, and visual style
- Specifies "clean, educational illustration style with a white background"
- Names the mathematical concept being illustrated
- Describes specific visual elements, their arrangement, and colors
- Mentions the target dimensions: max 4.5" on the longest side, for a 6×9 trade paperback
- Avoids any text or equations in the image (those go in the essay)

## CRITICAL RULES

- Do NOT remove, edit, or rewrite any essay text — only INSERT image tags
- Do NOT add more than 2 images (1 is often enough)
- Every image must EXPLAIN something — no decorative images
- Keep Mermaid diagrams simple: 4–6 nodes, legible at 4.5 inches wide
- No captions — use empty caption in the image tag or minimal alt-text
- Mermaid diagrams must be valid syntax
- R scripts must be syntactically correct and self-contained (except common.R)
- AI prompts must be detailed, specific, and describe the exact scene, style, and mathematical content
- Vary image types and placement across essays — do not be predictable
- Choose image placement for maximum explanatory impact
`, title, slug, draftContent, factcheckContent,
		"[[IMG:img_01_descriptive_name.png|Caption describing what the image shows]]",
		"---IMAGE-SOURCES---",
		"---IMAGE:img_01_descriptive_name.png|method:mermaid---\n(mermaid source code here)\n---END-IMAGE---\n\n---IMAGE:img_02_descriptive_name.png|method:r---\n(R script source code here)\n---END-IMAGE---\n\n---IMAGE:img_03_descriptive_name.png|method:ai---\n(detailed image generation prompt here)\n---END-IMAGE---")
}

func Draft2Prompt(title, draftContent, factcheckContent, illustrateContent string, targetWords int, arc NarrativeArc) string {
	arcDirective := ""
	if arc.DraftHint != "" {
		arcDirective = fmt.Sprintf("\nNARRATIVE ARC: %s\n%s\n", arc.Label, arc.DraftHint)
	}
	return fmt.Sprintf(`You are revising an Asimov-style popular mathematics essay based on a fact-check report and illustration plan.

ESSAY: "%s"

TARGET LENGTH: Approximately %d words (~%d minute read). If the original draft is longer than this target, trim it down — cut less essential examples, tighten prose, remove tangents. Every paragraph must earn its place.
%s
ORIGINAL DRAFT:
%s

FACT-CHECK REPORT:
%s

ILLUSTRATED VERSION (with image tags):
%s

Your task is to produce a REVISED DRAFT that incorporates all corrections from the fact-check report AND naturally integrates the illustrations. Specifically:

1. **Fix all factual errors** — correct mathematical statements, historical claims, and numbers as specified in the report.

2. **Address misleading simplifications** — add necessary caveats or nuance where the report flagged oversimplification.

3. **Fix attribution** — correct any misattributed results, add missing names, and properly cite theorems.

4. **Improve tone** — address any passages flagged as condescending, boring, or that fall flat.

5. **Preserve and reference images** — The illustrated version contains [[IMG:filename.png|caption]] tags. You MUST:
   - Keep ALL image tags exactly as they appear (do not modify the tag syntax)
   - Keep each image tag at the same location or move it to a better location if needed
   - Write prose that naturally references the figures (e.g., "As the diagram shows...", "The plot below reveals...", "Notice in the figure how...")
   - Make the images feel like an integral part of the explanation, not afterthoughts

6. **Preserve voice** — maintain the warm, conversational Asimov-style voice. The corrections should feel seamless, not patched in.

7. **Keep structure** — maintain the same overall structure and flow unless the report specifically calls for reorganization.

8. **Don't over-correct** — if the fact-check says PASS on something, leave it alone. Only change what was flagged.

OUTPUT: The complete revised essay in markdown with image tags preserved. Do NOT include a changelog or notes about what you changed — just the clean, corrected essay ready for human review.

FORMAT: Markdown with section headers. Use **bold** for emphasis, *italics* for book/paper titles and foreign words. Use $...$ for inline math and $$...$$ for display math. Target approximately %d words (~%d minute read).`, title, targetWords, targetWords/265, arcDirective, draftContent, factcheckContent, illustrateContent, targetWords, targetWords/265)
}

func SectionDraft2Prompt(title, ideaContent, partTitle string) string {
	return fmt.Sprintf(`You are writing a section divider page for a popular mathematics book.

SECTION: "%s"
PART TITLE: %s

IDEA FILE:
%s

Write a very short introductory paragraph for this section divider — 2 to 3 sentences maximum.
The paragraph sets up what this part of the book explores, in a warm, inviting voice. It does not
summarize the essays — it invites the reader in. Think of it as a breath between chapters, not an
announcement.

After the paragraph, include this image tag on its own line:

[[IMG:section-placeholder.png|Insert cartoon here]]

REQUIREMENTS:
- Maximum 3 sentences
- Warm, conversational tone — consistent with the book's character
- Do NOT list essays or topics explicitly
- Do NOT use phrases like "In this section..." or "The following chapters..."
- The reader should feel curiosity, not obligation

OUTPUT: A markdown document with:
1. A level-1 heading with the part title
2. The 2-3 sentence introductory paragraph
3. The image tag

Example format:
# What Your Kitchen Knows

Your kitchen is a physics laboratory with terrible safety protocols.
Every morning, your cereal, your coffee, and your shower conspire to
demonstrate graduate-level mathematics — without your permission.

[[IMG:section-placeholder.png|Insert cartoon here]]`, title, partTitle, ideaContent)
}

func IntroOutlinePrompt(title, ideaContent string) string {
	return fmt.Sprintf(`You are outlining a book introduction for a popular mathematics book series.

BOOK INTRODUCTION: "%s"

IDEA FILE:
%s

Create a short outline (4-5 beats) for a 400-500 word book introduction. This is the reader's
first encounter with this book's voice and territory.

REQUIREMENTS:
- Open with a concrete, sensory image (not an abstraction)
- Build toward the book's governing idea without naming it outright
- End with an invitation, not a table of contents
- Do NOT list chapters, parts, or essay titles
- Do NOT use phrases like "In this book..." or "The reader will learn..."

For each beat, provide:
- A one-line description of what happens
- The emotional note it strikes
- How it connects to the next beat

The introduction should feel like the opening of a conversation, not the preface of a textbook.

FORMAT: Markdown outline with clear beat numbers.`, title, ideaContent)
}

func IntroDraftPrompt(title, outlineContent, ideaContent string) string {
	return fmt.Sprintf(`You are writing a book introduction for a popular mathematics book series.

BOOK INTRODUCTION: "%s"

OUTLINE:
%s

IDEA FILE:
%s

Write the complete book introduction from the outline.

REQUIREMENTS:
- 400-500 words maximum (1.5-2 pages)
- Match the book's character as described in the idea file
- No images, no equations, no section headers (except the title)
- Hint at the journey ahead without mapping it
- The reader should feel curiosity, not obligation
- No meta-commentary about the book's structure
- Do NOT reference parts by name or number
- Do NOT use phrases like "Welcome to..." or "In the pages that follow..."

VOICE:
- Warm, clear, conversational — like the opening of a letter from a knowledgeable friend
- Specific and sensory — ground abstract ideas in physical experience
- Confident but not arrogant — share wonder, don't perform expertise

OUTPUT: A markdown document with a level-1 heading and the introduction prose. No section headers within the text.

FORMAT: Markdown. Use *italics* for emphasis where appropriate. No bold, no math notation.`, title, outlineContent, ideaContent)
}

func IntroDraft2Prompt(title, draftContent string) string {
	return fmt.Sprintf(`You are polishing a book introduction for a popular mathematics book series.

BOOK INTRODUCTION: "%s"

DRAFT:
%s

This is a polish pass. Tighten the language, cut anything that sounds like a sales pitch or a
course syllabus. The introduction should read like the opening of a conversation, not the
preface of a textbook.

REQUIREMENTS:
- Keep within 400-500 words
- Preserve the opening sensory image — it sets the tone
- Cut any sentence that tells the reader what to expect rather than making them curious
- Remove any phrase that sounds like marketing or academic boilerplate
- Ensure every sentence earns its place — if it doesn't deepen the invitation, cut it
- Maintain the book's character and voice from the draft

OUTPUT: The complete polished introduction in markdown. No changelog or notes.

FORMAT: Markdown with a level-1 heading. No section headers within the text. Use *italics* for emphasis.`, title, draftContent)
}
