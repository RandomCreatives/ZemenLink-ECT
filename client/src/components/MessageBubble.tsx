import React from 'react';
import { format } from 'date-fns';
import { CheckCheck } from 'lucide-react';

interface Message {
  id: string;
  content: string;
  sender_id: string;
  created_at: string;
}

interface MessageBubbleProps {
  message: Message;
  isSelf: boolean;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message, isSelf }) => {
  return (
    <div className={`flex ${isSelf ? 'justify-end' : 'justify-start'}`}>
      <div className={`max-w-[75%] rounded-xl px-3 py-1.5 shadow-sm relative ${
        isSelf ? 'bg-telegram-bubble rounded-tr-none' : 'bg-white rounded-tl-none'
      }`}>
        <p className="text-sm text-gray-800 break-words">{message.content}</p>
        <div className="flex items-center justify-end gap-1 text-[10px] text-gray-400 mt-1">
          <span>{format(new Date(message.created_at), 'HH:mm')}</span>
          {isSelf && <CheckCheck className="w-3 h-3 text-telegram-blue" />}
        </div>
      </div>
    </div>
  );
};

export default MessageBubble;
