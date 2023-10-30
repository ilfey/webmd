package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/kyoto-framework/zen/v2"
	"github.com/rotisserie/eris"

	"github.com/alecthomas/chroma"
	chtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const ()

var (
	htmlFormatter  *chtml.Formatter
	highlightStyle *chroma.Style
	scipChars      = []string{"/", "\\", "`", "*", "\"", "'", ":", ";", "&", "<", ">", "|", "!", "?", "(", ")", "[", "]", "{", "}"}
)

func init() {
	htmlFormatter = chtml.New(chtml.WithClasses(true), chtml.TabWidth(2), chtml.WithLineNumbers(true))
	if htmlFormatter == nil {
		panic(eris.Errorf("couldn't create html formatter"))
	}
	styleName := "borland"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {

		panic(eris.Errorf("didn't find style '%s'", styleName))
	}
}

func RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {

	switch v := node.(type) {
	// case *ast.Text:
	// case *ast.Softbreak:
	// case *ast.Hardbreak:
	// case *ast.NonBlockingSpace:
	// case *ast.Emph:
	// case *ast.Strong:
	// case *ast.Del:
	case *ast.BlockQuote:
		BlockQuote(w, v, entering)
	// case *ast.Aside:
	case *ast.Link:
		Link(w, v, entering)
	case *ast.CrossReference:
		link := &ast.Link{Destination: append([]byte("#"), v.Destination...)}
		Link(w, link, entering)
	case *ast.Citation:
		// case *ast.Image:
	case *ast.Code:
		Code(w, v, entering)
	case *ast.CodeBlock:
		CodeBlock(w, v, entering)
	case *ast.Caption:
		Caption(w, v, entering)
	// case *ast.CaptionFigure:
	// case *ast.Document:
	case *ast.Paragraph:
		Paragraph(w, v, entering)
	// case *ast.HTMLSpan:
	// case *ast.HTMLBlock:
	case *ast.Heading:
		Heading(w, v, entering)
	case *ast.HorizontalRule:
		Hr(w, v, entering)
	case *ast.List:
		List(w, v, entering)
	case *ast.ListItem:
		ListItem(w, v, entering)
	case *ast.Table:
		Table(w, v, entering)
	case *ast.TableCell:
		TableCell(w, v, entering)

	case *ast.TableHeader:
		TableHeader(w, v, entering)
	case *ast.TableBody:
		TableBody(w, v, entering)
	case *ast.TableRow:
		TableRow(w, v, entering)
	// case *ast.TableFooter:
	// case *ast.Math:
	// case *ast.MathBlock:
	// case *ast.DocumentMatter:
	// case *ast.Callout:
	// case *ast.Index:
	// case *ast.Subscript:
	// case *ast.Superscript:
	// case *ast.Footnotes:
	default:
		return ast.GoToNext, false
	}

	return ast.GoToNext, true
}

func Caption(w io.Writer, p *ast.Caption, entering bool) {
	if entering {
		io.WriteString(w, "<figcaption class=\"text-xs text-gray-700 dark:text-gray-400\">")
	} else {
		io.WriteString(w, "</figcaption>")
	}
}

func TableBody(w io.Writer, p *ast.TableBody, entering bool) {
	if entering {
		io.WriteString(w, "<tbody>")
	} else {
		io.WriteString(w, "</tbody>")
	}
}

func TableRow(w io.Writer, p *ast.TableRow, entering bool) {
	if _, ok := p.Parent.(*ast.TableHeader); ok {
		if entering {
			io.WriteString(w, "<tr>")
		} else {
			io.WriteString(w, "</tr>")
		}

		return
	}

	if entering {
		io.WriteString(w, "<tr class=\"bg-gray-200/50 border-b dark:bg-gray-900/10 border-gray-500/50\">")
	} else {
		io.WriteString(w, "</tr>")
	}
}

func TableHeader(w io.Writer, p *ast.TableHeader, entering bool) {
	if entering {
		io.WriteString(w, "<thead class=\"text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400\">")
	} else {
		io.WriteString(w, "</thead>")
	}
}

func TableCell(w io.Writer, p *ast.TableCell, entering bool) {
	if !entering {
		if p.IsHeader {
			io.WriteString(w, "</th>")
		} else {
			io.WriteString(w, "</td>")
		}

		return
	}

	var attrs []string
	openTag := "<td class=\"px-6 py-4\""
	if p.IsHeader {
		openTag = "<th class=\"px-6 py-3\""
	}

	align := p.Align.String()
	if align != "" {
		attrs = append(attrs, fmt.Sprintf(`align="%s"`, align))
	}

	if colspan := p.ColSpan; colspan > 0 {
		attrs = append(attrs, fmt.Sprintf(`colspan="%d"`, colspan))
	}

	io.WriteString(w, fmt.Sprintf("%s %s>", openTag, strings.Join(attrs, " ")))

}

func Table(w io.Writer, p *ast.Table, entering bool) {
	if entering {
		io.WriteString(w, "<table class=\"my-4 table-auto\">")
	} else {
		io.WriteString(w, "</table>")
	}
}

func renderCode(w io.Writer, p *ast.CodeBlock, entering bool) error {
	source := string(p.Literal)
	lang := string(p.Info)

	if lang == "" {
		html.EscapeHTML(w, p.Literal)

		return nil
	}

	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}

	if l == nil {
		l = lexers.Fallback
	}

	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return htmlFormatter.Format(w, highlightStyle, it)
}

func CodeBlock(w io.Writer, p *ast.CodeBlock, entering bool) {
	io.WriteString(w, "<pre class=\"my-2 rounded-md bg-gray-300/50 dark:bg-gray-700/50 border border-gray-500/50 px-1.5 py-0.5\"><code>")
	err := renderCode(w, p, entering)
	if err != nil {
		// TODO: handle error
		fmt.Println(err)

		html.EscapeHTML(w, p.Literal)
	}

	io.WriteString(w, "</code></pre>")
}

func Code(w io.Writer, p *ast.Code, entering bool) {
	io.WriteString(w, "<code class=\"rounded-md bg-gray-300/50 dark:bg-gray-700/50 border border-gray-500/50 px-1.5 py-0.5\">")
	html.EscapeHTML(w, p.Literal)
	io.WriteString(w, "</code>")
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
		io.WriteString(w, "<blockquote class=\"my-2 border-l-4 rounded-sm border-gray-500/50 pl-2\">")
	} else {
		io.WriteString(w, "</blockquote>")
	}
}

func ListItem(w io.Writer, p *ast.ListItem, entering bool) {
	if p.RefLink != nil {
		slug := html.Slugify(p.RefLink)
		io.WriteString(w, html.FootnoteItem("[^", slug))

		return
	}

	if entering {
		tag := "<li>"
		if p.ListFlags&ast.ListTypeDefinition != 0 {
			tag = "<dd class=\"ml-4\">"
		}

		if p.ListFlags&ast.ListTypeTerm != 0 {
			tag = "<dt>"
		}

		io.WriteString(w, tag)

	} else {
		tag := "</li>"
		if p.ListFlags&ast.ListTypeDefinition != 0 {
			tag = "</dd>"
		}

		if p.ListFlags&ast.ListTypeTerm != 0 {
			tag = "</dt>"
		}

		io.WriteString(w, tag)
	}
}

func List(w io.Writer, p *ast.List, entering bool) {
	var attrs []string
	var classes []string
	tag := "ul"

	if p.ListFlags&ast.ListTypeOrdered != 0 {
		tag = "ol"
		classes = append(classes, "list-disc")
	} else if p.ListFlags&ast.ListTypeDefinition != 0 {
		tag = "dl"
		classes = append(classes, "list-disc")
	} else {
		classes = append(classes, "list-decimal")
	}

	if tag == "ul" {
		if p.Start > 0 {
			attrs = append(attrs, fmt.Sprintf(`start="%d"`, p.Start))
		}
	}

	if _, ok := p.Parent.(*ast.ListItem); ok {
		classes = append(classes, "my-0.5")
	} else {
		classes = append(classes, "my-2")
	}

	if entering {
		io.WriteString(w, fmt.Sprintf("<%s class=\"pl-4 %s\" %s>", tag, strings.Join(classes, " "), strings.Join(attrs, " ")))
	} else {
		io.WriteString(w, fmt.Sprintf("</%s>", tag))
	}
}

func formatID(id string) (result string) {
	for _, c := range strings.Split(strings.ToLower(id), "") {
		if zen.In(c, scipChars) {
			continue
		}

		result += c
	}

	return
}

func Heading(w io.Writer, p *ast.Heading, entering bool) {
	var textSize string
	var attrs []string

	attrs = append(attrs, fmt.Sprintf("id=\"%s\"", formatID(p.HeadingID)))

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
		io.WriteString(w, fmt.Sprintf("<h%d class=\"px-2 my-4 font-bold %s\" %s>", p.Level, textSize, strings.Join(attrs, " ")))
	} else {
		io.WriteString(w, fmt.Sprintf("</h%d>", p.Level))
	}
}

func Hr(w io.Writer, p *ast.HorizontalRule, entering bool) {
	size := "h-0.5"

	if zen.All(p.Leaf.Literal, func(c byte) bool {
		return c == '-'
	}) {
		size = "h-1"
	} else if zen.All(p.Leaf.Literal, func(c byte) bool {
		return c == '*'
	}) {
		size = "h-1.5"
	}

	io.WriteString(w, fmt.Sprintf("<hr class=\"my-4 rounded-full %s bg-gray-500/50\"/>", size))
}

func Paragraph(w io.Writer, p *ast.Paragraph, entering bool) {
	if entering {
		io.WriteString(w, "<p class=\"\">")
	} else {
		io.WriteString(w, "</p>")
	}
}
