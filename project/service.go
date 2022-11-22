package project

import (
	"errors"
	"reflect"

	"github.com/kataras/iris/v12"
)

type ProjectService interface {
	Create(params map[string]interface{}) (Project, int16, error)
	GetbyID(pid uint16) (Project, int16, error)
	Update(params map[string]interface{}) (Project, int16, error)
}

func NewProjectService(repo ProjectRepository) ProjectService {
	return &projectService{
		repo: repo,
	}
}

type projectService struct {
	repo ProjectRepository
}

func (p *projectService) Create(params map[string]interface{}) (Project, int16, error) {
	// check if projectName is already exist
	var projectDesignation, projectDescription, projectReference string
	projectName, ok := params["name"].(string)
	if !ok {
		return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
	}
	if projectName == "" {
		return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
	}
	projectRank := uint8(params["rank"].(float64))
	if _, ok := params["designation"]; ok {
		projectDesignation, ok = params["designation"].(string)
		if !ok {
			return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
		}
	}
	if _, ok := params["description"]; ok {
		projectDescription, ok = params["description"].(string)
		if !ok {
			return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
		}
	}
	if _, ok := params["reference"]; ok {
		projectReference, ok = params["reference"].(string)
		if !ok {
			return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
		}
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
		return Project{}, iris.StatusNotAcceptable, errors.New(iris.StatusText(iris.StatusNotAcceptable))
	}
	if projectDesignation != "" {
		if _, found := p.repo.Select(map[string]interface{}{"designation": projectDesignation}); found {
			return Project{}, iris.StatusNotAcceptable, errors.New(iris.StatusText(iris.StatusNotAcceptable))
		}
	}
	insertProject, err := p.repo.Insert(insertProject)
	if err != nil {
		return Project{}, iris.StatusInternalServerError, errors.New(iris.StatusText(iris.StatusInternalServerError))
	}
	return insertProject, iris.StatusCreated, nil
}

func (p *projectService) GetbyID(pid uint16) (Project, int16, error) {
	selectProject, found := p.repo.Select(map[string]interface{}{"p_id": pid})
	if found {
		return selectProject, iris.StatusOK, nil
	}
	return Project{}, iris.StatusNotFound, errors.New(iris.StatusText(iris.StatusNotFound))
}

func (p *projectService) Update(params map[string]interface{}) (Project, int16, error) {
	if params["p_id"] == nil || reflect.TypeOf(params["p_id"]).Kind() == reflect.String {
		return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
	}
	updatePID := uint16(params["p_id"].(float64))
	if _, found := p.repo.Select(map[string]interface{}{"p_id": updatePID}); !found {
		return Project{}, iris.StatusNotFound, errors.New(iris.StatusText(iris.StatusNotFound))
	}
	if params["name"] == "" || params["name"] == nil {
		return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
	}
	updateInfo, err := projectValidate(params)
	if err != nil {
		return Project{}, iris.StatusBadRequest, errors.New(iris.StatusText(iris.StatusBadRequest))
	}
	updateProject, err := p.repo.Update(Project{PID: updatePID}, updateInfo)
	if err != nil {
		return Project{}, iris.StatusInternalServerError, errors.New(iris.StatusText(iris.StatusInternalServerError))
	}
	return updateProject, iris.StatusOK, nil
}
