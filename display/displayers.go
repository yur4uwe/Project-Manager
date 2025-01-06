package display

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
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

	fmt.Println(project.PrintProjectInfo(projects[selected]))

	fmt.Println("Press Enter to continue...")
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
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

	fmt.Println(project.PrintProjectInfo(projects[selected]))

	fmt.Println("Are you sure you want to remove this project? (y/n)")

	var char, key, err = keyboard.GetKey()

	if err != nil {
		log.Fatal(err)
	}

	if key == keyboard.KeyEnter || char == 'y' || char == 'Y' {
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

// (void) Mutates the projects slice
func AddProject(projects *[]project.Project) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Add Project")

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	path, err := getExecutablePath()
	if err != nil {
		log.Fatal(err)
	}

	path = PathChooser(path) + "/" + name

	project.AddProject(projects, name, description, path)
}