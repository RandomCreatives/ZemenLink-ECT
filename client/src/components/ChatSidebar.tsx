import React from 'react';
import { Menu, Search, ShieldCheck } from 'lucide-react';
import ChatList from './ChatList';

interface ChatSidebarProps {
  tenantID: string;
  isOpen: boolean;
  onSelectChat: (id: string) => void;
  activeChatID?: string;
}

const mockChats = [
  { id: 'global-chat', name: 'Enterprise Kernel Chat', lastMessage: 'Real-time sync enabled.', time: 'Just now', initials: 'ZL', unreadCount: 2 },
  { id: 'compliance-team', name: 'Compliance Team', lastMessage: 'Escrow request pending approval.', time: '12:45', initials: 'CT' },
  { id: 'it-support', name: 'IT Support (OIDC)', lastMessage: 'Auth logs look normal.', time: 'Yesterday', initials: 'IT' },
];

const ChatSidebar: React.FC<ChatSidebarProps> = ({ tenantID, isOpen, onSelectChat, activeChatID }) => {
  return (
    <div className={`
      fixed inset-0 z-40 md:relative md:flex w-full md:w-80 bg-white border-r flex-col transition-transform duration-300 ease-in-out
      ${isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
    `}>
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

      <ChatList
        chats={mockChats}
        onSelect={onSelectChat}
        activeID={activeChatID}
      />

      <div className="p-4 border-t bg-gray-50 flex items-center gap-2 text-xs text-gray-500">
        <ShieldCheck className="w-4 h-4 text-green-500" />
        <span>Siloed Isolation: {tenantID}</span>
      </div>
    </div>
  );
};

export default ChatSidebar;
