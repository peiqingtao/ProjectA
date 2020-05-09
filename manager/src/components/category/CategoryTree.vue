<template>
  <div>
    <!-- <el-table :data="tableData" style="width: 100%" :row-class-name="tableRowClassName">
            <el-table-column prop="ID" label="ID" width="180"></el-table-column>
            <el-table-column prop="Name" label="分类名" width="180"></el-table-column>
            <el-table-column prop="ParentId" label="上级分类ID"></el-table-column>
    </el-table>-->
    <el-row>
      <el-col :span="12">
        <el-page-header content="分类列表"></el-page-header>
      </el-col>

      <el-col :span="12">
        <el-button type="primary" class="float-right" @click="addDialogVisible = true">添加</el-button>
      </el-col>
      <el-dialog title="分类添加" :visible.sync="addDialogVisible">
        <el-form :model="categoryAdd" ref="categoryAddFrom" :rules="categoryAddRules" label-width="130px">
          <el-tabs type="border-card" v-model="addDialogActiveName">
            <el-tab-pane label="基础信息" name="general">
              <el-form-item label="分类名称" prop="Name">
                <el-input v-model="categoryAdd.Name" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="上级分类">
                <el-select
                  v-model="categoryAdd.ParentId"
                  placeholder="请选择上级分类"
                  @change="categoryOptionChangeHandle"
                >
                  <el-option
                    v-for="item in categoryOptions"
                    :key="item.ID"
                    :label="item.Name"
                    :value="item.ID"
                    :disabled="item.ID == 1"
                  >
                    <span v-bind:style="{paddingLeft: item.Deep*13 + 'px'}">{{ item.Name }}</span>
                  </el-option>
                </el-select>
                <el-checkbox
                  v-model="categoryAdd.IsTop"
                  style="margin-left:24px;"
                  @change="IsTopHandle"
                >是否顶级分类</el-checkbox>
              </el-form-item>
              <el-form-item label="分类 Logo">
                <el-upload action>
                  <el-button size="small" type="primary">点击上传</el-button>
                  <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
                </el-upload>
              </el-form-item>
              <el-form-item label="描述">
                <el-input
                  type="textarea"
                  :rows="3"
                  placeholder="请输入内容"
                  v-model="categoryAdd.Description"
                ></el-input>
              </el-form-item>
              <el-form-item label="排序" prop="SortOrder">
                <el-input v-model.number="categoryAdd.SortOrder" autocomplete="off"></el-input>
              </el-form-item>
            </el-tab-pane>
            <el-tab-pane label="seo 信息" name="seo">
              <el-form-item label="Meta Title">
                <el-input v-model="categoryAdd.MetaTitle" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="Meta Keywords">
                <el-input v-model="categoryAdd.MetaKeywords" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="Meta Description">
                <el-input type="textarea" v-model="categoryAdd.MetaDescription" autocomplete="off"></el-input>
              </el-form-item>
            </el-tab-pane>
          </el-tabs>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="addDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submitCategoryAddForm">确 定</el-button>
        </div>
      </el-dialog>
    </el-row>
    <el-row>
      <el-col :span="24">
        <el-tree :data="treeData" show-checkbox node-key="ID" :props="treeDataProps"></el-tree>
      </el-col>
    </el-row>
  </div>
</template>
<script>
// import { write } from 'fs';
import base from "../../api/uri.js";
export default {
  name: "CategoryTree",
  data() {
    return {
      //需要的树的数据
      treeData: [],
      //展示的字段名字
      treeDataProps: {
        children: "children",
        label: "Name"
      },
      //添加dialog是否可见
      addDialogVisible: false,
      categoryAdd: {
        IsTop: true,
        SortOrder: 0
      },
      addDialogActiveName: "general",
      categoryOptions: [],
      categoryAddRules: {
        Name: [
          { required: true, message: "请输入分类名称", trigger: ["blur", "change"]}
        ],
        SortOrder: [
          { required: true, message: "请输入排序值", trigger: ["blur", "change"]},
          { type:"integer", message: "排序值必须为整数", trigger: ["blur", "change"]}
        ]
      }
    };
  },
  //挂载后时件
  mounted() {
    this.getTree();
  },

  // 方法列表
  methods: {
    getTree() {
      let url = base + "category-tree"
      this.axios.get(url).then(resp => {
        if (resp.data.error == "") {
          this.treeData = this.tree(resp.data.data, 0)
          this.indentTree(resp.data.data, 0)
        } else {
          this.treeData = []
          this.categoryOptions = []
        }
      });
    },
    //递归的获取树状分类
    tree(all, id) {
      let items = []
      for (let item of all) {
        if (item.ParentId == id) {
          //存在子分类
          item.children = this.tree(all, item.ID);
          items.push(item)
        }
      }
      return items;
    },

    indentTree(all, id, deep = 0) {
      for (let item of all) {
        if (item.ParentId == id) {
          item.Deep = deep;
          this.categoryOptions.push(item)
          //存在子分类
          this.indentTree(all, item.ID, deep + 1)
        }
      }
    },

    categoryOptionChangeHandle(value) {
      this.categoryAdd.IsTop = false
    },

    IsTopHandle(value) {
      if (value) {
        this.categoryAdd.ParentId = null
      }
    },

    submitCategoryAddForm() {
      //表单检验
      this.$refs["categoryAddFrom"].validate(valid => {
        if (valid) {
          //验证通过
          //提交表单数据到服务端
          let url = base +"category";
          console.log(this.categoryAdd)
          this.axios.post(url,this.categoryAdd).then(resp => {
            window.console.log("添加完成")
          })
          //隐藏dialog
          this.addDialogVisible = false;

        } else {

          return false;

        }
      });
    },

    tableRowClassName({ rowIndex }) {
      if (rowIndex === 1) {
        return "warning-row";
      } else if (rowIndex === 3) {
        return "success-row";
      }
      return "";
    }
  }
};
</script>
<style scoped>
</style>>