package markdown

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
)

func Hook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// switch v := node.(type) {
	// case *ast.CodeBlock:
	// 	CodeBlock(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.Code:
	// 	Code(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.Link:
	// 	Link(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.BlockQuote:
	// 	BlockQuote(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.List:
	// 	List(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.Heading:
	// 	Heading(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.Paragraph:
	// 	Paragraph(w, v, entering)
	// 	return ast.GoToNext, true
	// case *ast.HorizontalRule:
	// 	Hr(w, v, entering)
	// 	return ast.GoToNext, true
	// default:
	// 	return ast.GoToNext, false
	// }

	return ast.GoToNext, false
}

func CodeBlock(w io.Writer, p *ast.CodeBlock, entering bool) {
	if entering {
		io.WriteString(w, "<pre><code>")
	} else {
		io.WriteString(w, "</code></pre>")
	}
}

func Code(w io.Writer, p *ast.Code, entering bool) {
	if entering {
		io.WriteString(w, "<code class=\"rounded-md border border-gray-600 px-1.5 py-0.5\">")
	} else {
		io.WriteString(w, "</code>")
	}
}

func Link(w io.Writer, p *ast.Link, entering bool) {
	if entering {
		io.WriteString(w, fmt.Sprintf("<a class=\"underline hover:text-primary-500 text-primary-600\" href=\"%s\">", p.Destination))
	} else {
		io.WriteString(w, "</a>")
	}
}

func BlockQuote(w io.Writer, p *ast.BlockQuote, entering bool) {
	if entering {
		io.WriteString(w, "<blockquote class=\"border-l-4 rounded-sm border-gray-600 pl-2\">")
	} else {
		io.WriteString(w, "</blockquote>")
	}
}

func List(w io.Writer, p *ast.List, entering bool) {

	var listType string
	if p.ListFlags == ast.ListTypeOrdered {
		listType = "ol"
	} else {
		listType = "ul"
	}

	var listClass string

	switch p.ListFlags {
	case ast.ListTypeOrdered:
		listClass = "list-decimal"
	case ast.ListTypeDefinition:
		listClass = "list-disc"
	case ast.ListTypeTerm:
		listClass = "list-decimal"
	case ast.ListItemContainsBlock:
		listClass = "list-disc"
	case ast.ListItemBeginningOfList:
		listClass = "list-disc"
	}

	if entering {
		io.WriteString(w, fmt.Sprintf("<%s class=\"pl-4 %s\">", listType, listClass))
	} else {
		io.WriteString(w, fmt.Sprintf("</%s>", listType))
	}
}

func Heading(w io.Writer, p *ast.Heading, entering bool) {

	var textSize string

	switch p.Level {
	case 1:
		textSize = "text-2xl"
	case 2:
		textSize = "text-xl"
	case 3:
		textSize = "text-lg"
	case 4:
		textSize = "text-base"
	case 5:
		textSize = "text-sm"
	case 6:
		textSize = "text-xs"
	}

	if entering {
		io.WriteString(w, "<h"+fmt.Sprint(p.Level)+" class=\"px-2 my-4 font-bold "+textSize+"\">")
	} else {
		io.WriteString(w, "<h"+fmt.Sprint(p.Level)+"/>")
	}
}

func Hr(w io.Writer, p *ast.HorizontalRule, entering bool) {
	if entering {
		io.WriteString(w, "<hr class=\"my-4\"/>")
	}
}

func Paragraph(w io.Writer, p *ast.Paragraph, entering bool) {
	if entering {
		io.WriteString(w, "<p class=\"\">")
	} else {
		io.WriteString(w, "</p>")
	}
}
