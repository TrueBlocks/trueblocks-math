package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type item struct {
	Type       string
	Slug       string
	Title      string
	Book       string
	Part       int
	PartTitle  string
	Order      int
	Hook       string
	HiddenMath string
	Arc        string
	Ending     string
}

func parsePlanFile(path string, bookNum string) ([]item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	var items []item
	var currentPart int
	var currentPartTitle string
	partRe := regexp.MustCompile(`^## Part (\d+):\s*(.+)$`)
	introRe := regexp.MustCompile(`^## Book Introduction`)
	scanner := bufio.NewScanner(f)
	inTable := false
	headerSeen := false

	for scanner.Scan() {
		line := scanner.Text()

		if m := partRe.FindStringSubmatch(line); m != nil {
			p, _ := strconv.Atoi(m[1])
			currentPart = p
			currentPartTitle = strings.TrimSpace(m[2])
			inTable = false
			headerSeen = false
			continue
		}

		if introRe.MatchString(line) {
			currentPart = 0
			currentPartTitle = "Book Introduction"
			inTable = false
			headerSeen = false
			continue
		}

		if strings.HasPrefix(line, "| #") || strings.HasPrefix(line, "| -") {
			if strings.HasPrefix(line, "| #") {
				inTable = true
				headerSeen = false
			} else if inTable {
				headerSeen = true
			}
			continue
		}

		if !inTable || !headerSeen {
			if !strings.HasPrefix(line, "|") {
				inTable = false
				headerSeen = false
			}
			continue
		}

		if !strings.HasPrefix(line, "|") {
			inTable = false
			headerSeen = false
			continue
		}

		cols := splitTableRow(line)
		if len(cols) < 8 {
			continue
		}

		num := strings.TrimSpace(cols[0])
		typ := strings.TrimSpace(cols[1])
		slug := strings.TrimSpace(cols[2])
		title := strings.TrimSpace(cols[3])
		hook := strings.TrimSpace(cols[4])
		hiddenMath := strings.TrimSpace(cols[5])
		arc := strings.TrimSpace(cols[6])
		ending := strings.TrimSpace(cols[7])

		if typ == "" || slug == "" {
			continue
		}

		order := 0
		if num != "—" && num != "-" {
			order, _ = strconv.Atoi(num)
		}

		if arc == "-" {
			arc = ""
		}
		if ending == "-" {
			ending = ""
		}

		it := item{
			Type:       typ,
			Slug:       slug,
			Title:      title,
			Book:       bookNum,
			Part:       currentPart,
			PartTitle:  currentPartTitle,
			Order:      order,
			Hook:       hook,
			HiddenMath: hiddenMath,
			Arc:        arc,
			Ending:     ending,
		}
		items = append(items, it)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning %s: %w", path, err)
	}

	return items, nil
}

func splitTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func toSlug(label string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(label), " ", "-"))
}

func loadAllItems(designDir string) ([]item, error) {
	plans := []struct {
		file string
		book string
	}{
		{"Plan for Book I.md", "I"},
		{"Plan for Book II.md", "II"},
		{"Plan for Book III.md", "III"},
	}

	var all []item
	for _, p := range plans {
		path := filepath.Join(designDir, p.file)
		items, err := parsePlanFile(path, p.book)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
	}
	return all, nil
}

func main() {
	designDir := "../../design"
	if len(os.Args) > 1 {
		designDir = os.Args[1]
	}

	items, err := loadAllItems(designDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading plans: %v\n", err)
		os.Exit(1)
	}

	baseDir := "pipeline/ideas"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().Format("2006-01-02")

	for _, it := range items {
		var ideaContent string
		switch it.Type {
		case "section":
			ideaContent = fmt.Sprintf("# %s\n\n"+
				"**Slug:** %s\n"+
				"**Type:** section\n\n"+
				"## Part Title Page\n\n"+
				"This is a section divider for Part %d: %s.\n\n"+
				"## Placement\n\n"+
				"- **Book:** %s\n"+
				"- **Part %d:** %s\n",
				it.Title, it.Slug, it.Part, it.PartTitle, it.Book, it.Part, it.PartTitle)
		case "introduction":
			ideaContent = fmt.Sprintf("# %s\n\n"+
				"**Slug:** %s\n"+
				"**Type:** introduction\n\n"+
				"## Book Introduction\n\n"+
				"%s\n\n"+
				"## The Hidden Math\n\n"+
				"%s\n\n"+
				"## Placement\n\n"+
				"- **Book:** %s\n",
				it.Title, it.Slug, it.Hook, it.HiddenMath, it.Book)
		default:
			ideaContent = fmt.Sprintf("# %s\n\n"+
				"**Slug:** %s\n"+
				"**Type:** essay\n\n"+
				"## The Everyday Experience\n\n"+
				"%s\n\n"+
				"## The Hidden Math\n\n"+
				"%s\n\n"+
				"## Placement\n\n"+
				"- **Book:** %s\n"+
				"- **Part %d:** %s\n"+
				"- **Order:** %d\n",
				it.Title, it.Slug, it.Hook, it.HiddenMath, it.Book, it.Part, it.PartTitle, it.Order)
		}

		arcSlug := toSlug(it.Arc)
		endingSlug := toSlug(it.Ending)

		metaContent := fmt.Sprintf("slug: %s\n"+
			"title: \"%s\"\n"+
			"type: %s\n"+
			"book: %s\n"+
			"part: %d\n"+
			"part_title: \"%s\"\n"+
			"order: %d\n"+
			"status: pending\n"+
			"model: claude-sonnet-4-20250514\n"+
			"arc: \"%s\"\n"+
			"ending: \"%s\"\n"+
			"created: %s\n",
			it.Slug, it.Title, it.Type, it.Book, it.Part, it.PartTitle, it.Order, arcSlug, endingSlug, now)

		mdPath := filepath.Join(baseDir, it.Slug+".md")
		metaPath := filepath.Join(baseDir, it.Slug+".meta.yaml")

		if err := os.WriteFile(mdPath, []byte(ideaContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", mdPath, err)
			os.Exit(1)
		}

		if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", metaPath, err)
			os.Exit(1)
		}

		fmt.Printf("Created: %s [%s] (.md + .meta.yaml)\n", it.Slug, it.Type)
	}

	fmt.Printf("\nTotal: %d items scaffolded in %s/\n", len(items), baseDir)

	typeCounts := map[string]map[string]int{}
	for _, it := range items {
		if typeCounts[it.Book] == nil {
			typeCounts[it.Book] = map[string]int{}
		}
		typeCounts[it.Book][it.Type]++
	}
	for _, b := range []string{"I", "II", "III"} {
		c := typeCounts[b]
		fmt.Printf("  Book %s: %d essays, %d sections, %d introductions\n",
			b, c["essay"], c["section"], c["introduction"])
	}
}
