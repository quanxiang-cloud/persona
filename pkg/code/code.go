package code

import "git.internal.yunify.com/qxp/persona/pkg/misc/error2"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// InvalidURI 无效的URI
	InvalidURI = 160014000000
	// InvalidParams 无效的参数
	InvalidParams = 160014000001
	// InvalidTimestamp 无效的时间格式
	InvalidTimestamp = 160014000002
	// NameExist 名字已经存在
	NameExist = 160014000003
	// TimeOut 超时
	TimeOut = 160014000004
	// Rollback 回滚
	Rollback = 160014000005
	// LockExpire 锁过期
	LockExpire = 160014000006
)

// CodeTable 码表
var CodeTable = map[int64]string{
	InvalidURI:       "无效的URI.",
	InvalidParams:    "无效的参数.",
	InvalidTimestamp: "无效的时间格式.",
	NameExist:        "名称已被使用！请检查后重试！",
	TimeOut:          "超时",
	Rollback:         "回滚",
	LockExpire:       "锁已过期",
}
