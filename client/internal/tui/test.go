package tui

//
//import (
//	"fmt"
//
//	tea "github.com/charmbracelet/bubbletea"
//)
//
//// Модель для BubbleTea
//type model struct {
//	choices  []string         // доступные варианты
//	cursor   int              // текущая позиция курсора
//	selected map[int]struct{} // выбранные варианты
//}
//
//// Инициализация модели
//func InitialModel() model {
//	return model{
//		choices:  []string{"Регистрация", "Логин"},
//		selected: make(map[int]struct{}),
//	}
//}
//
//// Init - инициализация BubbleTea (не требуется в этом примере)
//func (m model) Init() tea.Cmd {
//	return nil
//}
//
//// Update - обработка сообщений
//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "ctrl+c", "q":
//			return m, tea.Quit
//		case "up", "k":
//			if m.cursor > 0 {
//				m.cursor--
//			}
//		case "down", "j":
//			if m.cursor < len(m.choices)-1 {
//				m.cursor++
//			}
//		case "enter", " ":
//			// При выборе команды
//			choice := m.choices[m.cursor]
//			switch choice {
//			case "Регистрация":
//				fmt.Println("\nВы выбрали Регистрация 1 - Привет, мир!")
//			case "Команда 2":
//				fmt.Println("\nВы выбрали Логин 2 - Пока, мир!")
//			}
//			return m, tea.Quit
//		}
//	}
//	return m, nil
//}
//
//// View - отображение интерфейса
//func (m model) View() string {
//	s := "Выберите команду:\n\n"
//
//	for i, choice := range m.choices {
//		cursor := " " // нет курсора
//		if m.cursor == i {
//			cursor = ">" // курсор на этом элементе
//		}
//
//		s += fmt.Sprintf("%s %s\n", cursor, choice)
//	}
//
//	s += "\nНажмите q для выхода, ↑/↓ для перемещения, enter для выбора\n"
//	return s
//}
