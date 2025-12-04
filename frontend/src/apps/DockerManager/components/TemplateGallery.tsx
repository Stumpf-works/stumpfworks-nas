import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { dockerApi } from '@/api/docker';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import Modal from '@/components/ui/Modal';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { Play, Info, Tag, User, Calendar, Loader2 } from 'lucide-react';
import { toast } from 'react-hot-toast';

interface Template {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  author: string;
  version: string;
  compose: string;
  variables: Record<string, string>;
  requirements?: {
    min_memory_mb?: number;
    min_disk_gb?: number;
    ports?: number[];
    notes?: string[];
  };
}

interface DeployModalProps {
  template: Template | null;
  isOpen: boolean;
  onClose: () => void;
  onDeploy: (stackName: string, variables: Record<string, string>) => void;
  deploying?: boolean;
}

function DeployModal({ template, isOpen, onClose, onDeploy, deploying = false }: DeployModalProps) {
  const [stackName, setStackName] = useState('');
  const [variables, setVariables] = useState<Record<string, string>>({});

  useEffect(() => {
    if (template) {
      setStackName(template.id);
      setVariables({ ...template.variables });
    }
  }, [template]);

  if (!template) return null;

  const handleDeploy = () => {
    if (stackName.trim()) {
      onDeploy(stackName, variables);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Deploy ${template.name}`}>
      <div className="space-y-6">
        {/* Stack Name */}
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Stack Name
          </label>
          <Input
            type="text"
            value={stackName}
            onChange={(e) => setStackName(e.target.value)}
            placeholder="Enter stack name"
            required
          />
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Unique name for this deployment
          </p>
        </div>

        {/* Variables */}
        {Object.keys(template.variables).length > 0 && (
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              Configuration Variables
            </label>
            <div className="space-y-3">
              {Object.entries(template.variables).map(([key, defaultValue]) => (
                <div key={key}>
                  <label className="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                    {key}
                  </label>
                  <Input
                    type="text"
                    value={variables[key] || defaultValue}
                    onChange={(e) => setVariables({ ...variables, [key]: e.target.value })}
                    placeholder={defaultValue}
                  />
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Requirements */}
        {template.requirements && (
          <div className="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-lg">
            <h4 className="text-sm font-medium text-blue-900 dark:text-blue-100 mb-2">
              Requirements
            </h4>
            <ul className="space-y-1 text-sm text-blue-800 dark:text-blue-200">
              {template.requirements.min_memory_mb && (
                <li>• Minimum Memory: {template.requirements.min_memory_mb} MB</li>
              )}
              {template.requirements.min_disk_gb && (
                <li>• Minimum Disk: {template.requirements.min_disk_gb} GB</li>
              )}
              {template.requirements.ports && template.requirements.ports.length > 0 && (
                <li>• Ports: {template.requirements.ports.join(', ')}</li>
              )}
              {template.requirements.notes && template.requirements.notes.length > 0 && (
                <>
                  {template.requirements.notes.map((note, idx) => (
                    <li key={idx}>• {note}</li>
                  ))}
                </>
              )}
            </ul>
          </div>
        )}

        {/* Actions */}
        <div className="flex justify-end gap-3 pt-4">
          <Button variant="secondary" onClick={onClose} disabled={deploying}>
            Cancel
          </Button>
          <Button
            variant="primary"
            onClick={handleDeploy}
            disabled={!stackName.trim() || deploying}
          >
            {deploying ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                Deploying...
              </>
            ) : (
              <>
                <Play className="w-4 h-4 mr-2" />
                Deploy
              </>
            )}
          </Button>
        </div>
      </div>
    </Modal>
  );
}

export default function TemplateGallery() {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);
  const [showDeployModal, setShowDeployModal] = useState(false);
  const [showDetailsModal, setShowDetailsModal] = useState(false);
  const [deploying, setDeploying] = useState(false);

  useEffect(() => {
    fetchTemplates();
    fetchCategories();
  }, []);

  const fetchTemplates = async () => {
    try {
      setLoading(true);
      const response = await dockerApi.listTemplates();
      if (response.success && response.data) {
        setTemplates(response.data);
        setError(null);
      } else {
        setError(response.error?.message || 'Failed to load templates');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const response = await dockerApi.getTemplateCategories();
      if (response.success && response.data) {
        setCategories(['all', ...response.data]);
      }
    } catch (err) {
      console.error('Failed to fetch categories:', err);
    }
  };

  const handleDeploy = async (stackName: string, variables: Record<string, string>) => {
    if (!selectedTemplate) return;

    setDeploying(true);
    const toastId = toast.loading(`Deploying ${selectedTemplate.name}...`);

    try {
      const response = await dockerApi.deployTemplate(selectedTemplate.id, stackName, variables);
      if (response.success) {
        toast.success(`Successfully deployed ${selectedTemplate.name} as "${stackName}"!`, { id: toastId });
        setShowDeployModal(false);
        setSelectedTemplate(null);

        // Refresh templates after short delay
        setTimeout(() => fetchTemplates(), 1000);
      } else {
        toast.error(response.error?.message || 'Failed to deploy template', { id: toastId });
      }
    } catch (err) {
      toast.error(getErrorMessage(err), { id: toastId });
    } finally {
      setDeploying(false);
    }
  };

  const filteredTemplates = selectedCategory === 'all'
    ? templates
    : templates.filter(t => t.category === selectedCategory);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-red-600 dark:text-red-400">
        <p>{error}</p>
        <Button variant="secondary" onClick={fetchTemplates} className="mt-4">
          Retry
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Template Gallery</h2>
        <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
          One-click deployment of popular media servers and applications
        </p>
      </div>

      {/* Category Filter */}
      <div className="flex gap-2 flex-wrap">
        {categories.map((category) => (
          <button
            key={category}
            onClick={() => setSelectedCategory(category)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              selectedCategory === category
                ? 'bg-macos-blue text-white'
                : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
            }`}
          >
            {category.charAt(0).toUpperCase() + category.slice(1)}
          </button>
        ))}
      </div>

      {/* Template Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredTemplates.map((template) => (
          <motion.div
            key={template.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.2 }}
          >
            <Card className="hover:shadow-lg transition-shadow cursor-pointer">
              <div className="p-6 space-y-4">
                {/* Icon & Title */}
                <div className="flex items-start gap-4">
                  <div className="text-4xl">{template.icon}</div>
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900 dark:text-white">
                      {template.name}
                    </h3>
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      <Tag className="w-3 h-3 inline mr-1" />
                      {template.category}
                    </p>
                  </div>
                </div>

                {/* Description */}
                <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-3">
                  {template.description}
                </p>

                {/* Meta */}
                <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
                  <span className="flex items-center gap-1">
                    <User className="w-3 h-3" />
                    {template.author}
                  </span>
                  <span className="flex items-center gap-1">
                    <Calendar className="w-3 h-3" />
                    v{template.version}
                  </span>
                </div>

                {/* Actions */}
                <div className="flex gap-2 pt-2">
                  <Button
                    variant="primary"
                    size="sm"
                    onClick={() => {
                      setSelectedTemplate(template);
                      setShowDeployModal(true);
                    }}
                    className="flex-1"
                  >
                    <Play className="w-4 h-4 mr-1" />
                    Deploy
                  </Button>
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => {
                      setSelectedTemplate(template);
                      setShowDetailsModal(true);
                    }}
                  >
                    <Info className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </Card>
          </motion.div>
        ))}
      </div>

      {/* No templates */}
      {filteredTemplates.length === 0 && (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No templates found in this category
        </div>
      )}

      {/* Deploy Modal */}
      <DeployModal
        template={selectedTemplate}
        isOpen={showDeployModal}
        onClose={() => setShowDeployModal(false)}
        onDeploy={handleDeploy}
        deploying={deploying}
      />

      {/* Details Modal */}
      {selectedTemplate && (
        <Modal
          isOpen={showDetailsModal}
          onClose={() => setShowDetailsModal(false)}
          title={selectedTemplate.name}
        >
          <div className="space-y-4">
            <p className="text-gray-600 dark:text-gray-400">{selectedTemplate.description}</p>

            <div className="space-y-2">
              <h4 className="font-medium text-gray-900 dark:text-white">Author</h4>
              <p className="text-sm text-gray-600 dark:text-gray-400">{selectedTemplate.author}</p>
            </div>

            <div className="space-y-2">
              <h4 className="font-medium text-gray-900 dark:text-white">Docker Compose</h4>
              <pre className="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg text-xs overflow-x-auto">
                <code>{selectedTemplate.compose}</code>
              </pre>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
}
