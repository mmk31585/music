import {
  createRouter,
  createWebHistory,
  type NavigationFailure,
  type NavigationGuardNext,
  type RouteLocationNormalized,
  type RouteLocationNormalizedLoaded,
  type RouteRecordRaw,
} from 'vue-router'
import PageNotFound from '@/pages/errors/PageNotFound.vue'
import middleware from './middleware/register'
import isObject from 'lodash.isobject'
import allRoutes from './routes'

const routes: RouteRecordRaw[] = [
  ...allRoutes,
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: PageNotFound,
    meta: {
      title: 'صفحه مورد نظر پیدا نشد!',
      layout: 'layout-empty',
    },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes as RouteRecordRaw[],
})

//------------------------------------------------------------------------------

router.beforeEach(
  async (
    to: RouteLocationNormalized,
    from: RouteLocationNormalizedLoaded,
    next: NavigationGuardNext,
  ) => {
    let result = null

    for (const check of middleware.beforeEach) {
      if (typeof check === 'function') {
        //@ts-expect-error: Ignore parameters count
        result = await check(to, from, result)
      }
    }

    if (result !== null && isObject(result)) return next(result)

    return next()
  },
)

router.beforeResolve(
  (to: RouteLocationNormalized, from: RouteLocationNormalizedLoaded, next: NavigationGuardNext) => {
    // for (const check of middleware.beforeResolve) {
    //   if (typeof check === 'function') {
    //     //@ts-expect-error: Ignore parameters count
    //     check(to, from)
    //   }
    // }

    next()
  },
)

router.afterEach(
  (
    to: RouteLocationNormalized,
    from: RouteLocationNormalizedLoaded,
    failure: NavigationFailure | void,
  ) => {
    for (const check of middleware.afterEach as ((
      to?: RouteLocationNormalized,
      from?: RouteLocationNormalizedLoaded,
      failure?: NavigationFailure | void,
    ) => unknown)[]) {
      if (typeof check === 'function') {
        check(to, from, failure)
      }
    }
  },
)

export default router
