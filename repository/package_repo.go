package repository

import (
	"log"
	"os"

	"github.com/bongnv/rpkgmgr/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Repository provides interface to work with storage layers.
//go:generate mockery -inpkg -name Repository -case=underscore
type Repository interface {
	ExistByNameAndVersion(name, version string) (bool, error)
	Insert(pkg *model.Package) error
	Update(pkg *model.Package) error
	Shutdown()
}

// NewRepository returns an implementation of Repository.
// Currently, it's based on MySQL as storage layer.
func NewRepository() Repository {
	db, err := gorm.Open("mysql", os.Getenv("GORM_URL"))
	if err != nil {
		log.Fatalln("Failed to initialize db, err: ", err)
		return nil
	}

	db.SingularTable(true)
	return &repositoryImpl{
		db: db,
	}
}

type repositoryImpl struct {
	db *gorm.DB
}

func (r *repositoryImpl) ExistByNameAndVersion(name, version string) (bool, error) {
	count := 0
	db := r.db.
		Model(&model.Package{}).
		Where(&model.Package{
			Name:    name,
			Version: version,
		}).
		Count(&count)

	if db.Error != nil {
		return false, db.Error
	}

	return count != 0, nil
}

func (r *repositoryImpl) Insert(pkg *model.Package) error {
	return r.db.Create(pkg).Error
}

func (r *repositoryImpl) Shutdown() {
	r.db.Close()
}

func (r *repositoryImpl) Update(pkg *model.Package) error {
	return r.db.Save(pkg).Error
}
