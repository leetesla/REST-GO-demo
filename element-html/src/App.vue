<template>
<div id="app">
      <el-container>

        <el-header > <Header/> </el-header>

        <el-container>

          <el-aside width="200px"> <Aside/> </el-aside>
          <el-main> <router-view/> </el-main>

        </el-container>
      </el-container>
</div>
</template>

<script>
import Header from './components/Header.vue'
import Aside from './components/Aside.vue'
export default {
  name:"app",
  components: {
    Header,Aside
  },
}
</script>

<style>
#app {
  font-family: Helvetica, sans-serif;
  text-align: center;
}
</style>
