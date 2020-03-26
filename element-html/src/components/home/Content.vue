<template>

  <div>
    <el-breadcrumb separator-class="el-icon-arrow-right">
    <el-breadcrumb-item >首页</el-breadcrumb-item>
    <el-breadcrumb-item>策略列表</el-breadcrumb-item>
    </el-breadcrumb>
    <br><br>
    <!-- <img src="../../assets/background.jpg"> -->

  <template>
    <el-button @click="goAdd" type="primary">新增策略</el-button>
    <el-table
      :data="tableData"
      style="width: 100%">
      <el-table-column prop="Name" label="名称" width="80">
      </el-table-column>
      <el-table-column prop="MinPrice" label="区间min" width="80">
      </el-table-column>
      <el-table-column prop="MaxPrice" label="区间max" width="80">
      </el-table-column>
      <el-table-column prop="" label="交易对" width="90">
          <template slot-scope="scope">
            {{scope.row.QuoteCurrency+"/"+scope.row.BaseCurrency}}
        </template>
      </el-table-column>
      
      <el-table-column prop="HeightPrice" label="间距" width="80">
      </el-table-column>
      <el-table-column prop="SellPrice" label="平仓差价" width="80">
      </el-table-column>
      <el-table-column prop="start_time" label="运行状态" width="150">
        <template slot-scope="scope">
            {{scope.row.Run==1?"runing":"stop"}}
        </template>
      </el-table-column>
      <el-table-column prop="start_time" label="策略模式" width="150">
        <template slot-scope="scope">
            {{scope.row.Model==1?"经典":scope.row.Model==2?"只买":"只卖"}}
        </template>
      </el-table-column>
      
      <el-table-column prop="h24" label="1天收益" width="150">
        <template slot-scope="scope">
            {{scope.row.h24.toString().substring(0,6)}}
            {{"("+(Math.round(scope.row.h24/scope.row.Capital*365*10000)/100)+"%)"}}
        </template>
      </el-table-column>
      <el-table-column prop="d7" label="7天收益" width="150">
        <template slot-scope="scope">
            {{scope.row.d7.toString().substring(0,6)}}
            {{"("+(Math.round(scope.row.d7/scope.row.Capital*365*10000/7)/100)+"%)"}}
        </template>
      </el-table-column>
      <el-table-column prop="all" label="总收益" width="80">
        <template slot-scope="scope">
            {{scope.row.all.toString().substring(0,6)}}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220">
        <template slot-scope="scope">
          <el-button size="mini" @click="goEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button size="mini" v-if="scope.row.Run==1" type="primary" @click="handleStop(scope.$index, scope.row)">停止</el-button>
          <el-button size="mini" v-if="scope.row.Run==0" type="primary" @click="handleStart(scope.$index, scope.row)">启动</el-button>
          <el-button size="mini" v-if="scope.row.Model==2&&scope.row.Run==1" type="primary" @click="handleAction(scope.$index, scope.row)">卖出</el-button>
          <el-button size="mini" v-if="scope.row.Model==3&&scope.row.Run==1" type="primary" @click="handleAction(scope.$index, scope.row)">买入</el-button>
          <el-button size="mini" v-if="scope.row.Run==0" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
        </template>
      </el-table-column>
      
    </el-table>
  </template>

  </div>
</template>
<script>
  export default {
  data() {
        return {
          tableData: [
          ],
        }
      },
mounted(){
      var vue = this;
      get_data(vue);
  },
  methods: {
    goAdd(){
      this.$router.replace('/add');
    },
    goEdit(index,row){
      //this.$router.replace('/edit');
      this.$router.push({path:'/edit',query:{Id:row.Id}})

    },
    handleStart(index,row){
      var vue= this;
      handle_start(index,vue)
    },
    handleStop(index,row){
      var vue= this;
      handle_stop(index,vue)
    },
    handleDelete(index,row){
      var vue= this;
      handle_delete(index,vue)
    },
    handleAction(index,row){
      var vue= this;
      handle_action(index,vue)
    },
  }
  }
  
  
  function get_data(vue){
    vue.$http.get("/lianghua/list").then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        console.log(resp)
         for(var i=0;i< resp.data.data.length;i++){
           resp.data.data[i]["h24"]=0
           resp.data.data[i]["d7"]=0
           resp.data.data[i]["all"]=0
        }
        vue.tableData = resp.data.data;
        for(var i=0;i<vue.tableData.length;i++){
          get_profit(i,vue)
        }
    }).catch(error =>{
        vue.$message.error(error);
    })
  }
  function handle_start(index,vue){
    vue.$http.get("/lianghua/start?Id="+vue.tableData[index].Id).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        //vue.tableData[index].Run = 1;
        get_data(vue)
    }).catch(error =>{
        vue.$message.error(error);
    })

  }
  function handle_stop(index,vue){
    vue.$http.get("/lianghua/stop?Id="+vue.tableData[index].Id).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        get_data(vue)
    }).catch(error =>{
        vue.$message.error(error);
    })
  }
  function handle_action(index,vue){
    vue.$http.get("/lianghua/action?Id="+vue.tableData[index].id).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        get_data(vue)
    }).catch(error =>{
        vue.$message.error(error);
    })
  }
  function get_profit(index,vue){
    vue.$http.get("/lianghua/profit?Id="+vue.tableData[index].Id).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        console.log(resp.data.data)
        for (let it in resp.data.data) {
          vue.tableData[index][it] = resp.data.data[it]
        }
    }).catch(error =>{
        vue.$message.error(error);
    })

  }
  function handle_delete(index,vue){
    vue.$http.get("/lianghua/delete?Id="+vue.tableData[index].Id).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        get_data(vue)
    }).catch(error =>{
        vue.$message.error(error);
    })
  }
</script>
<style>
.el-popper{
  background:white !important;
}
.el-dialog__body{
  padding-top: 0px;
}
.el-table .acctive-row{
  background: oldlace;
}
.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
</style>