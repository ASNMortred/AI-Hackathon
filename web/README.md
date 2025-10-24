# 云觅 - AI 对话系统

这是一个基于 React + Vite + Ant Design 的 AI 对话系统，支持语音输入和朗读功能。

## 项目特性

- ✨ 基于 React 18 构建
- 🎨 使用 Ant Design 5 组件库
- 🎤 支持语音输入功能
- 🔊 支持消息朗读功能
- 💬 流式消息显示
- 📱 响应式设计
- ♿ 无障碍功能支持

## 技术栈

- **框架**: React 18.2.0
- **构建工具**: Vite 4.4.0
- **UI 组件库**: Ant Design 5.12.0
- **图标**: @ant-design/icons 5.2.6
- **HTTP 客户端**: Axios 1.6.8

## 项目结构

```
web/
├── src/
│   ├── App.jsx          # 主应用组件
│   ├── App.css          # 主应用样式
│   ├── main.jsx         # 应用入口
│   └── index.css        # 全局样式
├── img/                 # 图片资源
│   ├── avatar.png       # 头像图片
│   └── 麦克风.png       # 麦克风图标（备用）
├── index.html           # HTML 入口
├── vite.config.js       # Vite 配置
└── package.json         # 项目配置

```

## 安装依赖

```bash
npm install
```

## 运行项目

### 开发模式

```bash
npm run dev
```

开发服务器将在 http://localhost:5173 启动

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 功能说明

### 1. 聊天对话
- 在输入框中输入消息，点击"发送"按钮或按 Enter 键发送
- 支持流式消息接收，实时显示 AI 回复

### 2. 语音输入
- 点击麦克风按钮开启语音输入（需要浏览器支持）
- 说话内容会自动转换为文字
- 再次点击麦克风按钮关闭语音输入

### 3. 消息朗读
- 点击 AI 消息旁边的"播放"按钮可朗读消息内容
- 支持中文语音合成

### 4. 新会话
- 点击左侧"+ 新会话"按钮开始新的对话
- 会话历史会被清空

## 浏览器支持

- Chrome 60+
- Safari 14+
- Edge 79+
- Firefox 60+

**注意**: 语音功能需要浏览器支持 Web Speech API

## API 配置

项目使用代理配置将 `/api` 前缀的请求转发到后端服务器：

```javascript
// vite.config.js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8081',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```

如需修改后端地址，请编辑 `vite.config.js` 文件。

## 环境要求

- Node.js 14+ (推荐 16+)
- npm 6+ 或 yarn 1.22+

## 许可证

MIT
