import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { backupApi, BackupJob } from '../../../api/backup';

const BackupJobs: React.FC = () => {
  const [jobs, setJobs] = useState<BackupJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [selectedJob, setSelectedJob] = useState<BackupJob | null>(null);
  const [formData, setFormData] = useState<Partial<BackupJob>>({
    name: '',
    description: '',
    source: '',
    destination: '',
    type: 'full',
    schedule: '0 2 * * *',
    enabled: true,
    retention: 30,
    compression: true,
    encryption: false,
  });

  const fetchJobs = async () => {
    try {
      setLoading(true);
      const response = await backupApi.listJobs();
      setJobs(response.data || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch backup jobs');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchJobs();
    const interval = setInterval(fetchJobs, 10000); // Refresh every 10s
    return () => clearInterval(interval);
  }, []);

  const handleCreate = async () => {
    try {
      await backupApi.createJob(formData);
      setShowCreateModal(false);
      setFormData({
        name: '',
        description: '',
        source: '',
        destination: '',
        type: 'full',
        schedule: '0 2 * * *',
        enabled: true,
        retention: 30,
        compression: true,
        encryption: false,
      });
      fetchJobs();
    } catch (err: any) {
      setError(err.message || 'Failed to create backup job');
    }
  };

  const handleUpdate = async () => {
    if (!selectedJob) return;

    try {
      await backupApi.updateJob(selectedJob.id, formData);
      setShowEditModal(false);
      setSelectedJob(null);
      setFormData({});
      fetchJobs();
    } catch (err: any) {
      setError(err.message || 'Failed to update backup job');
    }
  };

  const handleDelete = async () => {
    if (!selectedJob) return;

    try {
      await backupApi.deleteJob(selectedJob.id);
      setShowDeleteModal(false);
      setSelectedJob(null);
      fetchJobs();
    } catch (err: any) {
      setError(err.message || 'Failed to delete backup job');
    }
  };

  const handleRunJob = async (jobId: string) => {
    try {
      await backupApi.runJob(jobId);
      fetchJobs();
    } catch (err: any) {
      setError(err.message || 'Failed to run backup job');
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'running':
        return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
      case 'success':
        return 'bg-green-500/20 text-green-400 border-green-500/30';
      case 'failed':
        return 'bg-red-500/20 text-red-400 border-red-500/30';
      default:
        return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
    }
  };

  if (loading && jobs.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-gray-400">Loading backup jobs...</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-4">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
            Backup Jobs
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {jobs.length} job(s) configured
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
        >
          Create Job
        </button>
      </div>

      {/* Error Display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-4 bg-red-500/10 border border-red-500/30 rounded-lg"
          >
            <div className="flex justify-between items-start">
              <p className="text-red-400">{error}</p>
              <button
                onClick={() => setError(null)}
                className="text-red-400 hover:text-red-300"
              >
                ✕
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Jobs List */}
      {jobs.length === 0 ? (
        <div className="text-center py-12 text-gray-400">
          No backup jobs configured. Create your first backup job to get started.
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {jobs.map((job) => (
            <motion.div
              key={job.id}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4"
            >
              {/* Job Header */}
              <div className="flex justify-between items-start mb-3">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                    {job.name}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    {job.description}
                  </p>
                </div>
                <span
                  className={`px-2 py-1 text-xs rounded-md border ${getStatusBadge(
                    job.status
                  )}`}
                >
                  {job.status}
                </span>
              </div>

              {/* Job Details */}
              <div className="space-y-2 text-sm mb-4">
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Type:</span>
                  <span className="text-gray-900 dark:text-gray-100">{job.type}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Schedule:</span>
                  <span className="text-gray-900 dark:text-gray-100 font-mono text-xs">
                    {job.schedule}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Source:</span>
                  <span className="text-gray-900 dark:text-gray-100 truncate max-w-[200px]">
                    {job.source}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Retention:</span>
                  <span className="text-gray-900 dark:text-gray-100">
                    {job.retention} days
                  </span>
                </div>
                {job.lastRun && (
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Last Run:</span>
                    <span className="text-gray-900 dark:text-gray-100">
                      {new Date(job.lastRun).toLocaleString()}
                    </span>
                  </div>
                )}
              </div>

              {/* Job Actions */}
              <div className="flex gap-2">
                <button
                  onClick={() => handleRunJob(job.id)}
                  disabled={job.status === 'running'}
                  className="flex-1 px-3 py-2 text-sm bg-green-500/20 text-green-600 dark:text-green-400 hover:bg-green-500/30 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Run Now
                </button>
                <button
                  onClick={() => {
                    setSelectedJob(job);
                    setFormData(job);
                    setShowEditModal(true);
                  }}
                  className="px-3 py-2 text-sm bg-blue-500/20 text-blue-600 dark:text-blue-400 hover:bg-blue-500/30 rounded-lg transition-colors"
                >
                  Edit
                </button>
                <button
                  onClick={() => {
                    setSelectedJob(job);
                    setShowDeleteModal(true);
                  }}
                  className="px-3 py-2 text-sm bg-red-500/20 text-red-600 dark:text-red-400 hover:bg-red-500/30 rounded-lg transition-colors"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          ))}
        </div>
      )}

      {/* Create/Edit Modal */}
      <AnimatePresence>
        {(showCreateModal || showEditModal) && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => {
              setShowCreateModal(false);
              setShowEditModal(false);
            }}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                {showCreateModal ? 'Create Backup Job' : 'Edit Backup Job'}
              </h3>

              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="col-span-2">
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Job Name
                    </label>
                    <input
                      type="text"
                      value={formData.name || ''}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div className="col-span-2">
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Description
                    </label>
                    <textarea
                      value={formData.description || ''}
                      onChange={(e) =>
                        setFormData({ ...formData, description: e.target.value })
                      }
                      rows={2}
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div>
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Source Path
                    </label>
                    <input
                      type="text"
                      value={formData.source || ''}
                      onChange={(e) => setFormData({ ...formData, source: e.target.value })}
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div>
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Destination Path
                    </label>
                    <input
                      type="text"
                      value={formData.destination || ''}
                      onChange={(e) =>
                        setFormData({ ...formData, destination: e.target.value })
                      }
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div>
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Backup Type
                    </label>
                    <select
                      value={formData.type || 'full'}
                      onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                    >
                      <option value="full">Full</option>
                      <option value="incremental">Incremental</option>
                      <option value="differential">Differential</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Schedule (Cron)
                    </label>
                    <input
                      type="text"
                      value={formData.schedule || ''}
                      onChange={(e) => setFormData({ ...formData, schedule: e.target.value })}
                      placeholder="0 2 * * *"
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100 font-mono"
                    />
                  </div>

                  <div>
                    <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Retention (days)
                    </label>
                    <input
                      type="number"
                      value={formData.retention || 30}
                      onChange={(e) =>
                        setFormData({ ...formData, retention: parseInt(e.target.value) })
                      }
                      className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div className="col-span-2 flex gap-4">
                    <label className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={formData.compression || false}
                        onChange={(e) =>
                          setFormData({ ...formData, compression: e.target.checked })
                        }
                        className="rounded"
                      />
                      <span className="text-sm text-gray-600 dark:text-gray-400">
                        Enable Compression
                      </span>
                    </label>

                    <label className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={formData.encryption || false}
                        onChange={(e) =>
                          setFormData({ ...formData, encryption: e.target.checked })
                        }
                        className="rounded"
                      />
                      <span className="text-sm text-gray-600 dark:text-gray-400">
                        Enable Encryption
                      </span>
                    </label>

                    <label className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={formData.enabled || false}
                        onChange={(e) =>
                          setFormData({ ...formData, enabled: e.target.checked })
                        }
                        className="rounded"
                      />
                      <span className="text-sm text-gray-600 dark:text-gray-400">Enabled</span>
                    </label>
                  </div>
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowCreateModal(false);
                    setShowEditModal(false);
                    setFormData({});
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={showCreateModal ? handleCreate : handleUpdate}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
                >
                  {showCreateModal ? 'Create' : 'Update'}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {showDeleteModal && selectedJob && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowDeleteModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Delete Backup Job
              </h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Are you sure you want to delete the backup job "{selectedJob.name}"? This cannot
                be undone.
              </p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => {
                    setShowDeleteModal(false);
                    setSelectedJob(null);
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default BackupJobs;
