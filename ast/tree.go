package ast

import "strings"

// treeChild is one labeled outgoing edge in a node's Tree() rendering.
type treeChild struct {
	label string
	node  Node
}

// renderTree builds a node's Tree() output: the node's own header line,
// followed by each child's own Tree() output nested underneath with
// box-drawing connectors ("├─ "/"└─ "). A child's continuation lines are
// re-indented so multi-line children still nest correctly under their
// parent.
func renderTree(header string, children ...treeChild) string {
	var out strings.Builder
	out.WriteString(header)
	for i, child := range children {
		connector, indent := "├─ ", "│  "
		if i == len(children)-1 {
			connector, indent = "└─ ", "   "
		}
		lines := strings.Split(child.node.Tree(), "\n")
		out.WriteString("\n")
		out.WriteString(connector)
		if child.label != "" {
			out.WriteString(child.label)
			out.WriteString(": ")
		}
		out.WriteString(lines[0])
		for _, line := range lines[1:] {
			out.WriteString("\n")
			out.WriteString(indent)
			out.WriteString(line)
		}
	}
	return out.String()
}
