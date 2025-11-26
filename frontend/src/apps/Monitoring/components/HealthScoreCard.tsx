import { motion } from 'framer-motion';
import { Heart } from 'lucide-react';
import type { HealthScore } from '@/api/monitoring';
import Card from '@/components/ui/Card';

interface HealthScoreCardProps {
  healthScore: HealthScore;
}

export default function HealthScoreCard({ healthScore }: HealthScoreCardProps) {
  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'healthy':
        return 'text-green-500 bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800';
      case 'warning':
        return 'text-yellow-500 bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800';
      case 'critical':
        return 'text-red-500 bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800';
      default:
        return 'text-gray-500 bg-gray-50 dark:bg-gray-900/20 border-gray-200 dark:border-gray-800';
    }
  };

  const getScoreColor = (score: number): string => {
    if (score >= 80) return 'text-green-500';
    if (score >= 60) return 'text-yellow-500';
    return 'text-red-500';
  };

  const healthDetails = [
    { label: 'CPU Health', value: healthScore.details.cpu_health },
    { label: 'Memory Health', value: healthScore.details.memory_health },
    { label: 'Disk Health', value: healthScore.details.disk_health },
    { label: 'Network Health', value: healthScore.details.network_health },
  ];

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Card className={`p-6 border-2 ${getStatusColor(healthScore.status)}`}>
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <Heart className={`w-8 h-8 ${getScoreColor(healthScore.score)}`} />
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">System Health</h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Status: <span className={`font-medium ${getScoreColor(healthScore.score)}`}>
                  {healthScore.status.charAt(0).toUpperCase() + healthScore.status.slice(1)}
                </span>
              </p>
            </div>
          </div>

          <div className="text-right">
            <div className={`text-4xl font-bold ${getScoreColor(healthScore.score)}`}>
              {healthScore.score}
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Health Score</p>
          </div>
        </div>

        {/* Progress bar */}
        <div className="mb-4">
          <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${healthScore.score}%` }}
              transition={{ duration: 0.8, ease: 'easeOut' }}
              className={`h-full ${healthScore.score >= 80 ? 'bg-green-500' : healthScore.score >= 60 ? 'bg-yellow-500' : 'bg-red-500'}`}
            />
          </div>
        </div>

        {/* Health details */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {healthDetails.map((detail) => (
            <div key={detail.label} className="text-center">
              <div className={`text-lg font-semibold ${getScoreColor(detail.value)}`}>
                {detail.value}
              </div>
              <div className="text-xs text-gray-600 dark:text-gray-400">{detail.label}</div>
            </div>
          ))}
        </div>
      </Card>
    </motion.div>
  );
}
