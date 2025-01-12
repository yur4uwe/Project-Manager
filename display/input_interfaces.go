package display

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"
)

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
		} else {
			Clear()
			input += string(char)
		}
	}

	return input, nil
}

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
			path = ""
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
