import { useEffect, useState } from 'react';
import { dockerApi, DockerNetwork } from '@/api/docker';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';

export default function NetworkManager() {
  const [networks, setNetworks] = useState<DockerNetwork[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadNetworks();
    const interval = setInterval(loadNetworks, 5000); // Refresh every 5s
    return () => clearInterval(interval);
  }, []);

  const loadNetworks = async () => {
    try {
      const response = await dockerApi.listNetworks();
      if (response.success && response.data) {
        setNetworks(response.data);
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load networks');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const getDriverColor = (driver: string) => {
    switch (driver.toLowerCase()) {
      case 'bridge':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
      case 'host':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'overlay':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      case 'macvlan':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400';
      case 'none':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
      default:
        return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400';
    }
  };

  const getScopeColor = (scope: string) => {
    switch (scope.toLowerCase()) {
      case 'local':
        return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400';
      case 'global':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400';
      case 'swarm':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
    }
  };

  const getContainerCount = (network: DockerNetwork) => {
    return network.containers ? Object.keys(network.containers).length : 0;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Error Display */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-600 dark:text-gray-400">
          {networks.length} network{networks.length !== 1 ? 's' : ''}
        </div>
      </div>

      {/* Networks Grid */}
      {networks.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🌐</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No networks found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            Docker networks will appear here when created
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
          {networks.map((network) => (
            <Card key={network.id} hoverable>
              <div className="p-6">
                {/* Header */}
                <div className="mb-4">
                  <h3 className="font-bold text-lg text-gray-900 dark:text-gray-100 mb-2">
                    {network.name}
                  </h3>
                  <p className="text-xs text-gray-600 dark:text-gray-400 font-mono mb-2">
                    {network.id.substring(0, 12)}
                  </p>
                  <div className="flex gap-2">
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getDriverColor(
                        network.driver
                      )}`}
                    >
                      {network.driver}
                    </span>
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getScopeColor(
                        network.scope
                      )}`}
                    >
                      {network.scope}
                    </span>
                  </div>
                </div>

                {/* Details */}
                <div className="space-y-2 mb-4">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Containers:</span>
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {getContainerCount(network)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Internal:</span>
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {network.internal ? 'Yes' : 'No'}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">IPv6:</span>
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {network.enableIPv6 ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                </div>

                {/* IPAM Configuration */}
                {network.ipam?.config && network.ipam.config.length > 0 && (
                  <div className="mb-4">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                      IPAM Config:
                    </div>
                    <div className="space-y-1">
                      {network.ipam.config.map((config, idx) => (
                        <div
                          key={idx}
                          className="text-xs p-2 bg-gray-50 dark:bg-gray-800 rounded"
                        >
                          {config.subnet && (
                            <div>
                              <span className="font-medium text-gray-900 dark:text-gray-100">
                                Subnet:
                              </span>{' '}
                              <span className="text-gray-600 dark:text-gray-400 font-mono">
                                {config.subnet}
                              </span>
                            </div>
                          )}
                          {config.gateway && (
                            <div>
                              <span className="font-medium text-gray-900 dark:text-gray-100">
                                Gateway:
                              </span>{' '}
                              <span className="text-gray-600 dark:text-gray-400 font-mono">
                                {config.gateway}
                              </span>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Options */}
                {network.options && Object.keys(network.options).length > 0 && (
                  <div className="mb-4">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                      Options:
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {Object.entries(network.options).map(([key, value]) => (
                        <span
                          key={key}
                          className="px-2 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded"
                          title={`${key}=${value}`}
                        >
                          {key}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Labels */}
                {network.labels && Object.keys(network.labels).length > 0 && (
                  <div>
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                      Labels:
                    </div>
                    <div className="space-y-1">
                      {Object.entries(network.labels).map(([key, value]) => (
                        <div
                          key={key}
                          className="text-xs p-2 bg-gray-50 dark:bg-gray-800 rounded"
                        >
                          <span className="font-medium text-gray-900 dark:text-gray-100">
                            {key}:
                          </span>{' '}
                          <span className="text-gray-600 dark:text-gray-400">{value}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
