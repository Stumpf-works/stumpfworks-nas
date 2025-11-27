import { useEffect, useRef, useState } from 'react';
import { motion } from 'framer-motion';

export function Terminal() {
  const terminalRef = useRef<HTMLDivElement>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const [output, setOutput] = useState<string[]>([
    '╔══════════════════════════════════════════════════════════════╗',
    '║         StumpfWorks NAS Terminal                             ║',
    '╚══════════════════════════════════════════════════════════════╝',
    '',
    'Type your commands below. Use Ctrl+C to interrupt running commands.',
    '',
  ]);
  const [currentInput, setCurrentInput] = useState('');
  const [commandHistory, setCommandHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [isConnected, setIsConnected] = useState(false);
  const [currentDir, setCurrentDir] = useState('~');

  useEffect(() => {
    // WebSocket connection for real-time terminal
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/terminal`;

    try {
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        setIsConnected(true);
        addOutput('✓ Connected to terminal session', 'success');
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);

          if (data.type === 'output') {
            addOutput(data.data, 'output');
          } else if (data.type === 'error') {
            addOutput(data.data, 'error');
          } else if (data.type === 'cwd') {
            setCurrentDir(data.data);
          }
        } catch (err) {
          // Raw output
          addOutput(event.data, 'output');
        }
      };

      ws.onerror = () => {
        setIsConnected(false);
        addOutput('✗ WebSocket connection error - falling back to simulation mode', 'error');
        addOutput('Note: Commands will be simulated. For real terminal access, ensure the backend WebSocket endpoint is running.', 'warning');
      };

      ws.onclose = () => {
        setIsConnected(false);
        addOutput('✗ Terminal session disconnected', 'warning');
      };

      return () => {
        ws.close();
      };
    } catch (err) {
      setIsConnected(false);
      addOutput('✗ WebSocket not available - using simulation mode', 'error');
      addOutput('Commands will be simulated. This is a demonstration of the Terminal UI.', 'warning');
    }
  }, []);

  useEffect(() => {
    // Auto-scroll to bottom when output changes
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [output]);

  useEffect(() => {
    // Focus input when component mounts
    inputRef.current?.focus();
  }, []);

  const addOutput = (text: string, type: 'output' | 'error' | 'success' | 'warning' = 'output') => {
    setOutput(prev => [...prev, `${type === 'error' ? '❌ ' : type === 'success' ? '✓ ' : type === 'warning' ? '⚠️  ' : ''}${text}`]);
  };

  const executeCommand = (cmd: string) => {
    if (!cmd.trim()) return;

    // Add to history
    setCommandHistory(prev => [...prev, cmd]);
    setHistoryIndex(-1);

    // Echo command
    addOutput(`${currentDir} $ ${cmd}`, 'output');

    if (isConnected && wsRef.current?.readyState === WebSocket.OPEN) {
      // Send to backend via WebSocket
      wsRef.current.send(JSON.stringify({
        type: 'command',
        data: cmd,
      }));
    } else {
      // Simulation mode - handle some basic commands locally
      handleSimulatedCommand(cmd);
    }

    setCurrentInput('');
  };

  const handleSimulatedCommand = (cmd: string) => {
    const trimmed = cmd.trim().toLowerCase();

    if (trimmed === 'help') {
      addOutput('Available simulated commands:', 'output');
      addOutput('  help          - Show this help message', 'output');
      addOutput('  clear         - Clear the terminal', 'output');
      addOutput('  date          - Show current date and time', 'output');
      addOutput('  whoami        - Show current user', 'output');
      addOutput('  uptime        - Show system uptime simulation', 'output');
      addOutput('', 'output');
      addOutput('Note: This is simulation mode. Connect to backend for full terminal access.', 'warning');
    } else if (trimmed === 'clear') {
      setOutput([]);
    } else if (trimmed === 'date') {
      addOutput(new Date().toString(), 'output');
    } else if (trimmed === 'whoami') {
      addOutput('root', 'output');
    } else if (trimmed === 'uptime') {
      addOutput(' 12:34:56 up 7 days, 3:21, 1 user, load average: 0.15, 0.20, 0.18', 'output');
    } else if (trimmed === 'ls' || trimmed.startsWith('ls ')) {
      addOutput('bin  boot  dev  etc  home  lib  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var', 'output');
    } else if (trimmed === 'pwd') {
      addOutput('/root', 'output');
    } else if (trimmed.startsWith('echo ')) {
      addOutput(cmd.substring(5), 'output');
    } else if (trimmed === '') {
      // Empty command, do nothing
    } else {
      addOutput(`${trimmed}: command not found (simulation mode)`, 'error');
      addOutput('Type "help" for available commands or connect backend for full terminal.', 'warning');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      executeCommand(currentInput);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (commandHistory.length > 0) {
        const newIndex = historyIndex === -1 ? commandHistory.length - 1 : Math.max(0, historyIndex - 1);
        setHistoryIndex(newIndex);
        setCurrentInput(commandHistory[newIndex]);
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (historyIndex !== -1) {
        const newIndex = historyIndex + 1;
        if (newIndex >= commandHistory.length) {
          setHistoryIndex(-1);
          setCurrentInput('');
        } else {
          setHistoryIndex(newIndex);
          setCurrentInput(commandHistory[newIndex]);
        }
      }
    } else if (e.key === 'c' && e.ctrlKey) {
      e.preventDefault();
      addOutput('^C', 'output');
      setCurrentInput('');

      if (isConnected && wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({
          type: 'interrupt',
        }));
      }
    } else if (e.key === 'l' && e.ctrlKey) {
      e.preventDefault();
      setOutput([]);
    }
  };

  const handleClear = () => {
    setOutput([
      '╔══════════════════════════════════════════════════════════════╗',
      '║         StumpfWorks NAS Terminal                             ║',
      '╚══════════════════════════════════════════════════════════════╝',
      '',
    ]);
  };

  return (
    <div className="flex flex-col h-full bg-gray-900 text-green-400 font-mono">
      {/* Header */}
      <div className="p-3 md:p-4 border-b border-gray-700 bg-gray-800 flex justify-between items-center">
        <div className="flex items-center gap-2 md:gap-3">
          <span className="text-xl md:text-2xl">💻</span>
          <div>
            <h1 className="text-base md:text-lg font-bold text-gray-100">Terminal</h1>
            <p className="text-xs text-gray-400">
              {isConnected ? (
                <span className="text-green-400">● Connected</span>
              ) : (
                <span className="text-yellow-400">● Simulation Mode</span>
              )}
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleClear}
            className="px-2 md:px-3 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded text-xs md:text-sm transition-colors"
            title="Clear terminal (Ctrl+L)"
          >
            <span className="hidden sm:inline">Clear</span>
            <span className="sm:hidden">🗑️</span>
          </button>
        </div>
      </div>

      {/* Terminal Output */}
      <div
        ref={terminalRef}
        className="flex-1 overflow-y-auto p-3 md:p-4 space-y-1 text-xs md:text-sm"
        onClick={() => inputRef.current?.focus()}
      >
        {output.map((line, index) => (
          <div key={index} className="whitespace-pre-wrap break-all">
            {line}
          </div>
        ))}

        {/* Current Input Line */}
        <div className="flex items-center gap-2">
          <span className="text-blue-400">{currentDir}</span>
          <span className="text-gray-400">$</span>
          <input
            ref={inputRef}
            type="text"
            value={currentInput}
            onChange={(e) => setCurrentInput(e.target.value)}
            onKeyDown={handleKeyDown}
            className="flex-1 bg-transparent outline-none text-green-400"
            autoComplete="off"
            spellCheck={false}
          />
          <motion.span
            animate={{ opacity: [1, 0] }}
            transition={{ duration: 0.8, repeat: Infinity }}
            className="w-2 h-4 bg-green-400"
          />
        </div>
      </div>

      {/* Footer Hints */}
      <div className="px-3 md:px-4 py-2 border-t border-gray-700 bg-gray-800 text-xs text-gray-500 flex flex-col sm:flex-row justify-between gap-2">
        <div className="flex flex-wrap gap-2 sm:gap-4">
          <span className="hidden sm:inline">Ctrl+C: Interrupt</span>
          <span className="hidden sm:inline">Ctrl+L: Clear</span>
          <span className="hidden sm:inline">↑↓: History</span>
          <span className="sm:hidden">^C: Interrupt | ^L: Clear | ↑↓: History</span>
        </div>
        <div>
          {commandHistory.length} command{commandHistory.length !== 1 ? 's' : ''} in history
        </div>
      </div>
    </div>
  );
}
