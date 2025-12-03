import { useState, useEffect } from 'react';
import Modal from '@/components/ui/Modal';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { timeMachineApi, TimeMachineConfig } from '@/api/timemachine';
import { getErrorMessage } from '@/api/client';
import { AlertCircle, Info } from 'lucide-react';

interface SettingsModalProps {
  config: TimeMachineConfig | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export default function SettingsModal({ config, isOpen, onClose, onSuccess }: SettingsModalProps) {
  const [formData, setFormData] = useState<Partial<TimeMachineConfig>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (config) {
      setFormData({
        share_name: config.share_name,
        base_path: config.base_path,
        default_quota_gb: config.default_quota_gb,
        auto_discovery: config.auto_discovery,
        use_smb: config.use_smb,
        use_afp: config.use_afp,
        smb_version: config.smb_version,
      });
    }
  }, [config]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      setLoading(true);
      const response = await timeMachineApi.updateConfig(formData);

      if (response.success) {
        onSuccess();
        onClose();
      } else {
        setError(response.error?.message || 'Failed to update configuration');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Time Machine Settings" size="large">
      <form onSubmit={handleSubmit}>
        <div className="space-y-6">
          {/* Error Message */}
          {error && (
            <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
              <div className="flex items-start gap-2">
                <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
                <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
              </div>
            </div>
          )}

          {/* Info Banner */}
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
            <div className="flex items-start gap-2">
              <Info className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-blue-800 dark:text-blue-200">
                Changing these settings will reload the Time Machine service and may temporarily interrupt backups.
              </p>
            </div>
          </div>

          {/* Share Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Share Name
            </label>
            <Input
              type="text"
              value={formData.share_name || ''}
              onChange={(e) => setFormData({ ...formData, share_name: e.target.value })}
              placeholder="TimeMachine"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Name shown in Finder sidebar on macOS
            </p>
          </div>

          {/* Base Path */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Base Path
            </label>
            <Input
              type="text"
              value={formData.base_path || ''}
              onChange={(e) => setFormData({ ...formData, base_path: e.target.value })}
              placeholder="/mnt/storage/timemachine"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Root directory for all Time Machine backups
            </p>
          </div>

          {/* Default Quota */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Default Quota (GB)
            </label>
            <Input
              type="number"
              value={formData.default_quota_gb || 0}
              onChange={(e) => setFormData({ ...formData, default_quota_gb: parseInt(e.target.value) || 0 })}
              placeholder="500"
              min={0}
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Default storage quota for new devices (0 = unlimited)
            </p>
          </div>

          {/* Protocol Settings */}
          <div className="space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              Protocol Settings
            </h3>

            {/* Use SMB */}
            <div className="flex items-start gap-3">
              <input
                type="checkbox"
                id="use-smb"
                checked={formData.use_smb || false}
                onChange={(e) => setFormData({ ...formData, use_smb: e.target.checked })}
                className="mt-1 w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
              />
              <div className="flex-1">
                <label htmlFor="use-smb" className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Enable SMB (Recommended)
                </label>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                  Modern protocol supported by macOS 10.9+ (Mountain Lion and later)
                </p>
              </div>
            </div>

            {/* SMB Version */}
            {formData.use_smb && (
              <div className="ml-7">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  SMB Version
                </label>
                <select
                  value={formData.smb_version || '3'}
                  onChange={(e) => setFormData({ ...formData, smb_version: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-macos-blue"
                >
                  <option value="1">SMB 1.0 (Legacy)</option>
                  <option value="2">SMB 2.0</option>
                  <option value="3">SMB 3.0 (Recommended)</option>
                </select>
              </div>
            )}

            {/* Use AFP */}
            <div className="flex items-start gap-3">
              <input
                type="checkbox"
                id="use-afp"
                checked={formData.use_afp || false}
                onChange={(e) => setFormData({ ...formData, use_afp: e.target.checked })}
                className="mt-1 w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
              />
              <div className="flex-1">
                <label htmlFor="use-afp" className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Enable AFP (Legacy)
                </label>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                  Legacy protocol for older macOS versions (deprecated by Apple)
                </p>
              </div>
            </div>
          </div>

          {/* Discovery Settings */}
          <div className="space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              Network Discovery
            </h3>

            {/* Auto Discovery */}
            <div className="flex items-start gap-3">
              <input
                type="checkbox"
                id="auto-discovery"
                checked={formData.auto_discovery || false}
                onChange={(e) => setFormData({ ...formData, auto_discovery: e.target.checked })}
                className="mt-1 w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
              />
              <div className="flex-1">
                <label htmlFor="auto-discovery" className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Enable Bonjour/Avahi Discovery
                </label>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                  Advertise this server on the local network so Macs can discover it automatically
                </p>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
            <Button variant="secondary" onClick={onClose} type="button">
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={loading}>
              {loading ? 'Saving...' : 'Save Settings'}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  );
}
