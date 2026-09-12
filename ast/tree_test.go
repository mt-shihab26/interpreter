package ast

import "testing"

// stubNode is a minimal Node with a fixed Tree() output, used to test renderTree in isolation from any real node's formatting.
type stubNode struct {
	tree string
}

// TokenLiteral always returns "" -- stubNode only needs to satisfy the Node interface.
func (s stubNode) TokenLiteral() string { return "" }

// String always returns "" -- stubNode only needs to satisfy the Node interface.
func (s stubNode) String() string { return "" }

// Tree returns the fixed tree string the stubNode was constructed with.
func (s stubNode) Tree() string { return s.tree }

// TestRenderTreeNoChildren checks that a header with no children renders as just that header line.
func TestRenderTreeNoChildren(t *testing.T) {
	if actual := renderTree("Header"); actual != "Header" {
		t.Errorf("renderTree(\"Header\") wrong, got=%q", actual)
	}
}

// TestRenderTreeSingleChild checks that a single child is rendered with a "└─" connector.
func TestRenderTreeSingleChild(t *testing.T) {
	actual := renderTree("Header", treeChild{"Only", stubNode{"Leaf"}})
	expected := "Header\n└─ Only: Leaf"
	if actual != expected {
		t.Errorf("renderTree() wrong.\nexpected=%q\ngot=%q", expected, actual)
	}
}

// TestRenderTreeMultipleChildren checks that non-last children get a "├─" connector and the last gets "└─".
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

// TestRenderTreeUnlabeledChild checks that a child with an empty label is rendered without a "label: " prefix.
func TestRenderTreeUnlabeledChild(t *testing.T) {
	actual := renderTree("Header", treeChild{"", stubNode{"Leaf"}})
	expected := "Header\n└─ Leaf"
	if actual != expected {
		t.Errorf("renderTree() wrong.\nexpected=%q\ngot=%q", expected, actual)
	}
}

// TestRenderTreeMultilineChild checks that only a multiline child's first line gets a connector, with later lines re-indented instead.
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
