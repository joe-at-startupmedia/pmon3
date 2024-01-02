package output

import (
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

func DescTable(tbData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"--", "Desc"})
	hColor := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor}
	table.SetHeaderColor(hColor, hColor)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	for index, row := range tbData {
		if index == 0 {
			switch row[1] {
			case "running":
				table.Rich(row, []tablewriter.Colors{{}, {tablewriter.Bold, tablewriter.FgHiGreenColor}})
			case "stopped":
				table.Rich(row, []tablewriter.Colors{{}, {tablewriter.Bold, tablewriter.FgHiYellowColor}})
			case "failed":
				table.Rich(row, []tablewriter.Colors{{}, {tablewriter.Bold, tablewriter.FgRedColor}})
			default:
				table.Append(row)
			}
		} else {
			table.Append(row)
		}
	}

	//table.SetRowLine(true)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.Render()
}

// render single one table row
func TableOne(data []string) {
	Table([][]string{data})
}

var (
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "╥",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "╨",
		InnerJunction:  "╫",

		InnerDivider: "║",
	}
)

type Model struct {
	tableModel table.Model
}

const (
	columnKeyID     = "id"
	columnKeyName   = "name"
	columnKeyPid    = "pid"
	columnKeyStatus = "status"
	columnKeyUser   = "user"
	columnKeyCpu    = "cpu"
	columnKeyMem    = "mem"
	columnKeyDate   = "date"
)

func getStatusColor(status string) string {
	switch status {
	case "running":
		return "#16ff16"
	case "stopped":
		return "#ffff00"
	case "failed":
		return "#ff0000"
	case "init":
		return "#808080"
	}
	return "#ffffff"
}

func NewModel(tbData [][]string) Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, columnKeyID, 5).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#afaf00")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyName, "Name", 15),
		table.NewColumn(columnKeyPid, "PID", 10),
		table.NewColumn(columnKeyStatus, columnKeyStatus, 10),
		table.NewColumn(columnKeyUser, columnKeyUser, 15),
		table.NewColumn(columnKeyCpu, "CPU", 5),
		table.NewColumn(columnKeyMem, columnKeyMem, 10),
		table.NewColumn(columnKeyDate, columnKeyDate, 20),
	}

	var rows []table.Row

	for _, row := range tbData {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyID:     row[0],
			columnKeyName:   row[1],
			columnKeyPid:    row[2],
			columnKeyStatus: table.NewStyledCell(row[3], lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(row[3])))),
			columnKeyUser:   row[4],
			columnKeyCpu:    row[5],
			columnKeyMem:    row[6],
			columnKeyDate:   row[7],
		}))
	}

	model := Model{
		// Throw features in... the point is not to look good, it's just reference!
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))).
			SelectableRows(false).
			Focused(false).
			Border(customBorder).
			WithBaseStyle(
				lipgloss.NewStyle().
					Align(lipgloss.Left).
					Bold(true),
			).
			SortByAsc(columnKeyID),
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

func Table(tbData [][]string) {

	p := tea.NewProgram(NewModel(tbData))

	if err := p.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
