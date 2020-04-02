<template>

  <div>
    <el-breadcrumb separator-class="el-icon-arrow-right">
    <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
    <el-breadcrumb-item>配置</el-breadcrumb-item>

    </el-breadcrumb>
    <br><br>
    <!-- <img src="../../assets/background.jpg"> -->

<el-form ref="form" :model="form" label-width="100px" >
    <el-form-item label="名称">
      <el-col :span="4">
    <el-input v-model="Name" autocomplete="off"></el-input>
    </el-col>
    
  </el-form-item>
  <br>
  <el-form-item label="AccessKey">
    <el-col :span="4">
    <el-input v-model="apiKey.AccessKey"></el-input>
    </el-col>
      <el-col :span="2">
    SecretKey
    </el-col>
    <el-col :span="4">
      <el-input v-model="apiKey.SecretKey"></el-input>
    </el-col>
  </el-form-item>

   <el-form-item label="AccountId">
    <el-col :span="3">
    <el-select v-model="form.AccountId" placeholder="请选择">
    <el-option 
      v-for="item in accounts"
      :key="item.id"
      :label="item.id"
      :value="item.id">
    </el-option>
  </el-select>
        </el-col>
  </el-form-item>

  <el-form-item label="交易对">
    <el-col :span="7">
    <el-select v-model="form.QuoteCurrency" placeholder="请选择">
    <el-option
      v-for="item in Quotes"
      :key="item"
      :label="item"
      :value="item">
    </el-option>
  </el-select>

    <el-select v-model="form.BaseCurrency" placeholder="请选择" filterable>
    <el-option
      v-for="item in Bases"
      :key="item"
      :label="item"
      :value="item">
    </el-option>
  </el-select>
  {{balance}}
  </el-col>
  
  </el-form-item>
  <br>

   <el-form-item label="价格区间">
     <el-col :span="4">
    <el-input v-model="stracy.MinPrice" ></el-input>
    </el-col>
     <el-col :span="1">
    ~
    </el-col>
    
    <el-col :span="4">
    <el-input v-model="stracy.MaxPrice"></el-input>
    </el-col>
    <el-col :span="2">
    {{price}}
    </el-col>
  </el-form-item>

  <el-form-item label="间距">
  <el-col :span="4">
    <el-input v-model="stracy.HeightPrice"></el-input>
      </el-col>
        <el-col :span="1">
  {{form.QuoteCurrency}}
</el-col>
  </el-form-item>
   <el-form-item label="下单量">
     <el-col :span="4">
    <el-input v-model="stracy.Amount"></el-input>
  </el-col>
  <el-col :span="1">
  {{form.QuoteCurrency}}
</el-col>
  </el-form-item>

   <el-form-item label="平仓差价">
     <el-col :span="4">
    <el-input v-model="stracy.SellPrice"></el-input>
      </el-col>
        <el-col :span="1">
  {{form.QuoteCurrency}}
</el-col>
  </el-form-item>
    <el-form-item label="本金">
     <el-col :span="4">
    <el-input v-model="stracy.Capital"></el-input>
      </el-col>
        <el-col :span="1">
    {{form.QuoteCurrency}}
    </el-col>
  </el-form-item>
   <el-form-item label="模式">
    <el-col :span="4">
      <el-select v-model="Model" placeholder="请选择">
      <el-option  label="经典" value="1"></el-option>
      <el-option  label="只买" value="2"></el-option>
        <el-option  label="只卖" value="3"></el-option>
      </el-select>
    </el-col>
  </el-form-item>



<br>
  <el-form-item>
    <el-col :span="3">
    <el-button type="primary" @click="onSubmit">立即提交</el-button>
    </el-col>
  </el-form-item>


</el-form>

  </div>
</template>
<script>
  export default {
  data() {
        return {
          Name:"forever",
          Model:"1",
          
          apiKey:{
            AccessKey: 'xxx',
            SecretKey: 'xxx',
          },
          form: {
            AccountId: '8101421',
            QuoteCurrency:"usdt",
            BaseCurrency:"btc",
          },
          accounts:[],
          stracy:{
            MinPrice:"6000",
            MaxPrice:"12000",
            HeightPrice:"2.1",
            Amount:"1000",
            SellPrice:"3",//
            Capital:"6000", //成本
          },
          options:[],//BuyMount BuyPrice SellPrice
          Quotes:["usdt","btc","eth","ht","husd","alts"],//交易对计价币
          Bases:[],//
          balance:"",
          symbols:{},
          price:"",
        }
      },
  methods: {
    onSubmit() {
      var vue = this;
      on_submit(vue)
    },
    getApiKey(){
      var vue = this;
      get_accounts(vue);
    },

    getBalance(){
      var vue = this;
      get_balance(vue);
    }
  },
  mounted(){
    var vue = this;
    get_symbols(vue)
  },
  watch: {
    apiKey: {
      handler: 'getApiKey',
      deep: true
    },
    form:{
       handler:'getBalance',
       deep: true
    },
    "form.QuoteCurrency":function(oldData,newData){
     // alert(vue.form.QuoteCurrency)
      var vue = this;
      get_bases(vue);
    },
    "form.BaseCurrency":function(oldData,newData){
      //alert(vue.form.BaseCurrency)
      var vue = this;
      //alert(vue.form.BaseCurrency)
      get_price(vue);
    },
  }

}

function empty(data){
  if(data==""|| data ==0 || data==undefined){
    return true;
  }
  return false;
}
function get_balance(vue){
  if(vue.form.AccountId==""){
    return;
  }
  vue.balance = "";
  if(vue.form.QuoteCurrency!="" || empty(vue.form.AccountId)){
    vue.$http.get("/lianghua/get_balance?AccessKey="+vue.apiKey.AccessKey+"&SecretKey="+vue.apiKey.SecretKey+"&AccountId="+vue.form.AccountId+"&Currency="+vue.form.QuoteCurrency).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        var balanceJson ={};
        if(vue.balance != ""){
          balanceJson = JSON.parse(vue.balance);
        }
        balanceJson[vue.form.QuoteCurrency]=resp.data.data;
        vue.balance = JSON.stringify(balanceJson);
    }).catch(error =>{
        vue.$message.error(error);
    })
  }

if(vue.form.BaseCurrency!=""){
    vue.$http.get("/lianghua/get_balance?AccessKey="+vue.apiKey.AccessKey+"&SecretKey="+vue.apiKey.SecretKey+"&AccountId="+vue.form.AccountId+"&Currency="+vue.form.BaseCurrency).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        var balanceJson ={};
        if(vue.balance != ""){
          balanceJson = JSON.parse(vue.balance);
        }
        balanceJson[vue.form.BaseCurrency]=resp.data.data;
        vue.balance = JSON.stringify(balanceJson);
    }).catch(error =>{
        vue.$message.error(error);
    })
  }
  
}
function get_accounts(vue){
  vue.accounts=[];
  vue.form.AccountId="";
  if(vue.apiKey.AccessKey == "" || vue.apiKey.SecretKey== ""){
    return;
  }
  vue.$http.get("/lianghua/get_accounts?AccessKey="+vue.apiKey.AccessKey+"&SecretKey="+vue.apiKey.SecretKey).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        vue.accounts = resp.data.data;
        vue.form.AccountId = resp.data.data[0].id;
    }).catch(error =>{
        vue.$message.error(error);
    })
  }


  function get_symbols(vue){
      vue.$http.get("/lianghua/get_symbols").then(resp => {
         var code = resp.data.status;     
         if(code != "ok"){
          return;
        }
        vue.symbols = resp.data.data;
    }).catch(error =>{
        vue.$message.error(error);
    })
  }
  function get_bases(vue){
    vue.form.BaseCurrency = "";
    if(empty(vue.form.QuoteCurrency)){
      return
    }
    vue.Bases=[];
    for(let symbol  in vue.symbols){
        if(vue.symbols[symbol]["quote-currency"]== vue.form.QuoteCurrency ){
          vue.Bases.push(vue.symbols[symbol]["base-currency"])
        }
    }
  }
  function on_submit(vue){
    var ret = "Name="+vue.Name+"&Model="+vue.Model+"&";
    for (let it in vue.form) {
      ret += it + '=' + vue.form[it] + '&';
    }
    for (let it in vue.apiKey) {
      ret += it + '=' + vue.apiKey[it] + '&';
    }
      for (let it in vue.stracy) {
      ret += it + '=' + vue.stracy[it] + '&';
    }
    vue.$http.get("/lianghua/add?"+ret).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        vue.$message({message: '操作成功!',type: 'success'});
        vue.$router.replace('/home');
    }).catch(error =>{
        vue.$message.error(error);
    })
  
  }

function get_price(vue){
  vue.price = "";
  if (empty(vue.form.BaseCurrency) || empty(vue.form.AccountId)){
    return
  }
  vue.$http.get("/lianghua/get_price?AccessKey="+vue.apiKey.AccessKey+"&SecretKey="+vue.apiKey.SecretKey+"&AccountId="+vue.form.AccountId+"&QuoteCurrency="+vue.form.QuoteCurrency+"&BaseCurrency="+vue.form.BaseCurrency).then(resp => {
        var code = resp.data.code;     
        //未登录            
        if(code != 0){
          vue.$message.error(resp.data.message);
          return;
        }
        vue.price = resp.data.data+" "+vue.form.QuoteCurrency;
    }).catch(error =>{
        vue.$message.error(error);
    })
}
</script>
<style>
/* .el-popper{
  background:white !important;
} */

</style>