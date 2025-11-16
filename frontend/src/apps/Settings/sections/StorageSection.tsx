// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { syslibApi, type ZFSPool } from '@/api/syslib';
import { getErrorMessage } from '@/api/client';

export function StorageSection() {
  const [pools, setPools] = useState<ZFSPool[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadPools();
  }, []);

  const loadPools = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await syslibApi.zfs.listPools();
      if (response.success && response.data) {
        setPools(response.data);
      } else {
        setError(response.error?.message || 'Failed to load ZFS pools');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleScrubPool = async (poolName: string) => {
    try {
      const response = await syslibApi.zfs.scrubPool(poolName);
      if (response.success) {
        await loadPools();
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Storage Management</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage ZFS pools, disks, and RAID arrays
          </p>
        </div>
        <Button variant="primary" onClick={loadPools}>
          Refresh
        </Button>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* ZFS Pools */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            ZFS Pools
          </h2>

          {loading ? (
            <p className="text-gray-600 dark:text-gray-400">Loading pools...</p>
          ) : pools.length === 0 ? (
            <p className="text-gray-600 dark:text-gray-400">No ZFS pools configured</p>
          ) : (
            <div className="space-y-4">
              {pools.map((pool) => (
                <div
                  key={pool.name}
                  className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                >
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                        {pool.name}
                      </h3>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        Health: <span className={pool.health === 'ONLINE' ? 'text-green-600' : 'text-red-600'}>{pool.health}</span>
                      </p>
                    </div>
                    <Button size="sm" variant="secondary" onClick={() => handleScrubPool(pool.name)}>
                      Scrub
                    </Button>
                  </div>

                  <div className="grid grid-cols-3 gap-4 text-sm">
                    <div>
                      <span className="text-gray-600 dark:text-gray-400">Size:</span>
                      <p className="font-medium text-gray-900 dark:text-gray-100">
                        {(pool.size / 1024 ** 3).toFixed(2)} GB
                      </p>
                    </div>
                    <div>
                      <span className="text-gray-600 dark:text-gray-400">Used:</span>
                      <p className="font-medium text-gray-900 dark:text-gray-100">
                        {(pool.allocated / 1024 ** 3).toFixed(2)} GB
                      </p>
                    </div>
                    <div>
                      <span className="text-gray-600 dark:text-gray-400">Capacity:</span>
                      <p className="font-medium text-gray-900 dark:text-gray-100">
                        {pool.capacity}%
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </Card>

      {/* Note about advanced management */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          For advanced storage management including creating pools, managing disks, and RAID configuration,
          use the dedicated Storage Manager app.
        </p>
      </div>
    </div>
  );
}
