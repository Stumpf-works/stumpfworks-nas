import { ReactNode } from 'react';
import { motion } from 'framer-motion';

interface CardProps {
  children: ReactNode;
  className?: string;
  hoverable?: boolean;
}

export default function Card({ children, className = '', hoverable = false }: CardProps) {
  return (
    <motion.div
      whileHover={hoverable ? { y: -2, boxShadow: '0 10px 40px rgba(0, 0, 0, 0.15)' } : {}}
      className={`bg-white dark:bg-macos-dark-100 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 ${className}`}
    >
      {children}
    </motion.div>
  );
}
