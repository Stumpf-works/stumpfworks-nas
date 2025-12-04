import { useState, useEffect } from 'react';
import { Server, AlertCircle, CheckCircle, Loader2 } from 'lucide-react';
import { dockerApi } from '@/api/docker';

interface HubStatusData {
  hub_url: string;
  is_online: boolean;
  template_count: number;
  error: string | null;
}

export default function HubStatus() {
  const [status, setStatus] = useState<HubStatusData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchHubStatus();
    // Refresh status every 30 seconds
    const interval = setInterval(fetchHubStatus, 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchHubStatus = async () => {
    try {
      const response = await dockerApi.getHubStatus();
      if (response.success && response.data) {
        setStatus(response.data as HubStatusData);
      }
    } catch (err) {
      console.error('Failed to fetch hub status:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center gap-2 px-3 py-1.5 bg-gray-100 dark:bg-gray-800 rounded-lg text-sm">
        <Loader2 className="w-4 h-4 animate-spin text-gray-500" />
        <span className="text-gray-600 dark:text-gray-400">Checking Hub...</span>
      </div>
    );
  }

  if (!status) return null;

  return (
    <div className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm transition-colors ${
      status.is_online
        ? 'bg-green-50 dark:bg-green-900/20'
        : 'bg-yellow-50 dark:bg-yellow-900/20'
    }`}>
      <Server className={`w-4 h-4 ${
        status.is_online
          ? 'text-green-600 dark:text-green-400'
          : 'text-yellow-600 dark:text-yellow-400'
      }`} />

      {status.is_online ? (
        <>
          <CheckCircle className="w-3.5 h-3.5 text-green-600 dark:text-green-400" />
          <span className="text-green-700 dark:text-green-300 font-medium">
            Hub Online
          </span>
          <span className="text-green-600 dark:text-green-400">
            • {status.template_count} templates
          </span>
        </>
      ) : (
        <>
          <AlertCircle className="w-3.5 h-3.5 text-yellow-600 dark:text-yellow-400" />
          <span className="text-yellow-700 dark:text-yellow-300 font-medium">
            Hub Offline
          </span>
          <span className="text-yellow-600 dark:text-yellow-400 text-xs">
            Using cached templates
          </span>
        </>
      )}

      {/* Tooltip with Hub URL */}
      <div className="hidden group-hover:block absolute top-full left-0 mt-2 px-2 py-1 bg-gray-900 text-white text-xs rounded shadow-lg whitespace-nowrap z-10">
        {status.hub_url}
        {status.error && (
          <div className="text-red-300 mt-1">{status.error}</div>
        )}
      </div>
    </div>
  );
}
