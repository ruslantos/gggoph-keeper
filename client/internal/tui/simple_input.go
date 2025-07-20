package tui

type textinput struct {
	placeholder string
	value       string
}

func (t *textinput) Value() string {
	return t.value
}

func (t *textinput) SetValue(s string) {
	t.value = s
}
