package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// adding new method to an existing type in Go?
// type nmDD list.DefaultDelegate

// func (d *nmDD) SetShowDescription(v bool) {
// 	d.ShowDescription = v
// }

type item string

func (i item) FilterValue() string { return "" }

type status int

const divisor = 4

var counter int

const (
	todo status = iota
	inProgress
	done
)

/* MODEL MANAGEMENT */

/* STYLING */
var (
	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, true, true, true).
			BorderBackground(lipgloss.AdaptiveColor{Light: "#43BFfD", Dark: "#73F5fd"}).
			BorderStyle(lipgloss.Border{Top: " "}).
			MarginRight(0)
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, true, true, true).
			BorderBackground(lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#0000ff"}).
			BorderStyle(lipgloss.Border{Top: " "}).
			MarginRight(0)

	// helpStyle = lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("241"))
)

/* CUSTOM ITEM */

type Task struct {
	status      status
	verified    bool
	title       string
	description string
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string { return "" }

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	loaded   bool
	focused  status
	lists    []list.Model
	err      error
	quitting bool
}

func New() *Model {
	return &Model{}
}

// func (m *Model) itemDone() tea.Msg {
// selectedItem := m.lists[m.focused].SelectedItem()
// selectedTask := selectedItem.(Task)
// list.Item(selectedTask)
// selectedItem.(status: todo, done: true, title: "buy milk", description: "strawberry milk")
// 	return nil
// }

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(Task)
	selectedTask.title = "jjjjj"
	fmt.Println(selectedItem)
	fmt.Println(selectedTask)
	// fmt.Println(selectedTask.description)
	// fmt.Println(m.lists[m.focused].Index())
	// m.lists[selectedTask.status] [m.lists[m.focused].Index()]

	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))

	return nil
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	lndd := list.NewDefaultDelegate()
	lndd.ShowDescription = false
	lndd.SetSpacing(0)

	defaultList := list.New([]list.Item{}, lndd, width/divisor, 15) //height/2)
	// defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, 15) //height/2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init To Do)
	m.lists[todo].SetFilteringEnabled(false)
	m.lists[todo].SetShowStatusBar(false)
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	m.lists[inProgress].SetFilteringEnabled(false)
	m.lists[inProgress].SetShowStatusBar(false)
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	m.lists[done].SetFilteringEnabled(false)
	m.lists[done].SetShowStatusBar(false)
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(15)     //msg.Height - divisor)
			focusedStyle.Height(15)    //(msg.Height - divisor)
			m.initLists(msg.Width, 15) //msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
			// return m, m.itemDone
		}
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var x int
	var xRender []string
	if m.loaded {

		for i := 0; i < 3; i++ {
			switch m.focused {
			case inProgress:
				x = 1
			case done:
				x = 2
			default:
				x = 0
			}

			xView := m.lists[i].View()

			if x == i {
				xRender = append(xRender, focusedStyle.Render(xView))
			} else {
				xRender = append(xRender, columnStyle.Render(xView))
			}
		}

	} else {
		counter++
		return "loading..." + strconv.Itoa(counter)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		xRender...,
	)

	// if m.loaded {
	// 	todoView := m.lists[todo].View()
	// 	inProgView := m.lists[inProgress].View()
	// 	doneView := m.lists[done].View()
	// 	switch m.focused {
	// 	case inProgress:
	// 		return lipgloss.JoinHorizontal(
	// 			lipgloss.Left,
	// 			columnStyle.Render(todoView),
	// 			focusedStyle.Render(inProgView),
	// 			columnStyle.Render(doneView),
	// 		)
	// 	case done:
	// 		return lipgloss.JoinHorizontal(
	// 			lipgloss.Left,
	// 			columnStyle.Render(todoView),
	// 			columnStyle.Render(inProgView),
	// 			focusedStyle.Render(doneView),
	// 		)
	// 	default:
	// 		return lipgloss.JoinHorizontal(
	// 			lipgloss.Left,
	// 			focusedStyle.Render(todoView),
	// 			columnStyle.Render(inProgView),
	// 			columnStyle.Render(doneView),
	// 		)
	// 	}
	// } else {
	// 	counter++
	// 	return "loading..." + strconv.Itoa(counter)
	// }
}

func main() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
