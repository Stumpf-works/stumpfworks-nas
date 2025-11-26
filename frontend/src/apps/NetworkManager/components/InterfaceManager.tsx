import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { networkApi, NetworkInterface } from '@/api/network';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';

// IP address validation regex
const IPV4_REGEX = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

// Validate IP address format
const isValidIP = (ip: string): boolean => {
  return IPV4_REGEX.test(ip.trim());
};

// Validate netmask format
const isValidNetmask = (netmask: string): boolean => {
  if (!IPV4_REGEX.test(netmask.trim())) return false;

  // Common netmasks
  const validNetmasks = [
    '255.255.255.255', '255.255.255.254', '255.255.255.252', '255.255.255.248',
    '255.255.255.240', '255.255.255.224', '255.255.255.192', '255.255.255.128',
    '255.255.255.0', '255.255.254.0', '255.255.252.0', '255.255.248.0',
    '255.255.240.0', '255.255.224.0', '255.255.192.0', '255.255.128.0',
    '255.255.0.0', '255.254.0.0', '255.252.0.0', '255.248.0.0',
    '255.240.0.0', '255.224.0.0', '255.192.0.0', '255.128.0.0',
    '255.0.0.0', '254.0.0.0', '252.0.0.0', '248.0.0.0',
    '240.0.0.0', '224.0.0.0', '192.0.0.0', '128.0.0.0'
  ];

  return validNetmasks.includes(netmask.trim());
};

export default function InterfaceManager() {
  const [interfaces, setInterfaces] = useState<NetworkInterface[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [configuring, setConfiguring] = useState<NetworkInterface | null>(null);
  const [configMode, setConfigMode] = useState<'dhcp' | 'static'>('dhcp');
  const [staticConfig, setStaticConfig] = useState({
    address: '',
    netmask: '255.255.255.0',
    gateway: '',
  });
  const [validationErrors, setValidationErrors] = useState<{
    address?: string;
    netmask?: string;
    gateway?: string;
  }>({});
  const [confirmAction, setConfirmAction] = useState<{
    title: string;
    message: string;
    action: () => void;
  } | null>(null);

  useEffect(() => {
    loadInterfaces();
    const interval = setInterval(loadInterfaces, 5000); // Refresh every 5s
    return () => clearInterval(interval);
  }, []);

  const loadInterfaces = async () => {
    try {
      const response = await networkApi.listInterfaces();
      if (response.success && response.data) {
        setInterfaces(response.data);
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load interfaces');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const toggleInterfaceState = (iface: NetworkInterface) => {
    const action = iface.isUp ? 'bring down' : 'bring up';
    const newState = iface.isUp ? 'down' : 'up';

    setConfirmAction({
      title: `${action.charAt(0).toUpperCase() + action.slice(1)} ${iface.name}?`,
      message: `Are you sure you want to ${action} interface ${iface.name}? This may interrupt network connectivity.`,
      action: async () => {
        try {
          const response = await networkApi.setInterfaceState(iface.name, newState);
          if (response.success) {
            loadInterfaces();
            setError('');
          } else {
            setError(response.error?.message || 'Failed to change interface state');
          }
        } catch (err) {
          setError(getErrorMessage(err));
        } finally {
          setConfirmAction(null);
        }
      },
    });
  };

  const openConfigModal = (iface: NetworkInterface) => {
    setConfiguring(iface);
    setConfigMode('dhcp');
    setStaticConfig({
      address: iface.addresses[0]?.split('/')[0] || '',
      netmask: '255.255.255.0',
      gateway: '',
    });
    setValidationErrors({});
  };

  const validateStaticConfig = (): boolean => {
    const errors: { address?: string; netmask?: string; gateway?: string } = {};
    let isValid = true;

    // Validate IP address
    if (!staticConfig.address.trim()) {
      errors.address = 'IP address is required';
      isValid = false;
    } else if (!isValidIP(staticConfig.address)) {
      errors.address = 'Invalid IP address format (e.g., 192.168.1.100)';
      isValid = false;
    }

    // Validate netmask
    if (!staticConfig.netmask.trim()) {
      errors.netmask = 'Netmask is required';
      isValid = false;
    } else if (!isValidNetmask(staticConfig.netmask)) {
      errors.netmask = 'Invalid netmask format (e.g., 255.255.255.0)';
      isValid = false;
    }

    // Validate gateway (optional, but must be valid if provided)
    if (staticConfig.gateway.trim() && !isValidIP(staticConfig.gateway)) {
      errors.gateway = 'Invalid gateway IP address format';
      isValid = false;
    }

    setValidationErrors(errors);
    return isValid;
  };

  const handleConfigure = () => {
    if (!configuring) return;

    // Validate static config if in static mode
    if (configMode === 'static' && !validateStaticConfig()) {
      return;
    }

    // Show confirmation dialog
    const modeDesc = configMode === 'dhcp' ? 'DHCP (automatic)' : `Static IP (${staticConfig.address})`;
    setConfirmAction({
      title: `Configure ${configuring.name}?`,
      message: `Are you sure you want to configure ${configuring.name} to use ${modeDesc}? This may interrupt network connectivity.`,
      action: async () => {
        try {
          let response;
          if (configMode === 'dhcp') {
            response = await networkApi.configureInterface(configuring.name, 'dhcp');
          } else {
            response = await networkApi.configureInterface(
              configuring.name,
              'static',
              staticConfig.address.trim(),
              staticConfig.netmask.trim(),
              staticConfig.gateway.trim()
            );
          }

          if (response.success) {
            setConfiguring(null);
            setError('');
            loadInterfaces();
          } else {
            setError(response.error?.message || 'Failed to configure interface');
          }
        } catch (err) {
          setError(getErrorMessage(err));
        } finally {
          setConfirmAction(null);
        }
      },
    });
  };

  const getInterfaceTypeIcon = (type: string) => {
    switch (type) {
      case 'wireless':
        return '📶';
      case 'loopback':
        return '🔄';
      case 'bridge':
        return '🌉';
      case 'virtual':
        return '💻';
      default:
        return '🌐';
    }
  };

  const getInterfaceTypeColor = (type: string) => {
    switch (type) {
      case 'wireless':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      case 'loopback':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
      case 'bridge':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400';
      case 'virtual':
        return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400';
      default:
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Error Display */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Interfaces Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        {interfaces.map((iface) => (
          <Card key={iface.name} hoverable>
            <div className="p-6">
              {/* Header */}
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="text-3xl">{getInterfaceTypeIcon(iface.type)}</div>
                  <div>
                    <h3 className="font-bold text-lg text-gray-900 dark:text-gray-100">
                      {iface.name}
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {iface.hardwareAddr || 'No MAC'}
                    </p>
                  </div>
                </div>
                {/* Status Toggle */}
                <button
                  onClick={() => toggleInterfaceState(iface)}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    iface.isUp ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
                  }`}
                  title={iface.isUp ? 'Click to bring down' : 'Click to bring up'}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                      iface.isUp ? 'translate-x-6' : 'translate-x-1'
                    }`}
                  />
                </button>
              </div>

              {/* Type Badge */}
              <div className="mb-4">
                <span
                  className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getInterfaceTypeColor(
                    iface.type
                  )}`}
                >
                  {iface.type}
                </span>
              </div>

              {/* Details */}
              <div className="space-y-2 mb-4">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Status:</span>
                  <span
                    className={`font-medium ${
                      iface.isUp
                        ? 'text-green-600 dark:text-green-400'
                        : 'text-gray-600 dark:text-gray-400'
                    }`}
                  >
                    {iface.isUp ? 'UP' : 'DOWN'}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Speed:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {iface.speed}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">MTU:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {iface.mtu}
                  </span>
                </div>
              </div>

              {/* IP Addresses */}
              {iface.addresses.length > 0 && (
                <div className="mb-4 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                  <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-2">
                    IP Addresses:
                  </div>
                  {iface.addresses.map((addr, idx) => (
                    <div
                      key={idx}
                      className="text-sm font-mono text-gray-900 dark:text-gray-100"
                    >
                      {addr}
                    </div>
                  ))}
                </div>
              )}

              {/* Flags */}
              {iface.flags.length > 0 && (
                <div className="mb-4">
                  <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                    Flags:
                  </div>
                  <div className="flex flex-wrap gap-1">
                    {iface.flags.map((flag, idx) => (
                      <span
                        key={idx}
                        className="px-2 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded"
                      >
                        {flag}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {/* Configure Button */}
              {iface.type !== 'loopback' && (
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={() => openConfigModal(iface)}
                  className="w-full"
                >
                  ⚙️ Configure
                </Button>
              )}
            </div>
          </Card>
        ))}
      </div>

      {/* Configuration Modal */}
      <AnimatePresence>
        {configuring && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setConfiguring(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Configure {configuring.name}
              </h2>

              {/* Mode Selection */}
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Configuration Mode
                </label>
                <div className="flex gap-2">
                  <button
                    onClick={() => setConfigMode('dhcp')}
                    className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                      configMode === 'dhcp'
                        ? 'bg-macos-blue text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                    }`}
                  >
                    DHCP
                  </button>
                  <button
                    onClick={() => setConfigMode('static')}
                    className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                      configMode === 'static'
                        ? 'bg-macos-blue text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                    }`}
                  >
                    Static IP
                  </button>
                </div>
              </div>

              {/* Static IP Configuration */}
              {configMode === 'static' && (
                <div className="space-y-4">
                  <div>
                    <Input
                      label="IP Address"
                      value={staticConfig.address}
                      onChange={(e) => {
                        setStaticConfig({ ...staticConfig, address: e.target.value });
                        // Clear error when user types
                        if (validationErrors.address) {
                          setValidationErrors({ ...validationErrors, address: undefined });
                        }
                      }}
                      placeholder="192.168.1.100"
                      required
                    />
                    {validationErrors.address && (
                      <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                        {validationErrors.address}
                      </p>
                    )}
                  </div>
                  <div>
                    <Input
                      label="Netmask"
                      value={staticConfig.netmask}
                      onChange={(e) => {
                        setStaticConfig({ ...staticConfig, netmask: e.target.value });
                        // Clear error when user types
                        if (validationErrors.netmask) {
                          setValidationErrors({ ...validationErrors, netmask: undefined });
                        }
                      }}
                      placeholder="255.255.255.0"
                      required
                    />
                    {validationErrors.netmask && (
                      <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                        {validationErrors.netmask}
                      </p>
                    )}
                  </div>
                  <div>
                    <Input
                      label="Gateway (Optional)"
                      value={staticConfig.gateway}
                      onChange={(e) => {
                        setStaticConfig({ ...staticConfig, gateway: e.target.value });
                        // Clear error when user types
                        if (validationErrors.gateway) {
                          setValidationErrors({ ...validationErrors, gateway: undefined });
                        }
                      }}
                      placeholder="192.168.1.1"
                    />
                    {validationErrors.gateway && (
                      <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                        {validationErrors.gateway}
                      </p>
                    )}
                  </div>
                </div>
              )}

              {configMode === 'dhcp' && (
                <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg text-blue-700 dark:text-blue-400 text-sm">
                  <p className="font-medium mb-1">DHCP Mode</p>
                  <p>IP address will be automatically assigned by the DHCP server.</p>
                </div>
              )}

              {/* Actions */}
              <div className="flex gap-3 mt-6">
                <Button
                  variant="secondary"
                  onClick={() => setConfiguring(null)}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button onClick={handleConfigure} className="flex-1">
                  Apply Configuration
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Confirmation Dialog */}
      <AnimatePresence>
        {confirmAction && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setConfirmAction(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                {confirmAction.title}
              </h2>
              <p className="text-gray-700 dark:text-gray-300 mb-6">
                {confirmAction.message}
              </p>
              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setConfirmAction(null)}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  onClick={confirmAction.action}
                  className="flex-1 bg-red-600 hover:bg-red-700"
                >
                  Confirm
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
