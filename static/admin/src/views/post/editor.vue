<template>
    <div id="post-editor">
        <el-form ref="form" :model="post" label-width="80px" class="editor-main">
            <el-form-item label="标题">
                <el-input v-model="post.title"></el-input>
            </el-form-item>
            <!--  内容 -->
            <el-form-item>
                <mavon-editor v-model="post.markdown"/>
            </el-form-item>
            <!-- 分类 -->
            <!-- 标签 -->
            <!-- 链接设定 -->
            <!-- meta -->
            <!-- 提交工具栏 -->

        </el-form>
    </div>
</template>

<script>


    export default {
        data(){
          return {
              post:{
                  title: "标题是什么",
                  markdown:"# 这是标题"
              }
          }
        },
        name: "editor"
    }
</script>

<style scoped>

</style>