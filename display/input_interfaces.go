package display

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"
	project "github.com/yur4uwe/cmd-project-manager/project_utils"
)

/*
ChoiceMenu displays a menu with the given options and returns the index of the selected option.

Parameters:
- options: A slice of strings representing the menu options.
- header: A string to display as the header for the menu.
- no_options: A string to display if there are no options available.
- termination_options: A variadic list of characters that can be used to terminate the menu.

Returns:
- int: The index of the selected option if the Enter key is pressed.
- -1: If the ESC key is pressed.
- -2: If a key from the termination_options slice is pressed.
*/
func ChoiceMenu(options []string, header string, no_options string, termination_options ...string) int {
	selected := 0

	for {
		display_string := header

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

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal("Error while getting keyboard key: ", err)
		}

		if key == keyboard.KeyArrowDown {
			selected = (selected + 1) % len(options)
			Clear()
		} else if key == keyboard.KeyArrowUp {
			selected = (selected - 1 + len(options)) % len(options)
			Clear()
		} else if key == keyboard.KeyEnter {
			return selected
		} else if key == keyboard.KeyEsc {
			return -1
		} else if strings.Contains(strings.Join(termination_options, ""), string(char)) {
			return -2
		} else {
			Clear()
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

/*
readInputWithCancel reads input from the user and allows canceling with the 'ESC' key.

Parameters:
- header: A string to display as the header for the input prompt.
- escape_keys: A variadic list of keys that can be used to cancel the input.

Returns:
- string: The validated input if the Enter key is pressed and the input is valid.
- error: An error if the input is canceled or if there is an issue with getting the keyboard input.

Possible Returns:
- If the input is canceled by pressing one of the escape keys, the function returns an empty string and an error with the message "input cancelled".
- If the Enter key is pressed and the input is valid according to the predicate, the function returns the input string and nil error.
- If there is an error while getting the keyboard input, the function returns an empty string and the error.
*/
func readInputWithCancel(header string, escape_keys ...keyboard.Key) (string, error) {
	defer Clear()
	input := ""

	for {
		fmt.Println(header, input)
		char, key, err := keyboard.GetKey()
		if err != nil {
			return "", err
		}

		if arrayContainsAtLeastOneKey(escape_keys, key) {
			return "", fmt.Errorf("input cancelled")
		} else if key == keyboard.KeyEnter {
			break
		} else if key == keyboard.KeyBackspace {
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
			Clear()
		} else if key == keyboard.KeySpace {
			input += " "
		} else {
			Clear()
			input += string(char)
		}
	}

	return input, nil
}

/*
PathChooser allows the user to choose a path by navigating through the filesystem.

Parameters:
- header: A string to display as the header for the path chooser.
- current_path: The current path to start from.

Returns:
- string: The chosen path if the Enter key is pressed and the path is valid.
- If the ESC key is pressed, the function returns an empty string.
*/
func PathChooser(header string, current_path string) string {
	var recent_path_options = GetMostRecentPaths()

	header += "\nEnter the absolute path to the project directory or choose already existing."

	var path_options = len(recent_path_options) + 1
	var selected = -1
	var path string = current_path

	for {
		fmt.Println(header)

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

		fmt.Printf("\nAbsolute Path: %s\n", path)

		split_path := strings.Split(path, "/")

		folders := MatchFoldersInPath(strings.Join(split_path[:len(split_path)-1], "/"), split_path[len(split_path)-1])

		fmt.Println("  ", strings.Join(folders, "\n  "))

		char, key, err := keyboard.GetKey()

		if err != nil {
			log.Fatal("Error getting keyboard key: ", err)
		}

		if key == keyboard.KeyEnter {
			if isValidPath(path) {
				break
			}
			Clear()
			fmt.Println("Invalid path. Please enter a valid filesystem path.")
		} else if key == keyboard.KeyEsc {
			return ""
		} else if key == keyboard.KeyBackspace {
			if len(path) > 0 {
				path = path[:len(path)-1]
			}
			Clear()
		} else if key == keyboard.KeyArrowUp {
			selected = (selected - 1 + path_options) % path_options
			Clear()
		} else if key == keyboard.KeyArrowDown {
			selected = (selected + 1) % path_options
			Clear()
		} else if key == keyboard.KeyTab {
			if len(folders) != 1 {
				continue
			}
			split_path = strings.Split(path, "/")
			path = strings.Join(split_path[:len(split_path)-1], "/") + "/" + folders[0] + "/"
			Clear()
		} else {
			path += string(char)
			Clear()
		}
	}

	return path
}

/*
AddProjectInterface provides an interface for adding a new project or linking an existing project.

Parameters:
- projects: A pointer to a slice of Project structs.

Returns:
- void: This function mutates the projects slice.
*/
func AddProjectInterface(projects *[]project.Project) {
	add_options := []string{
		"Create new Project", "Link an already created project",
	}

	option := ChoiceMenu(add_options, "", "")

	switch option {
	case -1:
		return
	case 0:
		CreateNewProject(projects)
	case 1:
		LinkProject(projects)
	}
}
