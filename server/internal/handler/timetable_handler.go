package handler

import (
	"strconv"

	"campus_collab/internal/service"
	"campus_collab/internal/service/dto"
	"campus_collab/pkg/response"

	"github.com/gin-gonic/gin"
)

// TimetableHandler 课表处理器
type TimetableHandler struct {
	ttService *service.TimetableService
}

func NewTimetableHandler(svc *service.TimetableService) *TimetableHandler {
	return &TimetableHandler{ttService: svc}
}

// ===== 班级公共课表 =====

// CreateClassTimetable 录入班级公共课表
func (h *TimetableHandler) CreateClassTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("classId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	var req dto.CreateClassTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.ttService.CreateClassTimetable(c.Request.Context(), userID, uint(classID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// GetClassTimetable 获取班级公共课表
func (h *TimetableHandler) GetClassTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("classId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	result, err := h.ttService.GetClassTimetable(c.Request.Context(), userID, uint(classID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// UpdateClassTimetable 更新班级公共课表
func (h *TimetableHandler) UpdateClassTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("classId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	var req dto.UpdateClassTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.ttService.UpdateClassTimetable(c.Request.Context(), userID, uint(classID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// ===== 个人课表 =====

// CreatePersonalTimetable 添加个人课表条目
func (h *TimetableHandler) CreatePersonalTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreatePersonalTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.ttService.CreatePersonalTimetable(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// GetPersonalTimetable 获取个人完整课表
func (h *TimetableHandler) GetPersonalTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	classIDStr := c.Query("class_id")
	if classIDStr == "" {
		response.BadRequest(c, "class_id 不能为空")
		return
	}
	classID, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "class_id 格式不正确")
		return
	}

	result, err := h.ttService.GetPersonalTimetable(c.Request.Context(), userID, uint(classID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// UpdatePersonalTimetable 修改个人课表条目
func (h *TimetableHandler) UpdatePersonalTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	entryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "条目ID格式不正确")
		return
	}

	var req dto.UpdatePersonalTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err = h.ttService.UpdatePersonalTimetable(c.Request.Context(), userID, uint(entryID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKWithMsg(c, "已更新")
}

// DeletePersonalTimetable 删除个人课表条目
func (h *TimetableHandler) DeletePersonalTimetable(c *gin.Context) {
	userID := c.GetUint("user_id")

	entryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "条目ID格式不正确")
		return
	}

	err = h.ttService.DeletePersonalTimetable(c.Request.Context(), userID, uint(entryID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKWithMsg(c, "已删除")
}

// ===== 纠错 =====

// CreateCorrection 提交纠错
func (h *TimetableHandler) CreateCorrection(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateCorrectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.ttService.CreateCorrection(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// ListCorrections 获取纠错列表
func (h *TimetableHandler) ListCorrections(c *gin.Context) {
	userID := c.GetUint("user_id")

	classIDStr := c.Query("class_id")
	if classIDStr == "" {
		response.BadRequest(c, "class_id 不能为空")
		return
	}
	classID, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "class_id 格式不正确")
		return
	}

	status := -1
	if s := c.Query("status"); s != "" {
		status, _ = strconv.Atoi(s)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.ttService.ListCorrections(c.Request.Context(), userID, uint(classID), status, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Page(c, result.List, result.Total, result.Page, result.PageSize)
}

// ReviewCorrection 处理纠错
func (h *TimetableHandler) ReviewCorrection(c *gin.Context) {
	userID := c.GetUint("user_id")

	correctionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "纠错ID格式不正确")
		return
	}

	var req dto.ReviewCorrectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err = h.ttService.ReviewCorrection(c.Request.Context(), userID, uint(correctionID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	msg := "纠错已驳回"
	if req.Action == "approve" {
		msg = "纠错已采纳，课表已更新"
	}
	response.OKWithMsg(c, msg)
}
