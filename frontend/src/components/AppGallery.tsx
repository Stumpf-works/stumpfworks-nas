import { useState, useMemo, memo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Search } from 'lucide-react';
import { registeredApps, appCategories, categoryNames, categoryIcons, type AppCategory } from '../apps';
import type { App } from '../types';

interface AppGalleryProps {
  isOpen: boolean;
  onClose: () => void;
  onLaunchApp: (appId: string) => void;
}

export const AppGallery = memo(function AppGallery({ isOpen, onClose, onLaunchApp }: AppGalleryProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<AppCategory | 'all'>('all');

  // Filter apps based on search and category
  const filteredApps = useMemo(() => {
    let apps = registeredApps;

    // Filter by category
    if (selectedCategory !== 'all') {
      const categoryAppIds = appCategories[selectedCategory] as readonly string[];
      apps = apps.filter((app) => categoryAppIds.includes(app.id));
    }

    // Filter by search query
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      apps = apps.filter(
        (app) =>
          app.name.toLowerCase().includes(query) ||
          app.id.toLowerCase().includes(query)
      );
    }

    return apps;
  }, [searchQuery, selectedCategory]);

  // Group apps by category for display
  const appsByCategory = useMemo(() => {
    const grouped: Record<AppCategory, App[]> = {
      system: [],
      management: [],
      security: [],
      tools: [],
      development: [],
    };

    filteredApps.forEach((app) => {
      // Find which category this app belongs to
      for (const [category, appIds] of Object.entries(appCategories)) {
        if ((appIds as readonly string[]).includes(app.id)) {
          grouped[category as AppCategory].push(app);
          break;
        }
      }
    });

    return grouped;
  }, [filteredApps]);

  const handleAppClick = (appId: string) => {
    onLaunchApp(appId);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 z-[100] bg-black/80 backdrop-blur-md"
        onClick={onClose}
      >
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          transition={{ type: 'spring', stiffness: 300, damping: 30 }}
          onClick={(e) => e.stopPropagation()}
          className="h-full overflow-auto p-8"
        >
          {/* Header */}
          <div className="max-w-7xl mx-auto mb-8">
            <div className="flex items-center justify-between mb-6">
              <h1 className="text-3xl font-bold text-white">App Gallery</h1>
              <button
                onClick={onClose}
                className="p-2 text-white/70 hover:text-white hover:bg-white/10 rounded-lg transition-colors"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            {/* Search Bar */}
            <div className="relative mb-6">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search apps..."
                className="w-full pl-12 pr-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                autoFocus
              />
            </div>

            {/* Category Filter */}
            <div className="flex gap-2 overflow-x-auto pb-2">
              <button
                onClick={() => setSelectedCategory('all')}
                className={`px-4 py-2 rounded-lg font-medium whitespace-nowrap transition-colors ${
                  selectedCategory === 'all'
                    ? 'bg-macos-blue text-white'
                    : 'bg-white/10 text-white/70 hover:bg-white/20'
                }`}
              >
                All Apps
              </button>
              {Object.entries(categoryNames).map(([key, name]) => (
                <button
                  key={key}
                  onClick={() => setSelectedCategory(key as AppCategory)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium whitespace-nowrap transition-colors ${
                    selectedCategory === key
                      ? 'bg-macos-blue text-white'
                      : 'bg-white/10 text-white/70 hover:bg-white/20'
                  }`}
                >
                  <span>{categoryIcons[key as AppCategory]}</span>
                  {name}
                </button>
              ))}
            </div>
          </div>

          {/* Apps Grid */}
          <div className="max-w-7xl mx-auto">
            {selectedCategory === 'all' ? (
              // Show all categories
              <>
                {Object.entries(appsByCategory).map(([category, apps]) => {
                  if (apps.length === 0) return null;
                  return (
                    <div key={category} className="mb-12">
                      <h2 className="text-xl font-semibold text-white/90 mb-4 flex items-center gap-2">
                        <span className="text-2xl">{categoryIcons[category as AppCategory]}</span>
                        {categoryNames[category as AppCategory]}
                      </h2>
                      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                        {apps.map((app) => (
                          <AppIcon
                            key={app.id}
                            app={app}
                            onClick={() => handleAppClick(app.id)}
                          />
                        ))}
                      </div>
                    </div>
                  );
                })}
              </>
            ) : (
              // Show single category
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                {filteredApps.map((app) => (
                  <AppIcon
                    key={app.id}
                    app={app}
                    onClick={() => handleAppClick(app.id)}
                  />
                ))}
              </div>
            )}

            {/* No Results */}
            {filteredApps.length === 0 && (
              <div className="text-center py-20">
                <div className="text-6xl mb-4">🔍</div>
                <p className="text-white/60 text-lg">No apps found</p>
                <p className="text-white/40 text-sm mt-2">
                  Try adjusting your search or category filter
                </p>
              </div>
            )}
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
});

interface AppIconProps {
  app: App;
  onClick: () => void;
}

const AppIcon = memo(function AppIcon({ app, onClick }: AppIconProps) {
  return (
    <motion.button
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      onClick={onClick}
      className="flex flex-col items-center gap-3 p-4 bg-white/5 hover:bg-white/10 rounded-2xl transition-colors group"
    >
      <div className="text-5xl group-hover:scale-110 transition-transform">{app.icon}</div>
      <span className="text-sm font-medium text-white/90 text-center line-clamp-2">
        {app.name}
      </span>
    </motion.button>
  );
});
