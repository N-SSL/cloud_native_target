<template>
  <div class="app-container">
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column align="center" label="ID" width="100px">
        <template slot-scope="scope">
          {{ scope.row.ID }}
        </template>
      </el-table-column>
      <el-table-column label="BindID">
        <template slot-scope="scope">
          {{ scope.row.BindID }}
        </template>
      </el-table-column>
      <el-table-column label="EnvName" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.EnvName }}</span>
        </template>
      </el-table-column>
      <el-table-column label="SessionID" align="center" width="100px">
        <template slot-scope="scope">
          {{ scope.row.SessionID }}
        </template>
      </el-table-column>
      <el-table-column label="ClusterIP" align="center">
        <template slot-scope="scope">
          {{ scope.row.ClusterIP }}
        </template>
      </el-table-column>
      <el-table-column label="Port" align="center" width="100px">
        <template slot-scope="scope">
          {{ scope.row.Port }}
        </template>
      </el-table-column>
      <el-table-column label="Url" align="center">
        <template slot-scope="scope">
          {{ scope.row.Url }}
        </template>
      </el-table-column>
      <el-table-column label="option" align="center">
        <template slot-scope="scope">
          <el-button size="mini" type="danger" @click.stop="deleteRunningEnv(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
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
      listLoading: true
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      this.$axios.get('/api/k8s/listAllRunningEnv').then(data => {
        this.list = data.data.data
        this.listLoading = false
      })
    },
    deleteRunningEnv(item) {
      const bodyFormData = new FormData()
      bodyFormData.append('Name', item.RunningName)
      this.$axios.post('/api/k8s/endEnvForce', bodyFormData).then(() => {
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
    }
  }
}
</script>
