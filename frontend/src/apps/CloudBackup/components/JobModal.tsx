import { useState, useEffect } from 'react';
import { cloudBackupApi, CloudProvider, CloudSyncJob } from '@/api/cloudbackup';
import { getErrorMessage } from '@/api/client';
import Modal from '@/components/ui/Modal';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import { RefreshCw } from 'lucide-react';

interface JobModalProps {
  job: CloudSyncJob | null;
  providers: CloudProvider[];
  onClose: () => void;
  onSaved: () => void;
}

export function JobModal({ job, providers, onClose, onSaved }: JobModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [providerId, setProviderId] = useState<number>(0);
  const [direction, setDirection] = useState<'upload' | 'download' | 'sync'>('upload');
  const [localPath, setLocalPath] = useState('');
  const [remotePath, setRemotePath] = useState('');
  const [schedule, setSchedule] = useState('');
  const [scheduleEnabled, setScheduleEnabled] = useState(false);
  const [bandwidthLimit, setBandwidthLimit] = useState('');
  const [encryptionEnabled, setEncryptionEnabled] = useState(false);
  const [deleteAfterUpload, setDeleteAfterUpload] = useState(false);
  const [retention, setRetention] = useState(0);
  const [enabled, setEnabled] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (job) {
      setName(job.name);
      setDescription(job.description);
      setProviderId(job.provider_id);
      setDirection(job.direction as 'upload' | 'download' | 'sync');
      setLocalPath(job.local_path);
      setRemotePath(job.remote_path);
      setSchedule(job.schedule);
      setScheduleEnabled(job.schedule_enabled);
      setBandwidthLimit(job.bandwidth_limit);
      setEncryptionEnabled(job.encryption_enabled);
      setDeleteAfterUpload(job.delete_after_upload);
      setRetention(job.retention);
      setEnabled(job.enabled);
    } else if (providers.length > 0) {
      setProviderId(providers[0].id);
    }
  }, [job, providers]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSaving(true);

    try {
      const jobData = {
        name,
        description,
        provider_id: providerId,
        direction,
        local_path: localPath,
        remote_path: remotePath,
        schedule,
        schedule_enabled: scheduleEnabled,
        bandwidth_limit: bandwidthLimit,
        encryption_enabled: encryptionEnabled,
        delete_after_upload: deleteAfterUpload,
        retention,
        enabled,
      };

      let response;
      if (job) {
        response = await cloudBackupApi.updateJob(job.id, jobData);
      } else {
        response = await cloudBackupApi.createJob(jobData);
      }

      if (response.success) {
        onSaved();
      } else {
        setError(response.error?.message || 'Failed to save sync job');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title={job ? 'Edit Sync Job' : 'Create Sync Job'}
      icon={<RefreshCw className="w-6 h-6" />}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded text-red-600 dark:text-red-400 text-sm">
            {error}
          </div>
        )}

        <Input
          label="Job Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="Daily Backup"
        />

        <Input
          label="Description (optional)"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Backup important files daily"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Cloud Provider
          </label>
          <select
            value={providerId}
            onChange={(e) => setProviderId(Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
            required
          >
            {providers.map((provider) => (
              <option key={provider.id} value={provider.id}>
                {provider.name} ({provider.type})
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Sync Direction
          </label>
          <select
            value={direction}
            onChange={(e) => setDirection(e.target.value as 'upload' | 'download' | 'sync')}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
          >
            <option value="upload">Upload (Local → Cloud)</option>
            <option value="download">Download (Cloud → Local)</option>
            <option value="sync">Bidirectional Sync</option>
          </select>
        </div>

        <Input
          label="Local Path"
          value={localPath}
          onChange={(e) => setLocalPath(e.target.value)}
          required
          placeholder="/var/backups/data"
        />

        <Input
          label="Remote Path"
          value={remotePath}
          onChange={(e) => setRemotePath(e.target.value)}
          required
          placeholder="backups/data"
        />

        <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
          <h3 className="text-sm font-medium text-gray-900 dark:text-gray-100 mb-3">
            Schedule (optional)
          </h3>
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="schedule-enabled"
                checked={scheduleEnabled}
                onChange={(e) => setScheduleEnabled(e.target.checked)}
                className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
              />
              <label htmlFor="schedule-enabled" className="text-sm text-gray-700 dark:text-gray-300">
                Enable automatic scheduling
              </label>
            </div>
            {scheduleEnabled && (
              <div>
                <Input
                  label="Cron Expression"
                  value={schedule}
                  onChange={(e) => setSchedule(e.target.value)}
                  placeholder="0 2 * * * (daily at 2 AM)"
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Format: minute hour day month weekday
                </p>
              </div>
            )}
          </div>
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
          <h3 className="text-sm font-medium text-gray-900 dark:text-gray-100 mb-3">
            Advanced Options
          </h3>
          <div className="space-y-3">
            <div>
              <Input
                label="Bandwidth Limit (optional)"
                value={bandwidthLimit}
                onChange={(e) => setBandwidthLimit(e.target.value)}
                placeholder="10M, 1G"
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Example: 10M = 10 MB/s, 1G = 1 GB/s
              </p>
            </div>

            <div>
              <Input
                label="Retention (days)"
                type="number"
                value={retention}
                onChange={(e) => setRetention(Number(e.target.value))}
                placeholder="0"
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                0 = keep forever
              </p>
            </div>

            <div className="space-y-2">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="encryption"
                  checked={encryptionEnabled}
                  onChange={(e) => setEncryptionEnabled(e.target.checked)}
                  className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
                />
                <label htmlFor="encryption" className="text-sm text-gray-700 dark:text-gray-300">
                  Enable encryption
                </label>
              </div>

              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="delete-after"
                  checked={deleteAfterUpload}
                  onChange={(e) => setDeleteAfterUpload(e.target.checked)}
                  className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
                />
                <label htmlFor="delete-after" className="text-sm text-gray-700 dark:text-gray-300">
                  Delete local files after upload
                </label>
              </div>

              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="job-enabled"
                  checked={enabled}
                  onChange={(e) => setEnabled(e.target.checked)}
                  className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
                />
                <label htmlFor="job-enabled" className="text-sm text-gray-700 dark:text-gray-300">
                  Enable this job
                </label>
              </div>
            </div>
          </div>
        </div>

        <div className="flex gap-3 pt-4">
          <Button type="submit" disabled={saving}>
            {saving ? 'Saving...' : job ? 'Update Job' : 'Create Job'}
          </Button>
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
        </div>
      </form>
    </Modal>
  );
}
