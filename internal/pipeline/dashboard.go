package pipeline

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LogEntry struct {
	Time    string `json:"time"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

type LogBuffer struct {
	mu      sync.Mutex
	entries []LogEntry
	maxSize int
	out     io.Writer
	version atomic.Int64
}

func NewLogBuffer(out io.Writer, maxSize int) *LogBuffer {
	return &LogBuffer{out: out, maxSize: maxSize}
}

func (lb *LogBuffer) Write(p []byte) (n int, err error) {
	n, err = lb.out.Write(p)
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	level := "info"
	if containsStr(msg, "ERROR") {
		level = "error"
	} else if containsStr(msg, "VERBOSE") {
		level = "verbose"
	}
	lb.mu.Lock()
	lb.entries = append(lb.entries, LogEntry{
		Time:    time.Now().Format("15:04:05"),
		Message: msg,
		Level:   level,
	})
	if len(lb.entries) > lb.maxSize {
		lb.entries = lb.entries[len(lb.entries)-lb.maxSize:]
	}
	lb.mu.Unlock()
	lb.version.Add(1)
	return
}

func (lb *LogBuffer) Entries(verbose bool) []LogEntry {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if verbose {
		out := make([]LogEntry, len(lb.entries))
		copy(out, lb.entries)
		return out
	}
	var out []LogEntry
	for _, e := range lb.entries {
		if e.Level != "verbose" {
			out = append(out, e)
		}
	}
	return out
}

func (lb *LogBuffer) Version() int64 {
	return lb.version.Load()
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

type Dashboard struct {
	Runner        *Runner
	Port          int
	CycleInterval int
	ConfigPath    string
	stepCh        chan struct{}
	intervalCh    chan int
	mu            sync.Mutex
	lastLog       []string
	nextCycleAt   time.Time
	cycleRunning  bool
	cycleStarted  time.Time
	retryMessage  string
	LogBuf        *LogBuffer
}

func NewDashboard(runner *Runner, port, cycleInterval int, configPath string, logBuf *LogBuffer) *Dashboard {
	return &Dashboard{
		Runner:        runner,
		Port:          port,
		CycleInterval: cycleInterval,
		ConfigPath:    configPath,
		stepCh:        make(chan struct{}, 1),
		intervalCh:    make(chan int, 1),
		LogBuf:        logBuf,
	}
}

func (d *Dashboard) SetNextCycleAt(t time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.nextCycleAt = t
}

func (d *Dashboard) SetCycleRunning() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cycleRunning = true
	d.cycleStarted = time.Now()
}

func (d *Dashboard) SetCycleFinished() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cycleRunning = false
}

func (d *Dashboard) SetRetryMessage(msg string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.retryMessage = msg
}

func (d *Dashboard) StepChannel() <-chan struct{} {
	return d.stepCh
}

func (d *Dashboard) IntervalChannel() <-chan int {
	return d.intervalCh
}

func (d *Dashboard) SetLastLog(actions []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lastLog = actions
}

func (d *Dashboard) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", d.handleIndex)
	mux.HandleFunc("/api/status", d.handleStatus)
	mux.HandleFunc("/api/step", d.handleStep)
	mux.HandleFunc("/api/essays", d.handleEssays)
	mux.HandleFunc("/api/open", d.handleOpen)
	mux.HandleFunc("/api/eject", d.handleEject)
	mux.HandleFunc("/api/open-docx", d.handleOpenDocx)
	mux.HandleFunc("/api/disk-stats", d.handleDiskStats)
	mux.HandleFunc("/api/logs", d.handleLogs)
	mux.HandleFunc("/api/settings", d.handleSettings)

	addr := fmt.Sprintf("127.0.0.1:%d", d.Port)
	d.Runner.Log.Printf("Dashboard: http://%s", addr)
	return http.ListenAndServe(addr, mux)
}

func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.New("dashboard").Parse(dashboardHTML)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.Execute(w, nil)
}

type statusResponse struct {
	DryRun        bool                    `json:"dry_run"`
	Verbose       bool                    `json:"verbose"`
	Cycle         int                     `json:"cycle"`
	TotalCost     float64                 `json:"total_cost"`
	CycleInterval int                     `json:"cycle_interval"`
	SecondsLeft   int                     `json:"seconds_left"`
	CycleRunning  bool                    `json:"cycle_running"`
	CycleElapsed  int                     `json:"cycle_elapsed"`
	RetryMessage  string                  `json:"retry_message"`
	ReadMean      float64                 `json:"read_mean"`
	ReadSpread    float64                 `json:"read_spread"`
	Projects      []projectStatusResponse `json:"projects"`
	Summary       map[string]int          `json:"summary"`
	LastLog       []string                `json:"last_log"`
}

type projectStatusResponse struct {
	Name    string         `json:"name"`
	Cycle   int            `json:"cycle"`
	Cost    float64        `json:"cost"`
	Summary map[string]int `json:"summary"`
}

func (d *Dashboard) handleStatus(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	lastLog := d.lastLog
	nextAt := d.nextCycleAt
	running := d.cycleRunning
	started := d.cycleStarted
	retryMsg := d.retryMessage
	d.mu.Unlock()

	secsLeft := int(time.Until(nextAt).Seconds())
	if secsLeft < 0 {
		secsLeft = 0
	}

	var elapsed int
	if running {
		elapsed = int(time.Since(started).Seconds())
	}

	totalSummary := map[string]int{
		"pending": 0, "research": 0, "outline": 0,
		"draft": 0, "factcheck": 0, "draft2": 0, "illustrate": 0, "done": 0, "error": 0,
	}

	var projectStatuses []projectStatusResponse
	for _, ps := range d.Runner.Projects {
		s := ps.Summary()
		projectStatuses = append(projectStatuses, projectStatusResponse{
			Name:    ps.Project,
			Cycle:   ps.CycleCount,
			Cost:    ps.TotalCost,
			Summary: s,
		})
		for k, v := range s {
			totalSummary[k] += v
		}
	}

	resp := statusResponse{
		DryRun:        d.Runner.Config.Pipeline.DryRun,
		Verbose:       d.Runner.Config.Pipeline.Verbose,
		Cycle:         d.Runner.TotalCycles(),
		TotalCost:     d.Runner.TotalCost(),
		CycleInterval: d.CycleInterval,
		SecondsLeft:   secsLeft,
		CycleRunning:  running,
		CycleElapsed:  elapsed,
		RetryMessage:  retryMsg,
		ReadMean:      d.Runner.Config.Pipeline.ReadMean,
		ReadSpread:    d.Runner.Config.Pipeline.ReadSpread,
		Projects:      projectStatuses,
		Summary:       totalSummary,
		LastLog:       lastLog,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (d *Dashboard) handleStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	select {
	case d.stepCh <- struct{}{}:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"ok":true}`)
	default:
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, `{"ok":false,"error":"step already pending"}`)
	}
}

func (d *Dashboard) handleOpen(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	slug := r.URL.Query().Get("slug")
	if project == "" || slug == "" {
		http.Error(w, "project and slug required", http.StatusBadRequest)
		return
	}

	var ps *PipelineState
	for _, p := range d.Runner.Projects {
		if p.Project == project {
			ps = p
			break
		}
	}
	if ps == nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	essay := ps.Essays[slug]
	if essay == nil {
		http.Error(w, "essay not found", http.StatusNotFound)
		return
	}

	stage := essay.CurrentStage
	mdFile := filepath.Join(ps.BaseDir, stage.Dir(), slug+".md")

	if err := exec.Command("open", mdFile).Run(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ok":true,"path":"%s"}`, mdFile)
}

func (d *Dashboard) handleEject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	project := r.URL.Query().Get("project")
	slug := r.URL.Query().Get("slug")
	if project == "" || slug == "" {
		http.Error(w, "project and slug required", http.StatusBadRequest)
		return
	}

	var ps *PipelineState
	for _, p := range d.Runner.Projects {
		if p.Project == project {
			ps = p
			break
		}
	}
	if ps == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"project not found"}`)
		return
	}

	essay := ps.Essays[slug]
	if essay == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"essay not found"}`)
		return
	}

	mdFile := filepath.Join(ps.BaseDir, "draft2", slug+".md")
	if _, err := os.Stat(mdFile); os.IsNotExist(err) {
		mdFile = filepath.Join(ps.BaseDir, "illustrate", slug+".md")
		if _, err := os.Stat(mdFile); os.IsNotExist(err) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"ok":false,"error":"neither draft2 nor illustrate file found"}`)
			return
		}
	}

	raw, err := os.ReadFile(mdFile)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"reading draft2: %s"}`, err)
		return
	}
	cleaned := sanitizeForDocx(string(raw))

	tmpFile, err := os.CreateTemp("", "eject-*.md")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"creating temp file: %s"}`, err)
		return
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(cleaned); err != nil {
		tmpFile.Close()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"writing temp file: %s"}`, err)
		return
	}
	tmpFile.Close()

	outDir := filepath.Join(ps.BaseDir, "draft3")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"creating draft3 dir: %s"}`, err)
		return
	}

	home, _ := os.UserHomeDir()
	templatePath := filepath.Join(home, ".works", "templates", "book-template.dotm")
	outFile := filepath.Join(outDir, slug+".docx")

	cmd := exec.Command("md2docx", templatePath, tmpFile.Name(), outFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"md2docx: %s %s"}`, err, string(output))
		return
	}

	dataDir := filepath.Join(d.Runner.BaseDir, "data")
	renderCmd := exec.Command("imagerender", "--data", dataDir, "--slug", slug, ps.BaseDir)
	if output, err := renderCmd.CombinedOutput(); err != nil {
		d.Runner.Log.Printf("[%s] WARNING: imagerender %s: %v %s", project, slug, err, string(output))
	} else {
		d.Runner.Log.Printf("[%s] imagerender %s: %s", project, slug, strings.TrimSpace(string(output)))
	}

	swapCmd := exec.Command("imageswap", "--images", filepath.Join(ps.BaseDir, "images"), outFile)
	if output, err := swapCmd.CombinedOutput(); err != nil {
		d.Runner.Log.Printf("[%s] WARNING: imageswap %s: %v %s", project, slug, err, string(output))
	}

	d.Runner.Log.Printf("[%s] Ejected %s → %s", project, slug, outFile)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ok":true,"path":"%s"}`, outFile)
}

func (d *Dashboard) handleOpenDocx(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	slug := r.URL.Query().Get("slug")
	if project == "" || slug == "" {
		http.Error(w, "project and slug required", http.StatusBadRequest)
		return
	}

	var ps *PipelineState
	for _, p := range d.Runner.Projects {
		if p.Project == project {
			ps = p
			break
		}
	}
	if ps == nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	docxFile := filepath.Join(ps.BaseDir, "draft3", slug+".docx")
	if err := exec.Command("open", docxFile).Run(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ok":true,"path":"%s"}`, docxFile)
}

var (
	reHTMLTag    = regexp.MustCompile(`<[^>]+>`)
	reFencedCode = regexp.MustCompile("(?s)```.*?```")
)

func sanitizeForDocx(md string) string {
	md = reFencedCode.ReplaceAllString(md, "[markdown removed]")
	md = reHTMLTag.ReplaceAllString(md, "")

	var lines []string
	for _, line := range strings.Split(md, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			lines = append(lines, "[markdown removed]")
			continue
		}
		if strings.HasPrefix(trimmed, "![") {
			lines = append(lines, "[markdown removed]")
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (d *Dashboard) handleLogs(w http.ResponseWriter, r *http.Request) {
	verbose := r.URL.Query().Get("verbose") == "true"
	ver := d.LogBuf.Version()
	entries := d.LogBuf.Entries(verbose)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Version int64      `json:"version"`
		Entries []LogEntry `json:"entries"`
	}{Version: ver, Entries: entries})
}

func (d *Dashboard) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CycleInterval *int     `json:"cycle_interval"`
		Verbose       *bool    `json:"verbose"`
		ReadMean      *float64 `json:"read_mean"`
		ReadSpread    *float64 `json:"read_spread"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.CycleInterval != nil && *req.CycleInterval >= 5 {
		d.CycleInterval = *req.CycleInterval
		d.Runner.Config.Pipeline.CycleInterval = *req.CycleInterval
		select {
		case d.intervalCh <- *req.CycleInterval:
		default:
		}
	}
	if req.Verbose != nil {
		d.Runner.Config.Pipeline.Verbose = *req.Verbose
	}
	if req.ReadMean != nil && *req.ReadMean >= 1 && *req.ReadMean <= 15 {
		d.Runner.Config.Pipeline.ReadMean = *req.ReadMean
	}
	if req.ReadSpread != nil && *req.ReadSpread >= 0.5 && *req.ReadSpread <= 3 {
		d.Runner.Config.Pipeline.ReadSpread = *req.ReadSpread
	}

	if err := SaveConfig(d.ConfigPath, d.Runner.Config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ok":true}`)
}

func (d *Dashboard) handleDiskStats(w http.ResponseWriter, r *http.Request) {
	stages := []string{"ideas", "research", "outline", "draft", "factcheck", "illustrate", "draft2", "draft3"}
	counts := make(map[string]int, len(stages))

	for _, ps := range d.Runner.Projects {
		for _, stage := range stages {
			dir := filepath.Join(ps.BaseDir, stage)
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range entries {
				ext := filepath.Ext(e.Name())
				if ext == ".md" || ext == ".docx" {
					counts[stage]++
				}
			}
		}
	}

	totalCost := 0.0
	for _, ps := range d.Runner.Projects {
		ps.mu.Lock()
		totalCost += ps.TotalCost + ps.SessionCost
		ps.mu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Counts    map[string]int `json:"counts"`
		TotalCost float64        `json:"total_cost"`
	}{Counts: counts, TotalCost: totalCost})
}

type essayJSON struct {
	Project   string `json:"project"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Book      string `json:"book"`
	Part      int    `json:"part"`
	PartTitle string `json:"part_title"`
	Order     int    `json:"order"`
	Stage     string `json:"stage"`
	Status    string `json:"status"`
	Arc       string `json:"arc,omitempty"`
	Error     string `json:"error,omitempty"`
	Ejected   string `json:"ejected"`
	WordCount int    `json:"word_count,omitempty"`
	ReadMins  int    `json:"read_mins,omitempty"`
}

func (d *Dashboard) handleEssays(w http.ResponseWriter, r *http.Request) {
	projectFilter := r.URL.Query().Get("project")

	var essays []essayJSON
	for _, ps := range d.Runner.Projects {
		if projectFilter != "" && ps.Project != projectFilter {
			continue
		}
		for _, e := range ps.SnapshotEssays() {
			stage := e.CurrentStage.String()
			status := e.Status
			if e.IsDone() {
				stage = "done"
				status = "done"
			}
			errMsg := ""
			if meta, ok := e.Meta[e.CurrentStage]; ok && meta.Error != "" {
				errMsg = meta.Error
			}
			ejected := ""
			var wordCount, readMins int
			if e.IsDone() {
				docxPath := filepath.Join(ps.BaseDir, "draft3", e.Slug+".docx")
				md2Path := filepath.Join(ps.BaseDir, "draft2", e.Slug+".md")
				illPath := filepath.Join(ps.BaseDir, "illustrate", e.Slug+".md")
				if di, err := os.Stat(docxPath); err == nil {
					ejected = "current"
					if mi, err := os.Stat(md2Path); err == nil && mi.ModTime().After(di.ModTime()) {
						ejected = "stale"
					} else if mi, err := os.Stat(illPath); err == nil && mi.ModTime().After(di.ModTime()) {
						ejected = "stale"
					}
				}
			}
			for _, s := range []string{"draft2", "illustrate", "draft", "outline"} {
				path := filepath.Join(ps.BaseDir, s, e.Slug+".md")
				if raw, err := os.ReadFile(path); err == nil {
					words := strings.Fields(string(raw))
					wordCount = len(words)
					readMins = (wordCount + 264) / 265
					break
				}
			}
			essays = append(essays, essayJSON{
				Project:   ps.Project,
				Slug:      e.Slug,
				Title:     e.Title,
				Book:      e.Book,
				Part:      e.Part,
				PartTitle: e.PartTitle,
				Order:     e.Order,
				Stage:     stage,
				Status:    status,
				Arc:       e.Arc,
				Error:     errMsg,
				Ejected:   ejected,
				WordCount: wordCount,
				ReadMins:  readMins,
			})
		}
	}

	sort.Slice(essays, func(i, j int) bool {
		ri, rj := stageRank(essays[i].Stage, essays[i].Ejected), stageRank(essays[j].Stage, essays[j].Ejected)
		if ri != rj {
			return ri < rj
		}
		if essays[i].Book != essays[j].Book {
			return essays[i].Book < essays[j].Book
		}
		if essays[i].Part != essays[j].Part {
			return essays[i].Part < essays[j].Part
		}
		return essays[i].Order < essays[j].Order
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

func stageRank(stage, ejected string) int {
	switch stage {
	case "research":
		return 1
	case "outline":
		return 2
	case "draft":
		return 3
	case "factcheck":
		return 4
	case "draft2":
		return 5
	case "illustrate":
		return 6
	case "done":
		if ejected == "current" {
			return 9
		}
		return 7
	default:
		return 8
	}
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Pipeline Dashboard</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
         background: #1a1a2e; color: #e0e0e0; }
  .layout { display: flex; height: 100vh; }
  .main-panel { flex: 1; min-width: 0; padding: 20px; overflow-y: auto; }
  .log-panel { width: 380px; min-width: 300px; background: #0f0f23; border-left: 1px solid #2a2a4e;
               display: flex; flex-direction: column; }
  .log-header { padding: 12px 16px; border-bottom: 1px solid #2a2a4e; display: flex;
                align-items: center; gap: 12px; flex-wrap: wrap; }
  .log-header h2 { color: #e94560; font-size: 1rem; margin: 0; white-space: nowrap; }
  .log-settings { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
  .log-settings label { font-size: 0.75rem; color: #8899aa; display: flex; align-items: center; gap: 4px; cursor: pointer; }
  .log-settings input[type="checkbox"] { accent-color: #e94560; }
  .interval-ctrl { display: flex; align-items: center; gap: 4px; }
  .interval-ctrl input[type="number"] { width: 60px; background: #16213e; border: 1px solid #2a2a4e;
    color: #e0e0e0; border-radius: 4px; padding: 3px 6px; font-size: 0.8rem; text-align: center; }
  .interval-ctrl .save-btn { background: #e94560; color: white; border: none; padding: 3px 10px;
    border-radius: 4px; cursor: pointer; font-size: 0.75rem; }
  .interval-ctrl .save-btn:hover { background: #c73650; }
  .interval-ctrl .save-btn.saved { background: #4caf50; }
  .log-body { flex: 1; overflow-y: auto; overflow-x: auto; padding: 8px 12px;
              font-family: "SF Mono", Monaco, monospace; font-size: 0.75rem; line-height: 1.5; }
  .log-line { padding: 1px 0; white-space: pre; }
  .log-line.error { color: #f44336; }
  .log-line.verbose { color: #546e7a; }
  .log-line.info { color: #6a9955; }
  .log-count { font-size: 0.7rem; color: #546e7a; padding: 4px 12px; border-top: 1px solid #2a2a4e;
               display: flex; justify-content: space-between; align-items: center; }
  .copy-btn { background: #16213e; color: #8899aa; border: 1px solid #2a2a4e; padding: 2px 10px;
              border-radius: 4px; cursor: pointer; font-size: 0.7rem; }
  .copy-btn:hover { border-color: #e94560; color: #e0e0e0; }
  .copy-btn.copied { background: #4caf50; color: white; border-color: #4caf50; }
  h1 { color: #e94560; margin-bottom: 4px; font-size: 1.5rem; }
  .subtitle { color: #8899aa; font-size: 0.85rem; margin-bottom: 16px; }
  .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
          gap: 10px; margin-bottom: 20px; }
  .card { background: #16213e; border-radius: 8px; padding: 12px; text-align: center; }
  .card .number { font-size: 1.6rem; font-weight: bold; color: #e94560; }
  .card .label { font-size: 0.7rem; color: #8899aa; text-transform: uppercase; letter-spacing: 1px; }
  .controls { margin-bottom: 20px; display: flex; align-items: center; gap: 16px; }
  .btn { background: #e94560; color: white; border: none; padding: 10px 24px;
         border-radius: 6px; cursor: pointer; font-size: 1rem; }
  .btn:hover { background: #c73650; }
  .btn:disabled { background: #444; cursor: not-allowed; }
  .mode { display: inline-block; padding: 4px 12px; border-radius: 4px;
          font-size: 0.85rem; font-weight: bold; }
  .mode.dry { background: #ff9800; color: #000; }
  .mode.live { background: #4caf50; color: #fff; }
  .countdown { font-family: "SF Mono", Monaco, monospace; font-size: 1.4rem;
               color: #e94560; min-width: 60px; }
  .countdown-label { color: #8899aa; font-size: 0.75rem; text-transform: uppercase; }
  .progress-ring { width: 48px; height: 48px; }
  .progress-ring circle { fill: none; stroke-width: 3; }
  .progress-ring .bg { stroke: #16213e; }
  .progress-ring .fg { stroke: #e94560; stroke-linecap: round;
                       transform: rotate(-90deg); transform-origin: center;
                       transition: stroke-dashoffset 0.5s ease; }
  .project-tabs { display: flex; gap: 8px; margin-bottom: 16px; flex-wrap: wrap; }
  .tab { padding: 6px 16px; border-radius: 4px; cursor: pointer; background: #16213e;
         color: #8899aa; font-size: 0.85rem; border: 1px solid #2a2a4e; }
  .tab.active { background: #e94560; color: white; border-color: #e94560; }
  .tab:hover { border-color: #e94560; }
  table { width: 100%; border-collapse: collapse; margin-top: 12px; }
  th { text-align: left; padding: 8px 12px; background: #16213e; color: #8899aa;
       font-size: 0.75rem; text-transform: uppercase; letter-spacing: 1px; }
  td { padding: 6px 12px; border-bottom: 1px solid #1a1a3e; font-size: 0.85rem; }
  tr:hover td { background: #16213e; }
  .stage-badge { display: inline-block; padding: 2px 8px; border-radius: 3px;
                 font-size: 0.75rem; font-weight: bold; }
  .stage-ideas { background: #333; }
  .stage-research { background: #1565c0; }
  .stage-outline { background: #6a1b9a; }
  .stage-draft { background: #2e7d32; }
  .stage-factcheck { background: #e65100; }
  .stage-draft2 { background: #00838f; }
  .stage-done { background: #2e7d32; }
  .status-pending { color: #888; }
  .status-done { color: #4caf50; }
  .status-final { color: #4caf50; }
  .status-in-progress { color: #ff9800; }
  .status-error { color: #f44336; }
  .cost { color: #ff9800; font-weight: bold; }
  .retry-banner { background: #e65100; color: white; padding: 10px 16px; border-radius: 6px;
                  margin-bottom: 16px; font-weight: bold; font-size: 0.95rem; }
  .eject-btn { background: #6a1b9a; color: white; border: none; padding: 3px 10px; border-radius: 3px;
               font-size: 0.75rem; cursor: pointer; font-weight: bold; }
  .eject-btn:hover { background: #8e24aa; }
  .eject-btn:disabled { background: #555; cursor: default; }
</style>
</head>
<body>
<div class="layout">
<div class="main-panel">
  <h1>Pipeline Dashboard <span id="mode" class="mode dry">DRY RUN</span></h1>
  <div class="subtitle" id="subtitle">Loading...</div>
  <div id="retryBanner" class="retry-banner" style="display:none"></div>

  <div class="controls">
    <svg class="progress-ring" viewBox="0 0 52 52">
      <circle class="bg" cx="26" cy="26" r="22"/>
      <circle class="fg" id="ring" cx="26" cy="26" r="22" stroke-dasharray="138.2" stroke-dashoffset="0"/>
    </svg>
    <div>
      <div class="countdown" id="countdown">--</div>
      <div class="countdown-label" id="countdownLabel">Next cycle</div>
    </div>
    <button class="btn" id="stepBtn" onclick="doStep()">Run Now</button>
  </div>

  <div class="grid">
    <div class="card"><div class="number" id="pending">-</div><div class="label">Pending</div></div>
    <div class="card"><div class="number" id="research">-</div><div class="label">Research</div></div>
    <div class="card"><div class="number" id="outline">-</div><div class="label">Outline</div></div>
    <div class="card"><div class="number" id="draft">-</div><div class="label">Draft</div></div>
    <div class="card"><div class="number" id="factcheck">-</div><div class="label">Factcheck</div></div>
    <div class="card"><div class="number" id="draft2">-</div><div class="label">Draft 2</div></div>
    <div class="card"><div class="number" id="done">-</div><div class="label">Done</div></div>
    <div class="card"><div class="number cost" id="cost">-</div><div class="label">Session Cost</div></div>
    <div class="card"><div class="number" id="cycle">-</div><div class="label">Cycles</div></div>
  </div>

  <div style="font-size:0.7rem;color:#546e7a;text-transform:uppercase;letter-spacing:1px;margin-bottom:6px;">On Disk (all time)</div>
  <div class="grid" style="margin-bottom:20px;">
    <div class="card"><div class="number" id="dk-ideas">-</div><div class="label">Ideas</div></div>
    <div class="card"><div class="number" id="dk-research">-</div><div class="label">Research</div></div>
    <div class="card"><div class="number" id="dk-outline">-</div><div class="label">Outline</div></div>
    <div class="card"><div class="number" id="dk-draft">-</div><div class="label">Draft</div></div>
    <div class="card"><div class="number" id="dk-factcheck">-</div><div class="label">Factcheck</div></div>
    <div class="card"><div class="number" id="dk-draft2">-</div><div class="label">Draft 2</div></div>
    <div class="card"><div class="number" id="dk-draft3">-</div><div class="label">Draft 3</div></div>
    <div class="card"><div class="number cost" id="dk-cost">-</div><div class="label">Total Cost</div></div>
    <div class="card"><div class="number">&nbsp;</div><div class="label">&nbsp;</div></div>
  </div>

  <div class="project-tabs" id="projectTabs"></div>

  <table>
    <thead>
      <tr><th>Project</th><th>#</th><th>Length</th><th>Title</th><th>Arc</th><th>Stage</th><th>Status</th><th></th></tr>
    </thead>
    <tbody id="essayTable"></tbody>
  </table>
</div>

<div class="log-panel">
  <div class="log-header">
    <h2>Logs</h2>
    <div class="log-settings">
      <label><input type="checkbox" id="verboseToggle" onchange="toggleVerbose()"> Verbose</label>
      <div class="interval-ctrl">
        <label style="margin:0">Interval:</label>
        <input type="number" id="intervalInput" min="5" step="5" style="width:50px">
      </div>
      <div class="interval-ctrl">
        <label style="margin:0">Read:</label>
        <input type="number" id="readMeanInput" min="1" max="15" step="0.5" style="width:50px" title="Target reading time (minutes)">
        <label style="margin:0">&plusmn;</label>
        <input type="number" id="readSpreadInput" min="0.5" max="3" step="0.5" style="width:45px" title="Std deviation (minutes)">
        <button class="save-btn" id="saveIntervalBtn" onclick="saveSettings()">Save</button>
      </div>
    </div>
  </div>
  <div class="log-body" id="logBody"></div>
  <div class="log-count"><span id="logCount"></span> <span id="idleIndicator" style="color:#ff9800;display:none"></span> <button class="copy-btn" onclick="copyLogs()">Copy</button></div>
</div>
</div>

<script>
let cycleInterval = 15, secondsLeft = 15, activeProject = '', isRunning = false;
let verbose = false;
const circ = 2 * Math.PI * 22;

function updateRing() {
  const label = document.getElementById('countdownLabel');
  const ring = document.getElementById('ring');
  const cd = document.getElementById('countdown');
  const pct = cycleInterval > 0 ? secondsLeft / cycleInterval : 0;
  ring.setAttribute('stroke-dashoffset', String(circ * (1 - pct)));
  cd.textContent = secondsLeft + 's';
  if (isRunning) {
    ring.style.stroke = '#ff9800';
    label.textContent = 'Running';
  } else {
    ring.style.stroke = '#e94560';
    label.textContent = 'Next cycle';
  }
}

setInterval(function() { if (secondsLeft > 0) secondsLeft--; updateRing(); }, 1000);

var lastLogText = '';
var lastLogVersion = -1;
var lastNewLogTime = Date.now();
async function refreshLogs() {
  try {
    var data = await fetch('/api/logs?verbose=' + verbose).then(function(r) { return r.json(); });
    if (data.version === lastLogVersion) {
      updateIdleIndicator();
      return;
    }
    lastLogVersion = data.version;
    lastNewLogTime = Date.now();
    var reversed = (data.entries || []).slice().reverse();
    var box = document.getElementById('logBody');
    box.innerHTML = reversed.map(function(e) {
      return '<div class="log-line ' + e.level + '"><span style="color:#546e7a">' + esc(e.time) + '</span> ' + esc(e.message) + '</div>';
    }).join('');
    lastLogText = reversed.map(function(e) { return e.time + ' ' + e.message; }).join('\n');
    document.getElementById('logCount').textContent = reversed.length + ' entries';
    updateIdleIndicator();
  } catch(e) {}
}
function updateIdleIndicator() {
  var el = document.getElementById('idleIndicator');
  var ago = Math.floor((Date.now() - lastNewLogTime) / 1000);
  if (ago >= 5) {
    el.textContent = 'idle ' + ago + 's';
    el.style.display = 'inline';
  } else {
    el.style.display = 'none';
  }
}

async function refresh() {
  try {
    var s = await fetch('/api/status').then(function(r) { return r.json(); });
  } catch(e) { return; }
  cycleInterval = s.cycle_interval || 15;
  secondsLeft = s.seconds_left;
  isRunning = s.cycle_running || false;
  if (!document.getElementById('intervalInput').matches(':focus')) {
    document.getElementById('intervalInput').value = cycleInterval;
  }
  if (!document.getElementById('readMeanInput').matches(':focus')) {
    document.getElementById('readMeanInput').value = s.read_mean || 5;
  }
  if (!document.getElementById('readSpreadInput').matches(':focus')) {
    document.getElementById('readSpreadInput').value = s.read_spread || 1;
  }
  if (s.verbose !== undefined) {
    verbose = s.verbose;
    document.getElementById('verboseToggle').checked = verbose;
  }
  document.getElementById('pending').textContent = s.summary.pending || 0;
  document.getElementById('research').textContent = s.summary.research || 0;
  document.getElementById('outline').textContent = s.summary.outline || 0;
  document.getElementById('draft').textContent = s.summary.draft || 0;
  document.getElementById('factcheck').textContent = s.summary.factcheck || 0;
  document.getElementById('draft2').textContent = s.summary.draft2 || 0;
  document.getElementById('done').textContent = s.summary.done || 0;
  document.getElementById('cost').textContent = '$' + s.total_cost.toFixed(2);
  document.getElementById('cycle').textContent = s.cycle;
  var mode = document.getElementById('mode');
  if (s.dry_run) { mode.textContent = 'DRY RUN'; mode.className = 'mode dry'; }
  else { mode.textContent = 'LIVE'; mode.className = 'mode live'; }
  var banner = document.getElementById('retryBanner');
  if (s.retry_message) { banner.textContent = s.retry_message; banner.style.display = 'block'; }
  else { banner.style.display = 'none'; }
  var projects = s.projects || [];
  document.getElementById('subtitle').textContent =
    projects.map(function(p) { return p.name; }).join(', ') + ' \u2014 ' + (s.summary.done||0) + ' done of ' +
    Object.values(s.summary).reduce(function(a,b) { return a+b; },0) + ' total';

  var tabs = document.getElementById('projectTabs');
  var tabsHtml = '<div class="tab' + (activeProject===''?' active':'') + '" onclick="filterProject(\x27\x27)">All</div>';
  projects.forEach(function(p) {
    tabsHtml += '<div class="tab' + (activeProject===p.name?' active':'') +
      '" onclick="filterProject(\x27' + esc(p.name) + '\x27)"> ' + esc(p.name) + '</div>';
  });
  tabs.innerHTML = tabsHtml;
  updateRing();

  var url = activeProject ? '/api/essays?project=' + encodeURIComponent(activeProject) : '/api/essays';
  var essays = await fetch(url).then(function(r) { return r.json(); });
  var tbody = document.getElementById('essayTable');
  tbody.innerHTML = (essays||[]).map(function(e) {
    var idStr = e.book + '.' + e.part + '.' + e.order;
    var lenStr = e.word_count ? e.word_count.toLocaleString() + '-' + e.read_mins + 'min' : '';
    return '<tr>' +
    '<td>' + esc(e.project) + '</td>' +
    '<td>' + idStr + '</td>' +
    '<td>' + lenStr + '</td>' +
    '<td><a href="#" onclick="openFolder(\x27' + encodeURIComponent(e.project) + '\x27,\x27' + encodeURIComponent(e.slug) + '\x27);return false" style="color:#e94560;text-decoration:none">' + esc(e.title) + '</a></td>' +
    '<td>' + (e.arc ? esc(e.arc) : '') + '</td>' +
    '<td><span class="stage-badge stage-' + e.stage + '">' + e.stage + '</span></td>' +
    '<td class="status-' + e.status + '">' + e.status + (e.error ? ' <span title="' + esc(e.error) + '" style="cursor:help">\u26a0</span>' : '') + '</td>' +
    '<td>' + (function() {
      if (e.stage !== 'done') return '';
      if (e.ejected === 'current') return '<button class="eject-btn" style="background:#2e7d32" onclick="openDocx(\x27' + encodeURIComponent(e.project) + '\x27,\x27' + encodeURIComponent(e.slug) + '\x27)">Open</button>';
      return '<button class="eject-btn" onclick="doEject(\x27' + encodeURIComponent(e.project) + '\x27,\x27' + encodeURIComponent(e.slug) + '\x27,this)">' + (e.ejected === 'stale' ? 'Re-eject' : 'Eject') + '</button>';
    })() + '</td>' +
    '</tr>';
  }).join('');
}
function esc(s) { var d = document.createElement('div'); d.textContent = s; return d.innerHTML; }
function filterProject(p) { activeProject = p; refresh(); }
async function openFolder(project, slug) {
  await fetch('/api/open?project=' + project + '&slug=' + slug);
}
async function doEject(project, slug, btn) {
  btn.disabled = true; btn.textContent = 'Ejecting...';
  try {
    var resp = await fetch('/api/eject?project=' + project + '&slug=' + slug, { method: 'POST' });
    var data = await resp.json();
    if (data.ok) {
      btn.textContent = 'Done';
      btn.style.background = '#4caf50';
      setTimeout(function() { refresh(); }, 1500);
    } else {
      btn.textContent = 'Error';
      btn.style.background = '#f44336';
      btn.title = data.error || 'unknown error';
      setTimeout(function() { btn.disabled = false; btn.textContent = 'Eject'; btn.style.background = ''; btn.title = ''; }, 3000);
    }
  } catch(e) {
    btn.textContent = 'Error';
    btn.style.background = '#f44336';
    setTimeout(function() { btn.disabled = false; btn.textContent = 'Eject'; btn.style.background = ''; btn.title = ''; }, 3000);
  }
}
function openDocx(project, slug) {
  fetch('/api/open-docx?project=' + project + '&slug=' + slug);
}
async function doStep() {
  var btn = document.getElementById('stepBtn');
  btn.disabled = true; btn.textContent = 'Running...';
  await fetch('/api/step', { method: 'POST' });
  setTimeout(async function() { await refresh(); btn.disabled = false; btn.textContent = 'Run Now'; }, 2000);
}
function toggleVerbose() {
  verbose = document.getElementById('verboseToggle').checked;
  saveSettings();
  refreshLogs();
}
function copyLogs() {
  var btn = document.querySelector('.copy-btn');
  navigator.clipboard.writeText(lastLogText).then(function() {
    btn.textContent = 'Copied'; btn.classList.add('copied');
    setTimeout(function() { btn.textContent = 'Copy'; btn.classList.remove('copied'); }, 1500);
  });
}
async function saveSettings() {
  var btn = document.getElementById('saveIntervalBtn');
  var intVal = parseInt(document.getElementById('intervalInput').value, 10);
  var meanVal = parseFloat(document.getElementById('readMeanInput').value);
  var spreadVal = parseFloat(document.getElementById('readSpreadInput').value);
  var body = { verbose: verbose };
  if (intVal >= 5) body.cycle_interval = intVal;
  if (meanVal >= 1 && meanVal <= 15) body.read_mean = meanVal;
  if (spreadVal >= 0.5 && spreadVal <= 3) body.read_spread = spreadVal;
  await fetch('/api/settings', { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body) });
  btn.textContent = 'Saved';
  btn.classList.add('saved');
  setTimeout(function() { btn.textContent = 'Save'; btn.classList.remove('saved'); }, 1500);
}
async function refreshDiskStats() {
  var res = await fetch('/api/disk-stats');
  var d = await res.json();
  var c = d.counts || {};
  document.getElementById('dk-ideas').textContent = c.ideas || 0;
  document.getElementById('dk-research').textContent = c.research || 0;
  document.getElementById('dk-outline').textContent = c.outline || 0;
  document.getElementById('dk-draft').textContent = c.draft || 0;
  document.getElementById('dk-factcheck').textContent = c.factcheck || 0;
  document.getElementById('dk-draft2').textContent = c.draft2 || 0;
  document.getElementById('dk-draft3').textContent = c.draft3 || 0;
  document.getElementById('dk-cost').textContent = '$' + d.total_cost.toFixed(2);
}
refresh();
refreshLogs();
refreshDiskStats();
setInterval(refresh, 3000);
setInterval(refreshLogs, 500);
setInterval(refreshDiskStats, 10000);
</script>
</body>
</html>`
