import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { pluginApi, Plugin } from '../../api/plugins';

export function PluginManager() {
  const [plugins, setPlugins] = useState<Plugin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showInstallModal, setShowInstallModal] = useState(false);
  const [showConfigModal, setShowConfigModal] = useState(false);
  const [showUninstallModal, setShowUninstallModal] = useState(false);
  const [selectedPlugin, setSelectedPlugin] = useState<Plugin | null>(null);
  const [installPath, setInstallPath] = useState('');
  const [configJson, setConfigJson] = useState('');

  const fetchPlugins = async () => {
    try {
      setLoading(true);
      const response = await pluginApi.listPlugins();
      setPlugins(response.data || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch plugins');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlugins();
    const interval = setInterval(fetchPlugins, 10000); // Refresh every 10s
    return () => clearInterval(interval);
  }, []);

  const handleInstall = async () => {
    if (!installPath.trim()) {
      return;
    }

    try {
      await pluginApi.installPlugin({ sourcePath: installPath });
      setShowInstallModal(false);
      setInstallPath('');
      fetchPlugins();
    } catch (err: any) {
      setError(err.message || 'Failed to install plugin');
    }
  };

  const handleUninstall = async () => {
    if (!selectedPlugin) return;

    try {
      await pluginApi.uninstallPlugin(selectedPlugin.id);
      setShowUninstallModal(false);
      setSelectedPlugin(null);
      fetchPlugins();
    } catch (err: any) {
      setError(err.message || 'Failed to uninstall plugin');
    }
  };

  const handleToggleEnabled = async (plugin: Plugin) => {
    try {
      if (plugin.enabled) {
        await pluginApi.disablePlugin(plugin.id);
      } else {
        await pluginApi.enablePlugin(plugin.id);
      }
      fetchPlugins();
    } catch (err: any) {
      setError(err.message || 'Failed to toggle plugin');
    }
  };

  const handleOpenConfig = (plugin: Plugin) => {
    setSelectedPlugin(plugin);
    setConfigJson(JSON.stringify(plugin.config || {}, null, 2));
    setShowConfigModal(true);
  };

  const handleSaveConfig = async () => {
    if (!selectedPlugin) return;

    try {
      const config = JSON.parse(configJson);
      await pluginApi.updatePluginConfig(selectedPlugin.id, { config });
      setShowConfigModal(false);
      setSelectedPlugin(null);
      setConfigJson('');
      fetchPlugins();
    } catch (err: any) {
      setError(err.message || 'Failed to update plugin config');
    }
  };

  const getStatusBadge = (plugin: Plugin) => {
    if (!plugin.installed) {
      return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
    }
    if (plugin.enabled) {
      return 'bg-green-500/20 text-green-400 border-green-500/30';
    }
    return 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30';
  };

  const getStatusText = (plugin: Plugin) => {
    if (!plugin.installed) return 'Not Installed';
    if (plugin.enabled) return 'Enabled';
    return 'Disabled';
  };

  if (loading && plugins.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-gray-400">Loading plugins...</div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Plugin Manager
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mt-1">
              Manage system plugins and extensions
            </p>
          </div>
          <button
            onClick={() => setShowInstallModal(true)}
            className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
          >
            Install Plugin
          </button>
        </div>
      </div>

      {/* Error Display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="mx-6 mt-4 p-4 bg-red-500/10 border border-red-500/30 rounded-lg"
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

      {/* Plugin Grid */}
      <div className="flex-1 overflow-auto p-6">
        {plugins.length === 0 ? (
          <div className="text-center py-12 text-gray-400">
            No plugins installed. Install your first plugin to get started.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {plugins.map((plugin) => (
              <motion.div
                key={plugin.id}
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-6 hover:shadow-lg transition-shadow"
              >
                {/* Plugin Header */}
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-lg flex items-center justify-center text-2xl">
                      {plugin.icon || '🔌'}
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                        {plugin.name}
                      </h3>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        v{plugin.version}
                      </p>
                    </div>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs rounded-md border ${getStatusBadge(
                      plugin
                    )}`}
                  >
                    {getStatusText(plugin)}
                  </span>
                </div>

                {/* Plugin Info */}
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-3">
                  {plugin.description}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-500 mb-4">
                  by {plugin.author}
                </p>

                {/* Actions */}
                <div className="flex gap-2">
                  {plugin.installed && (
                    <>
                      <button
                        onClick={() => handleToggleEnabled(plugin)}
                        className={`flex-1 px-3 py-2 text-sm rounded-lg transition-colors ${
                          plugin.enabled
                            ? 'bg-yellow-500/20 text-yellow-600 dark:text-yellow-400 hover:bg-yellow-500/30'
                            : 'bg-green-500/20 text-green-600 dark:text-green-400 hover:bg-green-500/30'
                        }`}
                      >
                        {plugin.enabled ? 'Disable' : 'Enable'}
                      </button>
                      <button
                        onClick={() => handleOpenConfig(plugin)}
                        className="px-3 py-2 text-sm bg-blue-500/20 text-blue-600 dark:text-blue-400 hover:bg-blue-500/30 rounded-lg transition-colors"
                        title="Configure"
                      >
                        ⚙️
                      </button>
                      <button
                        onClick={() => {
                          setSelectedPlugin(plugin);
                          setShowUninstallModal(true);
                        }}
                        className="px-3 py-2 text-sm bg-red-500/20 text-red-600 dark:text-red-400 hover:bg-red-500/30 rounded-lg transition-colors"
                        title="Uninstall"
                      >
                        🗑️
                      </button>
                    </>
                  )}
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </div>

      {/* Install Modal */}
      <AnimatePresence>
        {showInstallModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowInstallModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Install Plugin
              </h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                    Plugin Source Path
                  </label>
                  <input
                    type="text"
                    value={installPath}
                    onChange={(e) => setInstallPath(e.target.value)}
                    placeholder="/path/to/plugin"
                    className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Path to directory containing plugin.json
                  </p>
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowInstallModal(false);
                    setInstallPath('');
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleInstall}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
                  disabled={!installPath.trim()}
                >
                  Install
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Config Modal */}
      <AnimatePresence>
        {showConfigModal && selectedPlugin && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowConfigModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Configure {selectedPlugin.name}
              </h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                    Plugin Configuration (JSON)
                  </label>
                  <textarea
                    value={configJson}
                    onChange={(e) => setConfigJson(e.target.value)}
                    rows={15}
                    className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 font-mono text-sm"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowConfigModal(false);
                    setSelectedPlugin(null);
                    setConfigJson('');
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSaveConfig}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
                >
                  Save Config
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Uninstall Confirmation Modal */}
      <AnimatePresence>
        {showUninstallModal && selectedPlugin && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowUninstallModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Uninstall Plugin
              </h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Are you sure you want to uninstall "{selectedPlugin.name}"? This will
                remove all plugin files and cannot be undone.
              </p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => {
                    setShowUninstallModal(false);
                    setSelectedPlugin(null);
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleUninstall}
                  className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors"
                >
                  Uninstall
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
