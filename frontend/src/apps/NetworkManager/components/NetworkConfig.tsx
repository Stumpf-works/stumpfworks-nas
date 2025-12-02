// Revision: 2025-12-02 | Author: StumpfWorks AI | Version: 2.0.0
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Network,
  RefreshCw,
  Wifi,
  Cable,
  Activity,
  X,
  Info,
  Layers,
  Link2,
  Trash2,
  GitBranch,
  AlertTriangle,
  CheckCircle2,
} from 'lucide-react';
import { networkApi, type NetworkInterface, type PendingChangesResponse } from '@/api/network';
import { syslibApi, type CreateBondRequest, type CreateVLANRequest } from '@/api/syslib';

type DialogType = 'none' | 'bond' | 'vlan' | 'bridge';

export default function NetworkConfig() {
  const [interfaces, setInterfaces] = useState<NetworkInterface[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogType, setDialogType] = useState<DialogType>('none');
  const [showCreateMenu, setShowCreateMenu] = useState(false);
  const [pendingChanges, setPendingChanges] = useState<PendingChangesResponse>({ has_pending: false, count: 0, changes: [] });
  const [isApplying, setIsApplying] = useState(false);

  const [bondFormData, setBondFormData] = useState<CreateBondRequest>({
    name: 'bond0',
    mode: 'balance-rr',
    interfaces: [],
  });

  const [vlanFormData, setVlanFormData] = useState<CreateVLANRequest>({
    parent: '',
    vlan_id: 100,
  });

  const [bridgeFormData, setBridgeFormData] = useState({
    name: 'br0',
    description: '',
    ports: [] as string[],
    ipAddress: '',
    gateway: '',
    ipv6Address: '',
    ipv6Gateway: '',
    vlanAware: false,
    autostart: true,
  });

  // Bond modes
  const bondModes = [
    { value: 'balance-rr', label: 'Balance Round-Robin (0)', description: 'Sequential transmission across all slaves' },
    { value: 'active-backup', label: 'Active-Backup (1)', description: 'One slave active, others on standby' },
    { value: 'balance-xor', label: 'Balance XOR (2)', description: 'XOR hash-based distribution' },
    { value: 'broadcast', label: 'Broadcast (3)', description: 'Transmit on all slaves' },
    { value: '802.3ad', label: '802.3ad LACP (4)', description: 'IEEE 802.3ad Dynamic link aggregation' },
    { value: 'balance-tlb', label: 'Adaptive Transmit Load Balancing (5)', description: 'Outgoing traffic distribution' },
    { value: 'balance-alb', label: 'Adaptive Load Balancing (6)', description: 'TX and RX load balancing' },
  ];

  // Fetch network interfaces
  const fetchInterfaces = async () => {
    setIsLoading(true);
    try {
      const response = await networkApi.listInterfaces();
      if (response.success && response.data) {
        setInterfaces(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch network interfaces:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch pending changes
  const fetchPendingChanges = async () => {
    try {
      const response = await networkApi.getPendingChanges();
      if (response.success && response.data) {
        setPendingChanges(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch pending changes:', error);
    }
  };

  useEffect(() => {
    fetchInterfaces();
    fetchPendingChanges();
  }, []);

  const handleCreateBond = async () => {
    if (!bondFormData.name || bondFormData.interfaces.length < 2) {
      alert('Please provide bond name and select at least 2 interfaces');
      return;
    }

    try {
      const response = await syslibApi.network.createBond(bondFormData);
      if (response.success) {
        alert(`Bond interface created: ${bondFormData.name}`);
        setDialogType('none');
        setBondFormData({ name: 'bond0', mode: 'balance-rr', interfaces: [] });
        fetchInterfaces();
      }
    } catch (error) {
      console.error('Failed to create bond:', error);
      alert('Failed to create bond interface');
    }
  };

  const handleCreateVLAN = async () => {
    if (!vlanFormData.parent || !vlanFormData.vlan_id) {
      alert('Please select parent interface and provide VLAN ID');
      return;
    }

    try {
      const response = await syslibApi.network.createVLAN(vlanFormData);
      if (response.success) {
        alert(`VLAN interface created: ${vlanFormData.parent}.${vlanFormData.vlan_id}`);
        setDialogType('none');
        setVlanFormData({ parent: '', vlan_id: 100 });
        fetchInterfaces();
      }
    } catch (error) {
      console.error('Failed to create VLAN:', error);
      alert('Failed to create VLAN interface');
    }
  };

  const handleDeleteBond = async (name: string) => {
    if (!confirm(`Are you sure you want to delete bond interface "${name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const response = await syslibApi.network.deleteBond(name);
      if (response.success) {
        alert(`Bond interface "${name}" deleted successfully`);
        fetchInterfaces();
      }
    } catch (error) {
      console.error('Failed to delete bond:', error);
      alert('Failed to delete bond interface');
    }
  };

  const handleDeleteVLAN = async (name: string) => {
    if (!confirm(`Are you sure you want to delete VLAN interface "${name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      // Parse parent and VLAN ID from name (e.g., "eth0.100" -> parent: "eth0", vlanId: 100)
      const parts = name.split('.');
      if (parts.length !== 2) {
        alert('Invalid VLAN interface name');
        return;
      }
      const parent = parts[0];
      const vlanId = parseInt(parts[1], 10);

      const response = await syslibApi.network.deleteVLAN(parent, vlanId);
      if (response.success) {
        alert(`VLAN interface "${name}" deleted successfully`);
        fetchInterfaces();
      }
    } catch (error) {
      console.error('Failed to delete VLAN:', error);
      alert('Failed to delete VLAN interface');
    }
  };

  const toggleBondInterface = (ifName: string) => {
    setBondFormData((prev) => ({
      ...prev,
      interfaces: prev.interfaces.includes(ifName)
        ? prev.interfaces.filter((i) => i !== ifName)
        : [...prev.interfaces, ifName],
    }));
  };

  const toggleBridgePort = (ifName: string) => {
    setBridgeFormData((prev) => ({
      ...prev,
      ports: prev.ports.includes(ifName)
        ? prev.ports.filter((i) => i !== ifName)
        : [...prev.ports, ifName],
    }));
  };

  // CIDR validation helper
  const validateCIDR = (cidr: string): boolean => {
    if (!cidr) return true; // Empty is valid (optional field)
    const cidrRegex = /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/;
    if (!cidrRegex.test(cidr)) return false;

    // Validate IP octets
    const [ip, prefix] = cidr.split('/');
    const octets = ip.split('.').map(Number);
    if (octets.some(o => o < 0 || o > 255)) return false;

    // Validate prefix
    const prefixNum = parseInt(prefix);
    if (prefixNum < 0 || prefixNum > 32) return false;

    return true;
  };

  const validateIPv6CIDR = (cidr: string): boolean => {
    if (!cidr) return true; // Empty is valid (optional field)
    // Basic IPv6 CIDR validation (simplified)
    const ipv6CidrRegex = /^([0-9a-fA-F:]+)\/\d{1,3}$/;
    if (!ipv6CidrRegex.test(cidr)) return false;

    const [, prefix] = cidr.split('/');
    const prefixNum = parseInt(prefix);
    if (prefixNum < 0 || prefixNum > 128) return false;

    return true;
  };

  const handleCreateBridge = async () => {
    if (!bridgeFormData.name) {
      alert('Please provide bridge name');
      return;
    }

    // Validate CIDR formats
    if (!validateCIDR(bridgeFormData.ipAddress)) {
      alert('Invalid IPv4 CIDR format. Use format: 192.168.1.10/24');
      return;
    }

    if (!validateIPv6CIDR(bridgeFormData.ipv6Address)) {
      alert('Invalid IPv6 CIDR format. Use format: 2001:db8::1/64');
      return;
    }

    try {
      const response = await networkApi.createBridgeWithPendingChanges(
        bridgeFormData.name,
        bridgeFormData.description,
        bridgeFormData.ports,
        bridgeFormData.ipAddress || undefined,
        bridgeFormData.gateway || undefined,
        bridgeFormData.ipv6Address || undefined,
        bridgeFormData.ipv6Gateway || undefined,
        bridgeFormData.vlanAware,
        bridgeFormData.autostart
      );
      if (response.success) {
        alert(`Bridge "${bridgeFormData.name}" added to pending changes. Click "Apply Configuration" to create it.`);
        setDialogType('none');
        setBridgeFormData({
          name: 'br0',
          description: '',
          ports: [],
          ipAddress: '',
          gateway: '',
          ipv6Address: '',
          ipv6Gateway: '',
          vlanAware: false,
          autostart: true
        });
        fetchPendingChanges();
      }
    } catch (error) {
      console.error('Failed to create bridge:', error);
      alert('Failed to create bridge interface');
    }
  };

  const handleApplyChanges = async () => {
    if (!confirm(`Apply ${pendingChanges.count} pending network change(s)? This will modify your network configuration.`)) {
      return;
    }

    setIsApplying(true);
    try {
      const response = await networkApi.applyPendingChanges();
      if (response.success) {
        alert('All pending changes applied successfully!');
        fetchPendingChanges();
        fetchInterfaces();
      }
    } catch (error) {
      console.error('Failed to apply changes:', error);
      alert('Failed to apply pending changes. Network configuration has been rolled back.');
    } finally {
      setIsApplying(false);
    }
  };

  const handleDiscardChanges = async () => {
    if (!confirm(`Discard all ${pendingChanges.count} pending change(s)? This action cannot be undone.`)) {
      return;
    }

    try {
      const response = await networkApi.discardPendingChanges();
      if (response.success) {
        alert('All pending changes discarded.');
        fetchPendingChanges();
      }
    } catch (error) {
      console.error('Failed to discard changes:', error);
      alert('Failed to discard pending changes');
    }
  };

  const handleDeleteBridge = async (name: string) => {
    if (!confirm(`Are you sure you want to delete bridge interface "${name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const response = await networkApi.deleteBridge(name);
      if (response.success) {
        alert(`Bridge interface "${name}" deleted successfully`);
        fetchInterfaces();
      }
    } catch (error) {
      console.error('Failed to delete bridge:', error);
      alert('Failed to delete bridge interface');
    }
  };

  const getInterfaceIcon = (iface: NetworkInterface) => {
    if (iface.name.startsWith('wl')) {
      return <Wifi className="w-5 h-5 text-blue-500" />;
    } else if (iface.name.startsWith('br') || iface.name.startsWith('vmbr')) {
      return <GitBranch className="w-5 h-5 text-cyan-500" />;
    } else if (iface.name.startsWith('bond')) {
      return <Link2 className="w-5 h-5 text-purple-500" />;
    } else if (iface.name.includes('.')) {
      return <Layers className="w-5 h-5 text-orange-500" />;
    } else {
      return <Cable className="w-5 h-5 text-green-500" />;
    }
  };

  const getInterfaceType = (iface: NetworkInterface): string => {
    if (iface.name.startsWith('wl')) return 'Wireless';
    if (iface.name.startsWith('br') || iface.name.startsWith('vmbr')) return 'Bridge';
    if (iface.name.startsWith('bond')) return 'Bond';
    if (iface.name.includes('.')) return 'VLAN';
    if (iface.name.startsWith('lo')) return 'Loopback';
    if (iface.name.startsWith('en')) return 'Ethernet';
    if (iface.name.startsWith('eth')) return 'Ethernet';
    return 'Unknown';
  };

  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-3">
          <Network className="w-6 h-6 text-macos-blue" />
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
            Network Configuration
          </h1>
        </div>
        <div className="flex gap-2">
          <button
            onClick={fetchInterfaces}
            className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>

          {/* Proxmox-style Create Dropdown */}
          <div className="relative">
            <button
              onClick={() => setShowCreateMenu(!showCreateMenu)}
              className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              Create
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </button>

            {/* Dropdown Menu */}
            {showCreateMenu && (
              <>
                {/* Backdrop */}
                <div
                  className="fixed inset-0 z-10"
                  onClick={() => setShowCreateMenu(false)}
                />

                {/* Menu */}
                <motion.div
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="absolute right-0 mt-2 w-56 bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 z-20 overflow-hidden"
                >
                  {/* Linux Bridge */}
                  <button
                    onClick={() => {
                      setDialogType('bridge');
                      setShowCreateMenu(false);
                    }}
                    className="w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-100 dark:hover:bg-macos-dark-200 transition-colors text-left"
                  >
                    <GitBranch className="w-5 h-5 text-cyan-500" />
                    <div>
                      <div className="font-medium text-gray-900 dark:text-gray-100">Linux Bridge</div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">Network bridge interface</div>
                    </div>
                  </button>

                  {/* Linux Bond */}
                  <button
                    onClick={() => {
                      setDialogType('bond');
                      setShowCreateMenu(false);
                    }}
                    className="w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-100 dark:hover:bg-macos-dark-200 transition-colors text-left border-t border-gray-200 dark:border-gray-700"
                  >
                    <Link2 className="w-5 h-5 text-purple-500" />
                    <div>
                      <div className="font-medium text-gray-900 dark:text-gray-100">Linux Bond</div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">Link aggregation</div>
                    </div>
                  </button>

                  {/* Linux VLAN */}
                  <button
                    onClick={() => {
                      setDialogType('vlan');
                      setShowCreateMenu(false);
                    }}
                    className="w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-100 dark:hover:bg-macos-dark-200 transition-colors text-left border-t border-gray-200 dark:border-gray-700"
                  >
                    <Layers className="w-5 h-5 text-orange-500" />
                    <div>
                      <div className="font-medium text-gray-900 dark:text-gray-100">Linux VLAN</div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">VLAN interface (802.1Q)</div>
                    </div>
                  </button>
                </motion.div>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Pending Changes Banner */}
      {pendingChanges.has_pending && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mx-6 mt-4 bg-gradient-to-r from-orange-50 to-yellow-50 dark:from-orange-900/20 dark:to-yellow-900/20 border-2 border-orange-200 dark:border-orange-800 rounded-xl p-4"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-orange-100 dark:bg-orange-900/40 rounded-lg">
                <AlertTriangle className="w-5 h-5 text-orange-600 dark:text-orange-400" />
              </div>
              <div>
                <h3 className="font-bold text-gray-900 dark:text-gray-100">
                  {pendingChanges.count} Pending Network Change{pendingChanges.count !== 1 ? 's' : ''}
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Review and apply your changes to modify the network configuration
                </p>
              </div>
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleDiscardChanges}
                disabled={isApplying}
                className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-macos-dark-300 transition-colors border border-gray-300 dark:border-gray-600 disabled:opacity-50"
              >
                <X className="w-4 h-4" />
                Discard All
              </button>
              <button
                onClick={handleApplyChanges}
                disabled={isApplying}
                className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-lg hover:from-green-600 hover:to-emerald-600 transition-all shadow-lg disabled:opacity-50"
              >
                {isApplying ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                    Applying...
                  </>
                ) : (
                  <>
                    <CheckCircle2 className="w-4 h-4" />
                    Apply Configuration
                  </>
                )}
              </button>
            </div>
          </div>

          {/* Pending Changes List */}
          <div className="mt-4 space-y-2">
            {pendingChanges.changes.map((change) => (
              <div
                key={change.id}
                className="bg-white dark:bg-macos-dark-100 rounded-lg p-3 border border-orange-200 dark:border-orange-800"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="px-2 py-1 bg-macos-blue/10 dark:bg-macos-blue/20 text-macos-blue text-xs font-medium rounded">
                      {change.change_type}
                    </span>
                    <span className="px-2 py-1 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 text-xs font-medium rounded">
                      {change.action}
                    </span>
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {change.resource_id}
                    </span>
                  </div>
                  {change.description && (
                    <span className="text-xs text-gray-600 dark:text-gray-400">
                      {change.description}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </motion.div>
      )}

      {/* Content - Proxmox-style Table (StumpfWorks Design) */}
      <div className="flex-1 overflow-y-auto p-6">
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
          </div>
        ) : interfaces.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-64 text-center">
            <Network className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" />
            <p className="text-gray-500 dark:text-gray-400">No network interfaces found</p>
          </div>
        ) : (
          <div className="bg-white dark:bg-macos-dark-100 rounded-xl border border-gray-200 dark:border-gray-700 shadow-lg overflow-hidden">
            {/* Table */}
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-macos-dark-200 dark:to-macos-dark-300 border-b-2 border-gray-200 dark:border-gray-700">
                    <th className="px-4 py-3 text-left text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Name</th>
                    <th className="px-4 py-3 text-left text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Type</th>
                    <th className="px-4 py-3 text-center text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Active</th>
                    <th className="px-4 py-3 text-center text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Autostart</th>
                    <th className="px-4 py-3 text-center text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">VLAN aware</th>
                    <th className="px-4 py-3 text-left text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Ports/Slaves</th>
                    <th className="px-4 py-3 text-left text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">CIDR</th>
                    <th className="px-4 py-3 text-left text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Gateway</th>
                    <th className="px-4 py-3 text-center text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {interfaces
                    .filter((iface) => !iface.name.startsWith('lo')) // Hide loopback
                    .map((iface, index) => (
                      <motion.tr
                        key={iface.name}
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.03 }}
                        className={`hover:bg-gray-50 dark:hover:bg-macos-dark-200 transition-colors ${
                          iface.isUp
                            ? 'bg-green-50/30 dark:bg-green-900/10'
                            : 'bg-gray-50/50 dark:bg-macos-dark-200/50'
                        }`}
                      >
                        {/* Name */}
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            {getInterfaceIcon(iface)}
                            <span className="font-semibold text-gray-900 dark:text-gray-100">
                              {iface.name}
                            </span>
                          </div>
                        </td>

                        {/* Type */}
                        <td className="px-4 py-3">
                          <span className="text-sm text-gray-700 dark:text-gray-300">
                            {getInterfaceType(iface)}
                          </span>
                        </td>

                        {/* Active */}
                        <td className="px-4 py-3 text-center">
                          <span
                            className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${
                              iface.isUp
                                ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300'
                                : 'bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
                            }`}
                          >
                            <Activity className="w-3 h-3" />
                            {iface.isUp ? 'Yes' : 'No'}
                          </span>
                        </td>

                        {/* Autostart */}
                        <td className="px-4 py-3 text-center">
                          <span className="text-sm text-gray-700 dark:text-gray-300">-</span>
                        </td>

                        {/* VLAN aware */}
                        <td className="px-4 py-3 text-center">
                          <span className="text-sm text-gray-700 dark:text-gray-300">-</span>
                        </td>

                        {/* Ports/Slaves */}
                        <td className="px-4 py-3">
                          <span className="text-sm text-gray-700 dark:text-gray-300">-</span>
                        </td>

                        {/* CIDR */}
                        <td className="px-4 py-3">
                          {iface.addresses && iface.addresses.length > 0 ? (
                            <div className="space-y-1">
                              {iface.addresses.slice(0, 2).map((addr, idx) => (
                                <div key={idx} className="font-mono text-xs text-gray-900 dark:text-gray-100">
                                  {addr}
                                </div>
                              ))}
                              {iface.addresses.length > 2 && (
                                <div className="text-xs text-gray-500">
                                  +{iface.addresses.length - 2} more
                                </div>
                              )}
                            </div>
                          ) : (
                            <span className="text-sm text-gray-500 dark:text-gray-400">-</span>
                          )}
                        </td>

                        {/* Gateway */}
                        <td className="px-4 py-3">
                          <span className="text-sm text-gray-700 dark:text-gray-300">-</span>
                        </td>

                        {/* Actions */}
                        <td className="px-4 py-3">
                          <div className="flex items-center justify-center gap-2">
                            {/* Edit button - coming soon for physical interfaces */}
                            <button
                              className="p-1.5 bg-macos-blue/10 dark:bg-macos-blue/20 text-macos-blue rounded hover:bg-macos-blue/20 dark:hover:bg-macos-blue/30 transition-colors"
                              title="Edit interface"
                            >
                              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                              </svg>
                            </button>

                            {/* Delete button for Bond, Bridge, and VLAN interfaces */}
                            {(iface.name.startsWith('bond') || iface.name.startsWith('br') || iface.name.startsWith('vmbr') || iface.name.includes('.')) &&
                             !iface.name.startsWith('br-') && (
                              <button
                                onClick={() => {
                                  if (iface.name.startsWith('bond')) {
                                    handleDeleteBond(iface.name);
                                  } else if (iface.name.startsWith('br') || iface.name.startsWith('vmbr')) {
                                    handleDeleteBridge(iface.name);
                                  } else if (iface.name.includes('.')) {
                                    handleDeleteVLAN(iface.name);
                                  }
                                }}
                                className="p-1.5 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors"
                                title={`Delete ${getInterfaceType(iface)}`}
                              >
                                <Trash2 className="w-4 h-4" />
                              </button>
                            )}
                          </div>
                        </td>
                      </motion.tr>
                    ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>

      {/* Create Bond Dialog */}
      {dialogType === 'bond' && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full m-4 max-h-[80vh] overflow-y-auto"
          >
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create Bond Interface
              </h3>
              <button
                onClick={() => setDialogType('none')}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              {/* Bond Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Bond Name
                </label>
                <input
                  type="text"
                  value={bondFormData.name}
                  onChange={(e) => setBondFormData({ ...bondFormData, name: e.target.value })}
                  placeholder="bond0"
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                />
              </div>

              {/* Bond Mode */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Bonding Mode
                </label>
                <div className="space-y-2">
                  {bondModes.map((mode) => (
                    <button
                      key={mode.value}
                      onClick={() => setBondFormData({ ...bondFormData, mode: mode.value })}
                      className={`w-full p-3 rounded-lg border-2 transition-all text-left ${
                        bondFormData.mode === mode.value
                          ? 'border-macos-blue bg-macos-blue/10 dark:bg-macos-blue/20'
                          : 'border-gray-300 dark:border-gray-600 hover:border-macos-blue/50'
                      }`}
                    >
                      <div className="font-semibold text-gray-900 dark:text-gray-100">
                        {mode.label}
                      </div>
                      <div className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                        {mode.description}
                      </div>
                    </button>
                  ))}
                </div>
              </div>

              {/* Interface Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Select Interfaces ({bondFormData.interfaces.length} selected)
                </label>
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {interfaces
                    .filter((iface) => !iface.name.startsWith('lo') && !iface.name.startsWith('bond'))
                    .map((iface) => (
                      <div
                        key={iface.name}
                        className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                      >
                        <input
                          type="checkbox"
                          checked={bondFormData.interfaces.includes(iface.name)}
                          onChange={() => toggleBondInterface(iface.name)}
                          className="w-4 h-4 text-macos-blue rounded focus:ring-2 focus:ring-macos-blue"
                        />
                        <div className="flex items-center gap-2 flex-1">
                          {getInterfaceIcon(iface)}
                          <div>
                            <div className="font-medium text-gray-900 dark:text-gray-100">
                              {iface.name}
                            </div>
                            <div className="text-xs text-gray-600 dark:text-gray-400">
                              {iface.hardwareAddr} • {getInterfaceType(iface)}
                            </div>
                          </div>
                        </div>
                        <span
                          className={`text-xs px-2 py-1 rounded ${
                            iface.isUp
                              ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300'
                              : 'bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
                          }`}
                        >
                          {iface.isUp ? 'UP' : 'DOWN'}
                        </span>
                      </div>
                    ))}
                </div>
              </div>

              {/* Info */}
              <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                <div className="flex gap-2">
                  <Info className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <div className="text-sm text-gray-700 dark:text-gray-300">
                    <p className="font-medium mb-1">Bond Interface Notes:</p>
                    <ul className="list-disc list-inside space-y-1 text-xs">
                      <li>Select at least 2 interfaces to create a bond</li>
                      <li>All interfaces in a bond should have similar characteristics</li>
                      <li>802.3ad requires switch support for LACP</li>
                      <li>Active-backup provides the simplest failover</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <button
                onClick={() => setDialogType('none')}
                className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateBond}
                className="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors"
              >
                Create Bond
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Create VLAN Dialog */}
      {dialogType === 'vlan' && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-md w-full m-4"
          >
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create VLAN Interface
              </h3>
              <button
                onClick={() => setDialogType('none')}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              {/* Parent Interface */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Parent Interface
                </label>
                <select
                  value={vlanFormData.parent}
                  onChange={(e) => setVlanFormData({ ...vlanFormData, parent: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                >
                  <option value="">Select interface...</option>
                  {interfaces
                    .filter((iface) => !iface.name.startsWith('lo') && !iface.name.includes('.'))
                    .map((iface) => (
                      <option key={iface.name} value={iface.name}>
                        {iface.name} ({getInterfaceType(iface)})
                      </option>
                    ))}
                </select>
              </div>

              {/* VLAN ID */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  VLAN ID (1-4094)
                </label>
                <input
                  type="number"
                  min="1"
                  max="4094"
                  value={vlanFormData.vlan_id}
                  onChange={(e) =>
                    setVlanFormData({ ...vlanFormData, vlan_id: parseInt(e.target.value) || 100 })
                  }
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                />
              </div>

              {/* Preview */}
              {vlanFormData.parent && (
                <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-3">
                  <span className="text-sm text-gray-600 dark:text-gray-400">
                    Interface will be created as:
                  </span>
                  <div className="font-mono font-bold text-macos-blue mt-1">
                    {vlanFormData.parent}.{vlanFormData.vlan_id}
                  </div>
                </div>
              )}

              {/* Info */}
              <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                <div className="flex gap-2">
                  <Info className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <div className="text-sm text-gray-700 dark:text-gray-300">
                    <p className="font-medium mb-1">VLAN Notes:</p>
                    <ul className="list-disc list-inside space-y-1 text-xs">
                      <li>VLANs enable network segmentation</li>
                      <li>Switch must support 802.1Q tagging</li>
                      <li>Valid VLAN IDs: 1-4094</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <button
                onClick={() => setDialogType('none')}
                className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateVLAN}
                className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                Create VLAN
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Create Bridge Dialog */}
      {dialogType === 'bridge' && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full m-4 max-h-[80vh] overflow-y-auto"
          >
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create Bridge Interface
              </h3>
              <button
                onClick={() => setDialogType('none')}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              {/* Bridge Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Bridge Name
                </label>
                <input
                  type="text"
                  value={bridgeFormData.name}
                  onChange={(e) => setBridgeFormData({ ...bridgeFormData, name: e.target.value })}
                  placeholder="br0 or vmbr0"
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
                />
              </div>

              {/* Description */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Description (optional)
                </label>
                <input
                  type="text"
                  value={bridgeFormData.description}
                  onChange={(e) => setBridgeFormData({ ...bridgeFormData, description: e.target.value })}
                  placeholder="e.g., Main network bridge"
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
                />
              </div>

              {/* IPv4 Configuration (Proxmox-style) */}
              <div>
                <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100 mb-2">
                  IPv4 Configuration
                </h4>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      IPv4/CIDR (optional)
                    </label>
                    <input
                      type="text"
                      value={bridgeFormData.ipAddress}
                      onChange={(e) => setBridgeFormData({ ...bridgeFormData, ipAddress: e.target.value })}
                      placeholder="192.168.1.10/24"
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-cyan-500 focus:border-transparent text-sm font-mono"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Gateway (IPv4)
                    </label>
                    <input
                      type="text"
                      value={bridgeFormData.gateway}
                      onChange={(e) => setBridgeFormData({ ...bridgeFormData, gateway: e.target.value })}
                      placeholder="192.168.1.1"
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-cyan-500 focus:border-transparent text-sm font-mono"
                    />
                  </div>
                </div>
              </div>

              {/* IPv6 Configuration (Proxmox-style) */}
              <div>
                <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100 mb-2">
                  IPv6 Configuration
                </h4>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      IPv6/CIDR (optional)
                    </label>
                    <input
                      type="text"
                      value={bridgeFormData.ipv6Address}
                      onChange={(e) => setBridgeFormData({ ...bridgeFormData, ipv6Address: e.target.value })}
                      placeholder="2001:db8::1/64"
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-cyan-500 focus:border-transparent text-sm font-mono"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Gateway (IPv6)
                    </label>
                    <input
                      type="text"
                      value={bridgeFormData.ipv6Gateway}
                      onChange={(e) => setBridgeFormData({ ...bridgeFormData, ipv6Gateway: e.target.value })}
                      placeholder="2001:db8::ffff"
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-cyan-500 focus:border-transparent text-sm font-mono"
                    />
                  </div>
                </div>
              </div>

              {/* VLAN Aware (Proxmox feature) */}
              <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg">
                <input
                  type="checkbox"
                  checked={bridgeFormData.vlanAware}
                  onChange={(e) => setBridgeFormData({ ...bridgeFormData, vlanAware: e.target.checked })}
                  className="w-4 h-4 text-cyan-500 rounded focus:ring-2 focus:ring-cyan-500"
                />
                <div>
                  <label className="text-sm font-medium text-gray-900 dark:text-gray-100">
                    VLAN aware
                  </label>
                  <p className="text-xs text-gray-600 dark:text-gray-400">
                    Enable if you want to use VLANs on this bridge (Proxmox feature)
                  </p>
                </div>
              </div>

              {/* Autostart */}
              <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg">
                <input
                  type="checkbox"
                  checked={bridgeFormData.autostart}
                  onChange={(e) => setBridgeFormData({ ...bridgeFormData, autostart: e.target.checked })}
                  className="w-4 h-4 text-cyan-500 rounded focus:ring-2 focus:ring-cyan-500"
                />
                <div>
                  <label className="text-sm font-medium text-gray-900 dark:text-gray-100">
                    Auto-start on system boot
                  </label>
                  <p className="text-xs text-gray-600 dark:text-gray-400">
                    Automatically create and configure this bridge when the system starts
                  </p>
                </div>
              </div>

              {/* Port Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Bridge Ports ({bridgeFormData.ports.length} selected, optional)
                </label>
                <p className="text-xs text-gray-600 dark:text-gray-400 mb-3">
                  Select physical interfaces to attach to this bridge, or leave empty to create an isolated bridge
                </p>
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {interfaces
                    .filter((iface) =>
                      !iface.name.startsWith('lo') &&
                      !iface.name.startsWith('br') &&
                      !iface.name.startsWith('vmbr') &&
                      !iface.name.startsWith('bond') &&
                      !iface.name.startsWith('docker')
                    )
                    .map((iface) => (
                      <div
                        key={iface.name}
                        className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                      >
                        <input
                          type="checkbox"
                          checked={bridgeFormData.ports.includes(iface.name)}
                          onChange={() => toggleBridgePort(iface.name)}
                          className="w-4 h-4 text-cyan-500 rounded focus:ring-2 focus:ring-cyan-500"
                        />
                        <div className="flex items-center gap-2 flex-1">
                          {getInterfaceIcon(iface)}
                          <div>
                            <div className="font-medium text-gray-900 dark:text-gray-100">
                              {iface.name}
                            </div>
                            <div className="text-xs text-gray-600 dark:text-gray-400">
                              {iface.hardwareAddr} • {getInterfaceType(iface)}
                            </div>
                          </div>
                        </div>
                        <span
                          className={`text-xs px-2 py-1 rounded ${
                            iface.isUp
                              ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300'
                              : 'bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
                          }`}
                        >
                          {iface.isUp ? 'UP' : 'DOWN'}
                        </span>
                      </div>
                    ))}
                </div>
              </div>

              {/* Info */}
              <div className="bg-cyan-50 dark:bg-cyan-900/20 rounded-lg p-4">
                <div className="flex gap-2">
                  <Info className="w-5 h-5 text-cyan-600 dark:text-cyan-400 flex-shrink-0" />
                  <div className="text-sm text-gray-700 dark:text-gray-300">
                    <p className="font-medium mb-1">Proxmox-Style Pending Changes Workflow:</p>
                    <ul className="list-disc list-inside space-y-1 text-xs">
                      <li><strong>Create Without Applying:</strong> Configuration is saved but NOT applied to the system</li>
                      <li><strong>Review Changes:</strong> All pending changes are shown in a banner above</li>
                      <li><strong>Apply Configuration:</strong> Click "Apply Configuration" to apply all changes atomically</li>
                      <li><strong>Automatic Rollback:</strong> If anything fails, all changes are rolled back automatically</li>
                      <li><strong>Safe Operation:</strong> Network connectivity is protected during changes</li>
                      <li>Perfect for VMs and containers to share the same network as the host</li>
                      <li>Empty bridges are useful for isolated VM networks (like OPNsense WAN/LAN)</li>
                    </ul>
                  </div>
                </div>
              </div>

              {/* Warning for interfaces with IP addresses */}
              {bridgeFormData.ports.some((portName) => {
                const iface = interfaces.find((i) => i.name === portName);
                return iface && iface.addresses && iface.addresses.length > 0;
              }) && (
                <div className="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4 border border-yellow-200 dark:border-yellow-800">
                  <div className="flex gap-2">
                    <Info className="w-5 h-5 text-yellow-600 dark:text-yellow-400 flex-shrink-0" />
                    <div className="text-sm text-gray-700 dark:text-gray-300">
                      <p className="font-medium mb-1 text-yellow-800 dark:text-yellow-300">IP Migration Active:</p>
                      <p className="text-xs">
                        One or more selected ports have IP addresses. These will be automatically migrated to the bridge "{bridgeFormData.name}".
                        The physical interfaces will become bridge ports without IP addresses. This is the Proxmox-style configuration.
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <button
                onClick={() => setDialogType('none')}
                className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateBridge}
                className="px-4 py-2 bg-cyan-500 text-white rounded-lg hover:bg-cyan-600 transition-colors"
              >
                Add to Pending Changes
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
