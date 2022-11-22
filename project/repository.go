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
	log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	log.Printf("Select Project: %v", queryInfo)
	log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	result := database.Db.Model(&Project{}).Where(queryInfo).First(&p)
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
	result := database.Db.Model(&project).Updates(updateInfo)
	if result.Error != nil {
		return Project{}, result.Error
	}
	tx := database.Db.Model(&Project{}).Where("p_id = ?", project.PID).First(&updateProject)
	if tx.Error != nil {
		return Project{}, tx.Error
	}
	log.Println(updateProject)
	return
}
