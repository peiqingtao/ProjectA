import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)
import Layout from '../components/partial/Layout.vue'
import CategoryTree from '../components/category/CategoryTree.vue'
const routes = [
    //非常规的组件
    // {path:'/login',component:Login},
    //嵌套，常规的组件
    {
        path:'/', component:Layout,
        children:[
            { path:'category-tree',component:CategoryTree, },
        ]
    },
]

//定义路由
const router = new VueRouter({
    routes // (缩写) 相当于routes：routes
})

export default router