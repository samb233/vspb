package vspb

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// version control use a sqlite db to record plugin versions
type VersionControl struct {
	db *gorm.DB
}

func NewVersionControl(dbName string) (*VersionControl, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&PkgInfo{}); err != nil {
		return nil, err
	}

	return &VersionControl{
		db: db,
	}, nil
}

type PkgInfo struct {
	gorm.Model
	Name    string
	Version string
	Failed  bool
}

func (vc *VersionControl) GetPackage(name string) (*PkgInfo, error) {
	info := &PkgInfo{}
	if err := vc.db.Model(info).Where("name = ?", name).Find(info).Error; err != nil {
		return nil, err
	}
	return info, nil
}

func (vc *VersionControl) CreatePkgInfo(pkg *PkgInfo) error {
	return vc.db.Create(pkg).Error
}

func (vc *VersionControl) UpdatePkgInfo(id uint, update *PkgInfo) error {
	return vc.db.Model(&PkgInfo{}).Where("id = ?", id).Updates(update).Error
}

func (vc *VersionControl) DeletePkgInfo(id uint) error {
	return vc.db.Where("id = ?", id).Delete(&PkgInfo{}).Error
}

func (vc *VersionControl) GetAllPkgInfo() ([]*PkgInfo, error) {
	infos := make([]*PkgInfo, 0)
	if err := vc.db.Model(&PkgInfo{}).Find(&infos).Error; err != nil {
		return nil, err
	}

	return infos, nil
}
