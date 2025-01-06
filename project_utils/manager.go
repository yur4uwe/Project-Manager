package project

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Project struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Path        string `json:"Path"`
	TimeStamp   string `json:"TimeStamp"`
}

func ReadProjectsFromFile() []Project {
	fmt.Println("Read Projects From File")

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
	fmt.Println("Save Projects")

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
	fmt.Println("Add Project")

	new_project := Project{
		Name:        name,
		Description: description,
		Path:        path,
		TimeStamp:   time.Now().Format(time.RFC3339),
		ID:          len(*projects),
	}

	*projects = append(*projects, new_project)

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
	fmt.Println("Remove Project From Slice By ID")

	for i, project := range projects {
		if project.ID == id {
			projects = append(projects[:i], projects[i+1:]...)
			break
		} else if project.ID > id {
			projects[i].ID--
		}
	}

	return projects
}

func UpdateProject(projects []Project, id int, name, description, path string) []Project {
	fmt.Println("Update Project By ID")

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
