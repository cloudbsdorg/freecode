package template

import (
	"fmt"
	"regexp"
	"strings"
)

type ElementType string

const (
	TypeWindow  ElementType = "window"
	TypeVBox    ElementType = "vbox"
	TypeHBox    ElementType = "hbox"
	TypeGrid    ElementType = "grid"
	TypeText    ElementType = "text"
	TypeList    ElementType = "list"
	TypeButton  ElementType = "button"
	TypeInput   ElementType = "input"
	TypeSpacer  ElementType = "spacer"
	TypeDivider ElementType = "divider"
	TypeImage   ElementType = "image"
	TypeProgress ElementType = "progress"
)

type Element struct {
	Type           ElementType
	ID             string
	Attributes     map[string]string
	AttributeVars  map[string][]string
	Children       []*Element
	Content        string
	ContentIsVar   bool
	ContentVarName string
}

type Template struct {
	Root    *Element
	Vars    map[string]interface{}
}

type Parser struct {
	pos int
	src string
}

func Parse(src string) (*Template, error) {
	p := &Parser{src: src, pos: 0}
	elements, err := p.parseElements()
	if err != nil {
		return nil, err
	}

	var root *Element
	if len(elements) == 1 {
		root = elements[0]
	} else {
		root = &Element{
			Type:     TypeVBox,
			Children: elements,
		}
	}

	return &Template{
		Root: root,
		Vars: make(map[string]interface{}),
	}, nil
}

func (p *Parser) parseElements() ([]*Element, error) {
	var elements []*Element
	for p.pos < len(p.src) {
		p.skipWhitespace()
		if p.pos >= len(p.src) {
			break
		}

		if p.match("${") {
			elem, err := p.parseVariable()
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
			continue
		}

		if !p.match("<") {
			p.skipToNextTag()
			continue
		}

		elem, err := p.parseElement()
		if err != nil {
			return nil, err
		}
		if elem != nil {
			elements = append(elements, elem)
		}
	}
	return elements, nil
}

func (p *Parser) parseElement() (*Element, error) {
	name := p.parseName()
	if name == "" {
		p.skipToNextTag()
		return nil, nil
	}

	elem := &Element{
		Type:          ElementType(name),
		Attributes:    make(map[string]string),
		AttributeVars: make(map[string][]string),
	}

	for {
		p.skipWhitespace()
		if p.peek() == '>' || p.pos >= len(p.src) {
			break
		}
		if p.match("/>") {
			return elem, nil
		}
		if p.peek() == '<' {
			break
		}

		attrName := p.parseName()
		if attrName == "" {
			break
		}
		if attrName == "id" {
			p.skipWhitespace()
			if p.match("=") {
				elem.ID = p.parseString()
			}
			continue
		}
		p.skipWhitespace()
		if !p.match("=") {
			break
		}
		attrVal := p.parseString()
		elem.Attributes[attrName] = attrVal
		if strings.Contains(attrVal, "${") {
			p.extractAttributeVars(elem, attrName, attrVal)
		}
	}

	if p.match(">") {
		content, children, err := p.parseContent()
		if err != nil {
			return nil, err
		}
		elem.Content = content
		elem.Children = children

		if p.match("</") {
			p.parseName()
			p.match(">")
		}
	}

	return elem, nil
}

func (p *Parser) parseContent() (string, []*Element, error) {
	var content strings.Builder
	var children []*Element

	for p.pos < len(p.src) {
		if p.match("</") {
			p.parseName()
			p.match(">")
			break
		}

		if p.match("${") {
			elem, err := p.parseVariable()
			if err != nil {
				return "", nil, err
			}
			children = append(children, elem)
			continue
		}

		if p.match("<") {
			elem, err := p.parseElement()
			if err != nil {
				return "", nil, err
			}
			if elem != nil {
				children = append(children, elem)
			}
			continue
		}

		content.WriteByte(p.src[p.pos])
		p.pos++
	}

	return content.String(), children, nil
}

func (p *Parser) parseVariable() (*Element, error) {
	var name strings.Builder
	for p.pos < len(p.src) && p.src[p.pos] != '}' && p.src[p.pos] != '.' {
		name.WriteByte(p.src[p.pos])
		p.pos++
	}

	if p.match("}") {
		return &Element{
			Type:    "var",
			Content: name.String(),
		}, nil
	}

	if p.match(".") {
		prop := p.parseName()
		if p.match("}") {
			return &Element{
				Type:    "var",
				Content: name.String() + "." + prop,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid variable syntax at position %d", p.pos)
}

func (p *Parser) parseName() string {
	var name strings.Builder
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' {
			name.WriteByte(c)
			p.pos++
		} else {
			break
		}
	}
	return name.String()
}

func (p *Parser) parseString() string {
	var quote byte
	if p.peek() == '"' || p.peek() == '\'' {
		quote = p.src[p.pos]
		p.pos++
	}

	var val strings.Builder
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == quote {
			p.pos++
			break
		}
		if quote == 0 && (c == ' ' || c == '\t' || c == '\n' || c == '>' || c == '/') {
			break
		}
		val.WriteByte(c)
		p.pos++
	}
	return val.String()
}

func (p *Parser) skipWhitespace() {
	for p.pos < len(p.src) && (p.src[p.pos] == ' ' || p.src[p.pos] == '\t' || p.src[p.pos] == '\n' || p.src[p.pos] == '\r') {
		p.pos++
	}
}

func (p *Parser) peek() byte {
	if p.pos < len(p.src) {
		return p.src[p.pos]
	}
	return 0
}

func (p *Parser) match(s string) bool {
	if p.pos+len(s) <= len(p.src) && p.src[p.pos:p.pos+len(s)] == s {
		p.pos += len(s)
		return true
	}
	return false
}

func (p *Parser) skipToNextTag() {
	for p.pos < len(p.src) && p.src[p.pos] != '<' {
		p.pos++
	}
}

func (t *Template) Set(name string, value interface{}) {
	t.Vars[name] = value
}

func (t *Template) SetMap(vars map[string]interface{}) {
	for k, v := range vars {
		t.Vars[k] = v
	}
}

func interpolate(text string, vars map[string]interface{}) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		key := match[2 : len(match)-1]
		if val, ok := vars[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}

func (p *Parser) extractAttributeVars(elem *Element, attrName, attrVal string) {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(attrVal, -1)
	for _, match := range matches {
		if len(match) > 1 {
			elem.AttributeVars[attrName] = append(elem.AttributeVars[attrName], match[1])
		}
	}
}
