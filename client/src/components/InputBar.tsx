import React from 'react';
import { Smile, Paperclip, Send } from 'lucide-react';

interface InputBarProps {
  value: string;
  onChange: (val: string) => void;
  onSend: () => void;
  disabled: boolean;
}

const InputBar: React.FC<InputBarProps> = ({ value, onChange, onSend, disabled }) => {
  return (
    <div className="p-4 bg-white md:bg-transparent">
      <div className="max-w-3xl mx-auto flex items-end gap-2">
        <div className="flex-1 bg-white rounded-xl shadow-sm flex items-end p-2 border border-gray-200">
          <Smile className="text-gray-400 m-2 cursor-pointer" />
          <textarea
            rows={1}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && !e.shiftKey && (e.preventDefault(), onSend())}
            placeholder="Write a message..."
            className="flex-1 py-2 px-1 text-sm focus:outline-none resize-none max-h-32"
          />
          <Paperclip className="text-gray-400 m-2 cursor-pointer rotate-45" />
        </div>
        <button
          onClick={onSend}
          disabled={disabled || !value.trim()}
          className={`w-12 h-12 rounded-full flex items-center justify-center text-white shadow-md transition-all flex-shrink-0 ${
            !disabled && value.trim() ? 'bg-telegram-blue hover:scale-105' : 'bg-gray-300 cursor-not-allowed'
          }`}
        >
          <Send className="w-5 h-5 ml-0.5" />
        </button>
      </div>
    </div>
  );
};

export default InputBar;
