// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import Card from '@/components/ui/Card';
import { TwoFactorAuth } from '@/components/TwoFactorAuth/TwoFactorAuth';

export function GeneralSection({ user }: { user: any; systemInfo: any }) {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">General</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          User account and security settings
        </p>
      </div>

      {/* User Information */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            User Information
          </h2>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Username:</span>
              <span className="font-medium text-gray-900 dark:text-gray-100">
                {user?.username}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Email:</span>
              <span className="font-medium text-gray-900 dark:text-gray-100">
                {user?.email}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Role:</span>
              <span className="font-medium text-gray-900 dark:text-gray-100 capitalize">
                {user?.role}
              </span>
            </div>
          </div>
        </div>
      </Card>

      {/* Two-Factor Authentication */}
      <Card>
        <div className="p-6">
          <TwoFactorAuth />
        </div>
      </Card>
    </div>
  );
}
