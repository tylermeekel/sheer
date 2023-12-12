package fileselect

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylermeekel/sheer/internal/steps/style"
)

type fileSelectStyles struct {
	screenStyle   lipgloss.Style
	titleStyle    lipgloss.Style
	selectedStyle lipgloss.Style
	dirStyle      lipgloss.Style
	hintStyle     lipgloss.Style
}

type fileSelectStep struct {
	title string
	//numberOfOptions int
	files    []os.DirEntry
	filepath string
	selected int

	styles fileSelectStyles
}

func New(title string) *fileSelectStep {
	var styles fileSelectStyles

	styles.screenStyle = style.BaseScreenStyle
	styles.titleStyle = style.BaseTitleStyle
	styles.dirStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
	styles.selectedStyle = lipgloss.NewStyle().Foreground(style.AccentColor).Bold(true)
	styles.hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A9A9A9")).MarginTop(2)

	filepath := "."
	files, err := os.ReadDir(filepath)
	if err != nil {
		fmt.Println("Error reading files")
		os.Exit(1)
	}

	fs := fileSelectStep{
		title:    title,
		styles:   styles,
		files:    files,
		filepath: filepath,
	}

	return &fs
}

func (fs *fileSelectStep) next() {
	if fs.selected < len(fs.files)-1 {
		fs.selected++
	} else {
		fs.selected = 0
	}
}

func (fs *fileSelectStep) prev() {
	if fs.selected > 0 {
		fs.selected--
	} else {
		fs.selected = len(fs.files) - 1
	}
}

func (fs *fileSelectStep) Update(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			fs.next()
		case tea.KeyUp:
			fs.prev()
		case tea.KeyLeft:
			if fs.filepath != "." {
				fs.filepath = fs.filepath[:strings.LastIndex(fs.filepath, "/")]
				fs.selected = 0
				var err error
				fs.files, err = os.ReadDir(fs.filepath)
				if err != nil {
					fmt.Println("Error reading files")
					os.Exit(1)
				}
			}
		case tea.KeyEnter:
			file := fs.files[fs.selected]
			if file.IsDir() {
				fs.filepath += "/" + file.Name()
				fs.selected = 0
				var err error
				fs.files, err = os.ReadDir(fs.filepath)
				if err != nil {
					fmt.Println("Error reading files")
					os.Exit(1)
				}
			} else {
				fs.filepath += "/" + file.Name()
				return true
			}
		}
	}

	return false
}

func (fs *fileSelectStep) Render() string {
	var str string

	str += fs.styles.titleStyle.Render(fs.title)

	for i, file := range fs.files {
		str += "\n"
		if i == fs.selected {
			str += fs.styles.selectedStyle.Render("> " + file.Name())
		} else {
			if file.IsDir() {
				str += fs.styles.dirStyle.Render(file.Name())
			} else {
				str += file.Name()
			}
		}
	}

	if fs.filepath != "." {str += fs.styles.hintStyle.Render("(←) Go Back")}
	return fs.styles.screenStyle.Render(str)
}

func (fs *fileSelectStep) Result() string {
	return fs.filepath
}
