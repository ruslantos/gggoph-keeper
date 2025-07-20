package tui

//
//// Модель для BubbleTea
//type model struct {
//	choices  []string         // доступные варианты
//	cursor   int              // текущая позиция курсора
//	selected map[int]struct{} // выбранные варианты
//	logs     []string         // лог важных действий
//	baseURL  string           // базовый URL сервера
//}
//
//// Инициализация модели
//func InitialModel(baseURL string) model {
//	return model{
//		choices:  []string{"Регистрация", "Логин", "Получить данные"},
//		selected: make(map[int]struct{}),
//		logs:     []string{"Приложение запущено"},
//		baseURL:  baseURL,
//	}
//}
//
//// Init - инициализация BubbleTea
//func (m model) Init() tea.Cmd {
//	return nil
//}
//
//// Update - обработка сообщений
//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		return m.handleKeyMsg(msg)
//	case string: // Обработка результатов команд
//		m.addLog(msg)
//		return m, nil
//	}
//	return m, nil
//}
//
//// Обработка нажатий клавиш
//func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
//	switch msg.String() {
//	case "ctrl+c", "q":
//		m.addLog("Выход из приложения")
//		return m, tea.Quit
//	case "up", "k":
//		if m.cursor > 0 {
//			m.cursor--
//		}
//	case "down", "j":
//		if m.cursor < len(m.choices)-1 {
//			m.cursor++
//		}
//	case "enter", " ":
//		choice := m.choices[m.cursor]
//		m.addLog(fmt.Sprintf("Выполняется: %s", choice))
//		return m, m.handleCommand(choice)
//	}
//	return m, nil
//}
//
//// Обработка выбранной команды
//func (m model) handleCommand(choice string) tea.Cmd {
//	switch choice {
//	case "Регистрация":
//		return m.registerUser()
//	case "Логин":
//		return m.loginUser()
//	case "Получить данные":
//		return m.fetchData()
//	default:
//		return nil
//	}
//}
//
//// View - отображение интерфейса
//func (m model) View() string {
//	var buf bytes.Buffer
//
//	// Отображение меню
//	buf.WriteString("Выберите команду:\n\n")
//	for i, choice := range m.choices {
//		cursor := " "
//		if m.cursor == i {
//			cursor = ">"
//		}
//		buf.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
//	}
//
//	// Отображение логов
//	buf.WriteString("\n\nЛоги:\n")
//	for _, logEntry := range m.logs {
//		buf.WriteString(fmt.Sprintf("• %s\n", logEntry))
//	}
//
//	buf.WriteString("\nНажмите q для выхода, ↑/↓ для перемещения, enter для выбора\n")
//	return buf.String()
//}
//
//// Добавление записи в лог
//func (m *model) addLog(message string) {
//	m.logs = append(m.logs, message)
//	if len(m.logs) > 10 { // Ограничим количество логов на экране
//		m.logs = m.logs[1:]
//	}
//}
//
//// Регистрация пользователя
//func (m model) registerUser() tea.Cmd {
//	return func() tea.Msg {
//		requestBody, _ := json.Marshal(map[string]string{
//			"username": "user123",
//			"password": "pass123",
//		})
//
//		resp, err := http.Post(m.baseURL+"/register", "application/json", bytes.NewBuffer(requestBody))
//		if err != nil {
//			return fmt.Sprintf("Ошибка регистрации: %v", err)
//		}
//		defer resp.Body.Close()
//
//		body, _ := io.ReadAll(resp.Body)
//		return fmt.Sprintf("Регистрация успешна 🙂: %s", string(body))
//	}
//}
//
//// Логин пользователя
//func (m model) loginUser() tea.Cmd {
//	return func() tea.Msg {
//		requestBody, _ := json.Marshal(map[string]string{
//			"username": "user123",
//			"password": "pass123",
//		})
//
//		resp, err := http.Post(m.baseURL+"/login", "application/json", bytes.NewBuffer(requestBody))
//		if err != nil {
//			return fmt.Sprintf("Ошибка входа: %v", err)
//		}
//		defer resp.Body.Close()
//
//		body, _ := io.ReadAll(resp.Body)
//		return fmt.Sprintf("Авторизация: %s", string(body))
//	}
//}
//
//// Получение данных
//func (m model) fetchData() tea.Cmd {
//	return func() tea.Msg {
//		resp, err := http.Get(m.baseURL + "/data")
//		if err != nil {
//			return fmt.Sprintf("Ошибка получения данных: %v", err)
//		}
//		defer resp.Body.Close()
//
//		body, _ := io.ReadAll(resp.Body)
//		return fmt.Sprintf("Данные: %s", string(body))
//	}
//}
//
//// Cobra команда для запуска TUI
//func NewTUICommand() *cobra.Command {
//	var serverURL string
//
//	cmd := &cobra.Command{
//		Use:   "tui",
//		Short: "Запуск TUI клиента",
//		Run: func(cmd *cobra.Command, args []string) {
//			log.Println("Запуск TUI клиента")
//			p := tea.NewProgram(InitialModel(serverURL))
//			if _, err := p.Run(); err != nil {
//				fmt.Printf("Ошибка запуска TUI: %v", err)
//				os.Exit(1)
//			}
//		},
//	}
//
//	cmd.Flags().StringVarP(&serverURL, "server", "s", "http://localhost:8080", "URL сервера")
//	return cmd
//}
