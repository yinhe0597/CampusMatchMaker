import request from './request'

export const timetableApi = {
  // 班级公共课表
  createClassTimetable(classId, data) {
    return request.post(`/timetables/class/${classId}`, data)
  },
  getClassTimetable(classId) {
    return request.get(`/timetables/class/${classId}`)
  },
  updateClassTimetable(classId, data) {
    return request.put(`/timetables/class/${classId}`, data)
  },

  // 个人课表
  createPersonalTimetable(data) {
    return request.post('/timetables/personal', data)
  },
  getPersonalTimetable(classId) {
    return request.get('/timetables/personal', { params: { class_id: classId } })
  },
  updatePersonalTimetable(id, data) {
    return request.put(`/timetables/personal/${id}`, data)
  },
  deletePersonalTimetable(id) {
    return request.delete(`/timetables/personal/${id}`)
  },

  // 纠错
  createCorrection(data) {
    return request.post('/timetables/corrections', data)
  },
  listCorrections(params) {
    return request.get('/timetables/corrections', { params })
  },
  reviewCorrection(id, data) {
    return request.put(`/timetables/corrections/${id}`, data)
  },
}
