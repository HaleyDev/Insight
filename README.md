# Insight

前后端分离架构，后端基于 Go 开发，前端基于 Vue 3 开发。

## 技术栈

### 后端

Go 1.22 / Gin / GORM / Viper / JWT / Redis / Zap

### 前端

Vue 3 / Vuetify / Vite / TypeScript / Pinia / Axios

## 调试启动

### 后端

```bash
cd backend && go run main.go -c config/local
```

服务启动在 `http://localhost:8080`。

### 前端

```bash
cd frontend && npm install && npm run dev
```

服务启动在 `http://localhost:5173`，`/api` 请求自动代理到后端 `http://localhost:8080`。