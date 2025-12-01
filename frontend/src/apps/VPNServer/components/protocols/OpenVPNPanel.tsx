import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  ArrowLeft,
  Shield,
  Power,
  Users,
  Download,
  Settings,
  Activity,
  Lock,
  Globe,
  Server,
  Key,
  FileText
} from 'lucide-react';

interface OpenVPNClient {
  id: string;
  commonName: string;
  realAddress: string;
  virtualAddress: string;
  connectedSince: string;
  bytesReceived: number;
  bytesSent: number;
}

interface OpenVPNPanelProps {
  onBack: () => void;
}

const OpenVPNPanel: React.FC<OpenVPNPanelProps> = ({ onBack }) => {
  const [serverEnabled, setServerEnabled] = useState(true);
  const [serverRunning, setServerRunning] = useState(false);

  const [serverConfig] = useState({
    protocol: 'UDP',
    port: 1194,
    subnet: '10.9.0.0/24',
    cipher: 'AES-256-GCM',
    auth: 'SHA512',
    compression: 'lz4-v2'
  });

  const [clients, setClients] = useState<OpenVPNClient[]>([
    {
      id: '1',
      commonName: 'client1',
      realAddress: '203.0.113.45:51234',
      virtualAddress: '10.9.0.6',
      connectedSince: '2024-11-29 10:30:15',
      bytesReceived: 5432876,
      bytesSent: 3219876
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
            <div className="absolute inset-0 bg-gradient-to-r from-green-500 to-emerald-600 rounded-xl blur-lg opacity-50" />
            <div className="relative bg-gradient-to-br from-green-500 to-emerald-600 p-4 rounded-xl">
              <Shield className="w-8 h-8 text-white" />
            </div>
          </div>
          <div className="flex-1">
            <h2 className="text-3xl font-bold text-white">OpenVPN</h2>
            <p className="text-purple-200 mt-1">Industry standard with wide compatibility</p>
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
          <div className="absolute inset-0 bg-gradient-to-br from-green-500/20 to-emerald-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Activity className="w-8 h-8 text-green-400 mb-2" />
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
            <Server className="w-8 h-8 text-purple-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Protocol</p>
            <p className="text-2xl font-bold text-white">{serverConfig.protocol}</p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-cyan-500/20 to-blue-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Globe className="w-8 h-8 text-cyan-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Port</p>
            <p className="text-2xl font-bold text-white">{serverConfig.port}</p>
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
            <Settings className="w-6 h-6 text-green-400" />
            <h3 className="text-xl font-bold text-white">Server Configuration</h3>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-purple-200 mb-2">Protocol</label>
                <select
                  value={serverConfig.protocol}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-green-500 transition-colors"
                >
                  <option value="UDP" className="bg-slate-800">UDP (Recommended)</option>
                  <option value="TCP" className="bg-slate-800">TCP</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Listen Port</label>
                <input
                  type="number"
                  value={serverConfig.port}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-green-500 transition-colors"
                />
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Virtual Network</label>
                <input
                  type="text"
                  value={serverConfig.subnet}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-green-500 transition-colors"
                />
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-purple-200 mb-2">Cipher</label>
                <select
                  value={serverConfig.cipher}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-green-500 transition-colors"
                >
                  <option value="AES-256-GCM" className="bg-slate-800">AES-256-GCM (Recommended)</option>
                  <option value="AES-128-GCM" className="bg-slate-800">AES-128-GCM</option>
                  <option value="AES-256-CBC" className="bg-slate-800">AES-256-CBC</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Authentication</label>
                <select
                  value={serverConfig.auth}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-green-500 transition-colors"
                >
                  <option value="SHA512" className="bg-slate-800">SHA512</option>
                  <option value="SHA256" className="bg-slate-800">SHA256</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-purple-200 mb-2">Compression</label>
                <select
                  value={serverConfig.compression}
                  className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-green-500 transition-colors"
                >
                  <option value="lz4-v2" className="bg-slate-800">LZ4-v2</option>
                  <option value="lzo" className="bg-slate-800">LZO</option>
                  <option value="none" className="bg-slate-800">None</option>
                </select>
              </div>
            </div>
          </div>

          <div className="mt-6 flex gap-4">
            <motion.button
              className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-xl font-medium"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Download className="w-4 h-4" />
              Download Server Config
            </motion.button>
            <motion.button
              className="flex items-center gap-2 px-4 py-2 bg-white/5 border border-white/10 text-white rounded-xl font-medium hover:bg-white/10 transition-colors"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Key className="w-4 h-4" />
              Regenerate Certificates
            </motion.button>
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
            <Users className="w-6 h-6 text-green-400" />
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
                      <h4 className="text-lg font-semibold text-white mb-1">{client.commonName}</h4>
                      <div className="flex items-center gap-4 text-sm text-purple-200">
                        <span>Real: {client.realAddress}</span>
                        <span>Virtual: {client.virtualAddress}</span>
                      </div>
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

      {/* Certificate Management */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <Lock className="w-6 h-6 text-green-400" />
            <h3 className="text-xl font-bold text-white">Certificate Management</h3>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <motion.button
              className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-left"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <FileText className="w-6 h-6 text-green-400 mb-2" />
              <p className="text-white font-medium mb-1">CA Certificate</p>
              <p className="text-sm text-purple-200">View & download</p>
            </motion.button>

            <motion.button
              className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-left"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Key className="w-6 h-6 text-blue-400 mb-2" />
              <p className="text-white font-medium mb-1">Generate Client</p>
              <p className="text-sm text-purple-200">Create new certificate</p>
            </motion.button>

            <motion.button
              className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-left"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <FileText className="w-6 h-6 text-purple-400 mb-2" />
              <p className="text-white font-medium mb-1">Revoke Certificate</p>
              <p className="text-sm text-purple-200">Manage revocations</p>
            </motion.button>
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default OpenVPNPanel;
