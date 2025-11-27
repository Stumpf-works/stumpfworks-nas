import { useState } from 'react';
import { Shield, FileText, AlertTriangle, Activity } from 'lucide-react';
import { AuditLogs } from '../AuditLogs/AuditLogs';
import { Security } from '../Security/Security';
import { Alerts } from '../Alerts/Alerts';

type Tab = 'audit' | 'security' | 'alerts';

export function SecurityCenter() {
  const [activeTab, setActiveTab] = useState<Tab>('audit');

  const tabs = [
    { id: 'audit' as Tab, label: 'Audit Logs', icon: FileText },
    { id: 'security' as Tab, label: 'Failed Logins', icon: Shield },
    { id: 'alerts' as Tab, label: 'Alerts', icon: AlertTriangle },
  ];

  return (
    <div className="h-full flex flex-col bg-gray-50 dark:bg-macos-dark-100">
      {/* Header */}
      <div className="bg-white dark:bg-macos-dark-200 border-b border-gray-200 dark:border-gray-700">
        <div className="px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-500/10 rounded-lg">
              <Activity className="w-6 h-6 text-blue-500" />
            </div>
            <div>
              <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
                Security Center
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Monitor security events, audit logs, and system alerts
              </p>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 px-6">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            const isActive = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`
                  flex items-center gap-2 px-4 py-3 rounded-t-lg font-medium transition-colors
                  ${
                    isActive
                      ? 'bg-gray-50 dark:bg-macos-dark-100 text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 dark:border-blue-400'
                      : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 hover:bg-gray-50 dark:hover:bg-macos-dark-300'
                  }
                `}
              >
                <Icon className="w-4 h-4" />
                {tab.label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        {activeTab === 'audit' && <AuditLogs />}
        {activeTab === 'security' && <Security />}
        {activeTab === 'alerts' && <Alerts />}
      </div>
    </div>
  );
}
