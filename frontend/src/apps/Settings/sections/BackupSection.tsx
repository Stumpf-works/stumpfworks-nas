// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { backupApi } from '@/api/backup';
import { getErrorMessage } from '@/api/client';

export function BackupSection() {
  const [backups, setBackups] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadBackups();
  }, []);

  const loadBackups = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await backupApi.listJobs();
      if (response.success && response.data) {
        setBackups(response.data);
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Backup Configuration</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure snapshots, rsync, and cloud backup settings
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Backup Jobs
            </h2>
            <Button variant="primary" onClick={loadBackups}>
              Refresh
            </Button>
          </div>

          {loading ? (
            <p className="text-gray-600 dark:text-gray-400">Loading backups...</p>
          ) : backups.length === 0 ? (
            <p className="text-gray-600 dark:text-gray-400">No backup jobs configured</p>
          ) : (
            <div className="space-y-3">
              {backups.map((backup: any, idx: number) => (
                <div
                  key={backup.id || idx}
                  className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                >
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                    {backup.name || `Backup ${idx + 1}`}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {backup.type || 'Unknown type'}
                  </p>
                </div>
              ))}
            </div>
          )}
        </div>
      </Card>

      {/* Note */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          For creating and managing backup jobs, use the dedicated Backup Manager app.
        </p>
      </div>
    </div>
  );
}
