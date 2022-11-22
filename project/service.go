package project

import (
	"errors"
	"reflect"
)

type ProjectService interface {
	Create(params map[string]interface{}) (Project, error)
	GetbyID(pid uint16) (Project, error)
	Update(params map[string]interface{}) (Project, error)
}

func NewProjectService(repo ProjectRepository) ProjectService {
	return &projectService{
		repo: repo,
	}
}

type projectService struct {
	repo ProjectRepository
}

func (p *projectService) Create(params map[string]interface{}) (Project, error) {
	// check if projectName is already exist
	var projectDesignation, projectDescription, projectReference string
	projectName, ok := params["name"].(string)
	if !ok {
		return Project{}, errors.New("parameter name is not exist")
	}
	if projectName == "" {
		return Project{}, errors.New("parameter name invalid")
	}
	projectRank := uint8(params["rank"].(float64))
	if _, ok := params["designation"]; ok {
		projectDesignation = params["designation"].(string)
	}
	if _, ok := params["description"]; ok {
		projectDescription = params["description"].(string)
	}
	if _, ok := params["reference"]; ok {
		projectReference = params["reference"].(string)
	}
	insertProject := Project{
		Name:        projectName,
		Rank:        projectRank,
		Designation: projectDesignation,
		Description: projectDescription,
		Reference:   projectReference,
		Status:      1,
	}
	// construct Project struct into Insert operation
	if _, found := p.repo.Select(map[string]interface{}{"name": projectName}); found {
		return Project{}, errors.New("project name is already exist")
	}
	if projectDesignation != "" {
		if _, found := p.repo.Select(map[string]interface{}{"designation": projectDesignation}); found {
			return Project{}, errors.New("dumplicated project designation")
		}
	}
	insertProject, err := p.repo.Insert(insertProject)
	if err != nil {
		return Project{}, err
	}
	return insertProject, nil
}

func (p *projectService) GetbyID(pid uint16) (Project, error) {
	selectProject, found := p.repo.Select(map[string]interface{}{"p_id": pid})
	if found {
		return selectProject, nil
	}
	return Project{}, errors.New("project is not exist")
}

func (p *projectService) Update(params map[string]interface{}) (Project, error) {
	if params["p_id"] == nil || reflect.TypeOf(params["p_id"]).Kind() == reflect.String {
		return Project{}, errors.New("parameter error")
	}
	updatePID := uint16(params["p_id"].(float64))
	if _, found := p.repo.Select(map[string]interface{}{"p_id": updatePID}); !found {
		return Project{}, errors.New("project is not exist")
	}
	if err := projectValidate(params); err != nil {
		return Project{}, err
	}
	if params["name"] == "" || params["name"] == nil {
		return Project{}, errors.New("parameter name can not be null")
	}
	// log.Println("_____________________________________")
	// log.Println(project.Name)
	// updatePID := project.PID
	// projectMap, err := utils.ToMap(project, "json")
	// if err != nil {
	// 	return Project{}, err
	// }
	// if err := projectValidate(projectMap); err != nil {
	// 	return Project{}, err
	// }
	// updateProject, err = p.repo.Update(Project{PID: updatePID}, projectMap)
	return Project{}, nil
}
