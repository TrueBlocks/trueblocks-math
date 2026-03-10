package pipeline

import (
	"fmt"
	randv2 "math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Stage int

const (
	StageIdeas Stage = iota
	StageResearch
	StageOutline
	StageDraft
	StageFactcheck
	StageIllustrate
	StageDraft2
	StageExport
	StageDone
)

var stageNames = []string{"ideas", "research", "outline", "draft", "factcheck", "illustrate", "draft2", "export"}

func (s Stage) String() string {
	if int(s) < len(stageNames) {
		return stageNames[s]
	}
	return "done"
}

func (s Stage) Dir() string { return s.String() }

func NextStage(s Stage) Stage {
	if s < StageExport {
		return s + 1
	}
	return StageDone
}

func StageFromString(s string) Stage {
	for i, name := range stageNames {
		if name == s {
			return Stage(i)
		}
	}
	return StageIdeas
}

type EssayMeta struct {
	Slug      string  `yaml:"slug"`
	Title     string  `yaml:"title"`
	Type      string  `yaml:"type"`
	Book      string  `yaml:"book"`
	Part      int     `yaml:"part"`
	PartTitle string  `yaml:"part_title"`
	Order     int     `yaml:"order"`
	Status    string  `yaml:"status"`
	Model     string  `yaml:"model"`
	Arc       string  `yaml:"arc,omitempty"`
	Ending    string  `yaml:"ending,omitempty"`
	Created   string  `yaml:"created"`
	Started   string  `yaml:"started,omitempty"`
	Completed string  `yaml:"completed,omitempty"`
	Tokens    int     `yaml:"tokens,omitempty"`
	Cost      float64 `yaml:"cost,omitempty"`
	Error     string  `yaml:"error,omitempty"`
}

type EssayState struct {
	Slug         string
	Title        string
	Type         string
	Book         string
	Part         int
	PartTitle    string
	Order        int
	Arc          string
	Ending       string
	CurrentStage Stage
	Status       string
	Meta         map[Stage]*EssayMeta
}

func (e *EssayState) NextAction() Stage {
	if e.Status == "error" {
		return e.CurrentStage
	}
	if e.CurrentStage == StageIdeas && e.Status == "pending" {
		return StageResearch
	}
	if e.Status == "final" && e.CurrentStage < StageExport {
		return NextStage(e.CurrentStage)
	}
	return StageDone
}

func (e *EssayState) IsDone() bool {
	return e.CurrentStage == StageExport && e.Status == "final"
}

func (e *EssayState) IsAvailable() bool {
	next := e.NextAction()
	return next != StageDone && e.Status != "in-progress"
}

type PipelineState struct {
	Project     string
	BaseDir     string
	Essays      map[string]*EssayState
	CycleCount  int
	TotalCost   float64
	SessionCost float64
	mu          sync.Mutex
}

func NewPipelineState(project, baseDir string) *PipelineState {
	for _, name := range append(stageNames, "images") {
		os.MkdirAll(filepath.Join(baseDir, name), 0755)
	}
	return &PipelineState{
		Project: project,
		BaseDir: baseDir,
		Essays:  make(map[string]*EssayState),
	}
}

func (ps *PipelineState) LoadFromDisk() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.Essays = make(map[string]*EssayState)
	ps.TotalCost = 0

	for stageIdx, stageName := range stageNames {
		stage := Stage(stageIdx)
		dir := filepath.Join(ps.BaseDir, stageName)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("reading %s: %w", dir, err)
		}

		for _, entry := range entries {
			if filepath.Ext(entry.Name()) != ".yaml" {
				continue
			}
			slug := entry.Name()
			slug = slug[:len(slug)-len(".meta.yaml")]

			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				return fmt.Errorf("reading %s/%s: %w", stageName, entry.Name(), err)
			}

			var meta EssayMeta
			if err := yaml.Unmarshal(data, &meta); err != nil {
				return fmt.Errorf("parsing %s/%s: %w", stageName, entry.Name(), err)
			}

			essay, ok := ps.Essays[slug]
			if !ok {
				essay = &EssayState{
					Slug:         meta.Slug,
					Title:        meta.Title,
					Type:         meta.Type,
					Book:         meta.Book,
					Part:         meta.Part,
					PartTitle:    meta.PartTitle,
					Order:        meta.Order,
					Arc:          meta.Arc,
					Ending:       meta.Ending,
					CurrentStage: stage,
					Status:       meta.Status,
					Meta:         make(map[Stage]*EssayMeta),
				}
				ps.Essays[slug] = essay
			}

			essay.Meta[stage] = &meta
			if meta.Type != "" {
				essay.Type = meta.Type
			}
			if meta.Arc != "" {
				essay.Arc = meta.Arc
			}
			if meta.Ending != "" {
				essay.Ending = meta.Ending
			}
			if stage > essay.CurrentStage {
				essay.CurrentStage = stage
				essay.Status = meta.Status
			}
			ps.TotalCost += meta.Cost
		}
	}

	return nil
}

func (ps *PipelineState) WriteMeta(stage Stage, meta *EssayMeta) error {
	dir := filepath.Join(ps.BaseDir, stage.Dir())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}

	metaPath := filepath.Join(dir, meta.Slug+".meta.yaml")
	return os.WriteFile(metaPath, data, 0644)
}

func (ps *PipelineState) WriteContent(stage Stage, slug, content string) error {
	dir := filepath.Join(ps.BaseDir, stage.Dir())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	mdPath := filepath.Join(dir, slug+".md")
	return os.WriteFile(mdPath, []byte(content), 0644)
}

func ParseDebug(debug string) (book string, part, order int, ok bool) {
	parts := strings.Split(debug, ".")
	if len(parts) != 3 {
		return "", 0, 0, false
	}
	book = parts[0]
	part, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, 0, false
	}
	order, err = strconv.Atoi(parts[2])
	if err != nil {
		return "", 0, 0, false
	}
	return book, part, order, true
}

func (ps *PipelineState) SelectForCycle(cfg *PipelineConfig) []*EssayState {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if cfg.Debug != "" {
		book, part, order, ok := ParseDebug(cfg.Debug)
		if !ok {
			return nil
		}
		for _, essay := range ps.Essays {
			if essay.Book == book && essay.Part == part && essay.Order == order {
				return []*EssayState{essay}
			}
		}
		return nil
	}

	maxActions := cfg.MaxPerCycle
	if ps.CycleCount < maxActions {
		maxActions = ps.CycleCount + 1
	}

	var newEssays, continuingEssays []*EssayState
	for _, essay := range ps.Essays {
		if !essay.IsAvailable() {
			continue
		}
		if essay.CurrentStage == StageIdeas {
			newEssays = append(newEssays, essay)
		} else {
			continuingEssays = append(continuingEssays, essay)
		}
	}

	sortByBookOrder(newEssays)
	randv2.Shuffle(len(newEssays), func(i, j int) { newEssays[i], newEssays[j] = newEssays[j], newEssays[i] })
	sortByBookOrder(continuingEssays)

	var selected []*EssayState
	newCount := 0
	for _, e := range continuingEssays {
		if len(selected) >= maxActions {
			break
		}
		selected = append(selected, e)
	}
	for _, e := range newEssays {
		if len(selected) >= maxActions {
			break
		}
		if newCount >= cfg.NewPerCycle {
			break
		}
		selected = append(selected, e)
		newCount++
	}

	return selected
}

func sortByBookOrder(essays []*EssayState) {
	sort.Slice(essays, func(i, j int) bool {
		if essays[i].Book != essays[j].Book {
			return essays[i].Book < essays[j].Book
		}
		if essays[i].Part != essays[j].Part {
			return essays[i].Part < essays[j].Part
		}
		return essays[i].Order < essays[j].Order
	})
}

func (ps *PipelineState) SnapshotEssays() []*EssayState {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	out := make([]*EssayState, 0, len(ps.Essays))
	for _, e := range ps.Essays {
		out = append(out, e)
	}
	return out
}

func (ps *PipelineState) Summary() map[string]int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	counts := map[string]int{
		"pending": 0, "research": 0, "outline": 0,
		"draft": 0, "factcheck": 0, "draft2": 0, "illustrate": 0, "done": 0, "error": 0,
	}
	for _, e := range ps.Essays {
		if e.Status == "error" {
			counts["error"]++
		} else if e.IsDone() {
			counts["done"]++
		} else if e.CurrentStage == StageIdeas {
			counts["pending"]++
		} else {
			counts[e.CurrentStage.String()]++
		}
	}
	return counts
}

func (ps *PipelineState) RepairOrphans() []string {
	var repairs []string
	for stageIdx, stageName := range stageNames {
		stage := Stage(stageIdx)
		if stage == StageIdeas {
			continue
		}
		dir := filepath.Join(ps.BaseDir, stageName)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) != ".yaml" {
				continue
			}
			slug := entry.Name()[:len(entry.Name())-len(".meta.yaml")]
			yamlPath := filepath.Join(dir, entry.Name())
			mdPath := filepath.Join(dir, slug+".md")

			if _, err := os.Stat(mdPath); os.IsNotExist(err) {
				os.Remove(yamlPath)
				repairs = append(repairs, fmt.Sprintf("removed orphan %s/%s (no .md file)", stageName, entry.Name()))
				continue
			}

			data, err := os.ReadFile(yamlPath)
			if err != nil {
				continue
			}
			var meta EssayMeta
			if err := yaml.Unmarshal(data, &meta); err != nil {
				continue
			}
			if meta.Status == "in-progress" {
				meta.Status = "final"
				meta.Completed = nowString()
				if fixed, err := yaml.Marshal(&meta); err == nil {
					os.WriteFile(yamlPath, fixed, 0644)
					repairs = append(repairs, fmt.Sprintf("fixed %s/%s (in-progress → final)", stageName, entry.Name()))
				}
			}
		}
	}
	return repairs
}

func nowString() string {
	return time.Now().Format("2006-01-02T15:04:05")
}
