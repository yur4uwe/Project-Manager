package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
)

const (
	MAIN_MENU = iota
	ADD_PROJECT
	UPDATE_PROJECT
	REMOVE_PROJECT
	LIST_PROJECTS
)

func main() {
	var projects []project.Project = project.ReadProjectsFromFile()

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	defer project.SaveProjects(&projects)

	options := []string{
		"Add Project",
		"Update Project",
		"Remove Project",
		"List Projects",
		"Exit",
	}

	selected := 0

	var window_to_display int = MAIN_MENU

	var buffer = DisplayMainMenu(projects, keyboard.KeyBackspace, selected)

	fmt.Println(buffer)

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		buffer = DisplayMainMenu(projects, key, selected)

		if key == keyboard.KeyArrowDown && window_to_display == MAIN_MENU {
			selected = (selected + 1) % len(options)
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyArrowUp && window_to_display == MAIN_MENU {
			selected = (selected - 1 + len(options)) % len(options)
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyEnter {
			switch selected {
			case 0:
				window_to_display = ADD_PROJECT
			case 1:
				window_to_display = UPDATE_PROJECT
			case 2:
				window_to_display = REMOVE_PROJECT
			case 3:
				window_to_display = LIST_PROJECTS
			case 4:
				fmt.Println("Exiting...")
				return
			}
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyEsc {
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if char == 'q' || char == 'Q' {
			fmt.Println("Exiting...")
			return
		}

		switch window_to_display {
		case MAIN_MENU:
			buffer = DisplayMainMenu(projects, key, selected)
		case ADD_PROJECT:
			projects = DisplayAddProject(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		case UPDATE_PROJECT:
			projects = DisplayUpdateProject(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		case REMOVE_PROJECT:
			projects = DisplayRemoveProject(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		case LIST_PROJECTS:
			DisplayProjectsList(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		default:
			buffer = DisplayMainMenu(projects, key, selected)
		}

		fmt.Println(buffer)
	}
}

func DisplayMainMenu(projects []project.Project, key keyboard.Key, selected int) string {
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

func DisplayProjectsList(projects []project.Project) {
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

func DisplayRemoveProject(projects []project.Project) []project.Project {
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

func DisplayUpdateProject(projects []project.Project) []project.Project {
	var selected = PrintCompressedProjectList(projects)

	if selected < 0 || selected >= len(projects) {
		return projects
	}

	fmt.Println(project.PrintProjectInfo(projects[selected]))

	return projects
}

func DisplayAddProject(projects []project.Project) []project.Project {
	fmt.Println("Add Project")

	fmt.Print("Name: ")
	var name string
	fmt.Scanln(&name)

	fmt.Print("Description: ")
	var description string
	fmt.Scanln(&description)

	fmt.Print("Path: ")
	var path string
	fmt.Scanln(&path)

	project.AddProject(&projects, name, description, path)

	return projects
}

func PrintCompressedProjectList(projects []project.Project) int {
	var projects_slice = strings.Split(project.PrintCompressedProjectsSlice(projects), "\n")[:len(projects)]

	var joined string
	if len(projects_slice) == 0 {
		joined = ""
	} else {
		joined = "\n  " + strings.Join(projects_slice, "\n  ")
	}

	var display_string string = "Projects:  " + joined
	var selected = 0

	fmt.Println(display_string)

	if len(projects_slice) == 0 {
		fmt.Println("  No projects found.")
	}

	for {
		display_string = "Projects:\n"

		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowDown && len(projects_slice) > 0 {
			selected = (selected + 1) % len(projects_slice)
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyArrowUp && len(projects_slice) > 0 {
			selected = (selected - 1 + len(projects_slice)) % len(projects_slice)
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyEsc {
			return -1
		} else if key == keyboard.KeyEnter {
			return selected
		}

		for i, compressed_project_info := range projects_slice {
			if i == selected {
				display_string += fmt.Sprintf("> %s <\n", compressed_project_info)
			} else {
				display_string += fmt.Sprintf("  %s\n", compressed_project_info)
			}
		}

		if len(projects_slice) == 0 {
			display_string = "Projects:\n" + "  No projects found."
			fmt.Print("\033[H\033[2J") // Clear the screen
		}

		fmt.Println(display_string)
	}
}
