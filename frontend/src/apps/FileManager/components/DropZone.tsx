import { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface DropZoneProps {
  onFilesDropped: (files: FileList) => void;
  children: React.ReactNode;
  disabled?: boolean;
}

export default function DropZone({ onFilesDropped, children, disabled = false }: DropZoneProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [, setDragCounter] = useState(0); // Tracks nested drag events (value only used in setState callback)

  const handleDragEnter = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (disabled) return;

    setDragCounter(prev => prev + 1);
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  }, [disabled]);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (disabled) return;

    setDragCounter(prev => {
      const newCounter = prev - 1;
      if (newCounter === 0) {
        setIsDragging(false);
      }
      return newCounter;
    });
  }, [disabled]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (disabled) return;

    e.dataTransfer.dropEffect = 'copy';
  }, [disabled]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (disabled) return;

    setIsDragging(false);
    setDragCounter(0);

    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      onFilesDropped(e.dataTransfer.files);
    }
  }, [disabled, onFilesDropped]);

  return (
    <div
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
      className="relative h-full"
    >
      {children}

      <AnimatePresence>
        {isDragging && !disabled && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 z-50 pointer-events-none"
          >
            {/* Overlay */}
            <div className="absolute inset-0 bg-blue-500/20 dark:bg-blue-400/20 backdrop-blur-sm" />

            {/* Center content */}
            <div className="absolute inset-0 flex items-center justify-center">
              <motion.div
                initial={{ scale: 0.9 }}
                animate={{ scale: 1 }}
                className="bg-white dark:bg-macos-dark-100 rounded-2xl shadow-2xl p-12 border-4 border-dashed border-blue-500 dark:border-blue-400"
              >
                <div className="text-center">
                  <div className="text-6xl mb-4">📤</div>
                  <div className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                    Drop files here
                  </div>
                  <div className="text-gray-600 dark:text-gray-400">
                    Release to upload
                  </div>
                </div>
              </motion.div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
