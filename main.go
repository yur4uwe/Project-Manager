package main

import (
	"fmt"
	"log"
	"os"

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

// TODO: Add a way to open the project directory
// TODO: check if git repository exists before adding it
// TODO: Add a way to link the project to a existing folder (there is error at the moment)

func main() {
	log_file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)

	log.Println("Program start\n+-------------------+")

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
			display.Clear()
			display.AddProject(&projects)
			display.Clear()
		case UPDATE_PROJECT:
			display.Clear()
			projects = display.UpdateProject(projects)
			display.Clear()
		case REMOVE_PROJECT:
			display.Clear()
			projects = display.RemoveProject(projects)
			display.Clear()
		case LIST_PROJECTS:
			display.Clear()
			display.ProjectsList(projects)
			display.Clear()
		default:
			display.Clear()
		}
	}
}
