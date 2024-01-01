package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(writer io.Writer, path string, isPrintFiles bool) error {
	var (
		resultTree string
		err        error
	)

	if isPrintFiles {
		resultTree, err = dirTreeWithFiles(path, "")
		if err != nil {
			return err
		}
	} else {
		resultTree, err = dirTreeWithoutFiles(path, "")
		if err != nil {
			return err
		}
	}

	writer.Write([]byte(resultTree))

	return nil
}

func dirTreeWithFiles(path, prefix string) (string, error) {
	info, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer info.Close()

	dirEntries, err := info.ReadDir(0)
	if err != nil {
		return "", err
	}

	sort.Slice(dirEntries, func(i, j int) bool {
		return dirEntries[i].Name() < dirEntries[j].Name()
	})

	currentDirTree := ""
	prevPref := prefix

	for i, entry := range dirEntries {
		name := entry.Name()
		if entry.Name() == ".idea" || entry.Name() == ".DS_Store" {
			continue
		}

		isDir := entry.IsDir()

		if !isDir {
			fileInfo, err := entry.Info()
			if err != nil {
				return "", err
			}

			size := "empty"

			if fileInfo.Size() != 0 {
				size = fmt.Sprintf("%db", fileInfo.Size())
			}
			name = fmt.Sprintf("%s (%s)", name, size)
		}

		if i == len(dirEntries)-1 {
			currentDirTree += fmt.Sprintf("%s└───%s\n", prevPref, name)
			prefix = "\t"
		} else {
			currentDirTree += fmt.Sprintf("%s├───%s\n", prevPref, name)
			prefix = "│\t"
		}

		if isDir {
			dir, err := dirTreeWithFiles(fmt.Sprintf("%s/%s", path, name), prevPref+prefix)
			if err != nil {
				return "", err
			}

			currentDirTree += dir
		}
	}

	return currentDirTree, nil
}

func dirTreeWithoutFiles(path, prefix string) (string, error) {
	info, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer info.Close()

	dirEntries, err := info.ReadDir(0)
	if err != nil {
		return "", err
	}

	dirs := make([]string, 0, 0)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			if entry.Name() == ".idea" || entry.Name() == ".DS_Store" {
				continue
			}

			dirs = append(dirs, entry.Name())
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i] < dirs[j]
	})

	currentDirTree := ""
	prevPref := prefix

	for i, dirName := range dirs {
		if i == len(dirs)-1 {
			currentDirTree += fmt.Sprintf("%s└───%s\n", prevPref, dirName)
			prefix = "\t"
		} else {
			currentDirTree += fmt.Sprintf("%s├───%s\n", prevPref, dirName)
			prefix = "│\t"
		}

		dir, err := dirTreeWithoutFiles(fmt.Sprintf("%s/%s", path, dirName), prevPref+prefix)
		if err != nil {
			return "", err
		}

		currentDirTree += dir
	}

	return currentDirTree, nil
}
