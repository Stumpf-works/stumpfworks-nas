import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { pluginApi } from '../../api/plugins';
import { Search, Download, Star, X, Package } from 'lucide-react';
import { toast } from 'react-hot-toast';

interface StorePlugin {
  id: string;
  name: string;
  version: string;
  author: string;
  description: string;
  longDescription?: string;
  icon: string;
  category: string;
  tags: string[];
  downloadUrl: string;
  screenshots?: string[];
  rating?: number;
  downloads?: number;
  installed?: boolean;
}

// Sample plugin repository - in production this would come from an API
const SAMPLE_PLUGINS: StorePlugin[] = [
  {
    id: 'com.stumpfworks.hello-world',
    name: 'Hello World',
    version: '1.0.0',
    author: 'StumpfWorks Team',
    description: 'A simple example plugin that demonstrates the plugin system',
    longDescription: 'This is a comprehensive example plugin that shows how to build plugins for StumpfWorks NAS. It includes periodic logging, configuration management, and graceful shutdown handling.',
    icon: '🌍',
    category: 'Development',
    tags: ['example', 'tutorial', 'development'],
    downloadUrl: '/examples/plugins/hello-world',
    rating: 4.5,
    downloads: 1250,
    installed: false,
  },
  {
    id: 'com.stumpfworks.plex',
    name: 'Plex Media Server',
    version: '2.1.0',
    author: 'Community',
    description: 'Stream your media library to any device with Plex',
    longDescription: 'Plex organizes your video, music, and photo collections and streams them to all of your devices. This plugin integrates Plex Media Server directly into your StumpfWorks NAS.',
    icon: '🎬',
    category: 'Media',
    tags: ['media', 'streaming', 'entertainment'],
    downloadUrl: 'https://example.com/plugins/plex.tar.gz',
    rating: 4.8,
    downloads: 15420,
    installed: false,
  },
  {
    id: 'com.stumpfworks.syncthing',
    name: 'Syncthing',
    version: '1.5.3',
    author: 'Community',
    description: 'Continuous file synchronization across devices',
    longDescription: 'Syncthing is a continuous file synchronization program. It synchronizes files between two or more computers in real time, safely protected from prying eyes.',
    icon: '🔄',
    category: 'Utilities',
    tags: ['sync', 'backup', 'utilities'],
    downloadUrl: 'https://example.com/plugins/syncthing.tar.gz',
    rating: 4.6,
    downloads: 8930,
    installed: false,
  },
  {
    id: 'com.stumpfworks.monitoring',
    name: 'Advanced Monitoring',
    version: '3.2.1',
    author: 'StumpfWorks Team',
    description: 'Real-time system monitoring with graphs and alerts',
    longDescription: 'Get detailed insights into your system performance with real-time graphs, historical data, and customizable alerts for CPU, memory, disk, and network usage.',
    icon: '📊',
    category: 'System',
    tags: ['monitoring', 'performance', 'system'],
    downloadUrl: 'https://example.com/plugins/monitoring.tar.gz',
    rating: 4.7,
    downloads: 6740,
    installed: false,
  },
  {
    id: 'com.stumpfworks.vpn',
    name: 'VPN Server',
    version: '2.0.0',
    author: 'Community',
    description: 'Host your own WireGuard or OpenVPN server',
    longDescription: 'Set up a secure VPN server on your NAS using WireGuard or OpenVPN. Access your home network securely from anywhere in the world.',
    icon: '🔐',
    category: 'Security',
    tags: ['vpn', 'security', 'networking'],
    downloadUrl: 'https://example.com/plugins/vpn.tar.gz',
    rating: 4.9,
    downloads: 12340,
    installed: false,
  },
  {
    id: 'com.stumpfworks.photoprism',
    name: 'PhotoPrism',
    version: '1.8.0',
    author: 'Community',
    description: 'AI-powered photo management and organization',
    longDescription: 'PhotoPrism uses artificial intelligence to automatically organize and tag your photos. Browse your collection by location, date, or subject with ease.',
    icon: '📸',
    category: 'Media',
    tags: ['photos', 'ai', 'media'],
    downloadUrl: 'https://example.com/plugins/photoprism.tar.gz',
    rating: 4.8,
    downloads: 9820,
    installed: false,
  },
];

const CATEGORIES = ['All', 'Media', 'Utilities', 'System', 'Security', 'Development'];

export function AppStore() {
  const [plugins, setPlugins] = useState<StorePlugin[]>(SAMPLE_PLUGINS);
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedPlugin, setSelectedPlugin] = useState<StorePlugin | null>(null);
  const [installing, setInstalling] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchInstalledPlugins();
  }, []);

  const fetchInstalledPlugins = async () => {
    try {
      const response = await pluginApi.listPlugins();
      const installed = new Set(response.data?.map((p: any) => p.id) || []);

      // Update plugin installation status
      setPlugins(prev => prev.map(p => ({
        ...p,
        installed: installed.has(p.id)
      })));
    } catch (err: any) {
      console.error('Failed to fetch installed plugins:', err);
    }
  };

  const handleInstall = async (plugin: StorePlugin) => {
    setInstalling(plugin.id);
    setError(null);

    try {
      await pluginApi.installPlugin({ sourcePath: plugin.downloadUrl });

      // Update local state
      setPlugins(prev => prev.map(p =>
        p.id === plugin.id ? { ...p, installed: true } : p
      ));

      // Refresh installed plugins list
      await fetchInstalledPlugins();
    } catch (err: any) {
      setError(err.message || 'Failed to install plugin');
    } finally {
      setInstalling(null);
    }
  };

  const filteredPlugins = plugins.filter(plugin => {
    const matchesCategory = selectedCategory === 'All' || plugin.category === selectedCategory;
    const matchesSearch = !searchQuery ||
      plugin.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      plugin.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      plugin.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()));

    return matchesCategory && matchesSearch;
  });

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <div className="flex justify-between items-center mb-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              App Store
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mt-1">
              Discover and install plugins for your NAS
            </p>
          </div>
        </div>

        {/* Search Bar */}
        <div className="relative">
          <input
            type="text"
            placeholder="Search plugins..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full px-4 py-3 pl-12 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500"
          />
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        </div>

        {/* Categories */}
        <div className="flex gap-2 mt-4 overflow-x-auto pb-2">
          {CATEGORIES.map(category => (
            <button
              key={category}
              onClick={() => setSelectedCategory(category)}
              className={`px-4 py-2 rounded-lg whitespace-nowrap transition-colors ${
                selectedCategory === category
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-100 dark:bg-macos-dark-50 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
              }`}
            >
              {category}
            </button>
          ))}
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
        {filteredPlugins.length === 0 ? (
          <div className="text-center py-12">
            <Package className="w-16 h-16 mx-auto mb-4 text-gray-400 dark:text-gray-600" />
            <p className="text-gray-600 dark:text-gray-400">No plugins found matching your criteria.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredPlugins.map((plugin) => (
              <motion.div
                key={plugin.id}
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                whileHover={{ scale: 1.02 }}
                className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden hover:shadow-xl transition-all cursor-pointer"
                onClick={() => setSelectedPlugin(plugin)}
              >
                {/* Plugin Card */}
                <div className="p-6">
                  {/* Header */}
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl flex items-center justify-center text-3xl shadow-lg">
                        {plugin.icon}
                      </div>
                      <div>
                        <h3 className="font-semibold text-lg text-gray-900 dark:text-gray-100">
                          {plugin.name}
                        </h3>
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                          v{plugin.version}
                        </p>
                      </div>
                    </div>
                  </div>

                  {/* Description */}
                  <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
                    {plugin.description}
                  </p>

                  {/* Stats */}
                  <div className="flex items-center gap-4 mb-4 text-xs text-gray-500 dark:text-gray-400">
                    {plugin.rating && (
                      <div className="flex items-center gap-1">
                        <Star className="w-4 h-4 fill-yellow-500 text-yellow-500" />
                        <span>{plugin.rating.toFixed(1)}</span>
                      </div>
                    )}
                    {plugin.downloads && (
                      <div className="flex items-center gap-1">
                        <Download className="w-4 h-4" />
                        <span>{plugin.downloads.toLocaleString()}</span>
                      </div>
                    )}
                    <span className="px-2 py-1 bg-blue-500/10 text-blue-600 dark:text-blue-400 rounded-md">
                      {plugin.category}
                    </span>
                  </div>

                  {/* Author */}
                  <p className="text-xs text-gray-500 dark:text-gray-500 mb-4">
                    by {plugin.author}
                  </p>

                  {/* Install Button */}
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      if (!plugin.installed && !installing) {
                        handleInstall(plugin);
                      }
                    }}
                    disabled={plugin.installed || installing === plugin.id}
                    className={`w-full px-4 py-2 rounded-lg font-medium transition-colors ${
                      plugin.installed
                        ? 'bg-green-500/20 text-green-600 dark:text-green-400 cursor-default'
                        : installing === plugin.id
                        ? 'bg-blue-500/50 text-white cursor-wait'
                        : 'bg-blue-500 hover:bg-blue-600 text-white'
                    }`}
                  >
                    {plugin.installed ? '✓ Installed' : installing === plugin.id ? 'Installing...' : 'Install'}
                  </button>
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </div>

      {/* Plugin Detail Modal */}
      <AnimatePresence>
        {selectedPlugin && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setSelectedPlugin(null)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-8 max-w-3xl w-full max-h-[90vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Header */}
              <div className="flex items-start justify-between mb-6">
                <div className="flex items-center gap-4">
                  <div className="w-20 h-20 bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl flex items-center justify-center text-4xl shadow-lg">
                    {selectedPlugin.icon}
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                      {selectedPlugin.name}
                    </h2>
                    <p className="text-gray-500 dark:text-gray-400">
                      v{selectedPlugin.version} by {selectedPlugin.author}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => setSelectedPlugin(null)}
                  className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
                >
                  <X className="w-6 h-6" />
                </button>
              </div>

              {/* Stats */}
              <div className="flex items-center gap-6 mb-6 pb-6 border-b border-gray-200 dark:border-gray-700">
                {selectedPlugin.rating && (
                  <div className="flex items-center gap-2">
                    <Star className="w-5 h-5 fill-yellow-500 text-yellow-500" />
                    <span className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {selectedPlugin.rating.toFixed(1)}
                    </span>
                  </div>
                )}
                {selectedPlugin.downloads && (
                  <div className="flex items-center gap-2">
                    <Download className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                    <span className="text-gray-600 dark:text-gray-400">
                      {selectedPlugin.downloads.toLocaleString()} downloads
                    </span>
                  </div>
                )}
                <span className="px-3 py-1 bg-blue-500/10 text-blue-600 dark:text-blue-400 rounded-lg font-medium">
                  {selectedPlugin.category}
                </span>
              </div>

              {/* Description */}
              <div className="mb-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  About
                </h3>
                <p className="text-gray-600 dark:text-gray-400 leading-relaxed">
                  {selectedPlugin.longDescription || selectedPlugin.description}
                </p>
              </div>

              {/* Tags */}
              {selectedPlugin.tags.length > 0 && (
                <div className="mb-6">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-3">
                    Tags
                  </h3>
                  <div className="flex flex-wrap gap-2">
                    {selectedPlugin.tags.map(tag => (
                      <span
                        key={tag}
                        className="px-3 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-full text-sm"
                      >
                        #{tag}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {/* Install Button */}
              <div className="flex gap-3">
                <button
                  onClick={() => {
                    if (!selectedPlugin.installed && !installing) {
                      handleInstall(selectedPlugin);
                    }
                  }}
                  disabled={selectedPlugin.installed || installing === selectedPlugin.id}
                  className={`flex-1 px-6 py-3 rounded-lg font-medium transition-colors ${
                    selectedPlugin.installed
                      ? 'bg-green-500/20 text-green-600 dark:text-green-400 cursor-default'
                      : installing === selectedPlugin.id
                      ? 'bg-blue-500/50 text-white cursor-wait'
                      : 'bg-blue-500 hover:bg-blue-600 text-white'
                  }`}
                >
                  {selectedPlugin.installed ? '✓ Installed' : installing === selectedPlugin.id ? 'Installing...' : 'Install Plugin'}
                </button>
                <button
                  onClick={() => setSelectedPlugin(null)}
                  className="px-6 py-3 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Close
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
