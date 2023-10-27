package markdown

import (
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/ilfey/webmd/internal/fstree"
	"github.com/kyoto-framework/zen/v2"
)

func GetNav(node ast.Node) []*fstree.Link {

	var headings []*ast.Heading
	getHeadings(node, &headings)

	return zen.Map(headings, func(h *ast.Heading) *fstree.Link {

		var literals [][]byte
		getLiteral(h, &literals)

		text := strings.Join(zen.Map(literals, func(l []byte) string {
			return string(l)
		}), "")

		return fstree.NewLink(text, "#"+formatID(h.HeadingID))
	})
}

func getHeadings(node ast.Node, headings *[]*ast.Heading) {
	if v, ok := node.(*ast.Heading); ok {
		*headings = append(*headings, v)
	}

	for _, n := range node.GetChildren() {
		getHeadings(n, headings)
	}
}

func getLiteral(node ast.Node, literals *[][]byte) {
	children := node.GetChildren()

	if len(children) == 0 {
		*literals = append(*literals, node.AsLeaf().Literal)
	}

	for _, child := range children {
		getLiteral(child, literals)
	}
}
