# review_collect
**A review collect server for chrome extension [fetch_review](https://github.com/ix64/fetch_review)**

### 使用方法
1. 下载 review_collect 对应版本并运行
2. 访问 review_collect 显示的地址，确认软件正常
（默认监听 [http://127.0.0.1:6481](http://127.0.0.1:6481)）

### 更改配置(.env)
- BIND review_collect的监听地址
- DB_PATH 上传数据的储存数据库（使用BoltDB，不存在将自动创建）
- LOG_PATH 上传数据的log输出目录


### HTTP API
- `GET /` 程序信息
- `POST /:channel/:id` 上传单个数据
- `GET /:channel` 获取该channel所有数据
> 所有channel:
> - taobao (淘宝)
> - tmall (天猫)
> - jd (京东)
> - suning (苏宁)

