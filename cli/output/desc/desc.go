package table_desc

import (
	"log"
	"os"
	table_list "pmon3/cli/output/list"
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
	columns := []table.Column{
		table.NewColumn(columnKeyProperty, columnKeyProperty, 15).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#afaf00")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyDescription, columnKeyDescription, 55),
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
	}

	model := Model{
		// Throw features in... the point is not to look good, it's just reference!
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
