import { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  X,
  Server,
  AlertCircle,
  Network,
  Cpu,
  HardDrive,
  ChevronRight,
  ChevronLeft,
  Check,
  Sparkles,
  Disc,
  Eye,
  EyeOff,
  Lock,
  Key,
} from 'lucide-react';
import { vmsApi, type VMCreateRequest } from '@/api/vms';
import { networkApi } from '@/api/network';
import { getErrorMessage } from '@/api/client';

interface CreateVMWizardProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

type WizardStep = 'general' | 'resources' | 'network' | 'review';

interface StepConfig {
  id: WizardStep;
  title: string;
  description: string;
  icon: typeof Server;
}

const STEPS: StepConfig[] = [
  { id: 'general', title: 'General', description: 'Basic configuration', icon: Server },
  { id: 'resources', title: 'Resources', description: 'CPU, RAM & Storage', icon: Cpu },
  { id: 'network', title: 'Network', description: 'Network settings', icon: Network },
  { id: 'review', title: 'Review', description: 'Confirm & create', icon: Check },
];

export function CreateVMWizard({ isOpen, onClose, onSuccess }: CreateVMWizardProps) {
  const [currentStep, setCurrentStep] = useState<WizardStep>('general');
  const [formData, setFormData] = useState<VMCreateRequest & { password?: string; password_confirm?: string; ssh_key?: string }>({
    name: '',
    memory: 2048,
    vcpus: 2,
    disk_size: 20,
    os_type: 'linux',
    os_variant: 'ubuntu22.04',
    iso_path: '',
    network: 'default',
    autostart: false,
    password: '',
    password_confirm: '',
    ssh_key: '',
  });
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string>('');
  const [bridges, setBridges] = useState<string[]>(['default', 'br0']);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  // Fetch available bridges
  useEffect(() => {
    const fetchBridges = async () => {
      try {
        const response = await networkApi.listBridges();
        if (response.success && response.data && response.data.length > 0) {
          const availableBridges = ['default', ...response.data];
          setBridges(availableBridges);
        }
      } catch (err) {
        console.error('Failed to fetch bridges:', err);
      }
    };

    if (isOpen) {
      fetchBridges();
    }
  }, [isOpen]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isOpen) return;

      if (e.key === 'Escape') {
        e.preventDefault();
        if (currentStep === 'general') {
          onClose();
        } else {
          handleBack();
        }
      } else if (e.key === 'Enter' && !e.shiftKey && !creating) {
        e.preventDefault();
        handleNext();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, currentStep, creating]);

  const validateStep = useCallback((step: WizardStep): boolean => {
    const errors: Record<string, string> = {};

    if (step === 'general') {
      if (!formData.name.trim()) {
        errors.name = 'VM name is required';
      } else if (!/^[a-zA-Z0-9][a-zA-Z0-9_-]*$/.test(formData.name)) {
        errors.name = 'Invalid name format (alphanumeric, dash, underscore)';
      }
      if (!formData.os_type) {
        errors.os_type = 'OS type is required';
      }
      if (formData.password && formData.password.length < 8) {
        errors.password = 'Password must be at least 8 characters';
      }
      if (formData.password && formData.password !== formData.password_confirm) {
        errors.password_confirm = 'Passwords do not match';
      }
      if (formData.ssh_key && !formData.ssh_key.trim().startsWith('ssh-')) {
        errors.ssh_key = 'Invalid SSH key format (should start with ssh-rsa, ssh-ed25519, etc.)';
      }
    } else if (step === 'resources') {
      if (formData.memory < 512 || formData.memory > 32768) {
        errors.memory = 'Memory must be between 512 MB and 32 GB';
      }
      if (formData.vcpus < 1 || formData.vcpus > 16) {
        errors.vcpus = 'CPU cores must be between 1 and 16';
      }
      if (formData.disk_size < 10 || formData.disk_size > 500) {
        errors.disk_size = 'Disk size must be between 10 GB and 500 GB';
      }
    } else if (step === 'network') {
      if (!formData.network) {
        errors.network = 'Network bridge is required';
      }
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  }, [formData]);

  const handleNext = () => {
    if (!validateStep(currentStep)) {
      return;
    }

    const stepIndex = STEPS.findIndex((s) => s.id === currentStep);
    if (stepIndex < STEPS.length - 1) {
      setCurrentStep(STEPS[stepIndex + 1].id);
      setError('');
    } else {
      handleCreate();
    }
  };

  const handleBack = () => {
    const stepIndex = STEPS.findIndex((s) => s.id === currentStep);
    if (stepIndex > 0) {
      setCurrentStep(STEPS[stepIndex - 1].id);
      setError('');
    }
  };

  const handleCreate = async () => {
    if (!validateStep('review')) {
      return;
    }

    try {
      setCreating(true);
      setError('');
      const response = await vmsApi.createVM(formData);

      if (response.success) {
        // Success animation
        await new Promise((resolve) => setTimeout(resolve, 1000));
        onSuccess();
        onClose();
        // Reset form
        setFormData({
          name: '',
          memory: 2048,
          vcpus: 2,
          disk_size: 20,
          os_type: 'linux',
          os_variant: 'ubuntu22.04',
          iso_path: '',
          network: 'default',
          autostart: false,
          password: '',
          password_confirm: '',
          ssh_key: '',
        });
        setCurrentStep('general');
      } else {
        setError(response.error?.message || 'Failed to create VM');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setCreating(false);
    }
  };

  const getCurrentStepIndex = () => STEPS.findIndex((s) => s.id === currentStep);
  const canGoNext = validateStep(currentStep);

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      >
        <motion.div
          initial={{ scale: 0.95, opacity: 0, y: 20 }}
          animate={{ scale: 1, opacity: 1, y: 0 }}
          exit={{ scale: 0.95, opacity: 0, y: 20 }}
          transition={{ type: 'spring', duration: 0.5 }}
          className="bg-white/90 dark:bg-macos-dark-100/90 backdrop-blur-xl rounded-2xl shadow-2xl max-w-3xl w-full mx-4 max-h-[90vh] overflow-hidden border border-gray-200/50 dark:border-gray-700/50"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="relative p-6 border-b border-gray-200/50 dark:border-gray-700/50">
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-xl bg-macos-blue/10 dark:bg-macos-blue/20">
                  <Server className="w-6 h-6 text-macos-blue" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                    Create Virtual Machine
                  </h2>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Step {getCurrentStepIndex() + 1} of {STEPS.length}
                  </p>
                </div>
              </div>
              <button
                onClick={onClose}
                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-all hover:scale-110"
              >
                <X className="w-5 h-5 text-gray-600 dark:text-gray-400" />
              </button>
            </div>

            {/* Progress Steps */}
            <div className="flex items-center justify-between">
              {STEPS.map((step, index) => {
                const Icon = step.icon;
                const isActive = step.id === currentStep;
                const isCompleted = index < getCurrentStepIndex();

                return (
                  <div key={step.id} className="flex items-center flex-1">
                    <div className="flex flex-col items-center flex-1">
                      <motion.div
                        className={`
                          w-10 h-10 rounded-full flex items-center justify-center mb-2
                          transition-all duration-300
                          ${
                            isCompleted
                              ? 'bg-green-500 text-white'
                              : isActive
                              ? 'bg-macos-blue text-white shadow-lg shadow-macos-blue/30'
                              : 'bg-gray-200 dark:bg-gray-700 text-gray-400 dark:text-gray-500'
                          }
                        `}
                        animate={isActive ? { scale: [1, 1.1, 1] } : {}}
                        transition={{ duration: 0.3 }}
                      >
                        {isCompleted ? (
                          <Check className="w-5 h-5" />
                        ) : (
                          <Icon className="w-5 h-5" />
                        )}
                      </motion.div>
                      <div className="text-center">
                        <div
                          className={`text-xs font-medium ${
                            isActive
                              ? 'text-gray-900 dark:text-white'
                              : 'text-gray-500 dark:text-gray-400'
                          }`}
                        >
                          {step.title}
                        </div>
                      </div>
                    </div>
                    {index < STEPS.length - 1 && (
                      <div
                        className={`
                          h-0.5 flex-1 mx-2 transition-all duration-300
                          ${
                            isCompleted
                              ? 'bg-green-500'
                              : 'bg-gray-200 dark:bg-gray-700'
                          }
                        `}
                      />
                    )}
                  </div>
                );
              })}
            </div>
          </div>

          {/* Error Display */}
          <AnimatePresence>
            {error && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: 'auto' }}
                exit={{ opacity: 0, height: 0 }}
                className="mx-6 mt-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-start gap-3"
              >
                <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
                <div className="flex-1">
                  <h3 className="font-semibold text-red-900 dark:text-red-200">Error</h3>
                  <p className="text-sm text-red-700 dark:text-red-300 mt-1">{error}</p>
                </div>
                <button
                  onClick={() => setError('')}
                  className="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200"
                >
                  <X className="w-4 h-4" />
                </button>
              </motion.div>
            )}
          </AnimatePresence>

          {/* Step Content */}
          <div className="p-6 min-h-[400px]">
            <AnimatePresence mode="wait">
              <motion.div
                key={currentStep}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3 }}
              >
                {currentStep === 'general' && (
                  <GeneralStep
                    formData={formData}
                    setFormData={setFormData}
                    errors={validationErrors}
                  />
                )}
                {currentStep === 'resources' && (
                  <ResourcesStep
                    formData={formData}
                    setFormData={setFormData}
                    errors={validationErrors}
                  />
                )}
                {currentStep === 'network' && (
                  <NetworkStep
                    formData={formData}
                    setFormData={setFormData}
                    bridges={bridges}
                    errors={validationErrors}
                  />
                )}
                {currentStep === 'review' && (
                  <ReviewStep formData={formData} />
                )}
              </motion.div>
            </AnimatePresence>
          </div>

          {/* Footer with Navigation */}
          <div className="p-6 border-t border-gray-200/50 dark:border-gray-700/50 bg-gray-50/50 dark:bg-macos-dark-50/50 backdrop-blur-xl">
            <div className="flex items-center justify-between">
              <button
                onClick={getCurrentStepIndex() === 0 ? onClose : handleBack}
                disabled={creating}
                className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-all flex items-center gap-2 disabled:opacity-50"
              >
                <ChevronLeft className="w-4 h-4" />
                {getCurrentStepIndex() === 0 ? 'Cancel' : 'Back'}
              </button>

              <div className="text-sm text-gray-500 dark:text-gray-400">
                Press <kbd className="px-2 py-1 bg-gray-200 dark:bg-gray-700 rounded text-xs">Enter</kbd> to continue, <kbd className="px-2 py-1 bg-gray-200 dark:bg-gray-700 rounded text-xs">Esc</kbd> to go back
              </div>

              <button
                onClick={handleNext}
                disabled={!canGoNext || creating}
                className="px-6 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-macos-blue/30 hover:shadow-xl hover:shadow-macos-blue/40 hover:scale-105"
              >
                {creating ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                    Creating...
                  </>
                ) : currentStep === 'review' ? (
                  <>
                    <Sparkles className="w-4 h-4" />
                    Create VM
                  </>
                ) : (
                  <>
                    Next
                    <ChevronRight className="w-4 h-4" />
                  </>
                )}
              </button>
            </div>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
}

// Step Components
interface StepProps {
  formData: VMCreateRequest & { password?: string; password_confirm?: string; ssh_key?: string };
  setFormData: (data: any) => void;
  errors: Record<string, string>;
}

function GeneralStep({ formData, setFormData, errors }: StepProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [showPasswordConfirm, setShowPasswordConfirm] = useState(false);
  const osTypes = [
    { value: 'linux', label: 'Linux', description: 'Linux distributions' },
    { value: 'windows', label: 'Windows', description: 'Windows operating systems' },
    { value: 'unix', label: 'Unix', description: 'Unix-like systems' },
    { value: 'other', label: 'Other', description: 'Other operating systems' },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          General Configuration
        </h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Configure basic VM settings and choose your operating system
        </p>
      </div>

      {/* VM Name */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          VM Name <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="my-vm"
          className={`
            w-full px-4 py-3 bg-white dark:bg-macos-dark-50 border rounded-lg
            focus:ring-2 focus:ring-macos-blue focus:border-transparent
            text-gray-900 dark:text-white transition-all
            ${errors.name ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'}
          `}
          autoFocus
        />
        {errors.name && (
          <p className="text-sm text-red-600 dark:text-red-400 mt-1">{errors.name}</p>
        )}
        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Only alphanumeric characters, dashes, and underscores
        </p>
      </div>

      {/* OS Type Selection */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
          Operating System Type <span className="text-red-500">*</span>
        </label>
        <div className="grid grid-cols-2 gap-3">
          {osTypes.map((osType) => (
            <label
              key={osType.value}
              className={`
                p-4 border rounded-lg cursor-pointer transition-all
                ${
                  formData.os_type === osType.value
                    ? 'border-macos-blue bg-macos-blue/5 dark:bg-macos-blue/10'
                    : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
                }
              `}
            >
              <input
                type="radio"
                name="os_type"
                value={osType.value}
                checked={formData.os_type === osType.value}
                onChange={(e) => setFormData({ ...formData, os_type: e.target.value })}
                className="sr-only"
              />
              <div className="font-medium text-gray-900 dark:text-white">{osType.label}</div>
              <div className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                {osType.description}
              </div>
            </label>
          ))}
        </div>
      </div>

      {/* OS Variant */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          OS Variant (optional)
        </label>
        <input
          type="text"
          value={formData.os_variant}
          onChange={(e) => setFormData({ ...formData, os_variant: e.target.value })}
          placeholder="ubuntu22.04, win10, debian11"
          className="w-full px-4 py-3 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
        />
        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Helps optimize VM settings for specific OS versions
        </p>
      </div>

      {/* ISO Path */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          ISO Path (optional)
        </label>
        <div className="relative">
          <Disc className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
          <input
            type="text"
            value={formData.iso_path}
            onChange={(e) => setFormData({ ...formData, iso_path: e.target.value })}
            placeholder="/path/to/iso/file.iso"
            className="w-full pl-10 pr-4 py-3 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
          />
        </div>
        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Path to installation ISO on the host system
        </p>
      </div>

      {/* Password Section */}
      <div className="p-6 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/10 dark:to-indigo-900/10 rounded-xl border border-blue-200 dark:border-blue-800">
        <div className="flex items-center gap-2 mb-4">
          <Lock className="w-5 h-5 text-macos-blue" />
          <h4 className="font-semibold text-gray-900 dark:text-white">Root Password (optional)</h4>
        </div>
        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
          Set a root password for SSH access. Highly recommended for easy management!
        </p>

        <div className="grid grid-cols-2 gap-4">
          {/* Password */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Password
            </label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                value={formData.password || ''}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                placeholder="Enter password"
                className={`
                  w-full px-4 py-3 pr-12 bg-white dark:bg-macos-dark-50 border rounded-lg
                  focus:ring-2 focus:ring-macos-blue focus:border-transparent
                  text-gray-900 dark:text-white transition-all
                  ${errors.password ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'}
                `}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
            {errors.password && (
              <p className="text-sm text-red-600 dark:text-red-400 mt-1">{errors.password}</p>
            )}
          </div>

          {/* Password Confirm */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Confirm Password
            </label>
            <div className="relative">
              <input
                type={showPasswordConfirm ? 'text' : 'password'}
                value={formData.password_confirm || ''}
                onChange={(e) => setFormData({ ...formData, password_confirm: e.target.value })}
                placeholder="Confirm password"
                className={`
                  w-full px-4 py-3 pr-12 bg-white dark:bg-macos-dark-50 border rounded-lg
                  focus:ring-2 focus:ring-macos-blue focus:border-transparent
                  text-gray-900 dark:text-white transition-all
                  ${errors.password_confirm ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'}
                `}
              />
              <button
                type="button"
                onClick={() => setShowPasswordConfirm(!showPasswordConfirm)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                {showPasswordConfirm ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
            {errors.password_confirm && (
              <p className="text-sm text-red-600 dark:text-red-400 mt-1">{errors.password_confirm}</p>
            )}
          </div>
        </div>

        <div className="mt-3 flex items-start gap-2 text-xs text-gray-600 dark:text-gray-400">
          <div className="mt-0.5">ℹ️</div>
          <div>
            Minimum 8 characters. Leave empty to create VM without password (you can set it later).
          </div>
        </div>
      </div>

      {/* SSH Key Section */}
      <div className="p-6 bg-gradient-to-r from-green-50 to-emerald-50 dark:from-green-900/10 dark:to-emerald-900/10 rounded-xl border border-green-200 dark:border-green-800">
        <div className="flex items-center gap-2 mb-4">
          <Key className="w-5 h-5 text-green-600 dark:text-green-400" />
          <h4 className="font-semibold text-gray-900 dark:text-white">SSH Public Key (optional)</h4>
        </div>
        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
          Add your SSH public key for passwordless authentication. More secure than password-only!
        </p>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Public Key
          </label>
          <textarea
            value={formData.ssh_key || ''}
            onChange={(e) => setFormData({ ...formData, ssh_key: e.target.value })}
            placeholder="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB... user@host"
            rows={3}
            className={`
              w-full px-4 py-3 bg-white dark:bg-macos-dark-50 border rounded-lg
              focus:ring-2 focus:ring-green-500 focus:border-transparent
              text-gray-900 dark:text-white transition-all font-mono text-sm
              ${errors.ssh_key ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'}
            `}
          />
          {errors.ssh_key && (
            <p className="text-sm text-red-600 dark:text-red-400 mt-1">{errors.ssh_key}</p>
          )}
        </div>

        <div className="mt-3 space-y-2 text-xs text-gray-600 dark:text-gray-400">
          <div className="flex items-start gap-2">
            <div className="mt-0.5">✓</div>
            <div>Paste your public key from ~/.ssh/id_rsa.pub or ~/.ssh/id_ed25519.pub</div>
          </div>
          <div className="flex items-start gap-2">
            <div className="mt-0.5">✓</div>
            <div>You can use both password and SSH key for maximum flexibility</div>
          </div>
        </div>
      </div>
    </div>
  );
}

function ResourcesStep({ formData, setFormData }: StepProps) {
  const getMemoryRecommendation = (mb: number) => {
    if (mb < 1024) return '⚠️ Very low - may cause issues';
    if (mb < 2048) return '⚠️ Low - minimal for modern OS';
    if (mb < 4096) return '✓ Good - for typical workloads';
    if (mb < 8192) return '✓✓ Great - for demanding apps';
    return '✓✓ Excellent - for heavy workloads';
  };

  const getCpuRecommendation = (cores: number) => {
    if (cores === 1) return '✓ Single core - basic tasks';
    if (cores === 2) return '✓ Dual core - standard workloads';
    if (cores <= 4) return '✓✓ Quad core - good performance';
    return '✓✓ Multi-core - high performance';
  };

  const getDiskRecommendation = (gb: number) => {
    if (gb < 20) return '⚠️ Very small - may not be enough';
    if (gb < 40) return '✓ Minimal - for Linux VMs';
    if (gb < 80) return '✓ Good - for Windows or data';
    if (gb < 160) return '✓✓ Great - plenty of space';
    return '✓✓ Excellent - very spacious';
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Resource Allocation
        </h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Allocate CPU, memory, and storage resources for your VM
        </p>
      </div>

      {/* Memory Slider */}
      <div className="p-6 bg-gray-50 dark:bg-macos-dark-50 rounded-xl">
        <div className="flex items-center justify-between mb-4">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
            <HardDrive className="w-4 h-4 text-macos-blue" />
            Memory
          </label>
          <div className="flex items-baseline gap-2">
            <span className="text-2xl font-bold text-macos-blue">{formData.memory}</span>
            <span className="text-sm text-gray-600 dark:text-gray-400">MB</span>
            <span className="text-sm text-gray-500 dark:text-gray-400">
              ({(formData.memory / 1024).toFixed(1)} GB)
            </span>
          </div>
        </div>

        <input
          type="range"
          min="512"
          max="32768"
          step="512"
          value={formData.memory}
          onChange={(e) => setFormData({ ...formData, memory: parseInt(e.target.value) })}
          className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-macos-blue"
        />

        <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-2">
          <span>512 MB</span>
          <span>16 GB</span>
          <span>32 GB</span>
        </div>

        <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <p className="text-sm text-blue-900 dark:text-blue-200">
            {getMemoryRecommendation(formData.memory)}
          </p>
        </div>
      </div>

      {/* CPU Slider */}
      <div className="p-6 bg-gray-50 dark:bg-macos-dark-50 rounded-xl">
        <div className="flex items-center justify-between mb-4">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
            <Cpu className="w-4 h-4 text-macos-blue" />
            Virtual CPUs
          </label>
          <div className="flex items-baseline gap-2">
            <span className="text-2xl font-bold text-macos-blue">{formData.vcpus}</span>
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {formData.vcpus === 1 ? 'Core' : 'Cores'}
            </span>
          </div>
        </div>

        <input
          type="range"
          min="1"
          max="16"
          step="1"
          value={formData.vcpus}
          onChange={(e) => setFormData({ ...formData, vcpus: parseInt(e.target.value) })}
          className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-macos-blue"
        />

        <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-2">
          <span>1</span>
          <span>8</span>
          <span>16</span>
        </div>

        <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <p className="text-sm text-blue-900 dark:text-blue-200">
            {getCpuRecommendation(formData.vcpus)}
          </p>
        </div>
      </div>

      {/* Disk Size Slider */}
      <div className="p-6 bg-gray-50 dark:bg-macos-dark-50 rounded-xl">
        <div className="flex items-center justify-between mb-4">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
            <Disc className="w-4 h-4 text-macos-blue" />
            Disk Size
          </label>
          <div className="flex items-baseline gap-2">
            <span className="text-2xl font-bold text-macos-blue">{formData.disk_size}</span>
            <span className="text-sm text-gray-600 dark:text-gray-400">GB</span>
          </div>
        </div>

        <input
          type="range"
          min="10"
          max="500"
          step="10"
          value={formData.disk_size}
          onChange={(e) => setFormData({ ...formData, disk_size: parseInt(e.target.value) })}
          className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-macos-blue"
        />

        <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-2">
          <span>10 GB</span>
          <span>250 GB</span>
          <span>500 GB</span>
        </div>

        <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <p className="text-sm text-blue-900 dark:text-blue-200">
            {getDiskRecommendation(formData.disk_size)}
          </p>
        </div>
      </div>
    </div>
  );
}

interface NetworkStepProps extends StepProps {
  bridges: string[];
}

function NetworkStep({ formData, setFormData, bridges, errors }: NetworkStepProps) {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Network Configuration
        </h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Choose network bridge and autostart settings
        </p>
      </div>

      {/* Network Bridge */}
      <div className="p-6 bg-gray-50 dark:bg-macos-dark-50 rounded-xl">
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3 flex items-center gap-2">
          <Network className="w-4 h-4 text-macos-blue" />
          Network Bridge <span className="text-red-500">*</span>
        </label>
        <select
          value={formData.network}
          onChange={(e) => setFormData({ ...formData, network: e.target.value })}
          className={`
            w-full px-4 py-3 bg-white dark:bg-macos-dark-50 border rounded-lg
            focus:ring-2 focus:ring-macos-blue focus:border-transparent
            text-gray-900 dark:text-white
            ${errors.network ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'}
          `}
        >
          {bridges.map((bridge) => (
            <option key={bridge} value={bridge}>
              {bridge}
              {bridge === 'default' && ' (NAT - Internal network)'}
            </option>
          ))}
        </select>
        {errors.network && (
          <p className="text-sm text-red-600 dark:text-red-400 mt-1">{errors.network}</p>
        )}
        <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
          Available bridges: {bridges.join(', ')}
        </p>

        {formData.network === 'default' && (
          <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <p className="text-sm text-blue-900 dark:text-blue-200">
              ℹ️ Default network provides NAT with internal IP addressing
            </p>
          </div>
        )}
        {formData.network !== 'default' && (
          <div className="mt-4 p-3 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
            <p className="text-sm text-green-900 dark:text-green-200">
              ✓ Bridged network provides direct LAN access
            </p>
          </div>
        )}
      </div>

      {/* Autostart Option */}
      <div className="p-6 bg-gray-50 dark:bg-macos-dark-50 rounded-xl">
        <label className="flex items-start gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={formData.autostart}
            onChange={(e) => setFormData({ ...formData, autostart: e.target.checked })}
            className="mt-1 w-4 h-4 text-macos-blue rounded"
          />
          <div>
            <div className="font-medium text-gray-900 dark:text-white">
              Start VM automatically on boot
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              VM will start automatically when the host system boots
            </p>
          </div>
        </label>
      </div>
    </div>
  );
}

interface ReviewStepProps {
  formData: VMCreateRequest & { password?: string; ssh_key?: string };
}

function ReviewStep({ formData }: ReviewStepProps) {
  const sections = [
    {
      title: 'General',
      icon: Server,
      items: [
        { label: 'Name', value: formData.name },
        { label: 'OS Type', value: formData.os_type || 'linux' },
        { label: 'OS Variant', value: formData.os_variant || 'default' },
        ...(formData.iso_path ? [{ label: 'ISO Path', value: formData.iso_path }] : []),
      ],
    },
    {
      title: 'Security',
      icon: Lock,
      items: [
        {
          label: 'Root Password',
          value: formData.password ? '••••••••' : 'Not set'
        },
        {
          label: 'SSH Key',
          value: formData.ssh_key
            ? `${formData.ssh_key.substring(0, 30)}...`
            : 'Not configured'
        },
      ],
    },
    {
      title: 'Resources',
      icon: Cpu,
      items: [
        {
          label: 'Memory',
          value: `${formData.memory} MB (${(formData.memory / 1024).toFixed(1)} GB)`,
        },
        { label: 'Virtual CPUs', value: formData.vcpus.toString() },
        { label: 'Disk Size', value: `${formData.disk_size} GB` },
      ],
    },
    {
      title: 'Network',
      icon: Network,
      items: [
        { label: 'Network Bridge', value: formData.network || 'default' },
        { label: 'Autostart', value: formData.autostart ? 'Yes' : 'No' },
      ],
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Review Configuration
        </h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Review your VM configuration before creation
        </p>
      </div>

      <div className="space-y-4">
        {sections.map((section) => {
          const Icon = section.icon;
          return (
            <motion.div
              key={section.title}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="p-6 bg-gray-50 dark:bg-macos-dark-50 rounded-xl border border-gray-200 dark:border-gray-700"
            >
              <div className="flex items-center gap-2 mb-4">
                <Icon className="w-5 h-5 text-macos-blue" />
                <h4 className="font-semibold text-gray-900 dark:text-white">{section.title}</h4>
              </div>
              <div className="grid gap-3">
                {section.items.map((item, index) => (
                  <div
                    key={index}
                    className="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700 last:border-0"
                  >
                    <span className="text-sm text-gray-600 dark:text-gray-400">{item.label}</span>
                    <span className="text-sm font-medium text-gray-900 dark:text-white">
                      {item.value}
                    </span>
                  </div>
                ))}
              </div>
            </motion.div>
          );
        })}
      </div>

      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
        <p className="text-sm text-blue-900 dark:text-blue-200">
          ✓ Ready to create VM "{formData.name}"
        </p>
      </div>
    </div>
  );
}
