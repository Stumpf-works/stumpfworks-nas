import { useState } from 'react';
import { motion } from 'framer-motion';
import ContainerManager from './components/ContainerManager';
import ImageManager from './components/ImageManager';
import VolumeManager from './components/VolumeManager';
import NetworkManager from './components/NetworkManager';
import StackManager from './components/StackManager';

type Tab = 'containers' | 'images' | 'volumes' | 'networks' | 'stacks';

export function DockerManager() {
  const [activeTab, setActiveTab] = useState<Tab>('containers');

  const tabs = [
    { id: 'containers' as Tab, name: 'Containers', icon: '📦' },
    { id: 'images' as Tab, name: 'Images', icon: '💿' },
    { id: 'volumes' as Tab, name: 'Volumes', icon: '💾' },
    { id: 'networks' as Tab, name: 'Networks', icon: '🌐' },
    { id: 'stacks' as Tab, name: 'Stacks', icon: '📚' },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Docker Manager
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Manage Docker containers, images, volumes, networks, and compose stacks
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
                layoutId="dockerActiveTab"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-macos-blue"
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
              />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto">
        {activeTab === 'containers' && <ContainerManager />}
        {activeTab === 'images' && <ImageManager />}
        {activeTab === 'volumes' && <VolumeManager />}
        {activeTab === 'networks' && <NetworkManager />}
        {activeTab === 'stacks' && <StackManager />}
      </div>
    </div>
  );
}
