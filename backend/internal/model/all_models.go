package model

func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&Role{},
		&UserRole{},
		&Permission{},
		&PermissionModule{},
		&PermissionModuleRelation{},
		&RoleModule{},
		&SheetMusic{},
		&SheetFile{},
		&AudioFile{},
		&SheetComment{},
		&SheetLike{},
		&Tag{},
		&SheetTag{},
		&SheetOrder{},
		&SystemConfig{},
		&SystemLog{},
		&BannedIp{},
		&UserLoginLog{},
		&ContactMessage{},
		&FriendLink{},
		&DownloadRecord{},
		&SitePvHistory{},
	}
}
