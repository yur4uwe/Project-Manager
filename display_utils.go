package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/eiannone/keyboard"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
)

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

func isValidPath(path string) bool {
	// Define regex for Windows paths
	windowsPathPattern := `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	// Define regex for Unix-like paths
	unixPathPattern := `^(/[^/ ]*)+/?$`

	// Compile the regex patterns
	windowsPathRegex := regexp.MustCompile(windowsPathPattern)
	unixPathRegex := regexp.MustCompile(unixPathPattern)

	// Check if the path matches either pattern
	if !windowsPathRegex.MatchString(path) && !unixPathRegex.MatchString(path) {
		return false
	}

	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
