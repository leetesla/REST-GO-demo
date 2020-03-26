import Vue from "vue"
import Router from "vue-router"

import home from './home/Content'
import add from './add/Content'
import edit from './edit/Content'

Vue.use(Router)

export default new Router({
    routes:[
        {
            path:'/home',
            name:'home',
            component:home,
        },{
            path:'/',
            name:'home',
            component:home,
        },{
            path:'/add',
            name:'add',
            component:add,
        },{
            path:'/edit',
            name:'edit',
            component:edit,
        },
        
    ]
})