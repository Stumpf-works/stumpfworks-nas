import { useState } from 'react';
import { motion } from 'framer-motion';
import StorageOverview from './StorageOverview';
import DiskManager from './DiskManager';
import VolumeManager from './VolumeManager';
import ShareManager from './ShareManager';
import ZFSManager from './ZFSManager';
import RAIDManager from './RAIDManager';
import SMARTMonitor from './SMARTMonitor';
import SambaManager from './SambaManager';
import NFSManager from './NFSManager';

type Tab = 'overview' | 'disks' | 'volumes' | 'shares' | 'zfs' | 'raid' | 'smart' | 'samba' | 'nfs';

export function StorageManager() {
  const [activeTab, setActiveTab] = useState<Tab>('overview');

  const tabs = [
    { id: 'overview' as Tab, name: 'Overview', icon: '📊' },
    { id: 'disks' as Tab, name: 'Disks', icon: '💿' },
    { id: 'volumes' as Tab, name: 'Volumes', icon: '📦' },
    { id: 'shares' as Tab, name: 'Shares', icon: '📁' },
    { id: 'zfs' as Tab, name: 'ZFS', icon: '🗄️' },
    { id: 'raid' as Tab, name: 'RAID', icon: '🛡️' },
    { id: 'smart' as Tab, name: 'SMART', icon: '🔍' },
    { id: 'samba' as Tab, name: 'Samba', icon: '🔗' },
    { id: 'nfs' as Tab, name: 'NFS', icon: '🌐' },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Tab Navigation */}
      <div className="flex border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`
              relative px-6 py-3 text-sm font-medium transition-colors
              ${
                activeTab === tab.id
                  ? 'text-macos-blue dark:text-macos-blue'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
              }
            `}
          >
            <span className="mr-2">{tab.icon}</span>
            {tab.name}
            {activeTab === tab.id && (
              <motion.div
                layoutId="activeTab"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-macos-blue"
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
              />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto">
        {activeTab === 'overview' && (
          <div className="p-6">
            <StorageOverview />
          </div>
        )}
        {activeTab === 'disks' && (
          <div className="p-6">
            <DiskManager />
          </div>
        )}
        {activeTab === 'volumes' && (
          <div className="p-6">
            <VolumeManager />
          </div>
        )}
        {activeTab === 'shares' && (
          <div className="p-6">
            <ShareManager />
          </div>
        )}
        {activeTab === 'zfs' && <ZFSManager />}
        {activeTab === 'raid' && <RAIDManager />}
        {activeTab === 'smart' && <SMARTMonitor />}
        {activeTab === 'samba' && <SambaManager />}
        {activeTab === 'nfs' && <NFSManager />}
      </div>
    </div>
  );
}
