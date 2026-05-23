import { useState, useEffect, useRef } from 'react';

const API_URL = 'http://localhost:8080';
const WS_URL = 'ws://localhost:8080/ws';

export const useMessaging = (token) => {
  const [messages, setMessages] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef(null);

  useEffect(() => {
    if (!token) return;

    const socket = new WebSocket(`${WS_URL}?token=${token}`);

    socket.onopen = () => setIsConnected(true);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prev) => [...prev, msg]);
    };
    socket.onclose = () => setIsConnected(false);

    socketRef.current = socket;
    return () => socket.close();
  }, [token]);

  const sendMessage = async (content, chatID = 'global-chat') => {
    if (!token) return;

    const response = await fetch(`${API_URL}/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        chat_id: chatID,
        content: content,
        content_type: 'text'
      })
    });

    return response.ok;
  };

  return { messages, isConnected, sendMessage };
};
