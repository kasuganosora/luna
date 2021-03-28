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
            <el-form-item label="分类">
                <el-cascader v-model="post.catalog_id" :options="catalogs"  :clearable="true" :show-all-levels="false"/>
            </el-form-item>
            <!-- 标签 -->
            <el-form-item label="标签">
                <el-tag
                        effect="plain"
                        :key="tag"
                        v-for="tag in post.tags"
                        closable
                        @close="handleTagClose(post,tag)"
                        :disable-transitions="false">
                    {{tag}}
                </el-tag>
                <el-input
                        v-if="newTagInputVisible"
                        class="input-new-tag"
                        v-model="newTag"
                        ref="saveTagInput"
                        size="small"
                        @keyup.enter.native="handleNewTagInputConfirm"
                        @blur="handleNewTagInputConfirm"
                >
                </el-input>
                <el-button v-else class="button-new-tag" size="small" @click="showNewTagInput">+ 新标签</el-button>
            </el-form-item>
            <!-- 链接设定 -->
            <!-- meta -->
            <!-- 提交工具栏 -->

        </el-form>
    </div>
</template>

<script>
    import Catalog from '../../service/Catalog'
    export default {
        data(){
          return {
              catalogs: Catalog.getCascaderFormat(),
              newTag: "",
              newTagInputVisible: false,
              post:{
                  title: "标题是什么",
                  markdown:"# 这是标题",
                  catalog_id: null,
                  tags:["标签1","标签2"]
              }
          }
        },
        name: "editor",
        methods:{
            handleTagClose(post, tag) {
                this.post.tags.splice(post.tags.indexOf(tag), 1);
            },
            showNewTagInput() {
                this.newTagInputVisible = true;
                this.$nextTick(_ => {// eslint-disable-line no-unused-vars
                    this.$refs.saveTagInput.$refs.input.focus();
                });
            },
            handleNewTagInputConfirm() {
                if (this.newTag !== "") {
                    this.post.tags.push(this.newTag);
                }
                this.newTagInputVisible = false;
                this.newTag = '';
            },
        }
    }
</script>

<style scoped>
    .el-tag + .el-tag {
        margin-left: 10px;
    }
    .button-new-tag {
        margin-left: 10px;
        height: 32px;
        line-height: 30px;
        padding-top: 0;
        padding-bottom: 0;
    }
    .input-new-tag {
        width: 90px;
        margin-left: 10px;
        vertical-align: bottom;
    }
</style>