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

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func getExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		log.Println("Error while getting executable path: ", err)
		return "", err
	}

	dir, err := filepath.Abs(filepath.Dir(ex))

	return strings.ReplaceAll(dir, "\\", "/"), err
}

func PrintCompressedProjectList(projects []project.Project) int {
	var projects_slice = strings.Split(project.PrintCompressedProjectsSlice(projects), "\n")[:len(projects)]

	defer Clear()

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
		log.Println("GetMostRecentPaths: Error while reading path history file:", err)
		return []string{}
	}

	// Unmarshal the file
	var paths []string
	err = json.Unmarshal(file, &paths)
	if err != nil {
		log.Println("GetMostRecentPaths: Error while unmarshaling JSON file", err)
		return []string{}
	}

	// Return the 5 most recent paths
	if len(paths) > 5 {
		return paths[:5]
	}
	return paths
}

func MatchFoldersInPath(valid_path string, name_to_match string) []string {
	// List all folders in the given path
	// Return a slice of strings with the folder names

	var folders []string

	// Open the directory
	dir, err := os.ReadDir(valid_path)
	if err != nil {
		log.Println("Error while reading directory to match folders: ", err)
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

func arrayContainsAtLeastOneKey(array []keyboard.Key, args ...keyboard.Key) bool {
	for i := 0; i < len(array); i++ {
		for j := 0; j < len(args); j++ {
			if array[i] == args[j] {
				return true
			}
		}
	}
	return false
}
