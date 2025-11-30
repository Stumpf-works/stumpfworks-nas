import { useEffect, useState } from 'react';
import { addonsApi } from '@/api/addons';

// Map of addon IDs to their corresponding app IDs
const ADDON_TO_APP_MAP: Record<string, string> = {
  'vm-manager': 'vm-manager',
  'lxc-manager': 'lxc-manager',
};

/**
 * Hook to get list of installed addon app IDs
 * Returns array of app IDs that should be shown based on addon installation status
 */
export function useAddonApps() {
  const [installedAddonAppIds, setInstalledAddonAppIds] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchInstalledAddons = async () => {
      try {
        const response = await addonsApi.listAddons();
        if (response.success && response.data) {
          // Filter to only installed addons and map to their app IDs
          const installedAppIds = response.data
            .filter((addon) => addon.status.installed)
            .map((addon) => ADDON_TO_APP_MAP[addon.manifest.id])
            .filter((appId): appId is string => appId !== undefined);

          setInstalledAddonAppIds(installedAppIds);
        }
      } catch (error) {
        console.error('Failed to fetch installed addons:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchInstalledAddons();
  }, []);

  return { installedAddonAppIds, loading };
}
