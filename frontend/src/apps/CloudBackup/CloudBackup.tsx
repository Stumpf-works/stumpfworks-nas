import { useState, useEffect } from 'react';
import { cloudBackupApi, CloudProvider, CloudSyncJob, CloudSyncLog } from '@/api/cloudbackup';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { Cloud, Plus, Play, Trash2, Edit, CheckCircle, AlertTriangle, Clock, Upload, Download, RefreshCw, Server } from 'lucide-react';
import { ProviderModal } from './components/ProviderModal';
import { JobModal } from './components/JobModal';

export default function CloudBackup() {
  const [providers, setProviders] = useState<CloudProvider[]>([]);
  const [jobs, setJobs] = useState<CloudSyncJob[]>([]);
  const [logs, setLogs] = useState<CloudSyncLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<'providers' | 'jobs' | 'logs'>('providers');

  // Modals
  const [showProviderModal, setShowProviderModal] = useState(false);
  const [showJobModal, setShowJobModal] = useState(false);
  const [editingProvider, setEditingProvider] = useState<CloudProvider | null>(null);
  const [editingJob, setEditingJob] = useState<CloudSyncJob | null>(null);

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async () => {
    setLoading(true);
    try {
      const providersRes = await cloudBackupApi.listProviders();

      if (providersRes.success && providersRes.data) {
        setProviders(providersRes.data);
      }

      if (activeTab === 'jobs') {
        const jobsRes = await cloudBackupApi.listJobs();
        if (jobsRes.success && jobsRes.data) {
          setJobs(jobsRes.data);
        }
      } else if (activeTab === 'logs') {
        const logsRes = await cloudBackupApi.getLogs(undefined, 50);
        if (logsRes.success && logsRes.data) {
          setLogs(logsRes.data.logs);
        }
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleTestProvider = async (id: number) => {
    try {
      const response = await cloudBackupApi.testProvider(id);
      if (response.success) {
        alert('Connection test successful!');
        loadData();
      } else {
        alert(response.error?.message || 'Connection test failed');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleDeleteProvider = async (id: number) => {
    if (!confirm('Are you sure you want to delete this provider?')) return;

    try {
      const response = await cloudBackupApi.deleteProvider(id);
      if (response.success) {
        loadData();
      } else {
        alert(response.error?.message || 'Failed to delete provider');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleRunJob = async (id: number) => {
    try {
      const response = await cloudBackupApi.runJob(id);
      if (response.success) {
        alert('Sync job started successfully');
        loadData();
      } else {
        alert(response.error?.message || 'Failed to start sync job');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleDeleteJob = async (id: number) => {
    if (!confirm('Are you sure you want to delete this sync job?')) return;

    try {
      const response = await cloudBackupApi.deleteJob(id);
      if (response.success) {
        loadData();
      } else {
        alert(response.error?.message || 'Failed to delete job');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'failed':
        return <AlertTriangle className="w-5 h-5 text-red-500" />;
      case 'running':
        return <RefreshCw className="w-5 h-5 text-blue-500 animate-spin" />;
      default:
        return <Clock className="w-5 h-5 text-gray-500" />;
    }
  };

  const getDirectionIcon = (direction: string) => {
    switch (direction) {
      case 'upload':
        return <Upload className="w-4 h-4" />;
      case 'download':
        return <Download className="w-4 h-4" />;
      case 'sync':
        return <RefreshCw className="w-4 h-4" />;
      default:
        return <Cloud className="w-4 h-4" />;
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDuration = (seconds: number) => {
    if (seconds < 60) return `${seconds}s`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ${seconds % 60}s`;
    return `${Math.floor(seconds / 3600)}h ${Math.floor((seconds % 3600) / 60)}m`;
  };

  if (loading && providers.length === 0) {
    return <div className="text-center py-8">Loading cloud backup configuration...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Cloud Backup</h1>
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          Sync your data to cloud storage providers using rclone
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Tabs */}
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-8">
          <button
            onClick={() => setActiveTab('providers')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'providers'
                ? 'border-macos-blue text-macos-blue'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <Server className="w-4 h-4" />
              Providers ({providers.length})
            </div>
          </button>
          <button
            onClick={() => setActiveTab('jobs')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'jobs'
                ? 'border-macos-blue text-macos-blue'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <Cloud className="w-4 h-4" />
              Sync Jobs ({jobs.length})
            </div>
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'logs'
                ? 'border-macos-blue text-macos-blue'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <Clock className="w-4 h-4" />
              Recent Logs
            </div>
          </button>
        </nav>
      </div>

      {/* Providers Tab */}
      {activeTab === 'providers' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Configure cloud storage providers for backups
            </p>
            <Button onClick={() => { setEditingProvider(null); setShowProviderModal(true); }}>
              <Plus className="w-4 h-4 mr-2" />
              Add Provider
            </Button>
          </div>

          {providers.length === 0 ? (
            <Card>
              <div className="p-8 text-center">
                <Cloud className="w-12 h-12 mx-auto text-gray-400 mb-4" />
                <p className="text-gray-600 dark:text-gray-400 mb-4">No cloud providers configured</p>
                <Button onClick={() => setShowProviderModal(true)}>
                  <Plus className="w-4 h-4 mr-2" />
                  Add Your First Provider
                </Button>
              </div>
            </Card>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {providers.map((provider) => (
                <Card key={provider.id}>
                  <div className="p-6">
                    <div className="flex items-start justify-between">
                      <div className="flex items-start gap-3 flex-1">
                        <div className="p-2 bg-blue-100 dark:bg-blue-900/20 rounded-lg">
                          <Cloud className="w-6 h-6 text-blue-600 dark:text-blue-400" />
                        </div>
                        <div className="flex-1">
                          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                            {provider.name}
                          </h3>
                          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                            {provider.type.toUpperCase()}
                          </p>
                          {provider.description && (
                            <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                              {provider.description}
                            </p>
                          )}
                          <div className="flex items-center gap-2 mt-2">
                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                              provider.enabled
                                ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
                                : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-400'
                            }`}>
                              {provider.enabled ? 'Enabled' : 'Disabled'}
                            </span>
                            {provider.test_status && (
                              <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                                provider.test_status === 'success'
                                  ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
                                  : provider.test_status === 'failed'
                                  ? 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
                                  : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-400'
                              }`}>
                                {provider.test_status}
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2 mt-4">
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => handleTestProvider(provider.id)}
                      >
                        Test
                      </Button>
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => { setEditingProvider(provider); setShowProviderModal(true); }}
                      >
                        <Edit className="w-3 h-3" />
                      </Button>
                      <Button
                        variant="danger"
                        size="sm"
                        onClick={() => handleDeleteProvider(provider.id)}
                      >
                        <Trash2 className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Jobs Tab - Simplified for now */}
      {activeTab === 'jobs' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Manage scheduled and manual sync jobs
            </p>
            <Button onClick={() => { setEditingJob(null); setShowJobModal(true); }} disabled={providers.length === 0}>
              <Plus className="w-4 h-4 mr-2" />
              Add Sync Job
            </Button>
          </div>

          {jobs.length === 0 ? (
            <Card>
              <div className="p-8 text-center">
                <RefreshCw className="w-12 h-12 mx-auto text-gray-400 mb-4" />
                <p className="text-gray-600 dark:text-gray-400 mb-4">No sync jobs configured</p>
                {providers.length > 0 && (
                  <Button onClick={() => setShowJobModal(true)}>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Your First Sync Job
                  </Button>
                )}
              </div>
            </Card>
          ) : (
            <div className="space-y-3">
              {jobs.map((job) => (
                <Card key={job.id}>
                  <div className="p-6">
                    <div className="flex items-start justify-between">
                      <div className="flex items-start gap-3 flex-1">
                        {getStatusIcon(job.last_status)}
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                              {job.name}
                            </h3>
                            {getDirectionIcon(job.direction)}
                          </div>
                          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                            {job.description}
                          </p>
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2 mt-4">
                      <Button
                        size="sm"
                        onClick={() => handleRunJob(job.id)}
                        disabled={!job.enabled || job.last_status === 'running'}
                      >
                        <Play className="w-3 h-3 mr-2" />
                        Run Now
                      </Button>
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => { setEditingJob(job); setShowJobModal(true); }}
                      >
                        <Edit className="w-3 h-3" />
                      </Button>
                      <Button
                        variant="danger"
                        size="sm"
                        onClick={() => handleDeleteJob(job.id)}
                      >
                        <Trash2 className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Logs Tab */}
      {activeTab === 'logs' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              View recent sync job execution logs
            </p>
            <Button variant="secondary" onClick={loadData}>
              <RefreshCw className="w-4 h-4 mr-2" />
              Refresh
            </Button>
          </div>

          {logs.length === 0 ? (
            <Card>
              <div className="p-8 text-center">
                <Clock className="w-12 h-12 mx-auto text-gray-400 mb-4" />
                <p className="text-gray-600 dark:text-gray-400">No sync logs available</p>
              </div>
            </Card>
          ) : (
            <div className="space-y-3">
              {logs.map((log) => (
                <Card key={log.id}>
                  <div className="p-6">
                    <div className="flex items-start gap-3">
                      {getStatusIcon(log.status)}
                      <div className="flex-1">
                        <div className="flex items-center justify-between">
                          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                            {log.job_name}
                          </h3>
                          <span className="text-sm text-gray-500 dark:text-gray-400">
                            {new Date(log.started_at).toLocaleString()}
                          </span>
                        </div>
                        <div className="grid grid-cols-4 gap-4 mt-3 text-sm">
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Transferred:</span>
                            <span className="ml-2 text-gray-900 dark:text-gray-100">
                              {formatBytes(log.bytes_transferred)}
                            </span>
                          </div>
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Files:</span>
                            <span className="ml-2 text-gray-900 dark:text-gray-100">
                              {log.files_transferred}
                            </span>
                          </div>
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Duration:</span>
                            <span className="ml-2 text-gray-900 dark:text-gray-100">
                              {formatDuration(log.duration)}
                            </span>
                          </div>
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Triggered:</span>
                            <span className="ml-2 text-gray-900 dark:text-gray-100">
                              {log.triggered_by}
                            </span>
                          </div>
                        </div>
                        {log.error_message && (
                          <div className="mt-3 p-3 bg-red-50 dark:bg-red-900/20 rounded text-sm text-red-600 dark:text-red-400">
                            {log.error_message}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Modals */}
      {showProviderModal && (
        <ProviderModal
          provider={editingProvider}
          onClose={() => {
            setShowProviderModal(false);
            setEditingProvider(null);
          }}
          onSaved={() => {
            setShowProviderModal(false);
            setEditingProvider(null);
            loadData();
          }}
        />
      )}

      {showJobModal && (
        <JobModal
          job={editingJob}
          providers={providers}
          onClose={() => {
            setShowJobModal(false);
            setEditingJob(null);
          }}
          onSaved={() => {
            setShowJobModal(false);
            setEditingJob(null);
            loadData();
          }}
        />
      )}
    </div>
  );
}
