import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store/index'
import axios from '../axios'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  // {
  //   path: '/',
  //   component: Layout,
  //   redirect: '/dashboard',
  //   children: [{
  //     path: 'dashboard',
  //     name: '控制面板',
  //     component: () => import('@/views/dashboard/index'),
  //     meta: { title: '控制面板', icon: 'dashboard' }
  //   }],
  //   meta: {}
  // },

  {
    path: '/',
    component: Layout,
    redirect: '/envList',
    children: [{
      path: '/envList',
      name: '环境列表',
      component: () => import('@/views/envList/index'),
      meta: { title: '环境列表', icon: 'example' }
    }],
    meta: {}
  },

  // {
  //   path: 'envList',
  //   component: Layout,
  //   children: [{
  //     path: '/envList',
  //     name: '环境列表',
  //     component: () => import('@/views/envList/index'),
  //     meta: { title: '环境列表', icon: 'example' }
  //   }],
  //   meta: {}
  // },
  {
    path: 'curdEnv',
    component: Layout,
    children: [{
      path: '/curdEnv',
      name: '环境设置',
      component: () => import('@/views/curdEnv/index'),
      meta: { title: '环境设置', icon: 'table' }
    }],
    meta: {
      requiresAuth: true,
      role: ['admin']
    }
  },
  {
    path: 'runningList',
    component: Layout,
    children: [{
      path: '/runningList',
      name: '运行列表',
      component: () => import('@/views/runningList/index'),
      meta: { title: '运行列表', icon: 'table' }
    }],
    meta: {
      requiresAuth: true,
      role: ['admin']
    }
    // eslint-disable-next-line eqeqeq
    // hidden: store.state.userForm.role !== 'admin'
  },

  // {
  //   path: 'external-link',
  //   component: Layout,
  //   children: [
  //     {
  //       path: 'https://panjiachen.github.io/vue-element-admin-site/#/',
  //       meta: { title: 'External Link', icon: 'link' }
  //     }
  //   ],
  //   // eslint-disable-next-line eqeqeq
  //   // hidden: store.state.userForm.role != 'admin'
  // },

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  // scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
// export function resetRouter() {
//   const newRouter = createRouter()
//   router.matcher = newRouter.matcher // reset router
// }

function auth (to) {
  // eslint-disable-next-line no-prototype-builtins
  if (!to.matched[0].meta.hasOwnProperty('role') || to.matched[0].meta.role === 0) {
    return true
  }
  // 求交集
  const set1 = new Set(to.matched[0].meta.role)
  const ready = []
  ready.push(store.state.userForm.role)
  const set2 = new Set(ready)
  const result = new Set([...set1].filter(x => set2.has(x)))
  return result.size !== 0
}

router.beforeEach((to, from, next) => {
  console.log('switch router')
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.state.userForm.login) {
      axios.get('/api/k8s/home').then((data) => {
        store.dispatch('userForm/login', data.data.data).then(() => {
          if (auth(to)) {
            next()
          } else {
            next(false)
          }
        })
      }).catch(() => {
        next({ path: '/login' })
      })
    } else {
      if (auth(to)) {
        next()
      } else {
        next(false)
      }
    }
  } else {
    next()
  }
})

export default router
