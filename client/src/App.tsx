import React, { useState, useEffect, useRef } from 'react';
import { Send, Menu, Search, MoreVertical, Paperclip, Smile, CheckCheck, ShieldCheck } from 'lucide-react';
import { format } from 'date-fns';

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
  const [tenantID] = useState('tenant-abc');
  const [userID] = useState('user-123');
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

  // Simulate obtaining a JWT token for the session
  useEffect(() => {
    const fetchToken = async () => {
      // Use the actual token generated for the demo
      setToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk0ODM1MDEsInJvbGUiOiJhZG1pbiIsInRlbmFudF9pZCI6InRlbmFudC1hYmMiLCJ1c2VyX2lkIjoidXNlci0xMjMifQ.EDo_edjav7O2Rb0IYYKVXOz3LtExlS-a8m1omoeFwJU");
    };
    fetchToken();
  }, []);

  // WebSocket Connection
  useEffect(() => {
    if (!token) return;

    // Use absolute URL for the sandbox environment
    const wsUrl = `ws://${window.location.hostname}:8080/ws?token=${token}`;
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('Connected to ZemenLink Kernel');
      setIsConnected(true);
    };

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prev) => [...prev, msg]);
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
      metadata: { extra: { source: 'web-client' } }
    };

    try {
      const response = await fetch('http://localhost:8080/messages', {
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
      {/* Sidebar */}
      <div className="w-80 bg-white border-r flex flex-col hidden md:flex">
        <div className="p-4 border-b flex items-center gap-4 bg-telegram-blue text-white">
          <Menu className="cursor-pointer" />
          <span className="font-bold text-lg">ZemenLink</span>
        </div>

        <div className="p-3">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 text-gray-400 w-4 h-4" />
            <input
              type="text"
              placeholder="Search"
              className="w-full bg-gray-100 rounded-lg py-1.5 pl-10 pr-4 text-sm focus:outline-none"
            />
          </div>
        </div>

        <div className="flex-1 overflow-y-auto">
          <div className="p-4 flex items-center gap-3 hover:bg-gray-50 cursor-pointer">
            <div className="w-12 h-12 rounded-full bg-telegram-blue flex items-center justify-center text-white font-bold">ZL</div>
            <div className="flex-1 overflow-hidden">
              <div className="flex justify-between text-sm">
                <span className="font-semibold truncate">Enterprise Kernel</span>
                <span className="text-gray-500">Just now</span>
              </div>
              <p className="text-sm text-gray-500 truncate">Real-time sync enabled.</p>
            </div>
          </div>
        </div>

        <div className="p-4 border-t bg-gray-50 flex items-center gap-2 text-xs text-gray-500">
          <ShieldCheck className="w-4 h-4 text-green-500" />
          <span>Siloed Isolation: {tenantID}</span>
        </div>
      </div>

      {/* Main Chat */}
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

        {/* Messages Window */}
        <div className="flex-1 overflow-y-auto p-4 space-y-2 custom-scrollbar">
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex ${msg.sender_id === userID ? 'justify-end' : 'justify-start'}`}
            >
              <div className={`max-w-[75%] rounded-xl px-3 py-1.5 shadow-sm relative ${
                msg.sender_id === userID ? 'bg-telegram-bubble rounded-tr-none' : 'bg-white rounded-tl-none'
              }`}>
                <p className="text-sm text-gray-800 break-words">{msg.content}</p>
                <div className="flex items-center justify-end gap-1 text-[10px] text-gray-400 mt-1">
                  <span>{format(new Date(msg.created_at), 'HH:mm')}</span>
                  {msg.sender_id === userID && <CheckCheck className="w-3 h-3 text-telegram-blue" />}
                </div>
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>

        {/* Input Area */}
        <div className="p-4 bg-white md:bg-transparent">
          <div className="max-w-3xl mx-auto flex items-end gap-2">
            <div className="flex-1 bg-white rounded-xl shadow-sm flex items-end p-2 border border-gray-200">
              <Smile className="text-gray-400 m-2 cursor-pointer" />
              <textarea
                rows={1}
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && !e.shiftKey && (e.preventDefault(), handleSend())}
                placeholder="Write a message..."
                className="flex-1 py-2 px-1 text-sm focus:outline-none resize-none max-h-32"
              />
              <Paperclip className="text-gray-400 m-2 cursor-pointer rotate-45" />
            </div>
            <button
              onClick={handleSend}
              disabled={!isConnected}
              className={`w-12 h-12 rounded-full flex items-center justify-center text-white shadow-md transition-all flex-shrink-0 ${
                isConnected ? 'bg-telegram-blue hover:scale-105' : 'bg-gray-300 cursor-not-allowed'
              }`}
            >
              <Send className="w-5 h-5 ml-0.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default App;
