package project

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	path_manager "github.com/yur4uwe/cmd-project-manager/manage_paths"
)

type Project struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Path        string `json:"Path"`
	TimeStamp   string `json:"TimeStamp"`
}

func CheckDuplicateNames(projects *[]Project, name string) bool {
	new_name := strings.ToLower(name)

	for i := 0; i < len(*projects); i++ {
		existing_name := strings.ToLower((*projects)[i].Name)
		if new_name == existing_name {
			return false
		}
	}

	return true
}

func ReadProjectsFromFile() []Project {
	log.Println("Read Projects From File")

	file, err := os.ReadFile(".projects.json")
	if err != nil {
		log.Println(err)
		return nil
	}

	var projects []Project
	err = json.Unmarshal(file, &projects)
	if err != nil {
		log.Println(err)
		return nil
	}

	return projects
}

func SaveProjects(projects *[]Project) {
	log.Println("Save Projects\n+-------------------+")

	for i := range *projects {
		(*projects)[i].ID = i
	}

	projectsJSON, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	err = os.WriteFile(".projects.json", projectsJSON, 0644)
	if err != nil {
		log.Println(err)
	}
}

func AddProject(projects *[]Project, name, description, path string) Project {
	log.Println("Add Project")

	if path[len(path)-1] != '/' {
		path += "/"
	}

	path = path + name

	new_project := Project{
		Name:        name,
		Description: description,
		Path:        path,
		TimeStamp:   time.Now().Format(time.RFC3339),
		ID:          len(*projects),
	}

	*projects = append(*projects, new_project)

	if info, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Println("Error while creating project directory: ", err)
			return Project{}
		}
	} else if err != nil {
		log.Println("Error while checking if project directory exists: ", err)
		return Project{}
	} else if !info.IsDir() {
		log.Println("Error: selected path exists but is not a directory")
		return Project{}
	}

	if info, err := os.Stat(path + "/.git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = path
		err = cmd.Run()

		if err != nil {
			log.Println("Error while initializing git repository")
			return Project{}
		}
	} else if err != nil {
		log.Println("Error while checking if .git directory exists: ", err)
		return Project{}
	} else if !info.IsDir() {
		log.Println("WTF?!!!??! WHY IS .git NOT A DIRECTORY???!?!?!?")
		return Project{}
	}

	path_manager.AddRecentPath(path)

	SaveProjects(projects)

	return new_project
}

func PrintProjectsSlice(projects []Project) string {
	delimiter := "+--------------------------------------+\n"
	var display_string string = "Projects:\n" + delimiter

	for _, project := range projects {
		display_string += PrintProjectInfo(project) + delimiter
	}

	return display_string
}

// PrintCompressedProjectsSlice prints a compressed version of the projects slice
// with only the ID, Name, and Path of each project.
// Starts with "Projects:\n"
func PrintCompressedProjectsSlice(projects []Project) string {
	var display_string string

	for _, project := range projects {
		display_string += fmt.Sprintf("ID: %d, Name: %s, Path: %s\n", project.ID, project.Name, project.Path)
	}

	return display_string
}

func PrintProjectInfo(project Project) string {
	var project_info = fmt.Sprintf("Project Info:\nID: %d\nName: %s\nDescription: %s\nPath: %s\nCreate Timestamp: %s\n",
		project.ID, project.Name, project.Description, project.Path, project.TimeStamp)

	return project_info
}

func RemoveProject(projects []Project, id int) []Project {
	log.Println("Remove Project From Slice By ID")

	err := os.RemoveAll(projects[id].Path)

	if err != nil {
		log.Println("Failed to remove project folder", err)
	}

	for i := range projects {
		if projects[i].ID == id {
			projects = append(projects[:i], projects[i+1:]...)
			break
		}
	}

	for i := range projects {
		projects[i].ID = i
	}

	path_manager.RemovePath(projects[id].Path)

	return projects
}

func UpdateProject(projects []Project, id int, name, description, path string) []Project {
	log.Println("Update Project By ID")

	for i, project := range projects {
		if project.ID == id {
			if name != "" {
				projects[i].Name = name
			}
			if description != "" {
				projects[i].Description = description
			}
			if path != "" {
				projects[i].Path = path
			}
			projects[i].TimeStamp = time.Now().Format(time.RFC3339)
			break
		}
	}

	return projects
}

func OpenProjectInExplorer(path string) {
	log.Println("Open Project In Explorer")

	cmd := exec.Command("start", path)
	err := cmd.Run()

	if err != nil {
		log.Fatal("Error while opening project in file explorer: ", err)
	}
}

func OpenProjectInVSCode(path string) {
	log.Println("Open Project In VSCode")

	cmd := exec.Command("code", path)
	err := cmd.Run()

	if err != nil {
		log.Fatal("Error while opening project in vs code: ", err)
	}
}

func CopyProjectPath(path string) {
	err := clipboard.WriteAll(path)
	if err != nil {
		log.Printf("Failed to write to clipboard: %v\n", err)
		panic(err)
	}
	log.Println("Successfully copied to clipboard!")
}
