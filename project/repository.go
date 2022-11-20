package project

import (
	database "icarus/database/mariadb"
	"log"
)

type ProjectRepository interface {
	Select(queryInfo map[string]interface{}) (selectProject Project, found bool)
	Insert(project Project) (insertProject Project, err error)
	Update(project Project, updateInfo map[string]interface{}) (updateProject Project, err error)
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

type projectRepository struct {
}

func (r *projectRepository) Select(queryInfo map[string]interface{}) (p Project, found bool) {
	result := database.Db.Model(&Project{}).Where(queryInfo).First(&p)
	log.Println("________________________________")
	log.Println(p)
	if result.Error == nil {
		found = true
	}
	return
}

func (r *projectRepository) Insert(project Project) (Project, error) {
	result := database.Db.Model(&Project{}).Create(&project)
	if result.Error != nil {
		return Project{}, result.Error
	}
	return project, nil
}

func (r *projectRepository) Update(project Project, updateInfo map[string]interface{}) (updateProject Project, err error) {
	if result := database.Db.Model(&project).First(&updateProject); result.Error != nil {
		return Project{}, result.Error
	}
	log.Println("___________________________")
	log.Println(updateProject)
	tx := database.Db.Model(&updateProject).Updates(updateInfo)
	if tx.Error != nil {
		return Project{}, nil
	}
	return
}
