import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { adDCApi, type DCStatus, type ProvisionOptions, type ADDCUser, type ADGroup } from '@/api/ad-dc';
import toast from 'react-hot-toast';

type ADTab = 'domain' | 'users' | 'groups';

export function ADDomainController() {
  const [activeTab, setActiveTab] = useState<ADTab>('domain');
  const [status, setStatus] = useState<DCStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [showProvisionDialog, setShowProvisionDialog] = useState(false);

  const tabs = [
    { id: 'domain' as ADTab, name: 'Domain', icon: '🏰' },
    { id: 'users' as ADTab, name: 'Users', icon: '👤' },
    { id: 'groups' as ADTab, name: 'Groups', icon: '👥' },
  ];

  useEffect(() => {
    loadStatus();
  }, []);

  const loadStatus = async () => {
    try {
      setLoading(true);
      const response = await adDCApi.getStatus();
      if (response.success && response.data) {
        setStatus(response.data);
      }
    } catch (error: any) {
      toast.error(`Failed to load AD DC status: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleProvision = () => {
    setShowProvisionDialog(true);
  };

  const handleDemote = async () => {
    if (!confirm('Demote this domain controller? This will remove all AD data!')) {
      return;
    }

    try {
      const response = await adDCApi.demoteDomain();
      if (response.success) {
        toast.success('Domain controller demoted successfully');
        loadStatus();
      }
    } catch (error: any) {
      toast.error(`Failed to demote domain: ${error.message}`);
    }
  };

  const handleRestartService = async () => {
    try {
      const response = await adDCApi.restartService();
      if (response.success) {
        toast.success('Samba AD DC service restarted');
      }
    } catch (error: any) {
      toast.error(`Failed to restart service: ${error.message}`);
    }
  };

  return (
    <div className="h-full flex flex-col bg-gray-50 dark:bg-macos-dark-200">
      {/* Header */}
      <div className="bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Active Directory Domain Controller
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Manage Samba AD Domain Controller
            </p>
          </div>
          {status && (
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <span className={`w-3 h-3 rounded-full ${status.provisioned ? 'bg-green-500' : 'bg-gray-400'}`} />
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  {status.provisioned ? `${status.config.realm}` : 'Not Provisioned'}
                </span>
              </div>
              {status.provisioned ? (
                <>
                  <button
                    onClick={handleRestartService}
                    className="px-4 py-2 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg hover:bg-blue-200 dark:hover:bg-blue-900/50 transition-colors"
                  >
                    Restart Service
                  </button>
                  <button
                    onClick={handleDemote}
                    className="px-4 py-2 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors"
                  >
                    Demote DC
                  </button>
                </>
              ) : (
                <button
                  onClick={handleProvision}
                  className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  Provision Domain
                </button>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700">
        <div className="flex space-x-1 px-6">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-macos-blue text-macos-blue'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
              }`}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.name}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
          </div>
        ) : (
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.2 }}
          >
            {activeTab === 'domain' && <DomainTab status={status} onRefresh={loadStatus} />}
            {activeTab === 'users' && <UsersTab provisioned={status?.provisioned || false} />}
            {activeTab === 'groups' && <GroupsTab provisioned={status?.provisioned || false} />}
          </motion.div>
        )}
      </div>

      {/* Provision Dialog */}
      {showProvisionDialog && (
        <ProvisionDialog
          onClose={() => setShowProvisionDialog(false)}
          onSuccess={() => {
            setShowProvisionDialog(false);
            loadStatus();
          }}
        />
      )}
    </div>
  );
}

// ===== Domain Tab =====

interface DomainTabProps {
  status: DCStatus | null;
  onRefresh: () => void;
}

function DomainTab({ status }: DomainTabProps) {
  if (!status) {
    return <div className="text-center py-12 text-gray-500">No status information available</div>;
  }

  if (!status.provisioned) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">🏰</div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
          Domain Controller Not Provisioned
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-4">
          Click "Provision Domain" to set up an Active Directory domain
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Domain Info */}
      <div className="bg-white dark:bg-macos-dark-100 rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Domain Information</h2>
        <dl className="grid grid-cols-2 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Realm</dt>
            <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">{status.config.realm}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Domain</dt>
            <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">{status.config.domain}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Function Level</dt>
            <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">{status.config.function_level}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">DNS Backend</dt>
            <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">{status.config.dns_backend}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Service Status</dt>
            <dd className="mt-1 text-sm">
              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                status.service_status?.trim() === 'active'
                  ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
                  : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
              }`}>
                {status.service_status || 'unknown'}
              </span>
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Host IP</dt>
            <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">{status.config.host_ip || 'N/A'}</dd>
          </div>
        </dl>
      </div>

      {/* Domain Extended Info */}
      {status.domain_info && Object.keys(status.domain_info).length > 0 && (
        <div className="bg-white dark:bg-macos-dark-100 rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Extended Information</h2>
          <dl className="grid grid-cols-2 gap-4">
            {Object.entries(status.domain_info).map(([key, value]) => (
              <div key={key}>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">{key}</dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                  {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                </dd>
              </div>
            ))}
          </dl>
        </div>
      )}
    </div>
  );
}

// ===== Users Tab =====

interface UsersTabProps {
  provisioned: boolean;
}

function UsersTab({ provisioned }: UsersTabProps) {
  const [users, setUsers] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  useEffect(() => {
    if (provisioned) {
      loadUsers();
    } else {
      setLoading(false);
    }
  }, [provisioned]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const response = await adDCApi.listUsers();
      if (response.success && response.data) {
        setUsers(response.data);
      }
    } catch (error: any) {
      toast.error(`Failed to load users: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (username: string) => {
    if (!confirm(`Delete user ${username}?`)) {
      return;
    }

    try {
      const response = await adDCApi.deleteUser(username);
      if (response.success) {
        toast.success(`User ${username} deleted`);
        loadUsers();
      }
    } catch (error: any) {
      toast.error(`Failed to delete user: ${error.message}`);
    }
  };

  const handleEnable = async (username: string) => {
    try {
      const response = await adDCApi.enableUser(username);
      if (response.success) {
        toast.success(`User ${username} enabled`);
        loadUsers();
      }
    } catch (error: any) {
      toast.error(`Failed to enable user: ${error.message}`);
    }
  };

  const handleDisable = async (username: string) => {
    try {
      const response = await adDCApi.disableUser(username);
      if (response.success) {
        toast.success(`User ${username} disabled`);
        loadUsers();
      }
    } catch (error: any) {
      toast.error(`Failed to disable user: ${error.message}`);
    }
  };

  if (!provisioned) {
    return (
      <div className="text-center py-12 text-gray-500 dark:text-gray-400">
        Domain controller must be provisioned first
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          AD Users ({users.length})
        </h2>
        <button
          onClick={() => setShowCreateDialog(true)}
          className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          + Create User
        </button>
      </div>

      <div className="bg-white dark:bg-macos-dark-100 rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Username
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-macos-dark-100 divide-y divide-gray-200 dark:divide-gray-700">
            {users.map((username) => (
              <tr key={username} className="hover:bg-gray-50 dark:hover:bg-gray-800">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <span className="text-2xl mr-3">👤</span>
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {username}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                  <button
                    onClick={() => handleEnable(username)}
                    className="text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-500 mr-3"
                  >
                    Enable
                  </button>
                  <button
                    onClick={() => handleDisable(username)}
                    className="text-yellow-600 hover:text-yellow-700 dark:text-yellow-400 dark:hover:text-yellow-500 mr-3"
                  >
                    Disable
                  </button>
                  <button
                    onClick={() => handleDelete(username)}
                    className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-500"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showCreateDialog && (
        <CreateUserDialog
          onClose={() => setShowCreateDialog(false)}
          onSuccess={() => {
            setShowCreateDialog(false);
            loadUsers();
          }}
        />
      )}
    </div>
  );
}

// ===== Groups Tab =====

interface GroupsTabProps {
  provisioned: boolean;
}

function GroupsTab({ provisioned }: GroupsTabProps) {
  const [groups, setGroups] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  useEffect(() => {
    if (provisioned) {
      loadGroups();
    } else {
      setLoading(false);
    }
  }, [provisioned]);

  const loadGroups = async () => {
    try {
      setLoading(true);
      const response = await adDCApi.listGroups();
      if (response.success && response.data) {
        setGroups(response.data);
      }
    } catch (error: any) {
      toast.error(`Failed to load groups: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (groupName: string) => {
    if (!confirm(`Delete group ${groupName}?`)) {
      return;
    }

    try {
      const response = await adDCApi.deleteGroup(groupName);
      if (response.success) {
        toast.success(`Group ${groupName} deleted`);
        loadGroups();
      }
    } catch (error: any) {
      toast.error(`Failed to delete group: ${error.message}`);
    }
  };

  if (!provisioned) {
    return (
      <div className="text-center py-12 text-gray-500 dark:text-gray-400">
        Domain controller must be provisioned first
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          AD Groups ({groups.length})
        </h2>
        <button
          onClick={() => setShowCreateDialog(true)}
          className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          + Create Group
        </button>
      </div>

      <div className="bg-white dark:bg-macos-dark-100 rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Group Name
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-macos-dark-100 divide-y divide-gray-200 dark:divide-gray-700">
            {groups.map((groupName) => (
              <tr key={groupName} className="hover:bg-gray-50 dark:hover:bg-gray-800">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <span className="text-2xl mr-3">👥</span>
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {groupName}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                  <button
                    onClick={() => handleDelete(groupName)}
                    className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-500"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showCreateDialog && (
        <CreateGroupDialog
          onClose={() => setShowCreateDialog(false)}
          onSuccess={() => {
            setShowCreateDialog(false);
            loadGroups();
          }}
        />
      )}
    </div>
  );
}

// ===== Provision Dialog =====

interface ProvisionDialogProps {
  onClose: () => void;
  onSuccess: () => void;
}

function ProvisionDialog({ onClose, onSuccess }: ProvisionDialogProps) {
  const [options, setOptions] = useState<ProvisionOptions>({
    realm: '',
    domain: '',
    admin_password: '',
    dns_backend: 'SAMBA_INTERNAL',
    server_role: 'dc',
    function_level: '2008_R2',
  });
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!options.realm || !options.domain || !options.admin_password) {
      toast.error('Realm, Domain, and Admin Password are required');
      return;
    }

    try {
      setSubmitting(true);
      const response = await adDCApi.provisionDomain(options);
      if (response.success) {
        toast.success('Domain provisioned successfully');
        onSuccess();
      }
    } catch (error: any) {
      toast.error(`Failed to provision domain: ${error.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" onClick={onClose}>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl w-full max-w-md"
      >
        <form onSubmit={handleSubmit}>
          <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Provision AD Domain</h2>
          </div>

          <div className="p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Realm (e.g., EXAMPLE.COM) *
              </label>
              <input
                type="text"
                value={options.realm}
                onChange={(e) => setOptions({ ...options, realm: e.target.value.toUpperCase() })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Domain (NetBIOS, e.g., EXAMPLE) *
              </label>
              <input
                type="text"
                value={options.domain}
                onChange={(e) => setOptions({ ...options, domain: e.target.value.toUpperCase() })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Administrator Password *
              </label>
              <input
                type="password"
                value={options.admin_password}
                onChange={(e) => setOptions({ ...options, admin_password: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Function Level
              </label>
              <select
                value={options.function_level}
                onChange={(e) => setOptions({ ...options, function_level: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
              >
                <option value="2008_R2">2008 R2</option>
                <option value="2012">2012</option>
                <option value="2012_R2">2012 R2</option>
                <option value="2016">2016</option>
              </select>
            </div>

            <div className="text-xs text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 p-3 rounded">
              <strong>Warning:</strong> Provisioning will configure this server as an Active Directory Domain Controller. This operation may take several minutes.
            </div>
          </div>

          <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
            >
              {submitting ? 'Provisioning...' : 'Provision Domain'}
            </button>
          </div>
        </form>
      </motion.div>
    </div>
  );
}

// ===== Create User Dialog =====

interface CreateUserDialogProps {
  onClose: () => void;
  onSuccess: () => void;
}

function CreateUserDialog({ onClose, onSuccess }: CreateUserDialogProps) {
  const [user, setUser] = useState<ADDCUser>({ username: '' });
  const [password, setPassword] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!user.username || !password) {
      toast.error('Username and password are required');
      return;
    }

    try {
      setSubmitting(true);
      const response = await adDCApi.createUser({ user, password });
      if (response.success) {
        toast.success(`User ${user.username} created`);
        onSuccess();
      }
    } catch (error: any) {
      toast.error(`Failed to create user: ${error.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" onClick={onClose}>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl w-full max-w-md"
      >
        <form onSubmit={handleSubmit}>
          <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Create AD User</h2>
          </div>

          <div className="p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Username *
              </label>
              <input
                type="text"
                value={user.username}
                onChange={(e) => setUser({ ...user, username: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Password *
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Given Name
              </label>
              <input
                type="text"
                value={user.given_name || ''}
                onChange={(e) => setUser({ ...user, given_name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Surname
              </label>
              <input
                type="text"
                value={user.surname || ''}
                onChange={(e) => setUser({ ...user, surname: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Email
              </label>
              <input
                type="email"
                value={user.email || ''}
                onChange={(e) => setUser({ ...user, email: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
              />
            </div>
          </div>

          <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
            >
              {submitting ? 'Creating...' : 'Create User'}
            </button>
          </div>
        </form>
      </motion.div>
    </div>
  );
}

// ===== Create Group Dialog =====

interface CreateGroupDialogProps {
  onClose: () => void;
  onSuccess: () => void;
}

function CreateGroupDialog({ onClose, onSuccess }: CreateGroupDialogProps) {
  const [group, setGroup] = useState<ADGroup>({ name: '' });
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!group.name) {
      toast.error('Group name is required');
      return;
    }

    try {
      setSubmitting(true);
      const response = await adDCApi.createGroup(group);
      if (response.success) {
        toast.success(`Group ${group.name} created`);
        onSuccess();
      }
    } catch (error: any) {
      toast.error(`Failed to create group: ${error.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" onClick={onClose}>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl w-full max-w-md"
      >
        <form onSubmit={handleSubmit}>
          <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Create AD Group</h2>
          </div>

          <div className="p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Group Name *
              </label>
              <input
                type="text"
                value={group.name}
                onChange={(e) => setGroup({ ...group, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Description
              </label>
              <input
                type="text"
                value={group.description || ''}
                onChange={(e) => setGroup({ ...group, description: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
              />
            </div>
          </div>

          <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
            >
              {submitting ? 'Creating...' : 'Create Group'}
            </button>
          </div>
        </form>
      </motion.div>
    </div>
  );
}
