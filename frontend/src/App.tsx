import { useEffect, useState } from 'react';
import { useAuthStore } from '@/store';
import { authApi } from '@/api/auth';
import Login from '@/pages/Login';
import Desktop from '@/layout/Desktop';

export default function App() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const setAuth = useAuthStore((state) => state.setAuth);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    // Check if user is still authenticated on mount
    const checkAuth = async () => {
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
      }
      setIsChecking(false);
    };

    checkAuth();
  }, [setAuth, clearAuth]);

  if (isChecking) {
    return (
      <div className="w-screen h-screen flex items-center justify-center bg-gray-50 dark:bg-macos-dark-50">
        <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-macos-blue" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Login onSuccess={() => {}} />;
  }

  return <Desktop />;
}
