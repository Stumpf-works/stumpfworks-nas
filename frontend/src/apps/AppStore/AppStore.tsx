import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { pluginApi } from '../../api/plugins';
import { Search, Download, Star, X, Package, Globe, Film, RefreshCw, BarChart3, Shield, Camera, ChevronRight } from 'lucide-react';

interface StorePlugin {
  id: string;
  name: string;
  version: string;
  author: string;
  description: string;
  longDescription?: string;
  icon: string;
  iconColor: string;
  category: string;
  tags: string[];
  downloadUrl: string;
  screenshots?: string[];
  rating?: number;
  downloads?: number;
  installed?: boolean;
  featured?: boolean;
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
    icon: 'Globe',
    iconColor: 'from-blue-500 to-cyan-500',
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
    icon: 'Film',
    iconColor: 'from-purple-500 to-pink-500',
    category: 'Media',
    tags: ['media', 'streaming', 'entertainment'],
    downloadUrl: 'https://example.com/plugins/plex.tar.gz',
    rating: 4.8,
    downloads: 15420,
    installed: false,
    featured: true,
  },
  {
    id: 'com.stumpfworks.syncthing',
    name: 'Syncthing',
    version: '1.5.3',
    author: 'Community',
    description: 'Continuous file synchronization across devices',
    longDescription: 'Syncthing is a continuous file synchronization program. It synchronizes files between two or more computers in real time, safely protected from prying eyes.',
    icon: 'RefreshCw',
    iconColor: 'from-green-500 to-emerald-500',
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
    icon: 'BarChart3',
    iconColor: 'from-orange-500 to-red-500',
    category: 'System',
    tags: ['monitoring', 'performance', 'system'],
    downloadUrl: 'https://example.com/plugins/monitoring.tar.gz',
    rating: 4.7,
    downloads: 6740,
    installed: false,
    featured: true,
  },
  {
    id: 'com.stumpfworks.vpn',
    name: 'VPN Server',
    version: '2.0.0',
    author: 'Community',
    description: 'Host your own WireGuard or OpenVPN server',
    longDescription: 'Set up a secure VPN server on your NAS using WireGuard or OpenVPN. Access your home network securely from anywhere in the world.',
    icon: 'Shield',
    iconColor: 'from-indigo-500 to-blue-500',
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
    icon: 'Camera',
    iconColor: 'from-pink-500 to-rose-500',
    category: 'Media',
    tags: ['photos', 'ai', 'media'],
    downloadUrl: 'https://example.com/plugins/photoprism.tar.gz',
    rating: 4.8,
    downloads: 9820,
    installed: false,
  },
];

const CATEGORIES = ['All', 'Media', 'Utilities', 'System', 'Security', 'Development'];

// Icon mapping
const iconMap: { [key: string]: any } = {
  Globe,
  Film,
  RefreshCw,
  BarChart3,
  Shield,
  Camera,
};

const getIconComponent = (iconName: string) => {
  return iconMap[iconName] || Package;
};

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

  const featuredPlugins = plugins.filter(p => p.featured);

  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-50">
      {/* Header - Simplified Apple Style */}
      <div className="px-6 md:px-12 pt-8 pb-6 bg-white dark:bg-macos-dark-100">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-2">
            App Store
          </h1>
          <p className="text-lg text-gray-500 dark:text-gray-400">
            Discover powerful plugins for your NAS
          </p>
        </div>
      </div>

      {/* Search and Categories */}
      <div className="px-6 md:px-12 pb-6 bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-7xl mx-auto">
          {/* Search Bar */}
          <div className="relative mb-6">
            <input
              type="text"
              placeholder="Search apps and plugins"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full px-5 py-3.5 pl-12 bg-gray-100 dark:bg-macos-dark-50 border-none rounded-xl text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
            />
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          </div>

          {/* Categories - Subtle Pills */}
          <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
            {CATEGORIES.map(category => (
              <button
                key={category}
                onClick={() => setSelectedCategory(category)}
                className={`px-5 py-2 rounded-full whitespace-nowrap text-sm font-medium transition-all ${
                  selectedCategory === category
                    ? 'bg-macos-blue text-white shadow-sm'
                    : 'bg-gray-100 dark:bg-macos-dark-50 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
                }`}
              >
                {category}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto">
        <div className="max-w-7xl mx-auto px-6 md:px-12 py-8">
          {/* Error Display */}
          <AnimatePresence>
            {error && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl"
              >
                <div className="flex justify-between items-start">
                  <p className="text-red-600 dark:text-red-400">{error}</p>
                  <button
                    onClick={() => setError(null)}
                    className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300"
                  >
                    <X className="w-5 h-5" />
                  </button>
                </div>
              </motion.div>
            )}
          </AnimatePresence>

          {/* Featured Apps - Apple Style */}
          {!searchQuery && selectedCategory === 'All' && featuredPlugins.length > 0 && (
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Featured</h2>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {featuredPlugins.map((plugin) => {
                  const IconComponent = getIconComponent(plugin.icon);
                  return (
                    <motion.div
                      key={plugin.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="group relative bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-100 dark:to-macos-dark-200 rounded-2xl overflow-hidden hover:shadow-2xl transition-all duration-300 cursor-pointer"
                      onClick={() => setSelectedPlugin(plugin)}
                    >
                      <div className="p-8">
                        <div className="flex items-start gap-6">
                          {/* Large Icon */}
                          <div className={`w-24 h-24 bg-gradient-to-br ${plugin.iconColor} rounded-3xl flex items-center justify-center shadow-lg shrink-0`}>
                            <IconComponent className="w-12 h-12 text-white" />
                          </div>

                          {/* Content */}
                          <div className="flex-1 min-w-0">
                            <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
                              {plugin.name}
                            </h3>
                            <p className="text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
                              {plugin.description}
                            </p>

                            {/* Stats */}
                            <div className="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400 mb-4">
                              {plugin.rating && (
                                <div className="flex items-center gap-1">
                                  <Star className="w-4 h-4 fill-yellow-500 text-yellow-500" />
                                  <span className="font-medium">{plugin.rating.toFixed(1)}</span>
                                </div>
                              )}
                              <span className="px-2 py-1 bg-white/50 dark:bg-black/20 rounded-md text-xs font-medium">
                                {plugin.category}
                              </span>
                            </div>

                            {/* Get Button */}
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                if (!plugin.installed && !installing) {
                                  handleInstall(plugin);
                                }
                              }}
                              disabled={plugin.installed || installing === plugin.id}
                              className={`px-8 py-2.5 rounded-full font-semibold text-sm transition-all ${
                                plugin.installed
                                  ? 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-default'
                                  : installing === plugin.id
                                  ? 'bg-blue-400 text-white cursor-wait'
                                  : 'bg-macos-blue hover:bg-blue-600 text-white shadow-lg hover:shadow-xl'
                              }`}
                            >
                              {plugin.installed ? 'Installed' : installing === plugin.id ? 'Installing...' : 'GET'}
                            </button>
                          </div>

                          {/* Arrow */}
                          <ChevronRight className="w-6 h-6 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
                        </div>
                      </div>
                    </motion.div>
                  );
                })}
              </div>
            </div>
          )}

          {/* All Apps Grid */}
          <div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
              {searchQuery ? 'Search Results' : selectedCategory === 'All' ? 'All Apps' : selectedCategory}
            </h2>

            {filteredPlugins.length === 0 ? (
              <div className="text-center py-16">
                <Package className="w-20 h-20 mx-auto mb-4 text-gray-300 dark:text-gray-600" />
                <p className="text-gray-500 dark:text-gray-400 text-lg">No apps found matching your criteria</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredPlugins.map((plugin) => {
                  const IconComponent = getIconComponent(plugin.icon);
                  return (
                    <motion.div
                      key={plugin.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      whileHover={{ y: -4 }}
                      className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-2xl overflow-hidden hover:shadow-xl transition-all duration-300 cursor-pointer"
                      onClick={() => setSelectedPlugin(plugin)}
                    >
                      <div className="p-6">
                        {/* Icon and Title */}
                        <div className="flex items-center gap-4 mb-4">
                          <div className={`w-16 h-16 bg-gradient-to-br ${plugin.iconColor} rounded-2xl flex items-center justify-center shadow-md shrink-0`}>
                            <IconComponent className="w-8 h-8 text-white" />
                          </div>
                          <div className="flex-1 min-w-0">
                            <h3 className="font-bold text-lg text-gray-900 dark:text-white truncate">
                              {plugin.name}
                            </h3>
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                              {plugin.author}
                            </p>
                          </div>
                        </div>

                        {/* Description */}
                        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-2 min-h-[2.5rem]">
                          {plugin.description}
                        </p>

                        {/* Stats Row */}
                        <div className="flex items-center gap-3 mb-4 text-xs">
                          {plugin.rating && (
                            <div className="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                              <Star className="w-3.5 h-3.5 fill-yellow-500 text-yellow-500" />
                              <span className="font-medium">{plugin.rating.toFixed(1)}</span>
                            </div>
                          )}
                          <span className="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 rounded-md font-medium">
                            {plugin.category}
                          </span>
                        </div>

                        {/* Action Button */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            if (!plugin.installed && !installing) {
                              handleInstall(plugin);
                            }
                          }}
                          disabled={plugin.installed || installing === plugin.id}
                          className={`w-full px-6 py-2.5 rounded-full font-semibold text-sm transition-all ${
                            plugin.installed
                              ? 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 cursor-default'
                              : installing === plugin.id
                              ? 'bg-blue-400 text-white cursor-wait'
                              : 'bg-macos-blue hover:bg-blue-600 text-white shadow-md hover:shadow-lg'
                          }`}
                        >
                          {plugin.installed ? 'Installed' : installing === plugin.id ? 'Installing...' : 'GET'}
                        </button>
                      </div>
                    </motion.div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Plugin Detail Modal - Apple Style */}
      <AnimatePresence>
        {selectedPlugin && (() => {
          const IconComponent = getIconComponent(selectedPlugin.icon);
          return (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 bg-black/60 backdrop-blur-md flex items-center justify-center z-50 p-4"
              onClick={() => setSelectedPlugin(null)}
            >
              <motion.div
                initial={{ scale: 0.95, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                exit={{ scale: 0.95, opacity: 0 }}
                className="bg-white dark:bg-macos-dark-100 rounded-3xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden"
                onClick={(e) => e.stopPropagation()}
              >
                <div className="overflow-y-auto max-h-[90vh]">
                  {/* Header */}
                  <div className="p-8 pb-6">
                    <div className="flex items-start justify-between mb-6">
                      <div className="flex items-center gap-5">
                        <div className={`w-24 h-24 bg-gradient-to-br ${selectedPlugin.iconColor} rounded-3xl flex items-center justify-center shadow-xl`}>
                          <IconComponent className="w-12 h-12 text-white" />
                        </div>
                        <div>
                          <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-1">
                            {selectedPlugin.name}
                          </h2>
                          <p className="text-gray-500 dark:text-gray-400 mb-2">
                            {selectedPlugin.author}
                          </p>
                          <p className="text-sm text-gray-400 dark:text-gray-500">
                            Version {selectedPlugin.version}
                          </p>
                        </div>
                      </div>
                      <button
                        onClick={() => setSelectedPlugin(null)}
                        className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors p-2"
                      >
                        <X className="w-6 h-6" />
                      </button>
                    </div>

                    {/* Stats Row */}
                    <div className="flex items-center gap-6 pb-6 border-b border-gray-200 dark:border-gray-700">
                      {selectedPlugin.rating && (
                        <div className="flex items-center gap-2">
                          <div className="flex gap-0.5">
                            {[...Array(5)].map((_, i) => (
                              <Star
                                key={i}
                                className={`w-4 h-4 ${
                                  i < Math.floor(selectedPlugin.rating!)
                                    ? 'fill-yellow-500 text-yellow-500'
                                    : 'fill-gray-300 dark:fill-gray-600 text-gray-300 dark:text-gray-600'
                                }`}
                              />
                            ))}
                          </div>
                          <span className="text-sm font-semibold text-gray-900 dark:text-white">
                            {selectedPlugin.rating.toFixed(1)}
                          </span>
                        </div>
                      )}
                      {selectedPlugin.downloads && (
                        <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                          <Download className="w-4 h-4" />
                          <span>{selectedPlugin.downloads.toLocaleString()} downloads</span>
                        </div>
                      )}
                      <span className="px-3 py-1.5 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg text-sm font-medium">
                        {selectedPlugin.category}
                      </span>
                    </div>
                  </div>

                  {/* Content */}
                  <div className="px-8 pb-6">
                    {/* Description */}
                    <div className="mb-6">
                      <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-3">
                        About this app
                      </h3>
                      <p className="text-gray-600 dark:text-gray-400 leading-relaxed">
                        {selectedPlugin.longDescription || selectedPlugin.description}
                      </p>
                    </div>

                    {/* Tags */}
                    {selectedPlugin.tags.length > 0 && (
                      <div className="mb-6">
                        <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-3">
                          Keywords
                        </h3>
                        <div className="flex flex-wrap gap-2">
                          {selectedPlugin.tags.map(tag => (
                            <span
                              key={tag}
                              className="px-3 py-1.5 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg text-sm font-medium"
                            >
                              {tag}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Action Footer */}
                  <div className="px-8 py-6 bg-gray-50 dark:bg-macos-dark-50 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center justify-between">
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        {selectedPlugin.installed ? 'This app is installed on your NAS' : 'Click GET to install this app'}
                      </div>
                      <div className="flex gap-3">
                        <button
                          onClick={() => setSelectedPlugin(null)}
                          className="px-6 py-2.5 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-full transition-colors font-semibold text-sm"
                        >
                          Close
                        </button>
                        <button
                          onClick={() => {
                            if (!selectedPlugin.installed && !installing) {
                              handleInstall(selectedPlugin);
                            }
                          }}
                          disabled={selectedPlugin.installed || installing === selectedPlugin.id}
                          className={`px-8 py-2.5 rounded-full font-semibold text-sm transition-all ${
                            selectedPlugin.installed
                              ? 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-default'
                              : installing === selectedPlugin.id
                              ? 'bg-blue-400 text-white cursor-wait'
                              : 'bg-macos-blue hover:bg-blue-600 text-white shadow-lg hover:shadow-xl'
                          }`}
                        >
                          {selectedPlugin.installed ? 'Installed' : installing === selectedPlugin.id ? 'Installing...' : 'GET'}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </motion.div>
            </motion.div>
          );
        })()}
      </AnimatePresence>
    </div>
  );
}
