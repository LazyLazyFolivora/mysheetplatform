package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) ListModules() ([]model.PermissionModule, error) {
	var modules []model.PermissionModule
	err := r.db.Order("sort_order ASC, id ASC").Find(&modules).Error
	return modules, err
}

func (r *PermissionRepo) GetModuleByID(id uint) (*model.PermissionModule, error) {
	var m model.PermissionModule
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *PermissionRepo) CreateModule(m *model.PermissionModule) error {
	return r.db.Create(m).Error
}

func (r *PermissionRepo) UpdateModule(m *model.PermissionModule) error {
	return r.db.Save(m).Error
}

// DeleteModules 删除多个模块及其权限关联和角色分配
func (r *PermissionRepo) DeleteModules(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("module_id IN ?", ids).Delete(&model.PermissionModuleRelation{}).Error; err != nil {
			return err
		}
		if err := tx.Where("module_id IN ?", ids).Delete(&model.RoleModule{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.PermissionModule{}, ids).Error
	})
}

func (r *PermissionRepo) ListPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.Order("id ASC").Find(&perms).Error
	return perms, err
}

func (r *PermissionRepo) CreatePermission(p *model.Permission) error {
	return r.db.Create(p).Error
}

func (r *PermissionRepo) UpdatePermission(p *model.Permission) error {
	return r.db.Save(p).Error
}

// DeletePermission 删除权限及其模块关联
func (r *PermissionRepo) DeletePermission(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("perm_id = ?", id).Delete(&model.PermissionModuleRelation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Permission{}, id).Error
	})
}

func (r *PermissionRepo) GetModulePermissionIDs(moduleID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.PermissionModuleRelation{}).
		Where("module_id = ?", moduleID).
		Pluck("perm_id", &ids).Error
	return ids, err
}

// SetModulePermissions 全量替换模块的权限分配
func (r *PermissionRepo) SetModulePermissions(moduleID uint, permIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("module_id = ?", moduleID).Delete(&model.PermissionModuleRelation{}).Error; err != nil {
			return err
		}
		if len(permIDs) == 0 {
			return nil
		}
		relations := make([]model.PermissionModuleRelation, 0, len(permIDs))
		for _, pid := range permIDs {
			relations = append(relations, model.PermissionModuleRelation{ModuleID: moduleID, PermID: pid})
		}
		return tx.Create(&relations).Error
	})
}

// SetPermissionModule 设置权限所属的模块（单模块归属，全量替换）
func (r *PermissionRepo) SetPermissionModule(permID uint, moduleID *uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("perm_id = ?", permID).Delete(&model.PermissionModuleRelation{}).Error; err != nil {
			return err
		}
		if moduleID == nil || *moduleID == 0 {
			return nil
		}
		return tx.Create(&model.PermissionModuleRelation{ModuleID: *moduleID, PermID: permID}).Error
	})
}

func (r *PermissionRepo) GetModulePermissions(moduleID uint) ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.Table("permission").
		Joins("JOIN permission_module_relation ON permission_module_relation.perm_id = permission.id").
		Where("permission_module_relation.module_id = ?", moduleID).
		Find(&perms).Error
	return perms, err
}
