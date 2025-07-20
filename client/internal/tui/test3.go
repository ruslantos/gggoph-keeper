package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const (
	stateMenu = iota
	stateRegister
	stateLogin
)

type model struct {
	state      int
	choices    []string
	cursor     int
	baseURL    string
	logs       []string
	inputMode  bool
	inputText  string
	inputLabel string
	inputField int
	username   string
	password   string
	token      string
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func InitialModel(baseURL string) model {
	return model{
		state:   stateMenu,
		choices: []string{"Регистрация", "Логин", "Выход"},
		baseURL: baseURL,
		logs:    []string{"Добро пожаловать!"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case string:
		m.addLog(msg)
		return m, nil
	case model:
		// Обработка возврата новой модели
		return msg, nil
	}
	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.inputMode {
		switch msg.Type {
		case tea.KeyRunes:
			m.inputText += string(msg.Runes)
		case tea.KeySpace:
			m.inputText += " "
		case tea.KeyBackspace:
			if len(m.inputText) > 0 {
				m.inputText = m.inputText[:len(m.inputText)-1]
			}
		case tea.KeyEnter:
			if strings.TrimSpace(m.inputText) == "" {
				m.addLog("Ошибка: поле не может быть пустым")
				return m, nil
			}

			if m.inputField == 0 {
				m.username = m.inputText
				m.inputField = 1
				m.inputText = ""
				m.inputLabel = "Пароль: "
			} else {
				m.password = m.inputText
				m.inputMode = false
				if m.state == stateRegister {
					return m, m.registerUser()
				}
				return m, m.loginUser()
			}
		case tea.KeyEscape:
			m.resetInput()
			return m, nil
		}
		return m, nil
	}

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case tea.KeyEnter:
		if m.state == stateMenu {
			choice := m.choices[m.cursor]
			switch choice {
			case "Регистрация":
				m.state = stateRegister
				m.inputMode = true
				m.inputField = 0
				m.inputText = ""
				m.inputLabel = "Логин: "
				m.addLog("Начата регистрация")
			case "Логин":
				m.state = stateLogin
				m.inputMode = true
				m.inputField = 0
				m.inputText = ""
				m.inputLabel = "Логин: "
				m.addLog("Начата авторизация")
			case "Выход":
				return m, tea.Quit
			}
		}
	case tea.KeyRunes:
		// Обработка символов 'j' и 'k' для перемещения
		switch msg.String() {
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m *model) resetInput() {
	m.inputMode = false
	m.inputText = ""
	m.inputField = 0
	m.state = stateMenu
}

func (m model) View() string {
	if m.inputMode {
		return m.viewInput()
	}
	return m.viewMenu()
}

func (m model) viewMenu() string {
	var sb strings.Builder
	sb.WriteString("Выберите действие:\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	sb.WriteString("\nЛоги:\n")
	for _, logEntry := range m.logs {
		sb.WriteString(fmt.Sprintf("• %s\n", logEntry))
	}

	// Показываем статус авторизации
	if m.token != "" {
		sb.WriteString("\nСтатус: Авторизован\n")
	} else {
		sb.WriteString("\nСтатус: Не авторизован\n")
	}

	sb.WriteString("\nq/ctrl+c - выход | ↑↓/jk - перемещение | enter - выбор")
	return sb.String()
}

func (m model) viewInput() string {
	var sb strings.Builder
	title := "Регистрация"
	if m.state == stateLogin {
		title = "Вход"
	}
	sb.WriteString(fmt.Sprintf("%s\n\n", title))
	sb.WriteString(m.inputLabel)

	if m.inputField == 1 {
		sb.WriteString(strings.Repeat("*", len(m.inputText)))
	} else {
		sb.WriteString(m.inputText)
	}

	sb.WriteString("\n\n(enter - подтвердить, esc - отмена, backspace - удалить)\n")
	sb.WriteString("\nЛоги:\n")
	for _, logEntry := range m.logs {
		sb.WriteString(fmt.Sprintf("• %s\n", logEntry))
	}

	return sb.String()
}

func (m *model) addLog(message string) {
	m.logs = append(m.logs, message)
	if len(m.logs) > 10 {
		m.logs = m.logs[1:]
	}
}

func (m model) registerUser() tea.Cmd {
	return func() tea.Msg {
		requestBody, _ := json.Marshal(map[string]string{
			"login":    m.username,
			"password": m.password,
		})

		resp, err := http.Post(m.baseURL+"/register", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return fmt.Sprintf("Ошибка регистрации: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		m.resetInput()
		return fmt.Sprintf("🙂 Регистрация успешна: %s", string(body))
	}
}

func (m model) loginUser() tea.Cmd {
	return func() tea.Msg {
		requestBody, _ := json.Marshal(map[string]string{
			"login":    m.username,
			"password": m.password,
		})

		resp, err := http.Post(m.baseURL+"/login", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return fmt.Sprintf("Ошибка входа: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Sprintf("Ошибка входа: статус %d", resp.StatusCode)
		}

		var authResp AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			return fmt.Sprintf("Ошибка разбора ответа: %v", err)
		}

		// Создаем новую модель с обновленным токеном
		newModel := m
		newModel.token = authResp.Token
		newModel.resetInput()
		newModel.addLog("🙂 Авторизация успешна")
		return newModel
	}
}

func (m model) getProtectedData() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.makeAuthenticatedRequest("GET", "/protected/data", nil)
		if err != nil {
			return fmt.Sprintf("Ошибка запроса: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		return fmt.Sprintf("Защищенные данные: %s", string(body))
	}
}

// makeAuthenticatedRequest создает авторизованный HTTP-запрос с токеном из памяти
func (m *model) makeAuthenticatedRequest(method, path string, body io.Reader) (*http.Response, error) {
	if m.token == "" {
		return nil, fmt.Errorf("требуется авторизация")
	}

	req, err := http.NewRequest(method, m.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+m.token)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func NewTUICommand() *cobra.Command {
	var serverURL string

	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Запуск TUI клиента",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Запуск TUI клиента")
			p := tea.NewProgram(InitialModel(serverURL))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Ошибка запуска TUI: %v", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&serverURL, "server", "s", "http://localhost:8080", "URL сервера")
	return cmd
}
