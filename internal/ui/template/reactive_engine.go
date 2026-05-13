package template

import (
	"strconv"
	"strings"
	"sync"

	"github.com/freecode/freecode/internal/renderer"
)

type ReactiveEngine[R renderer.Renderer] struct {
	*ResponsiveEngine[R]
	State      map[string]interface{}
	dirty      map[string]bool
	mu         sync.RWMutex
	onChange   func(key string, value interface{})
	loader     *Loader
	templates  map[string]string
	lastRender string
}

func NewReactiveEngine[R renderer.Renderer]() *ReactiveEngine[R] {
	return &ReactiveEngine[R]{
		ResponsiveEngine: NewResponsiveEngine[R](),
		State:            make(map[string]interface{}),
		dirty:            make(map[string]bool),
		templates:        make(map[string]string),
	}
}

func NewReactiveEngineWithLoader[R renderer.Renderer](baseDir string) *ReactiveEngine[R] {
	loader := NewLoader(baseDir)
	re := NewReactiveEngine[R]()
	re.loader = loader
	re.templates = loader.LoadViews()
	if root, err := loader.LoadRoot(); err == nil {
		re.templates["ROOT"] = root
	}
	return re
}

func (re *ReactiveEngine[R]) LoadTemplate(name string) error {
	if re.loader == nil {
		return nil
	}
	content, err := re.loader.LoadView(name)
	if err != nil {
		content, err = re.loader.LoadView(name)
	}
	if err == nil {
		re.templates[name] = content
	}
	return err
}

func (re *ReactiveEngine[R]) LoadComponent(name string) error {
	if re.loader == nil {
		return nil
	}
	content, err := re.loader.LoadComponent(name)
	if err == nil {
		re.templates["component_"+name] = content
	}
	return err
}

func (re *ReactiveEngine[R]) GetTemplate(name string) string {
	return re.templates[name]
}

func (re *ReactiveEngine[R]) SetTemplate(name string, content string) {
	re.templates[name] = content
}

func (re *ReactiveEngine[R]) Set(key string, value interface{}) {
	re.mu.Lock()
	defer re.mu.Unlock()

	if re.State[key] != value {
		re.State[key] = value
		re.markAllDirty()
		if re.onChange != nil {
			re.onChange(key, value)
		}
	}
}

func (re *ReactiveEngine[R]) Get(key string) interface{} {
	re.mu.RLock()
	defer re.mu.RUnlock()
	return re.State[key]
}

func (re *ReactiveEngine[R]) GetString(key, def string) string {
	re.mu.RLock()
	defer re.mu.RUnlock()
	if val, ok := re.State[key].(string); ok {
		return val
	}
	return def
}

func (re *ReactiveEngine[R]) GetInt(key string, def int) int {
	re.mu.RLock()
	defer re.mu.RUnlock()
	if val, ok := re.State[key].(int); ok {
		return val
	}
	return def
}

func (re *ReactiveEngine[R]) GetBool(key string, def bool) bool {
	re.mu.RLock()
	defer re.mu.RUnlock()
	if val, ok := re.State[key].(bool); ok {
		return val
	}
	return def
}

func (re *ReactiveEngine[R]) SetMany(values map[string]interface{}) {
	re.mu.Lock()
	defer re.mu.Unlock()

	for key, value := range values {
		if re.State[key] != value {
			re.State[key] = value
			if re.onChange != nil {
				re.onChange(key, value)
			}
		}
	}
	re.markAllDirty()
}

func (re *ReactiveEngine[R]) markAllDirty() {
	for id := range re.components {
		re.dirty[id] = true
	}
	re.dirty["__root__"] = true
}

func (re *ReactiveEngine[R]) MarkDirty(id string) {
	re.mu.Lock()
	defer re.mu.Unlock()
	re.dirty[id] = true
}

func (re *ReactiveEngine[R]) IsDirty(id string) bool {
	re.mu.RLock()
	defer re.mu.RUnlock()
	return re.dirty[id]
}

func (re *ReactiveEngine[R]) ClearDirty(id string) {
	re.mu.Lock()
	defer re.mu.Unlock()
	delete(re.dirty, id)
}

func (re *ReactiveEngine[R]) ClearAllDirty() {
	re.mu.Lock()
	defer re.mu.Unlock()
	re.dirty = make(map[string]bool)
}

func (re *ReactiveEngine[R]) IsAnyDirty() bool {
	re.mu.RLock()
	defer re.mu.RUnlock()
	return len(re.dirty) > 0
}

func (re *ReactiveEngine[R]) OnStateChange(callback func(key string, value interface{})) {
	re.onChange = callback
}

func (re *ReactiveEngine[R]) Interpolate(text string) string {
	re.mu.RLock()
	defer re.mu.RUnlock()

	result := text
	for key, value := range re.State {
		placeholder := "${" + key + "}"
		valueStr := formatValue(value)
		result = strings.ReplaceAll(result, placeholder, valueStr)

		defaultPattern := "${" + key + ":"
		for strings.Contains(result, defaultPattern) {
			idx := strings.Index(result, defaultPattern)
			endIdx := strings.Index(result[idx:], "}")
			if endIdx == -1 {
				break
			}
			endIdx += idx
			defaultStart := idx + len(defaultPattern)
			defaultVal := result[defaultStart:endIdx]
			displayVal := valueStr
			if displayVal == "" && defaultVal != "" {
				displayVal = defaultVal
			}
			result = result[:idx] + displayVal + result[endIdx+1:]
		}
	}
	return result
}

func (re *ReactiveEngine[R]) InterpolateMap(attrs map[string]string) map[string]string {
	re.mu.RLock()
	defer re.mu.RUnlock()

	result := make(map[string]string)
	for key, value := range attrs {
		result[key] = re.Interpolate(value)
	}
	return result
}

func (re *ReactiveEngine[R]) MustRender(src string, width, height int, r R) string {
	re.mu.Lock()
	interpolated := re.Interpolate(src)
	re.ClearAllDirty()
	re.mu.Unlock()

	result, err := re.ResponsiveEngine.RenderAt(interpolated, width, height, r)
	if err != nil {
		return "Error: " + err.Error()
	}

	re.mu.Lock()
	re.lastRender = result
	re.mu.Unlock()

	return result
}

func (re *ReactiveEngine[R]) RenderTemplate(name string, width, height int, r R) string {
	re.mu.RLock()
	src, ok := re.templates[name]
	re.mu.RUnlock()
	if !ok {
		return "Template not found: " + name
	}
	return re.MustRender(src, width, height, r)
}

func (re *ReactiveEngine[R]) RenderDirty(src string, width, height int, r R) string {
	if !re.IsAnyDirty() {
		re.mu.RLock()
		defer re.mu.RUnlock()
		return re.lastRender
	}
	return re.MustRender(src, width, height, r)
}

func (re *ReactiveEngine[R]) GetState() map[string]interface{} {
	re.mu.RLock()
	defer re.mu.RUnlock()

	stateCopy := make(map[string]interface{}, len(re.State))
	for k, v := range re.State {
		stateCopy[k] = v
	}
	return stateCopy
}

func (re *ReactiveEngine[R]) Subscribe() <-chan StateChange {
	ch := make(chan StateChange, 100)
	re.mu.Lock()
	oldCallback := re.onChange
	re.onChange = func(key string, value interface{}) {
		ch <- StateChange{Key: key, Value: value}
		if oldCallback != nil {
			oldCallback(key, value)
		}
	}
	re.mu.Unlock()
	return ch
}

type StateChange struct {
	Key   string
	Value interface{}
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case []string:
		return strings.Join(val, ",")
	case []int:
		parts := make([]string, len(val))
		for i, v := range val {
			parts[i] = strconv.FormatInt(int64(v), 10)
		}
		return strings.Join(parts, ",")
	default:
		return ""
	}
}
