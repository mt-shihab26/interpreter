package ast

import "testing"

// stubNode is a minimal Node whose Tree() output is fixed, used to test
// renderTree's connector/indentation logic in isolation from any real AST
// node's own Tree() formatting.
type stubNode struct {
	tree string
}

func (s stubNode) TokenLiteral() string { return "" }
func (s stubNode) String() string       { return "" }
func (s stubNode) Tree() string         { return s.tree }

func TestRenderTreeNoChildren(t *testing.T) {
	if actual := renderTree("Header"); actual != "Header" {
		t.Errorf("renderTree(\"Header\") wrong, got=%q", actual)
	}
}

func TestRenderTreeSingleChild(t *testing.T) {
	actual := renderTree("Header", treeChild{"Only", stubNode{"Leaf"}})
	expected := "Header\n└─ Only: Leaf"
	if actual != expected {
		t.Errorf("renderTree() wrong.\nexpected=%q\ngot=%q", expected, actual)
	}
}

func TestRenderTreeMultipleChildren(t *testing.T) {
	actual := renderTree("Header",
		treeChild{"First", stubNode{"A"}},
		treeChild{"Second", stubNode{"B"}},
	)
	expected := "Header\n├─ First: A\n└─ Second: B"
	if actual != expected {
		t.Errorf("renderTree() wrong.\nexpected=%q\ngot=%q", expected, actual)
	}
}

func TestRenderTreeUnlabeledChild(t *testing.T) {
	actual := renderTree("Header", treeChild{"", stubNode{"Leaf"}})
	expected := "Header\n└─ Leaf"
	if actual != expected {
		t.Errorf("renderTree() wrong.\nexpected=%q\ngot=%q", expected, actual)
	}
}

// TestRenderTreeMultilineChild checks that when a child's own Tree() output
// spans multiple lines, only the first line gets the connector -- every
// following line is re-indented (not connected) so it nests correctly
// under the parent instead of looking like a sibling.
func TestRenderTreeMultilineChild(t *testing.T) {
	actual := renderTree("Header",
		treeChild{"Branch", stubNode{"Sub\n└─ Leaf"}},
		treeChild{"Last", stubNode{"End"}},
	)
	expected := "Header\n" +
		"├─ Branch: Sub\n" +
		"│  └─ Leaf\n" +
		"└─ Last: End"
	if actual != expected {
		t.Errorf("renderTree() wrong.\nexpected=%q\ngot=%q", expected, actual)
	}
}
