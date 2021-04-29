package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	_ "path/filepath"
	"strings"
)

type Node struct {
	File     os.DirEntry
	Children []Node
}

func (node Node) GetName() string {
	if node.File.IsDir() {
		return node.File.Name()
	} else {
		return fmt.Sprintf("%s (%s)", node.File.Name(), node.GetSize())
	}
}

func (node Node) GetSize() string {
	info, err := node.File.Info()
	if err == nil && info.Size() > 0 {
		return fmt.Sprintf("%db", info.Size())
	}
	return "empty"
}

func printTree(writer io.Writer, tree []Node, parentPrefix string) {
	var (
		prefix      = "├───"
		childPrefix = "│\t"
	)

	for index, node := range tree {
		if index == len(tree)-1 {
			prefix = "└───"
			childPrefix = "\t"
		}

		line := strings.Join([]string{parentPrefix, prefix, node.GetName(), "\n"}, "")
		writer.Write([]byte(line))

		if node.File.IsDir() {
			printTree(writer, node.Children, parentPrefix+childPrefix)
		}
	}
}

func buildTree(path string, printFiles bool) ([]Node, error) {
	var tree []Node
	dir, err := os.ReadDir(path)

	for _, entry := range dir {
		if !printFiles && !entry.IsDir() {
			continue
		}
		node := Node{File: entry}
		if entry.IsDir() {
			children, err := buildTree(filepath.Join(path, entry.Name()), printFiles)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}
		tree = append(tree, node)
	}

	return tree, err
}

func dirTree(writer io.Writer, path string, printFiles bool) error {
	tree, err := buildTree(path, printFiles)
	printTree(writer, tree, "")
	return err
}

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
