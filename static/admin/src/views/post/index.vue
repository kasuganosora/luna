<template>
  <div id="posts-list-page">
    <el-header>
      <h1>文章 <el-button size="small" @click="$router.push('/posts/new')">写文章</el-button></h1>
    </el-header>
    <!-- header toolbox 1 start -->
    <el-row type="flex" :gutter="20">
      <el-col>
        <el-breadcrumb separator="|">
          <el-breadcrumb-item>全部(1024)</el-breadcrumb-item>
          <el-breadcrumb-item>我的(2048)</el-breadcrumb-item>
          <el-breadcrumb-item>已发布(233)</el-breadcrumb-item>
          <el-breadcrumb-item>草稿(64)</el-breadcrumb-item>
          <el-breadcrumb-item>回收站(8848)</el-breadcrumb-item>
        </el-breadcrumb>
      </el-col>
      <el-col :span="5">
        <el-form>
          <el-input placeholder="请输入内容">
            <el-button slot="append" icon="el-icon-search"></el-button>
          </el-input>
        </el-form>
      </el-col>
    </el-row>
    <!-- header toolbox 1 end -->
    <!-- post table start -->
    <el-table stripe ref="multipleTable" style="width: 100%" :data="tableData">
      <el-table-column type="selection" width="55" />
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link :to="{ name: 'editPost', params: { id: scope.row.id }}">{{scope.row.title}}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="作者">
        <template slot-scope="scope">
          <el-link href="#" target="_blank">{{scope.row.author.name}}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="分类">
        <template slot-scope="scope">
          <el-breadcrumb separator=",">
            <el-breadcrumb-item v-for="catalog in scope.row.catalogs" :key="catalog.id">{{catalog.name}}</el-breadcrumb-item>
          </el-breadcrumb>
        </template>
      </el-table-column>
      <el-table-column label="标签">
        <template slot-scope="scope">
          <el-tag v-for="tag in scope.row.tags" :key="tag.id">{{tag.name}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="日期">
        <template slot-scope="scope">
          <el-row>
            <el-tag v-if="scope.row.status === 'published'" type="success">已发布</el-tag>
            <el-tag v-if="scope.row.status === 'draft'">草稿</el-tag>
            <el-tag v-if="scope.row.status === 'deleted'">已删除</el-tag>
          </el-row>
          <el-row>
            <span v-if="scope.row.status === 'published'">{{scope.row.published_at | fromNow}}</span>
            <span v-if="scope.row.status !== 'published'">{{scope.row.created_at | fromNow}}</span>
          </el-row>
        </template>
      </el-table-column>
    </el-table>
    <!-- post table end -->
    <!-- page footer start -->
    <el-row type="flex" :gutter="20">
      <el-pagination
          background
          layout="prev, pager, next"
          :total="1000">
      </el-pagination>
    </el-row>
    <!-- page footer end -->
  </div>
</template>

<script>
  import postSrv from '../../service/Post'

export default {
  name: "index",
  computed: {
    tableData() {
      return postSrv.getPosts()
    },
  }
}
</script>

<style scoped>

</style>