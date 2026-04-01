# 项目结构：

# 下载依赖：go mod tidy

# 启动命令: go run main.go -db=sqlite -port=8888

# 配置文件config.json
{
    "port": "3000",//端口
    "db_type": "sqlite",//数据库类型
    "sqlite_path": "./db/homework04.db",//sqlite路径
    "mysql_dsn": "root:root@tcp(127.0.0.1:3306)/homework04?charset=utf8mb4&parseTime=True&loc=Local",//mysql数据库连接
    "jwt_secret": "2026jwt_secret"//JWT密钥
  }

# 接口
# 用户注册：
curl --location 'http://127.0.0.1:8888/open/api/user/register' \
--header 'Content-Type: application/json' \
--data-raw '{
  "username": "ckp",
  "password": "123456",
  "email": "ckp@example.com"
}'

# 登录：
curl --location 'http://127.0.0.1:8888/open/api/user/login' \
--header 'Content-Type: application/json' \
--data '{
  "username": "ckp",
  "password": "123456"
}'

# 创建文章
curl --location 'http://127.0.0.1:8888/auth/api/post' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM0MTU4NzgsInVzZXJfaWQiOjQsInVzZXJuYW1lIjoiY2twIn0.7PFbBum5f20Puk09b5LqiN6Nl0QpEkxVZqXUkVqtNHo' \
--header 'Content-Type: application/json' \
--data '{
    "title":"今日头条22",
    "content":"这是文章内容22"
}'

# 查询文章
curl --location 'http://127.0.0.1:8888/open/api/posts/1'


# 修改文章
curl --location --request PUT 'http://127.0.0.1:8888/auth/api/post/1' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM0MTU4NzgsInVzZXJfaWQiOjQsInVzZXJuYW1lIjoiY2twIn0.7PFbBum5f20Puk09b5LqiN6Nl0QpEkxVZqXUkVqtNHo' \
--header 'Content-Type: application/json' \
--data '{
    "title":"今日头条11",
    "content":"这是文章内容11"
}'

#  删除文章
curl --location --request DELETE 'http://127.0.0.1:8888/auth/api/post/1' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM0MTU4NzgsInVzZXJfaWQiOjQsInVzZXJuYW1lIjoiY2twIn0.7PFbBum5f20Puk09b5LqiN6Nl0QpEkxVZqXUkVqtNHo'


# 评论
curl --location 'http://127.0.0.1:8888/auth/api/posts/comments/2' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM0MTU4NzgsInVzZXJfaWQiOjQsInVzZXJuYW1lIjoiY2twIn0.7PFbBum5f20Puk09b5LqiN6Nl0QpEkxVZqXUkVqtNHo' \
--header 'Content-Type: application/json' \
--data '{
    "content":"呵呵"
}'


# 查询评论
curl --location 'http://127.0.0.1:8888/open/api/posts/comments/2'