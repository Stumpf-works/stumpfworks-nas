import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Volume, CreateVolumeRequest, Disk } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { HardDrive, Shield, Package, TreePine, Database, RefreshCw, Plus, Trash2, X, CheckCircle, XCircle, AlertCircle, Folder, Server } from 'lucide-react';

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
    const iconClass = "w-8 h-8";
    switch (type) {
      case 'raid0':
      case 'raid1':
      case 'raid5':
      case 'raid6':
      case 'raid10':
        return <Shield className={`${iconClass} text-blue-500`} />;
      case 'lvm':
        return <Package className={`${iconClass} text-purple-500`} />;
      case 'zfs':
      case 'btrfs':
        return <TreePine className={`${iconClass} text-green-500`} />;
      default:
        return <Database className={`${iconClass} text-gray-500`} />;
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'online':
        return 'bg-gradient-to-r from-green-500 to-emerald-500 text-white';
      case 'degraded':
        return 'bg-gradient-to-r from-yellow-500 to-orange-500 text-white';
      case 'rebuilding':
        return 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white';
      case 'offline':
      case 'failed':
        return 'bg-gradient-to-r from-red-500 to-rose-600 text-white';
      default:
        return 'bg-gradient-to-r from-gray-400 to-gray-500 text-white';
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
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <Server className="w-6 h-6 text-white" />
          </div>
          <div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Storage Volumes
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {volumes.length} volume{volumes.length !== 1 ? 's' : ''} configured
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button
            onClick={loadVolumes}
            variant="secondary"
            className="flex items-center gap-2"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </Button>
          <Button
            onClick={() => setShowCreate(true)}
            className="flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Create Volume
          </Button>
        </div>
      </div>

      {/* Volume Grid */}
      {volumes.length > 0 && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {volumes.map((volume) => (
            <motion.div
              key={volume.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              whileHover={{ y: -2 }}
              transition={{ duration: 0.2 }}
            >
              <Card className="h-full hover:shadow-xl transition-shadow duration-200">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-3 flex-1 min-w-0">
                    <div className="flex-shrink-0 p-2 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl">
                      {getVolumeTypeIcon(volume.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100 truncate">
                        {volume.name}
                      </h3>
                      <p className="text-xs text-gray-600 dark:text-gray-400 uppercase truncate">
                        {volume.type}
                        {volume.raidLevel && ` (${volume.raidLevel})`}
                      </p>
                    </div>
                  </div>
                  <span className={`px-3 py-1 text-xs font-semibold rounded-full shadow-sm flex-shrink-0 ml-2 ${getStatusBadge(volume.status)}`}>
                    {volume.status}
                  </span>
                </div>

                {/* Storage Usage */}
                <div className="mb-4 p-4 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-xl">
                  <div className="flex justify-between text-sm mb-2">
                    <span className="text-gray-600 dark:text-gray-400 font-medium">Storage Usage</span>
                    <span className="font-semibold text-gray-900 dark:text-gray-100">
                      {formatBytes(volume.used)} / {formatBytes(volume.size)}
                    </span>
                  </div>
                  <div className="relative h-3 overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${(volume.used / volume.size) * 100}%` }}
                      transition={{ duration: 1, ease: 'easeOut' }}
                      className="h-full bg-gradient-to-r from-macos-blue to-macos-purple rounded-full shadow-sm"
                    />
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    {((volume.used / volume.size) * 100).toFixed(1)}% used
                  </p>
                </div>

                {/* Details Grid */}
                <div className="grid grid-cols-2 gap-3 mb-4">
                  <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
                    <div className="flex items-center gap-2 mb-1">
                      <Database className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Filesystem</span>
                    </div>
                    <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 uppercase">
                      {volume.filesystem}
                    </p>
                  </div>
                  <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
                    <div className="flex items-center gap-2 mb-1">
                      <CheckCircle className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Health</span>
                    </div>
                    <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                      {volume.health}%
                    </p>
                  </div>
                  <div className="col-span-2 p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
                    <div className="flex items-center gap-2 mb-1">
                      <Folder className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Mount Point</span>
                    </div>
                    <p className="text-sm font-mono font-semibold text-gray-900 dark:text-gray-100 truncate">
                      {volume.mountPoint}
                    </p>
                  </div>
                </div>

                {/* Disks */}
                {volume.disks && volume.disks.length > 0 && (
                  <div className="mb-4">
                    <div className="flex items-center gap-2 mb-2">
                      <HardDrive className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">
                        Disks ({volume.disks.length})
                      </span>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {volume.disks.map((disk) => (
                        <span
                          key={disk}
                          className="px-3 py-1.5 bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 border border-gray-200 dark:border-gray-700 rounded-lg text-xs font-mono font-medium text-gray-900 dark:text-gray-100 shadow-sm"
                        >
                          {disk}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Actions */}
                <div>
                  <Button
                    onClick={() => handleDelete(volume.id)}
                    variant="secondary"
                    size="sm"
                    className="w-full flex items-center justify-center gap-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
                  >
                    <Trash2 className="w-4 h-4" />
                    Delete Volume
                  </Button>
                </div>
              </Card>
            </motion.div>
          ))}
        </div>
      )}

      {/* Empty State */}
      {volumes.length === 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex flex-col items-center justify-center py-16 px-4"
        >
          <div className="p-6 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-2xl mb-6">
            <Server className="w-16 h-16 text-gray-400 dark:text-gray-600" />
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
            No Volumes Found
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6 text-center max-w-md">
            Create your first storage volume to start managing your data
          </p>
          <Button
            onClick={() => setShowCreate(true)}
            className="flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Create Volume
          </Button>
        </motion.div>
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

// Helper function to preview auto-generated mount point
function generateMountPointPreview(name: string): string {
  if (!name) return '/mnt/my-volume';

  // Same logic as backend
  let safeName = name.toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^a-z0-9\-_]/g, '')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '');

  if (!safeName) safeName = 'volume';

  return `/mnt/${safeName}`;
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
      className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-2xl shadow-2xl p-6 w-full max-w-2xl max-h-[80vh] overflow-auto"
      >
        {/* Header */}
        <div className="flex items-center justify-between mb-6 pb-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl">
              <Plus className="w-6 h-6 text-white" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create Volume
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Configure a new storage volume
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

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

          <div>
            <Input
              label="Mount Point (Optional)"
              value={formData.mountPoint}
              onChange={(e) => setFormData({ ...formData, mountPoint: e.target.value })}
              placeholder={formData.name ? generateMountPointPreview(formData.name) : "/mnt/my-volume"}
            />
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {formData.mountPoint
                ? `Will be mounted at: ${formData.mountPoint}`
                : formData.name
                  ? `Will auto-generate: ${generateMountPointPreview(formData.name)}`
                  : "Leave empty to auto-generate from volume name"}
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Select Disks ({formData.disks.length} selected)
            </label>
            <div className="space-y-2 max-h-40 overflow-y-auto p-2 bg-gray-50 dark:bg-macos-dark-200/50 rounded-xl">
              {availableDisks.map((disk) => (
                <label
                  key={disk.name}
                  className={`flex items-center p-3 rounded-lg cursor-pointer transition-all ${
                    formData.disks.includes(disk.name)
                      ? 'bg-gradient-to-r from-macos-blue/10 to-macos-purple/10 border-2 border-macos-blue dark:border-macos-purple'
                      : 'bg-white dark:bg-macos-dark-200 border-2 border-transparent hover:border-gray-300 dark:hover:border-gray-600'
                  }`}
                >
                  <input
                    type="checkbox"
                    checked={formData.disks.includes(disk.name)}
                    onChange={() => toggleDisk(disk.name)}
                    className="mr-3 w-4 h-4 text-macos-blue"
                  />
                  <HardDrive className="w-5 h-5 text-gray-500 dark:text-gray-400 mr-2" />
                  <div className="flex-1">
                    <div className="font-medium text-gray-900 dark:text-gray-100">
                      {disk.name}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-400">
                      {disk.model} • {(disk.size / 1024 / 1024 / 1024).toFixed(2)} GB
                    </div>
                  </div>
                </label>
              ))}
              {availableDisks.length === 0 && (
                <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                  <AlertCircle className="w-8 h-8 mx-auto mb-2 opacity-50" />
                  <p className="text-sm">No available disks found</p>
                </div>
              )}
            </div>
          </div>

          {error && (
            <div className="p-3 bg-gradient-to-br from-red-50 to-rose-50 dark:from-red-900/20 dark:to-rose-900/20 border border-red-200 dark:border-red-800 rounded-xl flex items-start gap-2">
              <XCircle className="w-4 h-4 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-600 dark:text-red-400">
                {error}
              </p>
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
