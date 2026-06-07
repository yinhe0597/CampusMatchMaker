package handler

import (
	"strconv"

	"campus_collab/internal/service"
	"campus_collab/internal/service/dto"
	"campus_collab/pkg/response"

	"github.com/gin-gonic/gin"
)

// ClassHandler 班级处理器
type ClassHandler struct {
	classService *service.ClassService
}

// NewClassHandler 创建班级处理器
func NewClassHandler(svc *service.ClassService) *ClassHandler {
	return &ClassHandler{classService: svc}
}

// ListMyClasses 获取我的班级列表
func (h *ClassHandler) ListMyClasses(c *gin.Context) {
	userID := c.GetUint("user_id")

	result, err := h.classService.ListMyClasses(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, result)
}

// LookupClassByCode 通过邀请码查找班级
func (h *ClassHandler) LookupClassByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "邀请码不能为空")
		return
	}

	result, err := h.classService.LookupClassByCode(c.Request.Context(), code)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, result)
}

// CreateClass 创建班级
func (h *ClassHandler) CreateClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.classService.CreateClass(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, result)
}

// GetClassDetail 获取班级详情
func (h *ClassHandler) GetClassDetail(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	result, err := h.classService.GetClassDetail(c.Request.Context(), userID, uint(classID))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, result)
}

// JoinClass 加入班级
func (h *ClassHandler) JoinClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	var req dto.JoinClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.classService.JoinClass(c.Request.Context(), userID, uint(classID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, result)
}

// ListMembers 成员列表
func (h *ClassHandler) ListMembers(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.classService.ListMembers(c.Request.Context(), userID, uint(classID), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Page(c, result.Members, result.Total, result.Page, result.PageSize)
}

// RemoveMember 移除成员
func (h *ClassHandler) RemoveMember(c *gin.Context) {
	operatorID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "班级ID格式不正确")
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "用户ID格式不正确")
		return
	}

	err = h.classService.RemoveMember(c.Request.Context(), operatorID, uint(classID), uint(targetUserID))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKWithMsg(c, "移除成功")
}
