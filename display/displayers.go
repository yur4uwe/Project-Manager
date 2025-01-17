package display

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	path_manager "github.com/yur4uwe/cmd-project-manager/manage_paths"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
)

// (int) Returns the index of the selected option
func MainMenu() int {
	var display_string string = "Main Menu\n"

	options := []string{
		"Add Project",
		"Update Project",
		"Remove Project",
		"List Projects",
		"Exit",
	}

	return ChoiceMenu(options, display_string, "", "Q", "q")
}

// (void) Lists Projects
func ProjectsList(projects []project.Project) {
	var selected = PrintCompressedProjectList(projects)

	if selected < 0 || selected >= len(projects) {
		return
	}

	Clear()

	header := project.PrintProjectInfo(projects[selected]) + "\nProject Options\n"
	options := []string{"Open in VS Code", "Open in File Explorer", "Copy Path to Projects Directory", "Back"}

	do_next := ChoiceMenu(options, header, "", "B", "b")

	switch do_next {
	case -1, -2, 3:
		return
	case 0:
		path_manager.IncrementAccess(projects[selected].Path)
		project.OpenProjectInVSCode(projects[selected].Path)
	case 1:
		project.OpenProjectInExplorer(projects[selected].Path)
	case 2:
		project.CopyProjectPath(projects[selected].Path)
	}

	fmt.Println("Press Enter to continue...")
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal("Error while getting keyboard key: ", err)
		}
		if key == keyboard.KeyEnter || key == keyboard.KeyEsc {
			break
		}
	}
}

// Returns mutated projects slice
func RemoveProject(projects []project.Project) []project.Project {
	var selected = PrintCompressedProjectList(projects)

	if selected < 0 || selected >= len(projects) {
		return projects
	}

	buffer := project.PrintProjectInfo(projects[selected]) +
		"\nAre you sure you want to remove this project? (y/n)\n" +
		"\n**This action only deletes link to a project" +
		"\nand not the directory with it**"

	fmt.Println(buffer)

	var char, key, err = keyboard.GetKey()

	if err != nil {
		log.Fatal("Error while getting keyboard key: ", err)
	}

	if key == keyboard.KeyEnter || char == 'y' || char == 'Y' {
		path_manager.RemovePath(projects[selected].Path)
		projects = append(projects[:selected], projects[selected+1:]...)
		return projects
	}

	return projects
}

// Returns mutated projects slice
func UpdateProject(projects []project.Project) []project.Project {
	var selected = PrintCompressedProjectList(projects)

	if selected < 0 || selected >= len(projects) {
		return projects
	}

	fmt.Println(project.PrintProjectInfo(projects[selected]))

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Update Project fields(leave empty to keep current value):")

	fmt.Printf("Old Name: %s\n", projects[selected].Name)
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	if strings.TrimSpace(name) != "" {
		name = strings.TrimSpace(name)
		projects[selected].Name = name
	}

	fmt.Printf("Old Description: %s\n", projects[selected].Description)
	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	if strings.TrimSpace(description) != "" {
		description = strings.TrimSpace(description)
		projects[selected].Description = description
	}

	return projects
}

func CreateNewProject(projects *[]project.Project) {
	header := "Create Project\n"

	header += "Name: "
	var name string
	var err error
	for {
		name, err = readInputWithCancel(header, keyboard.KeyEsc)
		if err != nil {
			return
		}
		name = strings.TrimSpace(name)
		if project.CheckDuplicateNames(projects, name) {

			break
		}
		header = "Invalid name\nAdd Project\nName: "
	}

	header = "Create Project\nName: " + name + "\nDescription: "
	description, err := readInputWithCancel(header, keyboard.KeyEsc)
	if err != nil {
		return
	}
	description = strings.TrimSpace(description)
	header += description + "\n"

	path, err := getExecutablePath()
	if err != nil {
		log.Fatal("Error while getting executable path", err)
	}

	path = PathChooser(header, path)

	if path != "" {
		project.AddProject(projects, name, description, path)
	}
}

func LinkProject(projects *[]project.Project) {
	header := "Navigate to the project directory"

	path, err := getExecutablePath()
	if err != nil {
		log.Fatal("Error while getting executable path", err)
		path = "~/"
	}

	path = PathChooser(header, path)

	if path == "" {
		return
	}

	name := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]

	header = fmt.Sprintf("Linking project\nName: %v", name)
	description, err := readInputWithCancel(header)
	if err != nil {
		return
	}

	project.AddProject(projects, name, description, path)
}
