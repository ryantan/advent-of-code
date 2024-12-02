package main

import (
	"bufio"
	"os"
	"strconv"
)

type Tree struct {
	Height int
	Row    int
	Col    int
	See    map[string]int
}

func (t *Tree) Score() int {
	return t.See["LEFT"] * t.See["RIGHT"] * t.See["TOP"] * t.See["BOTTOM"]
}

func NewTree(h int, i int, j int) *Tree {
	return &Tree{Height: h, Row: i, Col: j, See: map[string]int{
		"LEFT": 0, "RIGHT": 0, "TOP": 0, "BOTTOM": 0,
	}}
}

type Stack struct {
	items []*Tree
}

func (s *Stack) Push(t *Tree) {
	s.items = append(s.items, t)
}

func (s *Stack) Length() int {
	return len(s.items)
}

func (s *Stack) Pop() *Tree {
	item := s.items[s.Length()-1]
	s.items = s.items[:s.Length()-1]
	return item
}

func (s *Stack) Peek() *Tree {
	return s.items[s.Length()-1]
}

func readFile(filename string) [][]*Tree {
	grid := [][]*Tree{}
	file, _ := os.Open(filename)
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	i := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		row := []*Tree{}
		for j, c := range line {
			h, _ := strconv.Atoi(string(c))
			row = append(row, NewTree(h, i, j))
		}
		grid = append(grid, row)
		i += 1
	}
	return grid
}

func part1(grid [][]*Tree) {
	rows, cols := len(grid), len(grid[0])

	// initialise DP table
	tallestFrom := map[string][][]*Tree{
		"LEFT":   make([][]*Tree, rows),
		"RIGHT":  make([][]*Tree, rows),
		"TOP":    make([][]*Tree, rows),
		"BOTTOM": make([][]*Tree, rows),
	}
	for _, d := range []string{"LEFT", "RIGHT", "TOP", "BOTTOM"} {
		for i := 0; i < rows; i++ {
			row := make([]*Tree, cols)
			for j := 0; j < cols; j++ {
				row[j] = NewTree(-1, i, j)
			}
			tallestFrom[d][i] = row
		}
	}
	// perform dynamic programming
	for i := 1; i < rows-1; i++ {
		for j := 1; j < cols-1; j++ {
			if tallestFrom["LEFT"][i][j-1].Height > grid[i][j-1].Height {
				tallestFrom["LEFT"][i][j] = tallestFrom["LEFT"][i][j-1]
			} else {
				tallestFrom["LEFT"][i][j] = grid[i][j-1]
			}
		}
	}
	for i := 1; i < rows-1; i++ {
		for j := cols - 2; j > 0; j-- {
			if tallestFrom["RIGHT"][i][j+1].Height > grid[i][j+1].Height {
				tallestFrom["RIGHT"][i][j] = tallestFrom["RIGHT"][i][j+1]
			} else {
				tallestFrom["RIGHT"][i][j] = grid[i][j+1]
			}
		}
	}
	for i := 1; i < rows-1; i++ {
		for j := 1; j < cols-1; j++ {
			if tallestFrom["TOP"][i-1][j].Height > grid[i-1][j].Height {
				tallestFrom["TOP"][i][j] = tallestFrom["TOP"][i-1][j]
			} else {
				tallestFrom["TOP"][i][j] = grid[i-1][j]
			}
		}
	}
	for i := rows - 2; i > 0; i-- {
		for j := 1; j < cols-1; j++ {
			if tallestFrom["BOTTOM"][i+1][j].Height > grid[i+1][j].Height {
				tallestFrom["BOTTOM"][i][j] = tallestFrom["BOTTOM"][i+1][j]
			} else {
				tallestFrom["BOTTOM"][i][j] = grid[i+1][j]
			}
		}
	}
	// start counting number of visible trees
	count := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j].Height > tallestFrom["TOP"][i][j].Height || grid[i][j].Height > tallestFrom["BOTTOM"][i][j].Height || grid[i][j].Height > tallestFrom["LEFT"][i][j].Height || grid[i][j].Height > tallestFrom["RIGHT"][i][j].Height {
				count += 1
			}
		}
	}
	//fmt.Println("number of visible trees:", count)
}

func part2(grid [][]*Tree) {
	rows, cols := len(grid), len(grid[0])
	stack := &Stack{items: make([]*Tree, 0)}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for stack.Length() > 0 {
				if stack.Peek().Height <= grid[i][j].Height {
					tree := stack.Pop()
					tree.See["RIGHT"] = j - tree.Col
				} else {
					break
				}
			}
			stack.Push(grid[i][j])
		}
		for stack.Length() > 0 {
			tree := stack.Pop()
			tree.See["RIGHT"] = cols - 1 - tree.Col
		}
	}

	for i := 0; i < rows; i++ {
		for j := cols - 1; j > -1; j-- {
			for stack.Length() > 0 {
				if stack.Peek().Height <= grid[i][j].Height {
					tree := stack.Pop()
					tree.See["LEFT"] = tree.Col - j
				} else {
					break
				}
			}
			stack.Push(grid[i][j])
		}
		for stack.Length() > 0 {
			tree := stack.Pop()
			tree.See["LEFT"] = tree.Col
		}
	}

	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			for stack.Length() > 0 {
				if stack.Peek().Height <= grid[i][j].Height {
					tree := stack.Pop()
					tree.See["BOTTOM"] = i - tree.Row
				} else {
					break
				}
			}
			stack.Push(grid[i][j])
		}
		for stack.Length() > 0 {
			tree := stack.Pop()
			tree.See["BOTTOM"] = rows - 1 - tree.Row
		}
	}

	for j := 0; j < cols; j++ {
		for i := rows - 1; i > -1; i-- {
			for stack.Length() > 0 {
				if stack.Peek().Height <= grid[i][j].Height {
					tree := stack.Pop()
					tree.See["TOP"] = tree.Row - i
				} else {
					break
				}
			}
			stack.Push(grid[i][j])
		}
		for stack.Length() > 0 {
			tree := stack.Pop()
			tree.See["TOP"] = tree.Row
		}
	}

	max_score := -1
	max_i := 0
	max_j := 0
	for i := 1; i < rows-1; i++ {
		for j := 1; j < cols-1; j++ {
			score := grid[i][j].Score()
			if score > max_score {
				max_score = score
				max_i = i
				max_j = j
			}
		}
	}
	if max_j+max_i == 0 {

	}
	//fmt.Println("max_score", max_score, "for tree", max_i, max_j)
}

func main() {
	grid := readFile("../input.txt")
	part1(grid)
	part2(grid)
}
