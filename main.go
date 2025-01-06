package main

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	display "github.com/yur4uwe/cmd-project-manager/display"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
)

const (
	TERMINATE = iota - 2
	MAIN_MENU
	ADD_PROJECT
	UPDATE_PROJECT
	REMOVE_PROJECT
	LIST_PROJECTS
	EXIT_PROGRAM
)

// TODO: Add a way to open the project in the file explorer
// TODO: check if git repository exists before adding it
// TODO: Add a way to link the project to a existing folder (there is error at the moment)

func main() {
	var projects []project.Project = project.ReadProjectsFromFile()

	if err := keyboard.Open(); err != nil {
		log.Fatal("Error while opening the keyboard: ", err)
	}
	defer keyboard.Close()

	defer project.SaveProjects(&projects)
	for {
		selected := display.MainMenu()

		switch selected {
		case TERMINATE, EXIT_PROGRAM:
			fmt.Println("Exiting...")
			return
		case ADD_PROJECT:
			fmt.Print("\033[H\033[2J") // Clear the screen
			display.AddProject(&projects)
			fmt.Print("\033[H\033[2J") // Clear the screen
		case UPDATE_PROJECT:
			fmt.Print("\033[H\033[2J") // Clear the screen
			projects = display.UpdateProject(projects)
			fmt.Print("\033[H\033[2J") // Clear the screen
		case REMOVE_PROJECT:
			fmt.Print("\033[H\033[2J") // Clear the screen
			projects = display.RemoveProject(projects)
			fmt.Print("\033[H\033[2J") // Clear the screen
		case LIST_PROJECTS:
			fmt.Print("\033[H\033[2J") // Clear the screen
			display.ProjectsList(projects)
			fmt.Print("\033[H\033[2J") // Clear the screen
		default:
			fmt.Print("\033[H\033[2J") // Clear the screen
		}
	}
}
