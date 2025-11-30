import { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Terminal, Send, Trash2 } from 'lucide-react';
import { lxcApi } from '@/api/lxc';
import { getErrorMessage } from '@/api/client';

interface ConsoleModalProps {
  isOpen: boolean;
  onClose: () => void;
  containerName: string;
}

interface CommandEntry {
  command: string;
  output: string;
  error: string;
  timestamp: Date;
}

export function ConsoleModal({ isOpen, onClose, containerName }: ConsoleModalProps) {
  const [command, setCommand] = useState('');
  const [history, setHistory] = useState<CommandEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const outputRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when history updates
  useEffect(() => {
    if (outputRef.current) {
      outputRef.current.scrollTop = outputRef.current.scrollHeight;
    }
  }, [history]);

  const handleExecute = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();

    if (!command.trim()) return;

    try {
      setLoading(true);
      const response = await lxcApi.execCommand(containerName, command);

      if (response.success && response.data) {
        setHistory(prev => [
          ...prev,
          {
            command,
            output: response.data!.stdout,
            error: response.data!.stderr,
            timestamp: new Date(),
          },
        ]);
        setCommand('');
      }
    } catch (err) {
      setHistory(prev => [
        ...prev,
        {
          command,
          output: '',
          error: getErrorMessage(err),
          timestamp: new Date(),
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  const handleClear = () => {
    setHistory([]);
  };

  if (!isOpen) return null;

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
          className="bg-gray-900 rounded-xl shadow-2xl max-w-4xl w-full mx-4 h-[80vh] flex flex-col"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-gray-700">
            <div className="flex items-center gap-3">
              <Terminal className="w-5 h-5 text-green-400" />
              <div>
                <h2 className="text-lg font-bold text-white">Console</h2>
                <p className="text-sm text-gray-400">{containerName}</p>
              </div>
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleClear}
                className="p-2 hover:bg-gray-800 rounded-lg transition-colors text-gray-400 hover:text-white"
                title="Clear history"
              >
                <Trash2 className="w-4 h-4" />
              </button>
              <button
                onClick={onClose}
                className="p-2 hover:bg-gray-800 rounded-lg transition-colors text-gray-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
          </div>

          {/* Output Area */}
          <div
            ref={outputRef}
            className="flex-1 overflow-auto p-4 font-mono text-sm bg-black/50"
          >
            {history.length === 0 ? (
              <div className="text-gray-500 text-center py-8">
                <Terminal className="w-12 h-12 mx-auto mb-2 opacity-50" />
                <p>No commands executed yet</p>
                <p className="text-xs mt-1">Type a command below to execute in the container</p>
              </div>
            ) : (
              <div className="space-y-4">
                {history.map((entry, index) => (
                  <div key={index} className="space-y-1">
                    {/* Command */}
                    <div className="flex items-start gap-2">
                      <span className="text-green-400 select-none">$</span>
                      <span className="text-white">{entry.command}</span>
                      <span className="text-gray-600 text-xs ml-auto">
                        {entry.timestamp.toLocaleTimeString()}
                      </span>
                    </div>

                    {/* Output */}
                    {entry.output && (
                      <pre className="text-gray-300 whitespace-pre-wrap pl-4 border-l-2 border-gray-700">
                        {entry.output}
                      </pre>
                    )}

                    {/* Error */}
                    {entry.error && (
                      <pre className="text-red-400 whitespace-pre-wrap pl-4 border-l-2 border-red-900">
                        {entry.error}
                      </pre>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Input Area */}
          <form onSubmit={handleExecute} className="p-4 border-t border-gray-700">
            <div className="flex gap-2">
              <div className="flex-1 flex items-center gap-2 bg-gray-800 rounded-lg px-3 py-2">
                <span className="text-green-400 select-none">$</span>
                <input
                  type="text"
                  value={command}
                  onChange={(e) => setCommand(e.target.value)}
                  placeholder="Enter command..."
                  className="flex-1 bg-transparent border-none focus:outline-none text-white font-mono placeholder-gray-500"
                  autoFocus
                  disabled={loading}
                />
              </div>
              <button
                type="submit"
                disabled={!command.trim() || loading}
                className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              >
                {loading ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                    Running...
                  </>
                ) : (
                  <>
                    <Send className="w-4 h-4" />
                    Execute
                  </>
                )}
              </button>
            </div>
            <p className="text-xs text-gray-500 mt-2">
              Tip: Press Enter to execute, or type <code className="px-1 py-0.5 bg-gray-800 rounded">ls</code>,{' '}
              <code className="px-1 py-0.5 bg-gray-800 rounded">ps aux</code>, etc.
            </p>
          </form>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
}
