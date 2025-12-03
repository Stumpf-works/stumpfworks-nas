import { useState, useEffect } from 'react';
import { cloudBackupApi, CloudProvider } from '@/api/cloudbackup';
import { getErrorMessage } from '@/api/client';
import Modal from '@/components/ui/Modal';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import { Cloud } from 'lucide-react';

interface ProviderModalProps {
  provider: CloudProvider | null;
  onClose: () => void;
  onSaved: () => void;
}

interface ProviderField {
  key: string;
  label: string;
  type: string;
  required: boolean;
  placeholder?: string;
}

const PROVIDER_TYPES: Record<string, { name: string; fields: ProviderField[] }> = {
  s3: {
    name: 'Amazon S3',
    fields: [
      { key: 'access_key_id', label: 'Access Key ID', type: 'text', required: true },
      { key: 'secret_access_key', label: 'Secret Access Key', type: 'password', required: true },
      { key: 'region', label: 'Region', type: 'text', required: true, placeholder: 'us-east-1' },
      { key: 'endpoint', label: 'Endpoint (optional)', type: 'text', required: false, placeholder: 'https://s3.amazonaws.com' },
    ],
  },
  b2: {
    name: 'Backblaze B2',
    fields: [
      { key: 'account', label: 'Account ID', type: 'text', required: true },
      { key: 'key', label: 'Application Key', type: 'password', required: true },
    ],
  },
  gdrive: {
    name: 'Google Drive',
    fields: [
      { key: 'client_id', label: 'Client ID', type: 'text', required: true },
      { key: 'client_secret', label: 'Client Secret', type: 'password', required: true },
    ],
  },
  dropbox: {
    name: 'Dropbox',
    fields: [
      { key: 'client_id', label: 'App Key', type: 'text', required: true },
      { key: 'client_secret', label: 'App Secret', type: 'password', required: true },
    ],
  },
  onedrive: {
    name: 'Microsoft OneDrive',
    fields: [
      { key: 'client_id', label: 'Client ID', type: 'text', required: true },
      { key: 'client_secret', label: 'Client Secret', type: 'password', required: true },
    ],
  },
  azureblob: {
    name: 'Azure Blob Storage',
    fields: [
      { key: 'account', label: 'Storage Account Name', type: 'text', required: true },
      { key: 'key', label: 'Storage Account Key', type: 'password', required: true },
    ],
  },
  sftp: {
    name: 'SFTP',
    fields: [
      { key: 'host', label: 'Host', type: 'text', required: true },
      { key: 'user', label: 'Username', type: 'text', required: true },
      { key: 'pass', label: 'Password', type: 'password', required: false },
      { key: 'key_file', label: 'SSH Key File (optional)', type: 'text', required: false },
      { key: 'port', label: 'Port', type: 'number', required: false, placeholder: '22' },
    ],
  },
};

export function ProviderModal({ provider, onClose, onSaved }: ProviderModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [type, setType] = useState<keyof typeof PROVIDER_TYPES>('s3');
  const [enabled, setEnabled] = useState(true);
  const [config, setConfig] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (provider) {
      setName(provider.name);
      setDescription(provider.description);
      setType(provider.type as keyof typeof PROVIDER_TYPES);
      setEnabled(provider.enabled);
      try {
        setConfig(JSON.parse(provider.config));
      } catch {
        setConfig({});
      }
    }
  }, [provider]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSaving(true);

    try {
      const providerData = {
        name,
        description,
        type,
        enabled,
        config: JSON.stringify(config),
      };

      let response;
      if (provider) {
        response = await cloudBackupApi.updateProvider(provider.id, providerData);
      } else {
        response = await cloudBackupApi.createProvider(providerData);
      }

      if (response.success) {
        onSaved();
      } else {
        setError(response.error?.message || 'Failed to save provider');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setSaving(false);
    }
  };

  const handleConfigChange = (key: string, value: string) => {
    setConfig({ ...config, [key]: value });
  };

  const currentProviderType = PROVIDER_TYPES[type];

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title={provider ? 'Edit Cloud Provider' : 'Add Cloud Provider'}
      icon={<Cloud className="w-6 h-6" />}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded text-red-600 dark:text-red-400 text-sm">
            {error}
          </div>
        )}

        <Input
          label="Provider Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="My Cloud Storage"
        />

        <Input
          label="Description (optional)"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Personal backup storage"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Provider Type
          </label>
          <select
            value={type}
            onChange={(e) => {
              setType(e.target.value as keyof typeof PROVIDER_TYPES);
              setConfig({});
            }}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
            disabled={!!provider}
          >
            {Object.entries(PROVIDER_TYPES).map(([key, value]) => (
              <option key={key} value={key}>
                {value.name}
              </option>
            ))}
          </select>
          {provider && (
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Provider type cannot be changed after creation
            </p>
          )}
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
          <h3 className="text-sm font-medium text-gray-900 dark:text-gray-100 mb-3">
            {currentProviderType.name} Configuration
          </h3>
          <div className="space-y-3">
            {currentProviderType.fields.map((field) => (
              <Input
                key={field.key}
                label={field.label}
                type={field.type}
                value={config[field.key] || ''}
                onChange={(e) => handleConfigChange(field.key, e.target.value)}
                required={field.required}
                placeholder={field.placeholder || ''}
              />
            ))}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            id="enabled"
            checked={enabled}
            onChange={(e) => setEnabled(e.target.checked)}
            className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
          />
          <label htmlFor="enabled" className="text-sm text-gray-700 dark:text-gray-300">
            Enable this provider
          </label>
        </div>

        <div className="flex gap-3 pt-4">
          <Button type="submit" disabled={saving}>
            {saving ? 'Saving...' : provider ? 'Update Provider' : 'Add Provider'}
          </Button>
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
        </div>
      </form>
    </Modal>
  );
}
