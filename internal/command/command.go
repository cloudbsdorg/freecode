package command

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type Command struct {
	Name string
	Args []string
	Dir  string
	Env  []string
}

type Runner interface {
	Run(ctx context.Context, cmd Command) error
	RunWithOutput(ctx context.Context, cmd Command) (string, error)
}

type registry struct {
	commands map[string]func(args []string) error
}

var reg = &registry{commands: make(map[string]func(args []string) error)}

func Register(name string, fn func(args []string) error) {
	reg.commands[name] = fn
}

func Get(name string) (func(args []string) error, bool) {
	fn, ok := reg.commands[name]
	return fn, ok
}

func List() []string {
	var names []string
	for name := range reg.commands {
		names = append(names, name)
	}
	return names
}

func Run(ctx context.Context, cmd Command, runner Runner) error {
	return runner.Run(ctx, cmd)
}

func RunOutput(ctx context.Context, cmd Command, runner Runner) (string, error) {
	return runner.RunWithOutput(ctx, cmd)
}

type Template struct {
	Name        string
	Pattern    string
	Args       []TemplateArg
	Help       string
	Handler    func(ctx context.Context, args map[string]string) (string, error)
}

type TemplateArg struct {
	Name     string
	Type     string
	Required bool
	Default string
	Help    string
}

type templateRegistry struct {
	templates map[string]*Template
}

var templates = &templateRegistry{
	templates: make(map[string]*Template),
}

func RegisterTemplate(t *Template) {
	templates.templates[t.Name] = t
}

func GetTemplate(name string) (*Template, bool) {
	t, ok := templates.templates[name]
	return t, ok
}

func ListTemplates() []*Template {
	result := make([]*Template, 0, len(templates.templates))
	for _, t := range templates.templates {
		result = append(result, t)
	}
	return result
}

var templateVarPattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func (t *Template) Render(args map[string]string) (string, error) {
	result := t.Pattern

	matches := templateVarPattern.FindAllStringSubmatch(t.Pattern, -1)
	for _, match := range matches {
		varName := strings.TrimSpace(match[1])
		value, ok := args[varName]
		if !ok {
			for _, arg := range t.Args {
				if arg.Name == varName && arg.Default != "" {
					value = arg.Default
					break
				}
			}
		}
		if !ok && value == "" {
			for _, arg := range t.Args {
				if arg.Name == varName && arg.Required {
					return "", fmt.Errorf("required argument missing: %s", varName)
				}
			}
		}
		if value != "" {
			result = strings.Replace(result, match[0], value, 1)
		}
	}

	return result, nil
}

func (t *Template) Validate(args map[string]string) error {
	for _, arg := range t.Args {
		if arg.Required {
			if val, ok := args[arg.Name]; !ok || val == "" {
				if arg.Default == "" {
					return fmt.Errorf("required argument missing: %s", arg.Name)
				}
			}
		}
	}
	return nil
}

func (t *Template) Execute(ctx context.Context, args map[string]string) (string, error) {
	if err := t.Validate(args); err != nil {
		return "", err
	}

	if t.Handler != nil {
		return t.Handler(ctx, args)
	}

	return "", fmt.Errorf("template has no handler")
}

type TemplateParser struct{}

func NewTemplateParser() *TemplateParser {
	return &TemplateParser{}
}

func (p *TemplateParser) Parse(pattern string) ([]string, error) {
	matches := templateVarPattern.FindAllStringSubmatch(pattern, -1)
	varNames := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		varName := strings.TrimSpace(match[1])
		if !seen[varName] {
			varNames = append(varNames, varName)
			seen[varName] = true
		}
	}

	return varNames, nil
}

func (p *TemplateParser) Expand(pattern string, args map[string]string) (string, error) {
	t := &Template{Pattern: pattern}
	return t.Render(args)
}
