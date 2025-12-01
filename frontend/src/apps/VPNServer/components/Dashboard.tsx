import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import {
  Shield,
  Activity,
  Users,
  ArrowRight,
  Wifi,
  Lock,
  Zap,
  Server,
  Circle,
  TrendingUp,
  Download,
  Upload
} from 'lucide-react';

interface VPNProtocol {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  gradient: string;
  glowColor: string;
  enabled: boolean;
  running: boolean;
  connections: number;
  maxConnections: number;
  ipRange: string;
  bytesIn: number;
  bytesOut: number;
}

interface DashboardProps {
  onProtocolClick: (protocolId: string) => void;
}

const Dashboard: React.FC<DashboardProps> = ({ onProtocolClick }) => {
  const [protocols, setProtocols] = useState<VPNProtocol[]>([
    {
      id: 'wireguard',
      name: 'WireGuard',
      description: 'Modern, fast, and secure VPN protocol',
      icon: <Zap className="w-6 h-6" />,
      gradient: 'from-cyan-500 to-blue-600',
      glowColor: 'cyan',
      enabled: true,
      running: true,
      connections: 3,
      maxConnections: 50,
      ipRange: '10.8.0.0/24',
      bytesIn: 1245000000,
      bytesOut: 892000000
    },
    {
      id: 'openvpn',
      name: 'OpenVPN',
      description: 'Industry standard with wide compatibility',
      icon: <Shield className="w-6 h-6" />,
      gradient: 'from-green-500 to-emerald-600',
      glowColor: 'green',
      enabled: true,
      running: false,
      connections: 0,
      maxConnections: 100,
      ipRange: '10.9.0.0/24',
      bytesIn: 0,
      bytesOut: 0
    },
    {
      id: 'pptp',
      name: 'PPTP',
      description: 'Legacy protocol for older devices',
      icon: <Server className="w-6 h-6" />,
      gradient: 'from-orange-500 to-red-600',
      glowColor: 'orange',
      enabled: false,
      running: false,
      connections: 0,
      maxConnections: 20,
      ipRange: '10.10.0.0/24',
      bytesIn: 0,
      bytesOut: 0
    },
    {
      id: 'l2tp',
      name: 'L2TP/IPsec',
      description: 'Enterprise-grade security',
      icon: <Lock className="w-6 h-6" />,
      gradient: 'from-purple-500 to-pink-600',
      glowColor: 'purple',
      enabled: true,
      running: true,
      connections: 1,
      maxConnections: 30,
      ipRange: '10.11.0.0/24',
      bytesIn: 456000000,
      bytesOut: 234000000
    }
  ]);

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const totalConnections = protocols.reduce((sum, p) => sum + p.connections, 0);
  const totalBytesIn = protocols.reduce((sum, p) => sum + p.bytesIn, 0);
  const totalBytesOut = protocols.reduce((sum, p) => sum + p.bytesOut, 0);
  const activeProtocols = protocols.filter(p => p.running).length;

  return (
    <div className="space-y-6">
      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="relative overflow-hidden rounded-2xl"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-blue-500/20 to-purple-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <div className="flex items-center justify-between mb-2">
              <Activity className="w-8 h-8 text-blue-400" />
              <span className="text-3xl font-bold text-white">{activeProtocols}</span>
            </div>
            <p className="text-sm text-purple-200">Active Protocols</p>
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
            <div className="flex items-center justify-between mb-2">
              <Users className="w-8 h-8 text-green-400" />
              <span className="text-3xl font-bold text-white">{totalConnections}</span>
            </div>
            <p className="text-sm text-purple-200">Total Connections</p>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="relative overflow-hidden rounded-2xl"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-cyan-500/20 to-blue-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <div className="flex items-center justify-between mb-2">
              <Download className="w-8 h-8 text-cyan-400" />
              <span className="text-2xl font-bold text-white">{formatBytes(totalBytesIn)}</span>
            </div>
            <p className="text-sm text-purple-200">Data Received</p>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="relative overflow-hidden rounded-2xl"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/20 to-pink-600/20 backdrop-blur-xl" />
          <div className="relative p-6 border border-white/10">
            <div className="flex items-center justify-between mb-2">
              <Upload className="w-8 h-8 text-purple-400" />
              <span className="text-2xl font-bold text-white">{formatBytes(totalBytesOut)}</span>
            </div>
            <p className="text-sm text-purple-200">Data Sent</p>
          </div>
        </motion.div>
      </div>

      {/* Protocol Cards */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {protocols.map((protocol, index) => (
          <motion.div
            key={protocol.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 * index }}
            className="relative group cursor-pointer"
            onClick={() => onProtocolClick(protocol.id)}
          >
            {/* Glow Effect */}
            <div className={`absolute inset-0 bg-gradient-to-r ${protocol.gradient} rounded-2xl blur-xl opacity-0 group-hover:opacity-30 transition-opacity duration-300`} />

            {/* Card Content */}
            <div className="relative bg-white/5 backdrop-blur-xl rounded-2xl p-6 border border-white/10 hover:border-white/20 transition-all duration-300">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-4">
                  {/* Icon Badge */}
                  <div className={`relative`}>
                    <div className={`absolute inset-0 bg-gradient-to-r ${protocol.gradient} rounded-xl blur opacity-50`} />
                    <div className={`relative bg-gradient-to-br ${protocol.gradient} p-3 rounded-xl text-white`}>
                      {protocol.icon}
                    </div>
                  </div>

                  <div>
                    <h3 className="text-xl font-bold text-white mb-1">{protocol.name}</h3>
                    <p className="text-sm text-purple-200">{protocol.description}</p>
                  </div>
                </div>

                <ArrowRight className="w-5 h-5 text-purple-300 group-hover:text-white group-hover:translate-x-1 transition-all" />
              </div>

              {/* Status Indicators */}
              <div className="flex items-center gap-4 mb-4">
                <div className="flex items-center gap-2">
                  <Circle
                    className={`w-3 h-3 ${
                      protocol.running
                        ? 'fill-green-400 text-green-400'
                        : protocol.enabled
                          ? 'fill-yellow-400 text-yellow-400'
                          : 'fill-gray-400 text-gray-400'
                    }`}
                  />
                  <span className={`text-sm font-medium ${
                    protocol.running
                      ? 'text-green-400'
                      : protocol.enabled
                        ? 'text-yellow-400'
                        : 'text-gray-400'
                  }`}>
                    {protocol.running ? 'Running' : protocol.enabled ? 'Stopped' : 'Disabled'}
                  </span>
                </div>

                <div className="flex items-center gap-2 text-purple-200">
                  <Wifi className="w-4 h-4" />
                  <span className="text-sm">{protocol.ipRange}</span>
                </div>
              </div>

              {/* Connections Bar */}
              <div className="mb-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-purple-200">Connections</span>
                  <span className="text-sm font-medium text-white">
                    {protocol.connections} / {protocol.maxConnections}
                  </span>
                </div>
                <div className="h-2 bg-white/5 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${(protocol.connections / protocol.maxConnections) * 100}%` }}
                    transition={{ duration: 1, delay: 0.5 + index * 0.1 }}
                    className={`h-full bg-gradient-to-r ${protocol.gradient}`}
                  />
                </div>
              </div>

              {/* Traffic Stats */}
              {protocol.running && (
                <div className="grid grid-cols-2 gap-4 pt-4 border-t border-white/10">
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <Download className="w-4 h-4 text-cyan-400" />
                      <span className="text-xs text-purple-200">Received</span>
                    </div>
                    <p className="text-sm font-medium text-white">{formatBytes(protocol.bytesIn)}</p>
                  </div>
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <Upload className="w-4 h-4 text-purple-400" />
                      <span className="text-xs text-purple-200">Sent</span>
                    </div>
                    <p className="text-sm font-medium text-white">{formatBytes(protocol.bytesOut)}</p>
                  </div>
                </div>
              )}
            </div>
          </motion.div>
        ))}
      </div>
    </div>
  );
};

export default Dashboard;
