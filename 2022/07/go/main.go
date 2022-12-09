package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"strconv"
	"strings"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

var diskSpace = 70_000_000
var updateNeeds = 30_000_000

type directory struct {
	path             string
	pathComponents   []string
	parent           string
	parentComponents []string
	totalSize        int

	// Not sure if we'll get duplicate names, better track with file names.
	//fileSizes   map[string]int
	files       map[string]int
	directories map[string]*directory
}

type filesHashEntry struct {
	dirPath string
	size    int
}

func a() {
	scanner := common.GetLineScanner(fileName)

	// Working directory.
	wd := ""
	root := directory{
		path:             "",
		pathComponents:   []string{""},
		parent:           "",
		parentComponents: []string{""},
		files:            map[string]int{},
		directories:      map[string]*directory{},
		totalSize:        0,
	}
	currentDirectory := &root

	// Flat hash of directories to easily iterate through without walking the tree.
	directoryHash := map[string]*directory{
		"": &root,
	}

	// Flat hash of files to easily iterate through files in case there's a need.
	filesHash := map[string]filesHashEntry{}

	createNewDirectory := func(name string, parent *directory) *directory {
		if parent == nil {
			panic("[createNewDirectory] Parent should not be nil")
		}
		newDirectory := directory{
			path:             wd + "/" + name,
			pathComponents:   append(parent.pathComponents, name),
			parent:           wd,
			parentComponents: parent.pathComponents,
			files:            map[string]int{},
			directories:      map[string]*directory{},
			totalSize:        0,
		}
		currentDirectory.directories[name] = &newDirectory
		return &newDirectory
	}

	// Scan structure
	for scanner.Scan() {
		l := scanner.Text()

		parts := strings.Split(l, " ")

		if parts[0] == "$" {
			// Command

			// Ignore `ls` command.
			if parts[1] != "cd" {
				continue
			}

			// Change directory
			if parts[2] == "/" {
				// Change to root.
				wd = ""
				currentDirectory = &root
			} else if parts[2] == ".." {
				// Change to parent
				wd = currentDirectory.parent
				currentDirectory = directoryHash[wd]
			} else {
				dirName := parts[2]
				wd = wd + "/" + dirName

				pointer, exists := directoryHash[wd]
				if exists {
					currentDirectory = pointer
				} else {
					//fmt.Printf("Path %s does not exist yet\n", wd)
					newDirectory := createNewDirectory(dirName, currentDirectory)
					directoryHash[wd] = newDirectory
					currentDirectory = newDirectory
				}
			}
		} else if parts[0] == "dir" {
			// Directory, do nothing.
		} else {
			// File
			fileName := parts[1]
			size, err := strconv.Atoi(parts[0])
			if err != nil {
				panic("Could not parse file size.")
			}
			currentDirectory.files[fileName] = size
			filesHash[wd+"/"+fileName] = filesHashEntry{
				dirPath: wd,
				size:    size,
			}
		}

		//fmt.Printf("l: %s\n", l)
	}

	logFilesHash(filesHash)
	// Calculate directory sizes.
	totalSize := calculateSizeForDirectory(&root)
	fmt.Printf("totalSize: %d\n", totalSize)

	sizeOfDirectoriesBelow100k := 0
	for _, d := range directoryHash {
		//fmt.Printf("%s: %d\n", dirPath, d.totalSize)
		if d.totalSize < 100000 {
			sizeOfDirectoriesBelow100k += d.totalSize
		}
	}

	fmt.Printf("Part1: %d\n", sizeOfDirectoriesBelow100k)
	//fmt.Printf("Part1 done\n")

	sizeToDeleteToBeEnough := updateNeeds - (diskSpace - root.totalSize)
	fmt.Printf("sizeToDeleteToBeEnough: %d\n", sizeToDeleteToBeEnough)

	closestToSizeToDelete := diskSpace
	for _, d := range directoryHash {
		//fmt.Printf("%s: %d\n", dirPath, d.totalSize)
		if d.totalSize < sizeToDeleteToBeEnough {
			continue
		}
		if d.totalSize-sizeToDeleteToBeEnough < closestToSizeToDelete-sizeToDeleteToBeEnough {
			closestToSizeToDelete = d.totalSize
		}
	}
	fmt.Printf("Part2: %d\n", closestToSizeToDelete)
}

// Walk the tree and calculate total sizes.
func calculateSizeForDirectory(d *directory) int {
	if d == nil {
		panic("Directory cannot be nil")
	}

	totalFiles := 0
	for _, size := range d.files {
		totalFiles += size
	}

	totalDirectory := 0
	for _, childDirectory := range d.directories {
		totalDirectory += calculateSizeForDirectory(childDirectory)
	}

	d.totalSize = totalFiles + totalDirectory
	return d.totalSize
}

// region Debug helpers
func logFilesHash(hash map[string]filesHashEntry) {
	fmt.Println("========Logging files hash")
	for path, d := range hash {
		fmt.Printf("%s: %d %s\n", path, d.size, d.dirPath)
	}
	fmt.Println("========End Logging files hash")
}

// endregion Debug helpers

func main() {
	a()
}
