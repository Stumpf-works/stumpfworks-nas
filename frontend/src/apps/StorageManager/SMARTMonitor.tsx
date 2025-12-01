// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  HardDrive,
  RefreshCw,
  Activity,
  Thermometer,
  Clock,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Zap,
  PlayCircle,
} from 'lucide-react';
import { syslibApi, type SMARTInfo } from '@/api/syslib';
import { storageApi, type Disk } from '@/api/storage';

export default function SMARTMonitor() {
  const [disks, setDisks] = useState<Disk[]>([]);
  const [selectedDisk, setSelectedDisk] = useState<Disk | null>(null);
  const [smartInfo, setSmartInfo] = useState<SMARTInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingSMART, setIsLoadingSMART] = useState(false);
  const [isRunningTest, setIsRunningTest] = useState(false);

  // Fetch disks
  const fetchDisks = async () => {
    setIsLoading(true);
    try {
      const response = await storageApi.listDisks();
      if (response.success && response.data) {
        setDisks(response.data);
        if (response.data.length > 0 && !selectedDisk) {
          setSelectedDisk(response.data[0]);
        }
      }
    } catch (error) {
      console.error('Failed to fetch disks:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch SMART info for selected disk
  const fetchSMARTInfo = async (device: string) => {
    setIsLoadingSMART(true);
    try {
      const response = await syslibApi.smart.getInfo(device);
      if (response.success && response.data) {
        setSmartInfo(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch SMART info:', error);
    } finally {
      setIsLoadingSMART(false);
    }
  };

  useEffect(() => {
    fetchDisks();
  }, []);

  useEffect(() => {
    if (selectedDisk) {
      fetchSMARTInfo(selectedDisk.path);
    }
  }, [selectedDisk]);

  const handleRunTest = async (testType: 'short' | 'long' | 'conveyance') => {
    if (!selectedDisk) return;

    setIsRunningTest(true);
    try {
      const response = await syslibApi.smart.runTest(selectedDisk.path, testType);
      if (response.success) {
        alert(`${testType} test started for ${selectedDisk.name}`);
      }
    } catch (error) {
      console.error('Failed to run SMART test:', error);
      alert('Failed to start SMART test');
    } finally {
      setIsRunningTest(false);
    }
  };

  const getHealthColor = (score: number) => {
    if (score >= 80) return 'text-green-500';
    if (score >= 60) return 'text-yellow-500';
    return 'text-red-500';
  };

  const getHealthIcon = (score: number) => {
    if (score >= 80) return <CheckCircle className="w-6 h-6 text-green-500" />;
    if (score >= 60) return <AlertTriangle className="w-6 h-6 text-yellow-500" />;
    return <XCircle className="w-6 h-6 text-red-500" />;
  };

  const formatHours = (hours: number): string => {
    const days = Math.floor(hours / 24);
    if (days > 365) {
      const years = Math.floor(days / 365);
      return `${years}y ${days % 365}d`;
    }
    return `${days}d ${hours % 24}h`;
  };

  return (
    <div className="flex flex-col h-full bg-gradient-to-br from-gray-50 to-white dark:from-macos-dark-100 dark:to-macos-dark-200">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200/50 dark:border-gray-700/50 bg-white/50 dark:bg-macos-dark-100/50 backdrop-blur-sm">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <Activity className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              SMART Monitoring
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Real-time disk health monitoring
            </p>
          </div>
        </div>
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={fetchDisks}
          className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
        >
          <RefreshCw className="w-4 h-4" />
          Refresh
        </motion.button>
      </div>

      {/* Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Disks List */}
        <div className="w-1/3 border-r border-gray-200 dark:border-gray-700 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center p-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
            </div>
          ) : (
            <div className="p-4 space-y-2">
              {disks.map((disk) => (
                <motion.div
                  key={disk.path}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  whileHover={{ y: -2, scale: 1.02 }}
                  onClick={() => setSelectedDisk(disk)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedDisk?.path === disk.path
                      ? 'bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 border-2 border-macos-blue shadow-lg'
                      : 'bg-white dark:bg-macos-dark-200 hover:bg-gray-50 dark:hover:bg-macos-dark-300 border-2 border-gray-200 dark:border-gray-700 hover:shadow-md'
                  }`}
                >
                  <div className="flex items-center gap-3 mb-3">
                    <div className="p-2 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-300 dark:to-macos-dark-400 rounded-lg">
                      <HardDrive className="w-5 h-5 text-macos-blue" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="font-semibold text-gray-900 dark:text-gray-100 truncate">
                        {disk.name}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400 truncate">
                        {disk.model}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 px-2 py-1 bg-gray-50 dark:bg-macos-dark-300 rounded-lg">
                    <Thermometer className={`w-4 h-4 ${
                      disk.temperature > 50 ? 'text-red-500' :
                      disk.temperature > 45 ? 'text-orange-500' :
                      'text-green-500'
                    }`} />
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {disk.temperature}°C
                    </span>
                  </div>
                </motion.div>
              ))}
            </div>
          )}
        </div>

        {/* SMART Details */}
        <div className="flex-1 overflow-y-auto">
          {selectedDisk ? (
            <div className="p-6">
              {isLoadingSMART ? (
                <div className="flex items-center justify-center p-12">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
                </div>
              ) : smartInfo ? (
                <>
                  {/* Health Score Card */}
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-gradient-to-br from-macos-blue/10 via-macos-purple/10 to-pink-500/10 dark:from-macos-blue/20 dark:via-macos-purple/20 dark:to-pink-500/20 rounded-2xl p-6 mb-6 border border-gray-200/50 dark:border-gray-700/50 shadow-xl"
                  >
                    <div className="flex items-center justify-between mb-4">
                      <div className="flex-1">
                        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                          {selectedDisk.name}
                        </h2>
                        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                          {selectedDisk.model} • {selectedDisk.serial}
                        </p>
                      </div>
                      <div className="flex items-center gap-4">
                        {getHealthIcon(smartInfo.healthScore)}
                        <div className="text-right">
                          <div className={`text-4xl font-bold ${getHealthColor(smartInfo.healthScore)}`}>
                            {smartInfo.healthScore}
                          </div>
                          <div className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">
                            Health Score
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Health Progress Bar */}
                    <div className="h-4 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden shadow-inner">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${smartInfo.healthScore}%` }}
                        transition={{ duration: 1, ease: 'easeOut' }}
                        className={`h-full rounded-full shadow-sm ${
                          smartInfo.healthScore >= 80 ? 'bg-gradient-to-r from-green-500 to-emerald-500' :
                          smartInfo.healthScore >= 60 ? 'bg-gradient-to-r from-yellow-500 to-orange-500' :
                          'bg-gradient-to-r from-red-500 to-rose-600'
                        }`}
                      />
                    </div>

                    <div className="mt-4 flex items-center gap-2 px-3 py-2 bg-white/50 dark:bg-macos-dark-200/50 rounded-lg">
                      {smartInfo.smartStatus === 'PASSED' ? (
                        <>
                          <CheckCircle className="w-4 h-4 text-green-500" />
                          <span className="text-sm font-medium text-green-600 dark:text-green-400">
                            SMART Status: PASSED
                          </span>
                        </>
                      ) : (
                        <>
                          <XCircle className="w-4 h-4 text-red-500" />
                          <span className="text-sm font-medium text-red-600 dark:text-red-400">
                            SMART Status: {smartInfo.smartStatus}
                          </span>
                        </>
                      )}
                    </div>
                  </motion.div>

                  {/* Key Metrics */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                    <motion.div
                      initial={{ opacity: 0, scale: 0.9 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.1 }}
                      className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                    >
                      <div className="flex items-center gap-2 mb-2">
                        <Thermometer className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">Temperature</span>
                      </div>
                      <div className={`text-2xl font-bold ${
                        smartInfo.temperature > 50 ? 'text-red-500' :
                        smartInfo.temperature > 45 ? 'text-yellow-500' :
                        'text-green-500'
                      }`}>
                        {smartInfo.temperature}°C
                      </div>
                    </motion.div>

                    <motion.div
                      initial={{ opacity: 0, scale: 0.9 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.2 }}
                      className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                    >
                      <div className="flex items-center gap-2 mb-2">
                        <Clock className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">Power On</span>
                      </div>
                      <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                        {formatHours(smartInfo.powerOnHours)}
                      </div>
                    </motion.div>

                    <motion.div
                      initial={{ opacity: 0, scale: 0.9 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.3 }}
                      className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                    >
                      <div className="flex items-center gap-2 mb-2">
                        <Zap className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">Bad Sectors</span>
                      </div>
                      <div className={`text-2xl font-bold ${
                        smartInfo.reallocatedSectors > 0 ? 'text-red-500' : 'text-green-500'
                      }`}>
                        {smartInfo.reallocatedSectors}
                      </div>
                    </motion.div>

                    <motion.div
                      initial={{ opacity: 0, scale: 0.9 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.4 }}
                      className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                    >
                      <div className="flex items-center gap-2 mb-2">
                        <AlertTriangle className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">Errors</span>
                      </div>
                      <div className={`text-2xl font-bold ${
                        smartInfo.uncorrectableErrors > 0 ? 'text-red-500' : 'text-green-500'
                      }`}>
                        {smartInfo.uncorrectableErrors}
                      </div>
                    </motion.div>
                  </div>

                  {/* Detailed Attributes */}
                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-xl p-6 mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                      Detailed Attributes
                    </h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600 dark:text-gray-400">Reallocated Sectors</span>
                        <span className={`text-sm font-medium ${
                          smartInfo.reallocatedSectors > 0 ? 'text-red-500' : 'text-green-500'
                        }`}>
                          {smartInfo.reallocatedSectors}
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600 dark:text-gray-400">Pending Sectors</span>
                        <span className={`text-sm font-medium ${
                          smartInfo.pendingSectors > 0 ? 'text-yellow-500' : 'text-green-500'
                        }`}>
                          {smartInfo.pendingSectors}
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600 dark:text-gray-400">Uncorrectable Errors</span>
                        <span className={`text-sm font-medium ${
                          smartInfo.uncorrectableErrors > 0 ? 'text-red-500' : 'text-green-500'
                        }`}>
                          {smartInfo.uncorrectableErrors}
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600 dark:text-gray-400">Power On Hours</span>
                        <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                          {smartInfo.powerOnHours.toLocaleString()}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Self-Test Controls */}
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                      Self-Test Operations
                    </h3>
                    <div className="grid grid-cols-3 gap-3">
                      <motion.button
                        whileHover={{ scale: 1.05, y: -2 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={() => handleRunTest('short')}
                        disabled={isRunningTest}
                        className="flex flex-col items-center gap-2 p-4 bg-gradient-to-br from-macos-blue/10 to-macos-blue/20 dark:from-macos-blue/20 dark:to-macos-blue/30 text-macos-blue rounded-xl hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed border border-macos-blue/20"
                      >
                        <PlayCircle className="w-6 h-6" />
                        <span className="text-sm font-bold">Short Test</span>
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400">~2 min</span>
                      </motion.button>
                      <motion.button
                        whileHover={{ scale: 1.05, y: -2 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={() => handleRunTest('long')}
                        disabled={isRunningTest}
                        className="flex flex-col items-center gap-2 p-4 bg-gradient-to-br from-macos-purple/10 to-macos-purple/20 dark:from-macos-purple/20 dark:to-macos-purple/30 text-macos-purple rounded-xl hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed border border-macos-purple/20"
                      >
                        <PlayCircle className="w-6 h-6" />
                        <span className="text-sm font-bold">Long Test</span>
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400">~hours</span>
                      </motion.button>
                      <motion.button
                        whileHover={{ scale: 1.05, y: -2 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={() => handleRunTest('conveyance')}
                        disabled={isRunningTest}
                        className="flex flex-col items-center gap-2 p-4 bg-gradient-to-br from-green-500/10 to-emerald-500/20 dark:from-green-500/20 dark:to-emerald-500/30 text-green-600 dark:text-green-400 rounded-xl hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed border border-green-500/20"
                      >
                        <PlayCircle className="w-6 h-6" />
                        <span className="text-sm font-bold">Conveyance</span>
                        <span className="text-xs font-medium text-gray-600 dark:text-gray-400">~5 min</span>
                      </motion.button>
                    </div>
                  </div>
                </>
              ) : (
                <div className="flex flex-col items-center justify-center p-12 text-center">
                  <AlertTriangle className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" />
                  <p className="text-gray-500 dark:text-gray-400">
                    SMART data not available for this disk
                  </p>
                </div>
              )}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center h-full text-center p-12">
              <HardDrive className="w-24 h-24 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-lg text-gray-500 dark:text-gray-400">
                Select a disk to view SMART data
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
