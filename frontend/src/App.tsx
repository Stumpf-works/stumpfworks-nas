import { useEffect, useState } from 'react';
import { useAuthStore } from '@/store';
import { authApi } from '@/api/auth';
import { setupApi } from '@/api/setup';
import Login from '@/pages/Login';
import SetupWizard from '@/pages/SetupWizard';
import Desktop from '@/layout/Desktop';

export default function App() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const setAuth = useAuthStore((state) => state.setAuth);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const [isChecking, setIsChecking] = useState(true);
  const [setupRequired, setSetupRequired] = useState(false);

  useEffect(() => {
    // Check setup status and authentication on mount
    const checkSetupAndAuth = async () => {
      try {
        // First check if setup is required
        const setupStatus = await setupApi.getStatus();
        if (setupStatus.setupRequired) {
          setSetupRequired(true);
          setIsChecking(false);
          return;
        }
      } catch (error) {
        // If setup check fails, continue to auth check
        console.error('Setup check failed:', error);
      }

      // Check if user is still authenticated
      const token = localStorage.getItem('accessToken');
      if (token) {
        try {
          const response = await authApi.getCurrentUser();
          if (response.success && response.data) {
            // User is still authenticated
            setIsChecking(false);
            return;
          }
        } catch (error) {
          // Token is invalid, clear auth
          clearAuth();
        }
      } else {
        // No token found - clear persisted auth state
        clearAuth();
      }
      setIsChecking(false);
    };

    checkSetupAndAuth();
  }, [setAuth, clearAuth]);

  if (isChecking) {
    return (
      <div className="w-screen h-screen flex items-center justify-center bg-gray-50 dark:bg-macos-dark-50">
        <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-macos-blue" />
      </div>
    );
  }

  // Show setup wizard if setup is required
  if (setupRequired) {
    return <SetupWizard onComplete={() => {
      setSetupRequired(false);
      // Force a re-check of auth after setup
      setIsChecking(true);
      setTimeout(() => {
        window.location.reload();
      }, 1000);
    }} />;
  }

  if (!isAuthenticated) {
    return <Login onSuccess={() => {}} />;
  }

  return <Desktop />;
}
