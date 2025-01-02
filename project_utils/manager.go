package project

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/yur4uwe/cmd-project-manager/parse_json"
)

type Project struct {
	ID          int
	Name        string
	Description string
	Path        string
	TimeStamp   string
}

func ReadProjectsFromFile() []Project {
	fmt.Println("Read Projects From File")

	file, err := os.ReadFile(".projects.json")

	if err != nil {
		log.Println(err)
		return nil
	}

	projectsJSON := string(file)
	projects := []Project{}

	// fmt.Println("Projects on json:", projectsJSON)

	for i := 0; i < len(projectsJSON); i++ {
		if projectsJSON[i] == '{' {
			var j int = i + 1
			for projectsJSON[j] != '}' {
				j++
			}

			var projectJSON = projectsJSON[i : j+1]
			var projectAsJsonMap = parse_json.Parse(projectJSON)
			var project = CastJSONToProject(projectAsJsonMap)
			projects = append(projects, project)

			i = j + 1
		}
	}

	return projects
}

func SaveProjects(projects []Project) {
	fmt.Println("Save Projects")

	projectsJSON := "["

	for i, project := range projects {
		projectsJSON += fmt.Sprintf(`{"ID": %d, "Name": "%s", "Description": "%s", "Path": "%s", "TimeStamp": "%s"}`, project.ID, project.Name, project.Description, project.Path, project.TimeStamp)
		if i != len(projects)-1 {
			projectsJSON += ","
		}
	}

	projectsJSON += "]"

	err := os.WriteFile(".projects.json", []byte(projectsJSON), 0644)

	if err != nil {
		log.Println(err)
	}
}

func AddProject(prjs *[]Project, name, description, path string) Project {
	fmt.Println("Add Project")

	projects := *prjs

	prj := Project{
		Name:        name,
		Description: description,
		Path:        path,
		TimeStamp:   time.Now().Format(time.RFC3339),
		ID:          len(projects),
	}

	*prjs = append(projects, prj)

	SaveProjects(*prjs)

	return prj
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

func CastJSONToProject(json map[string]string) Project {
	var id int
	var err error
	var name, description, path, timestamp string

	for key, value := range json {
		switch key {
		case "ID":
			id, err = strconv.Atoi(value)
		case "Name":
			name = value
		case "Description":
			description = value
		case "Path":
			path = value
		case "TimeStamp":
			timestamp = value
		}
	}

	if err != nil {
		log.Println(err)
		return Project{}
	}

	return Project{
		ID:          id,
		Name:        name,
		Description: description,
		Path:        path,
		TimeStamp:   timestamp,
	}
}
