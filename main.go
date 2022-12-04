package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kancli-demo/llist"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// type item string

// func (i item) FilterValue() string { return "" }

const divisor = 4

var maxCol = 3

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
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	styledef = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#6ad257"))
	styleyew = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#Cad257"))
	stylered = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FCF6F5FF")).
			Background(lipgloss.Color("#990011FF"))

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
	t.status++
	if t.status >= maxCol-1 {
		t.status = 0
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
	loaded    bool
	focused   int //status
	lists     []list.Model
	err       error
	quitting  bool
	stopwatch stopwatch.Model
	linklsd   *llist.LinkedList[linklsdet]
}
type linklsdet struct {
	ticketitems     []list.Item
	ticketid        string
	ticketnumber    string
	ticketordertime time.Time
}

func New() *Model {
	return &Model{stopwatch: stopwatch.NewWithInterval(time.Second * 5)}
}

// func (m *Model) itemDone() tea.Msg {
// selectedItem := m.lists[m.focused].SelectedItem()
// selectedTask := selectedItem.(Task)
// list.Item(selectedTask)
// selectedItem.(status: todo, done: true, title: "buy milk", description: "strawberry milk")
// 	return nil
// }

// func (m *Model) MoveToNext() tea.Msg {
// 	selectedItem := m.lists[m.focused].SelectedItem()
// 	selectedTask := selectedItem.(Task)
// 	selectedTask.title = "jjjjj"
// 	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
// 	selectedTask.Next()
// 	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
// 	return nil
// }

func (m *Model) Next() {
	m.focused++
	if m.focused > m.linklsd.Length()-1 {
		// if m.focused > maxCol-1 {
		m.focused = 0
	}
}

func (m *Model) Prev() {
	m.focused--
	if m.focused < 0 {
		// m.focused = maxCol - 1
		m.focused = m.linklsd.Length() - 1
	}
}

func (m *Model) initLists(width, height int) {
	// Remove item description and set up spacing
	listNewDefaultDelegate := list.NewDefaultDelegate()
	listNewDefaultDelegate.ShowDescription = false
	listNewDefaultDelegate.SetSpacing(0)
	defaultList := list.New([]list.Item{}, listNewDefaultDelegate, width/divisor, 15) //height/2)
	defaultList.SetShowHelp(false)
	var listModel []list.Model
	for i := 0; i < maxCol; i++ {
		listModel = append(listModel, defaultList)
	}
	m.lists = listModel
}

func (m *Model) readData(filename string, li *llist.LinkedList[linklsdet]) {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
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
	// we iterate through every detail line within our tickets

	for i := 0; i < len(tickets.Tickets); i++ {
		m.lists[i].SetFilteringEnabled(false)
		m.lists[i].SetShowStatusBar(false)
		m.lists[i].Title = ""

		var taskslines []Task
		for j := 0; j < len(tickets.Tickets[i].DetailLine.Detail); j++ {
			taskslines = append(taskslines, Task{status: i, title: tickets.Tickets[i].DetailLine.Detail[j]})
		}
		var listItem []list.Item
		for j := 0; j < len(tickets.Tickets[i].DetailLine.Detail); j++ {
			listItem = append(listItem, taskslines[j])
		}
		li.PushBack(linklsdet{
			ticketitems:     listItem,
			ticketid:        tickets.Tickets[i].TicketId,
			ticketnumber:    tickets.Tickets[i].TicketNumber,
			ticketordertime: tickets.Tickets[i].TicketOrderTime,
		})

	}
}

func (m *Model) refreshData(li *llist.LinkedList[linklsdet]) {
	node := li.Head()
	i := 0
	for node != nil {
		// fmt.Println(nodei.Value().ticketid)
		m.lists[i].SetItems(node.Value().ticketitems)
		i++
		node = node.Next()
	}

}

type Tickets struct {
	Tickets []Ticket `json:"tickets"`
}

// User struct which contains a name
type Ticket struct {
	TicketId        string     `json:"ticketid"`
	TicketNumber    string     `json:"ticketnumber"`
	TicketOrderTime time.Time  `json:"ticketordertime"`
	DetailLine      DetailLine `json:"detailline"`
}

// list of links
type DetailLine struct {
	Detail []string `json:"detail"`
}

func (m Model) Init() tea.Cmd {
	return m.stopwatch.Init()
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
			m.linklsd = llist.New[linklsdet]()
			m.readData("Tickets.json", m.linklsd)
			m.refreshData(m.linklsd)

			// m.linklsd.DeleteAt(2)
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
		// The "up" and "k" keys move the cursor up
		case "up", "k":

		// The "down" and "j" keys move the cursor down
		case "down", "j":
		case "enter":
			m.linklsd.DeleteAt(m.focused)
			m.refreshData(m.linklsd)
			// return m, m.MoveToNext
			// return m, m.itemDone
		}

		// case tea.KeyMsg:
		// 	switch {
		// 	case key.Matches(msg, m.keymap.quit):
		// 		m.quitting = true
		// 		return m, tea.Quit
		// 	case key.Matches(msg, m.keymap.reset):
		// 		m.timer.Timeout = timeout
		// 	case key.Matches(msg, m.keymap.start, m.keymap.stop):
		// 		return m, m.timer.Toggle()
		// 	}
	}

	var cmd tea.Cmd
	if m.loaded {
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
		m.refreshData(m.linklsd)
	}
	m.stopwatch, cmd = m.stopwatch.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	var xRender []string
	if m.loaded {
		node := m.linklsd.Head()
		i := 0
		for node != nil {
			wticketid := node.Value().ticketid
			wtordert := node.Value().ticketordertime
			currentTime := time.Now()
			diff := currentTime.Sub(wtordert)
			out := time.Time{}.Add(diff)

			m.lists[i].Styles.Title = titleStyle
			m.lists[i].Title = wticketid + " "
			switch {
			case diff.Minutes() > 10:
				m.lists[i].Title += stylered.Render(out.Format("04:05"))
			case diff.Minutes() > 5:
				m.lists[i].Title += styleyew.Render(out.Format("04:05"))
			default:
				m.lists[i].Title += styledef.Render(out.Format("04:05"))
			}

			xView := m.lists[i].View()
			if i == m.focused {
				xRender = append(xRender, focusedStyle.Render(xView))
			} else {
				xRender = append(xRender, columnStyle.Render(xView))
			}
			i++
			node = node.Next()
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
