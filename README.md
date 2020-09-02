# HuobiProAPI - Go Language

## 依赖管理 
支持Dep和Glide管理依赖

## 实例代码 
### 配置文件  
修改 [配置文件](./config/config.go)中*todo*的部分，替换成自己的Key和自己要连接的服务地址
### 测试用例                                                                                                                                                  
可以参考[接口测试](./services/Market_test.go)
pkill -9 REST-GO-demo;rm -rf nohup.out;nohup ./REST-GO-demo &

### 编译web
cd ./element-html & npm install & npm run build & go-bindata-assetfs -pkg bindata ./dist/... -o ../app/bindata/bindata.go
