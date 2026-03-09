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

func OutlinePrompt(title, researchContent string, targetWords int) string {
	return fmt.Sprintf(`You are structuring an Asimov-style popular mathematics essay.

ESSAY: "%s"

Below is the research brief. Using this material, create a detailed outline for an essay of approximately %d words (~%d minute read).

RESEARCH:
%s

Create an outline with:

1. **Opening Hook** (1-2 sentences)
   - Start with the everyday experience, not the math
   - Make the reader curious: "Why does...?" or "You've noticed that..."
   - NO academic throat-clearing. Drop the reader into the scene.

2. **Section Breakdown** (4-6 sections)
   For each section provide:
   - Section title (informal, curiosity-driven)
   - 2-3 sentence summary of what this section covers
   - Key facts or examples to include
   - The "turn" — the moment of surprise or insight

3. **The Big Reveal**
   - Where does the mathematical punchline land?
   - What should the reader feel? (wonder, surprise, satisfaction)

4. **Closing** (1-2 sentences)
   - Circle back to the opening
   - Leave the reader seeing the everyday thing differently

5. **Tone Notes**
   - Conversational, warm, occasionally wry
   - Like explaining something fascinating to a bright friend over coffee
   - NEVER condescending. Assume the reader is intelligent but not specialized.

FORMAT: Use markdown with clear section headers. Be specific — don't say "discuss the math," say WHICH math and HOW.`, title, targetWords, targetWords/265, researchContent)
}

func DraftPrompt(title, outlineContent, researchContent string, targetWords int) string {
	return fmt.Sprintf(`You are drafting an Asimov-style popular mathematics essay.

ESSAY: "%s"

Below are the outline and research materials. Write a complete first draft.

TARGET LENGTH: Approximately %d words (~%d minute read). This should read like a focused blog post, not an exhaustive treatment. Be concise and engaging.

OUTLINE:
%s

RESEARCH:
%s

WRITING GUIDELINES:

1. **Voice**: Warm, clear, conversational. Think Isaac Asimov meets Mary Roach. You are a knowledgeable friend, not a lecturer.

2. **Structure**: Follow the outline closely but let the prose breathe. Transitions should feel natural, not mechanical.

3. **Mathematics**: Present equations when they illuminate, but always in service of understanding. Every equation must be preceded by intuition and followed by interpretation. If a reader skips the equation, they should still follow the argument.

4. **Examples**: Use concrete, specific examples with real numbers. "A 2-meter chain hanging 30 cm" not "a chain of length L."

5. **Surprise**: Every essay needs at least one moment where the reader thinks "Wait, really?" Build toward it. Don't give it away too early.

6. **Opening**: Start in the everyday world. No "Mathematics is..." or "In this essay..." — start with a scene, an observation, a question.

7. **Closing**: Return to where you started. The reader should now see the familiar thing through mathematical eyes.

8. **Length**: Target approximately %d words (~%d minute read). This is a focused blog post — be selective about which examples and details to include. Every paragraph must earn its place.

9. **Avoid**: Jargon without explanation, false enthusiasm ("Amazing!"), hedging ("It's kind of like..."), and the word "actually."

FORMAT: Markdown with section headers. Use **bold** for emphasis, *italics* for book/paper titles and foreign words. Use $...$ for inline math and $$...$$ for display math.`, title, targetWords, targetWords/265, outlineContent, researchContent, targetWords, targetWords/265)
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

func Draft2Prompt(title, draftContent, factcheckContent string, targetWords int) string {
	return fmt.Sprintf(`You are revising an Asimov-style popular mathematics essay based on a fact-check report.

ESSAY: "%s"

TARGET LENGTH: Approximately %d words (~%d minute read). If the original draft is longer than this target, trim it down — cut less essential examples, tighten prose, remove tangents. Every paragraph must earn its place.

ORIGINAL DRAFT:
%s

FACT-CHECK REPORT:
%s

Your task is to produce a REVISED DRAFT that incorporates all corrections and improvements identified in the fact-check report. Specifically:

1. **Fix all factual errors** — correct mathematical statements, historical claims, and numbers as specified in the report.

2. **Address misleading simplifications** — add necessary caveats or nuance where the report flagged oversimplification.

3. **Fix attribution** — correct any misattributed results, add missing names, and properly cite theorems.

4. **Improve tone** — address any passages flagged as condescending, boring, or that fall flat.

5. **Preserve voice** — maintain the warm, conversational Asimov-style voice. The corrections should feel seamless, not patched in.

6. **Keep structure** — maintain the same overall structure and flow unless the report specifically calls for reorganization.

7. **Don't over-correct** — if the fact-check says PASS on something, leave it alone. Only change what was flagged.

OUTPUT: The complete revised essay in markdown. Do NOT include a changelog or notes about what you changed — just the clean, corrected essay ready for human review.

FORMAT: Markdown with section headers. Use **bold** for emphasis, *italics* for book/paper titles and foreign words. Use $...$ for inline math and $$...$$ for display math. Target approximately %d words (~%d minute read).`, title, targetWords, targetWords/265, draftContent, factcheckContent, targetWords, targetWords/265)
}
