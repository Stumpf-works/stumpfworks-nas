import { useState } from 'react';
import { motion } from 'framer-motion';
import DRBDPanel from './tabs/DRBDPanel';

type HATab = 'drbd' | 'cluster' | 'vip';

export function HighAvailability() {
  const [activeTab, setActiveTab] = useState<HATab>('drbd');

  const tabs = [
    { id: 'drbd' as HATab, name: 'DRBD Replication', icon: '💿', enabled: true },
    { id: 'cluster' as HATab, name: 'Cluster (Pacemaker)', icon: '🔗', enabled: false },
    { id: 'vip' as HATab, name: 'Virtual IP (Keepalived)', icon: '🌐', enabled: false },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          High Availability
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure DRBD replication, cluster management, and virtual IPs for high availability
        </p>
      </div>

      {/* Tabs */}
      <div className="flex items-center px-6 bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => tab.enabled && setActiveTab(tab.id)}
            disabled={!tab.enabled}
            className={`flex items-center gap-2 px-4 py-3 font-medium text-sm whitespace-nowrap transition-colors relative ${
              activeTab === tab.id
                ? 'text-macos-blue dark:text-macos-blue'
                : tab.enabled
                ? 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
                : 'text-gray-400 dark:text-gray-600 cursor-not-allowed'
            }`}
          >
            <span className="text-lg">{tab.icon}</span>
            {tab.name}
            {!tab.enabled && (
              <span className="ml-1 text-xs bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400 px-2 py-0.5 rounded">
                Coming Soon
              </span>
            )}
            {activeTab === tab.id && tab.enabled && (
              <motion.div
                layoutId="haActiveTab"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-macos-blue"
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
              />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto">
        {activeTab === 'drbd' && <DRBDPanel />}
        {activeTab === 'cluster' && (
          <div className="p-6 text-center text-gray-500 dark:text-gray-400">
            Pacemaker/Corosync cluster management coming soon...
          </div>
        )}
        {activeTab === 'vip' && (
          <div className="p-6 text-center text-gray-500 dark:text-gray-400">
            Keepalived virtual IP management coming soon...
          </div>
        )}
      </div>
    </div>
  );
}
