package initcmd

import (
	"path/filepath"
	"sort"
	"strings"
)

type previewNode struct {
	name     string
	children map[string]*previewNode
	file     bool
}

// FormatPreviewTree renders relative file paths as a tree suitable for terminal output.
func FormatPreviewTree(paths []string) []string {
	root := &previewNode{
		children: map[string]*previewNode{},
	}

	for _, path := range paths {
		current := root

		for _, part := range strings.Split(filepath.ToSlash(path), "/") {
			child, ok := current.children[part]
			if !ok {
				child = &previewNode{
					name:     part,
					children: map[string]*previewNode{},
				}
				current.children[part] = child
			}

			current = child
		}

		current.file = true
	}

	var lines []string
	for _, child := range sortedChildren(root) {
		if child.file && len(child.children) == 0 {
			lines = append(lines, child.name)
			continue
		}

		lines = append(lines, child.name+"/")
		renderPreviewChildren(child, "", &lines)
	}

	return lines
}

func renderPreviewChildren(node *previewNode, prefix string, lines *[]string) {
	children := sortedChildren(node)

	for index, child := range children {
		branch := "├── "
		nextPrefix := prefix + "│   "

		if index == len(children)-1 {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		if child.file && len(child.children) == 0 {
			*lines = append(*lines, prefix+branch+child.name)
			continue
		}

		*lines = append(*lines, prefix+branch+child.name+"/")
		renderPreviewChildren(child, nextPrefix, lines)
	}
}

func sortedChildren(node *previewNode) []*previewNode {
	names := make([]string, 0, len(node.children))
	for name := range node.children {
		names = append(names, name)
	}

	sort.Slice(names, func(left, right int) bool {
		leftNode := node.children[names[left]]
		rightNode := node.children[names[right]]

		if leftNode.file != rightNode.file {
			return !leftNode.file
		}

		return names[left] < names[right]
	})

	children := make([]*previewNode, 0, len(names))
	for _, name := range names {
		children = append(children, node.children[name])
	}

	return children
}
