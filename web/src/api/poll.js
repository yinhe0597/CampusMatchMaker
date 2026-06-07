import request from './request'

export const pollApi = {
  // 创建投票
  createPoll(data) {
    return request.post('/polls', data)
  },

  // 获取投票列表
  listPolls(params = {}) {
    return request.get('/polls', { params })
  },

  // 获取投票详情
  getPollDetail(id) {
    return request.get(`/polls/${id}`)
  },

  // 编辑投票
  editPoll(id, data) {
    return request.put(`/polls/${id}`, data)
  },

  // 开启投票
  openPoll(id) {
    return request.post(`/polls/${id}/open`)
  },

  // 关闭投票
  closePoll(id) {
    return request.post(`/polls/${id}/close`)
  },

  // 获取选项
  getOptions(id) {
    return request.get(`/polls/${id}/options`)
  },

  // 提交投票
  submitVote(id, data) {
    return request.post(`/polls/${id}/vote`, data)
  },

  // 获取结果
  getResults(id) {
    return request.get(`/polls/${id}/results`)
  },

  // 确认最终时段
  finalizePoll(id, data) {
    return request.post(`/polls/${id}/finalize`, data)
  },
}
