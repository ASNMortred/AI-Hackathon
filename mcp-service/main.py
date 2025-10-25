from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
from typing import Optional, Dict, Any, List
import uvicorn
import logging
import sys
from datetime import datetime
import uuid
import os
try:
    from openai import OpenAI  # type: ignore
except Exception:
    OpenAI = None
import requests
import json

# Qiniu OpenAI wrapper that mimics OpenAI client's chat.completions.create
class ChatMessage:
    def __init__(self, content: str):
        self.content = content

class Choice:
    def __init__(self, message: 'ChatMessage'):
        self.message = message

class ChatCompletionsResponse:
    def __init__(self, choices: List['Choice']):
        self.choices = choices

class QiniuOpenAIClient:
    def __init__(self, api_key: str, base_url: str):
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.chat = self.Chat(self)

    class Chat:
        def __init__(self, parent: 'QiniuOpenAIClient'):
            self._parent = parent
            self.completions = self.Completions(parent)

        class Completions:
            def __init__(self, parent: 'QiniuOpenAIClient'):
                self._parent = parent

            def create(self, model: str, messages: List[Dict[str, str]], temperature: Optional[float] = None, max_tokens: Optional[int] = None, stream: bool = False):
                url = f"{self._parent.base_url}/chat/completions"
                headers = {
                    "Authorization": f"Bearer {self._parent.api_key}",
                    "Content-Type": "application/json",
                }
                payload: Dict[str, Any] = {
                    "stream": bool(stream),
                    "model": model,
                    "messages": messages,
                }
                if temperature is not None:
                    payload["temperature"] = temperature
                if max_tokens is not None:
                    payload["max_tokens"] = max_tokens

                try:
                    resp = requests.post(url, json=payload, headers=headers, timeout=30)
                    resp.raise_for_status()
                    data = resp.json()
                except Exception as e:
                    return ChatCompletionsResponse([Choice(ChatMessage(f"调用七牛AI出错: {str(e)}"))])

                content_text: str = ""
                if isinstance(data, dict):
                    if "choices" in data and isinstance(data["choices"], list) and data["choices"]:
                        ch0 = data["choices"][0]
                        if isinstance(ch0, dict):
                            if "message" in ch0 and isinstance(ch0["message"], dict):
                                content_text = ch0["message"].get("content", "") or ""
                            elif "text" in ch0:
                                content_text = ch0.get("text") or ""
                    if not content_text:
                        content_text = data.get("data") or data.get("message") or ""
                        if not content_text:
                            try:
                                content_text = json.dumps(data, ensure_ascii=False)
                            except Exception:
                                content_text = str(data)

                return ChatCompletionsResponse([Choice(ChatMessage(content_text))])

os.makedirs('./logs', exist_ok=True)
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler('./logs/mcp-service.log')
    ]
)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Chat MCP Service",
    description="AI Chat Service with MCP Protocol Support",
    version="1.0.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

sessions: Dict[str, List[Dict[str, str]]] = {}

openai_api_key = os.getenv("QiNiu_Ai_Key")
openai_base_url = os.getenv("OPENAI_BASE_URL", "https://openai.qiniu.com/v1")
openai_model = os.getenv("OPENAI_MODEL", "qwen-vl-max-2025-01-25")

if openai_api_key:
    client = QiniuOpenAIClient(api_key=openai_api_key, base_url=openai_base_url)
    logger.info(f"Qiniu OpenAI client initialized with model: {openai_model}")
else:
    client = None
    logger.warning("QiNiu_Ai_Key not set, using mock responses")

class ChatRequest(BaseModel):
    message: str = Field(..., description="User input message")
    session_id: Optional[str] = Field(None, description="Session ID for conversation context")
    temperature: Optional[float] = Field(0.7, ge=0.0, le=2.0, description="Temperature for response generation")

class ChatResponse(BaseModel):
    success: bool
    data: Optional[Dict[str, Any]] = None
    error: Optional[str] = None

class MCPTool(BaseModel):
    name: str
    description: str
    parameters: Dict[str, Any]

mcp_tools: List[MCPTool] = [
    MCPTool(
        name="get_current_time",
        description="Get the current time",
        parameters={
            "type": "object",
            "properties": {},
            "required": []
        }
    ),
    MCPTool(
        name="search_knowledge",
        description="Search knowledge base",
        parameters={
            "type": "object",
            "properties": {
                "query": {"type": "string", "description": "Search query"}
            },
            "required": ["query"]
        }
    )
]

def execute_tool(tool_name: str, parameters: Dict[str, Any]) -> str:
    if tool_name == "get_current_time":
        return datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    elif tool_name == "search_knowledge":
        query = parameters.get("query", "")
        return f"Search results for: {query} (Mock implementation)"
    return "Unknown tool"

def generate_response(message: str, session_id: str, temperature: float) -> str:
    if session_id not in sessions:
        sessions[session_id] = []
    
    sessions[session_id].append({"role": "user", "content": message})
    
    if len(sessions[session_id]) > 20:
        sessions[session_id] = sessions[session_id][-20:]
    
    if client:
        try:
            response = client.chat.completions.create(
                model=openai_model,
                messages=sessions[session_id],
                temperature=temperature,
                max_tokens=1000
            )
            ai_response = response.choices[0].message.content
            sessions[session_id].append({"role": "assistant", "content": ai_response})
            logger.info(f"Generated response for session {session_id}")
            return ai_response
        except Exception as e:
            logger.error(f"Error calling OpenAI API: {str(e)}")
            return f"抱歉，我遇到了一些问题: {str(e)}"
    else:
        ai_response = f"收到您的消息: {message}\n\n这是一个模拟响应，因为未配置OpenAI API密钥。"
        sessions[session_id].append({"role": "assistant", "content": ai_response})
        return ai_response

@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "service": "mcp-service"
    }

@app.get("/mcp/tools")
async def list_tools():
    return {
        "success": True,
        "tools": [tool.dict() for tool in mcp_tools]
    }

@app.post("/api/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    try:
        logger.info(f"Received chat request - message: {request.message[:50]}..., session_id: {request.session_id}")
        
        session_id = request.session_id or str(uuid.uuid4())
        
        response_text = generate_response(
            message=request.message,
            session_id=session_id,
            temperature=request.temperature
        )
        
        logger.info(f"Generated response for session {session_id}")
        
        return ChatResponse(
            success=True,
            data={
                "response": response_text,
                "session_id": session_id
            },
            error=None
        )
    except Exception as e:
        logger.error(f"Error processing chat request: {str(e)}", exc_info=True)
        return ChatResponse(
            success=False,
            data=None,
            error=str(e)
        )

@app.get("/")
async def root():
    return {
        "service": "Chat MCP Service",
        "version": "1.0.0",
        "status": "running"
    }

if __name__ == "__main__":
    port = int(os.getenv("MCP_SERVICE_PORT", "8000"))
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=port,
        log_level="info"
    )

# Qiniu OpenAI wrapper that mimics OpenAI client's chat.completions.create
class ChatMessage:
    def __init__(self, content: str):
        self.content = content

class Choice:
    def __init__(self, message: 'ChatMessage'):
        self.message = message

class ChatCompletionsResponse:
    def __init__(self, choices: List['Choice']):
        self.choices = choices

class QiniuOpenAIClient:
    def __init__(self, api_key: str, base_url: str):
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.chat = self.Chat(self)

    class Chat:
        def __init__(self, parent: 'QiniuOpenAIClient'):
            self._parent = parent
            self.completions = self.Completions(parent)

        class Completions:
            def __init__(self, parent: 'QiniuOpenAIClient'):
                self._parent = parent

            def create(self, model: str, messages: List[Dict[str, str]], temperature: Optional[float] = None, max_tokens: Optional[int] = None, stream: bool = False):
                url = f"{self._parent.base_url}/chat/completions"
                headers = {
                    "Authorization": f"Bearer {self._parent.api_key}",
                    "Content-Type": "application/json",
                }
                payload: Dict[str, Any] = {
                    "stream": bool(stream),
                    "model": model,
                    "messages": messages,
                }
                if temperature is not None:
                    payload["temperature"] = temperature
                if max_tokens is not None:
                    payload["max_tokens"] = max_tokens

                try:
                    resp = requests.post(url, json=payload, headers=headers, timeout=30)
                    resp.raise_for_status()
                    data = resp.json()
                except Exception as e:
                    return ChatCompletionsResponse([Choice(ChatMessage(f"调用七牛AI出错: {str(e)}"))])

                content_text: str = ""
                if isinstance(data, dict):
                    if "choices" in data and isinstance(data["choices"], list) and data["choices"]:
                        ch0 = data["choices"][0]
                        if isinstance(ch0, dict):
                            if "message" in ch0 and isinstance(ch0["message"], dict):
                                content_text = ch0["message"].get("content", "") or ""
                            elif "text" in ch0:
                                content_text = ch0.get("text") or ""
                    if not content_text:
                        content_text = data.get("data") or data.get("message") or ""
                        if not content_text:
                            try:
                                content_text = json.dumps(data, ensure_ascii=False)
                            except Exception:
                                content_text = str(data)

                return ChatCompletionsResponse([Choice(ChatMessage(content_text))])
