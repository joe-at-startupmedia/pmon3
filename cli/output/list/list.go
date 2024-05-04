package table_list

import (
	"log"
	"os"
	"pmon3/pmond/utils/conv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	CustomBorder = table.Border{
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
	columnKeyID           = "id"
	columnKeyName         = "name"
	columnKeyPid          = "pid"
	columnKeyRestartCount = "restart_count"
	columnKeyStatus       = "status"
	columnKeyUser         = "user"
	columnKeyCpu          = "cpu"
	columnKeyMem          = "mem"
	columnKeyDate         = "date"
)

func GetStatusColor(status string) string {
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
	return "#646464"
}

func NewModel(tbData [][]string) Model {

	//min column sizing
	widthData := [9]int{
		2,
		5,
		3,
		1,
		6,
		5,
		4,
		7,
		19,
	}

	var rows []table.Row
	for _, row := range tbData {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyID:           conv.StrToUint32(row[0]),
			columnKeyName:         row[1],
			columnKeyPid:          row[2],
			columnKeyRestartCount: row[3],
			columnKeyStatus:       table.NewStyledCell(row[4], lipgloss.NewStyle().Foreground(lipgloss.Color(GetStatusColor(row[4])))),
			columnKeyUser:         row[5],
			columnKeyCpu:          row[6],
			columnKeyMem:          row[7],
			columnKeyDate:         row[8],
		}))

		//peak finder
		n := 0
		for n < 9 {
			colLength := len(row[n]) + 1
			if colLength > widthData[n] {
				widthData[n] = colLength
			}
			n++
		}
	}

	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", widthData[0]).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#afaf00")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyName, "Name", widthData[1]),
		table.NewColumn(columnKeyPid, "PID", widthData[2]),
		table.NewColumn(columnKeyRestartCount, "⟳", widthData[3]),
		table.NewColumn(columnKeyStatus, "Status", widthData[4]),
		table.NewColumn(columnKeyUser, "User", widthData[5]),
		table.NewColumn(columnKeyCpu, "CPU", widthData[6]),
		table.NewColumn(columnKeyMem, "Memory", widthData[7]),
		table.NewColumn(columnKeyDate, "Date", widthData[8]),
	}

	model := Model{
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))).
			SelectableRows(false).
			Focused(false).
			Border(CustomBorder).
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

func Render(tbData [][]string) {

	p := tea.NewProgram(NewModel(tbData))

	if err := p.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
