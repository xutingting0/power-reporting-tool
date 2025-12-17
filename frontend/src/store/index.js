import { createStore } from 'vuex'

export default createStore({
  state: {
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}'),
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
      localStorage.setItem('token', token)
    },
    SET_USER_INFO(state, userInfo) {
      state.userInfo = userInfo
      localStorage.setItem('userInfo', JSON.stringify(userInfo))
    },
    CLEAR_AUTH(state) {
      state.token = ''
      state.userInfo = {}
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  },
  actions: {
    // 修改login action，接收真实token和userInfo
    login({ commit }, { token, userInfo }) {
      commit('SET_TOKEN', token)
      commit('SET_USER_INFO', userInfo)
    },
    logout({ commit }) {
      commit('CLEAR_AUTH')
    }
  },
  getters: {
    isLoggedIn: state => !!state.token,
    userRole: state => state.userInfo.role || '',
    userName: state => state.userInfo.username || '',
    userId: state => state.userInfo.userId || '',
    // 添加其他需要的getters
    userRealName: state => state.userInfo.realName || state.userInfo.real_name || '',
    userUnit: state => state.userInfo.unit || '',
    userPhone: state => state.userInfo.phone || '',
    userPosition: state => state.userInfo.position || ''
  }
})