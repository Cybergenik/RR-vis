package main

import (
	"fmt"
	"time"
    "strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
    TICKRATE = 1
    OFFSET = 400
    funnel = 
`                                                                                                  \   /
                                                                                                   \ /
                                                                                                    *
`
	TITLE = `
                 _______ _______ _______ _______ _      _________________           _______ _______ _______ _______ _______         _______________________ 
                (  ____ (  ____ (  ____ (  ___  ( \     \__   __\__   __|\     /|  (  ____ (  ____ (  ____ (  ____ (  ____ |\     /(  ___  \__   __(  ____ )
                | (    )| (    \| (    \| (   ) | (        ) (     ) (  | )   ( |  | (    )| (    \| (    \| (    \| (    )| )   ( | (   ) |  ) (  | (    )|
                | (____)| (__   | |     | |   | | |        | |     | |  | (___) |  | (____)| (__   | (_____| (__   | (____)| |   | | |   | |  | |  | (____)|
                |     __|  __)  | | ____| |   | | |        | |     | |  |  ___  |  |     __|  __)  (_____  |  __)  |     __( (   ) | |   | |  | |  |     __)
                | (\ (  | (     | | \_  | |   | | |        | |     | |  | (   ) |  | (\ (  | (           ) | (     | (\ (   \ \_/ /| |   | |  | |  | (\ (   
                | ) \ \_| (____/| (___) | (___) | (____/___) (___  | |  | )   ( |  | ) \ \_| (____//\____) | (____/| ) \ \__ \   / | (___) ___) (__| ) \ \__
                |/   \__(_______(_______(_______(_______\_______/  )_(  |/     \|  |/   \__(_______\_______(_______|/   \__/  \_/  (_______\_______|/   \__/
                
`
	hline = `============================================================================================================================`
)

// Style
const (
    //bebd8f
	sand       = lipgloss.Color("#C2B280")
    hardSand   = lipgloss.Color("#ffcf5c")
	rock       = lipgloss.Color("#D3D3D3")
)

var (
	sandStyle  = lipgloss.NewStyle().Foreground(sand)
    setSandStyle = lipgloss.NewStyle().Foreground(hardSand)
	rockStyle  = lipgloss.NewStyle().Foreground(rock)
)

type TickMsg time.Time

func tickStats() tea.Cmd {
	return tea.Every(
		TICKRATE * time.Nanosecond,
		func(t time.Time) tea.Msg {
			return TickMsg(t)
		},
	)
}

func (m Model) Init() tea.Cmd {
    cmds := []tea.Cmd{tea.ClearScreen, tickStats()}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case TickMsg:
        update := <-m.updateCh
        // grain hit bottom:
        if update.State == 0{
            m.total = update.Total
            m.cave[update.Y][update.X-OFFSET] = "O"
            m.caveStr[update.Y] = caveRow(m.cave[update.Y])
            m.cave[0][500-OFFSET] = "+"
            m.caveStr[0] = caveRow(m.cave[0])
            return m, tickStats()
        // grain falling:
        } else if update.State == -1{
            m.cave[update.Y][update.X-OFFSET] = " "
            m.caveStr[update.Y] = caveRow(m.cave[update.Y])
            m.cave[update.Dy][update.Dx-OFFSET] = "+"
            m.caveStr[update.Dy] = caveRow(m.cave[update.Dy])
            return m, tickStats()
        // reached the top:
        } else if update.State == 1{
            m.total = update.Total
            return m, nil
        }
	}
	return m, nil
}

func caveRow(row []string) string {
    cave := strings.Builder{}
    curr_s := strings.Builder{}
    prev := ""
    i := 0
    for _, c := range row {
        if prev != "" && prev != c{
            if prev == " "{
                cave.WriteString(curr_s.String())
            } else if prev == "+"{
                cave.WriteString(sandStyle.Width(i).Render(curr_s.String()))
            } else if prev == "O"{
                cave.WriteString(setSandStyle.Width(i).Render(curr_s.String()))
            } else if prev == "#"{
                cave.WriteString(rockStyle.Width(i).Render(curr_s.String()))
            }
            curr_s = strings.Builder{}
            i = 0
        }
        curr_s.WriteString(c)
        prev = c
        i++
    }
    if prev == " "{
        cave.WriteString(curr_s.String())
    } else if prev == "+"{
        cave.WriteString(sandStyle.Width(i).Render(curr_s.String()))
    } else if prev == "O"{
        cave.WriteString(setSandStyle.Width(i).Render(curr_s.String()))
    } else if prev == "#"{
        cave.WriteString(rockStyle.Width(i).Render(curr_s.String()))
    }
    return cave.String()
}

func (m Model) View() string {
    //fmt.Println(m.caveString())
	body := fmt.Sprintf(
		`
%s
                        %s
                                                                                %s
%s
%s
		                                                            %s
        `,
        // title
		sandStyle.Width(200).Render(TITLE),
        rockStyle.Width(150).Render(hline),
        rockStyle.Width(15).Bold(true).Render(fmt.Sprintf("Total: %v", m.total)),
        funnel,
        //grid,
        strings.Join(m.caveStr, "\n"),
		// quit
		rockStyle.Width(30).Render("Press Esc or Ctrl+C to quit"),
	)

	return body
}

type UpdateMsg struct {
    X       int
    Y       int
    Dx      int
    Dy      int
    State   int
    Total   int
}

type Model struct {
    cave     [][]string
    caveStr  []string
    updateCh chan UpdateMsg
    total    int
}


func InitModel(updateCh chan UpdateMsg, walls map[Coords]bool, max_x int, max_y int) Model {
    m := Model{
        cave: make([][]string, max_y),
        caveStr: make([]string, max_y),
        updateCh: updateCh,
	}
    for y := range m.cave {
        curr_row := make([]string, max_x-(OFFSET/2))
        for x:=0; x<(max_x-(OFFSET/2));x++ {
            loc := Coords{x: x+OFFSET, y: y}
            if walls[loc] || y == len(m.cave)-1 {
                curr_row[x] = "#"
            } else {
                curr_row[x] = " "
            }
        }
        m.caveStr[y] = caveRow(curr_row)
        m.cave[y] = curr_row
    }
    m.cave[0][500-OFFSET] = "+"
	return m
}

