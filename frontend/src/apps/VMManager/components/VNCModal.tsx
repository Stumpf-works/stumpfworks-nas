import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Monitor, AlertCircle, Copy, CheckCircle } from 'lucide-react';
import { vmsApi } from '@/api/vms';
import { getErrorMessage } from '@/api/client';

interface VNCModalProps {
  isOpen: boolean;
  onClose: () => void;
  vmId: string;
  vmName: string;
}

export function VNCModal({ isOpen, onClose, vmId, vmName }: VNCModalProps) {
  const [vncPort, setVncPort] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (isOpen) {
      loadVNCPort();
    }
  }, [isOpen, vmId]);

  const loadVNCPort = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await vmsApi.getVNCPort(vmId);

      if (response.success && response.data) {
        setVncPort(response.data.port);
      } else {
        setError(response.error?.message || 'Failed to get VNC port');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (!isOpen) return null;

  const serverHost = window.location.hostname;
  const vncUrl = vncPort ? `vnc://${serverHost}:${vncPort}` : '';

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      >
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl max-w-2xl w-full mx-4"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
            <div className="flex items-center gap-3">
              <Monitor className="w-6 h-6 text-macos-blue" />
              <div>
                <h2 className="text-xl font-bold text-gray-900 dark:text-white">VNC Console Access</h2>
                <p className="text-sm text-gray-600 dark:text-gray-400">{vmName}</p>
              </div>
            </div>
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            >
              <X className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            </button>
          </div>

          {/* Content */}
          <div className="p-6">
            {loading ? (
              <div className="flex flex-col items-center justify-center py-12">
                <div className="w-12 h-12 border-4 border-macos-blue border-t-transparent rounded-full animate-spin mb-4" />
                <p className="text-gray-600 dark:text-gray-400">Loading VNC information...</p>
              </div>
            ) : error ? (
              <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-start gap-3">
                <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
                <div>
                  <h3 className="font-semibold text-red-900 dark:text-red-200">Error</h3>
                  <p className="text-sm text-red-700 dark:text-red-300 mt-1">{error}</p>
                </div>
              </div>
            ) : (
              <div className="space-y-6">
                {/* VNC Port Info */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    VNC Port
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={vncPort || ''}
                      readOnly
                      className="flex-1 px-4 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white font-mono"
                    />
                    <button
                      onClick={() => handleCopy(String(vncPort))}
                      className="p-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
                      title="Copy port"
                    >
                      {copied ? (
                        <CheckCircle className="w-5 h-5 text-green-600" />
                      ) : (
                        <Copy className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Server Address */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Server Address
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={serverHost}
                      readOnly
                      className="flex-1 px-4 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white font-mono"
                    />
                    <button
                      onClick={() => handleCopy(serverHost)}
                      className="p-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
                      title="Copy address"
                    >
                      {copied ? (
                        <CheckCircle className="w-5 h-5 text-green-600" />
                      ) : (
                        <Copy className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                      )}
                    </button>
                  </div>
                </div>

                {/* VNC URL */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    VNC URL
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={vncUrl}
                      readOnly
                      className="flex-1 px-4 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white font-mono"
                    />
                    <button
                      onClick={() => handleCopy(vncUrl)}
                      className="p-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
                      title="Copy URL"
                    >
                      {copied ? (
                        <CheckCircle className="w-5 h-5 text-green-600" />
                      ) : (
                        <Copy className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Instructions */}
                <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                  <h3 className="font-semibold text-blue-900 dark:text-blue-200 mb-2">How to Connect</h3>
                  <ol className="text-sm text-blue-700 dark:text-blue-300 space-y-2 list-decimal list-inside">
                    <li>Install a VNC client (TigerVNC, RealVNC, etc.)</li>
                    <li>Open your VNC client</li>
                    <li>Connect to: <code className="px-1 py-0.5 bg-blue-100 dark:bg-blue-900/50 rounded font-mono">{serverHost}:{vncPort}</code></li>
                    <li>Enter credentials if prompted</li>
                  </ol>
                </div>
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-3 p-6 border-t border-gray-200 dark:border-gray-700">
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
            >
              Close
            </button>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
}
