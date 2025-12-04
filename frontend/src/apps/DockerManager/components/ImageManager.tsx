import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { dockerApi, DockerImage } from '@/api/docker';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';

export default function ImageManager() {
  const [images, setImages] = useState<DockerImage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [pullModal, setPullModal] = useState(false);
  const [pullImage, setPullImage] = useState('');
  const [pulling, setPulling] = useState(false);
  const [deleteModal, setDeleteModal] = useState<DockerImage | null>(null);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [buildModal, setBuildModal] = useState(false);
  const [buildDockerfile, setBuildDockerfile] = useState('');
  const [buildTags, setBuildTags] = useState('');
  const [building, setBuilding] = useState(false);
  const [pushModal, setPushModal] = useState<DockerImage | null>(null);
  const [pushing, setPushing] = useState(false);

  useEffect(() => {
    loadImages();
    const interval = setInterval(loadImages, 5000); // Refresh every 5s
    return () => clearInterval(interval);
  }, []);

  const loadImages = async () => {
    try {
      const response = await dockerApi.listImages();
      if (response.success && response.data) {
        setImages(response.data);
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load images');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handlePull = async () => {
    if (!pullImage.trim()) {
      alert('Please enter an image name');
      return;
    }

    setPulling(true);
    try {
      const response = await dockerApi.pullImage(pullImage.trim());
      if (response.success) {
        setPullModal(false);
        setPullImage('');
        loadImages();
      } else {
        alert(response.error?.message || 'Failed to pull image');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setPulling(false);
    }
  };

  const handleDelete = async (image: DockerImage) => {
    setActionLoading(image.id);
    try {
      const response = await dockerApi.removeImage(image.id);
      if (response.success) {
        setDeleteModal(null);
        loadImages();
      } else {
        alert(response.error?.message || 'Failed to remove image');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleBuild = async () => {
    if (!buildDockerfile.trim()) {
      alert('Please enter a Dockerfile');
      return;
    }

    if (!buildTags.trim()) {
      alert('Please enter at least one tag (e.g., myimage:latest)');
      return;
    }

    setBuilding(true);
    try {
      const tags = buildTags.split(',').map((t) => t.trim()).filter((t) => t);
      const response = await dockerApi.buildImage(buildDockerfile, tags);
      if (response.success) {
        setBuildModal(false);
        setBuildDockerfile('');
        setBuildTags('');
        loadImages();
        alert('Image built successfully!');
      } else {
        alert(response.error?.message || 'Failed to build image');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setBuilding(false);
    }
  };

  const handlePush = async (image: DockerImage) => {
    setPushing(true);
    try {
      const response = await dockerApi.pushImage(image.id);
      if (response.success) {
        setPushModal(null);
        alert('Image pushed successfully!');
      } else {
        alert(response.error?.message || 'Failed to push image');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setPushing(false);
    }
  };

  const formatSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (timestamp: number) => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
  };

  const getImageTag = (image: DockerImage) => {
    if (image.repoTags && image.repoTags.length > 0) {
      return image.repoTags[0];
    }
    return '<none>:<none>';
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Error Display */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Controls */}
      <div className="flex items-center justify-between">
        <div className="flex gap-2">
          <Button onClick={() => setPullModal(true)}>⬇️ Pull Image</Button>
          <Button onClick={() => setBuildModal(true)} variant="secondary">🔨 Build Image</Button>
        </div>
        <div className="text-sm text-gray-600 dark:text-gray-400">
          {images.length} image{images.length !== 1 ? 's' : ''}
        </div>
      </div>

      {/* Images Grid */}
      {images.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">💿</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No images found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            Pull an image to get started with Docker containers
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
          {images.map((image) => (
            <Card key={image.id} hoverable>
              <div className="p-6">
                {/* Header */}
                <div className="mb-4">
                  <h3 className="font-bold text-lg text-gray-900 dark:text-gray-100 mb-1">
                    {getImageTag(image)}
                  </h3>
                  <p className="text-xs text-gray-600 dark:text-gray-400 font-mono">
                    {image.id ? image.id.replace('sha256:', '').substring(0, 12) : 'N/A'}
                  </p>
                </div>

                {/* Details */}
                <div className="space-y-2 mb-4">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Size:</span>
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {formatSize(image.size)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Virtual Size:</span>
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {formatSize(image.virtualSize)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Containers:</span>
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {image.containers}
                    </span>
                  </div>
                  <div className="text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Created:</span>
                    <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                      {formatDate(image.created)}
                    </span>
                  </div>
                </div>

                {/* Tags */}
                {image.repoTags && image.repoTags.length > 1 && (
                  <div className="mb-4">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                      Additional Tags:
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {image.repoTags.slice(1).map((tag, idx) => (
                        <span
                          key={idx}
                          className="px-2 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 text-xs rounded"
                        >
                          {tag}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Actions */}
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    variant="secondary"
                    onClick={() => setPushModal(image)}
                    disabled={actionLoading === image.id}
                    className="flex-1"
                  >
                    ⬆️ Push
                  </Button>
                  <Button
                    size="sm"
                    variant="danger"
                    onClick={() => setDeleteModal(image)}
                    disabled={actionLoading === image.id}
                    className="flex-1"
                  >
                    🗑️ Delete
                  </Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Pull Image Modal */}
      <AnimatePresence>
        {pullModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !pulling && setPullModal(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Pull Docker Image
              </h2>

              <div className="mb-4">
                <Input
                  label="Image Name"
                  value={pullImage}
                  onChange={(e) => setPullImage(e.target.value)}
                  placeholder="nginx:latest"
                  required
                  disabled={pulling}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter' && !pulling) {
                      handlePull();
                    }
                  }}
                />
                <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                  Examples: nginx:latest, ubuntu:22.04, redis:alpine
                </p>
              </div>

              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setPullModal(false)}
                  className="flex-1"
                  disabled={pulling}
                >
                  Cancel
                </Button>
                <Button onClick={handlePull} className="flex-1" disabled={pulling}>
                  {pulling ? 'Pulling...' : 'Pull Image'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {deleteModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setDeleteModal(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Delete Image
              </h2>
              <p className="text-gray-600 dark:text-gray-400 mb-6">
                Are you sure you want to delete image <strong>{getImageTag(deleteModal)}</strong>?
                This action cannot be undone.
              </p>
              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setDeleteModal(null)}
                  className="flex-1"
                  disabled={actionLoading === deleteModal.id}
                >
                  Cancel
                </Button>
                <Button
                  variant="danger"
                  onClick={() => handleDelete(deleteModal)}
                  className="flex-1"
                  disabled={actionLoading === deleteModal.id}
                >
                  {actionLoading === deleteModal.id ? 'Deleting...' : 'Delete'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Build Image Modal */}
      <AnimatePresence>
        {buildModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !building && setBuildModal(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Build Docker Image
              </h2>

              <div className="mb-4">
                <Input
                  label="Image Tag(s)"
                  value={buildTags}
                  onChange={(e) => setBuildTags(e.target.value)}
                  placeholder="myimage:latest, myimage:v1.0"
                  required
                  disabled={building}
                />
                <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                  Enter one or more tags separated by commas
                </p>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Dockerfile Content <span className="text-red-500">*</span>
                </label>
                <textarea
                  value={buildDockerfile}
                  onChange={(e) => setBuildDockerfile(e.target.value)}
                  placeholder={'FROM node:18-alpine\nWORKDIR /app\nCOPY . .\nRUN npm install\nCMD ["npm", "start"]'}
                  disabled={building}
                  className="w-full h-64 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-macos-blue resize-none"
                />
                <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                  Paste your Dockerfile content here
                </p>
              </div>

              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setBuildModal(false)}
                  className="flex-1"
                  disabled={building}
                >
                  Cancel
                </Button>
                <Button onClick={handleBuild} className="flex-1" disabled={building}>
                  {building ? 'Building...' : 'Build Image'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Push Image Modal */}
      <AnimatePresence>
        {pushModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !pushing && setPushModal(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Push Image
              </h2>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Push image <strong>{getImageTag(pushModal)}</strong> to registry?
              </p>
              <p className="text-xs text-gray-600 dark:text-gray-400 mb-6">
                Note: Make sure you're logged in to the registry using 'docker login' on the host system.
              </p>
              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setPushModal(null)}
                  className="flex-1"
                  disabled={pushing}
                >
                  Cancel
                </Button>
                <Button
                  onClick={() => handlePush(pushModal)}
                  className="flex-1"
                  disabled={pushing}
                >
                  {pushing ? 'Pushing...' : 'Push Image'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
