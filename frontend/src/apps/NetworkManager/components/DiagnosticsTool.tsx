import { useState } from 'react';
import { networkApi, DiagnosticResult } from '@/api/network';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';

type DiagnosticTool = 'ping' | 'traceroute' | 'netstat' | 'wol';

export default function DiagnosticsTool() {
  const [activeTool, setActiveTool] = useState<DiagnosticTool>('ping');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<DiagnosticResult | null>(null);

  // Ping state
  const [pingHost, setPingHost] = useState('');
  const [pingCount, setPingCount] = useState(4);

  // Traceroute state
  const [traceHost, setTraceHost] = useState('');

  // Netstat state
  const [netstatOptions, setNetstatOptions] = useState('-tuln');

  // Wake-on-LAN state
  const [wolMac, setWolMac] = useState('');
  const [wolSuccess, setWolSuccess] = useState<string | null>(null);

  const handlePing = async () => {
    if (!pingHost.trim()) {
      alert('Please enter a host to ping');
      return;
    }

    setLoading(true);
    setResult(null);

    try {
      const response = await networkApi.ping(pingHost, pingCount);
      if (response.success && response.data) {
        setResult(response.data);
      } else {
        alert(response.error?.message || 'Ping failed');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleTraceroute = async () => {
    if (!traceHost.trim()) {
      alert('Please enter a host to traceroute');
      return;
    }

    setLoading(true);
    setResult(null);

    try {
      const response = await networkApi.traceroute(traceHost);
      if (response.success && response.data) {
        setResult(response.data);
      } else {
        alert(response.error?.message || 'Traceroute failed');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleNetstat = async () => {
    setLoading(true);
    setResult(null);

    try {
      const response = await networkApi.netstat(netstatOptions);
      if (response.success && response.data) {
        setResult(response.data);
      } else {
        alert(response.error?.message || 'Netstat failed');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleWakeOnLAN = async () => {
    if (!wolMac.trim()) {
      alert('Please enter a MAC address');
      return;
    }

    setLoading(true);
    setWolSuccess(null);

    try {
      const response = await networkApi.wakeOnLAN(wolMac);
      if (response.success) {
        setWolSuccess('Wake-on-LAN packet sent successfully! 📡');
      } else {
        alert(response.error?.message || 'Failed to send WOL packet');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const tools = [
    { id: 'ping' as DiagnosticTool, name: 'Ping', icon: '🏓', description: 'Test network connectivity' },
    { id: 'traceroute' as DiagnosticTool, name: 'Traceroute', icon: '🗺️', description: 'Trace packet route' },
    { id: 'netstat' as DiagnosticTool, name: 'Netstat', icon: '📊', description: 'Network connections' },
    { id: 'wol' as DiagnosticTool, name: 'Wake-on-LAN', icon: '⏰', description: 'Wake remote devices' },
  ];

  return (
    <div className="p-6 space-y-6 max-w-6xl">
      {/* Tool Selection */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {tools.map((tool) => (
          <button
            key={tool.id}
            onClick={() => {
              setActiveTool(tool.id);
              setResult(null);
              setWolSuccess(null);
            }}
            className={`p-4 rounded-lg border-2 transition-all ${
              activeTool === tool.id
                ? 'border-macos-blue bg-blue-50 dark:bg-blue-900/20'
                : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 bg-white dark:bg-macos-dark-100'
            }`}
          >
            <div className="text-3xl mb-2">{tool.icon}</div>
            <div className="font-semibold text-gray-900 dark:text-gray-100">{tool.name}</div>
            <div className="text-xs text-gray-600 dark:text-gray-400 mt-1">{tool.description}</div>
          </button>
        ))}
      </div>

      {/* Ping Tool */}
      {activeTool === 'ping' && (
        <Card>
          <div className="p-6">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
              <span>🏓</span>
              Ping
            </h2>
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="md:col-span-2">
                  <Input
                    label="Host or IP Address"
                    value={pingHost}
                    onChange={(e) => setPingHost(e.target.value)}
                    placeholder="google.com or 8.8.8.8"
                  />
                </div>
                <div>
                  <Input
                    label="Count"
                    type="number"
                    value={pingCount}
                    onChange={(e) => setPingCount(parseInt(e.target.value) || 4)}
                    min={1}
                    max={100}
                  />
                </div>
              </div>
              <Button onClick={handlePing} disabled={loading}>
                {loading ? 'Pinging...' : 'Run Ping'}
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* Traceroute Tool */}
      {activeTool === 'traceroute' && (
        <Card>
          <div className="p-6">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
              <span>🗺️</span>
              Traceroute
            </h2>
            <div className="space-y-4">
              <Input
                label="Host or IP Address"
                value={traceHost}
                onChange={(e) => setTraceHost(e.target.value)}
                placeholder="google.com or 8.8.8.8"
              />
              <Button onClick={handleTraceroute} disabled={loading}>
                {loading ? 'Tracing...' : 'Run Traceroute'}
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* Netstat Tool */}
      {activeTool === 'netstat' && (
        <Card>
          <div className="p-6">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
              <span>📊</span>
              Netstat
            </h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Options
                </label>
                <select
                  value={netstatOptions}
                  onChange={(e) => setNetstatOptions(e.target.value)}
                  className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                >
                  <option value="-tuln">All listening ports (-tuln)</option>
                  <option value="-tulpn">All listening ports with process (-tulpn)</option>
                  <option value="-tun">All connections (-tun)</option>
                  <option value="-r">Routing table (-r)</option>
                  <option value="-i">Interface statistics (-i)</option>
                  <option value="-s">Protocol statistics (-s)</option>
                </select>
              </div>
              <Button onClick={handleNetstat} disabled={loading}>
                {loading ? 'Running...' : 'Run Netstat'}
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* Wake-on-LAN Tool */}
      {activeTool === 'wol' && (
        <Card>
          <div className="p-6">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
              <span>⏰</span>
              Wake-on-LAN
            </h2>
            <div className="space-y-4">
              <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg text-blue-700 dark:text-blue-400 text-sm">
                <p className="font-medium mb-1">About Wake-on-LAN</p>
                <p>
                  Send a magic packet to wake a device from sleep or powered-off state. The target
                  device must support and have Wake-on-LAN enabled in BIOS/UEFI.
                </p>
              </div>
              <Input
                label="MAC Address"
                value={wolMac}
                onChange={(e) => setWolMac(e.target.value)}
                placeholder="AA:BB:CC:DD:EE:FF or AA-BB-CC-DD-EE-FF"
              />
              <Button onClick={handleWakeOnLAN} disabled={loading}>
                {loading ? 'Sending...' : 'Send Magic Packet'}
              </Button>
              {wolSuccess && (
                <div className="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg text-green-700 dark:text-green-400">
                  {wolSuccess}
                </div>
              )}
            </div>
          </div>
        </Card>
      )}

      {/* Results Display */}
      {result && (
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
              Results
            </h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600 dark:text-gray-400">Command:</span>
                <span className="font-mono text-gray-900 dark:text-gray-100">
                  {result.command}
                </span>
              </div>
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600 dark:text-gray-400">Status:</span>
                <span
                  className={`px-2 py-1 rounded text-xs font-medium ${
                    result.success
                      ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
                  }`}
                >
                  {result.success ? 'Success' : 'Failed'}
                </span>
              </div>
            </div>

            {result.error && (
              <div className="mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
                {result.error}
              </div>
            )}

            <div className="mt-4">
              <div className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Output:
              </div>
              <pre className="p-4 bg-gray-900 dark:bg-black text-green-400 rounded-lg overflow-x-auto text-xs font-mono max-h-96 overflow-y-auto">
                {result.output || 'No output'}
              </pre>
            </div>
          </div>
        </Card>
      )}
    </div>
  );
}
