import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { timeMachineApi, TimeMachineConfig, TimeMachineDevice } from '@/api/timemachine';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import AddDeviceModal from './components/AddDeviceModal';
import EditDeviceModal from './components/EditDeviceModal';
import SettingsModal from './components/SettingsModal';
import {
  HardDrive,
  Plus,
  Settings,
  RefreshCw,
  Trash2,
  Edit,
  Power,
  PowerOff,
  Clock,
  Database
} from 'lucide-react';

export function TimeMachine() {
  const [config, setConfig] = useState<TimeMachineConfig | null>(null);
  const [devices, setDevices] = useState<TimeMachineDevice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showAddDeviceModal, setShowAddDeviceModal] = useState(false);
  const [showConfigModal, setShowConfigModal] = useState(false);
  const [editingDevice, setEditingDevice] = useState<TimeMachineDevice | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      const [configResponse, devicesResponse] = await Promise.all([
        timeMachineApi.getConfig(),
        timeMachineApi.listDevices(),
      ]);

      if (configResponse.success && configResponse.data) {
        setConfig(configResponse.data);
      }

      if (devicesResponse.success && devicesResponse.data) {
        setDevices(devicesResponse.data);
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleToggleService = async () => {
    if (!config) return;

    try {
      if (config.enabled) {
        await timeMachineApi.disable();
      } else {
        await timeMachineApi.enable();
      }
      await loadData();
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDeleteDevice = async (id: number) => {
    if (!confirm('Are you sure you want to delete this device? Backup data will be preserved.')) {
      return;
    }

    try {
      await timeMachineApi.deleteDevice(id);
      await loadData();
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleUpdateUsage = async () => {
    try {
      setLoading(true);
      await timeMachineApi.updateAllDeviceUsages();
      await loadData();
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const formatBytes = (gb: number) => {
    if (gb >= 1000) {
      return `${(gb / 1000).toFixed(2)} TB`;
    }
    return `${gb.toFixed(2)} GB`;
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return 'Never';
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const getUsagePercentage = (used: number, quota: number) => {
    if (quota === 0) return 0;
    return Math.min((used / quota) * 100, 100);
  };

  if (loading && !config) {
    return (
      <div className="flex items-center justify-center h-full bg-gray-50 dark:bg-macos-dark-50">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue"></div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-3">
              <HardDrive className="w-8 h-8" />
              Time Machine Server
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mt-1">
              macOS Time Machine backup server for network backups
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button
              variant="secondary"
              onClick={() => setShowConfigModal(true)}
            >
              <Settings className="w-4 h-4 mr-2" />
              Settings
            </Button>
            <Button
              variant={config?.enabled ? 'danger' : 'primary'}
              onClick={handleToggleService}
            >
              {config?.enabled ? (
                <>
                  <PowerOff className="w-4 h-4 mr-2" />
                  Disable
                </>
              ) : (
                <>
                  <Power className="w-4 h-4 mr-2" />
                  Enable
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Status Banner */}
        {config && (
          <div className={`mt-4 p-4 rounded-lg ${
            config.enabled
              ? 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800'
              : 'bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700'
          }`}>
            <div className="flex items-center gap-2">
              <div className={`w-3 h-3 rounded-full ${
                config.enabled ? 'bg-green-500' : 'bg-gray-400'
              }`}></div>
              <span className={`text-sm font-medium ${
                config.enabled ? 'text-green-700 dark:text-green-300' : 'text-gray-600 dark:text-gray-400'
              }`}>
                Time Machine is {config.enabled ? 'enabled and accepting backups' : 'disabled'}
              </span>
            </div>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="m-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-red-600 dark:text-red-400 text-sm">{error}</p>
        </div>
      )}

      {/* Devices Section */}
      <div className="flex-1 overflow-auto p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
            Registered Devices
          </h2>
          <div className="flex gap-3">
            <Button
              variant="secondary"
              onClick={handleUpdateUsage}
              disabled={loading}
            >
              <RefreshCw className={`w-4 h-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
              Update Usage
            </Button>
            <Button
              variant="primary"
              onClick={() => setShowAddDeviceModal(true)}
            >
              <Plus className="w-4 h-4 mr-2" />
              Add Device
            </Button>
          </div>
        </div>

        {/* Device Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {devices.map((device) => {
            const usagePercentage = getUsagePercentage(device.used_gb, device.quota_gb);
            const isNearLimit = usagePercentage >= 80;

            return (
              <motion.div
                key={device.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.2 }}
              >
                <Card className={`${!device.enabled ? 'opacity-60' : ''}`}>
                  <div className="p-6">
                    {/* Device Header */}
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex items-center gap-3">
                        <div className={`p-3 rounded-lg ${
                          device.enabled
                            ? 'bg-blue-50 dark:bg-blue-900/20'
                            : 'bg-gray-100 dark:bg-gray-800'
                        }`}>
                          <HardDrive className={`w-6 h-6 ${
                            device.enabled
                              ? 'text-blue-600 dark:text-blue-400'
                              : 'text-gray-400'
                          }`} />
                        </div>
                        <div>
                          <h3 className="font-semibold text-gray-900 dark:text-white">
                            {device.device_name}
                          </h3>
                          {device.model_id && (
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                              {device.model_id}
                            </p>
                          )}
                        </div>
                      </div>
                      <div className={`w-2 h-2 rounded-full ${
                        device.enabled ? 'bg-green-500' : 'bg-gray-400'
                      }`}></div>
                    </div>

                    {/* Storage Usage */}
                    <div className="mb-4">
                      <div className="flex justify-between text-sm mb-2">
                        <span className="text-gray-600 dark:text-gray-400">
                          Storage Used
                        </span>
                        <span className={`font-medium ${
                          isNearLimit
                            ? 'text-orange-600 dark:text-orange-400'
                            : 'text-gray-900 dark:text-white'
                        }`}>
                          {formatBytes(device.used_gb)} / {device.quota_gb > 0 ? formatBytes(device.quota_gb) : 'Unlimited'}
                        </span>
                      </div>
                      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full transition-all ${
                            isNearLimit
                              ? 'bg-orange-500'
                              : 'bg-blue-500'
                          }`}
                          style={{ width: `${usagePercentage}%` }}
                        ></div>
                      </div>
                    </div>

                    {/* Last Backup */}
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-4">
                      <Clock className="w-4 h-4" />
                      <span>Last backup: {formatDate(device.last_backup)}</span>
                    </div>

                    {/* Actions */}
                    <div className="flex gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => setEditingDevice(device)}
                        className="flex-1"
                      >
                        <Edit className="w-4 h-4 mr-1" />
                        Edit
                      </Button>
                      <Button
                        variant="danger"
                        size="sm"
                        onClick={() => handleDeleteDevice(device.id)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </Card>
              </motion.div>
            );
          })}

          {/* Empty State */}
          {devices.length === 0 && !loading && (
            <div className="col-span-full flex flex-col items-center justify-center py-12 text-center">
              <Database className="w-16 h-16 text-gray-400 mb-4" />
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
                No Devices Registered
              </h3>
              <p className="text-gray-600 dark:text-gray-400 mb-6 max-w-md">
                Add your macOS devices to start backing up to this Time Machine server
              </p>
              <Button
                variant="primary"
                onClick={() => setShowAddDeviceModal(true)}
              >
                <Plus className="w-4 h-4 mr-2" />
                Add Your First Device
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* Modals */}
      <AddDeviceModal
        isOpen={showAddDeviceModal}
        onClose={() => setShowAddDeviceModal(false)}
        onSuccess={loadData}
        defaultQuota={config?.default_quota_gb || 500}
      />

      <EditDeviceModal
        device={editingDevice}
        isOpen={!!editingDevice}
        onClose={() => setEditingDevice(null)}
        onSuccess={loadData}
      />

      <SettingsModal
        config={config}
        isOpen={showConfigModal}
        onClose={() => setShowConfigModal(false)}
        onSuccess={loadData}
      />
    </div>
  );
}
