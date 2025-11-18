import { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store';
import { authApi } from '@/api/auth';
import Login from '@/pages/Login';
import ResetPassword from '@/pages/ResetPassword';
import Desktop from '@/layout/Desktop';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

export default function App() {
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
      } else {
        // No token found - clear persisted auth state
        clearAuth();
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

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login onSuccess={() => {}} />
            </PublicRoute>
          }
        />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route
          path="/*"
          element={
            <ProtectedRoute>
              <Desktop />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}
