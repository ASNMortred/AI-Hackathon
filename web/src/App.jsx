import React, { useState, useEffect, useRef } from 'react';
import { Input, Button, message, Form, Tabs, Card, Row, Col } from 'antd';
import { AudioOutlined, SendOutlined, UserOutlined, LockOutlined, UploadOutlined } from '@ant-design/icons';
import './App.css';

const { TabPane } = Tabs;

const App = () => {
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [memoryId, setMemoryId] = useState('');
  const [isListening, setIsListening] = useState(false);
  const [isPlaying, setIsPlaying] = useState({});
  const [isWaitingForVoiceResponse, setIsWaitingForVoiceResponse] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [username, setUsername] = useState('');
  const [token, setToken] = useState('');

  const recognitionRef = useRef(null);
  const messageInputRef = useRef(null);
  const messagesEndRef = useRef(null);
  const fileInputRef = useRef(null);

  // 浏览器语音功能支持性检测
  const speechSupported = 'speechSynthesis' in window;
  const recognitionSupported = 'webkitSpeechRecognition' in window;

  // 滚动到底部
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // 移除HTML标签，保留纯文本
  const stripHtml = (html) => {
    const doc = new DOMParser().parseFromString(html, 'text/html');
    return doc.body.textContent || '';
  };

  // 语音朗读消息内容
  const readMessage = (content, index) => {
    window.speechSynthesis.cancel();

    if (!('speechSynthesis' in window)) {
      setErrorMessage('您的浏览器不支持语音朗读功能（speechSynthesis）。建议使用最新版 Chrome、Safari 或 Edge。');
      return;
    }

    setIsPlaying({ ...isPlaying, [index]: true });

    const text = stripHtml(content);
    if (!text.trim()) return;

    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = 'zh-CN';

    window.speechSynthesis.speak(utterance);

    utterance.onend = () => {
      setIsPlaying((prev) => ({ ...prev, [index]: false }));
    };

    utterance.onerror = (event) => {
      setIsPlaying((prev) => ({ ...prev, [index]: false }));
      setErrorMessage('语音朗读失败: ' + (event.error || '未知错误') + '。建议检查浏览器语音包、网络环境或重启浏览器。');
      console.error('语音朗读错误:', event);
    };
  };

  // 格式化消息内容，将Markdown表格转换为HTML表格
  const formatMessageContent = (content) => {
    if (!content) return '';

    const tableRegex = /^\|.*\|(\r?\n\|.*\|)+/gm;
    if (tableRegex.test(content)) {
      let html = '<table class="bot-table">';
      const lines = content.split('\n');
      let i = 0;

      while (i < lines.length) {
        const line = lines[i].trim();
        if (line.startsWith('|')) {
          const headerCells = line.split('|').filter((cell) => cell.trim() !== '');
          html += '<tr>';
          headerCells.forEach((cell) => {
            html += `<th>${cell.trim()}</th>`;
          });
          html += '</tr>';
          i++;

          if (i < lines.length && lines[i].trim().startsWith('|') && /^[\|:-]+$/.test(lines[i].trim())) {
            i++;
          }

          while (i < lines.length && lines[i].trim().startsWith('|')) {
            const dataLine = lines[i].trim();
            const dataCells = dataLine.split('|').filter((cell) => cell.trim() !== '');
            html += '<tr>';
            dataCells.forEach((cell) => {
              html += `<td>${cell.trim()}</td>`;
            });
            html += '</tr>';
            i++;
          }
        } else {
          if (line) html += `<p>${line}</p>`;
          i++;
        }
      }

      html += '</table>';
      return html;
    }

    return content.replace(/\n/g, '<br>');
  };

  // 生成唯一memoryId
  const generateMemoryId = () => {
    const id = Date.now().toString();
    localStorage.setItem('memoryId', id);
    return id;
  };

  // 无障碍辅助函数
  const announce = (text) => {
    if ('speechSynthesis' in window) {
      const utterance = new SpeechSynthesisUtterance(text);
      utterance.lang = 'zh-CN';
      window.speechSynthesis.speak(utterance);
    }
  };

  // 初始化
  useEffect(() => {
    // 检查本地存储中是否有登录信息
    const savedToken = localStorage.getItem('token');
    const savedUsername = localStorage.getItem('username');
    if (savedToken && savedUsername) {
      setIsLoggedIn(true);
      setToken(savedToken);
      setUsername(savedUsername);
    }

    const savedId = localStorage.getItem('memoryId');
    setMemoryId(savedId || generateMemoryId());

    // 页面加载时的语音提示（仅在已登录时显示）
    if (isLoggedIn && 'speechSynthesis' in window) {
      const welcomeMessage = new SpeechSynthesisUtterance('欢迎使用语音输入功能，是否启用语音输入？');
      welcomeMessage.lang = 'zh-CN';
      window.speechSynthesis.speak(welcomeMessage);
      welcomeMessage.onend = () => {
        const promptMessage = new SpeechSynthesisUtterance('请说出启用或取消来选择是否使用语音输入功能');
        promptMessage.lang = 'zh-CN';
        window.speechSynthesis.speak(promptMessage);
        promptMessage.onend = () => {
          setIsWaitingForVoiceResponse(true);
          if (recognitionRef.current) {
            recognitionRef.current.start();
          }
        };
      };
    }
  }, [isLoggedIn]);

  // 初始化语音识别
  useEffect(() => {
    if ('webkitSpeechRecognition' in window) {
      const recognition = new window.webkitSpeechRecognition();
      recognition.continuous = false;
      recognition.interimResults = false;
      recognition.lang = 'zh-CN';

      recognition.onresult = (event) => {
        const transcript = event.results[0][0].transcript;
        if (isWaitingForVoiceResponse) {
          setIsWaitingForVoiceResponse(false);
          recognition.stop();
          const lowerTranscript = transcript.toLowerCase();
          if (
            lowerTranscript.includes('启用') ||
            lowerTranscript.includes('是') ||
            lowerTranscript.includes('开启')
          ) {
            toggleVoiceInput();
            message.success('已启用语音输入');
          } else if (
            lowerTranscript.includes('取消') ||
            lowerTranscript.includes('不') ||
            lowerTranscript.includes('否')
          ) {
            message.info('已取消语音输入');
          } else {
            const retryMessage = new SpeechSynthesisUtterance('未识别到指令，请重试');
            retryMessage.lang = 'zh-CN';
            window.speechSynthesis.speak(retryMessage);
          }
        } else {
          setInputMessage(transcript);
          if (messageInputRef.current) {
            messageInputRef.current.focus();
          }
        }
      };

      recognition.onerror = (event) => {
        let suggestion = '';
        if (event.error === 'not-allowed') {
          suggestion = '（未授权麦克风，请检查浏览器地址栏左侧的权限设置）';
        } else if (event.error === 'network') {
          suggestion = '（网络错误，语音识别依赖 Google 服务，国内网络可能无法访问）';
        } else if (event.error === 'no-speech') {
          suggestion = '（未检测到语音，请重试）';
        } else if (event.error === 'aborted') {
          suggestion = '（语音识别被中断，请重试）';
        }
        setErrorMessage('语音识别失败: ' + event.error + suggestion);
        console.error('语音识别错误:', event.error, event);
        setIsListening(false);
      };

      recognition.onend = () => {
        if (isListening) {
          recognition.start();
        }
      };

      recognitionRef.current = recognition;
    }
  }, [isWaitingForVoiceResponse, isListening]);

  // 切换语音输入状态
  const toggleVoiceInput = () => {
    if (!recognitionRef.current) {
      setErrorMessage('您的浏览器不支持语音识别');
      return;
    }

    if (isListening) {
      recognitionRef.current.stop();
      setIsListening(false);
      announce('语音输入已停止');
    } else {
      setInputMessage('');
      if (messageInputRef.current) {
        messageInputRef.current.focus();
      }
      setErrorMessage('');

      navigator.mediaDevices
        .getUserMedia({ audio: true })
        .then(() => {
          setIsListening(true);
          recognitionRef.current.start();
          announce('语音输入已启动，请开始说话');
        })
        .catch((err) => {
          console.error('麦克风权限请求失败:', err);
          setErrorMessage('无法访问麦克风: 请在浏览器设置中启用麦克风权限');
          setIsListening(false);
        });
    }
  };

  // 发送消息
  const sendMessage = async () => {
    if (isListening && recognitionRef.current) {
      recognitionRef.current.stop();
      setIsListening(false);
    }

    if (!inputMessage.trim() || loading) return;

    setMessages((prev) => [...prev, { content: inputMessage, isUser: true }]);
    announce('您的消息已发送');

    const userMessage = inputMessage;
    setInputMessage('');
    setLoading(true);
    setErrorMessage('');

    try {
      const botMessageIndex = messages.length + 1;
      setMessages((prev) => [...prev, { content: '', isUser: false }]);

      const response = await fetch('/api/v1/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}` // 添加认证头
        },
        body: JSON.stringify({
          memoryId: memoryId,
          message: userMessage,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder('utf-8');
      let done = false;

      let streamErrorOccurred = false;
      while (!done) {
        try {
          const { value, done: doneReading } = await reader.read();
          done = doneReading;
          if (value) {
            const chunk = decoder.decode(value, { stream: true });
            setMessages((prev) => {
              const newMessages = [...prev];
              newMessages[botMessageIndex].content += chunk;
              return newMessages;
            });
          }
        } catch (streamError) {
          console.error('Stream reading error:', streamError);
          setErrorMessage('流式传输失败，请重试');
          streamErrorOccurred = true;
          done = true;
        }
      }

      if (!streamErrorOccurred) {
        setErrorMessage('');
      }
    } catch (error) {
      setErrorMessage('请求失败，请重试');
      console.error('请求错误:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUploadClick = () => {
    if (!isLoggedIn) {
      message.warning('请先登录');
      return;
    }
    fileInputRef.current?.click();
  };

  const handleFileSelect = async (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    try {
      const resp = await fetch('/api/v1/upload', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });
      const data = await resp.json();
      if (resp.ok) {
        message.success('上传成功');
        setMessages((prev) => [...prev, { content: `已上传文件：${file.name}\n链接：${data.url || ''}` , isUser: true }]);
      } else {
        message.error(data.error || '上传失败');
      }
    } catch (err) {
      console.error('上传错误:', err);
      message.error('上传失败，请重试');
    } finally {
      e.target.value = '';
    }
  };

  // 新会话
  const newSession = () => {
    setMemoryId(generateMemoryId());
    setMessages([]);
    setErrorMessage('');
    announce('新会话已开始，聊天记录已清空');
    message.success('已开始新会话');
  };

  // 用户注册
  const handleRegister = async (values) => {
    try {
      const response = await fetch('/api/v1/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: values.username,
          password: values.password,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        message.success('注册成功，请登录');
      } else {
        message.error(data.error || '注册失败');
      }
    } catch (error) {
      console.error('注册错误:', error);
      message.error('注册失败，请重试');
    }
  };

  // 用户登录
  const handleLogin = async (values) => {
    try {
      const response = await fetch('/api/v1/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: values.username,
          password: values.password,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // 保存登录信息到本地存储
        localStorage.setItem('token', data.token);
        localStorage.setItem('username', data.username);
        
        // 更新状态
        setIsLoggedIn(true);
        setToken(data.token);
        setUsername(data.username);
        
        message.success('登录成功');
      } else {
        message.error(data.error || '登录失败');
      }
    } catch (error) {
      console.error('登录错误:', error);
      message.error('登录失败，请重试');
    }
  };

  // 用户登出
  const handleLogout = () => {
    // 清除本地存储的登录信息
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    
    // 更新状态
    setIsLoggedIn(false);
    setToken('');
    setUsername('');
    
    message.success('已退出登录');
  };

  // 如果用户未登录，显示登录/注册页面
  if (!isLoggedIn) {
    return (
      <div className="app-layout" style={{ justifyContent: 'center', alignItems: 'center' }}>
        <div style={{ width: '400px' }}>
          <Card title="云觅" bordered={false} style={{ width: '100%' }}>
            <Tabs defaultActiveKey="1">
              <TabPane tab="登录" key="1">
                <Form onFinish={handleLogin}>
                  <Form.Item
                    name="username"
                    rules={[{ required: true, message: '请输入用户名!' }]}
                  >
                    <Input prefix={<UserOutlined />} placeholder="用户名" />
                  </Form.Item>
                  <Form.Item
                    name="password"
                    rules={[{ required: true, message: '请输入密码!' }]}
                  >
                    <Input prefix={<LockOutlined />} type="password" placeholder="密码" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                      登录
                    </Button>
                  </Form.Item>
                </Form>
              </TabPane>
              <TabPane tab="注册" key="2">
                <Form onFinish={handleRegister}>
                  <Form.Item
                    name="username"
                    rules={[{ required: true, message: '请输入用户名!' }]}
                  >
                    <Input prefix={<UserOutlined />} placeholder="用户名" />
                  </Form.Item>
                  <Form.Item
                    name="password"
                    rules={[{ required: true, message: '请输入密码!' }]}
                  >
                    <Input prefix={<LockOutlined />} type="password" placeholder="密码" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                      注册
                    </Button>
                  </Form.Item>
                </Form>
              </TabPane>
            </Tabs>
          </Card>
        </div>
      </div>
    );
  }

  // 用户已登录，显示聊天界面
  return (
    <div className="app-layout">
      {/* 左侧侧边栏 */}
      <aside className="sidebar">
        <div className="avatar-block">
          <img src="/img/avatar.png" alt="云觅" className="avatar-img" />
          <div className="system-title">云觅</div>
          <div style={{ marginTop: '10px', fontSize: '14px' }}>欢迎, {username}!</div>
        </div>
        <Button className="new-session-btn" onClick={newSession} aria-label="新会话">
          + 新会话
        </Button>
        <Button 
          style={{ marginTop: '10px', width: '180px' }} 
          onClick={handleLogout}
        >
          退出登录
        </Button>
      </aside>

      {/* 右侧主内容区 */}
      <main id="main-content" className="main-content" tabIndex="-1">
        <section className="chat-container" aria-label="聊天消息" role="region">
          <div className="message-list" aria-live="polite" role="list">
            {messages.map((msg, index) => (
              <div key={index} className="message-item" role="listitem">
                <div
                  className={`message-bubble ${msg.isUser ? 'user-message' : 'bot-message'}`}
                  dangerouslySetInnerHTML={{ __html: formatMessageContent(msg.content) }}
                ></div>
                {!msg.isUser && speechSupported && (
                  <Button
                    className="read-button"
                    size="small"
                    onClick={() => readMessage(msg.content, index)}
                    disabled={isPlaying[index] || !speechSupported}
                    aria-label="朗读消息内容"
                    tabIndex="0"
                  >
                    播放
                  </Button>
                )}
              </div>
            ))}
            {errorMessage && <div className="error-message" aria-live="assertive">{errorMessage}</div>}
            <div ref={messagesEndRef} />
          </div>

          <div className="input-area">
            <Input
              ref={messageInputRef}
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onPressEnter={sendMessage}
              placeholder="请输入消息..."
              disabled={loading}
              aria-label="消息输入框"
              aria-describedby="input-hint"
              tabIndex="0"
            />
            {recognitionSupported && (
              <Button
                onClick={toggleVoiceInput}
                loading={isListening}
                type={isListening ? 'primary' : 'default'}
                aria-label={`切换语音输入，当前状态：${isListening ? '正在聆听' : '已停止'}`}
                tabIndex="0"
                title="语音输入 (Alt+V)"
                className="mic-button"
                icon={<AudioOutlined />}
              />
            )}
            <input ref={fileInputRef} type="file" style={{ display: 'none' }} onChange={handleFileSelect} />
            <Button
              onClick={handleUploadClick}
              aria-label="上传文件"
              tabIndex="0"
              title="上传文件"
              className="upload-button"
              icon={<UploadOutlined />}
            />
            <Button
              type="primary"
              onClick={sendMessage}
              loading={loading}
              aria-label="发送消息"
              tabIndex="0"
              title="发送消息 (Alt+S)"
              className="send-button"
              icon={<SendOutlined />}
            >
              发送
            </Button
            >
          </div>
        </section>
      </main>
    </div>
  );
};

export default App;


  const handleUploadClick = () => {
    if (!isLoggedIn) {
      message.warning('请先登录');
      return;
    }
    fileInputRef.current?.click();
  };

  const handleFileSelect = async (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    try {
      const resp = await fetch('/api/v1/upload', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });
      const data = await resp.json();
      if (resp.ok) {
        message.success('上传成功');
        setMessages((prev) => [...prev, { content: `已上传文件：${file.name}\n链接：${data.url || ''}` , isUser: true }]);
      } else {
        message.error(data.error || '上传失败');
      }
    } catch (err) {
      console.error('上传错误:', err);
      message.error('上传失败，请重试');
    } finally {
      e.target.value = '';
    }
  };