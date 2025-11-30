import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { vmsApi, type VM } from '@/api/vms';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import {
  Server,
  Play,
  Square,
  Trash2,
  Plus,
  RefreshCw,
  AlertCircle,
  Cpu,
  HardDrive,
  Calendar,
  Power,
  Monitor
} from 'lucide-react';
import { CreateVMModal } from './components/CreateVMModal';
import { VNCModal } from './components/VNCModal';

export function VMManager() {
  const [vms, setVms] = useState<VM[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [vncVM, setVncVM] = useState<{ id: string; name: string } | null>(null);

  useEffect(() => {
    loadVMs();
  }, []);

  const loadVMs = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await vmsApi.listVMs();
      if (response.success && response.data) {
        setVms(response.data);
      } else {
        setError(response.error?.message || 'Failed to load VMs');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleStart = async (vmId: string) => {
    try {
      setActionLoading(vmId);
      const response = await vmsApi.startVM(vmId);
      if (response.success) {
        await loadVMs();
      } else {
        setError(response.error?.message || 'Failed to start VM');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleStop = async (vmId: string) => {
    try {
      setActionLoading(vmId);
      const response = await vmsApi.stopVM(vmId);
      if (response.success) {
        await loadVMs();
      } else {
        setError(response.error?.message || 'Failed to stop VM');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (vmId: string, name: string) => {
    if (!confirm(`Are you sure you want to delete VM "${name}"? This cannot be undone.`)) {
      return;
    }

    try {
      setActionLoading(vmId);
      const response = await vmsApi.deleteVM(vmId, true);
      if (response.success) {
        await loadVMs();
      } else {
        setError(response.error?.message || 'Failed to delete VM');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const getStateColor = (state: string) => {
    switch (state.toLowerCase()) {
      case 'running': return 'text-green-600 dark:text-green-400';
      case 'stopped': return 'text-gray-600 dark:text-gray-400';
      case 'paused': return 'text-yellow-600 dark:text-yellow-400';
      default: return 'text-gray-600 dark:text-gray-400';
    }
  };

  const getStateBadge = (state: string) => {
    const baseClasses = 'px-2 py-1 rounded-full text-xs font-semibold';
    switch (state.toLowerCase()) {
      case 'running': return `${baseClasses} bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300`;
      case 'stopped': return `${baseClasses} bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300`;
      case 'paused': return `${baseClasses} bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300`;
      default: return `${baseClasses} bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300`;
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-full bg-gray-50 dark:bg-macos-dark-50">
        <RefreshCw className="w-12 h-12 text-macos-blue animate-spin" />
        <p className="mt-4 text-gray-600 dark:text-gray-400">Loading VMs...</p>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Server className="w-8 h-8 text-macos-blue" />
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">VM Manager</h1>
              <p className="text-sm text-gray-600 dark:text-gray-400">Manage virtual machines</p>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              onClick={loadVMs}
              className="px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center gap-2"
            >
              <RefreshCw className="w-4 h-4" />
              Refresh
            </button>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
            >
              <Plus className="w-4 h-4" />
              Create VM
            </button>
          </div>
        </div>
      </div>

      {/* Error Display */}
      {error && (
        <div className="mx-6 mt-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-start gap-3">
          <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <h3 className="font-semibold text-red-900 dark:text-red-200">Error</h3>
            <p className="text-sm text-red-700 dark:text-red-300 mt-1">{error}</p>
          </div>
          <button onClick={() => setError('')} className="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200">
            ✕
          </button>
        </div>
      )}

      {/* VMs List */}
      <div className="flex-1 overflow-auto p-6">
        {vms.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full">
            <Server className="w-16 h-16 text-gray-400 dark:text-gray-600 mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">No Virtual Machines</h3>
            <p className="text-gray-600 dark:text-gray-400 mb-6">Get started by creating your first VM</p>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
            >
              <Plus className="w-4 h-4" />
              Create VM
            </button>
          </div>
        ) : (
          <div className="grid gap-4">
            {vms.map((vm) => (
              <motion.div
                key={vm.uuid}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
              >
                <Card>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className={`p-3 rounded-lg ${vm.state.toLowerCase() === 'running' ? 'bg-green-100 dark:bg-green-900/30' : 'bg-gray-100 dark:bg-gray-800/30'}`}>
                        <Server className={`w-6 h-6 ${getStateColor(vm.state)}`} />
                      </div>
                      <div>
                        <div className="flex items-center gap-3">
                          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{vm.name}</h3>
                          <span className={getStateBadge(vm.state)}>{vm.state}</span>
                        </div>
                        <div className="flex items-center gap-4 mt-2 text-sm text-gray-600 dark:text-gray-400">
                          <div className="flex items-center gap-1">
                            <Cpu className="w-4 h-4" />
                            {vm.vcpus} vCPUs
                          </div>
                          <div className="flex items-center gap-1">
                            <HardDrive className="w-4 h-4" />
                            {(vm.memory / 1024).toFixed(1)} GB RAM
                          </div>
                          <div className="flex items-center gap-1">
                            <HardDrive className="w-4 h-4" />
                            {vm.disk_size} GB Disk
                          </div>
                          {vm.autostart && (
                            <div className="flex items-center gap-1 text-macos-blue">
                              <Power className="w-4 h-4" />
                              Autostart
                            </div>
                          )}
                        </div>
                        <div className="flex items-center gap-1 mt-1 text-xs text-gray-500 dark:text-gray-500">
                          <Calendar className="w-3 h-3" />
                          Created {new Date(vm.created_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      {vm.state.toLowerCase() === 'running' ? (
                        <button
                          onClick={() => handleStop(vm.uuid)}
                          disabled={actionLoading === vm.uuid}
                          className="p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors disabled:opacity-50"
                          title="Stop VM"
                        >
                          {actionLoading === vm.uuid ? (
                            <RefreshCw className="w-5 h-5 animate-spin" />
                          ) : (
                            <Square className="w-5 h-5" />
                          )}
                        </button>
                      ) : (
                        <button
                          onClick={() => handleStart(vm.uuid)}
                          disabled={actionLoading === vm.uuid}
                          className="p-2 text-green-600 dark:text-green-400 hover:bg-green-50 dark:hover:bg-green-900/20 rounded-lg transition-colors disabled:opacity-50"
                          title="Start VM"
                        >
                          {actionLoading === vm.uuid ? (
                            <RefreshCw className="w-5 h-5 animate-spin" />
                          ) : (
                            <Play className="w-5 h-5" />
                          )}
                        </button>
                      )}
                      {vm.state.toLowerCase() === 'running' && (
                        <button
                          onClick={() => setVncVM({ id: vm.uuid, name: vm.name })}
                          className="p-2 text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-colors"
                          title="VNC Console"
                        >
                          <Monitor className="w-5 h-5" />
                        </button>
                      )}
                      <button
                        onClick={() => handleDelete(vm.uuid, vm.name)}
                        disabled={actionLoading === vm.uuid}
                        className="p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors disabled:opacity-50"
                        title="Delete VM"
                      >
                        <Trash2 className="w-5 h-5" />
                      </button>
                    </div>
                  </div>
                </Card>
              </motion.div>
            ))}
          </div>
        )}
      </div>

      {/* Create VM Modal */}
      <CreateVMModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={loadVMs}
      />

      {/* VNC Modal */}
      {vncVM && (
        <VNCModal
          isOpen={!!vncVM}
          onClose={() => setVncVM(null)}
          vmId={vncVM.id}
          vmName={vncVM.name}
        />
      )}
    </div>
  );
}
