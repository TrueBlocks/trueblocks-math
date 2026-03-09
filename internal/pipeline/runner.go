package pipeline

import (
	"context"
	"fmt"
	"log"
	"math"
	randv2 "math/rand/v2"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Runner struct {
	Config   *Config
	Projects []*PipelineState
	Client   *AnthropicClient
	Log      *log.Logger
	BaseDir  string
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

func (r *Runner) RunCycle(ctx context.Context) ([]string, error) {
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

	selected := ps.SelectForCycle(&r.Config.Pipeline)
	if len(selected) == 0 {
		return nil, nil
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

func (r *Runner) processEssay(ctx context.Context, ps *PipelineState, essay *EssayState, targetStage Stage) error {
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
		prompt = OutlinePrompt(ideaMeta.Title, research, targetWords)

	case StageDraft:
		model = r.Config.Models.Draft
		outline := readContent(StageOutline)
		research := readContent(StageResearch)
		prompt = DraftPrompt(ideaMeta.Title, outline, research, targetWords)

	case StageFactcheck:
		model = r.Config.Models.Factcheck
		draft := readContent(StageDraft)
		research := readContent(StageResearch)
		prompt = FactcheckPrompt(ideaMeta.Title, draft, research)

	case StageDraft2:
		model = r.Config.Models.Draft2
		draft := readContent(StageDraft)
		factcheck := readContent(StageFactcheck)
		prompt = Draft2Prompt(ideaMeta.Title, draft, factcheck, targetWords)

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
