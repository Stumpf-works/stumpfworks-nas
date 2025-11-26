import { useState, useEffect } from 'react';
import { addcApi, DCStatus, ProvisionOptions } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { Server, AlertCircle, CheckCircle, RefreshCw, Power, Settings, Database } from 'lucide-react';

export default function DomainStatus() {
  const [status, setStatus] = useState<DCStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showProvisionForm, setShowProvisionForm] = useState(false);
  const [provisioning, setProvisioning] = useState(false);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  // Provision form state
  const [provisionForm, setProvisionForm] = useState<ProvisionOptions>({
    realm: '',
    domain: '',
    admin_password: '',
    dns_backend: 'SAMBA_INTERNAL',
    dns_forwarder: '',
    server_role: 'dc',
    function_level: '2008_R2',
    host_ip: '',
  });

  const loadStatus = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.getStatus();
      if (response.success && response.data) {
        setStatus(response.data);
      } else {
        setError(response.error?.message || 'Failed to load domain status');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load domain status');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadStatus();
  }, []);

  const handleProvision = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!provisionForm.realm || !provisionForm.domain || !provisionForm.admin_password) {
      setError('Realm, Domain, and Administrator Password are required');
      return;
    }

    try {
      setProvisioning(true);
      setError('');
      const response = await addcApi.provisionDomain(provisionForm);

      if (response.success) {
        alert(`Domain ${response.data?.realm} provisioned successfully!`);
        setShowProvisionForm(false);
        loadStatus();
      } else {
        setError(response.error?.message || 'Failed to provision domain');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to provision domain');
    } finally {
      setProvisioning(false);
    }
  };

  const handleDemote = async () => {
    if (!confirm('Are you sure you want to demote this domain controller? This action cannot be undone and will remove all AD data!')) {
      return;
    }

    try {
      setActionLoading('demote');
      setError('');
      const response = await addcApi.demoteDomain();

      if (response.success) {
        alert('Domain controller demoted successfully');
        loadStatus();
      } else {
        setError(response.error?.message || 'Failed to demote domain');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to demote domain');
    } finally {
      setActionLoading(null);
    }
  };

  const handleRestartService = async () => {
    try {
      setActionLoading('restart');
      setError('');
      const response = await addcApi.restartService();

      if (response.success) {
        alert('Samba AD DC service restarted successfully');
        loadStatus();
      } else {
        setError(response.error?.message || 'Failed to restart service');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to restart service');
    } finally {
      setActionLoading(null);
    }
  };

  const handleTestConfig = async () => {
    try {
      setActionLoading('test');
      setError('');
      const response = await addcApi.testConfiguration();

      if (response.success) {
        alert('Configuration is valid!');
      } else {
        setError(response.error?.message || 'Configuration test failed');
      }
    } catch (err: any) {
      setError(err.message || 'Configuration test failed');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDBCheck = async () => {
    try {
      setActionLoading('dbcheck');
      setError('');
      const response = await addcApi.showDBCheck();

      if (response.success && response.data) {
        alert(`Database Check Result:\n${response.data.result}`);
      } else {
        setError(response.error?.message || 'Database check failed');
      }
    } catch (err: any) {
      setError(err.message || 'Database check failed');
    } finally {
      setActionLoading(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="flex items-center gap-3 text-gray-600 dark:text-gray-400">
          <RefreshCw className="w-5 h-5 animate-spin" />
          <span>Loading domain status...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
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
            <div>
              <p className="font-medium text-red-900 dark:text-red-100">Error</p>
              <p className="text-red-700 dark:text-red-300 text-sm mt-1">{error}</p>
            </div>
            <button
              onClick={() => setError('')}
              className="ml-auto text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300"
            >
              ×
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Domain Status Card */}
      {status && (
        <div className="bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <Server className="w-6 h-6 text-macos-blue" />
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
                  Domain Controller Status
                </h2>
              </div>
              <button
                onClick={loadStatus}
                disabled={loading}
                className="p-2 text-gray-600 dark:text-gray-400 hover:text-macos-blue dark:hover:text-macos-blue transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
              >
                <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
              </button>
            </div>

            {/* Status Badge */}
            <div className="flex items-center gap-2 mb-6">
              {status.provisioned ? (
                <>
                  <CheckCircle className="w-5 h-5 text-green-500" />
                  <span className="text-green-700 dark:text-green-400 font-medium">Provisioned</span>
                  {status.service_status && (
                    <span className={`ml-2 px-2 py-1 rounded text-xs font-medium ${
                      status.service_status === 'active'
                        ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400'
                        : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-400'
                    }`}>
                      {status.service_status}
                    </span>
                  )}
                </>
              ) : (
                <>
                  <AlertCircle className="w-5 h-5 text-orange-500" />
                  <span className="text-orange-700 dark:text-orange-400 font-medium">Not Provisioned</span>
                </>
              )}
            </div>

            {/* Domain Information */}
            {status.provisioned && status.domain_info && (
              <div className="bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-4 space-y-2">
                <h3 className="font-medium text-gray-900 dark:text-gray-100 mb-3">Domain Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
                  {Object.entries(status.domain_info).map(([key, value]) => (
                    <div key={key} className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">{key}:</span>
                      <span className="text-gray-900 dark:text-gray-100 font-medium">{value as string}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Configuration */}
            {status.config && (
              <div className="mt-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg p-4">
                <h3 className="font-medium text-gray-900 dark:text-gray-100 mb-3">Configuration</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Realm:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-medium">
                      {status.config.realm || 'N/A'}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Domain:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-medium">
                      {status.config.domain || 'N/A'}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">DNS Backend:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-medium">
                      {status.config.dns_backend}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Function Level:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-medium">
                      {status.config.function_level}
                    </span>
                  </div>
                </div>
              </div>
            )}

            {/* Action Buttons */}
            <div className="mt-6 flex flex-wrap gap-3">
              {status.provisioned ? (
                <>
                  <button
                    onClick={handleRestartService}
                    disabled={actionLoading === 'restart'}
                    className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                  >
                    {actionLoading === 'restart' ? (
                      <RefreshCw className="w-4 h-4 animate-spin" />
                    ) : (
                      <Power className="w-4 h-4" />
                    )}
                    Restart Service
                  </button>
                  <button
                    onClick={handleTestConfig}
                    disabled={actionLoading === 'test'}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50"
                  >
                    {actionLoading === 'test' ? (
                      <RefreshCw className="w-4 h-4 animate-spin" />
                    ) : (
                      <Settings className="w-4 h-4" />
                    )}
                    Test Config
                  </button>
                  <button
                    onClick={handleDBCheck}
                    disabled={actionLoading === 'dbcheck'}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50"
                  >
                    {actionLoading === 'dbcheck' ? (
                      <RefreshCw className="w-4 h-4 animate-spin" />
                    ) : (
                      <Database className="w-4 h-4" />
                    )}
                    Check Database
                  </button>
                  <button
                    onClick={handleDemote}
                    disabled={actionLoading === 'demote'}
                    className="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50"
                  >
                    {actionLoading === 'demote' ? (
                      <RefreshCw className="w-4 h-4 animate-spin" />
                    ) : (
                      <AlertCircle className="w-4 h-4" />
                    )}
                    Demote DC
                  </button>
                </>
              ) : (
                <button
                  onClick={() => setShowProvisionForm(true)}
                  className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  <Server className="w-4 h-4" />
                  Provision Domain
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Provision Form Modal */}
      <AnimatePresence>
        {showProvisionForm && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !provisioning && setShowProvisionForm(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-auto"
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Provision Active Directory Domain
                </h2>

                <form onSubmit={handleProvision} className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* Realm */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Realm <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={provisionForm.realm}
                        onChange={(e) => setProvisionForm({ ...provisionForm, realm: e.target.value.toUpperCase() })}
                        placeholder="EXAMPLE.COM"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    {/* Domain (NetBIOS) */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Domain (NetBIOS) <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={provisionForm.domain}
                        onChange={(e) => setProvisionForm({ ...provisionForm, domain: e.target.value.toUpperCase() })}
                        placeholder="EXAMPLE"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    {/* Admin Password */}
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Administrator Password <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="password"
                        value={provisionForm.admin_password}
                        onChange={(e) => setProvisionForm({ ...provisionForm, admin_password: e.target.value })}
                        placeholder="Strong password required"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    {/* DNS Backend */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        DNS Backend
                      </label>
                      <select
                        value={provisionForm.dns_backend}
                        onChange={(e) => setProvisionForm({ ...provisionForm, dns_backend: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      >
                        <option value="SAMBA_INTERNAL">Samba Internal</option>
                        <option value="BIND9_DLZ">BIND9 DLZ</option>
                        <option value="NONE">None</option>
                      </select>
                    </div>

                    {/* Function Level */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Function Level
                      </label>
                      <select
                        value={provisionForm.function_level}
                        onChange={(e) => setProvisionForm({ ...provisionForm, function_level: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      >
                        <option value="2008_R2">2008 R2</option>
                        <option value="2012">2012</option>
                        <option value="2012_R2">2012 R2</option>
                        <option value="2016">2016</option>
                      </select>
                    </div>

                    {/* DNS Forwarder */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        DNS Forwarder
                      </label>
                      <input
                        type="text"
                        value={provisionForm.dns_forwarder}
                        onChange={(e) => setProvisionForm({ ...provisionForm, dns_forwarder: e.target.value })}
                        placeholder="8.8.8.8"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    {/* Host IP */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Host IP
                      </label>
                      <input
                        type="text"
                        value={provisionForm.host_ip}
                        onChange={(e) => setProvisionForm({ ...provisionForm, host_ip: e.target.value })}
                        placeholder="192.168.1.100"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <button
                      type="button"
                      onClick={() => setShowProvisionForm(false)}
                      disabled={provisioning}
                      className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors disabled:opacity-50"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={provisioning}
                      className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                    >
                      {provisioning ? (
                        <>
                          <RefreshCw className="w-4 h-4 animate-spin" />
                          Provisioning...
                        </>
                      ) : (
                        <>
                          <Server className="w-4 h-4" />
                          Provision Domain
                        </>
                      )}
                    </button>
                  </div>
                </form>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
