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

// TODO: something is wrong again with opening folder in explorer
// TODO: Clear() issues

func main() {
	logFilePath := "log.txt"

	// Open the file with the os.O_TRUNC flag to clear its contents and set it for logging
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("Program start\n+-------------------+")

	var projects []project.Project = project.ReadProjectsFromFile()

	if err := keyboard.Open(); err != nil {
		log.Fatal("Error while opening the keyboard: ", err)
	}

	defer keyboard.Close()
	defer project.SaveProjects(&projects)

	display.Clear()

outerLoop:
	for {
		selected := display.MainMenu()

		switch selected {
		case TERMINATE, EXIT_PROGRAM:
			fmt.Println("Exiting...")
			break outerLoop
		case ADD_PROJECT:
			display.Clear()
			display.AddProjectInterface(&projects)
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
