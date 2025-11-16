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

  const getHealthBgColor = (score: number) => {
    if (score >= 80) return 'bg-green-500';
    if (score >= 60) return 'bg-yellow-500';
    return 'bg-red-500';
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
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-3">
          <Activity className="w-6 h-6 text-macos-blue" />
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
            SMART Monitoring
          </h1>
        </div>
        <button
          onClick={fetchDisks}
          className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
        >
          <RefreshCw className="w-4 h-4" />
          Refresh
        </button>
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
                  onClick={() => setSelectedDisk(disk)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedDisk?.path === disk.path
                      ? 'bg-macos-blue/10 dark:bg-macos-blue/20 border-2 border-macos-blue'
                      : 'bg-gray-50 dark:bg-macos-dark-200 hover:bg-gray-100 dark:hover:bg-macos-dark-300 border-2 border-transparent'
                  }`}
                >
                  <div className="flex items-center gap-3 mb-2">
                    <HardDrive className="w-5 h-5 text-macos-blue" />
                    <div className="flex-1">
                      <div className="font-semibold text-gray-900 dark:text-gray-100">
                        {disk.name}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {disk.model}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Thermometer className="w-4 h-4 text-gray-400" />
                    <span className="text-sm text-gray-600 dark:text-gray-400">
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
                  <div className="bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 rounded-2xl p-6 mb-6">
                    <div className="flex items-center justify-between mb-4">
                      <div>
                        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                          {selectedDisk.name}
                        </h2>
                        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                          {selectedDisk.model} • {selectedDisk.serial}
                        </p>
                      </div>
                      <div className="flex items-center gap-3">
                        {getHealthIcon(smartInfo.healthScore)}
                        <div className="text-right">
                          <div className={`text-3xl font-bold ${getHealthColor(smartInfo.healthScore)}`}>
                            {smartInfo.healthScore}
                          </div>
                          <div className="text-xs text-gray-600 dark:text-gray-400">
                            Health Score
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Health Progress Bar */}
                    <div className="h-3 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${smartInfo.healthScore}%` }}
                        className={`h-full rounded-full ${getHealthBgColor(smartInfo.healthScore)}`}
                      />
                    </div>

                    <div className="mt-4 flex items-center gap-2">
                      {smartInfo.smartStatus === 'PASSED' ? (
                        <>
                          <CheckCircle className="w-4 h-4 text-green-500" />
                          <span className="text-sm text-green-600 dark:text-green-400">
                            SMART Status: PASSED
                          </span>
                        </>
                      ) : (
                        <>
                          <XCircle className="w-4 h-4 text-red-500" />
                          <span className="text-sm text-red-600 dark:text-red-400">
                            SMART Status: {smartInfo.smartStatus}
                          </span>
                        </>
                      )}
                    </div>
                  </div>

                  {/* Key Metrics */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
                      <div className="flex items-center gap-2 mb-2">
                        <Thermometer className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs text-gray-600 dark:text-gray-400">Temperature</span>
                      </div>
                      <div className={`text-2xl font-bold ${
                        smartInfo.temperature > 50 ? 'text-red-500' :
                        smartInfo.temperature > 45 ? 'text-yellow-500' :
                        'text-gray-900 dark:text-gray-100'
                      }`}>
                        {smartInfo.temperature}°C
                      </div>
                    </div>

                    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
                      <div className="flex items-center gap-2 mb-2">
                        <Clock className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs text-gray-600 dark:text-gray-400">Power On</span>
                      </div>
                      <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                        {formatHours(smartInfo.powerOnHours)}
                      </div>
                    </div>

                    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
                      <div className="flex items-center gap-2 mb-2">
                        <Zap className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs text-gray-600 dark:text-gray-400">Bad Sectors</span>
                      </div>
                      <div className={`text-2xl font-bold ${
                        smartInfo.reallocatedSectors > 0 ? 'text-red-500' : 'text-green-500'
                      }`}>
                        {smartInfo.reallocatedSectors}
                      </div>
                    </div>

                    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
                      <div className="flex items-center gap-2 mb-2">
                        <AlertTriangle className="w-4 h-4 text-macos-blue" />
                        <span className="text-xs text-gray-600 dark:text-gray-400">Errors</span>
                      </div>
                      <div className={`text-2xl font-bold ${
                        smartInfo.uncorrectableErrors > 0 ? 'text-red-500' : 'text-green-500'
                      }`}>
                        {smartInfo.uncorrectableErrors}
                      </div>
                    </div>
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
                      <button
                        onClick={() => handleRunTest('short')}
                        disabled={isRunningTest}
                        className="flex flex-col items-center gap-2 p-4 bg-macos-blue/10 dark:bg-macos-blue/20 text-macos-blue rounded-xl hover:bg-macos-blue/20 dark:hover:bg-macos-blue/30 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <PlayCircle className="w-6 h-6" />
                        <span className="text-sm font-medium">Short Test</span>
                        <span className="text-xs text-gray-600 dark:text-gray-400">~2 min</span>
                      </button>
                      <button
                        onClick={() => handleRunTest('long')}
                        disabled={isRunningTest}
                        className="flex flex-col items-center gap-2 p-4 bg-macos-purple/10 dark:bg-macos-purple/20 text-macos-purple rounded-xl hover:bg-macos-purple/20 dark:hover:bg-macos-purple/30 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <PlayCircle className="w-6 h-6" />
                        <span className="text-sm font-medium">Long Test</span>
                        <span className="text-xs text-gray-600 dark:text-gray-400">~hours</span>
                      </button>
                      <button
                        onClick={() => handleRunTest('conveyance')}
                        disabled={isRunningTest}
                        className="flex flex-col items-center gap-2 p-4 bg-macos-green/10 dark:bg-macos-green/20 text-macos-green rounded-xl hover:bg-macos-green/20 dark:hover:bg-macos-green/30 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <PlayCircle className="w-6 h-6" />
                        <span className="text-sm font-medium">Conveyance</span>
                        <span className="text-xs text-gray-600 dark:text-gray-400">~5 min</span>
                      </button>
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
