import { useState } from 'react';
import { motion } from 'framer-motion';
import { useAuthStore } from '@/store';
import { authApi } from '@/api/auth';
import { twofaApi } from '@/api/twofa';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

interface LoginProps {
  onSuccess: () => void;
}

export default function Login({ onSuccess }: LoginProps) {
  const setAuth = useAuthStore((state) => state.setAuth);

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // 2FA state
  const [requires2FA, setRequires2FA] = useState(false);
  const [userId, setUserId] = useState<number | null>(null);
  const [twoFACode, setTwoFACode] = useState('');
  const [useBackupCode, setUseBackupCode] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const response = await authApi.login({ username, password });

      if (response.success && response.data) {
        // Check if 2FA is required
        if (response.data.requires2FA) {
          setRequires2FA(true);
          setUserId(response.data.userId || null);
          setIsLoading(false);
          return;
        }

        // Normal login (no 2FA)
        const { accessToken, refreshToken, user } = response.data;
        if (accessToken && refreshToken && user) {
          setAuth(user, accessToken, refreshToken);
          onSuccess();
        } else {
          setError('Invalid login response');
        }
      } else {
        setError(response.error?.message || 'Login failed');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  const handle2FASubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const response = await twofaApi.loginWith2FA({
        userId: userId!,
        code: twoFACode,
        isBackupCode: useBackupCode,
      });

      if (response) {
        const { accessToken, refreshToken, user } = response;
        setAuth(user, accessToken, refreshToken);
        onSuccess();
      }
    } catch (err) {
      setError(getErrorMessage(err));
      setTwoFACode(''); // Clear the code on error
    } finally {
      setIsLoading(false);
    }
  };

  const handleBackToLogin = () => {
    setRequires2FA(false);
    setUserId(null);
    setTwoFACode('');
    setUseBackupCode(false);
    setError('');
  };

  return (
    <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 dark:from-macos-dark-50 dark:via-macos-dark-100 dark:to-macos-dark-200">
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
        className="w-full max-w-md"
      >
        <div className="glass-light dark:glass-dark rounded-2xl shadow-macos-xl p-8 border border-gray-200/20 dark:border-gray-700/20">
          {/* Logo/Header */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-macos-blue text-white text-3xl mb-4">
              🏠
            </div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Stumpf.Works NAS
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mt-2">
              Sign in to continue
            </p>
          </div>

          {/* Login Form or 2FA Form */}
          {!requires2FA ? (
            <form onSubmit={handleSubmit} className="space-y-4">
              <Input
                label="Username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter your username"
                required
                autoFocus
              />

              <Input
                label="Password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                required
              />

              {error && (
                <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
                  {error}
                </div>
              )}

              <Button
                type="submit"
                variant="primary"
                size="lg"
                isLoading={isLoading}
                className="w-full"
              >
                Sign In
              </Button>
            </form>
          ) : (
            <form onSubmit={handle2FASubmit} className="space-y-4">
              <div className="text-center mb-4">
                <p className="text-gray-700 dark:text-gray-300 text-sm">
                  Enter the {useBackupCode ? 'backup code' : '6-digit code'} from your
                  authenticator app
                </p>
              </div>

              <Input
                label="Verification Code"
                type="text"
                value={twoFACode}
                onChange={(e) =>
                  setTwoFACode(
                    useBackupCode ? e.target.value : e.target.value.replace(/\D/g, '').slice(0, 6)
                  )
                }
                placeholder={useBackupCode ? 'XXXXXXXX' : '000000'}
                required
                autoFocus
                className="text-center text-2xl tracking-widest font-mono"
              />

              {error && (
                <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
                  {error}
                </div>
              )}

              <div className="flex gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="lg"
                  onClick={handleBackToLogin}
                  className="flex-1"
                >
                  Back
                </Button>
                <Button
                  type="submit"
                  variant="primary"
                  size="lg"
                  isLoading={isLoading}
                  className="flex-1"
                  disabled={
                    useBackupCode ? twoFACode.length === 0 : twoFACode.length !== 6
                  }
                >
                  Verify
                </Button>
              </div>

              <button
                type="button"
                onClick={() => {
                  setUseBackupCode(!useBackupCode);
                  setTwoFACode('');
                  setError('');
                }}
                className="w-full text-sm text-macos-blue hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
              >
                {useBackupCode ? 'Use authenticator code' : 'Use backup code'}
              </button>
            </form>
          )}

          {/* Footer */}
          <div className="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
            <p>Default credentials:</p>
            <p className="font-mono">admin / admin</p>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
