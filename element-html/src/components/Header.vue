<template>
<div id="Head">
    <el-dropdown v-if="islogin" trigger="click" @command="logout">
    <span class="el-dropdown-link">
        {{account}}<i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
        <el-dropdown-item >退出</el-dropdown-item>
    </el-dropdown-menu>
    </el-dropdown>
    <el-link v-if="!islogin" :underline="false" @click="dialogFormVisible = true">请登录</el-link>

<!-- 登录框 -->
<el-dialog title="" :visible.sync="dialogFormVisible" width="300px" :showClose="false" :closeOnClickModal="false">
  <el-form :model="form">
    <el-form-item label="用户名" >
      <el-input v-model="form.account" autocomplete="off"></el-input>
    </el-form-item>
    <el-form-item label="密码" >
      <el-input v-model="form.pwd" autocomplete="off" show-password></el-input>
    </el-form-item>
  </el-form>
  <div slot="footer" class="dialog-footer">
    <el-button type="primary" @click="login">确 定</el-button>
  </div>
</el-dialog>

</div>




</template>
<script>
    export default {
        name: 'Header',
        data:()=>({
            dialogFormVisible:true,
            islogin:false,
            account:"",
            form: {
              account: '',
              pwd: '',
            },
        }),
        mounted(){
            var vue = this;
            checkloin(vue);
        },
        methods: {
        login(){
            var vue = this
            login(vue);
            },
        logout(cmd) {
            console.log(cmd)
            var vue = this
            logout(vue)
            }
        }
    }
    function checkloin(vue){
        vue.$http.get("/account/get").then(resp => {
            var code = resp.data.code;     
            //未登录            
            if(code != 0){
                return;
            }
           if( resp.data.account!=""){
               vue.account = resp.data.account;
               vue.islogin = true;
               vue.dialogFormVisible= false;
           }
            
        }).catch(error =>{
            vue.$message.error(error);
        })
    }
    function logout(vue){
        vue.$http.get("/account/logout").then(resp => {
            
            var code = resp.data.code;
            var message = resp.data.message;
            if(code != 0){
                vue.$message.error(message);
                return;
            }
            vue.islogin=false;
            vue.dialogFormVisible=true;
            vue.$message({message: '退出成功!',type: 'success'});
        }).catch(error =>{
            vue.$message.error(error);
        })
    }
    function login(vue){
        vue.$http.get("/account/login?account="+vue.form.account+"&pwd="+vue.form.pwd).then(resp => {
            console.log(resp.data.code)
            var code = resp.data.code;
            var message = resp.data.message;
            console.log("login",message)    
            if(code != 0){
                vue.$message.error(message);
                return;
            }
            vue.islogin=true;
            vue.dialogFormVisible=false;
            vue.account=vue.form.account;
            vue.$message({message: '登录成功!',type: 'success'});
        }).catch(error =>{
            vue.$message.error(error);
        })
    }
</script>
<style>
#Head {
    height: 100%;
    background-color:transparent;
    text-align: right;
    float: right;
    align-items: center;
    display: flex;
}
.el-dialog__wrapper{
    background: white;
}
</style>