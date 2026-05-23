import React from 'react';
import { Menu, Search, ShieldCheck } from 'lucide-react';

interface ChatSidebarProps {
  tenantID: string;
}

const ChatSidebar: React.FC<ChatSidebarProps> = ({ tenantID }) => {
  return (
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
  );
};

export default ChatSidebar;
