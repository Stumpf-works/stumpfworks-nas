import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
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
  Power,
  ChevronRight,
  Server
} from 'lucide-react';
import { CreateContainerModal } from './components/CreateContainerModal';
import { ContainerDetailView } from './components/ContainerDetailView';

export function LXCManager() {
  const [containers, setContainers] = useState<Container[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [selectedContainer, setSelectedContainer] = useState<Container | null>(null);

  useEffect(() => {
    loadContainers();
    const interval = setInterval(loadContainers, 5000);
    return () => clearInterval(interval);
  }, []);

  const loadContainers = async () => {
    try {
      if (!loading) setLoading(false);
      setError('');
      const response = await lxcApi.listContainers();
      if (response.success && response.data) {
        setContainers(response.data);
        // Update selected container if it's still in the list
        if (selectedContainer) {
          const updated = response.data.find(c => c.name === selectedContainer.name);
          if (updated) {
            setSelectedContainer(updated);
          }
        }
      } else {
        setError(response.error?.message || 'Failed to load containers');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleAction = async (action: 'start' | 'stop' | 'restart' | 'delete', name: string) => {
    try {
      setActionLoading(name);

      if (action === 'delete') {
        if (!confirm(`Are you sure you want to delete container "${name}"? This cannot be undone.`)) {
          return;
        }
      }

      let response;
      switch (action) {
        case 'start':
          response = await lxcApi.startContainer(name);
          break;
        case 'stop':
          response = await lxcApi.stopContainer(name);
          break;
        case 'restart':
          await lxcApi.stopContainer(name);
          await new Promise(resolve => setTimeout(resolve, 2000));
          response = await lxcApi.startContainer(name);
          break;
        case 'delete':
          response = await lxcApi.deleteContainer(name);
          setSelectedContainer(null);
          break;
      }

      if (response?.success) {
        await loadContainers();
      } else {
        setError(response?.error?.message || `Failed to ${action} container`);
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
    const baseClasses = 'px-2 py-1 rounded-full text-xs font-semibold flex items-center gap-1.5';
    const dotClasses = 'w-1.5 h-1.5 rounded-full';

    let colorClasses = '';
    let dotColor = '';

    switch (state.toLowerCase()) {
      case 'running':
        colorClasses = 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300';
        dotColor = 'bg-green-500';
        break;
      case 'stopped':
        colorClasses = 'bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300';
        dotColor = 'bg-gray-500';
        break;
      case 'frozen':
        colorClasses = 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300';
        dotColor = 'bg-blue-500';
        break;
      default:
        colorClasses = 'bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300';
        dotColor = 'bg-gray-500';
    }

    return (
      <span className={`${baseClasses} ${colorClasses}`}>
        <span className={`${dotClasses} ${dotColor}`} />
        {state}
      </span>
    );
  };

  if (loading && containers.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full bg-gray-50 dark:bg-macos-dark-50">
        <RefreshCw className="w-12 h-12 text-macos-blue animate-spin" />
        <p className="mt-4 text-gray-600 dark:text-gray-400">Loading containers...</p>
      </div>
    );
  }

  return (
    <div className="h-full flex bg-gray-50 dark:bg-macos-dark-50">
      {/* Left Panel - Container List */}
      <div className="w-80 flex flex-col bg-white dark:bg-macos-dark-100 border-r border-gray-200 dark:border-gray-700">
        {/* Header */}
        <div className="p-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-3 mb-4">
            <Box className="w-6 h-6 text-macos-blue" />
            <div className="flex-1">
              <h1 className="text-lg font-bold text-gray-900 dark:text-white">LXC Manager</h1>
              <p className="text-xs text-gray-600 dark:text-gray-400">{containers.length} container{containers.length !== 1 ? 's' : ''}</p>
            </div>
          </div>

          <div className="flex gap-2">
            <button
              onClick={loadContainers}
              className="flex-1 px-3 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center justify-center gap-2 text-sm"
            >
              <RefreshCw className="w-4 h-4" />
              Refresh
            </button>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="flex-1 px-3 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center justify-center gap-2 text-sm"
            >
              <Plus className="w-4 h-4" />
              Create
            </button>
          </div>
        </div>

        {/* Error Display */}
        {error && (
          <div className="mx-4 mt-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <div className="flex items-start gap-2">
              <AlertCircle className="w-4 h-4 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="text-xs text-red-700 dark:text-red-300">{error}</p>
              </div>
              <button onClick={() => setError('')} className="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200 text-xs">
                ✕
              </button>
            </div>
          </div>
        )}

        {/* Container List */}
        <div className="flex-1 overflow-auto p-4">
          {containers.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-center px-4">
              <Server className="w-12 h-12 text-gray-400 dark:text-gray-600 mb-3" />
              <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-1">No Containers</h3>
              <p className="text-xs text-gray-600 dark:text-gray-400 mb-4">Create your first container to get started</p>
              <button
                onClick={() => setIsCreateModalOpen(true)}
                className="px-3 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2 text-sm"
              >
                <Plus className="w-4 h-4" />
                Create Container
              </button>
            </div>
          ) : (
            <div className="space-y-2">
              {containers.map((container) => (
                <motion.button
                  key={container.name}
                  onClick={() => setSelectedContainer(container)}
                  className={`
                    w-full text-left p-3 rounded-lg transition-all
                    ${selectedContainer?.name === container.name
                      ? 'bg-macos-blue/10 border-macos-blue dark:bg-macos-blue/20 border'
                      : 'bg-gray-50 dark:bg-macos-dark-50 hover:bg-gray-100 dark:hover:bg-gray-800 border border-transparent'
                    }
                  `}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <div className={`p-1.5 rounded ${container.state.toLowerCase() === 'running' ? 'bg-green-100 dark:bg-green-900/30' : 'bg-gray-100 dark:bg-gray-800/30'}`}>
                        <Box className={`w-4 h-4 ${getStateColor(container.state)}`} />
                      </div>
                      <span className="font-medium text-gray-900 dark:text-white text-sm truncate">{container.name}</span>
                    </div>
                    <ChevronRight className="w-4 h-4 text-gray-400 flex-shrink-0" />
                  </div>
                  <div className="flex items-center gap-2 ml-7">
                    {getStateBadge(container.state)}
                    {container.ipv4 && (
                      <div className="flex items-center gap-1 text-xs text-gray-600 dark:text-gray-400">
                        <Network className="w-3 h-3" />
                        {container.ipv4}
                      </div>
                    )}
                  </div>
                </motion.button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Right Panel - Container Details */}
      <div className="flex-1 flex flex-col">
        <AnimatePresence mode="wait">
          {selectedContainer ? (
            <motion.div
              key={selectedContainer.name}
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              className="h-full"
            >
              <ContainerDetailView
                container={selectedContainer}
                onAction={handleAction}
                onClose={() => setSelectedContainer(null)}
                loading={actionLoading === selectedContainer.name}
              />
            </motion.div>
          ) : (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="h-full flex flex-col items-center justify-center bg-white dark:bg-macos-dark-100"
            >
              <Box className="w-24 h-24 text-gray-300 dark:text-gray-700 mb-4" />
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Select a Container</h2>
              <p className="text-gray-600 dark:text-gray-400 text-center max-w-md px-4">
                Choose a container from the list to view details, manage settings, and access the console
              </p>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Create Container Modal */}
      <CreateContainerModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={loadContainers}
      />
    </div>
  );
}
