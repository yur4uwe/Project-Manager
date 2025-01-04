package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
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

func DisplayAddProject(projects *[]project.Project) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Add Project")

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	var path string
	for {
		fmt.Println("Enter the absolute path to the project directory or choose already existing.")
		fmt.Println("  Use current directory")
		fmt.Print("Absolute Path: ")
		path, _ = reader.ReadString('\n')
		path = strings.TrimSpace(path)

		if isValidPath(path) {
			break
		} else {
			fmt.Println("Invalid path. Please enter a valid filesystem path.")
		}
	}

	project.AddProject(projects, name, description, path)
}
