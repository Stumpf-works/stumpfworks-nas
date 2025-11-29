import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { vmsApi, type VM } from '../../api/vms';
import { Play, Square, Trash2, Plus, Monitor, RefreshCw } from 'lucide-react';
import toast from 'react-hot-toast';

export function VMManager() {
  const [vms, setVMs] = useState<VM[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionInProgress, setActionInProgress] = useState<string | null>(null);

  useEffect(() => {
    loadVMs();
    const interval = setInterval(loadVMs, 5000);
    return () => clearInterval(interval);
  }, []);

  const loadVMs = async () => {
    try {
      const response = await vmsApi.listVMs();
      if (response.success && response.data) {
        setVMs(response.data);
      }
    } catch (error) {
      console.error('Failed to load VMs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleStart = async (vm: VM) => {
    setActionInProgress(vm.uuid);
    try {
      const response = await vmsApi.startVM(vm.uuid);
      if (response.success) {
        toast.success(`VM ${vm.name} started`);
        await loadVMs();
      } else {
        toast.error(response.error?.message || 'Failed to start VM');
      }
    } catch (error: any) {
      toast.error(error.message || 'Failed to start VM');
    } finally {
      setActionInProgress(null);
    }
  };

  const handleStop = async (vm: VM, force: boolean = false) => {
    setActionInProgress(vm.uuid);
    try {
      const response = await vmsApi.stopVM(vm.uuid, force);
      if (response.success) {
        toast.success(`VM ${vm.name} stopped`);
        await loadVMs();
      } else {
        toast.error(response.error?.message || 'Failed to stop VM');
      }
    } catch (error: any) {
      toast.error(error.message || 'Failed to stop VM');
    } finally {
      setActionInProgress(null);
    }
  };

  const handleDelete = async (vm: VM) => {
    if (!confirm(`Delete VM "${vm.name}"? This cannot be undone.`)) return;

    setActionInProgress(vm.uuid);
    try {
      const response = await vmsApi.deleteVM(vm.uuid, true);
      if (response.success) {
        toast.success(`VM ${vm.name} deleted`);
        await loadVMs();
      } else {
        toast.error(response.error?.message || 'Failed to delete VM');
      }
    } catch (error: any) {
      toast.error(error.message || 'Failed to delete VM');
    } finally {
      setActionInProgress(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <RefreshCw className="w-8 h-8 animate-spin text-macos-blue mx-auto mb-4" />
          <p className="text-gray-600 dark:text-gray-400">Loading VMs...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-50">
      {/* Header */}
      <div className="px-6 pt-8 pb-6 bg-white dark:bg-macos-dark-100">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">VM Manager</h1>
            <p className="text-gray-500 dark:text-gray-400 mt-1">Manage virtual machines with KVM/QEMU</p>
          </div>
          <button className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2">
            <Plus className="w-5 h-5" />
            Create VM
          </button>
        </div>
      </div>

      {/* VM List */}
      <div className="flex-1 overflow-auto px-6 py-4">
        {vms.length === 0 ? (
          <div className="text-center py-12">
            <Monitor className="w-16 h-16 text-gray-300 dark:text-gray-600 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">No VMs Found</h3>
            <p className="text-gray-500 dark:text-gray-400">Create your first virtual machine to get started</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
            {vms.map((vm) => {
              const isRunning = vm.state === 'running';
              const isInProgress = actionInProgress === vm.uuid;

              return (
                <motion.div
                  key={vm.uuid}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="bg-gray-50 dark:bg-macos-dark-100 rounded-xl p-6 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{vm.name}</h3>
                      <span className={`inline-block px-2 py-1 text-xs font-medium rounded mt-2 ${
                        isRunning
                          ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                          : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400'
                      }`}>
                        {vm.state}
                      </span>
                    </div>
                  </div>

                  <div className="space-y-2 text-sm text-gray-600 dark:text-gray-400 mb-4">
                    <div className="flex justify-between">
                      <span>CPUs:</span>
                      <span className="font-medium">{vm.vcpus}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Memory:</span>
                      <span className="font-medium">{(vm.memory / 1024).toFixed(1)} GB</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Disk:</span>
                      <span className="font-medium">{(vm.disk_size / 1024).toFixed(1)} GB</span>
                    </div>
                  </div>

                  <div className="flex gap-2">
                    {!isRunning ? (
                      <button
                        onClick={() => handleStart(vm)}
                        disabled={isInProgress}
                        className="flex-1 px-3 py-2 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                      >
                        <Play className="w-4 h-4" />
                        Start
                      </button>
                    ) : (
                      <button
                        onClick={() => handleStop(vm)}
                        disabled={isInProgress}
                        className="flex-1 px-3 py-2 bg-yellow-100 text-yellow-700 rounded-lg hover:bg-yellow-200 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                      >
                        <Square className="w-4 h-4" />
                        Stop
                      </button>
                    )}
                    <button
                      onClick={() => handleDelete(vm)}
                      disabled={isInProgress || isRunning}
                      className="px-3 py-2 bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition-colors disabled:opacity-50"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </motion.div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
