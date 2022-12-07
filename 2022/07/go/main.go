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
	path           string
	parent         string
	pathComponents []string
	totalSize      int

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
		path:           "",
		parent:         "",
		pathComponents: []string{""},
		files:          map[string]int{},
		directories:    map[string]*directory{},
		totalSize:      0,
	}
	currentDirectory := &root

	directoryHash := map[string]*directory{
		"": &root,
	}

	filesHash := map[string]filesHashEntry{}
	directorySizes := map[string]int{}

	createNewDirectory := func(name string, parent *directory) *directory {
		if parent == nil {
			panic("[createNewDirectory] Parent should not be nil")
		}
		newDirectory := directory{
			path:           wd + "/" + name,
			parent:         wd,
			pathComponents: append(parent.pathComponents, name),
			files:          map[string]int{},
			directories:    map[string]*directory{},
			totalSize:      0,
		}
		currentDirectory.directories[name] = &newDirectory
		return &newDirectory
	}

	// Scan structure
	for scanner.Scan() {
		l := scanner.Text()
		//fmt.Printf("l: %s\n", l)

		parts := strings.Split(l, " ")

		if parts[0] == "$" {
			// Command
			if parts[1] == "cd" {
				// Change directory
				if parts[2] == "/" {
					// Change to root.
					wd = ""
					currentDirectory = &root
				} else if parts[2] == ".." {
					// Change to parent
					wdComponents := strings.Split(wd, "/")
					newPathComponents := wdComponents[:len(wdComponents)-1]
					wd = strings.Join(newPathComponents, "/")
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
					//logDirectoryHash(directoryHash)
				}
			}
			// We can ignore 'ls'

		} else if parts[0] == "dir" {
			// Directory
			dirName := parts[1]
			newDirectory := createNewDirectory(dirName, currentDirectory)
			directoryHash[wd+"/"+dirName] = newDirectory
			//logDirectoryHash(directoryHash)
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
	for _, entry := range filesHash {
		directorySizes[entry.dirPath] += entry.size
	}
	for dir, size := range directorySizes {
		fmt.Printf("%s: %d\n", dir, size)
	}

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

func logDirectoryHash(directoryHash map[string]*directory) {
	fmt.Println("========Logging directory hash")
	for path, d := range directoryHash {
		fmt.Printf("%s: %+v\n", path, d)
	}
	fmt.Println("========End Logging directory hash")
}

func logFilesHash(hash map[string]filesHashEntry) {
	fmt.Println("========Logging files hash")
	for path, d := range hash {
		fmt.Printf("%s: %d %s\n", path, d.size, d.dirPath)
	}
	fmt.Println("========End Logging files hash")
}

func main() {
	a()
}
