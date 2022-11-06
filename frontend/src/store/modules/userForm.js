const state = {
  login: false, // 登录状态
  user: null, // 用户信息
  role: ''
}

const mutations = {
  LOGIN: (state, data) => {
    state.user = data.username
    state.role = data.role
    state.login = true
  },
  LOGOUT: (state) => {
    state.user = null
    state.role = ''
    state.login = false
  }
}

const actions = {
  login({ commit }, data) {
    commit('LOGIN', data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
