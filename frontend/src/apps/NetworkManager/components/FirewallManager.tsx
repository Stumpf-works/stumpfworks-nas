import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { networkApi, FirewallStatus } from '@/api/network';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';

export default function FirewallManager() {
  const [status, setStatus] = useState<FirewallStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAddRule, setShowAddRule] = useState(false);
  const [newRule, setNewRule] = useState({
    action: 'allow',
    port: '',
    protocol: 'tcp',
    from: '',
    to: '',
  });

  useEffect(() => {
    loadFirewallStatus();
  }, []);

  const loadFirewallStatus = async () => {
    try {
      const response = await networkApi.getFirewallStatus();
      if (response.success && response.data) {
        setStatus(response.data);
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load firewall status');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const toggleFirewall = async () => {
    if (!status) return;

    try {
      const response = await networkApi.setFirewallState(!status.enabled);
      if (response.success) {
        loadFirewallStatus();
      } else {
        alert(response.error?.message || 'Failed to toggle firewall');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleAddRule = async () => {
    try {
      const response = await networkApi.addFirewallRule(
        newRule.action,
        newRule.port,
        newRule.protocol,
        newRule.from,
        newRule.to
      );

      if (response.success) {
        setShowAddRule(false);
        setNewRule({
          action: 'allow',
          port: '',
          protocol: 'tcp',
          from: '',
          to: '',
        });
        loadFirewallStatus();
      } else {
        alert(response.error?.message || 'Failed to add rule');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleDeleteRule = async (ruleNumber: number) => {
    if (!confirm('Are you sure you want to delete this rule?')) return;

    try {
      const response = await networkApi.deleteFirewallRule(ruleNumber);
      if (response.success) {
        loadFirewallStatus();
      } else {
        alert(response.error?.message || 'Failed to delete rule');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleSetDefaultPolicy = async (direction: string, policy: string) => {
    try {
      const response = await networkApi.setDefaultPolicy(direction, policy);
      if (response.success) {
        loadFirewallStatus();
      } else {
        alert(response.error?.message || 'Failed to set default policy');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleResetFirewall = async () => {
    if (!confirm('Are you sure you want to reset the firewall? This will remove all rules!')) return;

    try {
      const response = await networkApi.resetFirewall();
      if (response.success) {
        loadFirewallStatus();
      } else {
        alert(response.error?.message || 'Failed to reset firewall');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const getActionColor = (action: string) => {
    switch (action.toLowerCase()) {
      case 'allow':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'deny':
        return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400';
      case 'reject':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  if (!status) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <div className="text-6xl mb-4">🛡️</div>
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
            Firewall Not Available
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            UFW (Uncomplicated Firewall) is not installed or accessible
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6 max-w-6xl">
      {/* Error Display */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Firewall Status Card */}
      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                <span>🛡️</span>
                Firewall Status
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                UFW (Uncomplicated Firewall) Management
              </p>
            </div>
            <div className="flex items-center gap-4">
              <span
                className={`px-3 py-1 rounded-full text-sm font-medium ${
                  status.enabled
                    ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                    : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400'
                }`}
              >
                {status.enabled ? 'Enabled' : 'Disabled'}
              </span>
              <button
                onClick={toggleFirewall}
                className={`relative inline-flex h-8 w-14 items-center rounded-full transition-colors ${
                  status.enabled ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
                }`}
              >
                <span
                  className={`inline-block h-6 w-6 transform rounded-full bg-white transition-transform ${
                    status.enabled ? 'translate-x-7' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>
          </div>

          {/* Default Policies */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <div className="text-sm text-gray-600 dark:text-gray-400 mb-2">Incoming</div>
              <select
                value={status.defaultIncoming}
                onChange={(e) => handleSetDefaultPolicy('incoming', e.target.value)}
                className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue text-sm"
              >
                <option value="allow">Allow</option>
                <option value="deny">Deny</option>
                <option value="reject">Reject</option>
              </select>
            </div>
            <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <div className="text-sm text-gray-600 dark:text-gray-400 mb-2">Outgoing</div>
              <select
                value={status.defaultOutgoing}
                onChange={(e) => handleSetDefaultPolicy('outgoing', e.target.value)}
                className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue text-sm"
              >
                <option value="allow">Allow</option>
                <option value="deny">Deny</option>
                <option value="reject">Reject</option>
              </select>
            </div>
            <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <div className="text-sm text-gray-600 dark:text-gray-400 mb-2">Routed</div>
              <select
                value={status.defaultRouted}
                onChange={(e) => handleSetDefaultPolicy('routed', e.target.value)}
                className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue text-sm"
              >
                <option value="allow">Allow</option>
                <option value="deny">Deny</option>
                <option value="reject">Reject</option>
              </select>
            </div>
          </div>
        </div>
      </Card>

      {/* Rules Card */}
      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
              Firewall Rules ({status.rules?.length || 0})
            </h2>
            <div className="flex gap-2">
              <Button variant="danger" onClick={handleResetFirewall}>
                Reset All
              </Button>
              <Button onClick={() => setShowAddRule(true)}>+ Add Rule</Button>
            </div>
          </div>

          {status.rules && status.rules.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead className="bg-gray-50 dark:bg-gray-800">
                  <tr>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      #
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Action
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Port
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Protocol
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      From
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      To
                    </th>
                    <th className="px-4 py-3 text-right font-medium text-gray-700 dark:text-gray-300">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {status.rules?.map((rule) => (
                    <tr
                      key={rule.number}
                      className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                    >
                      <td className="px-4 py-3 text-gray-900 dark:text-gray-100">
                        {rule.number}
                      </td>
                      <td className="px-4 py-3">
                        <span
                          className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getActionColor(
                            rule.action
                          )}`}
                        >
                          {rule.action.toUpperCase()}
                        </span>
                      </td>
                      <td className="px-4 py-3 font-mono text-gray-900 dark:text-gray-100">
                        {rule.port || '—'}
                      </td>
                      <td className="px-4 py-3 text-gray-900 dark:text-gray-100">
                        {rule.protocol || '—'}
                      </td>
                      <td className="px-4 py-3 font-mono text-gray-900 dark:text-gray-100">
                        {rule.from || '—'}
                      </td>
                      <td className="px-4 py-3 font-mono text-gray-900 dark:text-gray-100">
                        {rule.to || '—'}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <Button
                          size="sm"
                          variant="danger"
                          onClick={() => handleDeleteRule(rule.number)}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="text-center py-12 text-gray-500 dark:text-gray-400">
              <div className="text-4xl mb-2">📋</div>
              <p>No firewall rules configured</p>
              <p className="text-sm mt-1">Click "Add Rule" to create your first rule</p>
            </div>
          )}
        </div>
      </Card>

      {/* Add Rule Modal */}
      <AnimatePresence>
        {showAddRule && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setShowAddRule(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Add Firewall Rule
              </h2>

              <div className="space-y-4">
                {/* Action */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Action
                  </label>
                  <select
                    value={newRule.action}
                    onChange={(e) => setNewRule({ ...newRule, action: e.target.value })}
                    className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                  >
                    <option value="allow">Allow</option>
                    <option value="deny">Deny</option>
                    <option value="reject">Reject</option>
                  </select>
                </div>

                {/* Port */}
                <Input
                  label="Port"
                  value={newRule.port}
                  onChange={(e) => setNewRule({ ...newRule, port: e.target.value })}
                  placeholder="80, 443, or 8000:8100"
                />

                {/* Protocol */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Protocol
                  </label>
                  <select
                    value={newRule.protocol}
                    onChange={(e) => setNewRule({ ...newRule, protocol: e.target.value })}
                    className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                  >
                    <option value="tcp">TCP</option>
                    <option value="udp">UDP</option>
                    <option value="">Any</option>
                  </select>
                </div>

                {/* From */}
                <Input
                  label="From (Source IP/Network)"
                  value={newRule.from}
                  onChange={(e) => setNewRule({ ...newRule, from: e.target.value })}
                  placeholder="any or 192.168.1.0/24"
                />

                {/* To */}
                <Input
                  label="To (Destination IP/Network)"
                  value={newRule.to}
                  onChange={(e) => setNewRule({ ...newRule, to: e.target.value })}
                  placeholder="any or 192.168.1.100"
                />
              </div>

              {/* Actions */}
              <div className="flex gap-3 mt-6">
                <Button
                  variant="secondary"
                  onClick={() => setShowAddRule(false)}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button onClick={handleAddRule} className="flex-1">
                  Add Rule
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
