import { useState } from 'react';
import Modal from '@/components/ui/Modal';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { timeMachineApi, CreateDeviceRequest } from '@/api/timemachine';
import { getErrorMessage } from '@/api/client';
import { AlertCircle } from 'lucide-react';

interface AddDeviceModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  defaultQuota: number;
}

export default function AddDeviceModal({ isOpen, onClose, onSuccess, defaultQuota }: AddDeviceModalProps) {
  const [formData, setFormData] = useState<CreateDeviceRequest>({
    device_name: '',
    mac_address: '',
    model_id: '',
    quota_gb: defaultQuota,
    username: 'stumpfs',
    password: '',
    enabled: true,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Validation
    if (!formData.device_name?.trim()) {
      setError('Device name is required');
      return;
    }

    try {
      setLoading(true);
      const response = await timeMachineApi.createDevice(formData);

      if (response.success) {
        onSuccess();
        onClose();
        // Reset form
        setFormData({
          device_name: '',
          mac_address: '',
          model_id: '',
          quota_gb: defaultQuota,
          username: 'stumpfs',
          password: '',
          enabled: true,
        });
      } else {
        setError(response.error?.message || 'Failed to create device');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Add Time Machine Device">
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

          {/* Device Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Device Name <span className="text-red-500">*</span>
            </label>
            <Input
              type="text"
              value={formData.device_name}
              onChange={(e) => setFormData({ ...formData, device_name: e.target.value })}
              placeholder="e.g., MacBook-Pro"
              required
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Mac hostname (will be shown in Finder sidebar)
            </p>
          </div>

          {/* MAC Address */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              MAC Address (Optional)
            </label>
            <Input
              type="text"
              value={formData.mac_address}
              onChange={(e) => setFormData({ ...formData, mac_address: e.target.value })}
              placeholder="e.g., 00:1A:2B:3C:4D:5E"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Network interface MAC address for device identification
            </p>
          </div>

          {/* Model ID */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Model ID (Optional)
            </label>
            <Input
              type="text"
              value={formData.model_id}
              onChange={(e) => setFormData({ ...formData, model_id: e.target.value })}
              placeholder="e.g., MacBookPro18,1"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Mac model identifier (from System Information)
            </p>
          </div>

          {/* Quota */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Storage Quota (GB)
            </label>
            <Input
              type="number"
              value={formData.quota_gb}
              onChange={(e) => setFormData({ ...formData, quota_gb: parseInt(e.target.value) || 0 })}
              placeholder="500"
              min={0}
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Maximum storage for this device (0 = unlimited)
            </p>
          </div>

          {/* Username */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              SMB Username
            </label>
            <Input
              type="text"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              placeholder="stumpfs"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              User account for SMB access
            </p>
          </div>

          {/* Password */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              SMB Password (Optional)
            </label>
            <Input
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              placeholder="Leave empty to use existing password"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Password for SMB authentication (leave empty if user already exists)
            </p>
          </div>

          {/* Enabled */}
          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="enabled"
              checked={formData.enabled}
              onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
              className="w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
            />
            <label htmlFor="enabled" className="text-sm font-medium text-gray-700 dark:text-gray-300">
              Enable device immediately
            </label>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
            <Button variant="secondary" onClick={onClose} type="button">
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={loading}>
              {loading ? 'Adding...' : 'Add Device'}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  );
}
