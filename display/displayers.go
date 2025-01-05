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

func MainMenu(projects []project.Project, key keyboard.Key, selected int) string {
	var display_string string = "Main Menu\n"

	options := []string{
		"Add Project",
		"Update Project",
		"Remove Project",
		"List Projects",
		"Exit",
	}

	for i, option := range options {
		if i == selected {
			display_string += fmt.Sprintf("> %s <\n", option)
		} else {
			display_string += fmt.Sprintf("  %s\n", option)
		}
	}

	return display_string
}

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

func UpdateProject(projects []project.Project) []project.Project {
	var selected = PrintCompressedProjectList(projects)

	if selected < 0 || selected >= len(projects) {
		return projects
	}

	fmt.Println(project.PrintProjectInfo(projects[selected]))

	return projects
}

func AddProject(projects *[]project.Project) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Add Project")

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	var recent_path_options = GetMostRecentPaths()

	var path_options = len(recent_path_options) + 1
	var selected = -1
	var path string = ""
	for {
		fmt.Println("Enter the absolute path to the project directory or choose already existing.")

		if selected == 0 {
			fmt.Println("> Use current directory <")
		}

		for i, option := range recent_path_options {
			if i == selected {
				fmt.Printf("> %s <\n", option)
			} else {
				fmt.Printf("  %s\n", option)
			}
		}

		fmt.Println(strings.Join(recent_path_options, "\n  "))
		fmt.Printf("Absolute Path: %s\n", path)

		char, key, err := keyboard.GetKey()

		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEnter {
			if isValidPath(path) {
				break
			}
			fmt.Print("\033[H\033[2J") // Clear the screen
			fmt.Println("Invalid path. Please enter a valid filesystem path.")
			path = ""
		} else if key == keyboard.KeyEsc {
			return
		} else if key == keyboard.KeyBackspace {
			if len(path) > 0 {
				path = path[:len(path)-1]
			}
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyArrowUp {
			selected = (selected - 1 + path_options) % path_options
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyArrowDown {
			selected = (selected + 1) % path_options
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else {
			path += string(char)
			fmt.Print("\033[H\033[2J") // Clear the screen
		}

	}

	project.AddProject(projects, name, description, path)
}
