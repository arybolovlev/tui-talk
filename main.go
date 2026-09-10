package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	cyan   = "\033[36m"
	yellow = "\033[93m"
	muted  = "\033[90m"
	white  = "\033[97m"
)

type model struct {
	page   int
	width  int
	height int
}

func initialModel() model {
	return model{page: 0, width: 80, height: 24}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "left", "h":
			if m.page > 0 {
				m.page--
			}
		case "right", "l":
			if m.page < pageCount()-1 {
				m.page++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch p := m.page; p {
	case 0:
		return renderFrame(titleArt(), m.width, m.height, "KubeCon Europe 27", true, p+1, "by Eleni and Aleksandr", true)
	case 1:
		return renderFrame(nil, m.width, m.height, "What is wrong with Ingress?", false, p+1, "", true)
	default:
		return renderFrame(nil, m.width, m.height, "Thank you! Questions?", false, p+1, "", true)
	}
}

func titleArt() []string {
	return []string{
		cyan + "An overheard conversation about" + reset,
		"",
		bold + yellow + "GatewayAPI" + reset,
	}
}

func renderFrame(content []string, width, height int, title string, centered bool, page int, bottomLabel string, showFooter bool) string {
	for _, line := range content {
		if lineWidth := len([]rune(stripANSI(line))) + 2; lineWidth > width {
			width = lineWidth
		}
	}
	if width < 20 {
		width = 20
	}
	if height < 6 {
		height = 6
	}

	innerWidth := width - 2
	footerRows := 0
	if showFooter {
		footerRows = 1
	}
	bodyHeight := height - 2 - footerRows // top and bottom borders, plus optional footer
	if len(content) > bodyHeight {
		content = content[:bodyHeight]
	}

	lines := []string{topBorder(title, width)}
	for i := 0; i < (bodyHeight-len(content))/2; i++ {
		lines = append(lines, frameRow("", innerWidth))
	}
	for _, line := range content {
		if centered {
			line = center(line, innerWidth)
		}
		lines = append(lines, frameRow(line, innerWidth))
	}
	for len(lines) < height-1-footerRows {
		lines = append(lines, frameRow("", innerWidth))
	}

	if showFooter {
		footer := fmt.Sprintf("%s←/→%s navigate  %s•%s  %s%d/%d%s  %sq%s quit",
			muted, reset, muted, reset, muted, page, pageCount(), reset, muted, reset)
		lines = append(lines, frameRow(center(footer, innerWidth), innerWidth))
	}
	lines = append(lines, bottomBorder(width, bottomLabel))

	return strings.Join(lines, "\n")
}

func topBorder(title string, width int) string {
	fill := width - len([]rune(title)) - 5
	if fill < 0 {
		fill = 0
	}
	return cyan + "╔═ " + bold + white + title + reset + cyan + " " + strings.Repeat("═", fill) + "╗" + reset
}

func bottomBorder(width int, label string) string {
	if label == "" {
		return cyan + "╚" + strings.Repeat("═", width-2) + "╝" + reset
	}

	label = " " + label + " "
	fill := width - 2 - len([]rune(label))
	if fill < 0 {
		return cyan + "╚" + strings.Repeat("═", width-2) + "╝" + reset
	}
	return cyan + "╚" + strings.Repeat("═", fill) + bold + white + label + reset + cyan + "╝" + reset
}

func frameRow(s string, width int) string {
	visible := []rune(stripANSI(s))
	if len(visible) > width {
		return cyan + "║" + reset + string(visible[:width]) + cyan + "║" + reset
	}
	return cyan + "║" + reset + s + strings.Repeat(" ", width-len(visible)) + cyan + "║" + reset
}

func center(s string, width int) string {
	visibleWidth := len([]rune(stripANSI(s)))
	if visibleWidth >= width {
		return s
	}
	return strings.Repeat(" ", (width-visibleWidth)/2) + s
}

func stripANSI(s string) string {
	return strings.NewReplacer(cyan, "", yellow, "", muted, "", white, "", bold, "", reset, "").Replace(s)
}

func pageCount() int {
	return 3
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("error:", err)
	}
}
