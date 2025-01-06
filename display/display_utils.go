package display

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
)

func getExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, err := filepath.Abs(filepath.Dir(ex))

	return strings.ReplaceAll(dir, "\\", "/"), err
}

func PrintCompressedProjectList(projects []project.Project) int {
	var projects_slice = strings.Split(project.PrintCompressedProjectsSlice(projects), "\n")[:len(projects)]

	defer fmt.Print("\033[H\033[2J") // Clear the screen

	return ChoiceMenu(projects_slice, "Projects:\n", "  No projects found.")
}

func isValidPath(path string) bool {
	// // Define regex for Windows paths
	// windowsPathPattern := `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	// // Define regex for Unix-like paths
	// unixPathPattern := `^(/[^/ ]*)+/?$`

	// // Compile the regex patterns
	// windowsPathRegex := regexp.MustCompile(windowsPathPattern)
	// unixPathRegex := regexp.MustCompile(unixPathPattern)

	// // Check if the path matches either pattern
	// if !windowsPathRegex.MatchString(path) && !unixPathRegex.MatchString(path) {
	// 	return false
	// }

	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func GetMostRecentPaths() []string {

	// Get the most recent paths from the .directory_history.json file
	// and return them as a slice of strings
	// If the file does not exist, return an empty slice
	// If the file exists but is empty, return an empty slice
	// If the file exists and contains paths, return them as a slice of strings

	// Read the file
	file, err := os.ReadFile(".directory_history.json")
	if err != nil {
		log.Println(err)
		return []string{}
	}

	// Unmarshal the file
	var paths []string
	err = json.Unmarshal(file, &paths)
	if err != nil {
		log.Println(err)
		return []string{}
	}

	// Return the 5 most recent paths
	if len(paths) > 5 {
		return paths[:5]
	}
	return paths
}

/*
ChoiceMenu displays a menu with the given options and returns the index of the selected option.

Returns:

index of the selected option if the user presses the Enter key.

0 is the default option if the user presses the Enter key without selecting an option

-1 if the user presses the ESC key.

-2 if the user presses a key from the termination_options slice.
*/
func ChoiceMenu(options []string, header string, no_options string, termination_options ...string) int {
	selected := 0

	fmt.Print(header)
	for i, option := range options {
		if i == selected {
			fmt.Printf("> %s <\n", option)
		} else {
			fmt.Printf("  %s\n", option)
		}
	}

	if len(options) == 0 {
		fmt.Println(no_options)
	}

	for {
		display_string := header

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowDown {
			selected = (selected + 1) % len(options)
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyArrowUp {
			selected = (selected - 1 + len(options)) % len(options)
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else if key == keyboard.KeyEnter {
			return selected
		} else if key == keyboard.KeyEsc {
			return -1
		} else if strings.Contains(strings.Join(termination_options, ""), string(char)) {
			return -2
		} else {
			fmt.Print("\033[H\033[2J") // Clear the screen
		}

		for i, option := range options {
			if i == selected {
				display_string += fmt.Sprintf("> %s <\n", option)
			} else {
				display_string += fmt.Sprintf("  %s\n", option)
			}
		}

		if len(options) == 0 {
			display_string = header + no_options
		}

		fmt.Println(display_string)
	}
}

func MatchFoldersInPath(valid_path string, name_to_match string) []string {
	// List all folders in the given path
	// Return a slice of strings with the folder names

	var folders []string

	// Open the directory
	dir, err := os.ReadDir(valid_path)
	if err != nil {
		log.Println(err)
		return folders
	}

	// Read the directory
	for _, entry := range dir {
		if entry.IsDir() && strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(name_to_match)) && !strings.HasPrefix(entry.Name(), ".") {
			folders = append(folders, entry.Name())
		}
	}

	if len(folders) == 0 {
		folders = append(folders, "No folders found.")
	} else if len(folders) > 5 {
		folders = folders[:5]
	}
	return folders
}

func PathChooser(current_path string) string {
	var recent_path_options = GetMostRecentPaths()

	fmt.Println("Enter the absolute path to the project directory or choose already existing.")

	var path_options = len(recent_path_options) + 1
	var selected = -1
	var path string = current_path

	for {

		if selected == 0 {
			fmt.Println("> Use current directory <")
		} else {
			fmt.Println("  Use current directory")
		}

		for i, option := range recent_path_options {
			if i == selected {
				fmt.Printf("> %s <\n", option)
			} else {
				fmt.Printf("  %s\n", option)
			}
		}

		fmt.Printf("Absolute Path: %s\n", path)

		split_path := strings.Split(path, "/")

		folders := MatchFoldersInPath(strings.Join(split_path[:len(split_path)-1], "/"), split_path[len(split_path)-1])

		fmt.Println("\n  ", strings.Join(folders, "\n  "))

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
			return ""
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
		} else if key == keyboard.KeyTab {
			if len(folders) != 1 {
				continue
			}
			split_path = strings.Split(path, "/")
			path = strings.Join(split_path[:len(split_path)-1], "/") + "/" + folders[0] + "/"
			fmt.Print("\033[H\033[2J") // Clear the screen
		} else {
			path += string(char)
			fmt.Print("\033[H\033[2J") // Clear the screen
		}
	}

	return path
}
