// Package html provides extensions to Go's standard html package
package html

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Element represents an HTML element with attributes and children
type Element struct {
	Tag        string
	Attributes map[string]string
	Children   []interface{} // Can be *Element or string
}

// NewElement creates a new HTML element with the given tag
func NewElement(tag string) *Element {
	return &Element{
		Tag:        tag,
		Attributes: make(map[string]string),
		Children:   make([]interface{}, 0),
	}
}

// Attr sets an attribute on the element and returns the element for chaining
func (e *Element) Attr(key, value string) *Element {
	e.Attributes[key] = value
	return e
}

// Class adds a CSS class to the element and returns the element for chaining
func (e *Element) Class(class string) *Element {
	if existing, ok := e.Attributes["class"]; ok {
		e.Attributes["class"] = existing + " " + class
	} else {
		e.Attributes["class"] = class
	}
	return e
}

// ID sets the ID of the element and returns the element for chaining
func (e *Element) ID(id string) *Element {
	e.Attributes["id"] = id
	return e
}

// AppendChild adds a child element or text node and returns the parent for chaining
func (e *Element) AppendChild(child interface{}) *Element {
	e.Children = append(e.Children, child)
	return e
}

// Text adds a text node as a child and returns the element for chaining
func (e *Element) Text(text string) *Element {
	return e.AppendChild(text)
}

// String renders the element and its children as an HTML string
func (e *Element) String() string {
	var buf bytes.Buffer
	e.render(&buf)
	return buf.String()
}

// render writes the HTML representation of the element to the given writer
func (e *Element) render(w io.Writer) {
	w.Write([]byte("<" + e.Tag))
	for key, value := range e.Attributes {
		w.Write([]byte(" " + key + "=\"" + value + "\""))
	}

	if len(e.Children) == 0 && isVoidElement(e.Tag) {
		w.Write([]byte(">"))
		return
	}

	w.Write([]byte(">"))

	for _, child := range e.Children {
		switch c := child.(type) {
		case *Element:
			c.render(w)
		case string:
			w.Write([]byte(c))
		}
	}

	w.Write([]byte("</" + e.Tag + ">"))
}

// isVoidElement returns true if the tag is an HTML void element
func isVoidElement(tag string) bool {
	switch strings.ToLower(tag) {
	case "area", "base", "br", "col", "embed", "hr", "img", "input",
		"link", "meta", "param", "source", "track", "wbr":
		return true
	default:
		return false
	}
}

// FindElementByID finds an element with the given ID in a parsed HTML node tree
func FindElementByID(node *html.Node, id string) *html.Node {
	if node == nil {
		return nil
	}

	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == id {
				return node
			}
		}
	}

	// Search children
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result := FindElementByID(child, id); result != nil {
			return result
		}
	}

	return nil
}

// FindElementsByClass finds all elements with the given class in a parsed HTML node tree
func FindElementsByClass(node *html.Node, class string) []*html.Node {
	var results []*html.Node
	findElementsByClass(node, class, &results)
	return results
}

// findElementsByClass is a helper for FindElementsByClass
func findElementsByClass(node *html.Node, class string, results *[]*html.Node) {
	if node == nil {
		return
	}

	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				classes := strings.Fields(attr.Val)
				for _, c := range classes {
					if c == class {
						*results = append(*results, node)
						break
					}
				}
			}
		}
	}

	// Search children
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findElementsByClass(child, class, results)
	}
}

// ExtractText extracts all text content from an HTML node and its descendants
func ExtractText(node *html.Node) string {
	var buf bytes.Buffer
	extractText(node, &buf)
	return buf.String()
}

// extractText is a helper for ExtractText
func extractText(node *html.Node, buf *bytes.Buffer) {
	if node == nil {
		return
	}

	if node.Type == html.TextNode {
		buf.WriteString(node.Data)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		extractText(child, buf)
	}
}

// NodeToString converts an html.Node back to its string representation
func NodeToString(node *html.Node) (string, error) {
	var buf bytes.Buffer
	err := html.Render(&buf, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
