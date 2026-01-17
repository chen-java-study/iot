# Windows 手动启动指南（最简单）

## 第一步：启动数据库（只需做一次）

### 方式A：使用Docker（推荐）

打开PowerShell或CMD，运行：

```cmd
docker run -d --name iot_postgres -e POSTGRES_DB=iot_card_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -p 5432:5432 postgres:15-alpine
```

等待10秒，然后初始化数据库：

```cmd
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\001_create_admin_users.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\002_create_sim_cards.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\003_create_recharge_records.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\004_create_system_config.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\insert_test_data.sql
```

### 方式B：使用本地PostgreSQL

如果已安装PostgreSQL：

```cmd
REM 创建数据库
createdb -U postgres iot_card_db

REM 执行迁移
cd backend\migrations
psql -U postgres -d iot_card_db -f 001_create_admin_users.sql
psql -U postgres -d iot_card_db -f 002_create_sim_cards.sql
psql -U postgres -d iot_card_db -f 003_create_recharge_records.sql
psql -U postgres -d iot_card_db -f 004_create_system_config.sql
psql -U postgres -d iot_card_db -f insert_test_data.sql
cd ..\..
```

---

## 第二步：启动后端

**打开第一个CMD窗口：**

```cmd
cd D:\workspace\backend
go run cmd/server/main.go
```

**看到这个说明成功：**
```
数据库连接成功
服务器启动在端口 8080
[GIN-debug] Listening and serving HTTP on :8080
```

**保持这个窗口运行！**

---

## 第三步：启动H5前端

**打开第二个CMD窗口：**

```cmd
cd D:\workspace\frontend\h5
npm install
npm run dev
```

**看到这个说明成功：**
```
VITE v5.x.x  ready in xxx ms
➜  Local:   http://localhost:3000/
```

**保持这个窗口运行！**

**在浏览器打开：** http://localhost:3000

---

## 第四步：启动管理端

**打开第三个CMD窗口：**

```cmd
cd D:\workspace\frontend\admin
npm install
npm run dev
```

**看到这个说明成功：**
```
VITE v5.x.x  ready in xxx ms
➜  Local:   http://localhost:3001/
```

**保持这个窗口运行！**

**在浏览器打开：** http://localhost:3001

---

## 现在开始测试！

### 测试H5 (http://localhost:3000)

1. 打开浏览器访问 http://localhost:3000
2. 输入测试卡号：`89860123456789012345`
3. 点击"查询"
4. 应该能看到卡片信息
5. 点击"充值续费"看充值页面

### 测试管理端 (http://localhost:3001)

1. 打开浏览器访问 http://localhost:3001
2. 输入账号：`admin`
3. 输入密码：`admin123`
4. 点击"登录"
5. 应该能看到首页统计数据
6. 点击左侧"卡片管理"看卡片列表

---

## 如果遇到问题

### 问题1：后端启动失败 - "连接数据库失败"

**检查数据库是否运行：**
```cmd
docker ps
```

应该能看到 `iot_postgres` 容器在运行

**如果没有，启动它：**
```cmd
docker start iot_postgres
```

### 问题2：前端启动失败 - "npm install失败"

**使用国内镜像：**
```cmd
npm config set registry https://registry.npmmirror.com
npm install
```

### 问题3：前端页面空白

1. 按F12打开浏览器开发者工具
2. 查看Console标签的错误信息
3. 查看Network标签，确认请求是否发送

**最常见原因：后端没启动**
- 确认后端窗口显示"服务器启动在端口 8080"
- 在浏览器访问：http://localhost:8080/api/v1/card/query?keyword=test
- 应该能看到JSON响应

### 问题4：端口被占用

**查看哪个程序占用端口：**
```cmd
netstat -ano | findstr :8080
netstat -ano | findstr :3000
netstat -ano | findstr :3001
```

**修改端口：**
- 后端端口：编辑 `backend\configs\config.yaml`，修改 `port: 8080`
- H5端口：编辑 `frontend\h5\vite.config.js`，修改 `port: 3000`
- 管理端端口：编辑 `frontend\admin\vite.config.js`，修改 `port: 3001`

---

## 停止服务

**要停止服务，直接关闭对应的CMD窗口即可！**

或者在窗口中按 `Ctrl+C`

---

## 下次启动

数据库只需启动一次，下次只需要：

```cmd
REM 窗口1：启动后端
cd D:\workspace\backend
go run cmd/server/main.go

REM 窗口2：启动H5
cd D:\workspace\frontend\h5
npm run dev

REM 窗口3：启动管理端
cd D:\workspace\frontend\admin
npm run dev
```

---

## 测试数据

系统已经插入了5张测试卡片：

1. `89860123456789012345` - 正常状态（中国移动）
2. `89860123456789012346` - 即将到期（中国联通）
3. `89860123456789012347` - 已过期（中国电信）
4. `89860123456789012348` - 正常状态（中国移动）
5. `89860123456789012349` - 正常状态（中国联通）

你可以在H5端查询任意一张卡！

---

## 快捷命令汇总

**一键复制：初始化数据库**
```cmd
docker run -d --name iot_postgres -e POSTGRES_DB=iot_card_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -p 5432:5432 postgres:15-alpine && timeout /t 10 && docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\001_create_admin_users.sql && docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\002_create_sim_cards.sql && docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\003_create_recharge_records.sql && docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\004_create_system_config.sql && docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend\migrations\insert_test_data.sql
```

**测试后端是否运行：**
```cmd
curl http://localhost:8080/api/v1/card/query?keyword=test
```

或在浏览器直接打开：http://localhost:8080/api/v1/card/query?keyword=89860123456789012345
