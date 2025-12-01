import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  Settings,
  Network,
  Shield,
  Users,
  Globe,
  Server,
  Key,
  Clock,
  Save,
  RefreshCw
} from 'lucide-react';

interface NetworkInterface {
  name: string;
  ipAddress: string;
  type: string;
}

const GeneralSettings: React.FC = () => {
  const [networkInterfaces] = useState<NetworkInterface[]>([
    { name: 'eth0', ipAddress: '192.168.1.100', type: 'Ethernet' },
    { name: 'wlan0', ipAddress: '192.168.1.101', type: 'WiFi' }
  ]);

  const [settings, setSettings] = useState({
    defaultInterface: 'eth0',
    accountSource: 'local',
    enableLogging: true,
    logLevel: 'info',
    maxConcurrentConnections: 100,
    connectionTimeout: 300,
    enableIPv6: false,
    dnsServers: '8.8.8.8, 8.8.4.4',
    defaultGateway: 'auto',
    enableNAT: true,
    forwardingRules: 'allow-all'
  });

  const handleSave = () => {
    // Save settings logic
    console.log('Saving settings:', settings);
  };

  return (
    <div className="space-y-6">
      {/* Network Configuration */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-purple-600/10 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl">
              <Network className="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 className="text-xl font-bold text-white">Network Configuration</h3>
              <p className="text-sm text-purple-200">Configure network interfaces and routing</p>
            </div>
          </div>

          <div className="space-y-4">
            {/* Network Interface Selection */}
            <div>
              <label className="block text-sm font-medium text-purple-200 mb-2">
                Default Network Interface
              </label>
              <select
                value={settings.defaultInterface}
                onChange={(e) => setSettings({ ...settings, defaultInterface: e.target.value })}
                className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-purple-500 transition-colors"
              >
                {networkInterfaces.map(iface => (
                  <option key={iface.name} value={iface.name} className="bg-slate-800">
                    {iface.name} - {iface.ipAddress} ({iface.type})
                  </option>
                ))}
              </select>
              <p className="text-xs text-purple-300 mt-2">
                Select the network interface for VPN traffic
              </p>
            </div>

            {/* DNS Servers */}
            <div>
              <label className="block text-sm font-medium text-purple-200 mb-2">
                DNS Servers
              </label>
              <input
                type="text"
                value={settings.dnsServers}
                onChange={(e) => setSettings({ ...settings, dnsServers: e.target.value })}
                placeholder="8.8.8.8, 8.8.4.4"
                className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-purple-300 focus:outline-none focus:border-purple-500 transition-colors"
              />
              <p className="text-xs text-purple-300 mt-2">
                Comma-separated list of DNS servers for VPN clients
              </p>
            </div>

            {/* Enable NAT */}
            <div className="flex items-center justify-between p-4 bg-white/5 rounded-xl border border-white/10">
              <div>
                <p className="text-white font-medium">Enable Network Address Translation (NAT)</p>
                <p className="text-sm text-purple-300">Allow VPN clients to access the internet</p>
              </div>
              <label className="relative inline-block w-12 h-6 cursor-pointer">
                <input
                  type="checkbox"
                  checked={settings.enableNAT}
                  onChange={(e) => setSettings({ ...settings, enableNAT: e.target.checked })}
                  className="sr-only peer"
                />
                <div className="w-12 h-6 bg-white/10 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-gradient-to-r peer-checked:from-blue-500 peer-checked:to-purple-600"></div>
              </label>
            </div>

            {/* Enable IPv6 */}
            <div className="flex items-center justify-between p-4 bg-white/5 rounded-xl border border-white/10">
              <div>
                <p className="text-white font-medium">Enable IPv6 Support</p>
                <p className="text-sm text-purple-300">Allow IPv6 traffic through VPN tunnels</p>
              </div>
              <label className="relative inline-block w-12 h-6 cursor-pointer">
                <input
                  type="checkbox"
                  checked={settings.enableIPv6}
                  onChange={(e) => setSettings({ ...settings, enableIPv6: e.target.checked })}
                  className="sr-only peer"
                />
                <div className="w-12 h-6 bg-white/10 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-gradient-to-r peer-checked:from-blue-500 peer-checked:to-purple-600"></div>
              </label>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Account Management */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-green-500/10 to-emerald-600/10 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl">
              <Users className="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 className="text-xl font-bold text-white">Account Management</h3>
              <p className="text-sm text-purple-200">Configure user authentication sources</p>
            </div>
          </div>

          <div className="space-y-4">
            {/* Account Source */}
            <div>
              <label className="block text-sm font-medium text-purple-200 mb-2">
                Account Source
              </label>
              <select
                value={settings.accountSource}
                onChange={(e) => setSettings({ ...settings, accountSource: e.target.value })}
                className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-green-500 transition-colors"
              >
                <option value="local" className="bg-slate-800">Local Users</option>
                <option value="ldap" className="bg-slate-800">LDAP</option>
                <option value="ad" className="bg-slate-800">Active Directory</option>
                <option value="radius" className="bg-slate-800">RADIUS</option>
              </select>
              <p className="text-xs text-purple-300 mt-2">
                Choose where VPN users are authenticated
              </p>
            </div>

            {/* Max Concurrent Connections */}
            <div>
              <label className="block text-sm font-medium text-purple-200 mb-2">
                Maximum Concurrent Connections
              </label>
              <input
                type="number"
                value={settings.maxConcurrentConnections}
                onChange={(e) => setSettings({ ...settings, maxConcurrentConnections: parseInt(e.target.value) })}
                min="1"
                max="1000"
                className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-green-500 transition-colors"
              />
              <p className="text-xs text-purple-300 mt-2">
                Total number of simultaneous VPN connections allowed
              </p>
            </div>

            {/* Connection Timeout */}
            <div>
              <label className="block text-sm font-medium text-purple-200 mb-2">
                Connection Timeout (seconds)
              </label>
              <input
                type="number"
                value={settings.connectionTimeout}
                onChange={(e) => setSettings({ ...settings, connectionTimeout: parseInt(e.target.value) })}
                min="60"
                max="3600"
                className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-green-500 transition-colors"
              />
              <p className="text-xs text-purple-300 mt-2">
                Idle connection timeout in seconds
              </p>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Logging & Monitoring */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="relative overflow-hidden rounded-2xl"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-pink-600/10 backdrop-blur-xl" />
        <div className="relative p-6 border border-white/10">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl">
              <Shield className="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 className="text-xl font-bold text-white">Logging & Monitoring</h3>
              <p className="text-sm text-purple-200">Configure system logging and alerts</p>
            </div>
          </div>

          <div className="space-y-4">
            {/* Enable Logging */}
            <div className="flex items-center justify-between p-4 bg-white/5 rounded-xl border border-white/10">
              <div>
                <p className="text-white font-medium">Enable Connection Logging</p>
                <p className="text-sm text-purple-300">Log all VPN connection attempts and activities</p>
              </div>
              <label className="relative inline-block w-12 h-6 cursor-pointer">
                <input
                  type="checkbox"
                  checked={settings.enableLogging}
                  onChange={(e) => setSettings({ ...settings, enableLogging: e.target.checked })}
                  className="sr-only peer"
                />
                <div className="w-12 h-6 bg-white/10 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-gradient-to-r peer-checked:from-purple-500 peer-checked:to-pink-600"></div>
              </label>
            </div>

            {/* Log Level */}
            <div>
              <label className="block text-sm font-medium text-purple-200 mb-2">
                Log Level
              </label>
              <select
                value={settings.logLevel}
                onChange={(e) => setSettings({ ...settings, logLevel: e.target.value })}
                disabled={!settings.enableLogging}
                className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-purple-500 transition-colors disabled:opacity-50"
              >
                <option value="error" className="bg-slate-800">Error</option>
                <option value="warning" className="bg-slate-800">Warning</option>
                <option value="info" className="bg-slate-800">Info</option>
                <option value="debug" className="bg-slate-800">Debug</option>
              </select>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Save Button */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="flex items-center justify-end gap-4"
      >
        <motion.button
          onClick={handleSave}
          className="flex items-center gap-2 px-8 py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl font-medium shadow-lg hover:shadow-xl transition-shadow"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <Save className="w-5 h-5" />
          Save Settings
        </motion.button>
      </motion.div>
    </div>
  );
};

export default GeneralSettings;
