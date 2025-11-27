import { useState, useEffect } from 'react';
import { addcApi, FSMORoles as FSMORolesType } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { Settings, RefreshCw, AlertCircle, Server } from 'lucide-react';

export default function FSMORoles() {
  const [roles, setRoles] = useState<FSMORolesType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadRoles = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.showFSMORoles();
      if (response.success && response.data) {
        setRoles(response.data);
      } else {
        setError(response.error?.message || 'Failed to load FSMO roles');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load FSMO roles');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadRoles();
  }, []);

  const roleDescriptions: Record<string, string> = {
    'SchemaMasterRole': 'Controls modifications to the Active Directory schema',
    'DomainNamingMasterRole': 'Controls the addition or removal of domains in the forest',
    'PDCEmulatorRole': 'Handles password changes, time synchronization, and Group Policy',
    'RIDMasterRole': 'Allocates blocks of Relative IDs (RIDs) to domain controllers',
    'InfrastructureMasterRole': 'Updates cross-domain group-to-user references',
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Settings className="w-6 h-6 text-macos-blue" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">FSMO Roles</h2>
        </div>
        <button
          onClick={loadRoles}
          disabled={loading}
          className="p-2 text-gray-600 dark:text-gray-400 hover:text-macos-blue dark:hover:text-macos-blue transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
        >
          <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
        </button>
      </div>

      {/* Error Message */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3"
          >
            <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <p className="font-medium text-red-900 dark:text-red-100">Error</p>
              <p className="text-red-700 dark:text-red-300 text-sm mt-1">{error}</p>
            </div>
            <button onClick={() => setError('')} className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300">
              ×
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Info Banner */}
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <Server className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
          <div>
            <p className="font-medium text-blue-900 dark:text-blue-100">About FSMO Roles</p>
            <p className="text-blue-700 dark:text-blue-300 text-sm mt-1">
              Flexible Single Master Operations (FSMO) roles are special tasks performed by specific domain controllers in Active Directory.
              There are five FSMO roles that ensure proper functioning of the directory service.
            </p>
          </div>
        </div>
      </div>

      {/* FSMO Roles List */}
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
        </div>
      ) : !roles || Object.keys(roles).length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No FSMO roles found
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {Object.entries(roles).map(([roleName, holder]) => (
            <div
              key={roleName}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-5 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start gap-4">
                <div className="p-3 bg-macos-blue/10 rounded-lg">
                  <Settings className="w-6 h-6 text-macos-blue" />
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-1">
                    {roleName.replace(/([A-Z])/g, ' $1').trim()}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                    {roleDescriptions[roleName] || 'FSMO role'}
                  </p>
                  <div className="flex items-center gap-2 text-sm">
                    <Server className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                    <span className="text-gray-700 dark:text-gray-300">
                      Holder: <span className="font-medium text-macos-blue">{holder as string}</span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Help Section */}
      <div className="bg-gray-50 dark:bg-macos-dark-50 border border-gray-200 dark:border-gray-700 rounded-lg p-5">
        <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">FSMO Role Management</h3>
        <div className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
          <p>
            <strong className="text-gray-900 dark:text-gray-100">Transfer:</strong> Use this when the current role holder is online and functioning properly. The transfer is graceful and coordinated between domain controllers.
          </p>
          <p>
            <strong className="text-gray-900 dark:text-gray-100">Seize:</strong> Use this only when the current role holder is permanently offline or cannot be contacted. Seizing a role is a forceful operation and should be used as a last resort.
          </p>
          <p className="text-orange-600 dark:text-orange-400 mt-3">
            <strong>Warning:</strong> Improperly managing FSMO roles can cause serious issues in your Active Directory environment. Always ensure you understand the implications before transferring or seizing roles.
          </p>
        </div>
      </div>
    </div>
  );
}
