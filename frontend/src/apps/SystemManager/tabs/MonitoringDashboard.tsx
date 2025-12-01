import { MonitoringWidgets } from '@/components/MonitoringWidgets/MonitoringWidgets';

export default function MonitoringDashboard() {
  return (
    <div className="p-6 space-y-6 bg-gray-50 dark:bg-macos-dark-50">
      <div>
        <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
          System Monitoring
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Real-time system metrics, health score, and performance trends
        </p>
      </div>

      <MonitoringWidgets timeRange="24h" />
    </div>
  );
}
