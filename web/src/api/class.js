import request from './request'

export const classApi = {
  // 获取我的班级列表
  listMyClasses() {
    return request.get('/classes')
  },

  // 创建班级
  createClass(data) {
    return request.post('/classes', data)
  },

  // 通过邀请码查找班级
  lookupByCode(code) {
    return request.get(`/classes/by-code/${code}`)
  },

  // 获取班级详情
  getClassDetail(id) {
    return request.get(`/classes/${id}`)
  },

  // 加入班级（需要 classID + 邀请码）
  joinClass(id, inviteCode) {
    return request.post(`/classes/${id}/join`, { invite_code: inviteCode })
  },

  // 获取成员列表
  listMembers(id, params = {}) {
    return request.get(`/classes/${id}/members`, { params })
  },

  // 移除成员
  removeMember(classId, userId) {
    return request.delete(`/classes/${classId}/members/${userId}`)
  },
}
