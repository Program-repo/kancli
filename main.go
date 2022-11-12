package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item string

func (i item) FilterValue() string { return "" }

const divisor = 4

var maxCol = 2

var counter int

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
	status      int //status
	verified    bool
	title       string
	description string
}

func (t *Task) Next() {
	if t.status == maxCol {
		t.status = 0
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
	focused  int //status
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
	// fmt.Println(selectedItem)
	// fmt.Println(selectedTask)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))

	return nil
}

func (m *Model) Next() {
	m.focused++
	if m.focused > maxCol {
		m.focused = 0
	}
}

func (m *Model) Prev() {
	m.focused--
	if m.focused < 0 {
		m.focused = maxCol
	}
}

func (m *Model) initLists(width, height int) {
	// Remove item description and spacing
	lndd := list.NewDefaultDelegate()
	lndd.ShowDescription = false
	lndd.SetSpacing(0)

	// Open our jsonFile
	jsonFile, err := os.Open("Tickets.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize our Tickets array
	var tickets Tickets
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'tickets' which we defined above
	json.Unmarshal(byteValue, &tickets)

	// we iterate through every user within our tickets array and
	// print out the user Type, their name, and their facebook url
	// as just an example

	defaultList := list.New([]list.Item{}, lndd, width/divisor, 15) //height/2)
	// defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, 15) //height/2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init To Do)
	// m.lists[0].SetFilteringEnabled(false)
	// m.lists[0].SetShowStatusBar(false)
	// m.lists[0].Title = "To Do"
	// m.lists[0].SetItems([]list.Item{
	// 	Task{status: 0, title: "buy milk", description: "strawberry milk"},
	// 	Task{status: 0, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
	// 	Task{status: 0, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	// })
	// // Init in progress
	// m.lists[1].SetFilteringEnabled(false)
	// m.lists[1].SetShowStatusBar(false)
	// m.lists[1].Title = "In Progress"
	// m.lists[1].SetItems([]list.Item{
	// 	Task{status: 1, title: "write code", description: "don't worry, it's Go"},
	// })
	// // Init done
	// m.lists[2].SetFilteringEnabled(false)
	// m.lists[2].SetShowStatusBar(false)
	// m.lists[2].Title = "Done"
	// m.lists[2].SetItems([]list.Item{
	// 	Task{status: 2, title: "stay cool", description: "as a cucumber"},
	// })

	for i := 0; i < len(tickets.Tickets); i++ {
		// fmt.Println("User Ticketid:" + tickets.Tickets[i].TicketId)
		// fmt.Println("User TicketN. " + tickets.Tickets[i].TicketNumber)
		// fmt.Println("User Detail:  " + tickets.Tickets[i].DetailLine.Detail)

		m.lists[i].SetFilteringEnabled(false)
		m.lists[i].SetShowStatusBar(false)
		m.lists[i].Title = tickets.Tickets[i].TicketId
		m.lists[i].SetItems([]list.Item{
			Task{status: i, title: tickets.Tickets[i].DetailLine.Detail},
		})

	}
}

type Tickets struct {
	Tickets []Ticket `json:"tickets"`
}

// User struct which contains a name
// a type and a list of social links
type Ticket struct {
	TicketId     string     `json:"ticketid"`
	TicketNumber string     `json:"ticketnumber"`
	DetailLine   DetailLine `json:"detailline"`
}

// list of links
type DetailLine struct {
	Detail string `json:"detail"`
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
			// return m, m.MoveToNext
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

	var xRender []string
	if m.loaded {
		for i := 0; i <= maxCol; i++ {
			xView := m.lists[i].View()
			if i == m.focused {
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
}

func main() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
