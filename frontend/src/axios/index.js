import axios from 'axios'

const config = {
  // 请求超时时间
  timeout: 10000,
  // 每次请求携带cookie
  withCredentials: true
}

const instance = axios.create(config)

export default instance
