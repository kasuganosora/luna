import Vue from 'vue'
import VueRouter from 'vue-router'

const routes = [
    {path: '/', component: () => import('@/views/index'), hidden: true},
    {path: '/posts', component: () => import('@/views/post/index'), hidden: true},
    {path: '/posts/:id', component: ()=> import('@/views/post/editor'), hidden: true, name: 'editPost'},
    {path: '/posts/new', component: ()=> import('@/views/post/editor'), hidden: true, name: 'createPost'},
    {path: '/files', component: ()=> import('@/views/file/index'), hidden: true, name: 'filesList'},
    {
        path: '/login',
        component: () => import('@/views/login'),
        hidden: true,
        meta: {
            layout: "empty"
        }
    },
]

Vue.use(VueRouter)
const router = new VueRouter({
    routes
})
export default router