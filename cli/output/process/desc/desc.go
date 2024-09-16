package table_desc

import (
	"log"
	"os"
	"pmon3/cli/output/process/list"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	tableModel table.Model
}

const (
	columnKeyProperty    = "property"
	columnKeyDescription = "description"
)

func NewModel(tbData [][]string) Model {

	//min column sizing
	widthData := [9]int{
		15,
		15,
	}

	var rows []table.Row
	for i, row := range tbData {
		if i == 0 {
			rows = append(rows, table.NewRow(table.RowData{
				columnKeyProperty:    row[0],
				columnKeyDescription: table.NewStyledCell(row[1], lipgloss.NewStyle().Foreground(lipgloss.Color(table_list.GetStatusColor(row[1])))),
			}))
		} else {
			rows = append(rows, table.NewRow(table.RowData{
				columnKeyProperty:    row[0],
				columnKeyDescription: row[1],
			}))
		}

		//peak finder
		n := 0
		for n < 2 {
			colLength := len(row[n]) + 1
			if colLength > widthData[n] {
				widthData[n] = colLength
			}
			n++
		}
	}

	columns := []table.Column{
		table.NewColumn(columnKeyProperty, columnKeyProperty, widthData[0]).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#afaf00")).
				Align(lipgloss.Left)),
		table.NewColumn(columnKeyDescription, columnKeyDescription, widthData[1]),
	}

	model := Model{
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))).
			Focused(false).
			Border(table_list.CustomBorder).
			WithBaseStyle(
				lipgloss.NewStyle().
					Align(lipgloss.Left).
					Bold(true),
			).
			WithMultiline(true),
	}

	return model
}

func (m Model) Init() tea.Cmd {
	//non-interactive
	return tea.Quit
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}
	body.WriteString(m.tableModel.View())
	body.WriteString("\n")
	return body.String()
}

func Render(tbData [][]string) {

	p := tea.NewProgram(NewModel(tbData))

	if err := p.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
