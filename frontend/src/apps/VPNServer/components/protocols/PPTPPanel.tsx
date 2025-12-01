import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  ArrowLeft,
  Server,
  Power,
  Users,
  Activity,
  Globe,
  AlertTriangle,
  Settings,
  Info
} from 'lucide-react';

interface PPTPClient {
  id: string;
  username: string;
  ipAddress: string;
  connectedSince: string;
  bytesReceived: number;
  bytesSent: number;
}

interface PPTPPanelProps {
  onBack: () => void;
}

const PPTPPanel: React.FC<PPTPPanelProps> = ({ onBack }) => {
  const [serverEnabled, setServerEnabled] = useState(false);
  const [serverRunning, setServerRunning] = useState(false);

  const [serverConfig] = useState({
    port: 1723,
    subnet: '10.10.0.0/24',
    encryption: 'MPPE-128',
    authentication: 'MS-CHAPv2'
  });

  const [clients, setClients] = useState<PPTPClient[]>([]);

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
            <div className="absolute inset-0 bg-gradient-to-r from-orange-500 to-red-600 rounded-xl blur-lg opacity-50" />
            <div className="relative bg-gradient-to-br from-orange-500 to-red-600 p-4 rounded-xl">
              <Server className="w-8 h-8 text-white" />
            </div>
          </div>
          <div className="flex-1">
            <h2 className="text-3xl font-bold text-white">PPTP VPN</h2>
            <p className="text-purple-200 mt-1">Legacy protocol for older devices</p>
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

      {/* Security Warning */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-orange-500/20 to-red-600/20 backdrop-blur-xl" />
        <div className="relative p-6 border border-orange-500/30">
          <div className="flex items-start gap-4">
            <div className="p-3 bg-orange-500/20 rounded-xl">
              <AlertTriangle className="w-6 h-6 text-orange-400" />
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-bold text-orange-400 mb-2">Security Notice</h3>
              <p className="text-purple-200 text-sm">
                PPTP is an outdated VPN protocol with known security vulnerabilities. It should only be used for compatibility with older devices that cannot support modern protocols like WireGuard or OpenVPN. Consider upgrading to a more secure protocol whenever possible.
              </p>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Server Status */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="grid grid-cols-1 md:grid-cols-4 gap-4"
      >
        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-orange-500/20 to-red-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Activity className="w-8 h-8 text-orange-400 mb-2" />
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
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/20 to-pink-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Globe className="w-8 h-8 text-purple-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Port</p>
            <p className="text-2xl font-bold text-white">{serverConfig.port}</p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-cyan-500/20 to-blue-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Server className="w-8 h-8 text-cyan-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Subnet</p>
            <p className="text-lg font-bold text-white">{serverConfig.subnet}</p>
          </div>
        </div>
      </motion.div>

      {/* Server Configuration */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <Settings className="w-6 h-6 text-orange-400" />
            <h3 className="text-xl font-bold text-white">Server Configuration</h3>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-purple-200 mb-2">Listen Port</label>
                <input
                  type="number"
                  value={serverConfig.port}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-orange-500 transition-colors"
                />
                <p className="text-xs text-purple-300 mt-1">Standard PPTP port is 1723</p>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Virtual Network</label>
                <input
                  type="text"
                  value={serverConfig.subnet}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-orange-500 transition-colors"
                />
                <p className="text-xs text-purple-300 mt-1">IP range for PPTP clients</p>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-purple-200 mb-2">Encryption</label>
                <select
                  value={serverConfig.encryption}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-orange-500 transition-colors"
                >
                  <option value="MPPE-128" className="bg-slate-800">MPPE-128</option>
                  <option value="MPPE-40" className="bg-slate-800">MPPE-40 (Weak)</option>
                </select>
                <p className="text-xs text-purple-300 mt-1">Microsoft Point-to-Point Encryption</p>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Authentication</label>
                <select
                  value={serverConfig.authentication}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-orange-500 transition-colors"
                >
                  <option value="MS-CHAPv2" className="bg-slate-800">MS-CHAPv2</option>
                  <option value="CHAP" className="bg-slate-800">CHAP</option>
                  <option value="PAP" className="bg-slate-800">PAP (Not Recommended)</option>
                </select>
                <p className="text-xs text-purple-300 mt-1">Challenge Handshake Authentication</p>
              </div>
            </div>
          </div>

          <div className="mt-6 p-4 bg-blue-500/10 border border-blue-500/20 rounded-xl">
            <div className="flex items-start gap-3">
              <Info className="w-5 h-5 text-blue-400 mt-0.5" />
              <div>
                <p className="text-blue-400 font-medium mb-1">Compatibility Notes</p>
                <p className="text-sm text-purple-200">
                  PPTP is natively supported on Windows (7-11), macOS, iOS, and Android without additional software. However, it provides weaker security compared to modern protocols.
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
            <Users className="w-6 h-6 text-orange-400" />
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
                      <p className="text-sm text-purple-200">IP: {client.ipAddress}</p>
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
    </div>
  );
};

export default PPTPPanel;
