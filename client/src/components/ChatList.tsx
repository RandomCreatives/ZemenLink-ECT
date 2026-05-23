import React from 'react';

interface ChatItem {
  id: string;
  name: string;
  lastMessage: string;
  time: string;
  unreadCount?: number;
  initials: string;
}

interface ChatListProps {
  chats: ChatItem[];
  onSelect: (id: string) => void;
  activeID?: string;
}

const ChatList: React.FC<ChatListProps> = ({ chats, onSelect, activeID }) => {
  return (
    <div className="flex-1 overflow-y-auto">
      {chats.map((chat) => (
        <div
          key={chat.id}
          onClick={() => onSelect(chat.id)}
          className={`
            p-3 md:p-4 flex items-center gap-3 cursor-pointer transition-colors active:bg-gray-100
            ${activeID === chat.id ? 'bg-telegram-blue bg-opacity-10 border-r-2 border-telegram-blue' : 'hover:bg-gray-50'}
          `}
        >
          <div className="w-12 h-12 rounded-full bg-telegram-blue flex items-center justify-center text-white font-bold text-lg flex-shrink-0">
            {chat.initials}
          </div>
          <div className="flex-1 overflow-hidden">
            <div className="flex justify-between items-center mb-0.5">
              <span className="font-semibold text-gray-900 truncate">{chat.name}</span>
              <span className="text-xs text-gray-400 whitespace-nowrap">{chat.time}</span>
            </div>
            <div className="flex justify-between items-center">
              <p className="text-sm text-gray-500 truncate pr-2">{chat.lastMessage}</p>
              {chat.unreadCount && (
                <div className="bg-telegram-blue text-white text-[10px] font-bold rounded-full h-5 min-w-5 px-1.5 flex items-center justify-center">
                  {chat.unreadCount}
                </div>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ChatList;
