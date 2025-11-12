import React from 'react';
import { motion } from 'framer-motion';

interface BreadcrumbsProps {
  currentPath: string;
  onNavigate: (path: string) => void;
}

const Breadcrumbs: React.FC<BreadcrumbsProps> = ({ currentPath, onNavigate }) => {
  const pathParts = currentPath.split('/').filter(Boolean);

  const buildPath = (index: number): string => {
    if (index === -1) return '/';
    return '/' + pathParts.slice(0, index + 1).join('/');
  };

  return (
    <div className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 text-sm overflow-x-auto">
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => onNavigate('/')}
        className="flex items-center gap-1 px-2 py-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>
        Root
      </motion.button>

      {pathParts.map((part, index) => (
        <React.Fragment key={index}>
          <svg className="w-4 h-4 text-gray-400 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => onNavigate(buildPath(index))}
            className={`px-2 py-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 ${
              index === pathParts.length - 1
                ? 'text-blue-600 dark:text-blue-400 font-medium'
                : 'text-gray-700 dark:text-gray-300'
            }`}
          >
            {part}
          </motion.button>
        </React.Fragment>
      ))}
    </div>
  );
};

export default Breadcrumbs;
