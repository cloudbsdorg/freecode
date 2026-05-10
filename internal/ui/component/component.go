package component

type Component[R any] struct {
	X, Y    int
	Width   int
	Height  int
	Visible bool
}

func (c *Component[R]) SetPosition(x, y int) {
	c.X = x
	c.Y = y
}

func (c *Component[R]) SetSize(w, h int) {
	c.Width = w
	c.Height = h
}

func (c *Component[R]) Show() {
	c.Visible = true
}

func (c *Component[R]) Hide() {
	c.Visible = false
}

func (c *Component[R]) IsVisible() bool {
	return c.Visible
}
