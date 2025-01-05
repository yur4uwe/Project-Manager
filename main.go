package main

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	display "github.com/yur4uwe/cmd-project-manager/display"
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

	var buffer = display.MainMenu(projects, keyboard.KeyBackspace, selected)

	fmt.Println(buffer)

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		buffer = display.MainMenu(projects, key, selected)

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
			buffer = display.MainMenu(projects, key, selected)
		case ADD_PROJECT:
			display.AddProject(&projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		case UPDATE_PROJECT:
			projects = display.UpdateProject(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		case REMOVE_PROJECT:
			projects = display.RemoveProject(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		case LIST_PROJECTS:
			display.ProjectsList(projects)
			window_to_display = MAIN_MENU
			fmt.Print("\033[H\033[2J") // Clear the screen
		default:
			buffer = display.MainMenu(projects, key, selected)
		}

		fmt.Println(buffer)
	}
}
