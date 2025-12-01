import { useEffect, useRef, useState } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';
import { lxcApi } from '@/api/lxc';

interface WebTerminalProps {
  containerName: string;
}

export function WebTerminal({ containerName }: WebTerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [commandHistory, setCommandHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);

  useEffect(() => {
    if (!terminalRef.current) return;

    // Create terminal instance
    const term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#d4d4d4',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#e5e5e5',
      },
      scrollback: 1000,
      tabStopWidth: 4,
    });

    // Add addons
    const fit = new FitAddon();
    const webLinks = new WebLinksAddon();

    term.loadAddon(fit);
    term.loadAddon(webLinks);

    // Open terminal
    term.open(terminalRef.current);
    fit.fit();

    // Welcome message
    term.writeln('\x1b[1;32m╔══════════════════════════════════════════════════════════════╗\x1b[0m');
    term.writeln('\x1b[1;32m║\x1b[0m  \x1b[1;36mStumpfWorks NAS - LXC Container Terminal\x1b[0m                    \x1b[1;32m║\x1b[0m');
    term.writeln('\x1b[1;32m╚══════════════════════════════════════════════════════════════╝\x1b[0m');
    term.writeln('');
    term.writeln(`\x1b[1;33mConnected to container:\x1b[0m ${containerName}`);
    term.writeln('\x1b[90mType your commands below. Use ↑/↓ for command history.\x1b[0m');
    term.writeln('');

    // Show prompt
    term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# `);

    // Handle resize
    const handleResize = () => {
      fit.fit();
    };
    window.addEventListener('resize', handleResize);

    // Handle input
    let currentLine = '';
    let cursorPosition = 0;

    term.onData((data) => {
      const code = data.charCodeAt(0);

      // Enter key
      if (code === 13) {
        const command = currentLine.trim();
        term.write('\r\n');

        if (command) {
          // Add to history
          setCommandHistory(prev => [...prev, command]);
          setHistoryIndex(-1);

          // Execute command
          executeCommand(term, containerName, command);
          currentLine = '';
          cursorPosition = 0;
        } else {
          term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# `);
        }
        return;
      }

      // Backspace
      if (code === 127) {
        if (cursorPosition > 0) {
          currentLine = currentLine.slice(0, cursorPosition - 1) + currentLine.slice(cursorPosition);
          cursorPosition--;
          term.write('\b \b');
          if (cursorPosition < currentLine.length) {
            term.write(currentLine.slice(cursorPosition) + ' ');
            term.write('\x1b[' + (currentLine.length - cursorPosition + 1) + 'D');
          }
        }
        return;
      }

      // Arrow up (previous command)
      if (data === '\x1b[A') {
        if (commandHistory.length > 0) {
          const newIndex = historyIndex < commandHistory.length - 1 ? historyIndex + 1 : historyIndex;
          const cmd = commandHistory[commandHistory.length - 1 - newIndex];
          if (cmd) {
            // Clear current line
            term.write('\r\x1b[K');
            term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# ${cmd}`);
            currentLine = cmd;
            cursorPosition = cmd.length;
            setHistoryIndex(newIndex);
          }
        }
        return;
      }

      // Arrow down (next command)
      if (data === '\x1b[B') {
        if (historyIndex > 0) {
          const newIndex = historyIndex - 1;
          const cmd = commandHistory[commandHistory.length - 1 - newIndex];
          // Clear current line
          term.write('\r\x1b[K');
          term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# ${cmd}`);
          currentLine = cmd;
          cursorPosition = cmd.length;
          setHistoryIndex(newIndex);
        } else if (historyIndex === 0) {
          // Clear line
          term.write('\r\x1b[K');
          term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# `);
          currentLine = '';
          cursorPosition = 0;
          setHistoryIndex(-1);
        }
        return;
      }

      // Ctrl+C
      if (code === 3) {
        term.write('^C\r\n');
        term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# `);
        currentLine = '';
        cursorPosition = 0;
        return;
      }

      // Ctrl+L (clear screen)
      if (code === 12) {
        term.clear();
        term.write(`\x1b[1;32mroot@${containerName}\x1b[0m:\x1b[1;34m~\x1b[0m# ${currentLine}`);
        if (cursorPosition < currentLine.length) {
          term.write('\x1b[' + (currentLine.length - cursorPosition) + 'D');
        }
        return;
      }

      // Regular characters
      if (code >= 32 && code < 127) {
        currentLine = currentLine.slice(0, cursorPosition) + data + currentLine.slice(cursorPosition);
        cursorPosition++;
        term.write(data);
        if (cursorPosition < currentLine.length) {
          term.write(currentLine.slice(cursorPosition) + '\x1b[' + (currentLine.length - cursorPosition) + 'D');
        }
      }
    });

    // Cleanup
    return () => {
      window.removeEventListener('resize', handleResize);
      term.dispose();
    };
  }, [containerName]);

  const executeCommand = async (term: Terminal, container: string, command: string) => {
    try {
      term.write(`\x1b[90mExecuting: ${command}\x1b[0m\r\n`);

      const response = await lxcApi.execCommand(container, command);

      if (response.success && response.data) {
        const { stdout, stderr, exit_code } = response.data;

        // Write stdout
        if (stdout) {
          const lines = stdout.split('\n');
          lines.forEach(line => {
            term.write(line + '\r\n');
          });
        }

        // Write stderr in red
        if (stderr) {
          const lines = stderr.split('\n');
          lines.forEach(line => {
            term.write(`\x1b[31m${line}\x1b[0m\r\n`);
          });
        }

        // Show exit code if non-zero
        if (exit_code !== 0) {
          term.write(`\x1b[31m[Exit code: ${exit_code}]\x1b[0m\r\n`);
        }
      } else {
        term.write(`\x1b[31mError: ${response.error?.message || 'Command failed'}\x1b[0m\r\n`);
      }
    } catch (err) {
      term.write(`\x1b[31mError: ${err instanceof Error ? err.message : 'Unknown error'}\x1b[0m\r\n`);
    }

    // Show prompt again
    term.write(`\x1b[1;32mroot@${container}\x1b[0m:\x1b[1;34m~\x1b[0m# `);
  };

  return (
    <div className="relative w-full h-full bg-[#1e1e1e] rounded-lg overflow-hidden">
      <div ref={terminalRef} className="w-full h-full p-2" />
    </div>
  );
}
