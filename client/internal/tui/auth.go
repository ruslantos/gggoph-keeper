package tui

//
//import (
//	"context"
//	"fmt"
//	"strings"
//
//	"github.com/charmbracelet/bubbles/key"
//	"github.com/charmbracelet/bubbles/list"
//	"github.com/charmbracelet/bubbles/textinput"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//
//	api "goph-keeper/client/internal/client"
//)
//
//// Стили
//var (
//	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
//	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
//	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
//	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
//)
//
//// item представляет элемент списка
//type item struct {
//	title, desc string
//}
//
//func (i item) Title() string       { return i.title }
//func (i item) Description() string { return i.desc }
//func (i item) FilterValue() string { return i.title }
//
//// modelList реализует модель для выбора режима
//type modelList struct {
//	list   list.Model
//	choice string
//}
//
//func (m modelList) Init() tea.Cmd { return nil }
//
//func (m modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "enter":
//			if selected, ok := m.list.SelectedItem().(item); ok {
//				m.choice = selected.title
//				return m, tea.Quit
//			}
//		case "q", "esc":
//			return m, tea.Quit
//		}
//	}
//
//	var cmd tea.Cmd
//	m.list, cmd = m.list.Update(msg)
//	return m, cmd
//}
//
//func (m modelList) View() string {
//	return m.list.View()
//}
//
//// showModeSelector отображает выбор режима
//func showModeSelector() (string, error) {
//	items := []list.Item{
//		item{title: "Login", desc: "Sign in to existing account"},
//		item{title: "Register", desc: "Create new account"},
//	}
//
//	// Настройка делегата для элементов списка
//	delegate := list.NewDefaultDelegate()
//	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
//		Border(lipgloss.NormalBorder(), false, false, false, true).
//		BorderForeground(lipgloss.Color("205")).
//		Foreground(lipgloss.Color("205")).
//		Padding(0, 0, 0, 1)
//
//	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy().
//		Foreground(lipgloss.Color("240"))
//
//	// Настройка списка
//	l := list.New(items, delegate, 0, 0)
//	l.Title = "Select Auth Mode"
//	l.SetShowStatusBar(false)
//	l.SetShowHelp(true)
//	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)
//	l.Styles.HelpStyle = lipgloss.NewStyle().
//		Foreground(lipgloss.Color("240")).
//		Padding(1, 0, 0, 2)
//
//	// Явная настройка клавиш
//	l.KeyMap.Quit.SetKeys("q", "esc")
//	l.KeyMap.CursorUp.SetKeys("up", "k")
//	l.KeyMap.CursorDown.SetKeys("down", "j")
//	l.KeyMap.Filter.SetKeys("/")
//	l.KeyMap.ShowFullHelp.SetKeys("?")
//	l.KeyMap.CloseFullHelp.SetKeys("?")
//
//	// Дополнительные подсказки
//	l.AdditionalShortHelpKeys = func() []key.Binding {
//		return []key.Binding{
//			key.NewBinding(
//				key.WithKeys("enter"),
//				key.WithHelp("enter", "select"),
//			),
//		}
//	}
//
//	p := tea.NewProgram(modelList{list: l}, tea.WithAltScreen())
//
//	finalModel, err := p.Run()
//	if err != nil {
//		return "", err
//	}
//
//	if selected, ok := finalModel.(modelList).list.SelectedItem().(item); ok {
//		return strings.ToLower(selected.title), nil
//	}
//	return "", fmt.Errorf("mode selection cancelled")
//}
//
//// authModel реализует модель формы аутентификации
//type authModel struct {
//	inputs  []textinput.Model
//	focused int
//	mode    string
//	client  *api.Client
//	err     error
//}
//
//// Start запускает TUI интерфейс
//func Start(client *api.Client) error {
//	mode, err := showModeSelector()
//	if err != nil {
//		return err
//	}
//
//	m := &authModel{
//		inputs:  createAuthForm(),
//		focused: 0,
//		mode:    mode,
//		client:  client,
//	}
//
//	p := tea.NewProgram(m, tea.WithAltScreen())
//	_, err = p.Run()
//	return err
//}
//
//func createAuthForm() []textinput.Model {
//	inputs := make([]textinput.Model, 2)
//
//	inputs[0] = textinput.New()
//	inputs[0].Placeholder = "Username"
//	inputs[0].Focus()
//	inputs[0].PromptStyle = focusedStyle
//	inputs[0].TextStyle = focusedStyle
//
//	inputs[1] = textinput.New()
//	inputs[1].Placeholder = "Password"
//	inputs[1].EchoMode = textinput.EchoPassword
//	inputs[1].EchoCharacter = '•'
//	inputs[1].PromptStyle = blurredStyle
//	inputs[1].TextStyle = blurredStyle
//
//	return inputs
//}
//
//func (m *authModel) Init() tea.Cmd {
//	return textinput.Blink
//}
//
//func (m *authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	var cmd tea.Cmd
//
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "ctrl+c", "esc":
//			return m, tea.Quit
//		case "tab", "shift+tab", "enter", "up", "down":
//			s := msg.String()
//			if s == "enter" && m.focused == len(m.inputs)-1 {
//				return m, m.authenticate()
//			}
//			if s == "enter" || s == "down" || s == "tab" {
//				m.nextInput()
//			} else if s == "up" || s == "shift+tab" {
//				m.prevInput()
//			}
//			return m, nil
//		}
//	}
//
//	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
//	return m, cmd
//}
//
//func (m *authModel) View() string {
//	var b strings.Builder
//	b.WriteString(fmt.Sprintf("\n %s\n\n", strings.Title(m.mode)))
//
//	for i := range m.inputs {
//		b.WriteString(m.inputs[i].View())
//		if i < len(m.inputs)-1 {
//			b.WriteRune('\n')
//		}
//	}
//
//	if m.err != nil {
//		b.WriteString(fmt.Sprintf("\n\n%s", errorStyle.Render(m.err.Error())))
//	}
//
//	b.WriteString("\n\n" + m.helpView())
//	return b.String()
//}
//
//func (m *authModel) helpView() string {
//	return helpStyle.Render("Tab/↑↓: Navigate • Enter: Submit • Ctrl+C/Esc: Quit")
//}
//
//func (m *authModel) nextInput() {
//	m.focused = (m.focused + 1) % len(m.inputs)
//	m.updateFocus()
//}
//
//func (m *authModel) prevInput() {
//	m.focused--
//	if m.focused < 0 {
//		m.focused = len(m.inputs) - 1
//	}
//	m.updateFocus()
//}
//
//func (m *authModel) updateFocus() {
//	for i := range m.inputs {
//		if i == m.focused {
//			m.inputs[i].Focus()
//			m.inputs[i].PromptStyle = focusedStyle
//			m.inputs[i].TextStyle = focusedStyle
//		} else {
//			m.inputs[i].Blur()
//			m.inputs[i].PromptStyle = blurredStyle
//			m.inputs[i].TextStyle = blurredStyle
//		}
//	}
//}
//
//func (m *authModel) authenticate() tea.Cmd {
//	return func() tea.Msg {
//		ctx := context.Background()
//		username := m.inputs[0].Value()
//		password := m.inputs[1].Value()
//
//		var err error
//		if m.mode == "login" {
//			_, err = m.client.Login(ctx, username, password)
//		} else {
//			err = m.client.Register(ctx, username, password)
//		}
//
//		if err != nil {
//			m.err = err
//			return err
//		}
//		return tea.Quit
//	}
//}
