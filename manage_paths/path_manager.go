package path_manager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type RecentPath struct {
	Path        string `json:"Path"`
	LastAccess  string `json:"LastAccess"`
	TimesOpened int    `json:"TimesOpened"`
}

func ReadRecentPathsFromFile() ([]RecentPath, error) {
	fmt.Println("Read Recent Paths From File")

	file, err := os.ReadFile(".directory_history.json")
	if err != nil {
		log.Println("Error while reading path history file: ", err)
		return nil, fmt.Errorf("os: path history file doesn't exist:\n %w", err)
	}

	var recent_paths []RecentPath
	err = json.Unmarshal(file, &recent_paths)
	if err != nil {
		log.Printf("Error while unmarshaling JSON file: %v\nJSON data: %s\n", err, string(file))
		return nil, fmt.Errorf("json: failed to unmarshal json file:\n %w", err)
	}

	return recent_paths, nil
}

func AddRecentPath(path string) {
	fmt.Println("Add Recent Path")

	recent_paths, err := ReadRecentPathsFromFile()
	if err != nil {
		log.Println("read path history: failed to read path history\n", err)
		return
	}

	var new_path = RecentPath{
		Path:        path,
		LastAccess:  time.Now().Format(time.RFC3339),
		TimesOpened: 1,
	}

	recent_paths = append(recent_paths, new_path)

	recent_paths = RemoveDuplicatePaths(recent_paths)

	SaveRecentPaths(recent_paths)
}

func SaveRecentPaths(recent_paths []RecentPath) {
	fmt.Println("Save Recent Paths")

	recentPathsJSON, err := json.MarshalIndent(recent_paths, "", "  ")
	if err != nil {
		log.Println("Error while marshaling JSON file: ", err)
		return
	}

	err = os.WriteFile(".directory_history.json", recentPathsJSON, 0644)
	if err != nil {
		log.Println("Error while writing to JSON file: ", err)
	}
}

func RemoveDuplicatePaths(recent_paths []RecentPath) []RecentPath {
	fmt.Println("Remove Duplicate Paths")

	var new_recent_paths []RecentPath

	for _, recent_path := range recent_paths {
		var found bool = false

		for _, new_recent_path := range new_recent_paths {
			if recent_path.Path == new_recent_path.Path {
				found = true
				break
			}
		}

		if !found {
			new_recent_paths = append(new_recent_paths, recent_path)
		}
	}

	return new_recent_paths
}

func IncrementAccess(path string) {
	fmt.Println("Increment Access")

	recent_paths, err := ReadRecentPathsFromFile()
	if err != nil {
		log.Println("read path history: failed to read path history\n", err)
		return
	}

	for i, recent_path := range recent_paths {
		if recent_path.Path == path {
			recent_paths[i].TimesOpened++
			recent_paths[i].LastAccess = time.Now().Format(time.RFC3339)
			break
		}
	}

	SaveRecentPaths(recent_paths)
}

func GetMostRecentPaths() []string {
	fmt.Println("Get Most Recent Paths")

	recent_paths, err := ReadRecentPathsFromFile()
	if err != nil {
		log.Println("read path history: failed to read path history\n", err)
		return nil
	}

	// Sort the recent paths by the last access time
	for i := 0; i < len(recent_paths); i++ {
		for j := i + 1; j < len(recent_paths); j++ {
			if recent_paths[i].LastAccess < recent_paths[j].LastAccess {
				recent_paths[i], recent_paths[j] = recent_paths[j], recent_paths[i]
			}
		}
	}

	// Get the 5 most recent paths
	if len(recent_paths) > 5 {
		recent_paths = recent_paths[:5]
	}

	var recent_paths_strings []string

	for _, recent_path := range recent_paths {
		recent_paths_strings = append(recent_paths_strings, recent_path.Path)
	}

	return recent_paths_strings
}

func RemovePath(path string) []string {

	recent_paths := GetMostRecentPaths()

	for i := 0; i < len(recent_paths); i++ {
		if recent_paths[i] == path {
			recent_paths = append(recent_paths[:i], recent_paths[i+1:]...)
		}
	}

	return recent_paths
}
