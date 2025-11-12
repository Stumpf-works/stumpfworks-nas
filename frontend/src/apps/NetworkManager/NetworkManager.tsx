import { useState } from 'react';
import { motion } from 'framer-motion';
import InterfaceManager from './components/InterfaceManager';
import DNSSettings from './components/DNSSettings';
import FirewallManager from './components/FirewallManager';
import DiagnosticsTool from './components/DiagnosticsTool';
import BandwidthMonitor from './components/BandwidthMonitor';

type Tab = 'interfaces' | 'dns' | 'firewall' | 'diagnostics' | 'bandwidth';

export function NetworkManager() {
  const [activeTab, setActiveTab] = useState<Tab>('interfaces');

  const tabs = [
    { id: 'interfaces' as Tab, name: 'Interfaces', icon: '🌐' },
    { id: 'dns' as Tab, name: 'DNS & Routes', icon: '📡' },
    { id: 'firewall' as Tab, name: 'Firewall', icon: '🛡️' },
    { id: 'diagnostics' as Tab, name: 'Diagnostics', icon: '🔍' },
    { id: 'bandwidth' as Tab, name: 'Bandwidth', icon: '📊' },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Network Manager
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure network interfaces, firewall, and diagnostics
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
        {activeTab === 'interfaces' && <InterfaceManager />}
        {activeTab === 'dns' && <DNSSettings />}
        {activeTab === 'firewall' && <FirewallManager />}
        {activeTab === 'diagnostics' && <DiagnosticsTool />}
        {activeTab === 'bandwidth' && <BandwidthMonitor />}
      </div>
    </div>
  );
}
