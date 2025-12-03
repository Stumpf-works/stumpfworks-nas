import { useState, useEffect } from 'react';
import Modal from '@/components/ui/Modal';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { timeMachineApi, TimeMachineDevice, UpdateDeviceRequest } from '@/api/timemachine';
import { getErrorMessage } from '@/api/client';
import { AlertCircle } from 'lucide-react';

interface EditDeviceModalProps {
  device: TimeMachineDevice | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export default function EditDeviceModal({ device, isOpen, onClose, onSuccess }: EditDeviceModalProps) {
  const [formData, setFormData] = useState<UpdateDeviceRequest>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (device) {
      setFormData({
        device_name: device.device_name,
        mac_address: device.mac_address || '',
        model_id: device.model_id || '',
        quota_gb: device.quota_gb,
        username: device.username || 'stumpfs',
        enabled: device.enabled,
      });
    }
  }, [device]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!device) return;

    setError(null);

    // Validation
    if (!formData.device_name?.trim()) {
      setError('Device name is required');
      return;
    }

    try {
      setLoading(true);
      const response = await timeMachineApi.updateDevice(device.id, formData);

      if (response.success) {
        onSuccess();
        onClose();
      } else {
        setError(response.error?.message || 'Failed to update device');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!device) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Edit ${device.device_name}`}>
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

          {/* Current Usage Info */}
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-gray-600 dark:text-gray-400">Storage Used:</span>
                <span className="ml-2 font-medium text-gray-900 dark:text-white">
                  {device.used_gb.toFixed(2)} GB
                </span>
              </div>
              <div>
                <span className="text-gray-600 dark:text-gray-400">Last Backup:</span>
                <span className="ml-2 font-medium text-gray-900 dark:text-white">
                  {device.last_backup ? new Date(device.last_backup).toLocaleString() : 'Never'}
                </span>
              </div>
            </div>
          </div>

          {/* Device Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Device Name <span className="text-red-500">*</span>
            </label>
            <Input
              type="text"
              value={formData.device_name || ''}
              onChange={(e) => setFormData({ ...formData, device_name: e.target.value })}
              placeholder="e.g., MacBook-Pro"
              required
            />
          </div>

          {/* MAC Address */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              MAC Address
            </label>
            <Input
              type="text"
              value={formData.mac_address || ''}
              onChange={(e) => setFormData({ ...formData, mac_address: e.target.value })}
              placeholder="e.g., 00:1A:2B:3C:4D:5E"
            />
          </div>

          {/* Model ID */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Model ID
            </label>
            <Input
              type="text"
              value={formData.model_id || ''}
              onChange={(e) => setFormData({ ...formData, model_id: e.target.value })}
              placeholder="e.g., MacBookPro18,1"
            />
          </div>

          {/* Quota */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Storage Quota (GB)
            </label>
            <Input
              type="number"
              value={formData.quota_gb || 0}
              onChange={(e) => setFormData({ ...formData, quota_gb: parseInt(e.target.value) || 0 })}
              placeholder="500"
              min={0}
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Current usage: {device.used_gb.toFixed(2)} GB
            </p>
          </div>

          {/* Username */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              SMB Username
            </label>
            <Input
              type="text"
              value={formData.username || ''}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              placeholder="stumpfs"
            />
          </div>

          {/* Password */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              New SMB Password (Optional)
            </label>
            <Input
              type="password"
              value={formData.password || ''}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              placeholder="Leave empty to keep current password"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Only fill if you want to change the password
            </p>
          </div>

          {/* Enabled */}
          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="enabled-edit"
              checked={formData.enabled || false}
              onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
              className="w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
            />
            <label htmlFor="enabled-edit" className="text-sm font-medium text-gray-700 dark:text-gray-300">
              Enable device
            </label>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
            <Button variant="secondary" onClick={onClose} type="button">
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={loading}>
              {loading ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  );
}
