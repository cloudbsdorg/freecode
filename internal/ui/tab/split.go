package tab

type SplitDirection bool

const (
	SplitVertical   SplitDirection = true
	SplitHorizontal SplitDirection = false
)

type Split struct {
	ID        string
	Direction SplitDirection
	Ratio     float64
	Left      interface{}
	Right     interface{}
}

type SplitState struct {
	splits map[string]*Split
}

func NewSplitState() *SplitState {
	return &SplitState{
		splits: make(map[string]*Split),
	}
}

func (s *SplitState) Create(id string, direction SplitDirection, ratio float64) *Split {
	spl := &Split{
		ID:        id,
		Direction: direction,
		Ratio:     ratio,
	}
	s.splits[id] = spl
	return spl
}

func (s *SplitState) Get(id string) (*Split, bool) {
	spl, ok := s.splits[id]
	return spl, ok
}

func (s *SplitState) Remove(id string) {
	delete(s.splits, id)
}

func (s *SplitState) SetRatio(id string, ratio float64) bool {
	spl, ok := s.splits[id]
	if !ok {
		return false
	}
	spl.Ratio = ratio
	return true
}

func (s *SplitState) Toggle(id string) {
	spl, ok := s.splits[id]
	if !ok {
		return
	}
	if spl.Direction == SplitVertical {
		spl.Direction = SplitHorizontal
	} else {
		spl.Direction = SplitVertical
	}
}

func (s *SplitState) List() []*Split {
	splits := make([]*Split, 0, len(s.splits))
	for _, spl := range s.splits {
		splits = append(splits, spl)
	}
	return splits
}
