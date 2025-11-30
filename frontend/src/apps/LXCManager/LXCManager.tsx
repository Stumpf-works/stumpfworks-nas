import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { lxcApi, type Container } from '@/api/lxc';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import {
  Box,
  Play,
  Square,
  Trash2,
  Plus,
  RefreshCw,
  AlertCircle,
  Network,
  Activity,
  Power
} from 'lucide-react';

export function LXCManager() {
  const [containers, setContainers] = useState<Container[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  useEffect(() => {
    loadContainers();
  }, []);

  const loadContainers = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await lxcApi.listContainers();
      if (response.success && response.data) {
        setContainers(response.data);
      } else {
        setError(response.error?.message || 'Failed to load containers');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleStart = async (name: string) => {
    try {
      setActionLoading(name);
      const response = await lxcApi.startContainer(name);
      if (response.success) {
        await loadContainers();
      } else {
        setError(response.error?.message || 'Failed to start container');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleStop = async (name: string) => {
    try {
      setActionLoading(name);
      const response = await lxcApi.stopContainer(name);
      if (response.success) {
        await loadContainers();
      } else {
        setError(response.error?.message || 'Failed to stop container');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (name: string) => {
    if (!confirm(`Are you sure you want to delete container "${name}"? This cannot be undone.`)) {
      return;
    }

    try {
      setActionLoading(name);
      const response = await lxcApi.deleteContainer(name);
      if (response.success) {
        await loadContainers();
      } else {
        setError(response.error?.message || 'Failed to delete container');
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
      case 'frozen': return 'text-blue-600 dark:text-blue-400';
      default: return 'text-gray-600 dark:text-gray-400';
    }
  };

  const getStateBadge = (state: string) => {
    const baseClasses = 'px-2 py-1 rounded-full text-xs font-semibold';
    switch (state.toLowerCase()) {
      case 'running': return `${baseClasses} bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300`;
      case 'stopped': return `${baseClasses} bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300`;
      case 'frozen': return `${baseClasses} bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300`;
      default: return `${baseClasses} bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300`;
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-full bg-gray-50 dark:bg-macos-dark-50">
        <RefreshCw className="w-12 h-12 text-macos-blue animate-spin" />
        <p className="mt-4 text-gray-600 dark:text-gray-400">Loading containers...</p>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Box className="w-8 h-8 text-macos-blue" />
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">LXC Manager</h1>
              <p className="text-sm text-gray-600 dark:text-gray-400">Manage Linux containers</p>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              onClick={loadContainers}
              className="px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center gap-2"
            >
              <RefreshCw className="w-4 h-4" />
              Refresh
            </button>
            <button className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2">
              <Plus className="w-4 h-4" />
              Create Container
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

      {/* Containers List */}
      <div className="flex-1 overflow-auto p-6">
        {containers.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full">
            <Box className="w-16 h-16 text-gray-400 dark:text-gray-600 mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">No Containers</h3>
            <p className="text-gray-600 dark:text-gray-400 mb-6">Get started by creating your first container</p>
            <button className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2">
              <Plus className="w-4 h-4" />
              Create Container
            </button>
          </div>
        ) : (
          <div className="grid gap-4">
            {containers.map((container) => (
              <motion.div
                key={container.name}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
              >
                <Card>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className={`p-3 rounded-lg ${container.state.toLowerCase() === 'running' ? 'bg-green-100 dark:bg-green-900/30' : 'bg-gray-100 dark:bg-gray-800/30'}`}>
                        <Box className={`w-6 h-6 ${getStateColor(container.state)}`} />
                      </div>
                      <div>
                        <div className="flex items-center gap-3">
                          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{container.name}</h3>
                          <span className={getStateBadge(container.state)}>{container.state}</span>
                        </div>
                        <div className="flex items-center gap-4 mt-2 text-sm text-gray-600 dark:text-gray-400">
                          {container.ipv4 && (
                            <div className="flex items-center gap-1">
                              <Network className="w-4 h-4" />
                              {container.ipv4}
                            </div>
                          )}
                          {container.state.toLowerCase() === 'running' && container.memory && (
                            <div className="flex items-center gap-1">
                              <Activity className="w-4 h-4" />
                              {container.memory} / {container.memory_limit}
                            </div>
                          )}
                          {container.state.toLowerCase() === 'running' && container.cpu_usage && (
                            <div className="flex items-center gap-1">
                              <Activity className="w-4 h-4" />
                              CPU: {container.cpu_usage}
                            </div>
                          )}
                          {container.autostart && (
                            <div className="flex items-center gap-1 text-macos-blue">
                              <Power className="w-4 h-4" />
                              Autostart
                            </div>
                          )}
                        </div>
                        {container.state.toLowerCase() === 'running' && container.pid > 0 && (
                          <div className="flex items-center gap-1 mt-1 text-xs text-gray-500 dark:text-gray-500">
                            PID: {container.pid}
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="flex gap-2">
                      {container.state.toLowerCase() === 'running' ? (
                        <button
                          onClick={() => handleStop(container.name)}
                          disabled={actionLoading === container.name}
                          className="p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors disabled:opacity-50"
                          title="Stop Container"
                        >
                          {actionLoading === container.name ? (
                            <RefreshCw className="w-5 h-5 animate-spin" />
                          ) : (
                            <Square className="w-5 h-5" />
                          )}
                        </button>
                      ) : (
                        <button
                          onClick={() => handleStart(container.name)}
                          disabled={actionLoading === container.name}
                          className="p-2 text-green-600 dark:text-green-400 hover:bg-green-50 dark:hover:bg-green-900/20 rounded-lg transition-colors disabled:opacity-50"
                          title="Start Container"
                        >
                          {actionLoading === container.name ? (
                            <RefreshCw className="w-5 h-5 animate-spin" />
                          ) : (
                            <Play className="w-5 h-5" />
                          )}
                        </button>
                      )}
                      <button
                        onClick={() => handleDelete(container.name)}
                        disabled={actionLoading === container.name}
                        className="p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors disabled:opacity-50"
                        title="Delete Container"
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
    </div>
  );
}
