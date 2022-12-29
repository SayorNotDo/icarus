package user

func BuildUserInfo(u *User) *UserInfo {
	if u == nil {
		return nil
	}
	ret := &UserInfo{
		UID:         u.UID,
		Username:    u.Username,
		ChineseName: u.ChineseName,
		EmployeeId:  u.EmployeeId,
		Email:       u.Email,
		Phone:       u.Phone,
		Department:  u.Department,
	}
	return ret
}
