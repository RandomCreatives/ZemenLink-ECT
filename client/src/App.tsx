import React, { useState, useEffect, useRef } from 'react';
import { Search, MoreVertical } from 'lucide-react';
import ChatSidebar from './components/ChatSidebar';
import MessageBubble from './components/MessageBubble';
import InputBar from './components/InputBar';
import Login from './components/Login';
import './App.css';

interface Message {
  id: string;
  content: string;
  sender_id: string;
  created_at: string;
  metadata?: any;
}

const App: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([
    { id: '1', content: 'Welcome to ZemenLink! This is your enterprise messaging kernel.', sender_id: 'system', created_at: new Date().toISOString() },
  ]);
  const [inputValue, setInputValue] = useState('');
  const [tenantID, setTenantID] = useState<string | null>(null);
  const [userID, setUserID] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  const socketRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleLogin = (tenant: string, user: string, authToken: string) => {
    setTenantID(tenant);
    setUserID(user);
    setToken(authToken);
  };

  useEffect(() => {
    if (!token) return;

    const wsUrl = `ws://${window.location.hostname}:8080/ws?token=${token}`;
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('Connected to ZemenLink Kernel');
      setIsConnected(true);
    };

    socket.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        setMessages((prev) => [...prev, msg]);
      } catch (err) {
        console.error("Failed to parse message", err);
      }
    };

    socket.onclose = () => {
      console.log('Disconnected from ZemenLink Kernel');
      setIsConnected(false);
    };

    socketRef.current = socket;

    return () => {
      socket.close();
    };
  }, [token]);

  const handleSend = async () => {
    if (!inputValue.trim() || !token) return;

    const payload = {
      chat_id: 'global-chat',
      content: inputValue,
      content_type: 'text',
    };

    try {
      // PROXY PATH: Use /api prefix to hit Vite proxy
      const response = await fetch('/api/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      });

      if (response.ok) {
        setInputValue('');
      } else {
        console.error("Failed to send message");
      }
    } catch (err) {
      console.error("Error sending message", err);
    }
  };

  return (
    <div className="flex h-screen bg-gray-100 overflow-hidden">
      {!token && <Login onLogin={handleLogin} />}

      <ChatSidebar tenantID={tenantID || ''} />

      {/* Main Chat Window */}
      <div className="flex-1 flex flex-col min-w-0 bg-[#e7ebf0]">
        {/* Chat Header */}
        <div className="h-16 bg-white border-b px-4 flex items-center justify-between z-10 shadow-sm">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-telegram-blue flex items-center justify-center text-white font-bold">ZL</div>
            <div>
              <h2 className="font-semibold text-gray-900 leading-tight">Enterprise Kernel Chat</h2>
              <span className={`text-xs ${isConnected ? 'text-green-500' : 'text-red-500'}`}>
                {isConnected ? 'connected' : 'connecting...'}
              </span>
            </div>
          </div>
          <div className="flex items-center gap-4 text-gray-400">
            <Search className="w-5 h-5 cursor-pointer" />
            <MoreVertical className="w-5 h-5 cursor-pointer" />
          </div>
        </div>

        {/* Messages List */}
        <div className="flex-1 overflow-y-auto p-4 space-y-2 custom-scrollbar">
          {messages.map((msg) => (
            <MessageBubble
              key={msg.id}
              message={msg}
              isSelf={msg.sender_id === userID}
            />
          ))}
          <div ref={messagesEndRef} />
        </div>

        <InputBar
          value={inputValue}
          onChange={setInputValue}
          onSend={handleSend}
          disabled={!isConnected}
        />
      </div>
    </div>
  );
};

export default App;
