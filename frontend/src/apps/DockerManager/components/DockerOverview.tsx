import { useEffect, useState } from 'react';
import { dockerApi } from '@/api/docker';
import Card from '@/components/ui/Card';
import { Container, Disc, HardDrive, Network, Activity, AlertTriangle, CheckCircle, XCircle } from 'lucide-react';

interface DockerStats {
  containers: {
    total: number;
    running: number;
    stopped: number;
    paused: number;
  };
  images: {
    total: number;
    size: number;
  };
  volumes: {
    total: number;
  };
  networks: {
    total: number;
  };
}

export default function DockerOverview() {
  const [stats, setStats] = useState<DockerStats>({
    containers: { total: 0, running: 0, stopped: 0, paused: 0 },
    images: { total: 0, size: 0 },
    volumes: { total: 0 },
    networks: { total: 0 },
  });
  const [loading, setLoading] = useState(true);

  const loadStats = async () => {
    try {
      const [containers, images, volumes, networks] = await Promise.all([
        dockerApi.listContainers(true),
        dockerApi.listImages(),
        dockerApi.listVolumes(),
        dockerApi.listNetworks(),
      ]);

      if (containers.success && containers.data) {
        const data = containers.data;
        const running = data.filter((c) => c.state === 'running').length;
        const stopped = data.filter((c) => c.state === 'exited').length;
        const paused = data.filter((c) => c.state === 'paused').length;

        setStats((prev) => ({
          ...prev,
          containers: {
            total: data.length,
            running,
            stopped,
            paused,
          },
        }));
      }

      if (images.success && images.data) {
        const data = images.data;
        const totalSize = data.reduce((acc, img) => acc + (img.size || 0), 0);
        setStats((prev) => ({
          ...prev,
          images: {
            total: data.length,
            size: totalSize,
          },
        }));
      }

      if (volumes.success && volumes.data) {
        const data = volumes.data;
        setStats((prev) => ({
          ...prev,
          volumes: { total: data.length },
        }));
      }

      if (networks.success && networks.data) {
        const data = networks.data;
        setStats((prev) => ({
          ...prev,
          networks: { total: data.length },
        }));
      }
    } catch (err) {
      console.error('Failed to load Docker stats:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadStats();
    const interval = setInterval(loadStats, 5000);
    return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
  };

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 md:gap-6 animate-pulse">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <div className="p-6 h-32 bg-gray-100 dark:bg-gray-800 rounded-lg"></div>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 md:gap-6">
        {/* Containers */}
        <Card className="hover:shadow-lg transition-shadow">
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-3 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                  <Container className="w-6 h-6 text-blue-600 dark:text-blue-400" />
                </div>
                <h3 className="font-semibold text-gray-900 dark:text-gray-100">Containers</h3>
              </div>
              <span className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                {stats.containers.total}
              </span>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  <span className="text-gray-600 dark:text-gray-400">Running</span>
                </div>
                <span className="font-medium text-green-600 dark:text-green-400">
                  {stats.containers.running}
                </span>
              </div>
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  <XCircle className="w-4 h-4 text-gray-500" />
                  <span className="text-gray-600 dark:text-gray-400">Stopped</span>
                </div>
                <span className="font-medium text-gray-600 dark:text-gray-400">
                  {stats.containers.stopped}
                </span>
              </div>
              {stats.containers.paused > 0 && (
                <div className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="w-4 h-4 text-yellow-500" />
                    <span className="text-gray-600 dark:text-gray-400">Paused</span>
                  </div>
                  <span className="font-medium text-yellow-600 dark:text-yellow-400">
                    {stats.containers.paused}
                  </span>
                </div>
              )}
            </div>
          </div>
        </Card>

        {/* Images */}
        <Card className="hover:shadow-lg transition-shadow">
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-3 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
                  <Disc className="w-6 h-6 text-purple-600 dark:text-purple-400" />
                </div>
                <h3 className="font-semibold text-gray-900 dark:text-gray-100">Images</h3>
              </div>
              <span className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                {stats.images.total}
              </span>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  <Activity className="w-4 h-4 text-purple-500" />
                  <span className="text-gray-600 dark:text-gray-400">Total Size</span>
                </div>
                <span className="font-medium text-purple-600 dark:text-purple-400">
                  {formatBytes(stats.images.size)}
                </span>
              </div>
            </div>
          </div>
        </Card>

        {/* Volumes */}
        <Card className="hover:shadow-lg transition-shadow">
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-3 bg-green-100 dark:bg-green-900/30 rounded-lg">
                  <HardDrive className="w-6 h-6 text-green-600 dark:text-green-400" />
                </div>
                <h3 className="font-semibold text-gray-900 dark:text-gray-100">Volumes</h3>
              </div>
              <span className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                {stats.volumes.total}
              </span>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Persistent data storage for containers
              </p>
            </div>
          </div>
        </Card>

        {/* Networks */}
        <Card className="hover:shadow-lg transition-shadow">
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-3 bg-orange-100 dark:bg-orange-900/30 rounded-lg">
                  <Network className="w-6 h-6 text-orange-600 dark:text-orange-400" />
                </div>
                <h3 className="font-semibold text-gray-900 dark:text-gray-100">Networks</h3>
              </div>
              <span className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                {stats.networks.total}
              </span>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Virtual networks for container communication
              </p>
            </div>
          </div>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <div className="p-6">
          <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-4">Quick Actions</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <button className="px-4 py-3 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors text-sm font-medium">
              Pull Image
            </button>
            <button className="px-4 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-sm font-medium">
              Create Container
            </button>
            <button className="px-4 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors text-sm font-medium">
              Create Volume
            </button>
            <button className="px-4 py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 transition-colors text-sm font-medium">
              Create Network
            </button>
          </div>
        </div>
      </Card>
    </div>
  );
}
