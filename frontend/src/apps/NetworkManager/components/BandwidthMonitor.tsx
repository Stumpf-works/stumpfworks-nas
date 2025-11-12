import { useEffect, useState } from 'react';
import { networkApi, InterfaceStats } from '@/api/network';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';

interface HistoricalData {
  timestamp: string;
  rxBytes: number;
  txBytes: number;
  rxRate: number;
  txRate: number;
}

export default function BandwidthMonitor() {
  const [stats, setStats] = useState<InterfaceStats[]>([]);
  const [history, setHistory] = useState<Record<string, HistoricalData[]>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedInterface, setSelectedInterface] = useState<string>('');

  useEffect(() => {
    loadStats();
    const interval = setInterval(loadStats, 2000); // Update every 2 seconds
    return () => clearInterval(interval);
  }, []);

  const loadStats = async () => {
    try {
      const response = await networkApi.getInterfaceStats();
      if (response.success && response.data) {
        const newStats = response.data;
        setStats(newStats);

        // Select first non-loopback interface by default
        if (!selectedInterface && newStats.length > 0) {
          const firstNonLoopback = newStats.find((s: InterfaceStats) => !s.name.startsWith('lo')) || newStats[0];
          setSelectedInterface(firstNonLoopback.name);
        }

        // Update history
        const timestamp = new Date().toLocaleTimeString();
        setHistory((prev) => {
          const updated = { ...prev };

          newStats.forEach((stat: InterfaceStats) => {
            if (!updated[stat.name]) {
              updated[stat.name] = [];
            }

            const prevData = updated[stat.name];
            let rxRate = 0;
            let txRate = 0;

            if (prevData.length > 0) {
              const lastEntry = prevData[prevData.length - 1];
              const timeDiff = 2; // 2 seconds
              rxRate = (stat.rxBytes - lastEntry.rxBytes) / timeDiff;
              txRate = (stat.txBytes - lastEntry.txBytes) / timeDiff;
            }

            updated[stat.name] = [
              ...prevData.slice(-29), // Keep last 29 entries (total 30 with new one)
              {
                timestamp,
                rxBytes: stat.rxBytes,
                txBytes: stat.txBytes,
                rxRate,
                txRate,
              },
            ];
          });

          return updated;
        });

        setError('');
      } else {
        setError(response.error?.message || 'Failed to load interface stats');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
  };

  const formatRate = (bytesPerSec: number) => {
    if (bytesPerSec === 0) return '0 B/s';
    const k = 1024;
    const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s'];
    const i = Math.floor(Math.log(bytesPerSec) / Math.log(k));
    return `${(bytesPerSec / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
  };

  const selectedStat = stats.find((s) => s.name === selectedInterface);
  const selectedHistory = history[selectedInterface] || [];
  const latestHistory = selectedHistory[selectedHistory.length - 1];

  const maxRate = Math.max(
    ...selectedHistory.map((h) => Math.max(h.rxRate, h.txRate)),
    1
  );

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

      {/* Interface Selection */}
      <Card>
        <div className="p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Select Interface
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
            {stats.map((stat) => (
              <button
                key={stat.name}
                onClick={() => setSelectedInterface(stat.name)}
                className={`p-4 rounded-lg border-2 transition-all text-left ${
                  selectedInterface === stat.name
                    ? 'border-macos-blue bg-blue-50 dark:bg-blue-900/20'
                    : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 bg-white dark:bg-macos-dark-100'
                }`}
              >
                <div className="font-semibold text-gray-900 dark:text-gray-100">
                  {stat.name}
                </div>
                <div className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                  {formatBytes(stat.rxBytes + stat.txBytes)} total
                </div>
              </button>
            ))}
          </div>
        </div>
      </Card>

      {selectedStat && (
        <>
          {/* Real-time Stats */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <div className="p-6">
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Download Speed
                </div>
                <div className="text-2xl font-bold text-green-600 dark:text-green-400">
                  {formatRate(latestHistory?.rxRate || 0)}
                </div>
              </div>
            </Card>
            <Card>
              <div className="p-6">
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Upload Speed
                </div>
                <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">
                  {formatRate(latestHistory?.txRate || 0)}
                </div>
              </div>
            </Card>
            <Card>
              <div className="p-6">
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Total Downloaded
                </div>
                <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                  {formatBytes(selectedStat.rxBytes)}
                </div>
              </div>
            </Card>
            <Card>
              <div className="p-6">
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Total Uploaded
                </div>
                <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                  {formatBytes(selectedStat.txBytes)}
                </div>
              </div>
            </Card>
          </div>

          {/* Bandwidth Graph */}
          <Card>
            <div className="p-6">
              <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
                Bandwidth Usage (Last 60 seconds)
              </h3>

              {selectedHistory.length > 0 ? (
                <div className="space-y-4">
                  {/* Download Graph */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                        ⬇️ Download
                      </span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">
                        {formatRate(latestHistory?.rxRate || 0)}
                      </span>
                    </div>
                    <div className="relative h-24 bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden">
                      <svg className="w-full h-full">
                        <polyline
                          fill="none"
                          stroke="rgb(34, 197, 94)"
                          strokeWidth="2"
                          points={selectedHistory
                            .map((h, i) => {
                              const x = (i / (selectedHistory.length - 1 || 1)) * 100;
                              const y = 100 - (h.rxRate / maxRate) * 90;
                              return `${x}%,${y}%`;
                            })
                            .join(' ')}
                        />
                        <polyline
                          fill="rgba(34, 197, 94, 0.1)"
                          stroke="none"
                          points={[
                            '0%,100%',
                            ...selectedHistory.map((h, i) => {
                              const x = (i / (selectedHistory.length - 1 || 1)) * 100;
                              const y = 100 - (h.rxRate / maxRate) * 90;
                              return `${x}%,${y}%`;
                            }),
                            '100%,100%',
                          ].join(' ')}
                        />
                      </svg>
                    </div>
                  </div>

                  {/* Upload Graph */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                        ⬆️ Upload
                      </span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">
                        {formatRate(latestHistory?.txRate || 0)}
                      </span>
                    </div>
                    <div className="relative h-24 bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden">
                      <svg className="w-full h-full">
                        <polyline
                          fill="none"
                          stroke="rgb(59, 130, 246)"
                          strokeWidth="2"
                          points={selectedHistory
                            .map((h, i) => {
                              const x = (i / (selectedHistory.length - 1 || 1)) * 100;
                              const y = 100 - (h.txRate / maxRate) * 90;
                              return `${x}%,${y}%`;
                            })
                            .join(' ')}
                        />
                        <polyline
                          fill="rgba(59, 130, 246, 0.1)"
                          stroke="none"
                          points={[
                            '0%,100%',
                            ...selectedHistory.map((h, i) => {
                              const x = (i / (selectedHistory.length - 1 || 1)) * 100;
                              const y = 100 - (h.txRate / maxRate) * 90;
                              return `${x}%,${y}%`;
                            }),
                            '100%,100%',
                          ].join(' ')}
                        />
                      </svg>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-12 text-gray-500 dark:text-gray-400">
                  Collecting data...
                </div>
              )}
            </div>
          </Card>

          {/* Packet Statistics */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <div className="p-6">
                <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
                  Packet Statistics
                </h3>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">RX Packets:</span>
                    <span className="font-mono text-gray-900 dark:text-gray-100">
                      {selectedStat.rxPackets.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">TX Packets:</span>
                    <span className="font-mono text-gray-900 dark:text-gray-100">
                      {selectedStat.txPackets.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">RX Errors:</span>
                    <span
                      className={`font-mono ${
                        selectedStat.rxErrors > 0
                          ? 'text-red-600 dark:text-red-400'
                          : 'text-gray-900 dark:text-gray-100'
                      }`}
                    >
                      {selectedStat.rxErrors.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">TX Errors:</span>
                    <span
                      className={`font-mono ${
                        selectedStat.txErrors > 0
                          ? 'text-red-600 dark:text-red-400'
                          : 'text-gray-900 dark:text-gray-100'
                      }`}
                    >
                      {selectedStat.txErrors.toLocaleString()}
                    </span>
                  </div>
                </div>
              </div>
            </Card>

            <Card>
              <div className="p-6">
                <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
                  Drop Statistics
                </h3>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">RX Dropped:</span>
                    <span
                      className={`font-mono ${
                        selectedStat.rxDropped > 0
                          ? 'text-orange-600 dark:text-orange-400'
                          : 'text-gray-900 dark:text-gray-100'
                      }`}
                    >
                      {selectedStat.rxDropped.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">TX Dropped:</span>
                    <span
                      className={`font-mono ${
                        selectedStat.txDropped > 0
                          ? 'text-orange-600 dark:text-orange-400'
                          : 'text-gray-900 dark:text-gray-100'
                      }`}
                    >
                      {selectedStat.txDropped.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Total Errors:</span>
                    <span
                      className={`font-mono ${
                        selectedStat.rxErrors + selectedStat.txErrors > 0
                          ? 'text-red-600 dark:text-red-400'
                          : 'text-green-600 dark:text-green-400'
                      }`}
                    >
                      {(selectedStat.rxErrors + selectedStat.txErrors).toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Total Dropped:</span>
                    <span
                      className={`font-mono ${
                        selectedStat.rxDropped + selectedStat.txDropped > 0
                          ? 'text-orange-600 dark:text-orange-400'
                          : 'text-green-600 dark:text-green-400'
                      }`}
                    >
                      {(selectedStat.rxDropped + selectedStat.txDropped).toLocaleString()}
                    </span>
                  </div>
                </div>
              </div>
            </Card>
          </div>
        </>
      )}
    </div>
  );
}
