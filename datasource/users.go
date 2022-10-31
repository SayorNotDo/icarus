package datasource

import (
	"icaru/user"
	"time"
)

var Users = map[int64]user.User{
	1: {
		UID:            1,
		Username:       "wentao3.chen",
		HashedPassword: []byte("test1234"),
		ChineseName:    "陈文滔",
		RoleId:         5,
		EmployeeId:     "151926",
		Position:       "测试开发工程师",
		Email:          "wentao3.chen@test.com",
		Phone:          "15710790761",
		JoinDate:       time.Date(2022, time.Month(11), 11, 11, 11, 11, 11, time.UTC),
		LastLoginTime:  time.Date(2022, time.Month(11), 11, 11, 11, 11, 11, time.UTC),
		Status:         true,
		Department:     "测试部",
	},
}
