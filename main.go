package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    
    tea "github.com/charmbracelet/bubbletea"
)

type Coords struct {
    x int
    y int
}

func min(x int, y int) int{
    if x < y{
        return x
    } else {
        return y
    }
}

func max(x int, y int) int{
    if x > y{
        return x
    } else {
        return y
    }
}

func coordsToInt(str string) []int{
    coords := []int{}
    for _, c := range strings.Split(str, ","){
        val, _ := strconv.Atoi(c);
        coords = append(coords, val)
    }
    return coords
}

func parse_grid(in string) (map[Coords]bool, int, int){
    walls := make(map[Coords]bool)
    max_y, max_x := 0, 0
    for _, line := range strings.Split(in, "\n"){
        coords := strings.Split(strings.Trim(line, "\n"), " -> ")
        for i:=0; i < len(coords)-1; i++{
            s := coordsToInt(coords[i])
            e := coordsToInt(coords[i+1])
            if s[0] == e[0]{
                for j:=min(s[1], e[1]); j<=max(s[1], e[1]); j++{
                    wall := Coords{x:s[0],y:j}
                    if j > max_y{
                        max_y = j
                    }
                    if s[0] > max_x{
                        max_x = s[0]
                    }
                    walls[wall] = true
                }
            } else {
                for j:=min(s[0], e[0]); j<=max(s[0], e[0]); j++{
                    wall := Coords{x:j,y:s[1]}
                    if s[1] > max_y{
                        max_y = s[1]
                    }
                    if j > max_x{
                        max_x = j
                    }
                    walls[wall] = true
                }
                
            }
        }
    }
    return walls, max_y+2, max_x
}

func pour_sand(walls map[Coords]bool, max_y int, updateCh chan UpdateMsg){
    total := 0
    curr_grain := Coords{x:500, y:0}
    for !walls[Coords{x:500, y:0}] {
        y := curr_grain.y+1
        if !walls[Coords{x:curr_grain.x, y:y}] && y != max_y{
            updateCh<-UpdateMsg{
                X: curr_grain.x,
                Y: curr_grain.y,
                Dx: curr_grain.x,
                Dy: y,
                State: -1,
            }
            curr_grain = Coords{x:curr_grain.x, y:y}

        } else if !walls[Coords{x:curr_grain.x-1, y:y}] && y != max_y{
            updateCh<-UpdateMsg{
                X: curr_grain.x,
                Y: curr_grain.y,
                Dx: curr_grain.x-1,
                Dy: y,
                State: -1,
            }
            curr_grain = Coords{x:curr_grain.x-1, y:y}

        } else if !walls[Coords{x:curr_grain.x+1, y:y}] && y != max_y{
            updateCh<-UpdateMsg{
                X: curr_grain.x,
                Y: curr_grain.y,
                Dx: curr_grain.x+1,
                Dy: y,
                State: -1,
            }
            curr_grain = Coords{x:curr_grain.x+1, y:y}

        } else {
            total += 1
            updateCh<-UpdateMsg{
                X: curr_grain.x,
                Y: curr_grain.y,
                State: 0,
                Total: total
            }
            walls[curr_grain] = true
            curr_grain = Coords{x:500, y:0}
        }
    }
    updateCh<-UpdateMsg{
        State: 1,
        Total: total
    }
}

func main() {
    file, err := os.ReadFile("input.txt")
    if err != nil{
        fmt.Println("Couldn't open input")
        os.Exit(-1)
    }
    walls, max_y, max_x := parse_grid(string(file))
    updateCh := make(chan UpdateMsg, 1)
    go pour_sand(walls, max_y, updateCh)
    tui_model := InitModel(updateCh, max_x, max_y)
    p := tea.NewProgram(tui_model)
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
    //fmt.Printf("Dims: %vx%v\n", max_x, max_y)
    //fmt.Printf("Grains of Sand: %v\n", pour_sand(walls, max_y))
}
