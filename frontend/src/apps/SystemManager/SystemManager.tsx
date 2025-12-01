import { useState } from 'react';
import { motion } from 'framer-motion';
import MonitoringDashboard from './tabs/MonitoringDashboard';
import ScheduledTasks from './tabs/ScheduledTasks';
import BackupManager from './tabs/BackupManager';

type SystemTab = 'monitoring' | 'tasks' | 'backups';

export function SystemManager() {
  const [activeTab, setActiveTab] = useState<SystemTab>('monitoring');

  const tabs = [
    { id: 'monitoring' as SystemTab, name: 'Monitoring', icon: '📊' },
    { id: 'tasks' as SystemTab, name: 'Scheduled Tasks', icon: '📅' },
    { id: 'backups' as SystemTab, name: 'Backups', icon: '⏱️' },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          System Management
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Monitor system health, manage scheduled tasks and backups
        </p>
      </div>

      {/* Tabs */}
      <div className="flex items-center px-6 bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`flex items-center gap-2 px-4 py-3 font-medium text-sm whitespace-nowrap transition-colors relative ${
              activeTab === tab.id
                ? 'text-macos-blue dark:text-macos-blue'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
            }`}
          >
            <span className="text-lg">{tab.icon}</span>
            {tab.name}
            {activeTab === tab.id && (
              <motion.div
                layoutId="systemManagerActiveTab"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-macos-blue"
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
              />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto">
        {activeTab === 'monitoring' && <MonitoringDashboard />}
        {activeTab === 'tasks' && <ScheduledTasks />}
        {activeTab === 'backups' && <BackupManager />}
      </div>
    </div>
  );
}
