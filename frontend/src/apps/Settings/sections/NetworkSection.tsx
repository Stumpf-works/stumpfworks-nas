// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { networkApi, type NetworkInterface, type DNSConfig, type FirewallStatus } from '@/api/network';
import { getErrorMessage } from '@/api/client';

export function NetworkSection({ user, systemInfo }: { user: any; systemInfo: any }) {
  const [activeTab, setActiveTab] = useState<'interfaces' | 'dns' | 'firewall'>('interfaces');
  const [interfaces, setInterfaces] = useState<NetworkInterface[]>([]);
  const [dnsConfig, setDnsConfig] = useState<DNSConfig | null>(null);
  const [firewallStatus, setFirewallStatus] = useState<FirewallStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (activeTab === 'interfaces') {
      loadInterfaces();
    } else if (activeTab === 'dns') {
      loadDNS();
    } else if (activeTab === 'firewall') {
      loadFirewall();
    }
  }, [activeTab]);

  const loadInterfaces = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await networkApi.listInterfaces();
      if (response.success && response.data) {
        setInterfaces(response.data);
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const loadDNS = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await networkApi.getDNS();
      if (response.success && response.data) {
        setDnsConfig(response.data);
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const loadFirewall = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await networkApi.getFirewallStatus();
      if (response.success && response.data) {
        setFirewallStatus(response.data);
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
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Network Configuration</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Manage network interfaces, DNS, and firewall settings
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-macos-dark-300">
        <button
          onClick={() => setActiveTab('interfaces')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'interfaces'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Interfaces
        </button>
        <button
          onClick={() => setActiveTab('dns')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'dns'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          DNS
        </button>
        <button
          onClick={() => setActiveTab('firewall')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'firewall'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Firewall
        </button>
      </div>

      {/* Interfaces Tab */}
      {activeTab === 'interfaces' && (
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Network Interfaces
            </h2>
            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading interfaces...</p>
            ) : (
              <div className="space-y-3">
                {interfaces.map((iface) => (
                  <div
                    key={iface.name}
                    className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">{iface.name}</h3>
                      <span
                        className={`px-2 py-1 text-xs rounded ${
                          iface.isUp
                            ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-200'
                            : 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
                        }`}
                      >
                        {iface.isUp ? 'UP' : 'DOWN'}
                      </span>
                    </div>
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div>
                        <span className="text-gray-600 dark:text-gray-400">MAC:</span>
                        <p className="font-medium text-gray-900 dark:text-gray-100">{iface.hardwareAddr}</p>
                      </div>
                      <div>
                        <span className="text-gray-600 dark:text-gray-400">IP:</span>
                        <p className="font-medium text-gray-900 dark:text-gray-100">
                          {iface.addresses.length > 0 ? iface.addresses[0] : 'N/A'}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </Card>
      )}

      {/* DNS Tab */}
      {activeTab === 'dns' && (
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              DNS Configuration
            </h2>
            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading DNS configuration...</p>
            ) : dnsConfig ? (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Nameservers
                  </label>
                  <div className="space-y-2">
                    {dnsConfig.nameservers.map((ns, idx) => (
                      <p key={idx} className="text-gray-900 dark:text-gray-100">{ns}</p>
                    ))}
                  </div>
                </div>
                {dnsConfig.searchDomains.length > 0 && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Search Domains
                    </label>
                    <div className="space-y-2">
                      {dnsConfig.searchDomains.map((domain, idx) => (
                        <p key={idx} className="text-gray-900 dark:text-gray-100">{domain}</p>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <p className="text-gray-600 dark:text-gray-400">No DNS configuration available</p>
            )}
          </div>
        </Card>
      )}

      {/* Firewall Tab */}
      {activeTab === 'firewall' && (
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Firewall Status
              </h2>
              {firewallStatus && (
                <span
                  className={`px-3 py-1 text-sm rounded ${
                    firewallStatus.enabled
                      ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-200'
                      : 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-200'
                  }`}
                >
                  {firewallStatus.enabled ? 'Enabled' : 'Disabled'}
                </span>
              )}
            </div>
            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading firewall status...</p>
            ) : firewallStatus ? (
              <div className="space-y-4">
                <div className="grid grid-cols-3 gap-4">
                  <div>
                    <span className="text-sm text-gray-600 dark:text-gray-400">Default Incoming:</span>
                    <p className="font-medium text-gray-900 dark:text-gray-100">{firewallStatus.defaultIncoming}</p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600 dark:text-gray-400">Default Outgoing:</span>
                    <p className="font-medium text-gray-900 dark:text-gray-100">{firewallStatus.defaultOutgoing}</p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600 dark:text-gray-400">Rules:</span>
                    <p className="font-medium text-gray-900 dark:text-gray-100">{firewallStatus.rules.length}</p>
                  </div>
                </div>
              </div>
            ) : (
              <p className="text-gray-600 dark:text-gray-400">No firewall status available</p>
            )}
          </div>
        </Card>
      )}

      {/* Note */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          For advanced network configuration, use the dedicated Network Manager app.
        </p>
      </div>
    </div>
  );
}
