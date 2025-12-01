import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  ArrowLeft,
  Zap,
  Power,
  Users,
  UserPlus,
  QrCode,
  Download,
  Copy,
  Check,
  Settings,
  Activity,
  Trash2,
  RefreshCw,
  Key,
  Globe
} from 'lucide-react';

interface WireGuardClient {
  id: string;
  name: string;
  publicKey: string;
  allowedIPs: string;
  endpoint: string;
  latestHandshake?: string;
  bytesReceived: number;
  bytesSent: number;
  enabled: boolean;
}

interface WireGuardPanelProps {
  onBack: () => void;
}

const WireGuardPanel: React.FC<WireGuardPanelProps> = ({ onBack }) => {
  const [serverEnabled, setServerEnabled] = useState(true);
  const [serverRunning, setServerRunning] = useState(true);
  const [showAddClient, setShowAddClient] = useState(false);
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const [serverConfig] = useState({
    listenPort: 51820,
    publicKey: 'xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Ds=',
    endpoint: '192.168.1.100:51820',
    subnet: '10.8.0.0/24',
    dns: '8.8.8.8, 8.8.4.4'
  });

  const [clients, setClients] = useState<WireGuardClient[]>([
    {
      id: '1',
      name: 'Johns Laptop',
      publicKey: 'gN65BkIKy1eCE9pP1wdc8ROUtkHLF2PfAqYdyYBz6EA=',
      allowedIPs: '10.8.0.2/32',
      endpoint: '203.0.113.45:51820',
      latestHandshake: '2 minutes ago',
      bytesReceived: 1245678,
      bytesSent: 892345,
      enabled: true
    },
    {
      id: '2',
      name: 'Mobile Phone',
      publicKey: 'aM12CdJFx8wBC7nO9vab4QNRstGID1OeZpXcxYAy3DB=',
      allowedIPs: '10.8.0.3/32',
      endpoint: '198.51.100.22:51820',
      latestHandshake: '5 minutes ago',
      bytesReceived: 456789,
      bytesSent: 234567,
      enabled: true
    },
    {
      id: '3',
      name: 'Work Desktop',
      publicKey: 'zK78DeLGy4zHC2mN5tcd9PORqtEIF8NdYoWbzXBx7EA=',
      allowedIPs: '10.8.0.4/32',
      endpoint: '',
      latestHandshake: undefined,
      bytesReceived: 0,
      bytesSent: 0,
      enabled: false
    }
  ]);

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const handleCopy = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
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
            <div className="absolute inset-0 bg-gradient-to-r from-cyan-500 to-blue-600 rounded-xl blur-lg opacity-50" />
            <div className="relative bg-gradient-to-br from-cyan-500 to-blue-600 p-4 rounded-xl">
              <Zap className="w-8 h-8 text-white" />
            </div>
          </div>
          <div className="flex-1">
            <h2 className="text-3xl font-bold text-white">WireGuard VPN</h2>
            <p className="text-purple-200 mt-1">Modern, fast, and secure VPN protocol</p>
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
          <div className="absolute inset-0 bg-gradient-to-br from-cyan-500/20 to-blue-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Activity className="w-8 h-8 text-cyan-400 mb-2" />
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
            <p className="text-2xl font-bold text-white">
              {clients.filter(c => c.latestHandshake).length} / {clients.length}
            </p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/20 to-pink-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Globe className="w-8 h-8 text-purple-400 mb-2" />
            <p className="text-sm text-purple-200 mb-1">Listen Port</p>
            <p className="text-2xl font-bold text-white">{serverConfig.listenPort}</p>
          </div>
        </div>

        <div className="relative overflow-hidden rounded-2xl">
          <div className="absolute inset-0 bg-gradient-to-br from-green-500/20 to-emerald-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <Key className="w-8 h-8 text-green-400 mb-2" />
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
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <Settings className="w-6 h-6 text-cyan-400" />
              <h3 className="text-xl font-bold text-white">Server Configuration</h3>
            </div>
            <button className="text-purple-200 hover:text-white transition-colors">
              <RefreshCw className="w-5 h-5" />
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm text-purple-200">Public Endpoint</label>
              <div className="flex items-center gap-2 p-3 bg-white/5 rounded-xl border border-white/10">
                <input
                  type="text"
                  value={serverConfig.endpoint}
                  readOnly
                  className="flex-1 bg-transparent text-white outline-none"
                />
                <button
                  onClick={() => handleCopy(serverConfig.endpoint, 'endpoint')}
                  className="p-2 hover:bg-white/10 rounded-lg transition-colors"
                >
                  {copiedId === 'endpoint' ? (
                    <Check className="w-4 h-4 text-green-400" />
                  ) : (
                    <Copy className="w-4 h-4 text-purple-300" />
                  )}
                </button>
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-sm text-purple-200">DNS Servers</label>
              <input
                type="text"
                value={serverConfig.dns}
                className="w-full p-3 bg-white/5 border border-white/10 rounded-xl text-white outline-none focus:border-cyan-500 transition-colors"
              />
            </div>
          </div>
        </div>
      </motion.div>

      {/* Clients Section */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-white/5 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <Users className="w-6 h-6 text-cyan-400" />
              <h3 className="text-xl font-bold text-white">Clients</h3>
            </div>
            <motion.button
              onClick={() => setShowAddClient(true)}
              className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-cyan-500 to-blue-600 text-white rounded-xl font-medium"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <UserPlus className="w-4 h-4" />
              Add Client
            </motion.button>
          </div>

          <div className="space-y-4">
            {clients.map((client, index) => (
              <motion.div
                key={client.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.05 * index }}
                className="p-4 bg-white/5 rounded-xl border border-white/10 hover:border-white/20 transition-all"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <h4 className="text-lg font-semibold text-white mb-1">{client.name}</h4>
                    <div className="flex items-center gap-4 text-sm text-purple-200">
                      <span>{client.allowedIPs}</span>
                      {client.latestHandshake && (
                        <span className="flex items-center gap-1">
                          <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                          Last seen {client.latestHandshake}
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <button
                      className="p-2 text-cyan-400 hover:bg-cyan-500/20 rounded-lg transition-colors"
                      title="Show QR Code"
                    >
                      <QrCode className="w-5 h-5" />
                    </button>
                    <button
                      className="p-2 text-blue-400 hover:bg-blue-500/20 rounded-lg transition-colors"
                      title="Download Config"
                    >
                      <Download className="w-5 h-5" />
                    </button>
                    <button
                      className="p-2 text-red-400 hover:bg-red-500/20 rounded-lg transition-colors"
                      title="Delete Client"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
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
        </div>
      </motion.div>
    </div>
  );
};

export default WireGuardPanel;
