package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) List() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *RoleRepo) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepo) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepo) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *RoleRepo) GetUserRoles(userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Table("role").
		Joins("JOIN user_role ON user_role.role_id = role.id").
		Where("user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) AssignUserRole(userID, roleID uint) error {
	ur := model.UserRole{UserID: userID, RoleID: roleID}
	return r.db.Create(&ur).Error
}

func (r *RoleRepo) RemoveUserRole(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}

func (r *RoleRepo) GetRoleModules(roleID uint) ([]model.RoleModule, error) {
	var rms []model.RoleModule
	err := r.db.Where("role_id = ?", roleID).Find(&rms).Error
	return rms, err
}

// SetRoleModules 全量替换角色的模块分配
func (r *RoleRepo) SetRoleModules(roleID uint, moduleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleModule{}).Error; err != nil {
			return err
		}
		if len(moduleIDs) == 0 {
			return nil
		}
		rms := make([]model.RoleModule, 0, len(moduleIDs))
		for _, mid := range moduleIDs {
			rms = append(rms, model.RoleModule{RoleID: roleID, ModuleID: mid})
		}
		return tx.Create(&rms).Error
	})
}
