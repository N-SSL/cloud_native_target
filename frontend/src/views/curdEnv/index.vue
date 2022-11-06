<template>
  <div class="app-container">
    <el-dialog :visible.sync="centerDialogVisible"
               title="添加环境">
      <el-form ref="form"
               :model="form"
               label-width="80px">
        <el-form-item label="环境名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="环境类型">
          <el-input v-model="form.type" />
        </el-form-item>
        <el-form-item label="环境描述">
          <el-input v-model="form.description"
                    type="textarea" />
        </el-form-item>
        <el-form-item>
          <el-row>
            <el-col :span="12">
              <el-upload ref="upload"
                         :file-list="DeployFile"
                         :auto-upload="false"
                         action="null"
                         :on-change="deployFileUpload">
                <el-button slot="trigger"
                           size="small"
                           type="success">Deployment 文件</el-button>
              </el-upload>
            </el-col>
            <el-col :span="12">
              <el-upload ref="upload"
                         :file-list="ServiceFile"
                         :auto-upload="false"
                         action="null"
                         :on-change="serviceFileUpload">
                <el-button slot="trigger"
                           size="small"
                           type="success">Service 文件</el-button>
              </el-upload>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item>
          <el-button type="primary"
                     @click="onSubmit">立即创建</el-button>
          <el-button @click="cancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
    <el-table v-loading="listLoading"
              :data="list"
              element-loading-text="Loading"
              border
              fit
              highlight-current-row>
      <el-table-column align="center"
                       label="ID"
                       width="100px">
        <template slot-scope="scope">
          {{ scope.row.ID }}
        </template>
      </el-table-column>
      <el-table-column label="EnvName"
                       align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.EnvName }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Description">
        <template slot-scope="scope">
          {{ scope.row.Description }}
        </template>
      </el-table-column>
      <el-table-column label="Type"
                       align="center">
        <template slot-scope="scope">
          {{ scope.row.Type }}
        </template>
      </el-table-column>
      <el-table-column label="option"
                       align="center">
        <template slot-scope="scope">
          <el-button size="mini"
                     type="danger"
                     @click.stop="deleteRunningEnv(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div slot="header"
         style="position: fixed;right: 50px;bottom: 50px"
         class="clearfix">
      <div style="position: static;right: 10px;top: -5px">
        <el-button icon="el-icon-edit"
                   circle
                   type="primary"
                   @click="addEnv">添加</el-button>
      </div>
    </div>
  </div>
</template>

<script>
var YAML = require('js-yaml')
var json2yaml = require('json2yaml')
export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      list: null,
      listLoading: true,
      centerDialogVisible: false,
      form: {
        name: '',
        type: '',
        description: ''
      },
      DeployFile: [],
      ServiceFile: [],
      singleDeployFile: '',
      singleServiceFile: ''
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      this.$axios.get('/api/k8s/listAll').then(data => {
        console.log(data)
        this.list = data.data.data
        this.listLoading = false
      })
    },
    deleteRunningEnv(item) {
      const bodyFormData = new FormData()
      bodyFormData.append('Name', item.EnvName)
      this.$axios.post('/api/k8s/deleteEnv', bodyFormData).then(() => {
        this.$message({
          type: 'success',
          message: `删除成功`
        })
        this.fetchData()
        this.$forceUpdate()
      }).catch(err => {
        this.$message({
          type: 'error',
          message: `删除失败`
        })
        console.log(err)
      })
    },
    addEnv() {
      this.centerDialogVisible = true
    },
    onSubmit() {
      // let reader = new FileReader()
      // reader.readAsText(this.singleDeployFile.raw, 'utf-8')
      // let hh = reader.result
      // let a = YAML.load(this.singleDeployFile.raw)
      // let bstr = json2yaml.stringify(this.singleDeployFile.raw)
      // // console.log(hh)
      // console.log(hh)
      // console.log(a)
      // console.log(bstr)
      const bodyFormData = new FormData()
      bodyFormData.append('Name', this.form.name)
      bodyFormData.append('Type', this.form.type)
      bodyFormData.append('Description', this.form.description)

      bodyFormData.append('DeployFile', this.singleDeployFile.raw)
      bodyFormData.append('ServiceFile', this.singleServiceFile.raw)

      this.$axios.post('/api/k8s/addNewEnv', bodyFormData).then(data => {
        console.log(data)
        if (data.data.code === 200) {
          this.$message({
            type: 'success',
            message: `创建成功`
          })
        } else {
          console.log(data.data.data)
          this.$message({
            type: 'error',
            message: `创建失败`
          })
        }
      }).catch(err => {
        console.log(err)
      })
    },
    cancel() {
      this.form.description = ''
      this.form.name = ''
      this.form.type = ''
      this.singleServiceFile = ''
      this.singleDeployFile = ''
      this.DeployFile = []
      this.ServiceFile = []
    },
    deployFileUpload(file) {
      this.singleDeployFile = file
      this.DeployName = file.name
    },
    serviceFileUpload(file) {
      this.singleServiceFile = file
      this.ServiceName = file.name
    }
  }
}
</script>
