package main

import (
	"fmt"
	"log"

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

func main() {
	var projects []project.Project = project.ReadProjectsFromFile()

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

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
				project.SaveProjects(projects)
				return
			}
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyEsc {
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if char == 'q' || char == 'Q' {
			fmt.Println("Exiting...")
			project.SaveProjects(projects)
			return
		}

		switch window_to_display {
		case MAIN_MENU:
			buffer = DisplayMainMenu(projects, key, selected)
		case ADD_PROJECT:
			buffer = DisplayAddProject(projects, selected)
		case UPDATE_PROJECT:
			buffer = DisplayUpdateProject(projects, selected)
		case REMOVE_PROJECT:
			buffer = DisplayRemoveProject(projects, selected)
		case LIST_PROJECTS:
			buffer = DisplayListProjects(projects, selected)
		default:
			buffer = DisplayMainMenu(projects, key, selected)
		}

		fmt.Println(buffer)
	}
}

func DisplayListProjects(projects []project.Project, selected int) string {
	return project.PrintCompressedProjectsSlice(projects)
}

func DisplayRemoveProject(projects []project.Project, selected int) string {
	panic("unimplemented")
}

func DisplayUpdateProject(projects []project.Project, selected int) string {
	panic("unimplemented")
}

func DisplayAddProject(projects []project.Project, selected int) string {
	panic("unimplemented")
}
