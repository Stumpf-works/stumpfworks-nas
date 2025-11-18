import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { authApi } from '@/api/auth';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { CheckCircle, AlertCircle } from 'lucide-react';

export default function ResetPassword() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token');

  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!token) {
      setError('Invalid reset link. No token provided.');
    }
  }, [token]);

  const validatePassword = (password: string): string | null => {
    if (password.length < 8) {
      return 'Password must be at least 8 characters long';
    }
    return null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validate token
    if (!token) {
      setError('Invalid reset link');
      return;
    }

    // Validate password
    const passwordError = validatePassword(newPassword);
    if (passwordError) {
      setError(passwordError);
      return;
    }

    // Check if passwords match
    if (newPassword !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setIsLoading(true);

    try {
      const response = await authApi.resetPassword({
        token,
        newPassword,
      });

      if (response.success) {
        setSuccess(true);
        // Redirect to login after 3 seconds
        setTimeout(() => {
          navigate('/login');
        }, 3000);
      } else {
        setError(response.error?.message || 'Password reset failed');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  if (success) {
    return (
      <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 dark:from-macos-dark-50 dark:via-macos-dark-100 dark:to-macos-dark-200">
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          className="w-full max-w-md"
        >
          <div className="bg-white dark:bg-macos-dark-300 rounded-2xl shadow-2xl p-8">
            <div className="text-center">
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: 0.2, type: 'spring', stiffness: 200 }}
                className="mx-auto mb-6 w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center"
              >
                <CheckCircle className="w-12 h-12 text-green-600 dark:text-green-400" />
              </motion.div>

              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
                Password Reset Successful!
              </h2>

              <p className="text-gray-600 dark:text-gray-400 mb-6">
                Your password has been changed successfully.
                <br />
                You can now log in with your new password.
              </p>

              <p className="text-sm text-gray-500 dark:text-gray-500">
                Redirecting to login page...
              </p>
            </div>
          </div>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 dark:from-macos-dark-50 dark:via-macos-dark-100 dark:to-macos-dark-200">
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        className="w-full max-w-md"
      >
        <div className="bg-white dark:bg-macos-dark-300 rounded-2xl shadow-2xl p-8">
          <div className="text-center mb-8">
            <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
              Reset Your Password
            </h2>
            <p className="text-gray-600 dark:text-gray-400">
              Enter your new password below
            </p>
          </div>

          {error && (
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-start gap-3"
            >
              <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-800 dark:text-red-300">{error}</p>
            </motion.div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <Input
              label="New Password"
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              placeholder="Enter new password"
              required
              autoFocus
              disabled={isLoading || !token}
            />

            <Input
              label="Confirm Password"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Confirm new password"
              required
              disabled={isLoading || !token}
            />

            <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
              <p className="text-sm text-blue-800 dark:text-blue-300 font-medium mb-2">
                Password Requirements:
              </p>
              <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-1 list-disc list-inside">
                <li>At least 8 characters long</li>
                <li>Mix of letters, numbers, and special characters recommended</li>
              </ul>
            </div>

            <Button
              type="submit"
              variant="primary"
              className="w-full"
              disabled={isLoading || !token}
            >
              {isLoading ? 'Resetting Password...' : 'Reset Password'}
            </Button>

            <Button
              type="button"
              variant="secondary"
              className="w-full"
              onClick={() => navigate('/login')}
              disabled={isLoading}
            >
              Back to Login
            </Button>
          </form>
        </div>
      </motion.div>
    </div>
  );
}
