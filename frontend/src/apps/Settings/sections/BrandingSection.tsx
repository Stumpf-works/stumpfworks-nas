// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState } from 'react';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';

export function BrandingSection({ user, systemInfo }: { user: any; systemInfo: any }) {
  const [companyName, setCompanyName] = useState('StumpfWorks NAS');
  const [logoUrl, setLogoUrl] = useState('');
  const [primaryColor, setPrimaryColor] = useState('#0071e3');
  const [accentColor, setAccentColor] = useState('#5e5ce6');

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Branding</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Customize logo, colors, and theme for your organization
        </p>
      </div>

      {/* Company Information */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Company Information
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Company Name
              </label>
              <Input
                type="text"
                value={companyName}
                onChange={(e) => setCompanyName(e.target.value)}
                placeholder="Your Company Name"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Logo URL
              </label>
              <Input
                type="url"
                value={logoUrl}
                onChange={(e) => setLogoUrl(e.target.value)}
                placeholder="https://example.com/logo.png"
              />
              <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                Recommended size: 200x50 pixels
              </p>
            </div>
          </div>
        </div>
      </Card>

      {/* Colors */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Color Theme
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Primary Color
              </label>
              <div className="flex gap-3 items-center">
                <Input
                  type="color"
                  value={primaryColor}
                  onChange={(e) => setPrimaryColor(e.target.value)}
                  className="w-20 h-10"
                />
                <Input
                  type="text"
                  value={primaryColor}
                  onChange={(e) => setPrimaryColor(e.target.value)}
                  placeholder="#0071e3"
                  className="flex-1"
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Accent Color
              </label>
              <div className="flex gap-3 items-center">
                <Input
                  type="color"
                  value={accentColor}
                  onChange={(e) => setAccentColor(e.target.value)}
                  className="w-20 h-10"
                />
                <Input
                  type="text"
                  value={accentColor}
                  onChange={(e) => setAccentColor(e.target.value)}
                  placeholder="#5e5ce6"
                  className="flex-1"
                />
              </div>
            </div>
          </div>
        </div>
      </Card>

      {/* Preview */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Preview
          </h2>
          <div className="p-4 border-2 border-dashed border-gray-300 dark:border-macos-dark-300 rounded-lg">
            <div className="flex items-center gap-3 mb-3">
              {logoUrl ? (
                <img src={logoUrl} alt="Logo" className="h-8" />
              ) : (
                <div className="h-8 w-32 bg-gray-200 dark:bg-macos-dark-200 rounded flex items-center justify-center text-sm">
                  Logo
                </div>
              )}
              <span className="font-bold text-gray-900 dark:text-gray-100">{companyName}</span>
            </div>
            <div className="flex gap-2">
              <div
                className="px-4 py-2 rounded text-white font-medium"
                style={{ backgroundColor: primaryColor }}
              >
                Primary Button
              </div>
              <div
                className="px-4 py-2 rounded text-white font-medium"
                style={{ backgroundColor: accentColor }}
              >
                Accent Button
              </div>
            </div>
          </div>
        </div>
      </Card>

      {/* Actions */}
      <div className="flex gap-3">
        <Button variant="primary">
          Save Branding
        </Button>
        <Button variant="secondary">
          Reset to Defaults
        </Button>
      </div>

      {/* Note */}
      <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
        <p className="text-sm text-yellow-800 dark:text-yellow-200">
          Changes will be applied after saving and may require a page refresh.
        </p>
      </div>
    </div>
  );
}
