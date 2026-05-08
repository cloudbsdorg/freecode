package ui

type AnimationSetting int

const (
	AnimationNone AnimationSetting = iota
	AnimationMinimal
	AnimationFull
)

func (a AnimationSetting) String() string {
	switch a {
	case AnimationNone:
		return "none"
	case AnimationMinimal:
		return "minimal"
	case AnimationFull:
		return "full"
	default:
		return "unknown"
	}
}

func ParseAnimationSetting(s string) AnimationSetting {
	switch s {
	case "none":
		return AnimationNone
	case "minimal":
		return AnimationMinimal
	case "full":
		return AnimationFull
	default:
		return AnimationFull
	}
}

type AnimationManager struct {
	level AnimationSetting
}

const (
	AnimationDurationFast    = 100
	AnimationDurationNormal  = 200
	AnimationDurationSlow    = 400
)

func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		level: AnimationFull,
	}
}

func NewAnimationManagerWithLevel(level AnimationSetting) *AnimationManager {
	return &AnimationManager{
		level: level,
	}
}

func (a *AnimationManager) IsEnabled() bool {
	return a.level != AnimationNone
}

func (a *AnimationManager) Level() AnimationSetting {
	return a.level
}

func (a *AnimationManager) SetLevel(level AnimationSetting) {
	a.level = level
}

func (a *AnimationManager) Toggle() {
	switch a.level {
	case AnimationFull:
		a.level = AnimationMinimal
	case AnimationMinimal:
		a.level = AnimationNone
	case AnimationNone:
		a.level = AnimationFull
	}
}

func (a *AnimationManager) GetDuration(ms int) int {
	if a.level == AnimationNone {
		return 0
	}

	var scale float64
	switch a.level {
	case AnimationMinimal:
		scale = 0.5
	case AnimationFull:
		scale = 1.0
	default:
		scale = 1.0
	}

	return int(float64(ms) * scale)
}

func (a *AnimationManager) GetFadeDuration() int {
	return a.GetDuration(AnimationDurationNormal)
}

func (a *AnimationManager) GetTransitionDuration() int {
	return a.GetDuration(AnimationDurationFast)
}

func (a *AnimationManager) GetToastDuration(baseMs int) int {
	return a.GetDuration(baseMs)
}

func (a *AnimationManager) FastDuration() int {
	return a.GetDuration(AnimationDurationFast)
}

func (a *AnimationManager) NormalDuration() int {
	return a.GetDuration(AnimationDurationNormal)
}

func (a *AnimationManager) SlowDuration() int {
	return a.GetDuration(AnimationDurationSlow)
}
