import { useState, useEffect } from 'react';
import { addcApi, ADDNSRecord } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { Globe, Plus, Trash2, RefreshCw, AlertCircle, FolderOpen } from 'lucide-react';

export default function DNSManagement() {
  const [zones, setZones] = useState<string[]>([]);
  const [selectedZone, setSelectedZone] = useState<string>('');
  const [records, setRecords] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingRecords, setLoadingRecords] = useState(false);
  const [error, setError] = useState('');
  const [showCreateZoneForm, setShowCreateZoneForm] = useState(false);
  const [showCreateRecordForm, setShowCreateRecordForm] = useState(false);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [newZoneName, setNewZoneName] = useState('');

  const [createRecordForm, setCreateRecordForm] = useState<ADDNSRecord>({
    name: '',
    type: 'A',
    value: '',
    ttl: 3600,
  });

  const loadZones = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.listDNSZones();
      if (response.success && response.data) {
        setZones(response.data);
      } else {
        setError(response.error?.message || 'Failed to load DNS zones');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load DNS zones');
    } finally {
      setLoading(false);
    }
  };

  const loadRecords = async (zone: string) => {
    try {
      setLoadingRecords(true);
      setError('');
      const response = await addcApi.listDNSRecords(zone);
      if (response.success && response.data) {
        setRecords(response.data);
        setSelectedZone(zone);
      } else {
        setError(response.error?.message || 'Failed to load DNS records');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load DNS records');
    } finally {
      setLoadingRecords(false);
    }
  };

  useEffect(() => {
    loadZones();
  }, []);

  const handleCreateZone = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!newZoneName.trim()) {
      setError('Zone name is required');
      return;
    }

    try {
      setActionLoading('create-zone');
      setError('');
      const response = await addcApi.createDNSZone(newZoneName);

      if (response.success) {
        alert(`DNS zone ${newZoneName} created successfully!`);
        setShowCreateZoneForm(false);
        setNewZoneName('');
        loadZones();
      } else {
        setError(response.error?.message || 'Failed to create DNS zone');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create DNS zone');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDeleteZone = async (zone: string) => {
    if (!confirm(`Are you sure you want to delete DNS zone "${zone}"?`)) {
      return;
    }

    try {
      setActionLoading(`delete-zone-${zone}`);
      setError('');
      const response = await addcApi.deleteDNSZone(zone);

      if (response.success) {
        alert(`DNS zone ${zone} deleted successfully`);
        if (selectedZone === zone) {
          setSelectedZone('');
          setRecords([]);
        }
        loadZones();
      } else {
        setError(response.error?.message || 'Failed to delete DNS zone');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to delete DNS zone');
    } finally {
      setActionLoading(null);
    }
  };

  const handleCreateRecord = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!createRecordForm.name || !createRecordForm.type || !createRecordForm.value) {
      setError('Name, type, and value are required');
      return;
    }

    try {
      setActionLoading('create-record');
      setError('');
      const response = await addcApi.addDNSRecord(selectedZone, createRecordForm);

      if (response.success) {
        alert(`DNS record created successfully!`);
        setShowCreateRecordForm(false);
        setCreateRecordForm({
          name: '',
          type: 'A',
          value: '',
          ttl: 3600,
        });
        loadRecords(selectedZone);
      } else {
        setError(response.error?.message || 'Failed to create DNS record');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create DNS record');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Globe className="w-6 h-6 text-macos-blue" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">DNS Management</h2>
        </div>
        <div className="flex gap-3">
          <button
            onClick={loadZones}
            disabled={loading}
            className="p-2 text-gray-600 dark:text-gray-400 hover:text-macos-blue dark:hover:text-macos-blue transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
          >
            <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
          </button>
          <button
            onClick={() => setShowCreateZoneForm(true)}
            className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Create Zone
          </button>
        </div>
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

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* DNS Zones */}
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">DNS Zones</h3>

          {loading && zones.length === 0 ? (
            <div className="flex items-center justify-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700">
              <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
            </div>
          ) : zones.length === 0 ? (
            <div className="text-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400">
              No DNS zones found
            </div>
          ) : (
            <div className="space-y-2">
              {zones.map((zone) => (
                <div
                  key={zone}
                  className={`bg-white dark:bg-macos-dark-100 border ${
                    selectedZone === zone ? 'border-macos-blue' : 'border-gray-200 dark:border-gray-700'
                  } rounded-lg p-4 hover:shadow-md transition-shadow`}
                >
                  <div className="flex items-center justify-between">
                    <button
                      onClick={() => loadRecords(zone)}
                      className="flex items-center gap-2 flex-1 text-left"
                    >
                      <Globe className="w-5 h-5 text-macos-blue" />
                      <span className="font-medium text-gray-900 dark:text-gray-100">{zone}</span>
                    </button>
                    <button
                      onClick={() => handleDeleteZone(zone)}
                      disabled={actionLoading === `delete-zone-${zone}`}
                      className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50"
                      title="Delete Zone"
                    >
                      {actionLoading === `delete-zone-${zone}` ? (
                        <RefreshCw className="w-4 h-4 animate-spin" />
                      ) : (
                        <Trash2 className="w-4 h-4" />
                      )}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* DNS Records */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              DNS Records {selectedZone && `- ${selectedZone}`}
            </h3>
            {selectedZone && (
              <button
                onClick={() => setShowCreateRecordForm(true)}
                className="flex items-center gap-2 px-3 py-1.5 text-sm bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                <Plus className="w-3 h-3" />
                Add Record
              </button>
            )}
          </div>

          {!selectedZone ? (
            <div className="text-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400">
              <FolderOpen className="w-12 h-12 mx-auto mb-3 opacity-50" />
              Select a zone to view records
            </div>
          ) : loadingRecords ? (
            <div className="flex items-center justify-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700">
              <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
            </div>
          ) : records.length === 0 ? (
            <div className="text-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400">
              No records in this zone
            </div>
          ) : (
            <div className="space-y-2">
              {records.map((record, index) => (
                <div
                  key={index}
                  className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-3"
                >
                  <span className="font-mono text-sm text-gray-900 dark:text-gray-100">{record}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Create Zone Modal */}
      <AnimatePresence>
        {showCreateZoneForm && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !actionLoading && setShowCreateZoneForm(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-lg w-full"
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Create DNS Zone
                </h2>

                <form onSubmit={handleCreateZone} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Zone Name <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={newZoneName}
                      onChange={(e) => setNewZoneName(e.target.value)}
                      placeholder="example.com"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      required
                    />
                  </div>

                  <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <button
                      type="button"
                      onClick={() => setShowCreateZoneForm(false)}
                      disabled={!!actionLoading}
                      className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors disabled:opacity-50"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={!!actionLoading}
                      className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                    >
                      {actionLoading === 'create-zone' ? (
                        <>
                          <RefreshCw className="w-4 h-4 animate-spin" />
                          Creating...
                        </>
                      ) : (
                        <>
                          <Plus className="w-4 h-4" />
                          Create Zone
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

      {/* Create Record Modal */}
      <AnimatePresence>
        {showCreateRecordForm && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !actionLoading && setShowCreateRecordForm(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-lg w-full"
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Add DNS Record to {selectedZone}
                </h2>

                <form onSubmit={handleCreateRecord} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Record Name <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={createRecordForm.name}
                      onChange={(e) => setCreateRecordForm({ ...createRecordForm, name: e.target.value })}
                      placeholder="www"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Type <span className="text-red-500">*</span>
                    </label>
                    <select
                      value={createRecordForm.type}
                      onChange={(e) => setCreateRecordForm({ ...createRecordForm, type: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                    >
                      <option value="A">A</option>
                      <option value="AAAA">AAAA</option>
                      <option value="CNAME">CNAME</option>
                      <option value="MX">MX</option>
                      <option value="TXT">TXT</option>
                      <option value="SRV">SRV</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Value <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={createRecordForm.value}
                      onChange={(e) => setCreateRecordForm({ ...createRecordForm, value: e.target.value })}
                      placeholder="192.168.1.1"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      TTL (seconds)
                    </label>
                    <input
                      type="number"
                      value={createRecordForm.ttl}
                      onChange={(e) => setCreateRecordForm({ ...createRecordForm, ttl: parseInt(e.target.value) })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <button
                      type="button"
                      onClick={() => setShowCreateRecordForm(false)}
                      disabled={!!actionLoading}
                      className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors disabled:opacity-50"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={!!actionLoading}
                      className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                    >
                      {actionLoading === 'create-record' ? (
                        <>
                          <RefreshCw className="w-4 h-4 animate-spin" />
                          Adding...
                        </>
                      ) : (
                        <>
                          <Plus className="w-4 h-4" />
                          Add Record
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
