import React, { useState } from 'react';
import { ShieldAlert } from 'lucide-react';

interface LoginProps {
  onLogin: (tenantID: string, userID: string, role: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const [tenant, setTenant] = useState('tenant-abc');
  const [user, setUser] = useState('user-123');
  const [role, setRole] = useState('user');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await fetch('/api/auth/token', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ tenant_id: tenant, user_id: user, role: role })
      });

      if (response.ok) {
        const { token } = await response.json();
        // Normally we'd store this in localStorage/secure cookie
        // But for this session, we pass it back to the main app state
        onLogin(tenant, user, token);
      } else {
        alert("Login failed");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div className="bg-telegram-blue p-6 text-white text-center">
          <div className="w-16 h-16 bg-white bg-opacity-20 rounded-full flex items-center justify-center mx-auto mb-4">
            <ShieldAlert className="w-8 h-8" />
          </div>
          <h1 className="text-2xl font-bold text-white">ZemenLink Access</h1>
          <p className="text-blue-100 mt-1">Enterprise Kernel Demo</p>
        </div>

        <form onSubmit={handleSubmit} className="p-8 space-y-6">
          <div>
            <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Tenant ID</label>
            <input
              type="text"
              value={tenant}
              onChange={(e) => setTenant(e.target.value)}
              className="w-full border-b-2 border-gray-100 focus:border-telegram-blue py-2 outline-none transition-colors"
              required
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">User ID</label>
            <input
              type="text"
              value={user}
              onChange={(e) => setUser(e.target.value)}
              className="w-full border-b-2 border-gray-100 focus:border-telegram-blue py-2 outline-none transition-colors"
              required
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Role</label>
            <select
              value={role}
              onChange={(e) => setRole(e.target.value)}
              className="w-full border-b-2 border-gray-100 focus:border-telegram-blue py-2 outline-none transition-colors bg-white"
            >
              <option value="user">Standard User</option>
              <option value="admin">Administrator</option>
              <option value="compliance_officer">Compliance Officer</option>
            </select>
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-telegram-blue text-white rounded-lg py-3 font-semibold shadow-lg hover:bg-opacity-90 transition-all transform active:scale-[0.98]"
          >
            {isLoading ? 'Connecting...' : 'Enter Messaging Kernel'}
          </button>
        </form>

        <div className="bg-gray-50 p-4 text-center text-[10px] text-gray-400 border-t">
          Siloed Data Isolation Active • E2EE Escrow Gated
        </div>
      </div>
    </div>
  );
};

export default Login;
