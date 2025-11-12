import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { dockerApi, ComposeStack } from '../../../api/docker';

const StackManager: React.FC = () => {
  const [stacks, setStacks] = useState<ComposeStack[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showLogsModal, setShowLogsModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [selectedStack, setSelectedStack] = useState<ComposeStack | null>(null);
  const [stackLogs, setStackLogs] = useState<string>('');
  const [newStackName, setNewStackName] = useState('');
  const [newStackCompose, setNewStackCompose] = useState('');
  const [editStackCompose, setEditStackCompose] = useState('');
  const [removeVolumes, setRemoveVolumes] = useState(false);

  const fetchStacks = async () => {
    try {
      setLoading(true);
      const response = await dockerApi.listStacks();
      setStacks(response.data || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch stacks');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStacks();
    const interval = setInterval(fetchStacks, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleCreateStack = async () => {
    if (!newStackName.trim() || !newStackCompose.trim()) {
      return;
    }

    try {
      await dockerApi.createStack({
        name: newStackName,
        compose: newStackCompose,
      });
      setShowCreateModal(false);
      setNewStackName('');
      setNewStackCompose('');
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to create stack');
    }
  };

  const handleUpdateStack = async () => {
    if (!selectedStack || !editStackCompose.trim()) {
      return;
    }

    try {
      await dockerApi.updateStack(selectedStack.name, {
        compose: editStackCompose,
      });
      setShowEditModal(false);
      setSelectedStack(null);
      setEditStackCompose('');
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to update stack');
    }
  };

  const handleDeployStack = async (name: string) => {
    try {
      await dockerApi.deployStack(name);
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to deploy stack');
    }
  };

  const handleStopStack = async (name: string) => {
    try {
      await dockerApi.stopStack(name);
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to stop stack');
    }
  };

  const handleRestartStack = async (name: string) => {
    try {
      await dockerApi.restartStack(name);
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to restart stack');
    }
  };

  const handleRemoveStack = async (name: string) => {
    try {
      await dockerApi.removeStack(name, removeVolumes);
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to remove stack');
    }
  };

  const handleDeleteStack = async () => {
    if (!selectedStack) return;

    try {
      await dockerApi.deleteStack(selectedStack.name);
      setShowDeleteModal(false);
      setSelectedStack(null);
      fetchStacks();
    } catch (err: any) {
      setError(err.message || 'Failed to delete stack');
    }
  };

  const handleViewLogs = async (stack: ComposeStack) => {
    try {
      const response = await dockerApi.getStackLogs(stack.name);
      setStackLogs(response.data || 'No logs available');
      setSelectedStack(stack);
      setShowLogsModal(true);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch logs');
    }
  };

  const handleEditStack = async (stack: ComposeStack) => {
    try {
      const response = await dockerApi.getStackCompose(stack.name);
      setEditStackCompose(response.data || '');
      setSelectedStack(stack);
      setShowEditModal(true);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch compose file');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running':
        return 'text-green-400';
      case 'stopped':
        return 'text-gray-400';
      case 'partial':
        return 'text-yellow-400';
      default:
        return 'text-gray-400';
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'running':
        return 'bg-green-500/20 text-green-400 border-green-500/30';
      case 'stopped':
        return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
      case 'partial':
        return 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30';
      default:
        return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
    }
  };

  if (loading && stacks.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-gray-400">Loading stacks...</div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold text-white">Docker Compose Stacks</h2>
          <p className="text-sm text-gray-400 mt-1">{stacks.length} stack(s)</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
        >
          Create Stack
        </button>
      </div>

      {/* Error Display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-4 bg-red-500/10 border border-red-500/30 rounded-lg"
          >
            <div className="flex justify-between items-start">
              <p className="text-red-400">{error}</p>
              <button
                onClick={() => setError(null)}
                className="text-red-400 hover:text-red-300"
              >
                ✕
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Stacks List */}
      {stacks.length === 0 ? (
        <div className="text-center py-12 text-gray-400">
          No stacks found. Create your first Docker Compose stack to get started.
        </div>
      ) : (
        <div className="space-y-3">
          {stacks.map((stack) => (
            <motion.div
              key={stack.name}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-gray-800/50 backdrop-blur-sm border border-gray-700/50 rounded-lg p-4"
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-lg font-medium text-white">{stack.name}</h3>
                    <span
                      className={`px-2 py-1 text-xs rounded-md border ${getStatusBadge(
                        stack.status
                      )}`}
                    >
                      {stack.status}
                    </span>
                  </div>
                  <div className="text-sm text-gray-400 space-y-1">
                    <div>Path: {stack.path}</div>
                    <div>Services: {stack.services?.length || 0}</div>
                    {stack.updatedAt && (
                      <div>Updated: {new Date(stack.updatedAt).toLocaleString()}</div>
                    )}
                  </div>

                  {/* Services */}
                  {stack.services && stack.services.length > 0 && (
                    <div className="mt-3 space-y-2">
                      {stack.services.map((service) => (
                        <div
                          key={service.name}
                          className="flex items-center justify-between p-2 bg-gray-900/50 rounded"
                        >
                          <div className="flex items-center gap-3">
                            <span className="text-white font-medium">{service.name}</span>
                            <span
                              className={`text-xs ${getStatusColor(service.status)}`}
                            >
                              {service.status}
                            </span>
                          </div>
                          {service.containers && service.containers.length > 0 && (
                            <span className="text-xs text-gray-500">
                              {service.containers.length} container(s)
                            </span>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                {/* Actions */}
                <div className="flex gap-2 ml-4">
                  {stack.status !== 'running' && (
                    <button
                      onClick={() => handleDeployStack(stack.name)}
                      className="p-2 hover:bg-green-500/20 text-green-400 rounded transition-colors"
                      title="Deploy"
                    >
                      ▶
                    </button>
                  )}
                  {stack.status === 'running' && (
                    <button
                      onClick={() => handleStopStack(stack.name)}
                      className="p-2 hover:bg-yellow-500/20 text-yellow-400 rounded transition-colors"
                      title="Stop"
                    >
                      ⏸
                    </button>
                  )}
                  <button
                    onClick={() => handleRestartStack(stack.name)}
                    className="p-2 hover:bg-blue-500/20 text-blue-400 rounded transition-colors"
                    title="Restart"
                  >
                    🔄
                  </button>
                  <button
                    onClick={() => handleViewLogs(stack)}
                    className="p-2 hover:bg-purple-500/20 text-purple-400 rounded transition-colors"
                    title="View Logs"
                  >
                    📋
                  </button>
                  <button
                    onClick={() => handleEditStack(stack)}
                    className="p-2 hover:bg-blue-500/20 text-blue-400 rounded transition-colors"
                    title="Edit"
                  >
                    ✏️
                  </button>
                  <button
                    onClick={() => {
                      setSelectedStack(stack);
                      setShowDeleteModal(true);
                    }}
                    className="p-2 hover:bg-red-500/20 text-red-400 rounded transition-colors"
                    title="Delete"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      )}

      {/* Create Stack Modal */}
      <AnimatePresence>
        {showCreateModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowCreateModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-gray-800 rounded-lg p-6 max-w-3xl w-full max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-white mb-4">Create New Stack</h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-gray-400 mb-2">Stack Name</label>
                  <input
                    type="text"
                    value={newStackName}
                    onChange={(e) => setNewStackName(e.target.value)}
                    placeholder="my-stack"
                    className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm text-gray-400 mb-2">
                    Docker Compose YAML
                  </label>
                  <textarea
                    value={newStackCompose}
                    onChange={(e) => setNewStackCompose(e.target.value)}
                    placeholder="version: '3.8'&#10;services:&#10;  web:&#10;    image: nginx:latest&#10;    ports:&#10;      - '80:80'"
                    rows={15}
                    className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 font-mono text-sm"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowCreateModal(false);
                    setNewStackName('');
                    setNewStackCompose('');
                  }}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateStack}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
                  disabled={!newStackName.trim() || !newStackCompose.trim()}
                >
                  Create Stack
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Edit Stack Modal */}
      <AnimatePresence>
        {showEditModal && selectedStack && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowEditModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-gray-800 rounded-lg p-6 max-w-3xl w-full max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-white mb-4">
                Edit Stack: {selectedStack.name}
              </h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-gray-400 mb-2">
                    Docker Compose YAML
                  </label>
                  <textarea
                    value={editStackCompose}
                    onChange={(e) => setEditStackCompose(e.target.value)}
                    rows={15}
                    className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 font-mono text-sm"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowEditModal(false);
                    setSelectedStack(null);
                    setEditStackCompose('');
                  }}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleUpdateStack}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
                  disabled={!editStackCompose.trim()}
                >
                  Update Stack
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Logs Modal */}
      <AnimatePresence>
        {showLogsModal && selectedStack && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowLogsModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-gray-800 rounded-lg p-6 max-w-4xl w-full max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-white mb-4">
                Stack Logs: {selectedStack.name}
              </h3>
              <pre className="bg-gray-900 p-4 rounded-lg text-xs text-gray-300 overflow-x-auto whitespace-pre-wrap">
                {stackLogs}
              </pre>
              <div className="flex justify-end mt-4">
                <button
                  onClick={() => setShowLogsModal(false)}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                >
                  Close
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {showDeleteModal && selectedStack && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowDeleteModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-gray-800 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-white mb-4">Delete Stack</h3>
              <p className="text-gray-300 mb-4">
                Are you sure you want to delete the stack "{selectedStack.name}"? This will
                remove the stack directory and compose file.
              </p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => {
                    setShowDeleteModal(false);
                    setSelectedStack(null);
                  }}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDeleteStack}
                  className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default StackManager;
