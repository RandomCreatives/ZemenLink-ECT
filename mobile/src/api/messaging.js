import { useState, useEffect, useRef } from 'react';

// When testing on physical devices, replace 'localhost' with your machine's IP address.
const BASE_HOST = 'localhost:8080';
const API_URL = `http://${BASE_HOST}`;
const WS_URL = `ws://${BASE_HOST}/ws`;

export const useMessaging = (token) => {
  const [messages, setMessages] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef(null);

  useEffect(() => {
    if (!token) return;

    const socket = new WebSocket(`${WS_URL}?token=${token}`);

    socket.onopen = () => setIsConnected(true);
    socket.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        setMessages((prev) => [...prev, msg]);
      } catch (err) {
        console.error("Failed to parse message", err);
      }
    };
    socket.onclose = () => setIsConnected(false);

    socketRef.current = socket;
    return () => socket.close();
  }, [token]);

  const sendMessage = async (content, chatID = 'global-chat') => {
    if (!token) return;

    try {
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
    } catch (err) {
      console.error("Failed to send message", err);
      return false;
    }
  };

  return { messages, isConnected, sendMessage };
};
