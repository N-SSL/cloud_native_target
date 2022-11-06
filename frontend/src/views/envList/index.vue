<template>
  <div class="dashboard-container">
    <el-dialog :visible.sync="centerDialogVisible" title="镜像信息">
      <div v-if="running_env_state" v-loading="startCon" class="text item" element-loading-text="环境启动中">
        <div class="text item">
          访问地址: <a :href="running_env_state.url">{{ running_env_state? running_env_state.url:"" }}</a>
        </div>
        <div class="text item">
          环境名：
          <el-tag style="margin-right: 5px;">
            {{ running_env_state.EnvName }}
          </el-tag>
        </div>
        <div class="text item">
          部署名称: {{ running_env_state.deployName }}
        </div>
        <div class="text item">
          描述: {{ running_env_state.Description }}
        </div>
      </div>
    </el-dialog>
    <el-card class="box-card">
      <el-row :gutter="6">
        <el-input v-model="search" style="width: 230px;margin-left: 6px" size="medium" @keyup.enter.native="handleQuery(1)" />
        <el-button class="filter-item" size="medium" style="margin-left: 10px;margin-bottom: 10px" type="primary" icon="el-icon-search" @click="handleQuery(1)">
          查询
        </el-button>
      </el-row>
    </el-card>
    <el-divider style="margin-top: 1px" />
    <el-row id="first-bmh3" :gutter="24">
      <el-col v-for="(item,index) in listdata" :key="index" :span="8" style="padding-bottom: 18px;">
        <el-card
          :body-style="{ padding: '8px',marginLeft: '20px' }"
          shadow="hover"
          @click.native=" item.status.status === 'running'"
        >
          <div class="clearfix" style="position: relative">
            <div class="container-title">
              <span style="font-size: 30px">{{ item.EnvName }}</span>
            </div>
            <div style="margin-top: 7px;">
              <el-rate v-model="item.rank" disabled show-score text-color="#ff9900" score-template="{value}" />
            </div>
          </div>
          <div style="padding: 5px;">
            <div class="container-type">
              <span>类型：{{ item.Type }}</span>
            </div>
            <div class="bottom clearfix">
              <div class="time container-title">描述： {{ item.Description }}</div>
            </div>
            <el-row>
              <el-button v-if="item.status.status === 'running'" type="primary" :disabled="item.status.stop_flag" size="mini" @click.stop="startEnv(item)">进入</el-button>
              <el-button v-if="item.status.status === 'running'" type="primary" :disabled="item.status.stop_flag" size="mini" @click.stop="restart(item)">重启</el-button>
              <el-button v-else type="primary" :disabled="item.status.start_flag" size="mini" @click.stop="startEnv(item)">启动</el-button>
              <el-button v-if="item.status.status === 'running' || item.status.status === 'stop'" type="primary" :disabled="item.status.delete_flag" size="mini" @click.stop="endEnv(item)">删除</el-button>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
// import { Loading } from 'element-ui'
export default {
  name: 'Index',
  data() {
    return {
      page: {
        total: 0,
        size: 20
      },
      name: 114514,
      local: 'base',
      listdata: [],
      search: '',
      loading: true,
      drawerFlag: false,
      drawer: false,
      input: '',
      images_id: '',
      container_id: '',
      images_name: '',
      images_desc: '',
      writeup_date: '',
      startCon: '',
      vul_host: '',
      is_flag: true,
      expire: true,
      item_raw_data: '',
      centerDialogVisible: false,
      cStatus: true,
      running_env_state: null
    }
  },
  mounted() {
    // this.$axios.get('/api/k8s/ping/pong')
    //   .then(res => console.log(res))
    //   .catch(err => console.log(err))
  },
  created() {
    this.listData(1)
  },
  methods: {
    listData() {
      const that = this
      this.$axios.get('/api/k8s/listAll').then(response => {
        that.listdata = response.data.data
        for (let i = 0; i < this.listdata.length; i++) {
          that.listdata[i].checkId = i
          that.listdata[i].status = {}
          that.listdata[i].status.status = 'null'
          that.listdata[i].status.start_flag = false
          that.listdata[i].status.stop_flag = false
          that.listdata[i].rank = 2.5
        }
      })
    },
    endEnv(item) {
      console.log('end')
      const bodyFormData = new FormData()
      bodyFormData.append('Name', item.EnvName)
      this.$axios.post('/api/k8s/endEnv', bodyFormData).then(data => {
        if (data.data.code === 200) {
          this.listdata[item.checkId].status.status = 'null'
          this.listdata[item.checkId].status.start_flag = false
          this.running_env_state = null
          this.$forceUpdate()
          this.$message({
            type: 'success',
            message: `删除成功`
          })
        }
      }).catch(err => {
        this.listdata[item.checkId].status.status = 'null'
        this.listdata[item.checkId].status.start_flag = false
        this.running_env_state = null
        this.$forceUpdate()
        console.log(err)
      })
    },
    restart(item) {
      console.log('restart')
      const bodyFormData = new FormData()
      bodyFormData.append('Name', item.EnvName)
      this.$axios.post('/api/k8s/restartEnv', bodyFormData).then(data => {
        if (data.data.code === 200) {
          this.listdata[item.checkId].status.status = 'running'
          this.listdata[item.checkId].status.start_flag = true
          this.running_env_state = null
          this.running_env_state = data.data.data
          this.$forceUpdate()
          this.$message({
            type: 'success',
            message: `重启成功`
          })
          this.openDialog()
        }
      }).catch(err => {
        console.log(err)
      })
    },
    startEnv(item) {
      console.log('start')
      const bodyFormData = new FormData()
      bodyFormData.append('Name', item.EnvName)
      this.$axios.post('/api/k8s/startEnv', bodyFormData).then(data => {
        if (data.data.code === 200) {
          this.listdata[item.checkId].status.status = 'running'
          this.listdata[item.checkId].status.start_flag = true
          this.running_env_state = null
          this.running_env_state = data.data.data
          this.$forceUpdate()
          this.$message({
            type: 'success',
            message: `启动成功`
          })
          this.openDialog()
        }
      }).catch(err => {
        this.$message({
          type: 'error',
          message: `启动失败`
        })
        console.log(err)
      })
    },
    openDialog() {
      console.log('openDialog')
      this.centerDialogVisible = true
    }
  }

}
</script>

<style lang="scss" scoped>
  .dashboard {
    &-container {
      margin: 30px;
    }
    &-text {
      font-size: 30px;
      line-height: 46px;
    }
  }
  .time {
    font-size: 13px;
    color: #999;
  }
  .bottom {
    margin-top: 5px;
    margin-bottom: 13px;
    line-height: 12px;
  }
  .button {
    padding: 5px;
    float: right;
  }
  .image {
    width: 100%;
    display: block;
  }
  .clearfix:before,
  .clearfix:after {
    display: table;
    content: "";
  }
  .clearfix:after {
    clear: both
  }
  .text {
    font-size: 14px;
  }
  .item {
    margin-bottom: 18px;
  }
  .container-title{
    margin-top: 10px;
    width: 100%;    /*根据自己项目进行定义宽度*/
    overflow: hidden;     /*设置超出的部分进行影藏*/
    text-overflow: ellipsis;     /*设置超出部分使用省略号*/
    white-space:nowrap ;    /*设置为单行*/
  }
  .date {
  }
  .date p{
    height: 20px;
    line-height: 20px;
    margin: 0;
    margin-block-end: 0em;
  }
  .el-row {
    display: flex;
    flex-wrap: wrap;
  }
  .container-type{
    width: 100%;    /*根据自己项目进行定义宽度*/
    overflow: hidden;     /*设置超出的部分进行影藏*/
    text-overflow: ellipsis;     /*设置超出部分使用省略号*/
    white-space:nowrap ;    /*设置为单行*/
  }
  /*p {*/
  /*  height: 20px;*/
  /*  line-height: 20px;*/
  /*}*/
</style>

<style rel="stylesheet/scss" lang="scss">
  .el-drawer{
    overflow: scroll
  }
</style>
