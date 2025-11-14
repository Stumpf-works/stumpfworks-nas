import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Volume, CreateVolumeRequest, Disk } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export default function VolumeManager() {
  const [volumes, setVolumes] = useState<Volume[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadVolumes();
  }, []);

  const loadVolumes = async () => {
    try {
      const response = await storageApi.listVolumes();
      if (response.success) {
        setVolumes(response.data || []);
      } else {
        console.error('Failed to load volumes:', response.error);
      }
    } catch (error) {
      console.error('Failed to load volumes:', error);
      alert('Failed to load volumes. Please check the console for details.');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this volume? All data will be lost!')) {
      return;
    }

    try {
      await storageApi.deleteVolume(id);
      loadVolumes();
    } catch (error) {
      console.error('Failed to delete volume:', error);
      alert('Failed to delete volume');
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const getVolumeTypeIcon = (type: string) => {
    switch (type) {
      case 'raid0':
      case 'raid1':
      case 'raid5':
      case 'raid6':
      case 'raid10':
        return '🛡️';
      case 'lvm':
        return '📦';
      case 'zfs':
      case 'btrfs':
        return '🌲';
      default:
        return '💾';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'degraded': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400';
      case 'rebuilding': return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
      case 'offline':
      case 'failed': return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400';
      default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Volumes ({volumes.length})
        </h2>
        <div className="flex space-x-2">
          <Button onClick={loadVolumes} variant="secondary">
            🔄 Refresh
          </Button>
          <Button onClick={() => setShowCreate(true)}>
            ➕ Create Volume
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {volumes.map((volume) => (
          <Card key={volume.id}>
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center space-x-3">
                <div className="text-3xl">{getVolumeTypeIcon(volume.type)}</div>
                <div>
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                    {volume.name}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 uppercase">
                    {volume.type}
                    {volume.raidLevel && ` (${volume.raidLevel})`}
                  </p>
                </div>
              </div>
              <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(volume.status)}`}>
                {volume.status}
              </span>
            </div>

            {/* Storage Usage */}
            <div className="mb-4">
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600 dark:text-gray-400">Storage</span>
                <span className="font-medium text-gray-900 dark:text-gray-100">
                  {formatBytes(volume.used)} / {formatBytes(volume.size)}
                </span>
              </div>
              <div className="overflow-hidden h-2 text-xs flex rounded-full bg-gray-200 dark:bg-gray-700">
                <div
                  style={{ width: `${(volume.used / volume.size) * 100}%` }}
                  className="flex flex-col text-center whitespace-nowrap text-white justify-center bg-macos-blue"
                />
              </div>
            </div>

            {/* Details Grid */}
            <div className="grid grid-cols-2 gap-3 text-sm mb-4">
              <div>
                <span className="text-gray-600 dark:text-gray-400">Filesystem:</span>
                <span className="ml-2 font-medium text-gray-900 dark:text-gray-100 uppercase">
                  {volume.filesystem}
                </span>
              </div>
              <div>
                <span className="text-gray-600 dark:text-gray-400">Health:</span>
                <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                  {volume.health}%
                </span>
              </div>
              <div className="col-span-2">
                <span className="text-gray-600 dark:text-gray-400">Mount:</span>
                <span className="ml-2 font-mono text-sm text-gray-900 dark:text-gray-100">
                  {volume.mountPoint}
                </span>
              </div>
            </div>

            {/* Disks */}
            {volume.disks && volume.disks.length > 0 && (
              <div className="mb-4">
                <div className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase mb-2">
                  Disks ({volume.disks.length})
                </div>
                <div className="flex flex-wrap gap-2">
                  {volume.disks.map((disk) => (
                    <span
                      key={disk}
                      className="px-2 py-1 bg-gray-100 dark:bg-macos-dark-200 rounded text-xs font-mono text-gray-900 dark:text-gray-100"
                    >
                      {disk}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {/* Actions */}
            <div className="flex space-x-2">
              <Button
                onClick={() => handleDelete(volume.id)}
                variant="secondary"
                size="sm"
                className="flex-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
              >
                🗑️ Delete
              </Button>
            </div>
          </Card>
        ))}
      </div>

      {volumes.length === 0 && (
        <div className="text-center py-12 text-gray-600 dark:text-gray-400">
          <div className="text-6xl mb-4">📦</div>
          <p className="text-lg font-medium mb-2">No volumes found</p>
          <p className="text-sm mb-4">Create your first volume to get started</p>
          <Button onClick={() => setShowCreate(true)}>
            ➕ Create Volume
          </Button>
        </div>
      )}

      {/* Create Volume Modal */}
      <AnimatePresence>
        {showCreate && (
          <CreateVolumeModal
            onClose={() => setShowCreate(false)}
            onSuccess={() => {
              setShowCreate(false);
              loadVolumes();
            }}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

interface CreateVolumeModalProps {
  onClose: () => void;
  onSuccess: () => void;
}

function CreateVolumeModal({ onClose, onSuccess }: CreateVolumeModalProps) {
  const [availableDisks, setAvailableDisks] = useState<Disk[]>([]);
  const [formData, setFormData] = useState<CreateVolumeRequest>({
    name: '',
    type: 'single',
    disks: [],
    filesystem: 'ext4',
    mountPoint: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadAvailableDisks();
  }, []);

  const loadAvailableDisks = async () => {
    try {
      const response = await storageApi.listDisks();
      if (response.success && response.data) {
        // Filter out system disks
        setAvailableDisks(response.data.filter((d) => !d.isSystem));
      }
    } catch (error) {
      console.error('Failed to load disks:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await storageApi.createVolume(formData);
      if (response.success) {
        onSuccess();
      } else {
        setError(response.error?.message || 'Failed to create volume');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create volume');
    } finally {
      setLoading(false);
    }
  };

  const toggleDisk = (diskName: string) => {
    setFormData((prev) => ({
      ...prev,
      disks: prev.disks.includes(diskName)
        ? prev.disks.filter((d) => d !== diskName)
        : [...prev.disks, diskName],
    }));
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl max-h-[80vh] overflow-auto"
      >
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">
          Create Volume
        </h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Volume Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="my-volume"
            required
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Volume Type
            </label>
            <select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value as any })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100"
            >
              <option value="single">Single Disk</option>
              <option value="raid0">RAID 0 (Striping)</option>
              <option value="raid1">RAID 1 (Mirroring)</option>
              <option value="raid5">RAID 5</option>
              <option value="raid6">RAID 6</option>
              <option value="raid10">RAID 10</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Filesystem
            </label>
            <select
              value={formData.filesystem}
              onChange={(e) => setFormData({ ...formData, filesystem: e.target.value as any })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100"
            >
              <option value="ext4">ext4</option>
              <option value="xfs">XFS</option>
              <option value="btrfs">Btrfs</option>
            </select>
          </div>

          <Input
            label="Mount Point"
            value={formData.mountPoint}
            onChange={(e) => setFormData({ ...formData, mountPoint: e.target.value })}
            placeholder="/mnt/volume1"
            required
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Select Disks
            </label>
            <div className="space-y-2 max-h-40 overflow-y-auto">
              {availableDisks.map((disk) => (
                <label
                  key={disk.name}
                  className="flex items-center p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg cursor-pointer hover:bg-gray-100 dark:hover:bg-macos-dark-300"
                >
                  <input
                    type="checkbox"
                    checked={formData.disks.includes(disk.name)}
                    onChange={() => toggleDisk(disk.name)}
                    className="mr-3"
                  />
                  <div className="flex-1">
                    <div className="font-medium text-gray-900 dark:text-gray-100">
                      {disk.name} ({disk.model})
                    </div>
                    <div className="text-sm text-gray-600 dark:text-gray-400">
                      {(disk.size / 1024 / 1024 / 1024).toFixed(2)} GB
                    </div>
                  </div>
                </label>
              ))}
            </div>
          </div>

          {error && (
            <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
              {error}
            </div>
          )}

          <div className="flex space-x-3 pt-4">
            <Button type="button" onClick={onClose} variant="secondary" className="flex-1">
              Cancel
            </Button>
            <Button type="submit" isLoading={loading} className="flex-1">
              Create Volume
            </Button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
}
