// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { systemApi, type UpdateCheckResult } from '@/api/system';
import { getErrorMessage } from '@/api/client';

export function UpdatesSection() {
  const [updateCheckResult, setUpdateCheckResult] = useState<UpdateCheckResult | null>(null);
  const [checkingUpdates, setCheckingUpdates] = useState(false);
  const [updateError, setUpdateError] = useState<string | null>(null);

  const handleCheckForUpdates = async (forceCheck = false) => {
    setCheckingUpdates(true);
    setUpdateError(null);
    try {
      const response = await systemApi.checkForUpdates(forceCheck);
      if (response.success && response.data) {
        setUpdateCheckResult(response.data);
      } else {
        setUpdateError(response.error?.message || 'Failed to check for updates');
      }
    } catch (error) {
      console.error('Failed to check for updates:', error);
      setUpdateError(getErrorMessage(error));
    } finally {
      setCheckingUpdates(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Updates</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Check for software updates and view release information
        </p>
      </div>

      {/* Application Info */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Application Information
          </h2>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Application:</span>
              <span className="font-medium text-gray-900 dark:text-gray-100">
                Stumpf.Works NAS
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Version:</span>
              <span className="font-medium text-gray-900 dark:text-gray-100">v1.1.1</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">License:</span>
              <span className="font-medium text-gray-900 dark:text-gray-100">MIT</span>
            </div>
          </div>
        </div>
      </Card>

      {/* Software Updates */}
      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Software Updates
            </h2>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => handleCheckForUpdates(true)}
              disabled={checkingUpdates}
            >
              {checkingUpdates ? 'Checking...' : 'Check for Updates'}
            </Button>
          </div>

          {/* Update Check Result */}
          {updateCheckResult && (
            <div className="space-y-3">
              {updateCheckResult.updateAvailable ? (
                <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
                  <div className="flex items-start justify-between mb-2">
                    <div>
                      <p className="font-semibold text-blue-900 dark:text-blue-100">
                        Update Available!
                      </p>
                      <p className="text-sm text-blue-800 dark:text-blue-200 mt-1">
                        {updateCheckResult.currentVersion} → {updateCheckResult.latestVersion}
                      </p>
                    </div>
                  </div>

                  {/* Release Notes */}
                  {updateCheckResult.releaseInfo && (
                    <div className="mt-3">
                      <p className="text-sm font-medium text-blue-900 dark:text-blue-100 mb-1">
                        Release Notes:
                      </p>
                      <div className="text-xs text-blue-800 dark:text-blue-200 max-h-32 overflow-y-auto bg-blue-100/50 dark:bg-blue-900/10 p-2 rounded">
                        <pre className="whitespace-pre-wrap font-mono">
                          {updateCheckResult.releaseInfo.body || 'No release notes available'}
                        </pre>
                      </div>
                    </div>
                  )}

                  {/* Download Button */}
                  {updateCheckResult.releaseInfo?.html_url && (
                    <div className="mt-3">
                      <Button
                        variant="primary"
                        size="sm"
                        onClick={() =>
                          window.open(updateCheckResult.releaseInfo!.html_url, '_blank')
                        }
                        className="w-full"
                      >
                        Download Update from GitHub
                      </Button>
                    </div>
                  )}
                </div>
              ) : (
                <div className="p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
                  <p className="text-sm text-green-800 dark:text-green-200">
                    {updateCheckResult.message}
                  </p>
                </div>
              )}
            </div>
          )}

          {/* Error Message */}
          {updateError && (
            <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
              <p className="text-sm text-red-800 dark:text-red-200">{updateError}</p>
            </div>
          )}

          {/* Initial State */}
          {!updateCheckResult && !updateError && !checkingUpdates && (
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Click "Check for Updates" to see if a new version is available.
            </p>
          )}
        </div>
      </Card>
    </div>
  );
}
