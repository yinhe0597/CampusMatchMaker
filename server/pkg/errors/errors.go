package errors

import "fmt"

// AppError 自定义业务错误类型
type AppError struct {
	Code    int    // 业务错误码
	Message string // 用户可见的错误信息
	Err     error  // 原始错误（用于日志，不对外暴露）
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建业务错误
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap 包装已有错误
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// ===== 通用错误 =====
var (
	ErrInternal     = New(1000, "服务器内部错误")
	ErrBadRequest   = New(1001, "请求参数错误")
	ErrNotFound     = New(1002, "资源不存在")
	ErrRateLimited  = New(1003, "请求频率超限")
	ErrNotAvailable = New(1004, "功能暂未开放")
)

// ===== 认证错误 (1100-1199) =====
var (
	ErrUnauthorized    = New(1100, "未登录或登录已过期")
	ErrStudentIDFormat = New(1101, "学号格式不正确")
	ErrPasswordWrong   = New(1102, "密码错误")
	ErrStudentIDExists = New(1103, "学号已被注册")
	ErrAuthFailed      = New(1104, "学号认证失败")
	ErrTokenInvalid    = New(1105, "Token 无效")
	ErrForbidden       = New(1106, "无操作权限")
)

// ===== 班级错误 (1200-1299) =====
var (
	ErrClassNotFound  = New(1200, "班级不存在")
	ErrInviteCode     = New(1201, "邀请码错误")
	ErrAlreadyMember  = New(1202, "已是班级成员")
	ErrNoPermission   = New(1203, "无管理权限")
	ErrClassCodeExist = New(1204, "班级代码已存在")
)

// ===== 课表错误 (1300-1399) =====
var (
	ErrTimetableNotFound  = New(1300, "课表条目不存在")
	ErrTimetableExists    = New(1301, "班级课表已存在")
	ErrPeriodInvalid      = New(1302, "节次范围无效")
	ErrTimeConflict       = New(1303, "时间冲突")
	ErrCorrectionNotFound = New(1304, "纠错记录不存在")
	ErrCorrectionHandled  = New(1305, "纠错已处理")
)

// ===== 投票错误 (1400-1499) =====
var (
	ErrPollNotFound  = New(1400, "投票不存在")
	ErrPollClosed    = New(1401, "投票已关闭")
	ErrAlreadyVoted  = New(1402, "已投过票")
	ErrOptionNotBelong = New(1403, "选项不属于该投票")
	ErrOnlyCreator   = New(1404, "仅创建者可操作")
	ErrPollNotOpen     = New(1405, "投票尚未开启")
	ErrPollDeadline    = New(1406, "已达截止时间")
	ErrPollNoAccess    = New(1407, "无权限查看")
	ErrPollNotEditable = New(1408, "仅草稿状态可编辑")
	ErrPollCannotOpen  = New(1409, "仅草稿状态可开启投票")
)
