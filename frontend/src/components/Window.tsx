import { useRef, useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useWindowStore } from '@/store';
import { getAppById } from '@/apps';
import type { Window as WindowType } from '@/types';

interface WindowProps {
  window: WindowType;
}

export default function Window({ window }: WindowProps) {
  const closeWindow = useWindowStore((state) => state.closeWindow);
  const focusWindow = useWindowStore((state) => state.focusWindow);
  const minimizeWindow = useWindowStore((state) => state.minimizeWindow);
  const maximizeWindow = useWindowStore((state) => state.maximizeWindow);
  const updateWindowPosition = useWindowStore((state) => state.updateWindowPosition);
  const updateWindowSize = useWindowStore((state) => state.updateWindowSize);

  const [isHoveringTitleBar, setIsHoveringTitleBar] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const windowRef = useRef<HTMLDivElement>(null);
  const titleBarRef = useRef<HTMLDivElement>(null);

  const handleMouseDown = () => {
    focusWindow(window.id);
  };

  // Custom drag implementation for title bar only
  const handleTitleBarMouseDown = (e: React.MouseEvent) => {
    if (!isMaximized && titleBarRef.current) {
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(true);
      setDragStart({
        x: e.clientX - window.position.x,
        y: e.clientY - window.position.y,
      });
      focusWindow(window.id);
    }
  };

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (isDragging && !isMaximized) {
        const newX = e.clientX - dragStart.x;
        const newY = Math.max(0, e.clientY - dragStart.y); // Prevent dragging above screen

        updateWindowPosition(window.id, {
          x: Math.max(0, Math.min(newX, window.innerWidth - window.size.width)),
          y: newY,
        });
      }
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging, dragStart, window.id, window.size.width, isMaximized, updateWindowPosition]);

  const isMaximized = window.state === 'maximized';

  // Get the app component
  const app = getAppById(window.appId);
  const AppComponent = app?.component;

  return (
    <motion.div
      ref={windowRef}
      initial={{ scale: 0.95, opacity: 0, y: 20 }}
      animate={{
        scale: 1,
        opacity: 1,
        x: isMaximized ? 0 : window.position.x,
        y: isMaximized ? 0 : window.position.y,
        width: isMaximized ? '100%' : window.size.width,
        height: isMaximized ? 'calc(100vh - 64px)' : window.size.height,
      }}
      exit={{ scale: 0.95, opacity: 0, y: 20 }}
      onMouseDown={handleMouseDown}
      style={{
        position: 'absolute',
        zIndex: window.zIndex,
        top: isMaximized ? '64px' : 0,
        left: isMaximized ? 0 : 0,
      }}
      className={`flex flex-col bg-white dark:bg-gray-800 rounded-xl shadow-2xl overflow-hidden border ${
        window.isFocused
          ? 'border-blue-500/50 shadow-blue-500/20'
          : 'border-gray-200 dark:border-gray-700'
      }`}
      transition={{ type: 'spring', stiffness: 400, damping: 35 }}
    >
      {/* Title Bar - macOS Style */}
      <div
        ref={titleBarRef}
        className={`flex items-center justify-between h-11 px-4 ${
          window.isFocused
            ? 'bg-gradient-to-b from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-800'
            : 'bg-gray-100 dark:bg-gray-800'
        } border-b border-gray-200 dark:border-gray-700 select-none ${
          !isMaximized ? 'cursor-move' : ''
        }`}
        onMouseDown={handleTitleBarMouseDown}
        onMouseEnter={() => setIsHoveringTitleBar(true)}
        onMouseLeave={() => setIsHoveringTitleBar(false)}
      >
        {/* Traffic Lights - macOS Style */}
        <div className="flex items-center space-x-2">
          <button
            onClick={(e) => {
              e.stopPropagation();
              closeWindow(window.id);
            }}
            className="group w-3 h-3 rounded-full bg-red-500 hover:bg-red-600 flex items-center justify-center transition-all duration-150 hover:scale-110"
            aria-label="Close"
          >
            {isHoveringTitleBar && (
              <svg
                className="w-2 h-2 text-red-900 opacity-0 group-hover:opacity-100 transition-opacity"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clipRule="evenodd"
                />
              </svg>
            )}
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              minimizeWindow(window.id);
            }}
            className="group w-3 h-3 rounded-full bg-yellow-500 hover:bg-yellow-600 flex items-center justify-center transition-all duration-150 hover:scale-110"
            aria-label="Minimize"
          >
            {isHoveringTitleBar && (
              <div className="w-2 h-0.5 bg-yellow-900 opacity-0 group-hover:opacity-100 transition-opacity" />
            )}
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              maximizeWindow(window.id);
            }}
            className="group w-3 h-3 rounded-full bg-green-500 hover:bg-green-600 flex items-center justify-center transition-all duration-150 hover:scale-110"
            aria-label="Maximize"
          >
            {isHoveringTitleBar && (
              <svg
                className="w-2 h-2 text-green-900 opacity-0 group-hover:opacity-100 transition-opacity"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 12v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"
                />
              </svg>
            )}
          </button>
        </div>

        {/* Title - Centered */}
        <div className="absolute left-1/2 transform -translate-x-1/2 text-sm font-semibold text-gray-700 dark:text-gray-200 flex items-center">
          {window.icon && <span className="mr-2 text-base">{window.icon}</span>}
          {window.title}
        </div>

        <div className="w-16" /> {/* Spacer for symmetry */}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        {AppComponent ? (
          <AppComponent />
        ) : (
          <div className="flex items-center justify-center h-full bg-white dark:bg-gray-900">
            <p className="text-gray-500 dark:text-gray-400">App not found</p>
          </div>
        )}
      </div>
    </motion.div>
  );
}
