package service

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/repository"
)

type PermissionService struct {
	permRepo *repository.PermissionRepo
	roleRepo *repository.RoleRepo
	logger   *zap.Logger
}

func NewPermissionService(permRepo *repository.PermissionRepo, roleRepo *repository.RoleRepo, logger *zap.Logger) *PermissionService {
	return &PermissionService{permRepo: permRepo, roleRepo: roleRepo, logger: logger}
}

type ModuleTreeNode struct {
	model.PermissionModule
	Children []*ModuleTreeNode `json:"children"`
}

func (s *PermissionService) BuildTree() ([]*ModuleTreeNode, error) {
	modules, err := s.permRepo.ListModules()
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[uint]*ModuleTreeNode)
	for i := range modules {
		nodeMap[modules[i].ID] = &ModuleTreeNode{PermissionModule: modules[i], Children: nil}
	}

	var roots []*ModuleTreeNode
	for _, node := range nodeMap {
		if node.ParentID == nil || *node.ParentID == 0 {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[*node.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

func (s *PermissionService) CheckAccess(userID uint, moduleCode string) (bool, error) {
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	modules, err := s.permRepo.ListModules()
	if err != nil {
		return false, err
	}

	var targetModule *model.PermissionModule
	for _, m := range modules {
		if m.Path == moduleCode {
			targetModule = &m
			break
		}
	}
	if targetModule == nil {
		return false, nil
	}

	for _, role := range roles {
		roleModules, err := s.roleRepo.GetRoleModules(role.ID)
		if err != nil {
			continue
		}
		for _, rm := range roleModules {
			if rm.ModuleID == targetModule.ID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (s *PermissionService) CreateModule(m *model.PermissionModule) error {
	return s.permRepo.CreateModule(m)
}

func (s *PermissionService) GetModule(id uint) (*model.PermissionModule, error) {
	return s.permRepo.GetModuleByID(id)
}

func (s *PermissionService) UpdateModule(m *model.PermissionModule) error {
	return s.permRepo.UpdateModule(m)
}

// DeleteModule 删除模块及其所有子孙模块（连同权限关联、角色分配）
func (s *PermissionService) DeleteModule(id uint) error {
	modules, err := s.permRepo.ListModules()
	if err != nil {
		return err
	}

	childrenMap := make(map[uint][]uint)
	for _, m := range modules {
		if m.ParentID != nil && *m.ParentID != 0 {
			childrenMap[*m.ParentID] = append(childrenMap[*m.ParentID], m.ID)
		}
	}

	ids := []uint{id}
	for i := 0; i < len(ids); i++ {
		ids = append(ids, childrenMap[ids[i]]...)
	}

	return s.permRepo.DeleteModules(ids)
}

// Role methods
func (s *PermissionService) ListRoles() ([]model.Role, error) {
	return s.roleRepo.List()
}

func (s *PermissionService) CreateRole(name, description string) (*model.Role, error) {
	role := &model.Role{Name: name, Description: description}
	err := s.roleRepo.Create(role)
	return role, err
}

func (s *PermissionService) UpdateRole(id uint, name, description string) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	role.Name = name
	role.Description = description
	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *PermissionService) DeleteRole(id uint) error {
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}
	return s.roleRepo.Delete(id)
}

func (s *PermissionService) GetRoleModuleIDs(roleID uint) ([]uint, error) {
	rms, err := s.roleRepo.GetRoleModules(roleID)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rms))
	for _, rm := range rms {
		ids = append(ids, rm.ModuleID)
	}
	return ids, nil
}

func (s *PermissionService) SetRoleModules(roleID uint, moduleIDs []uint) error {
	if _, err := s.roleRepo.GetByID(roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}
	return s.roleRepo.SetRoleModules(roleID, moduleIDs)
}

// Permission methods
func (s *PermissionService) ListPermissions() ([]model.Permission, error) {
	return s.permRepo.ListPermissions()
}

func (s *PermissionService) CreatePermission(p *model.Permission, moduleID *uint) error {
	if err := s.permRepo.CreatePermission(p); err != nil {
		return err
	}
	if moduleID != nil && *moduleID != 0 {
		return s.permRepo.SetPermissionModule(p.ID, moduleID)
	}
	return nil
}

func (s *PermissionService) UpdatePermission(p *model.Permission, moduleID *uint) error {
	if err := s.permRepo.UpdatePermission(p); err != nil {
		return err
	}
	if moduleID != nil {
		return s.permRepo.SetPermissionModule(p.ID, moduleID)
	}
	return nil
}

func (s *PermissionService) DeletePermission(id uint) error {
	return s.permRepo.DeletePermission(id)
}

func (s *PermissionService) GetModulePermissionIDs(moduleID uint) ([]uint, error) {
	return s.permRepo.GetModulePermissionIDs(moduleID)
}

func (s *PermissionService) SetModulePermissions(moduleID uint, permIDs []uint) error {
	if _, err := s.permRepo.GetModuleByID(moduleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模块不存在")
		}
		return err
	}
	return s.permRepo.SetModulePermissions(moduleID, permIDs)
}
