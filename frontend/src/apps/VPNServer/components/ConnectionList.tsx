import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  Users,
  UserPlus,
  Search,
  Check,
  X,
  Edit,
  Trash2,
  Shield,
  Key,
  Mail,
  Calendar,
  Activity
} from 'lucide-react';

interface VPNUser {
  id: string;
  username: string;
  email: string;
  createdAt: string;
  lastConnection?: string;
  protocols: {
    wireguard: boolean;
    openvpn: boolean;
    pptp: boolean;
    l2tp: boolean;
  };
  enabled: boolean;
}

const ConnectionList: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [users, setUsers] = useState<VPNUser[]>([
    {
      id: '1',
      username: 'john.doe',
      email: 'john@example.com',
      createdAt: '2024-01-15',
      lastConnection: '2024-11-28 14:32',
      protocols: { wireguard: true, openvpn: true, pptp: false, l2tp: true },
      enabled: true
    },
    {
      id: '2',
      username: 'jane.smith',
      email: 'jane@example.com',
      createdAt: '2024-02-20',
      lastConnection: '2024-11-27 09:15',
      protocols: { wireguard: true, openvpn: false, pptp: false, l2tp: false },
      enabled: true
    },
    {
      id: '3',
      username: 'mike.johnson',
      email: 'mike@example.com',
      createdAt: '2024-03-10',
      protocols: { wireguard: false, openvpn: true, pptp: true, l2tp: true },
      enabled: false
    },
    {
      id: '4',
      username: 'sarah.wilson',
      email: 'sarah@example.com',
      createdAt: '2024-04-05',
      lastConnection: '2024-11-29 16:45',
      protocols: { wireguard: true, openvpn: true, pptp: false, l2tp: true },
      enabled: true
    }
  ]);

  const [showAddUser, setShowAddUser] = useState(false);

  const protocols = [
    { id: 'wireguard', name: 'WireGuard', color: 'cyan' },
    { id: 'openvpn', name: 'OpenVPN', color: 'green' },
    { id: 'pptp', name: 'PPTP', color: 'orange' },
    { id: 'l2tp', name: 'L2TP/IPsec', color: 'purple' }
  ];

  const toggleProtocol = (userId: string, protocol: keyof VPNUser['protocols']) => {
    setUsers(users.map(user =>
      user.id === userId
        ? { ...user, protocols: { ...user.protocols, [protocol]: !user.protocols[protocol] } }
        : user
    ));
  };

  const toggleUserEnabled = (userId: string) => {
    setUsers(users.map(user =>
      user.id === userId ? { ...user, enabled: !user.enabled } : user
    ));
  };

  const filteredUsers = users.filter(user =>
    user.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
    user.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Header Actions */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex items-center justify-between gap-4"
      >
        {/* Search */}
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-purple-300" />
          <input
            type="text"
            placeholder="Search users..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-12 pr-4 py-3 bg-white/5 backdrop-blur-xl border border-white/10 rounded-xl text-white placeholder-purple-300 focus:outline-none focus:border-purple-500 transition-colors"
          />
        </div>

        {/* Add User Button */}
        <motion.button
          onClick={() => setShowAddUser(true)}
          className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl font-medium shadow-lg hover:shadow-xl transition-shadow"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <UserPlus className="w-5 h-5" />
          Add User
        </motion.button>
      </motion.div>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="relative overflow-hidden rounded-2xl"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-blue-500/20 to-purple-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-purple-200 mb-1">Total Users</p>
                <p className="text-3xl font-bold text-white">{users.length}</p>
              </div>
              <Users className="w-8 h-8 text-blue-400" />
            </div>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="relative overflow-hidden rounded-2xl"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-green-500/20 to-emerald-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-purple-200 mb-1">Active Users</p>
                <p className="text-3xl font-bold text-white">{users.filter(u => u.enabled).length}</p>
              </div>
              <Activity className="w-8 h-8 text-green-400" />
            </div>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="relative overflow-hidden rounded-2xl"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/20 to-pink-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-purple-200 mb-1">Connected Now</p>
                <p className="text-3xl font-bold text-white">{users.filter(u => u.lastConnection).length}</p>
              </div>
              <Shield className="w-8 h-8 text-purple-400" />
            </div>
          </div>
        </motion.div>
      </div>

      {/* Users Table */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative border border-white/10">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/10">
                  <th className="px-6 py-4 text-left text-sm font-semibold text-purple-200">User</th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-purple-200">Email</th>
                  <th className="px-6 py-4 text-center text-sm font-semibold text-purple-200">Status</th>
                  {protocols.map(protocol => (
                    <th key={protocol.id} className="px-6 py-4 text-center text-sm font-semibold text-purple-200">
                      {protocol.name}
                    </th>
                  ))}
                  <th className="px-6 py-4 text-center text-sm font-semibold text-purple-200">Last Connection</th>
                  <th className="px-6 py-4 text-right text-sm font-semibold text-purple-200">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredUsers.map((user, index) => (
                  <motion.tr
                    key={user.id}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 0.05 * index }}
                    className="border-b border-white/5 hover:bg-white/5 transition-colors"
                  >
                    {/* User Info */}
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                          {user.username.charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <p className="text-white font-medium">{user.username}</p>
                          <p className="text-xs text-purple-300">Created {user.createdAt}</p>
                        </div>
                      </div>
                    </td>

                    {/* Email */}
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2 text-purple-200">
                        <Mail className="w-4 h-4" />
                        <span className="text-sm">{user.email}</span>
                      </div>
                    </td>

                    {/* Status */}
                    <td className="px-6 py-4">
                      <div className="flex justify-center">
                        <button
                          onClick={() => toggleUserEnabled(user.id)}
                          className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
                            user.enabled
                              ? 'bg-green-500/20 text-green-400 hover:bg-green-500/30'
                              : 'bg-gray-500/20 text-gray-400 hover:bg-gray-500/30'
                          }`}
                        >
                          {user.enabled ? 'Active' : 'Disabled'}
                        </button>
                      </div>
                    </td>

                    {/* Protocol Access Checkboxes */}
                    {protocols.map(protocol => (
                      <td key={protocol.id} className="px-6 py-4">
                        <div className="flex justify-center">
                          <button
                            onClick={() => toggleProtocol(user.id, protocol.id as keyof VPNUser['protocols'])}
                            disabled={!user.enabled}
                            className={`w-8 h-8 rounded-lg flex items-center justify-center transition-all ${
                              user.protocols[protocol.id as keyof VPNUser['protocols']]
                                ? `bg-${protocol.color}-500/20 text-${protocol.color}-400 border border-${protocol.color}-400/50 hover:bg-${protocol.color}-500/30`
                                : 'bg-white/5 text-gray-500 border border-white/10 hover:bg-white/10'
                            } ${!user.enabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
                          >
                            {user.protocols[protocol.id as keyof VPNUser['protocols']] && (
                              <Check className="w-5 h-5" />
                            )}
                          </button>
                        </div>
                      </td>
                    ))}

                    {/* Last Connection */}
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-center gap-2 text-purple-200 text-sm">
                        {user.lastConnection ? (
                          <>
                            <Calendar className="w-4 h-4" />
                            {user.lastConnection}
                          </>
                        ) : (
                          <span className="text-gray-500">Never</span>
                        )}
                      </div>
                    </td>

                    {/* Actions */}
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-2">
                        <motion.button
                          whileHover={{ scale: 1.1 }}
                          whileTap={{ scale: 0.9 }}
                          className="p-2 text-blue-400 hover:bg-blue-500/20 rounded-lg transition-colors"
                        >
                          <Edit className="w-4 h-4" />
                        </motion.button>
                        <motion.button
                          whileHover={{ scale: 1.1 }}
                          whileTap={{ scale: 0.9 }}
                          className="p-2 text-red-400 hover:bg-red-500/20 rounded-lg transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </motion.button>
                      </div>
                    </td>
                  </motion.tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default ConnectionList;
