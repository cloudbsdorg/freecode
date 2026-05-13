package args

type Args struct {
	Continue  bool
	SessionID string
	Agent     string
	Model     string
	Prompt    string
	Fork      bool
	Setup     bool
	Renderer  string
	Headless  bool
}