import { useState } from 'react';
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
      {/* Tab Navigation - Modern Pill Style */}
      <div className="p-4 bg-gradient-to-r from-gray-50 to-gray-100 dark:from-macos-dark-100 dark:to-macos-dark-200 border-b border-gray-200/50 dark:border-gray-700/50">
        <div className="flex gap-2 overflow-x-auto pb-1">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                relative px-5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 whitespace-nowrap
                ${
                  activeTab === tab.id
                    ? 'bg-gradient-to-r from-macos-blue to-macos-purple text-white shadow-lg shadow-macos-blue/30 scale-105'
                    : 'bg-white/50 dark:bg-macos-dark-200/50 text-gray-700 dark:text-gray-300 hover:bg-white dark:hover:bg-macos-dark-200 hover:shadow-md backdrop-blur-sm'
                }
              `}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.name}
            </button>
          ))}
        </div>
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
