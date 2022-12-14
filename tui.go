package main

import (
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
    TICKRATE = 100
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
	hline = `=======================================================================================`
)

// Style
const (
	sand       = lipgloss.Color("#C2B280")
	rock       = lipgloss.Color("#918e7d")
	empty      = lipgloss.Color("#333333")
    background = lipgloss.Color("#d3d3d3")
)

var (
	sandStyle  = lipgloss.NewStyle().Foreground(sand).Background(background)
	rockStyle  = lipgloss.NewStyle().Foreground(rock).Background(background)
	emptyStyle =  lipgloss.NewStyle().Foreground(background).Background(background)
)

func tickStats() tea.Cmd {
	return tea.Every(
		TICKRATE * time.Millisecond,
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
        msg := <-updateCh
        // grain hit bottom:
        if x.State == 0{
            m.cave[msg.Y][msg.X] = "O"
            m.cave[0][500] = "+"
            return m, tickStats()
        // grain falling:
        } else if x.State == -1{
            m.cave[msg.Y][msg.X] = ""
            m.cave[msg.Dy][msg.Dx] = "+"
            return m, tickStats()
        // reached the top:
        } else if x.State == 1{
            m.total = x.Total
            return m, nil
        }
	}
	return m, nil
}

func (m Model) caveString() string {
    cave := ``
    //Do stuff:
    for _, row := range m.cave {
        curr_s := ""
        i := 0
        for _, c := range row {
            if len(curr_s) > 0 && curr_s[len(curr_s)-1] != c{
                if curr_s[len(curr_s)-1] == ""{
                    cave += emptyStyle.Width(i).Render(curr_s)
                    curr_s = ""
                } else if curr_s[len(curr_s)-1] == "+" || curr_s[len(curr_s)-1] == "O"{
                    cave += sandStyle.Width(i).Render(curr_s)
                    curr_s = ""
                } else if curr_s[len(curr_s)-1] == "#"{
                    cave += rockStyle.Width(i).Render(curr_s)
                    curr_s = ""
                }
                i = 0
            }
            // Empty
            if c == ""{
                curr_s += '.'
            } else if c == "+" || c == "O"{
                curr_s += c
            } else if c == "#"{
                curr_s += c
            }
            i++
        }
        if curr_s[len(curr_s)-1] == ""{
            cave += emptyStyle.Width(i).Render(curr_s)
            curr_s = ""
        } else if curr_s[len(curr_s)-1] == "+" || curr_s[len(curr_s)-1] == "O"{
            cave += sandStyle.Width(i).Render(curr_s)
            curr_s = ""
        } else if curr_s[len(curr_s)-1] == "#"{
            cave += rockStyle.Width(i).Render(curr_s)
            curr_s = ""
        }
        curr_s += '\n'
    }
}

func (m Model) View() string {
	body := fmt.Sprintf(
		`
            %s
	%s
%s
		    %s
        `,
        // title
		sandStyle.Width(150).Render(TITLE),
        rockStyle.Width(150).Render(hline),
        // grid
        m.caveString(),
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
    cave   [][]string
    caveCh chan UpdateMsg
    total  int
}

func InitModel(updateCh chan UpdateMsg, x int, y int) Model {
    m := Model{
        cave: [y][x]string{},
        caveCh: updateCh,
	}
    m.cave[0][500] = "+"
	return m
}
