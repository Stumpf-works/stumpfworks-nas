import { useRef, useState } from 'react';
import { motion } from 'framer-motion';
import { useWindowStore } from '@/store';
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
  const windowRef = useRef<HTMLDivElement>(null);

  const handleMouseDown = () => {
    focusWindow(window.id);
  };

  const handleDrag = (event: any, info: any) => {
    setIsDragging(true);
    updateWindowPosition(window.id, {
      x: window.position.x + info.delta.x,
      y: window.position.y + info.delta.y,
    });
  };

  const handleDragEnd = () => {
    setIsDragging(false);
  };

  const isMaximized = window.state === 'maximized';

  return (
    <motion.div
      ref={windowRef}
      initial={{ scale: 0.9, opacity: 0 }}
      animate={{
        scale: 1,
        opacity: 1,
        x: isMaximized ? 0 : window.position.x,
        y: isMaximized ? 0 : window.position.y,
        width: isMaximized ? '100%' : window.size.width,
        height: isMaximized ? '100%' : window.size.height,
      }}
      exit={{ scale: 0.9, opacity: 0 }}
      drag={!isMaximized}
      dragMomentum={false}
      dragElastic={0}
      dragConstraints={{ left: 0, top: 0, right: window.innerWidth - 400, bottom: window.innerHeight - 300 }}
      onDrag={handleDrag}
      onDragEnd={handleDragEnd}
      onMouseDown={handleMouseDown}
      style={{
        position: 'absolute',
        zIndex: window.zIndex,
      }}
      className={`flex flex-col bg-white dark:bg-macos-dark-100 rounded-lg shadow-window overflow-hidden ${
        window.isFocused ? 'ring-2 ring-macos-blue/30' : ''
      }`}
      transition={{ type: 'spring', stiffness: 300, damping: 30 }}
    >
      {/* Title Bar */}
      <div
        className="flex items-center justify-between h-10 px-4 bg-gray-100 dark:bg-macos-dark-200 border-b border-gray-200 dark:border-gray-700 cursor-move select-none"
        onMouseEnter={() => setIsHoveringTitleBar(true)}
        onMouseLeave={() => setIsHoveringTitleBar(false)}
      >
        {/* Traffic Lights */}
        <div className="flex items-center space-x-2">
          <button
            onClick={() => closeWindow(window.id)}
            className="group w-3 h-3 rounded-full bg-red-500 hover:bg-red-600 flex items-center justify-center"
          >
            {isHoveringTitleBar && (
              <svg className="w-2 h-2 text-red-900 opacity-0 group-hover:opacity-100" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clipRule="evenodd"
                />
              </svg>
            )}
          </button>
          <button
            onClick={() => minimizeWindow(window.id)}
            className="group w-3 h-3 rounded-full bg-yellow-500 hover:bg-yellow-600 flex items-center justify-center"
          >
            {isHoveringTitleBar && (
              <div className="w-2 h-0.5 bg-yellow-900 opacity-0 group-hover:opacity-100" />
            )}
          </button>
          <button
            onClick={() => maximizeWindow(window.id)}
            className="group w-3 h-3 rounded-full bg-green-500 hover:bg-green-600 flex items-center justify-center"
          >
            {isHoveringTitleBar && (
              <svg className="w-2 h-2 text-green-900 opacity-0 group-hover:opacity-100" fill="none" stroke="currentColor" viewBox="0 0 20 20">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 12v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
              </svg>
            )}
          </button>
        </div>

        {/* Title */}
        <div className="absolute left-1/2 transform -translate-x-1/2 text-sm font-medium text-gray-700 dark:text-gray-300">
          {window.title}
        </div>

        <div />
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto">
        {/* Window content will be rendered here dynamically */}
        <div className="w-full h-full">
          {/* Placeholder for now */}
        </div>
      </div>
    </motion.div>
  );
}
