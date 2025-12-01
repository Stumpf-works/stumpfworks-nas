import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  ArrowLeft,
  Lock,
  Power,
  Users,
  Activity,
  Globe,
  Settings,
  Key,
  Shield,
  Info,
  Server
} from 'lucide-react';

interface L2TPClient {
  id: string;
  username: string;
  ipAddress: string;
  connectedSince: string;
  bytesReceived: number;
  bytesSent: number;
}

interface L2TPPanelProps {
  onBack: () => void;
}

const L2TPPanel: React.FC<L2TPPanelProps> = ({ onBack }) => {
  const [serverEnabled, setServerEnabled] = useState(true);
  const [serverRunning, setServerRunning] = useState(true);

  const [serverConfig] = useState({
    port: 1701,
    ipsecPort: 500,
    subnet: '10.11.0.0/24',
    psk: '••••••••••••••••',
    encryption: 'AES-256',
    authentication: 'SHA2-256'
  });

  const [clients, setClients] = useState<L2TPClient[]>([
    {
      id: '1',
      username: 'corporate-user',
      ipAddress: '10.11.0.5',
      connectedSince: '2024-11-29 08:15:42',
      bytesReceived: 2345678,
      bytesSent: 1234567
    }
  ]);

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const toggleServer = () => {
    setServerRunning(!serverRunning);
  };

  return (
    <div className="space-y-6">
      {/* Back Button & Header */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
      >
        <button
          onClick={onBack}
          className="flex items-center gap-2 text-purple-200 hover:text-white mb-6 transition-colors"
        >
          <ArrowLeft className="w-5 h-5" />
          Back to Dashboard
        </button>

        <div className="flex items-center gap-4">
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-purple-500 to-pink-600 rounded-xl blur-lg opacity-50" />
            <div className="relative bg-gradient-to-br from-purple-500 to-pink-600 p-4 rounded-xl">
              <Lock className="w-8 h-8 text-white" />
            </div>
          </div>
          <div className="flex-1">
            <h2 className="text-3xl font-bold text-white">L2TP/IPsec VPN</h2>
            <p className="text-purple-200 mt-1">Enterprise-grade security with IPsec encryption</p>
          </div>
          <motion.button
            onClick={toggleServer}
            className={`flex items-center gap-2 px-6 py-3 rounded-xl font-medium shadow-lg transition-all ${
              serverRunning
                ? 'bg-gradient-to-r from-green-500 to-emerald-600 text-white'
                : 'bg-gradient-to-r from-gray-500 to-gray-600 text-white'
            }`}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <Power className="w-5 h-5" />
            {serverRunning ? 'Stop Server' : 'Start Server'}
          </motion.button>
        </div>
      </motion.div>

      {/* Server Status */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="grid grid-cols-1 md:grid-cols-4 gap-4"
      >
        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/20 to-pink-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Activity className="w-8 h-8 text-purple-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Status</p>
            <p className={`text-2xl font-bold ${serverRunning ? 'text-green-400' : 'text-gray-400'}`}>
              {serverRunning ? 'Running' : 'Stopped'}
            </p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-blue-500/20 to-purple-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Users className="w-8 h-8 text-blue-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Connected</p>
            <p className="text-2xl font-bold text-white">{clients.length}</p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-pink-500/20 to-rose-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Globe className="w-8 h-8 text-pink-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">L2TP Port</p>
            <p className="text-2xl font-bold text-white">{serverConfig.port}</p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-cyan-500/20 to-blue-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Shield className="w-8 h-8 text-cyan-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">IPsec Port</p>
            <p className="text-2xl font-bold text-white">{serverConfig.ipsecPort}</p>
          </div>
        </div>
      </motion.div>

      {/* IPsec Configuration */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl">
              <Shield className="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 className="text-xl font-bold text-white">IPsec Configuration</h3>
              <p className="text-sm text-purple-200">Configure IPsec security parameters</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-purple-200 mb-2">Pre-Shared Key (PSK)</label>
                <div className="relative">
                  <input
                    type="password"
                    value={serverConfig.psk}
                    className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-purple-500 transition-colors pr-12"
                  />
                  <button className="absolute right-3 top-1/2 -translate-y-1/2 text-purple-300 hover:text-white transition-colors">
                    <Key className="w-5 h-5" />
                  </button>
                </div>
                <p className="text-xs text-purple-300 mt-1">Shared secret for IPsec authentication</p>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Encryption Algorithm</label>
                <select
                  value={serverConfig.encryption}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-purple-500 transition-colors"
                >
                  <option value="AES-256" className="bg-slate-800">AES-256 (Recommended)</option>
                  <option value="AES-192" className="bg-slate-800">AES-192</option>
                  <option value="AES-128" className="bg-slate-800">AES-128</option>
                  <option value="3DES" className="bg-slate-800">3DES (Legacy)</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Authentication Algorithm</label>
                <select
                  value={serverConfig.authentication}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-purple-500 transition-colors"
                >
                  <option value="SHA2-256" className="bg-slate-800">SHA2-256</option>
                  <option value="SHA2-512" className="bg-slate-800">SHA2-512</option>
                  <option value="SHA1" className="bg-slate-800">SHA1 (Legacy)</option>
                </select>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-purple-200 mb-2">L2TP Port</label>
                <input
                  type="number"
                  value={serverConfig.port}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-purple-500 transition-colors"
                />
                <p className="text-xs text-purple-300 mt-1">Standard L2TP port is 1701</p>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">IPsec Port (UDP)</label>
                <input
                  type="number"
                  value={serverConfig.ipsecPort}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-purple-500 transition-colors"
                />
                <p className="text-xs text-purple-300 mt-1">Standard IPsec IKE port is 500</p>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Virtual Network</label>
                <input
                  type="text"
                  value={serverConfig.subnet}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-purple-500 transition-colors"
                />
                <p className="text-xs text-purple-300 mt-1">IP range for L2TP clients</p>
              </div>
            </div>
          </div>

          <div className="mt-6 p-4 bg-purple-500/10 border border-purple-500/20 rounded-xl">
            <div className="flex items-start gap-3">
              <Info className="w-5 h-5 text-purple-400 mt-0.5" />
              <div>
                <p className="text-purple-400 font-medium mb-1">NAT Traversal</p>
                <p className="text-sm text-purple-200">
                  L2TP/IPsec supports NAT-T (NAT Traversal) on UDP port 4500, allowing connections through NAT devices. This is automatically enabled for maximum compatibility.
                </p>
              </div>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Connected Clients */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <Users className="w-6 h-6 text-purple-400" />
            <h3 className="text-xl font-bold text-white">Connected Clients</h3>
          </div>

          {clients.length === 0 ? (
            <div className="text-center py-12 text-purple-200">
              <Users className="w-12 h-12 mx-auto mb-4 opacity-50" />
              <p>No clients currently connected</p>
            </div>
          ) : (
            <div className="space-y-4">
              {clients.map((client, index) => (
                <motion.div
                  key={client.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.05 * index }}
                  className="p-4 bg-white/5 rounded-xl border border-white/10"
                >
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <h4 className="text-lg font-semibold text-white mb-1">{client.username}</h4>
                      <p className="text-sm text-purple-200">Virtual IP: {client.ipAddress}</p>
                    </div>
                    <div className="flex items-center gap-2 text-xs text-purple-300">
                      <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                      Connected since {client.connectedSince}
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="flex items-center justify-between p-2 bg-white/5 rounded-lg">
                      <span className="text-purple-200">Received</span>
                      <span className="text-white font-medium">{formatBytes(client.bytesReceived)}</span>
                    </div>
                    <div className="flex items-center justify-between p-2 bg-white/5 rounded-lg">
                      <span className="text-purple-200">Sent</span>
                      <span className="text-white font-medium">{formatBytes(client.bytesSent)}</span>
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          )}
        </div>
      </motion.div>

      {/* Client Setup Instructions */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <Server className="w-6 h-6 text-purple-400" />
            <h3 className="text-xl font-bold text-white">Client Setup</h3>
          </div>

          <div className="space-y-3 text-sm">
            <div className="p-4 bg-white/5 rounded-xl">
              <p className="text-purple-200 mb-2">
                <span className="font-semibold text-white">Windows:</span> Native support in Windows 7-11
              </p>
              <p className="text-purple-300 text-xs">
                Settings → Network & Internet → VPN → Add VPN → Type: L2TP/IPsec with pre-shared key
              </p>
            </div>

            <div className="p-4 bg-white/5 rounded-xl">
              <p className="text-purple-200 mb-2">
                <span className="font-semibold text-white">macOS:</span> Native support in macOS
              </p>
              <p className="text-purple-300 text-xs">
                System Preferences → Network → + → Interface: VPN → VPN Type: L2TP over IPSec
              </p>
            </div>

            <div className="p-4 bg-white/5 rounded-xl">
              <p className="text-purple-200 mb-2">
                <span className="font-semibold text-white">iOS/iPadOS:</span> Native support
              </p>
              <p className="text-purple-300 text-xs">
                Settings → General → VPN → Add VPN Configuration → Type: L2TP
              </p>
            </div>

            <div className="p-4 bg-white/5 rounded-xl">
              <p className="text-purple-200 mb-2">
                <span className="font-semibold text-white">Android:</span> Native support
              </p>
              <p className="text-purple-300 text-xs">
                Settings → Network & Internet → VPN → + → Type: L2TP/IPsec PSK
              </p>
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default L2TPPanel;
