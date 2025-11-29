import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { addonsApi, type AddonWithStatus } from '../../api/addons';
import { Search, X, Package, AlertCircle, ChevronRight } from 'lucide-react';
import toast from 'react-hot-toast';


const getCategoryColor = (category: string) => {
  const colors: { [key: string]: string } = {
    virtualization: 'from-purple-500 to-indigo-500',
    storage: 'from-blue-500 to-cyan-500',
    media: 'from-pink-500 to-rose-500',
    security: 'from-green-500 to-emerald-500',
    networking: 'from-orange-500 to-red-500',
  };
  return colors[category.toLowerCase()] || 'from-gray-500 to-gray-600';
};

const CATEGORIES = ['All', 'Virtualization', 'Storage', 'Media', 'Security', 'Networking'];

export function AppStore() {
  const [addons, setAddons] = useState<AddonWithStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedAddon, setSelectedAddon] = useState<AddonWithStatus | null>(null);
  const [installing, setInstalling] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAddons();
  }, []);

  const loadAddons = async () => {
    try {
      setLoading(true);
      const response = await addonsApi.listAddons();

      if (response.success && response.data) {
        setAddons(response.data);
      } else {
        setError(response.error?.message || 'Failed to load addons');
      }
    } catch (err: any) {
      console.error('Failed to load addons:', err);
      setError('Failed to load addons from server');
    } finally {
      setLoading(false);
    }
  };

  const handleInstall = async (addon: AddonWithStatus) => {
    setInstalling(addon.manifest.id);
    setError(null);

    try {
      const response = await addonsApi.installAddon(addon.manifest.id);

      if (response.success) {
        toast.success(`${addon.manifest.name} installed successfully!`);

        // Refresh addons list to update status
        await loadAddons();

        // Close modal if open
        if (selectedAddon?.manifest.id === addon.manifest.id) {
          setSelectedAddon(null);
        }
      } else {
        const errorMsg = response.error?.message || 'Failed to install addon';
        setError(errorMsg);
        toast.error(errorMsg);
      }
    } catch (err: any) {
      const errorMsg = err.message || 'Failed to install addon';
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setInstalling(null);
    }
  };

  const handleUninstall = async (addon: AddonWithStatus) => {
    if (!confirm(`Are you sure you want to uninstall ${addon.manifest.name}?`)) {
      return;
    }

    setInstalling(addon.manifest.id);
    setError(null);

    try {
      const response = await addonsApi.uninstallAddon(addon.manifest.id);

      if (response.success) {
        toast.success(`${addon.manifest.name} uninstalled successfully`);

        // Refresh addons list
        await loadAddons();

        // Close modal if open
        if (selectedAddon?.manifest.id === addon.manifest.id) {
          setSelectedAddon(null);
        }
      } else {
        const errorMsg = response.error?.message || 'Failed to uninstall addon';
        setError(errorMsg);
        toast.error(errorMsg);
      }
    } catch (err: any) {
      const errorMsg = err.message || 'Failed to uninstall addon';
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setInstalling(null);
    }
  };

  const filteredAddons = addons.filter(addon => {
    const matchesCategory = selectedCategory === 'All' ||
      addon.manifest.category.toLowerCase() === selectedCategory.toLowerCase();

    const matchesSearch = !searchQuery ||
      addon.manifest.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      addon.manifest.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      addon.manifest.category.toLowerCase().includes(searchQuery.toLowerCase());

    return matchesCategory && matchesSearch;
  });

  // Featured addons (e.g., VM Manager and MinIO)
  const featuredAddons = addons.filter(a =>
    ['vm-manager', 'minio'].includes(a.manifest.id)
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full bg-white dark:bg-macos-dark-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue mx-auto mb-4"></div>
          <p className="text-gray-500 dark:text-gray-400">Loading App Store...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-50">
      {/* Header - Apple Style */}
      <div className="px-6 md:px-12 pt-8 pb-6 bg-white dark:bg-macos-dark-100">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-2">
            App Store
          </h1>
          <p className="text-lg text-gray-500 dark:text-gray-400">
            Discover and install powerful addons for your NAS
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
              placeholder="Search apps and addons"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full px-5 py-3.5 pl-12 bg-gray-100 dark:bg-macos-dark-50 border-none rounded-xl text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
            />
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          </div>

          {/* Categories */}
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
                className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl flex items-start gap-3"
              >
                <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 shrink-0 mt-0.5" />
                <div className="flex-1">
                  <p className="text-red-600 dark:text-red-400">{error}</p>
                </div>
                <button
                  onClick={() => setError(null)}
                  className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300"
                >
                  <X className="w-5 h-5" />
                </button>
              </motion.div>
            )}
          </AnimatePresence>

          {/* Featured Apps */}
          {!searchQuery && selectedCategory === 'All' && featuredAddons.length > 0 && (
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Featured</h2>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {featuredAddons.map((addon) => {
                  const isInstalled = addon.status.installed;
                  const isInstalling = installing === addon.manifest.id;

                  return (
                    <motion.div
                      key={addon.manifest.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="group relative bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-100 dark:to-macos-dark-200 rounded-2xl overflow-hidden hover:shadow-2xl transition-all duration-300 cursor-pointer"
                      onClick={() => setSelectedAddon(addon)}
                    >
                      <div className="p-8">
                        <div className="flex items-start gap-6">
                          {/* Large Icon */}
                          <div className={`w-24 h-24 bg-gradient-to-br ${getCategoryColor(addon.manifest.category)} rounded-3xl flex items-center justify-center shadow-lg shrink-0 text-5xl`}>
                            {addon.manifest.icon}
                          </div>

                          {/* Content */}
                          <div className="flex-1 min-w-0">
                            <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
                              {addon.manifest.name}
                            </h3>
                            <p className="text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
                              {addon.manifest.description}
                            </p>

                            {/* Stats */}
                            <div className="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400 mb-4">
                              <span className="px-2 py-1 bg-white/50 dark:bg-black/20 rounded-md text-xs font-medium">
                                {addon.manifest.category}
                              </span>
                              <span className="text-xs">v{addon.manifest.version}</span>
                            </div>

                            {/* Action Buttons */}
                            <div className="flex gap-3">
                              <button
                                onClick={(e) => {
                                  e.stopPropagation();
                                  if (!isInstalled && !isInstalling) {
                                    handleInstall(addon);
                                  }
                                }}
                                disabled={isInstalled || isInstalling}
                                className={`px-8 py-2.5 rounded-full font-semibold text-sm transition-all ${
                                  isInstalled
                                    ? 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-default'
                                    : isInstalling
                                    ? 'bg-blue-400 text-white cursor-wait'
                                    : 'bg-macos-blue hover:bg-blue-600 text-white shadow-lg hover:shadow-xl'
                                }`}
                              >
                                {isInstalled ? 'Installed' : isInstalling ? 'Installing...' : 'GET'}
                              </button>

                              {isInstalled && (
                                <button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleUninstall(addon);
                                  }}
                                  disabled={isInstalling}
                                  className="px-6 py-2.5 rounded-full font-semibold text-sm bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 hover:bg-red-200 dark:hover:bg-red-900/50 transition-all"
                                >
                                  Uninstall
                                </button>
                              )}
                            </div>
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

            {filteredAddons.length === 0 ? (
              <div className="text-center py-16">
                <Package className="w-20 h-20 mx-auto mb-4 text-gray-300 dark:text-gray-600" />
                <p className="text-gray-500 dark:text-gray-400 text-lg">No apps found matching your criteria</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredAddons.map((addon) => {
                  const isInstalled = addon.status.installed;
                  const isInstalling = installing === addon.manifest.id;

                  return (
                    <motion.div
                      key={addon.manifest.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      whileHover={{ y: -4 }}
                      className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-2xl overflow-hidden hover:shadow-xl transition-all duration-300 cursor-pointer"
                      onClick={() => setSelectedAddon(addon)}
                    >
                      <div className="p-6">
                        {/* Icon and Title */}
                        <div className="flex items-center gap-4 mb-4">
                          <div className={`w-16 h-16 bg-gradient-to-br ${getCategoryColor(addon.manifest.category)} rounded-2xl flex items-center justify-center shadow-md shrink-0 text-3xl`}>
                            {addon.manifest.icon}
                          </div>
                          <div className="flex-1 min-w-0">
                            <h3 className="font-bold text-lg text-gray-900 dark:text-white truncate">
                              {addon.manifest.name}
                            </h3>
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                              {addon.manifest.author}
                            </p>
                          </div>
                        </div>

                        {/* Description */}
                        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-2 min-h-[2.5rem]">
                          {addon.manifest.description}
                        </p>

                        {/* Stats */}
                        <div className="flex items-center gap-3 mb-4 text-xs">
                          <span className="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 rounded-md font-medium">
                            {addon.manifest.category}
                          </span>
                          <span className="text-gray-500 dark:text-gray-400">v{addon.manifest.version}</span>
                        </div>

                        {/* Action Buttons */}
                        <div className="flex gap-2">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              if (!isInstalled && !isInstalling) {
                                handleInstall(addon);
                              }
                            }}
                            disabled={isInstalled || isInstalling}
                            className={`flex-1 px-6 py-2.5 rounded-full font-semibold text-sm transition-all ${
                              isInstalled
                                ? 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 cursor-default'
                                : isInstalling
                                ? 'bg-blue-400 text-white cursor-wait'
                                : 'bg-macos-blue hover:bg-blue-600 text-white shadow-md hover:shadow-lg'
                            }`}
                          >
                            {isInstalled ? 'Installed' : isInstalling ? 'Installing...' : 'GET'}
                          </button>

                          {isInstalled && (
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                handleUninstall(addon);
                              }}
                              disabled={isInstalling}
                              className="px-4 py-2.5 rounded-full font-semibold text-sm bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 hover:bg-red-200 dark:hover:bg-red-900/50 transition-all"
                            >
                              Remove
                            </button>
                          )}
                        </div>
                      </div>
                    </motion.div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Addon Detail Modal */}
      <AnimatePresence>
        {selectedAddon && (() => {
          const addon = selectedAddon;
          const isInstalled = addon.status.installed;
          const isInstalling = installing === addon.manifest.id;

          return (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 bg-black/60 backdrop-blur-md flex items-center justify-center z-50 p-4"
              onClick={() => setSelectedAddon(null)}
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
                        <div className={`w-24 h-24 bg-gradient-to-br ${getCategoryColor(addon.manifest.category)} rounded-3xl flex items-center justify-center shadow-xl text-5xl`}>
                          {addon.manifest.icon}
                        </div>
                        <div>
                          <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-1">
                            {addon.manifest.name}
                          </h2>
                          <p className="text-gray-500 dark:text-gray-400 mb-2">
                            {addon.manifest.author}
                          </p>
                          <p className="text-sm text-gray-400 dark:text-gray-500">
                            Version {addon.manifest.version}
                          </p>
                        </div>
                      </div>
                      <button
                        onClick={() => setSelectedAddon(null)}
                        className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors p-2"
                      >
                        <X className="w-6 h-6" />
                      </button>
                    </div>

                    {/* Stats Row */}
                    <div className="flex items-center gap-6 pb-6 border-b border-gray-200 dark:border-gray-700">
                      <span className="px-3 py-1.5 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg text-sm font-medium">
                        {addon.manifest.category}
                      </span>
                      {isInstalled && (
                        <span className="px-3 py-1.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded-lg text-sm font-medium">
                          ✓ Installed
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Content */}
                  <div className="px-8 pb-6">
                    {/* Description */}
                    <div className="mb-6">
                      <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-3">
                        About this app
                      </h3>
                      <p className="text-gray-600 dark:text-gray-400 leading-relaxed mb-4">
                        {addon.manifest.description}
                      </p>
                    </div>

                    {/* System Requirements */}
                    {addon.manifest.system_packages.length > 0 && (
                      <div className="mb-6">
                        <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-3">
                          System Packages
                        </h3>
                        <div className="flex flex-wrap gap-2">
                          {addon.manifest.system_packages.map(pkg => (
                            <span
                              key={pkg}
                              className="px-3 py-1.5 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg text-sm font-mono"
                            >
                              {pkg}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Requirements */}
                    <div className="mb-6">
                      <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-3">
                        Requirements
                      </h3>
                      <div className="grid grid-cols-2 gap-4">
                        <div className="p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
                          <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Minimum Memory</p>
                          <p className="text-lg font-semibold text-gray-900 dark:text-white">
                            {addon.manifest.minimum_memory} MB
                          </p>
                        </div>
                        <div className="p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
                          <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Minimum Disk</p>
                          <p className="text-lg font-semibold text-gray-900 dark:text-white">
                            {addon.manifest.minimum_disk} GB
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Action Footer */}
                  <div className="px-8 py-6 bg-gray-50 dark:bg-macos-dark-50 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center justify-between">
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        {isInstalled ? 'This app is installed on your NAS' : 'Click GET to install this app'}
                      </div>
                      <div className="flex gap-3">
                        {isInstalled && (
                          <button
                            onClick={() => handleUninstall(addon)}
                            disabled={isInstalling}
                            className="px-6 py-2.5 bg-red-100 dark:bg-red-900/30 hover:bg-red-200 dark:hover:bg-red-900/50 text-red-600 dark:text-red-400 rounded-full transition-colors font-semibold text-sm"
                          >
                            Uninstall
                          </button>
                        )}
                        <button
                          onClick={() => {
                            if (!isInstalled && !isInstalling) {
                              handleInstall(addon);
                            }
                          }}
                          disabled={isInstalled || isInstalling}
                          className={`px-8 py-2.5 rounded-full font-semibold text-sm transition-all ${
                            isInstalled
                              ? 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-default'
                              : isInstalling
                              ? 'bg-blue-400 text-white cursor-wait'
                              : 'bg-macos-blue hover:bg-blue-600 text-white shadow-lg hover:shadow-xl'
                          }`}
                        >
                          {isInstalled ? 'Installed' : isInstalling ? 'Installing...' : 'GET'}
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
