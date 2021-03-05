<template>
  <el-container id="index-page" class="admin-page">
    <!-- Sidebar start -->
    <el-aside class="sidebar" :style="sidebarStyle">
      <el-menu default-active="2" text-color="#fff" active-text-color="#ffd04b" background-color="#545c64">
        <!-- Sidebar content start -->
        <el-menu-item><i class="el-icon-odometer"></i>仪表盘</el-menu-item>
        <!-- Sub Sidebar start -->
        <el-submenu index="1">
          <template slot="title">
            <i class="el-icon-edit-outline"></i>
            <span>文章</span>
          </template>
          <!-- Sub Sidebar content start -->
          <el-menu-item><i class="el-icon-document"></i>所有文章</el-menu-item>
          <el-menu-item><i class="el-icon-collection"></i>文章分类</el-menu-item>
          <el-menu-item><i class="el-icon-collection-tag"></i>标签</el-menu-item>
          <!-- Sub Sidebar content end -->
        </el-submenu>
        <!-- Sub Sidebar end -->
        <!-- Sub Sidebar start -->
        <el-submenu index="2">
          <template slot="title">
            <i class="el-icon-paperclip"></i>
            <span>文件</span>
          </template>
          <!-- Sub Sidebar content start -->
          <el-menu-item><i class="el-icon-files"></i>所有文件</el-menu-item>
          <el-menu-item><i class="el-icon-upload2"></i>添加文件</el-menu-item>
          <!-- Sub Sidebar content end -->
        </el-submenu>
        <!-- Sub Sidebar end -->

        <el-menu-item><i class="el-icon-chat-square"></i>评论</el-menu-item>
        <!-- Sub Sidebar start -->
        <el-submenu index="3">
          <template slot="title">
            <i class="el-icon-user"></i>
            <span>用户</span>
          </template>
          <!-- Sub Sidebar content start -->
          <el-menu-item><i class="el-icon-user"></i>所有用户</el-menu-item>
          <el-menu-item><i class="el-icon-document-add"></i>添加用户</el-menu-item>
          <el-menu-item><i class="el-icon-setting"></i>个人资料</el-menu-item>
          <!-- Sub Sidebar content end -->
        </el-submenu>
        <!-- Sub Sidebar end -->
        <!-- Sub Sidebar start -->
        <el-submenu index="4">
          <template slot="title">
            <i class="el-icon-setting"></i>
            <span>设置</span>
          </template>
          <!-- Sub Sidebar content start -->
          <el-menu-item><i class="el-icon-setting"></i>常规</el-menu-item>
          <!-- Sub Sidebar content end -->
        </el-submenu>
        <!-- Sub Sidebar end -->

        <!-- Sidebar content end -->
      </el-menu>
    </el-aside>
    <!-- Sidebar END -->
    <el-main>
      <!-- Page content start -->
      <router-view />
      <!-- Page content end -->
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: "layout",
  computed: {
    sidebarStyle() {
      let style = "height: " + this.$store.state.screenHeight + 'px'
      return style
    }
  }
}
</script>

<style scoped type="scss">
.sidebar {
  min-height: 100%;
  background-color: #545c64;
  overflow-y: scroll;
}
</style>