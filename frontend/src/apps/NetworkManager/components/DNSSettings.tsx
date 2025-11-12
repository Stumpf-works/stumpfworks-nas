import { useEffect, useState } from 'react';
import { networkApi, DNSConfig, Route } from '@/api/network';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';

export default function DNSSettings() {
  const [dnsConfig, setDnsConfig] = useState<DNSConfig | null>(null);
  const [routes, setRoutes] = useState<Route[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editing, setEditing] = useState(false);

  const [editDNS, setEditDNS] = useState({
    nameservers: [''],
    searchDomains: [''],
  });

  useEffect(() => {
    loadDNSConfig();
    loadRoutes();
  }, []);

  const loadDNSConfig = async () => {
    try {
      const response = await networkApi.getDNS();
      if (response.success && response.data) {
        setDnsConfig(response.data);
        setEditDNS({
          nameservers: response.data.nameservers.length > 0 ? response.data.nameservers : [''],
          searchDomains: response.data.searchDomains.length > 0 ? response.data.searchDomains : [''],
        });
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load DNS config');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const loadRoutes = async () => {
    try {
      const response = await networkApi.getRoutes();
      if (response.success && response.data) {
        setRoutes(response.data);
      }
    } catch (err) {
      console.error('Failed to load routes:', err);
    }
  };

  const handleSaveDNS = async () => {
    try {
      const nameservers = editDNS.nameservers.filter((ns) => ns.trim() !== '');
      const searchDomains = editDNS.searchDomains.filter((sd) => sd.trim() !== '');

      const response = await networkApi.setDNS(nameservers, searchDomains);
      if (response.success) {
        setEditing(false);
        loadDNSConfig();
      } else {
        alert(response.error?.message || 'Failed to save DNS config');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const addNameserver = () => {
    setEditDNS({
      ...editDNS,
      nameservers: [...editDNS.nameservers, ''],
    });
  };

  const removeNameserver = (index: number) => {
    setEditDNS({
      ...editDNS,
      nameservers: editDNS.nameservers.filter((_, i) => i !== index),
    });
  };

  const updateNameserver = (index: number, value: string) => {
    const updated = [...editDNS.nameservers];
    updated[index] = value;
    setEditDNS({ ...editDNS, nameservers: updated });
  };

  const addSearchDomain = () => {
    setEditDNS({
      ...editDNS,
      searchDomains: [...editDNS.searchDomains, ''],
    });
  };

  const removeSearchDomain = (index: number) => {
    setEditDNS({
      ...editDNS,
      searchDomains: editDNS.searchDomains.filter((_, i) => i !== index),
    });
  };

  const updateSearchDomain = (index: number, value: string) => {
    const updated = [...editDNS.searchDomains];
    updated[index] = value;
    setEditDNS({ ...editDNS, searchDomains: updated });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
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

      {/* DNS Configuration */}
      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                DNS Configuration
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Manage DNS nameservers and search domains
              </p>
            </div>
            {!editing ? (
              <Button onClick={() => setEditing(true)}>Edit DNS</Button>
            ) : (
              <div className="flex gap-2">
                <Button variant="secondary" onClick={() => { setEditing(false); loadDNSConfig(); }}>
                  Cancel
                </Button>
                <Button onClick={handleSaveDNS}>Save Changes</Button>
              </div>
            )}
          </div>

          {!editing ? (
            <div className="space-y-4">
              {/* Nameservers Display */}
              <div>
                <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Nameservers
                </h3>
                {dnsConfig && dnsConfig.nameservers.length > 0 ? (
                  <div className="space-y-2">
                    {dnsConfig.nameservers.map((ns, idx) => (
                      <div
                        key={idx}
                        className="px-4 py-2 bg-gray-50 dark:bg-gray-800 rounded-lg font-mono text-sm text-gray-900 dark:text-gray-100"
                      >
                        {ns}
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-gray-500 dark:text-gray-400 text-sm">No nameservers configured</p>
                )}
              </div>

              {/* Search Domains Display */}
              <div>
                <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Search Domains
                </h3>
                {dnsConfig && dnsConfig.searchDomains.length > 0 ? (
                  <div className="space-y-2">
                    {dnsConfig.searchDomains.map((sd, idx) => (
                      <div
                        key={idx}
                        className="px-4 py-2 bg-gray-50 dark:bg-gray-800 rounded-lg font-mono text-sm text-gray-900 dark:text-gray-100"
                      >
                        {sd}
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-gray-500 dark:text-gray-400 text-sm">No search domains configured</p>
                )}
              </div>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Nameservers Edit */}
              <div>
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Nameservers
                  </h3>
                  <Button size="sm" onClick={addNameserver}>
                    + Add Nameserver
                  </Button>
                </div>
                <div className="space-y-2">
                  {editDNS.nameservers.map((ns, idx) => (
                    <div key={idx} className="flex gap-2">
                      <Input
                        value={ns}
                        onChange={(e) => updateNameserver(idx, e.target.value)}
                        placeholder="8.8.8.8"
                        className="flex-1"
                      />
                      {editDNS.nameservers.length > 1 && (
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={() => removeNameserver(idx)}
                        >
                          Remove
                        </Button>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              {/* Search Domains Edit */}
              <div>
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Search Domains
                  </h3>
                  <Button size="sm" onClick={addSearchDomain}>
                    + Add Domain
                  </Button>
                </div>
                <div className="space-y-2">
                  {editDNS.searchDomains.map((sd, idx) => (
                    <div key={idx} className="flex gap-2">
                      <Input
                        value={sd}
                        onChange={(e) => updateSearchDomain(idx, e.target.value)}
                        placeholder="example.com"
                        className="flex-1"
                      />
                      {editDNS.searchDomains.length > 1 && (
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={() => removeSearchDomain(idx)}
                        >
                          Remove
                        </Button>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Routing Table */}
      <Card>
        <div className="p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Routing Table
          </h2>

          {routes.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead className="bg-gray-50 dark:bg-gray-800">
                  <tr>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Destination
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Gateway
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Interface
                    </th>
                    <th className="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">
                      Metric
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {routes.map((route, idx) => (
                    <tr
                      key={idx}
                      className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                    >
                      <td className="px-4 py-3 font-mono text-gray-900 dark:text-gray-100">
                        {route.destination}
                      </td>
                      <td className="px-4 py-3 font-mono text-gray-900 dark:text-gray-100">
                        {route.gateway || '—'}
                      </td>
                      <td className="px-4 py-3 font-mono text-gray-900 dark:text-gray-100">
                        {route.iface}
                      </td>
                      <td className="px-4 py-3 text-gray-900 dark:text-gray-100">
                        {route.metric || '—'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="text-gray-500 dark:text-gray-400 text-center py-8">
              No routes found
            </p>
          )}
        </div>
      </Card>
    </div>
  );
}
