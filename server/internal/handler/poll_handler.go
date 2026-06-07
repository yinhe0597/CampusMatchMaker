package handler

import (
	"strconv"

	"campus_collab/internal/service"
	"campus_collab/internal/service/dto"
	"campus_collab/pkg/response"

	"github.com/gin-gonic/gin"
)

// PollHandler 投票处理器
type PollHandler struct {
	pollService *service.PollService
}

func NewPollHandler(svc *service.PollService) *PollHandler {
	return &PollHandler{pollService: svc}
}

// CreatePoll 创建投票
// POST /api/v1/polls
func (h *PollHandler) CreatePoll(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.pollService.CreatePoll(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// ListPolls 获取投票列表
// GET /api/v1/polls?scope_type=class&scope_id=1&page=1&page_size=20
func (h *PollHandler) ListPolls(c *gin.Context) {
	userID := c.GetUint("user_id")
	scopeType := c.Query("scope_type")
	scopeIDStr := c.Query("scope_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	scopeID, err := strconv.ParseUint(scopeIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "缺少有效的 scope_id 参数")
		return
	}
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	result, err := h.pollService.ListPolls(c.Request.Context(), userID, uint(scopeID), scopeType, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// GetPollDetail 获取投票详情
// GET /api/v1/polls/:id
func (h *PollHandler) GetPollDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	result, err := h.pollService.GetPollDetail(c.Request.Context(), userID, uint(pollID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// EditPoll 编辑投票
// PUT /api/v1/polls/:id
func (h *PollHandler) EditPoll(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	var req dto.EditPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.pollService.EditPoll(c.Request.Context(), userID, uint(pollID), &req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "编辑成功"})
}

// OpenPoll 开启投票
// POST /api/v1/polls/:id/open
func (h *PollHandler) OpenPoll(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	if err := h.pollService.OpenPoll(c.Request.Context(), userID, uint(pollID)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "投票已开启"})
}

// ClosePoll 关闭投票
// POST /api/v1/polls/:id/close
func (h *PollHandler) ClosePoll(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	if err := h.pollService.ClosePoll(c.Request.Context(), userID, uint(pollID)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "投票已关闭"})
}

// GetOptions 获取投票选项
// GET /api/v1/polls/:id/options
func (h *PollHandler) GetOptions(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	items, err := h.pollService.GetOptions(c.Request.Context(), userID, uint(pollID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"options": items})
}

// SubmitVote 提交投票
// POST /api/v1/polls/:id/vote
func (h *PollHandler) SubmitVote(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	var req dto.SubmitVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.pollService.SubmitVote(c.Request.Context(), userID, uint(pollID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// GetResults 获取投票结果
// GET /api/v1/polls/:id/results
func (h *PollHandler) GetResults(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	result, err := h.pollService.GetResults(c.Request.Context(), userID, uint(pollID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, result)
}

// FinalizePoll 确认最终时段
// POST /api/v1/polls/:id/finalize
func (h *PollHandler) FinalizePoll(c *gin.Context) {
	userID := c.GetUint("user_id")
	pollID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的投票ID")
		return
	}

	var req dto.FinalizePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.pollService.FinalizePoll(c.Request.Context(), userID, uint(pollID), &req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "已确认最终时段"})
}
