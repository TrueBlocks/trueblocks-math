package pipeline

import (
	"context"
	"fmt"
	"log"
	"math"
	randv2 "math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	Config     *Config
	Projects   []*PipelineState
	Client     *AnthropicClient
	Log        *log.Logger
	BaseDir    string
	ConfigPath string
	CLIDryRun  bool
}

func NewRunner(cfg *Config, baseDir string) *Runner {
	return &Runner{
		Config:  cfg,
		Client:  &AnthropicClient{APIKey: cfg.API.AnthropicKey},
		Log:     log.New(os.Stdout, "", 0),
		BaseDir: baseDir,
	}
}

func (r *Runner) DiscoverProjects() error {
	projectsDir := filepath.Join(r.BaseDir, "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return fmt.Errorf("reading projects dir: %w", err)
	}

	r.Projects = nil
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		ideasDir := filepath.Join(projectsDir, entry.Name(), "ideas")
		if _, err := os.Stat(ideasDir); os.IsNotExist(err) {
			continue
		}
		ps := NewPipelineState(entry.Name(), filepath.Join(projectsDir, entry.Name()))
		r.Projects = append(r.Projects, ps)
	}

	if len(r.Projects) == 0 {
		return fmt.Errorf("no projects found in %s", projectsDir)
	}

	return nil
}

func (r *Runner) LoadState() error {
	for _, ps := range r.Projects {
		repairs := ps.RepairOrphans()
		for _, msg := range repairs {
			r.Log.Printf("[%s] REPAIR: %s", ps.Project, msg)
		}
		if err := ps.LoadFromDisk(); err != nil {
			return fmt.Errorf("project %s: %w", ps.Project, err)
		}
	}
	return nil
}

func (r *Runner) ReloadConfig() {
	if r.ConfigPath == "" {
		return
	}
	updated, err := LoadConfig(r.ConfigPath)
	if err != nil {
		r.Log.Printf("Config reload failed: %v (keeping previous settings)", err)
		return
	}
	port := r.Config.Dashboard.Port
	wasCliDry := r.CLIDryRun
	*r.Config = *updated
	r.Config.Dashboard.Port = port
	if wasCliDry {
		r.Config.Pipeline.DryRun = true
	}
	r.Client.APIKey = r.Config.API.AnthropicKey
}

func (r *Runner) RunCycle(ctx context.Context) ([]string, error) {
	r.ReloadConfig()
	var allActions []string
	for _, ps := range r.Projects {
		if ctx.Err() != nil {
			return allActions, ctx.Err()
		}
		actions, err := r.runProjectCycle(ctx, ps)
		if err != nil {
			return allActions, fmt.Errorf("project %s: %w", ps.Project, err)
		}
		allActions = append(allActions, actions...)
	}
	return allActions, nil
}

func (r *Runner) runProjectCycle(ctx context.Context, ps *PipelineState) ([]string, error) {
	if err := ps.LoadFromDisk(); err != nil {
		return nil, fmt.Errorf("loading state: %w", err)
	}

	r.renderAllImages(ps)

	selected := ps.SelectForCycle(&r.Config.Pipeline)
	if len(selected) == 0 {
		return nil, nil
	}

	if r.Config.Pipeline.Debug != "" {
		return r.runDebugCycle(ctx, ps, selected[0])
	}

	ps.CycleCount++
	r.Log.Printf("[%s] Cycle %d: processing %d essays", ps.Project, ps.CycleCount, len(selected))

	concurrency := r.Config.Pipeline.Concurrency
	if concurrency <= 0 {
		concurrency = 3
	}

	type result struct {
		action string
		err    error
	}

	results := make([]result, len(selected))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, essay := range selected {
		nextStage := essay.NextAction()
		if nextStage == StageDone {
			continue
		}

		action := fmt.Sprintf("[%s] %s → %s", ps.Project, essay.Slug, nextStage)
		r.Log.Printf("  %s → %s", essay.Slug, nextStage)

		wg.Add(1)
		go func(idx int, e *EssayState, ns Stage, act string) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[idx] = result{action: fmt.Sprintf("%s (cancelled)", act)}
				return
			}
			defer func() { <-sem }()

			if err := r.processEssay(ctx, ps, e, ns); err != nil {
				r.Log.Printf("  ERROR: %s: %v", e.Slug, err)
				r.markError(ps, e, ns, err)
				results[idx] = result{action: fmt.Sprintf("%s (ERROR: %v)", act, err), err: err}
				return
			}
			results[idx] = result{action: act}
		}(i, essay, nextStage, action)
	}

	wg.Wait()

	var actions []string
	for _, res := range results {
		if res.action != "" {
			actions = append(actions, res.action)
		}
	}

	return actions, nil
}

func (r *Runner) runDebugCycle(ctx context.Context, ps *PipelineState, essay *EssayState) ([]string, error) {
	r.Log.Printf("[%s] DEBUG MODE: %s (%s)", ps.Project, essay.Slug, r.Config.Pipeline.Debug)

	var actions []string
	for {
		if ctx.Err() != nil {
			return actions, ctx.Err()
		}

		nextStage := essay.NextAction()
		if nextStage == StageDone {
			r.Log.Printf("  %s: all stages complete", essay.Slug)
			break
		}

		action := fmt.Sprintf("[%s] %s → %s", ps.Project, essay.Slug, nextStage)
		r.Log.Printf("  %s → %s", essay.Slug, nextStage)

		if err := r.processEssay(ctx, ps, essay, nextStage); err != nil {
			r.Log.Printf("  ERROR: %s: %v", essay.Slug, err)
			r.markError(ps, essay, nextStage, err)
			actions = append(actions, fmt.Sprintf("%s (ERROR: %v)", action, err))
			return actions, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}

func (r *Runner) processEssay(ctx context.Context, ps *PipelineState, essay *EssayState, targetStage Stage) error {
	if targetStage == StageEject {
		return r.ejectEssay(ps, essay)
	}
	if targetStage == StageExport {
		return r.exportEssay(ps, essay)
	}

	prompt, model, err := r.buildPrompt(ps, essay, targetStage)
	if err != nil {
		return fmt.Errorf("building prompt: %w", err)
	}

	r.markInProgress(ps, essay, targetStage, model)

	var result *APIResult
	if r.Config.Pipeline.DryRun {
		result = DryRunResult(targetStage, essay.Title)
		r.Log.Printf("    [dry-run] %s", targetStage)
	} else {
		timeout := time.Duration(r.Config.Pipeline.APITimeout) * time.Second
		if timeout <= 0 {
			timeout = 300 * time.Second
		}
		result, err = r.Client.Call(ctx, model, prompt, timeout)
		if err != nil {
			return err
		}
		r.Log.Printf("    tokens: %d in / %d out, cost: $%.4f",
			result.InputTokens, result.OutputTokens, result.Cost)
	}

	if err := ps.WriteContent(targetStage, essay.Slug, result.Content); err != nil {
		return fmt.Errorf("writing content: %w", err)
	}

	if targetStage == StageIllustrate {
		if err := r.writeImageSources(ps, essay.Slug, result.Content); err != nil {
			r.Log.Printf("    WARNING: writing image sources: %v", err)
		} else if !r.Config.Pipeline.DryRun {
			r.renderImages(ps, essay.Slug)
		}
	}

	r.markComplete(ps, essay, targetStage, model, result)
	ps.TotalCost += result.Cost
	ps.SessionCost += result.Cost

	return nil
}

func (r *Runner) buildPrompt(ps *PipelineState, essay *EssayState, targetStage Stage) (string, string, error) {
	ideaMeta := essay.Meta[StageIdeas]
	if ideaMeta == nil {
		return "", "", fmt.Errorf("no idea metadata for %s", essay.Slug)
	}

	readContent := func(stage Stage) string {
		path := filepath.Join(ps.BaseDir, stage.Dir(), essay.Slug+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			return ""
		}
		return string(data)
	}

	var prompt, model string
	targetWords := r.sampleWordTarget()
	switch targetStage {
	case StageResearch:
		model = r.Config.Models.Research
		ideaContent := readContent(StageIdeas)
		hook := ideaMeta.Title
		hiddenMath := ""
		if len(ideaContent) > 0 {
			hook = ideaContent
		}
		prompt = ResearchPrompt(ideaMeta.Title, hook, hiddenMath)

	case StageOutline:
		model = r.Config.Models.Outline
		research := readContent(StageResearch)
		arc := RandomArc()
		essay.Arc = arc.Name
		prompt = OutlinePrompt(ideaMeta.Title, research, targetWords, arc)

	case StageDraft:
		model = r.Config.Models.Draft
		outline := readContent(StageOutline)
		research := readContent(StageResearch)
		arc, _ := ArcByName(essay.Arc)
		prompt = DraftPrompt(ideaMeta.Title, outline, research, targetWords, arc)

	case StageFactcheck:
		model = r.Config.Models.Factcheck
		draft := readContent(StageDraft)
		research := readContent(StageResearch)
		prompt = FactcheckPrompt(ideaMeta.Title, draft, research)

	case StageDraft2:
		model = r.Config.Models.Draft2
		draft := readContent(StageDraft)
		factcheck := readContent(StageFactcheck)
		illustrate := readContent(StageIllustrate)
		arc, _ := ArcByName(essay.Arc)
		prompt = Draft2Prompt(ideaMeta.Title, draft, factcheck, illustrate, targetWords, arc)

	case StageIllustrate:
		model = r.Config.Models.Illustrate
		draft := readContent(StageDraft)
		factcheck := readContent(StageFactcheck)
		prompt = IllustratePrompt(ideaMeta.Title, draft, factcheck, essay.Slug)

	default:
		return "", "", fmt.Errorf("unknown target stage: %s", targetStage)
	}

	return prompt, model, nil
}

func (r *Runner) sampleWordTarget() int {
	mean := r.Config.Pipeline.ReadMean
	spread := r.Config.Pipeline.ReadSpread
	if mean <= 0 {
		mean = 5.0
	}
	if spread <= 0 {
		spread = 1.0
	}
	minutes := randv2.NormFloat64()*spread + mean
	minutes = math.Max(mean-2*spread, math.Min(mean+2*spread, minutes))
	if minutes < 1 {
		minutes = 1
	}
	return int(minutes * 265)
}

func (r *Runner) markInProgress(ps *PipelineState, essay *EssayState, stage Stage, model string) {
	meta := &EssayMeta{
		Slug:      essay.Slug,
		Title:     essay.Title,
		Book:      essay.Book,
		Part:      essay.Part,
		PartTitle: essay.PartTitle,
		Order:     essay.Order,
		Status:    "in-progress",
		Model:     model,
		Arc:       essay.Arc,
		Created:   essay.Meta[StageIdeas].Created,
		Started:   nowString(),
	}
	ps.WriteMeta(stage, meta)

	ps.mu.Lock()
	essay.Status = "in-progress"
	essay.Meta[stage] = meta
	ps.mu.Unlock()
}

func (r *Runner) markComplete(ps *PipelineState, essay *EssayState, stage Stage, model string, result *APIResult) {
	meta := &EssayMeta{
		Slug:      essay.Slug,
		Title:     essay.Title,
		Book:      essay.Book,
		Part:      essay.Part,
		PartTitle: essay.PartTitle,
		Order:     essay.Order,
		Status:    "final",
		Model:     model,
		Arc:       essay.Arc,
		Created:   essay.Meta[StageIdeas].Created,
		Started:   nowString(),
		Completed: nowString(),
		Tokens:    result.InputTokens + result.OutputTokens,
		Cost:      result.Cost,
	}
	ps.WriteMeta(stage, meta)

	ps.mu.Lock()
	essay.CurrentStage = stage
	essay.Status = "final"
	essay.Meta[stage] = meta
	ps.mu.Unlock()
}

func (r *Runner) markError(ps *PipelineState, essay *EssayState, stage Stage, err error) {
	meta := &EssayMeta{
		Slug:      essay.Slug,
		Title:     essay.Title,
		Book:      essay.Book,
		Part:      essay.Part,
		PartTitle: essay.PartTitle,
		Order:     essay.Order,
		Status:    "error",
		Model:     "",
		Arc:       essay.Arc,
		Created:   essay.Meta[StageIdeas].Created,
		Started:   nowString(),
		Error:     err.Error(),
	}
	ps.WriteMeta(stage, meta)

	ps.mu.Lock()
	essay.Status = "error"
	essay.Meta[stage] = meta
	ps.mu.Unlock()
}

func (r *Runner) TotalCycles() int {
	total := 0
	for _, ps := range r.Projects {
		total += ps.CycleCount
	}
	return total
}

func (r *Runner) TotalCost() float64 {
	total := 0.0
	for _, ps := range r.Projects {
		total += ps.SessionCost
	}
	return total
}

type imageSource struct {
	filename string
	method   string
	source   string
}

func parseImageSources(content string) []imageSource {
	var sources []imageSource
	sepIdx := strings.Index(content, "---IMAGE-SOURCES---")
	if sepIdx < 0 {
		return nil
	}
	remainder := content[sepIdx:]

	re := regexp.MustCompile(`---IMAGE:([^|]+)\|method:(\w+)---\n([\s\S]*?)---END-IMAGE---`)
	matches := re.FindAllStringSubmatch(remainder, -1)
	for _, m := range matches {
		sources = append(sources, imageSource{
			filename: strings.TrimSpace(m[1]),
			method:   strings.TrimSpace(m[2]),
			source:   strings.TrimSpace(m[3]),
		})
	}
	return sources
}

func (r *Runner) writeImageSources(ps *PipelineState, slug, content string) error {
	sources := parseImageSources(content)
	if len(sources) == 0 {
		r.Log.Printf("    no image sources found in illustrate output")
		return nil
	}

	imgDir := filepath.Join(ps.BaseDir, "images", slug)
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return fmt.Errorf("creating images dir: %w", err)
	}

	for _, src := range sources {
		baseName := strings.TrimSuffix(src.filename, filepath.Ext(src.filename))

		var ext string
		switch src.method {
		case "mermaid":
			ext = ".mermaid"
		case "r":
			ext = ".R"
		case "ai":
			ext = ".ai-prompt.txt"
		default:
			r.Log.Printf("    unknown image method: %s for %s", src.method, src.filename)
			continue
		}

		srcPath := filepath.Join(imgDir, baseName+ext)
		if err := os.WriteFile(srcPath, []byte(src.source), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", srcPath, err)
		}

		metaContent := fmt.Sprintf("filename: %s\nmethod: %s\ndescription: %s\n",
			src.filename, src.method, baseName)
		metaPath := filepath.Join(imgDir, baseName+".meta.yaml")
		if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", metaPath, err)
		}

		r.Log.Printf("    image source: %s (%s)", baseName+ext, src.method)
	}

	essayText := content
	if sepIdx := strings.Index(content, "---IMAGE-SOURCES---"); sepIdx > 0 {
		essayText = strings.TrimSpace(content[:sepIdx])
	}
	if err := ps.WriteContent(StageIllustrate, slug, essayText); err != nil {
		return fmt.Errorf("rewriting illustrate content: %w", err)
	}

	return nil
}

func (r *Runner) renderAllImages(ps *PipelineState) {
	dataDir := filepath.Join(r.BaseDir, "data")
	cmd := exec.Command("imagerender", "--data", dataDir, ps.BaseDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		r.Log.Printf("[%s] WARNING: imagerender: %v\n%s", ps.Project, err, string(output))
		return
	}
	lines := strings.TrimSpace(string(output))
	for _, line := range strings.Split(lines, "\n") {
		if strings.Contains(line, "FAILED") || strings.Contains(line, "rendered") {
			r.Log.Printf("[%s] imagerender: %s", ps.Project, line)
		}
	}
}

func (r *Runner) renderImages(ps *PipelineState, slug string) {
	dataDir := filepath.Join(r.BaseDir, "data")
	cmd := exec.Command("imagerender", "--data", dataDir, "--slug", slug, ps.BaseDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		r.Log.Printf("    WARNING: imagerender %s: %v\n%s", slug, err, string(output))
		return
	}
	lines := strings.TrimSpace(string(output))
	if lines != "" {
		for _, line := range strings.Split(lines, "\n") {
			r.Log.Printf("    imagerender: %s", line)
		}
	}
}
