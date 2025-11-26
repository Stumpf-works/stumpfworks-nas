import { useState } from 'react';
import { motion } from 'framer-motion';
import DomainStatus from './components/DomainStatus';
import UserManagement from './components/UserManagement';
import GroupManagement from './components/GroupManagement';
import ComputerManagement from './components/ComputerManagement';
import OUManagement from './components/OUManagement';
import GPOManagement from './components/GPOManagement';
import DNSManagement from './components/DNSManagement';
import FSMORoles from './components/FSMORoles';

type Tab = 'status' | 'users' | 'groups' | 'computers' | 'ous' | 'gpos' | 'dns' | 'fsmo';

export function ADDCManager() {
  const [activeTab, setActiveTab] = useState<Tab>('status');

  const tabs = [
    { id: 'status' as Tab, name: 'Domain Status', icon: '🌐' },
    { id: 'users' as Tab, name: 'Users', icon: '👤' },
    { id: 'groups' as Tab, name: 'Groups', icon: '👥' },
    { id: 'computers' as Tab, name: 'Computers', icon: '💻' },
    { id: 'ous' as Tab, name: 'OUs', icon: '📁' },
    { id: 'gpos' as Tab, name: 'Group Policies', icon: '📋' },
    { id: 'dns' as Tab, name: 'DNS', icon: '🌍' },
    { id: 'fsmo' as Tab, name: 'FSMO Roles', icon: '⚙️' },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-macos-dark-100">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Active Directory Domain Controller
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Manage your Samba Active Directory Domain Controller
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
                layoutId="activeDCTab"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-macos-blue"
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
              />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto">
        {activeTab === 'status' && <DomainStatus />}
        {activeTab === 'users' && <UserManagement />}
        {activeTab === 'groups' && <GroupManagement />}
        {activeTab === 'computers' && <ComputerManagement />}
        {activeTab === 'ous' && <OUManagement />}
        {activeTab === 'gpos' && <GPOManagement />}
        {activeTab === 'dns' && <DNSManagement />}
        {activeTab === 'fsmo' && <FSMORoles />}
      </div>
    </div>
  );
}
