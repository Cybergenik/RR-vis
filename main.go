package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type State uint8

type CaveType uint8

const (
	POURING State = 0
	DONE    State = 1

	SAND CaveType = 0
	WALL CaveType = 1
)

type Cave struct {
	State  State
	Total  int
	Grains []Grain
}

type Grain struct {
	X    int
	Y    int
	Dx   int
	Dy   int
	Type CaveType
}

type Coord struct {
	x int
	y int
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func coordsToInt(str string) []int {
	coords := []int{}
	for _, c := range strings.Split(str, ",") {
		val, _ := strconv.Atoi(c)
		coords = append(coords, val)
	}
	return coords
}

func parse_grid(in string) (map[Coord]bool, int, int) {
	walls := make(map[Coord]bool)
	max_y, max_x := 0, 0
	for _, line := range strings.Split(in, "\n") {
		coords := strings.Split(strings.Trim(line, "\n"), " -> ")
		for i := 0; i < len(coords)-1; i++ {
			s := coordsToInt(coords[i])
			e := coordsToInt(coords[i+1])
			if s[0] == e[0] {
				for j := min(s[1], e[1]); j <= max(s[1], e[1]); j++ {
					wall := Coord{x: s[0], y: j}
					if j > max_y {
						max_y = j
					}
					if s[0] > max_x {
						max_x = s[0]
					}
					walls[wall] = true
				}
			} else {
				for j := min(s[0], e[0]); j <= max(s[0], e[0]); j++ {
					wall := Coord{x: j, y: s[1]}
					if s[1] > max_y {
						max_y = s[1]
					}
					if j > max_x {
						max_x = j
					}
					walls[wall] = true
				}

			}
		}
	}
	return walls, max_y + 2, max_x
}

func pour_sand(walls map[Coord]bool, max_y int, updateCh chan Cave, doUpdate bool) {
	total := 0
	grains := []Coord{{x: 500, y: 0}}
	for !walls[Coord{x: 500, y: 0}] {
		walls_i := -1
		updates := []Grain{}
		for i, curr_grain := range grains {
			y := curr_grain.y + 1
			if !walls[Coord{x: curr_grain.x, y: y}] && y != max_y {
				updates = append(updates, Grain{
					X:    curr_grain.x,
					Y:    curr_grain.y,
					Dx:   curr_grain.x,
					Dy:   y,
					Type: SAND,
				})
				grains[i] = Coord{x: curr_grain.x, y: y}

			} else if !walls[Coord{x: curr_grain.x - 1, y: y}] && y != max_y {
				updates = append(updates, Grain{
					X:    curr_grain.x,
					Y:    curr_grain.y,
					Dx:   curr_grain.x - 1,
					Dy:   y,
					Type: SAND,
				})
				grains[i] = Coord{x: curr_grain.x - 1, y: y}

			} else if !walls[Coord{x: curr_grain.x + 1, y: y}] && y != max_y {
				updates = append(updates, Grain{
					X:    curr_grain.x,
					Y:    curr_grain.y,
					Dx:   curr_grain.x + 1,
					Dy:   y,
					Type: SAND,
				})
				grains[i] = Coord{x: curr_grain.x + 1, y: y}

			} else {
				total += 1
				updates = append(updates, Grain{
					X:    curr_grain.x,
					Y:    curr_grain.y,
					Type: WALL,
				})
				walls[curr_grain] = true
				walls_i = i
			}
		}
		if walls_i >= 0 {
			grains = grains[walls_i+1:]
			grains = append(grains, Coord{x: 500, y: 0})
		}
		if doUpdate {
			updateCh <- Cave{Total: total, Grains: updates, State: POURING}
		}
	}
	//No grains, end of animation.
	updateCh <- Cave{Total: total, State: DONE}
}

func main() {
	input := flag.String("input", "", "input file")
	noTui := flag.Bool("no-tui", false, "do Tui or not (default false)")
	flag.Parse()
	if *input == "" {
		fmt.Println("No input provided")
		os.Exit(1)
	}
	file, err := os.ReadFile(*input)
	if err != nil {
		fmt.Println("Couldn't open input")
		os.Exit(1)
	}
	walls, max_y, max_x := parse_grid(string(file))
	updateCh := make(chan Cave, 10)
	if *noTui {
		go pour_sand(walls, max_y, updateCh, false)
		cave := <-updateCh
		fmt.Printf("Total grains of sand: %v\n", cave.Total)
	} else {
		tui_model := InitModel(updateCh, walls, max_x+1, max_y)
		go pour_sand(walls, max_y, updateCh, true)
		p := tea.NewProgram(tui_model)
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
