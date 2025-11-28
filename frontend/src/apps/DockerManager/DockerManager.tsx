import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { dockerApi } from '@/api/docker';
import { getErrorMessage } from '@/api/client';
import ContainerManager from './components/ContainerManager';
import ImageManager from './components/ImageManager';
import VolumeManager from './components/VolumeManager';
import NetworkManager from './components/NetworkManager';
import StackManager from './components/StackManager';
import Card from '@/components/ui/Card';
import { Disc, HardDrive, Network, Layers, Container as ContainerIcon, AlertCircle, RefreshCw } from 'lucide-react';

type Tab = 'containers' | 'images' | 'volumes' | 'networks' | 'stacks';

export function DockerManager() {
  const [activeTab, setActiveTab] = useState<Tab>('containers');
  const [dockerAvailable, setDockerAvailable] = useState<boolean | null>(null);
  const [dockerError, setDockerError] = useState<string>('');

  useEffect(() => {
    checkDockerAvailability();
  }, []);

  const checkDockerAvailability = async () => {
    try {
      const response = await dockerApi.listContainers();
      if (response.success) {
        setDockerAvailable(true);
      } else {
        setDockerAvailable(false);
        setDockerError(response.error?.message || 'Docker is not available');
      }
    } catch (err) {
      setDockerAvailable(false);
      setDockerError(getErrorMessage(err));
    }
  };

  const tabs = [
    { id: 'containers' as Tab, name: 'Containers', icon: ContainerIcon },
    { id: 'images' as Tab, name: 'Images', icon: Disc },
    { id: 'volumes' as Tab, name: 'Volumes', icon: HardDrive },
    { id: 'networks' as Tab, name: 'Networks', icon: Network },
    { id: 'stacks' as Tab, name: 'Stacks', icon: Layers },
  ];

  // Show loading state
  if (dockerAvailable === null) {
    return (
      <div className="flex flex-col items-center justify-center h-full bg-gray-50 dark:bg-macos-dark-50">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue"></div>
        <p className="mt-4 text-gray-600 dark:text-gray-400">Checking Docker availability...</p>
      </div>
    );
  }

  // Show Docker not available message
  if (dockerAvailable === false) {
    return (
      <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50 p-4 md:p-6">
        <div className="max-w-2xl mx-auto mt-12 w-full">
          <Card>
            <div className="p-6 md:p-8 text-center">
              <div className="flex justify-center mb-4">
                <ContainerIcon className="w-16 h-16 md:w-24 md:h-24 text-gray-400 dark:text-gray-600" />
              </div>
              <h2 className="text-xl md:text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                Docker Not Available
              </h2>
              <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mb-6">
                Docker is not running or not installed on this system.
              </p>
              <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4 mb-6">
                <div className="flex items-start gap-2">
                  <AlertCircle className="w-5 h-5 text-yellow-600 dark:text-yellow-400 flex-shrink-0 mt-0.5" />
                  <p className="text-xs md:text-sm text-yellow-800 dark:text-yellow-200 text-left break-words">
                    {dockerError}
                  </p>
                </div>
              </div>
              <div className="text-left space-y-3">
                <p className="text-sm md:text-base text-gray-700 dark:text-gray-300 font-semibold">
                  To use Docker features:
                </p>
                <ol className="list-decimal list-inside text-xs md:text-sm text-gray-600 dark:text-gray-400 space-y-2">
                  <li>Install Docker on your NAS system</li>
                  <li>
                    Start the Docker daemon:{' '}
                    <code className="bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded text-xs">
                      sudo systemctl start docker
                    </code>
                  </li>
                  <li>
                    Enable Docker on boot:{' '}
                    <code className="bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded text-xs">
                      sudo systemctl enable docker
                    </code>
                  </li>
                  <li>Refresh this page</li>
                </ol>
              </div>
              <button
                onClick={checkDockerAvailability}
                className="mt-6 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors text-sm md:text-base flex items-center gap-2 mx-auto"
              >
                <RefreshCw className="w-4 h-4" />
                Check Again
              </button>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-4 md:p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <h1 className="text-xl md:text-2xl font-bold text-gray-900 dark:text-gray-100">
          Docker Manager
        </h1>
        <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mt-1">
          Manage Docker containers, images, volumes, networks, and compose stacks
        </p>
      </div>

      {/* Tabs - Responsive */}
      <div className="flex items-center px-2 md:px-6 bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex items-center gap-1 md:gap-2 px-2 md:px-4 py-3 font-medium text-xs md:text-sm whitespace-nowrap transition-colors relative ${
                activeTab === tab.id
                  ? 'text-macos-blue dark:text-macos-blue'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
              }`}
            >
              <Icon className="w-4 h-4 md:w-5 md:h-5" />
              <span className="hidden sm:inline">{tab.name}</span>
              {activeTab === tab.id && (
                <motion.div
                  layoutId="dockerActiveTab"
                  className="absolute bottom-0 left-0 right-0 h-0.5 bg-macos-blue"
                  transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                />
              )}
            </button>
          );
        })}
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
