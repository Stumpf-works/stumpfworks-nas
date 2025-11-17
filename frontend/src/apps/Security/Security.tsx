import { useEffect, useState } from 'react';
import { securityApi, FailedLoginAttempt, IPBlock, FailedLoginStats } from '@/api/security';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import { motion } from 'framer-motion';

export function Security() {
  const [attempts, setAttempts] = useState<FailedLoginAttempt[]>([]);
  const [blockedIPs, setBlockedIPs] = useState<IPBlock[]>([]);
  const [stats, setStats] = useState<FailedLoginStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<'attempts' | 'blocked'>('attempts');

  // Pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [totalAttempts, setTotalAttempts] = useState(0);
  const [pageSize] = useState(50);

  const loadAttempts = async () => {
    setIsLoading(true);
    setError('');
    try {
      const offset = (currentPage - 1) * pageSize;
      const response = await securityApi.getFailedLogins(pageSize, offset);
      if (response.success && response.data) {
        setAttempts(response.data.attempts || []);
        setTotalAttempts(response.data.total || 0);
      } else {
        setError(response.error?.message || 'Failed to load failed login attempts');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  const loadBlockedIPs = async () => {
    try {
      const response = await securityApi.getBlockedIPs();
      if (response.success && response.data) {
        setBlockedIPs(response.data);
      }
    } catch (err) {
      console.error('Failed to load blocked IPs:', err);
    }
  };

  const loadStats = async () => {
    try {
      const response = await securityApi.getStats();
      if (response.success && response.data) {
        setStats(response.data);
      }
    } catch (err) {
      console.error('Failed to load stats:', err);
    }
  };

  useEffect(() => {
    loadAttempts();
    loadBlockedIPs();
    loadStats();
  }, [currentPage]);

  const handleUnblockIP = async (ipAddress: string) => {
    if (!confirm(`Are you sure you want to unblock ${ipAddress}?`)) {
      return;
    }

    try {
      const response = await securityApi.unblockIP(ipAddress);
      if (response.success) {
        // Reload data
        loadBlockedIPs();
        loadStats();
        alert('IP address unblocked successfully');
      } else {
        alert('Failed to unblock IP: ' + response.error?.message);
      }
    } catch (err) {
      alert('Failed to unblock IP: ' + getErrorMessage(err));
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const getReasonBadge = (reason: string) => {
    const colors: Record<string, string> = {
      invalid_password: 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-200',
      user_not_found: 'bg-orange-100 dark:bg-orange-900/20 text-orange-800 dark:text-orange-200',
      account_disabled: 'bg-yellow-100 dark:bg-yellow-900/20 text-yellow-800 dark:text-yellow-200',
    };
    return colors[reason] || 'bg-gray-100 dark:bg-gray-900/20 text-gray-800 dark:text-gray-200';
  };

  const totalPages = Math.ceil(totalAttempts / pageSize);

  return (
    <div className="p-4 md:p-6 h-full overflow-auto bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="mb-4 md:mb-6">
        <h1 className="text-xl md:text-2xl font-bold text-gray-900 dark:text-gray-100">
          Security Dashboard
        </h1>
        <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mt-1">
          Monitor failed login attempts and manage blocked IPs
        </p>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4 mb-4 md:mb-6">
          <Card>
            <div className="p-3 md:p-4">
              <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400">Total Failed Attempts</p>
              <p className="text-lg md:text-2xl font-bold text-gray-900 dark:text-gray-100">
                {stats.total_attempts.toLocaleString()}
              </p>
            </div>
          </Card>
          <Card>
            <div className="p-3 md:p-4">
              <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400">Last 24 Hours</p>
              <p className="text-lg md:text-2xl font-bold text-orange-600 dark:text-orange-400">
                {stats.last_24h_attempts.toLocaleString()}
              </p>
            </div>
          </Card>
          <Card>
            <div className="p-3 md:p-4">
              <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400">Blocked IPs</p>
              <p className="text-lg md:text-2xl font-bold text-red-600 dark:text-red-400">
                {stats.blocked_ips.toLocaleString()}
              </p>
            </div>
          </Card>
          <Card>
            <div className="p-3 md:p-4">
              <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400">Top Failed User</p>
              <p className="text-sm md:text-lg font-bold text-gray-900 dark:text-gray-100 truncate">
                {stats.top_failed_usernames?.[0]?.username || 'N/A'}
              </p>
              <p className="text-xs md:text-sm text-gray-500">
                {stats.top_failed_usernames?.[0]?.count || 0} attempts
              </p>
            </div>
          </Card>
        </div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 mb-4">
        <button
          onClick={() => setActiveTab('attempts')}
          className={`px-3 md:px-4 py-2 rounded-lg text-sm md:text-base font-medium transition-colors ${
            activeTab === 'attempts'
              ? 'bg-macos-blue text-white'
              : 'bg-white dark:bg-macos-dark-100 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-macos-dark-200'
          }`}
        >
          <span className="hidden sm:inline">Failed Login Attempts</span>
          <span className="sm:hidden">Attempts</span>
        </button>
        <button
          onClick={() => setActiveTab('blocked')}
          className={`px-3 md:px-4 py-2 rounded-lg text-sm md:text-base font-medium transition-colors ${
            activeTab === 'blocked'
              ? 'bg-macos-blue text-white'
              : 'bg-white dark:bg-macos-dark-100 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-macos-dark-200'
          }`}
        >
          <span className="hidden sm:inline">Blocked IPs ({blockedIPs.length})</span>
          <span className="sm:hidden">Blocked ({blockedIPs.length})</span>
        </button>
      </div>

      {/* Error Message */}
      {error && (
        <div className="mb-4 md:mb-6 p-3 md:p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm md:text-base text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Content */}
      {activeTab === 'attempts' ? (
        <Card>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-macos-dark-200 border-b border-gray-200 dark:border-macos-dark-300">
                <tr>
                  <th className="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">
                    Timestamp
                  </th>
                  <th className="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">
                    Username
                  </th>
                  <th className="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">
                    IP Address
                  </th>
                  <th className="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">
                    Reason
                  </th>
                  <th className="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-macos-dark-100 divide-y divide-gray-200 dark:divide-macos-dark-300">
                {isLoading ? (
                  <tr>
                    <td colSpan={5} className="px-4 md:px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                      Loading...
                    </td>
                  </tr>
                ) : attempts.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="px-4 md:px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                      No failed login attempts
                    </td>
                  </tr>
                ) : (
                  attempts.map((attempt) => (
                    <tr key={attempt.id} className="hover:bg-gray-50 dark:hover:bg-macos-dark-200">
                      <td className="px-4 md:px-6 py-4 whitespace-nowrap text-xs md:text-sm text-gray-900 dark:text-gray-100">
                        {formatDate(attempt.createdAt)}
                      </td>
                      <td className="px-4 md:px-6 py-4 whitespace-nowrap text-xs md:text-sm font-medium text-gray-900 dark:text-gray-100">
                        {attempt.username}
                      </td>
                      <td className="px-4 md:px-6 py-4 whitespace-nowrap text-xs md:text-sm text-gray-900 dark:text-gray-100">
                        {attempt.ipAddress}
                      </td>
                      <td className="px-4 md:px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 py-1 text-xs font-medium rounded-full ${getReasonBadge(
                            attempt.reason
                          )}`}
                        >
                          {attempt.reason.replace(/_/g, ' ')}
                        </span>
                      </td>
                      <td className="px-4 md:px-6 py-4 whitespace-nowrap">
                        {attempt.blocked ? (
                          <span className="px-2 py-1 text-xs font-medium rounded-full bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-200">
                            Blocked
                          </span>
                        ) : (
                          <span className="px-2 py-1 text-xs font-medium rounded-full bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-200">
                            Active
                          </span>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="px-4 md:px-6 py-4 border-t border-gray-200 dark:border-macos-dark-300 flex flex-col sm:flex-row items-center justify-between gap-3">
              <div className="text-xs md:text-sm text-gray-600 dark:text-gray-400">
                Showing {(currentPage - 1) * pageSize + 1} to{' '}
                {Math.min(currentPage * pageSize, totalAttempts)} of {totalAttempts} attempts
              </div>
              <div className="flex gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                <div className="flex items-center gap-2 px-3">
                  <span className="text-xs md:text-sm text-gray-600 dark:text-gray-400">
                    Page {currentPage} of {totalPages}
                  </span>
                </div>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </Card>
      ) : (
        <div className="space-y-4">
          {blockedIPs.length === 0 ? (
            <Card>
              <div className="p-6 md:p-8 text-center">
                <p className="text-sm md:text-base text-gray-500 dark:text-gray-400">No blocked IPs</p>
              </div>
            </Card>
          ) : (
            blockedIPs.map((block) => (
              <motion.div
                key={block.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
              >
                <Card>
                  <div className="p-4 md:p-6">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <h3 className="text-base md:text-lg font-bold text-gray-900 dark:text-gray-100">
                            {block.ipAddress}
                          </h3>
                          {block.isPermanent && (
                            <span className="px-2 py-1 text-xs font-medium rounded-full bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-200">
                              Permanent
                            </span>
                          )}
                        </div>
                        <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mb-1">
                          <strong>Reason:</strong> {block.reason}
                        </p>
                        <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mb-1">
                          <strong>Failed Attempts:</strong> {block.attempts}
                        </p>
                        <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mb-1">
                          <strong>Blocked At:</strong> {formatDate(block.createdAt)}
                        </p>
                        {!block.isPermanent && (
                          <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400">
                            <strong>Expires At:</strong> {formatDate(block.expiresAt)}
                          </p>
                        )}
                      </div>
                      <Button
                        variant="danger"
                        size="sm"
                        onClick={() => handleUnblockIP(block.ipAddress)}
                      >
                        Unblock
                      </Button>
                    </div>
                  </div>
                </Card>
              </motion.div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
