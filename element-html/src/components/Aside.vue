<template>
<div id="leaft">
    <el-menu @select="handleSelect">

      <el-menu-item index="2">
        <i class="el-icon-s-home"></i>
        <span slot="title">首页</span>
      </el-menu-item>
    </el-menu>
</div>
</template>
<script>
  export default {
    data:()=>({
            isCollapse:true,
        }),
    methods: {
      handleOpen(key, keyPath) {
        console.log(key, keyPath);
      },
      handleClose(key, keyPath) {
        console.log(key, keyPath);
      },
      handleSelect(key, keyPath) {
        console.log(key, keyPath);
        switch (key){
          case "2": //首页
          this.$router.push({path:'/home'});break;
          case "1-1"://kal分布
          this.$router.push({path:'kal_spread'});break;
          case "1-2": //邀请积分

          this.$router.push({path:'statis_inviter'});break;
          case "1-3": //赠送积分
          this.$router.push({path:'statis_point'});break;
          case "1-4": //时段分红
          this.$router.push({path:'statis_divider'});break;
          case "1-5": //积分竞拍
          this.$router.push({path:'statis_buykey'});break;
          case "2-1": //积分竞拍
          this.$router.push({path:'bvex_balance'});break;
          case "2-2": //交易
          this.$router.push({path:'trades'});break;
          case "2-3": //持仓
          this.$router.push({path:'holders'});break;
          default:
            break;
        }

      }
    }
  }
</script>