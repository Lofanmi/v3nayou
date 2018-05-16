# v3nayou

微信哪有服务号, 做最好用的校园服务平台!

广州大学与中山大学的工具类应用, 支持广州大学教务成绩课表, 校历, 图书, 实时公交, 四六级, 公选课查询.

欢迎添砖加瓦~ 本项目为公众号自定义菜单的Go后端接口, 具体实现包含`首屏界面`, `微信授权`和`教务系统`.

对应的前端项目(Vue2): [https://github.com/Lofanmi/v3nayou-spa](https://github.com/Lofanmi/v3nayou-spa)

## Requirements

- Go 1.8+
- MySQL 5.5+
- Nginx

## Deploy

#### 前端

```bash
# 克隆前端项目到本地
cd /home/wwwroot/
git clone https://github.com/Lofanmi/v3nayou-spa.git
```

#### 后端

```bash
# GOPATH 路径可以自己定义
export GOPATH=/home/golang
# 导入SQL
cd $GOPATH/src/github.com/Lofanmi/v3nayou/db_schema/
mysql -h 127.0.0.1 -u root -p < v3nayou.sql
# 拉取
go get -v github.com/Lofanmi/v3nayou
# 编译
cd $GOPATH/src/github.com/Lofanmi/v3nayou
go build
# 配置环境变量
cp .env.example .env
vi .env
# 运行 (可配合 supervisor 等工具实现常驻)
./v3nayou
```

#### Nginx 配置文件

```
server
    {
        listen 80;
        server_name localhost;
        root  /path/to/v3nayou-spa/dist;

        location ~ .*\.(gif|jpg|jpeg|png|bmp|swf)$
        {
            expires      30d;
        }
        location ~ .*\.(js|css)?$
        {
            expires      30d;
        }

        location ^~ /school/ {
            proxy_pass http://127.0.0.1:5666;
        }
        location ^~ /api/ {
            proxy_pass http://127.0.0.1:5666;
        }

        location ~ /\.
        {
            deny all;
        }
        location = / {
            index index.html;
        }

        access_log /path/to/log/v3nayou-spa.log;
    }
```

## License

Apache-2.0
