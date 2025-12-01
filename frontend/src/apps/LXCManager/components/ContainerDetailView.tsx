import { useState } from 'react';
import { motion } from 'framer-motion';
import {
  Activity,
  Terminal,
  Settings,
  Camera,
  Database,
  Play,
  Square,
  RefreshCw,
  Trash2,
  Copy,
  Download,
  Info,
  Cpu,
  MemoryStick,
  Network,
  Clock
} from 'lucide-react';
import { type Container } from '@/api/lxc';
import Card from '@/components/ui/Card';
import { WebTerminal } from './WebTerminal';

interface ContainerDetailViewProps {
  container: Container;
  onAction: (action: 'start' | 'stop' | 'restart' | 'delete', name: string) => Promise<void>;
  onClose: () => void;
  loading?: boolean;
}

type TabType = 'summary' | 'console' | 'resources' | 'options' | 'snapshots' | 'backup';

export function ContainerDetailView({ container, onAction, onClose, loading }: ContainerDetailViewProps) {
  const [activeTab, setActiveTab] = useState<TabType>('summary');

  const isRunning = container.state.toLowerCase() === 'running';

  const tabs = [
    { id: 'summary' as TabType, label: 'Summary', icon: Info, enabled: true },
    { id: 'console' as TabType, label: 'Console', icon: Terminal, enabled: isRunning },
    { id: 'resources' as TabType, label: 'Resources', icon: Activity, enabled: isRunning },
    { id: 'options' as TabType, label: 'Options', icon: Settings, enabled: true },
    { id: 'snapshots' as TabType, label: 'Snapshots', icon: Camera, enabled: true },
    { id: 'backup' as TabType, label: 'Backup', icon: Database, enabled: true },
  ];

  const getStateColor = () => {
    switch (container.state.toLowerCase()) {
      case 'running': return 'bg-green-500';
      case 'stopped': return 'bg-gray-500';
      case 'frozen': return 'bg-blue-500';
      default: return 'bg-gray-500';
    }
  };

  const getStateBadge = () => {
    const baseClasses = 'px-3 py-1 rounded-full text-sm font-semibold flex items-center gap-2';
    switch (container.state.toLowerCase()) {
      case 'running':
        return `${baseClasses} bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300`;
      case 'stopped':
        return `${baseClasses} bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300`;
      case 'frozen':
        return `${baseClasses} bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300`;
      default:
        return `${baseClasses} bg-gray-100 dark:bg-gray-800/30 text-gray-700 dark:text-gray-300`;
    }
  };

  const renderSummaryTab = () => (
    <div className="space-y-6">
      {/* Status Overview */}
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Status Overview</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                <Cpu className="w-4 h-4" />
                CPU Usage
              </div>
              <div className="text-2xl font-bold text-gray-900 dark:text-white">
                {isRunning ? (container.cpu_usage || '0%') : 'N/A'}
              </div>
            </div>
            <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                <MemoryStick className="w-4 h-4" />
                Memory
              </div>
              <div className="text-2xl font-bold text-gray-900 dark:text-white">
                {isRunning ? (container.memory || '0 MB') : 'N/A'}
              </div>
              {container.memory_limit && (
                <div className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                  of {container.memory_limit}
                </div>
              )}
            </div>
            <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                <Network className="w-4 h-4" />
                Network
              </div>
              <div className="text-lg font-bold text-gray-900 dark:text-white">
                {container.ipv4 || 'No IP'}
              </div>
            </div>
            <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                <Clock className="w-4 h-4" />
                Uptime
              </div>
              <div className="text-lg font-bold text-gray-900 dark:text-white">
                {isRunning ? 'Running' : 'Stopped'}
              </div>
              {isRunning && container.pid > 0 && (
                <div className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                  PID: {container.pid}
                </div>
              )}
            </div>
          </div>
        </div>
      </Card>

      {/* Container Information */}
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Container Information</h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">Container ID</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">{container.name}</span>
            </div>
            <div className="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">Status</span>
              <span className={getStateBadge()}>
                <span className={`w-2 h-2 rounded-full ${getStateColor()}`} />
                {container.state}
              </span>
            </div>
            {('template' in container) && (
              <div className="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-700">
                <span className="text-sm text-gray-600 dark:text-gray-400">Template</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{(container as any).template || 'N/A'}</span>
              </div>
            )}
            <div className="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">IPv4 Address</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">{container.ipv4 || 'N/A'}</span>
            </div>
            <div className="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">IPv6 Address</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">{container.ipv6 || 'N/A'}</span>
            </div>
            <div className="flex items-center justify-between py-2">
              <span className="text-sm text-gray-600 dark:text-gray-400">Autostart</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {container.autostart ? 'Enabled' : 'Disabled'}
              </span>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );

  const renderConsoleTab = () => (
    <div className="h-full flex flex-col">
      {isRunning ? (
        <div className="flex-1 min-h-[500px]">
          <WebTerminal containerName={container.name} />
        </div>
      ) : (
        <Card>
          <div className="p-6 text-center py-12 text-gray-500 dark:text-gray-400">
            Container must be running to access console
          </div>
        </Card>
      )}
    </div>
  );

  const renderResourcesTab = () => (
    <div className="space-y-6">
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Resource Usage</h3>
          {isRunning ? (
            <div className="space-y-6">
              {/* CPU Usage */}
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">CPU Usage</span>
                  <span className="text-sm font-bold text-gray-900 dark:text-white">{container.cpu_usage || '0%'}</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
                  <div
                    className="bg-blue-500 h-3 rounded-full transition-all duration-300"
                    style={{ width: container.cpu_usage || '0%' }}
                  />
                </div>
              </div>

              {/* Memory Usage */}
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Memory Usage</span>
                  <span className="text-sm font-bold text-gray-900 dark:text-white">
                    {container.memory} / {container.memory_limit}
                  </span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
                  <div
                    className="bg-green-500 h-3 rounded-full transition-all duration-300"
                    style={{ width: '45%' }}
                  />
                </div>
              </div>

              {/* Network */}
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Network</span>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-3">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">RX (Incoming)</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-white">0 MB/s</div>
                  </div>
                  <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-3">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">TX (Outgoing)</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-white">0 MB/s</div>
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="text-center py-12 text-gray-500 dark:text-gray-400">
              Container must be running to view resources
            </div>
          )}
        </div>
      </Card>
    </div>
  );

  const renderOptionsTab = () => (
    <Card>
      <div className="p-6">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Container Options</h3>
        <div className="space-y-4">
          <div className="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
            <div>
              <div className="font-medium text-gray-900 dark:text-white">Start at boot</div>
              <div className="text-sm text-gray-600 dark:text-gray-400">Automatically start this container on system boot</div>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" checked={container.autostart} className="sr-only peer" disabled />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
            </label>
          </div>
          <div className="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
            <div>
              <div className="font-medium text-gray-900 dark:text-white">Memory Limit</div>
              <div className="text-sm text-gray-600 dark:text-gray-400">Maximum memory allocation</div>
            </div>
            <div className="text-sm font-medium text-gray-900 dark:text-white">{container.memory_limit || 'Unlimited'}</div>
          </div>
          {('template' in container) && (
            <div className="flex items-center justify-between py-3">
              <div>
                <div className="font-medium text-gray-900 dark:text-white">Template</div>
                <div className="text-sm text-gray-600 dark:text-gray-400">OS template used</div>
              </div>
              <div className="text-sm font-medium text-gray-900 dark:text-white">{(container as any).template || 'N/A'}</div>
            </div>
          )}
        </div>
      </div>
    </Card>
  );

  const renderSnapshotsTab = () => (
    <Card>
      <div className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Snapshots</h3>
          <button className="px-3 py-1.5 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors text-sm flex items-center gap-2">
            <Camera className="w-4 h-4" />
            Take Snapshot
          </button>
        </div>
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          <Camera className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No snapshots available</p>
          <p className="text-sm mt-2">Create your first snapshot to save the container state</p>
        </div>
      </div>
    </Card>
  );

  const renderBackupTab = () => (
    <Card>
      <div className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Backup & Restore</h3>
          <button className="px-3 py-1.5 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors text-sm flex items-center gap-2">
            <Download className="w-4 h-4" />
            Create Backup
          </button>
        </div>
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          <Database className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No backups available</p>
          <p className="text-sm mt-2">Create a backup to protect your container data</p>
        </div>
      </div>
    </Card>
  );

  return (
    <div className="h-full flex flex-col bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="border-b border-gray-200 dark:border-gray-700 p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            >
              ←
            </button>
            <div>
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white">{container.name}</h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">Container {container.name}</p>
            </div>
            <span className={getStateBadge()}>
              <span className={`w-2 h-2 rounded-full ${getStateColor()}`} />
              {container.state}
            </span>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-2">
            {isRunning ? (
              <>
                <button
                  onClick={() => onAction('stop', container.name)}
                  disabled={loading}
                  className="px-4 py-2 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors flex items-center gap-2 disabled:opacity-50"
                >
                  {loading ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Square className="w-4 h-4" />}
                  Stop
                </button>
                <button
                  onClick={() => onAction('restart', container.name)}
                  disabled={loading}
                  className="px-4 py-2 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg hover:bg-blue-200 dark:hover:bg-blue-900/50 transition-colors flex items-center gap-2 disabled:opacity-50"
                >
                  {loading ? <RefreshCw className="w-4 h-4 animate-spin" /> : <RefreshCw className="w-4 h-4" />}
                  Restart
                </button>
              </>
            ) : (
              <button
                onClick={() => onAction('start', container.name)}
                disabled={loading}
                className="px-4 py-2 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300 rounded-lg hover:bg-green-200 dark:hover:bg-green-900/50 transition-colors flex items-center gap-2 disabled:opacity-50"
              >
                {loading ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Play className="w-4 h-4" />}
                Start
              </button>
            )}
            <button className="px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center gap-2">
              <Copy className="w-4 h-4" />
              Clone
            </button>
            <button
              onClick={() => onAction('delete', container.name)}
              disabled={loading}
              className="px-4 py-2 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors flex items-center gap-2 disabled:opacity-50"
            >
              <Trash2 className="w-4 h-4" />
              Delete
            </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 border-b border-gray-200 dark:border-gray-700 -mb-px">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            const isActive = activeTab === tab.id;
            const isEnabled = tab.enabled;

            return (
              <button
                key={tab.id}
                onClick={() => isEnabled && setActiveTab(tab.id)}
                disabled={!isEnabled}
                className={`
                  flex items-center gap-2 px-4 py-2 border-b-2 transition-colors
                  ${isActive
                    ? 'border-macos-blue text-macos-blue'
                    : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                  }
                  ${!isEnabled && 'opacity-50 cursor-not-allowed'}
                `}
              >
                <Icon className="w-4 h-4" />
                {tab.label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto p-6">
        <motion.div
          key={activeTab}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
        >
          {activeTab === 'summary' && renderSummaryTab()}
          {activeTab === 'console' && renderConsoleTab()}
          {activeTab === 'resources' && renderResourcesTab()}
          {activeTab === 'options' && renderOptionsTab()}
          {activeTab === 'snapshots' && renderSnapshotsTab()}
          {activeTab === 'backup' && renderBackupTab()}
        </motion.div>
      </div>
    </div>
  );
}
