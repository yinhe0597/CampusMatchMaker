<script setup>
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useClassStore } from '@/stores/class'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const classStore = useClassStore()

const isLoggedIn = computed(() => userStore.isLoggedIn)
const nickname = computed(() => userStore.userInfo?.nickname || '')

// 初始化：自动获取用户信息
onMounted(() => {
  userStore.init()
})

async function goToTimetable() {
  try {
    await classStore.fetchMyClasses()
  } catch {
    // fetchMyClasses 失败时降级跳转
  }
  if (classStore.myClasses.length > 0) {
    router.push({ name: 'Timetable', query: { class_id: classStore.myClasses[0].id } })
  } else {
    router.push({ name: 'ClassList' })
  }
}

function navigate(path) {
  router.push(path)
}

function handleLogout() {
  userStore.logout()
  router.push({ name: 'Login' })
}
</script>

<template>
  <el-container style="min-height: 100vh">
    <!-- 顶部导航 -->
    <el-header
      v-if="isLoggedIn"
      style="
        display: flex;
        align-items: center;
        justify-content: space-between;
        background: #fff;
        box-shadow: 0 1px 4px rgba(0,0,0,0.08);
      "
    >
      <div style="display: flex; align-items: center; gap: 24px;">
        <h2 style="margin: 0; color: #409eff;">校园协作平台</h2>
        <el-menu
          mode="horizontal"
          :default-active="route.path"
          :ellipsis="false"
        >
          <el-menu-item index="/polls" @click="navigate('/polls')">投票</el-menu-item>
          <el-menu-item index="/classes" @click="navigate('/classes')">班级</el-menu-item>
          <el-menu-item index="/timetable" @click="goToTimetable">课表</el-menu-item>
        </el-menu>
      </div>

      <div style="display: flex; align-items: center; gap: 12px;">
        <span style="font-size: 14px; color: #606266;">
          <el-icon><User /></el-icon>
          {{ nickname }}
        </span>
        <el-button type="danger" text @click="handleLogout">退出登录</el-button>
      </div>
    </el-header>

    <!-- 主内容区 -->
    <el-main style="padding: 20px; background: #f5f7fa;">
      <router-view />
    </el-main>
  </el-container>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>
