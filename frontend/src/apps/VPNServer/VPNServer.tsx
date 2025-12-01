import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Shield,
  Activity,
  Settings,
  Download,
  Play,
  Square,
  Loader,
} from 'lucide-react';
import { vpnApi, ProtocolStatus, WireGuardPeer, VPNProtocol } from '../../api/vpn';
import { toast } from 'react-hot-toast';

interface TabProps {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}

const Tab: React.FC<TabProps> = ({ active, onClick, children }) => (
  <button
    onClick={onClick}
    className={`
      relative px-6 py-3 rounded-xl font-medium transition-all duration-300
      ${
        active
          ? 'text-white'
          : 'text-gray-400 hover:text-gray-200'
      }
    `}
  >
    {active && (
      <motion.div
        layoutId="activeTab"
        className="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl"
        transition={{ type: 'spring', bounce: 0.2, duration: 0.6 }}
      />
    )}
    <span className="relative z-10">{children}</span>
  </button>
);

const VPNServer: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'wireguard' | 'openvpn' | 'settings'>('overview');
  const [protocols, setProtocols] = useState<ProtocolStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [installing, setInstalling] = useState<string | null>(null);

  useEffect(() => {
    loadProtocolStatus();
    const interval = setInterval(loadProtocolStatus, 5000); // Refresh every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const loadProtocolStatus = async () => {
    try {
      const response = await vpnApi.getStatus();
      if (response.success && response.data) {
        setProtocols(response.data);
      }
    } catch (error) {
      console.error('Failed to load VPN status:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleInstall = async (protocol: VPNProtocol) => {
    setInstalling(protocol);
    const toastId = toast.loading(`Installing ${protocol.toUpperCase()}...`);

    try {
      const response = await vpnApi.installProtocol(protocol);
      if (response.success) {
        toast.success(`${protocol.toUpperCase()} installed successfully!`, { id: toastId });
        await loadProtocolStatus();
      } else {
        throw new Error(response.message || 'Installation failed');
      }
    } catch (error: any) {
      toast.error(`Failed to install ${protocol}: ${error.message}`, { id: toastId });
    } finally {
      setInstalling(null);
    }
  };

  const handleEnable = async (protocol: VPNProtocol) => {
    const toastId = toast.loading(`Enabling ${protocol.toUpperCase()}...`);

    try {
      const response = await vpnApi.enableProtocol(protocol);
      if (response.success) {
        toast.success(`${protocol.toUpperCase()} enabled successfully!`, { id: toastId });
        await loadProtocolStatus();
      } else {
        throw new Error(response.message || 'Enable failed');
      }
    } catch (error: any) {
      toast.error(`Failed to enable ${protocol}: ${error.message}`, { id: toastId });
    }
  };

  const handleDisable = async (protocol: VPNProtocol) => {
    const toastId = toast.loading(`Disabling ${protocol.toUpperCase()}...`);

    try {
      const response = await vpnApi.disableProtocol(protocol);
      if (response.success) {
        toast.success(`${protocol.toUpperCase()} disabled successfully!`, { id: toastId });
        await loadProtocolStatus();
      } else {
        throw new Error(response.message || 'Disable failed');
      }
    } catch (error: any) {
      toast.error(`Failed to disable ${protocol}: ${error.message}`, { id: toastId });
    }
  };

  const ProtocolIcon: React.FC<{ protocol: string }> = ({ protocol }) => {
    const iconClass = "w-10 h-10";
    switch (protocol.toLowerCase()) {
      case 'wireguard':
        return <Shield className={`${iconClass} text-blue-400`} />;
      case 'openvpn':
        return <Shield className={`${iconClass} text-green-400`} />;
      case 'pptp':
        return <Shield className={`${iconClass} text-purple-400`} />;
      case 'l2tp':
        return <Shield className={`${iconClass} text-orange-400`} />;
      default:
        return <Shield className={`${iconClass} text-gray-400`} />;
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-8">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-8"
      >
        <div className="flex items-center gap-4 mb-2">
          <div className="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl">
            <Shield className="w-8 h-8 text-white" />
          </div>
          <div>
            <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              VPN Server
            </h1>
            <p className="text-gray-400 mt-1">
              Multi-protocol VPN server management
            </p>
          </div>
        </div>
      </motion.div>

      {/* Tabs */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="flex gap-2 mb-8 bg-gray-800/50 backdrop-blur-xl p-2 rounded-2xl border border-gray-700/50"
      >
        <Tab active={activeTab === 'overview'} onClick={() => setActiveTab('overview')}>
          <Activity className="w-4 h-4 inline mr-2" />
          Overview
        </Tab>
        <Tab active={activeTab === 'wireguard'} onClick={() => setActiveTab('wireguard')}>
          <Shield className="w-4 h-4 inline mr-2" />
          WireGuard
        </Tab>
        <Tab active={activeTab === 'openvpn'} onClick={() => setActiveTab('openvpn')}>
          <Shield className="w-4 h-4 inline mr-2" />
          OpenVPN
        </Tab>
        <Tab active={activeTab === 'settings'} onClick={() => setActiveTab('settings')}>
          <Settings className="w-4 h-4 inline mr-2" />
          Settings
        </Tab>
      </motion.div>

      {/* Content */}
      <AnimatePresence mode="wait">
        {activeTab === 'overview' && (
          <OverviewTab
            key="overview"
            protocols={protocols}
            loading={loading}
            installing={installing}
            onInstall={handleInstall}
            onEnable={handleEnable}
            onDisable={handleDisable}
            ProtocolIcon={ProtocolIcon}
          />
        )}
        {activeTab === 'wireguard' && (
          <WireGuardTab key="wireguard" />
        )}
        {activeTab === 'openvpn' && (
          <OpenVPNTab key="openvpn" />
        )}
        {activeTab === 'settings' && (
          <SettingsTab key="settings" />
        )}
      </AnimatePresence>
    </div>
  );
};

// Overview Tab Component
interface OverviewTabProps {
  protocols: ProtocolStatus[];
  loading: boolean;
  installing: string | null;
  onInstall: (protocol: VPNProtocol) => void;
  onEnable: (protocol: VPNProtocol) => void;
  onDisable: (protocol: VPNProtocol) => void;
  ProtocolIcon: React.FC<{ protocol: string }>;
}

const OverviewTab: React.FC<OverviewTabProps> = ({
  protocols,
  loading,
  installing,
  onInstall,
  onEnable,
  onDisable,
  ProtocolIcon,
}) => {
  if (loading) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="flex items-center justify-center h-64"
      >
        <Loader className="w-8 h-8 animate-spin text-blue-400" />
      </motion.div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="grid grid-cols-1 md:grid-cols-2 gap-6"
    >
      {protocols.map((protocol, index) => (
        <motion.div
          key={protocol.protocol}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: index * 0.1 }}
          className="relative group"
        >
          {/* Gradient background */}
          <div className="absolute inset-0 bg-gradient-to-br from-blue-600/20 to-purple-600/20 rounded-3xl blur-xl group-hover:blur-2xl transition-all duration-500" />

          {/* Card content */}
          <div className="relative bg-gray-800/80 backdrop-blur-xl rounded-3xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all duration-300">
            {/* Header */}
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <ProtocolIcon protocol={protocol.protocol} />
                <div>
                  <h3 className="text-xl font-bold text-white capitalize">
                    {protocol.protocol}
                  </h3>
                  <div className="flex items-center gap-2 mt-1">
                    {protocol.running ? (
                      <span className="flex items-center gap-1 text-xs text-green-400">
                        <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                        Running
                      </span>
                    ) : protocol.enabled ? (
                      <span className="flex items-center gap-1 text-xs text-yellow-400">
                        <div className="w-2 h-2 bg-yellow-400 rounded-full" />
                        Enabled
                      </span>
                    ) : (
                      <span className="flex items-center gap-1 text-xs text-gray-500">
                        <div className="w-2 h-2 bg-gray-500 rounded-full" />
                        Disabled
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>

            {/* Stats */}
            {protocol.installed && (
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div className="bg-gray-900/50 rounded-xl p-3">
                  <div className="text-2xl font-bold text-blue-400">
                    {protocol.connections}
                  </div>
                  <div className="text-xs text-gray-400">Active Connections</div>
                </div>
                <div className="bg-gray-900/50 rounded-xl p-3">
                  <div className="text-2xl font-bold text-purple-400">
                    {protocol.running ? 'Online' : 'Offline'}
                  </div>
                  <div className="text-xs text-gray-400">Status</div>
                </div>
              </div>
            )}

            {/* Actions */}
            <div className="flex gap-2">
              {!protocol.installed ? (
                <button
                  onClick={() => onInstall(protocol.protocol as VPNProtocol)}
                  disabled={installing === protocol.protocol}
                  className="flex-1 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-4 py-2 rounded-xl font-medium transition-all duration-300 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {installing === protocol.protocol ? (
                    <>
                      <Loader className="w-4 h-4 animate-spin" />
                      Installing...
                    </>
                  ) : (
                    <>
                      <Download className="w-4 h-4" />
                      Install
                    </>
                  )}
                </button>
              ) : (
                <>
                  {protocol.running ? (
                    <button
                      onClick={() => onDisable(protocol.protocol as VPNProtocol)}
                      className="flex-1 bg-red-600/20 hover:bg-red-600/30 border border-red-500/50 text-red-400 px-4 py-2 rounded-xl font-medium transition-all duration-300 flex items-center justify-center gap-2"
                    >
                      <Square className="w-4 h-4" />
                      Disable
                    </button>
                  ) : (
                    <button
                      onClick={() => onEnable(protocol.protocol as VPNProtocol)}
                      className="flex-1 bg-green-600/20 hover:bg-green-600/30 border border-green-500/50 text-green-400 px-4 py-2 rounded-xl font-medium transition-all duration-300 flex items-center justify-center gap-2"
                    >
                      <Play className="w-4 h-4" />
                      Enable
                    </button>
                  )}
                  <button className="bg-gray-700/50 hover:bg-gray-700 text-gray-300 px-4 py-2 rounded-xl transition-all duration-300">
                    <Settings className="w-4 h-4" />
                  </button>
                </>
              )}
            </div>
          </div>
        </motion.div>
      ))}
    </motion.div>
  );
};

// WireGuard Tab Component (placeholder for now)
const WireGuardTab: React.FC = () => {
  const [peers, setPeers] = useState<WireGuardPeer[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPeers();
  }, []);

  const loadPeers = async () => {
    try {
      const response = await vpnApi.wireguard.getPeers();
      if (response.success && response.data) {
        setPeers(response.data);
      }
    } catch (error) {
      console.error('Failed to load peers:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="space-y-6"
    >
      {/* Peer management interface will go here */}
      <div className="bg-gray-800/80 backdrop-blur-xl rounded-3xl p-6 border border-gray-700/50">
        <h2 className="text-2xl font-bold text-white mb-4">WireGuard Peers</h2>
        <p className="text-gray-400">Peer management interface coming soon...</p>
      </div>
    </motion.div>
  );
};

// OpenVPN Tab Component (placeholder)
const OpenVPNTab: React.FC = () => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="space-y-6"
    >
      <div className="bg-gray-800/80 backdrop-blur-xl rounded-3xl p-6 border border-gray-700/50">
        <h2 className="text-2xl font-bold text-white mb-4">OpenVPN Certificates</h2>
        <p className="text-gray-400">Certificate management interface coming soon...</p>
      </div>
    </motion.div>
  );
};

// Settings Tab Component (placeholder)
const SettingsTab: React.FC = () => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="space-y-6"
    >
      <div className="bg-gray-800/80 backdrop-blur-xl rounded-3xl p-6 border border-gray-700/50">
        <h2 className="text-2xl font-bold text-white mb-4">VPN Settings</h2>
        <p className="text-gray-400">Settings interface coming soon...</p>
      </div>
    </motion.div>
  );
};

export default VPNServer;
