import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { haApi, ClusterStatus } from '@/api/ha';
import toast from 'react-hot-toast';

export default function ClusterPanel() {
  const [clusterStatus, setClusterStatus] = useState<ClusterStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [maintenanceMode, setMaintenanceMode] = useState(false);

  const loadClusterStatus = async () => {
    try {
      const response = await haApi.getClusterStatus();
      if (response.success && response.data) {
        setClusterStatus(response.data);
        setMaintenanceMode(response.data.maintenance_mode);
      } else {
        toast.error(response.error?.message || 'Failed to load cluster status');
      }
    } catch (error) {
      toast.error('Failed to load cluster status');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadClusterStatus();
    const interval = setInterval(loadClusterStatus, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const toggleMaintenanceMode = async () => {
    try {
      const newMode = !maintenanceMode;
      const response = await haApi.setMaintenanceMode(newMode);
      if (response.success) {
        toast.success(`Maintenance mode ${newMode ? 'enabled' : 'disabled'}`);
        setMaintenanceMode(newMode);
        await loadClusterStatus();
      } else {
        toast.error(response.error?.message || 'Failed to toggle maintenance mode');
      }
    } catch (error) {
      toast.error('Failed to toggle maintenance mode');
      console.error(error);
    }
  };

  const handleNodeAction = async (nodeName: string, action: 'standby' | 'unstandby') => {
    try {
      const actionFn = action === 'standby' ? haApi.standbyNode : haApi.unstandbyNode;
      const response = await actionFn(nodeName);
      if (response.success) {
        toast.success(`Node ${action === 'standby' ? 'put in standby' : 'activated'}`);
        await loadClusterStatus();
      } else {
        toast.error(response.error?.message || `Failed to ${action} node`);
      }
    } catch (error) {
      toast.error(`Failed to ${action} node`);
      console.error(error);
    }
  };

  const handleResourceAction = async (action: string, actionFn: () => Promise<any>) => {
    try {
      const response = await actionFn();
      if (response.success) {
        toast.success(`${action} successful`);
        await loadClusterStatus();
      } else {
        toast.error(response.error?.message || `${action} failed`);
      }
    } catch (error) {
      toast.error(`${action} failed`);
      console.error(error);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue"></div>
      </div>
    );
  }

  if (!clusterStatus) {
    return (
      <div className="p-6 text-center text-gray-500 dark:text-gray-400">
        No cluster configured or Pacemaker/Corosync not available
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Cluster Overview */}
      <div className="bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
              Cluster Status
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              {clusterStatus.name || 'Unnamed Cluster'}
            </p>
          </div>
          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2 cursor-pointer">
              <span className="text-sm text-gray-700 dark:text-gray-300">Maintenance Mode</span>
              <input
                type="checkbox"
                checked={maintenanceMode}
                onChange={toggleMaintenanceMode}
                className="w-4 h-4 text-macos-blue rounded focus:ring-macos-blue"
              />
            </label>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-4">
          <div className="p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
            <p className="text-sm text-gray-600 dark:text-gray-400">Quorum</p>
            <p className={`text-lg font-semibold ${clusterStatus.quorum ? 'text-green-600' : 'text-red-600'}`}>
              {clusterStatus.quorum ? 'Yes' : 'No'}
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
            <p className="text-sm text-gray-600 dark:text-gray-400">Nodes Online</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {clusterStatus.nodes.filter(n => n.online).length} / {clusterStatus.nodes.length}
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
            <p className="text-sm text-gray-600 dark:text-gray-400">Resources</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {clusterStatus.resources.filter(r => r.active).length} / {clusterStatus.resources.length} Active
            </p>
          </div>
        </div>
      </div>

      {/* Nodes */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Cluster Nodes</h3>
        <div className="grid gap-4">
          {clusterStatus.nodes.map((node) => (
            <motion.div
              key={node.name}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 p-4"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="text-lg font-semibold text-gray-900 dark:text-gray-100">{node.name}</h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">{node.ip || 'N/A'}</p>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                    node.online ? 'text-green-600 bg-green-100' : 'text-red-600 bg-red-100'
                  }`}>
                    {node.online ? 'Online' : 'Offline'}
                  </span>
                  {node.online && (
                    <button
                      onClick={() => handleNodeAction(node.name, 'standby')}
                      className="px-3 py-1.5 bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200 transition-colors text-sm"
                    >
                      Standby
                    </button>
                  )}
                  {!node.online && (
                    <button
                      onClick={() => handleNodeAction(node.name, 'unstandby')}
                      className="px-3 py-1.5 bg-green-100 text-green-700 rounded hover:bg-green-200 transition-colors text-sm"
                    >
                      Activate
                    </button>
                  )}
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Resources */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Cluster Resources</h3>
        {clusterStatus.resources.length === 0 ? (
          <div className="text-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-6xl mb-4">🔗</div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
              No Cluster Resources
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Configure resources using pcs command line tool
            </p>
          </div>
        ) : (
          <div className="grid gap-4">
            {clusterStatus.resources.map((resource) => (
              <motion.div
                key={resource.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 p-4"
              >
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h4 className="text-lg font-semibold text-gray-900 dark:text-gray-100">{resource.id}</h4>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{resource.type}:{resource.agent}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">Running on: {resource.node || 'N/A'}</p>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                      resource.active ? 'text-green-600 bg-green-100' : 'text-gray-600 bg-gray-100'
                    }`}>
                      {resource.active ? 'Active' : 'Inactive'}
                    </span>
                    {resource.failed && (
                      <span className="px-3 py-1 rounded-full text-sm font-medium text-red-600 bg-red-100">
                        Failed
                      </span>
                    )}
                  </div>
                </div>
                <div className="flex flex-wrap gap-2">
                  <button
                    onClick={() => handleResourceAction('Enable', () => haApi.enableClusterResource(resource.id))}
                    className="px-3 py-1.5 bg-green-100 text-green-700 rounded hover:bg-green-200 transition-colors text-sm"
                  >
                    Enable
                  </button>
                  <button
                    onClick={() => handleResourceAction('Disable', () => haApi.disableClusterResource(resource.id))}
                    className="px-3 py-1.5 bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200 transition-colors text-sm"
                  >
                    Disable
                  </button>
                  {resource.failed && (
                    <button
                      onClick={() => handleResourceAction('Clear Failed State', () => haApi.clearClusterResource(resource.id))}
                      className="px-3 py-1.5 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors text-sm"
                    >
                      Clear
                    </button>
                  )}
                  <button
                    onClick={() => handleResourceAction('Delete', () => haApi.deleteClusterResource(resource.id))}
                    className="px-3 py-1.5 bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors text-sm"
                  >
                    Delete
                  </button>
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
