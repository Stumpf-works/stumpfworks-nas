import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { dockerApi, DockerContainer } from '@/api/docker';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';

export default function ContainerManager() {
  const [containers, setContainers] = useState<DockerContainer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAll, setShowAll] = useState(true);
  const [logsModal, setLogsModal] = useState<{ container: DockerContainer; logs: string } | null>(
    null
  );
  const [deleteModal, setDeleteModal] = useState<DockerContainer | null>(null);
  const [resourcesModal, setResourcesModal] = useState<DockerContainer | null>(null);
  const [execModal, setExecModal] = useState<DockerContainer | null>(null);
  const [execCommand, setExecCommand] = useState('ls -la');
  const [execOutput, setExecOutput] = useState('');
  const [resourceLimits, setResourceLimits] = useState({
    memory: '',
    cpuShares: '',
  });
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  useEffect(() => {
    loadContainers();
    const interval = setInterval(loadContainers, 5000); // Refresh every 5s
    return () => clearInterval(interval);
  }, [showAll]);

  const loadContainers = async () => {
    try {
      const response = await dockerApi.listContainers(showAll);
      if (response.success && response.data) {
        setContainers(response.data);
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load containers');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleStart = async (container: DockerContainer) => {
    setActionLoading(container.id);
    try {
      const response = await dockerApi.startContainer(container.id);
      if (response.success) {
        loadContainers();
      } else {
        alert(response.error?.message || 'Failed to start container');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleStop = async (container: DockerContainer) => {
    setActionLoading(container.id);
    try {
      const response = await dockerApi.stopContainer(container.id);
      if (response.success) {
        loadContainers();
      } else {
        alert(response.error?.message || 'Failed to stop container');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleRestart = async (container: DockerContainer) => {
    setActionLoading(container.id);
    try {
      const response = await dockerApi.restartContainer(container.id);
      if (response.success) {
        loadContainers();
      } else {
        alert(response.error?.message || 'Failed to restart container');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (container: DockerContainer) => {
    setActionLoading(container.id);
    try {
      const response = await dockerApi.removeContainer(container.id);
      if (response.success) {
        setDeleteModal(null);
        loadContainers();
      } else {
        alert(response.error?.message || 'Failed to remove container');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleViewLogs = async (container: DockerContainer) => {
    try {
      const response = await dockerApi.getContainerLogs(container.id);
      if (response.success && response.data) {
        setLogsModal({ container, logs: response.data });
      } else {
        alert(response.error?.message || 'Failed to load logs');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleUpdateResources = async () => {
    if (!resourcesModal) return;

    try {
      const resources: any = {};

      if (resourceLimits.memory) {
        // Convert MB to bytes
        resources.memory = parseInt(resourceLimits.memory) * 1024 * 1024;
      }

      if (resourceLimits.cpuShares) {
        resources.cpuShares = parseInt(resourceLimits.cpuShares);
      }

      const response = await dockerApi.updateContainerResources(resourcesModal.id, resources);
      if (response.success) {
        setResourcesModal(null);
        setResourceLimits({ memory: '', cpuShares: '' });
        loadContainers();
      } else {
        alert(response.error?.message || 'Failed to update resources');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleExecCommand = async () => {
    if (!execModal || !execCommand.trim()) return;

    try {
      const response = await dockerApi.execContainer(execModal.id, execCommand.split(' '));
      if (response.success && response.data) {
        setExecOutput(response.data.output);
      } else {
        alert(response.error?.message || 'Failed to execute command');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const getStatusColor = (state: string) => {
    switch (state.toLowerCase()) {
      case 'running':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'exited':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
      case 'paused':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400';
      case 'restarting':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
    }
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const formatPorts = (ports?: DockerContainer['ports']) => {
    if (!ports || ports.length === 0) return 'No ports exposed';
    return ports
      .map((p) => (p.publicPort ? `${p.publicPort}:${p.privatePort}/${p.type}` : `${p.privatePort}/${p.type}`))
      .join(', ');
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

      {/* Controls */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowAll(true)}
            className={`px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
              showAll
                ? 'bg-macos-blue text-white'
                : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
            }`}
          >
            All Containers
          </button>
          <button
            onClick={() => setShowAll(false)}
            className={`px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
              !showAll
                ? 'bg-macos-blue text-white'
                : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
            }`}
          >
            Running Only
          </button>
        </div>
        <div className="text-sm text-gray-600 dark:text-gray-400">
          {containers.length} container{containers.length !== 1 ? 's' : ''}
        </div>
      </div>

      {/* Containers Grid */}
      {containers.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">📦</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No containers found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            {showAll
              ? 'No Docker containers are available'
              : 'No running containers. Try showing all containers.'}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
          {containers.map((container) => (
            <Card key={container.id} hoverable>
              <div className="p-6">
                {/* Header */}
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="font-bold text-lg text-gray-900 dark:text-gray-100">
                      {container.name}
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400 font-mono">
                      {container.id.substring(0, 12)}
                    </p>
                  </div>
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium ${getStatusColor(
                      container.state
                    )}`}
                  >
                    {container.state}
                  </span>
                </div>

                {/* Details */}
                <div className="space-y-2 mb-4">
                  <div className="text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Image:</span>
                    <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                      {container.image}
                    </span>
                  </div>
                  <div className="text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Status:</span>
                    <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                      {container.status}
                    </span>
                  </div>
                  <div className="text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Ports:</span>
                    <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                      {formatPorts(container.ports)}
                    </span>
                  </div>
                  <div className="text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Created:</span>
                    <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                      {formatDate(container.created)}
                    </span>
                  </div>
                </div>

                {/* Actions */}
                <div className="flex flex-wrap gap-2">
                  {container.state.toLowerCase() !== 'running' && (
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => handleStart(container)}
                      disabled={actionLoading === container.id}
                      className="flex-1 min-w-[80px]"
                    >
                      ▶️ Start
                    </Button>
                  )}
                  {container.state.toLowerCase() === 'running' && (
                    <>
                      <Button
                        size="sm"
                        variant="secondary"
                        onClick={() => handleStop(container)}
                        disabled={actionLoading === container.id}
                        className="flex-1 min-w-[80px]"
                      >
                        ⏸️ Stop
                      </Button>
                      <Button
                        size="sm"
                        variant="secondary"
                        onClick={() => handleRestart(container)}
                        disabled={actionLoading === container.id}
                        className="flex-1 min-w-[80px]"
                      >
                        🔄 Restart
                      </Button>
                    </>
                  )}
                  <Button
                    size="sm"
                    variant="secondary"
                    onClick={() => handleViewLogs(container)}
                    className="flex-1 min-w-[80px]"
                  >
                    📄 Logs
                  </Button>
                  {container.state.toLowerCase() === 'running' && (
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => {
                        setExecModal(container);
                        setExecOutput('');
                      }}
                      className="flex-1 min-w-[80px]"
                    >
                      💻 Exec
                    </Button>
                  )}
                  <Button
                    size="sm"
                    variant="secondary"
                    onClick={() => setResourcesModal(container)}
                    className="flex-1 min-w-[80px]"
                  >
                    ⚙️ Resources
                  </Button>
                  <Button
                    size="sm"
                    variant="danger"
                    onClick={() => setDeleteModal(container)}
                    disabled={actionLoading === container.id}
                    className="flex-1 min-w-[80px]"
                  >
                    🗑️ Delete
                  </Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Logs Modal */}
      <AnimatePresence>
        {logsModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setLogsModal(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl w-full max-w-4xl max-h-[80vh] flex flex-col"
            >
              <div className="p-6 border-b border-gray-200 dark:border-gray-700">
                <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                  Container Logs: {logsModal.container.name}
                </h2>
              </div>
              <div className="flex-1 overflow-auto p-6">
                <pre className="text-xs font-mono text-gray-900 dark:text-gray-100 whitespace-pre-wrap bg-gray-50 dark:bg-gray-900 p-4 rounded-lg">
                  {logsModal.logs || 'No logs available'}
                </pre>
              </div>
              <div className="p-6 border-t border-gray-200 dark:border-gray-700">
                <Button onClick={() => setLogsModal(null)} className="w-full">
                  Close
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {deleteModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setDeleteModal(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Delete Container
              </h2>
              <p className="text-gray-600 dark:text-gray-400 mb-6">
                Are you sure you want to delete container <strong>{deleteModal.name}</strong>? This
                action cannot be undone.
              </p>
              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setDeleteModal(null)}
                  className="flex-1"
                  disabled={actionLoading === deleteModal.id}
                >
                  Cancel
                </Button>
                <Button
                  variant="danger"
                  onClick={() => handleDelete(deleteModal)}
                  className="flex-1"
                  disabled={actionLoading === deleteModal.id}
                >
                  {actionLoading === deleteModal.id ? 'Deleting...' : 'Delete'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Resource Limits Modal */}
      <AnimatePresence>
        {resourcesModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
            onClick={() => setResourcesModal(null)}
          >
            <motion.div
              initial={{ scale: 0.95 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0.95 }}
              className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
                Update Resource Limits
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                Container: {resourcesModal.name}
              </p>

              <div className="space-y-4 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Memory Limit (MB)
                  </label>
                  <input
                    type="number"
                    value={resourceLimits.memory}
                    onChange={(e) =>
                      setResourceLimits({ ...resourceLimits, memory: e.target.value })
                    }
                    placeholder="e.g., 512"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    CPU Shares (Relative Weight)
                  </label>
                  <input
                    type="number"
                    value={resourceLimits.cpuShares}
                    onChange={(e) =>
                      setResourceLimits({ ...resourceLimits, cpuShares: e.target.value })
                    }
                    placeholder="e.g., 1024"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  />
                </div>
              </div>

              <div className="flex gap-2">
                <Button
                  variant="secondary"
                  onClick={() => {
                    setResourcesModal(null);
                    setResourceLimits({ memory: '', cpuShares: '' });
                  }}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button onClick={handleUpdateResources} className="flex-1">
                  Update
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Exec Command Modal */}
      <AnimatePresence>
        {execModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
            onClick={() => setExecModal(null)}
          >
            <motion.div
              initial={{ scale: 0.95 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0.95 }}
              className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-2xl w-full max-h-[80vh] overflow-hidden flex flex-col"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
                Execute Command
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                Container: {execModal.name}
              </p>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Command
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={execCommand}
                    onChange={(e) => setExecCommand(e.target.value)}
                    onKeyPress={(e) => {
                      if (e.key === 'Enter') {
                        handleExecCommand();
                      }
                    }}
                    placeholder="ls -la"
                    className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 font-mono text-sm"
                  />
                  <Button onClick={handleExecCommand}>Run</Button>
                </div>
              </div>

              {execOutput && (
                <div className="flex-1 overflow-y-auto">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Output
                  </label>
                  <pre className="bg-gray-900 text-green-400 p-4 rounded-md text-xs font-mono whitespace-pre-wrap overflow-x-auto">
                    {execOutput}
                  </pre>
                </div>
              )}

              <div className="mt-4 flex justify-end">
                <Button
                  variant="secondary"
                  onClick={() => {
                    setExecModal(null);
                    setExecCommand('ls -la');
                    setExecOutput('');
                  }}
                >
                  Close
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
