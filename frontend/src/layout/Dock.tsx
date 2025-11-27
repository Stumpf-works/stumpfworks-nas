import { useState } from 'react';
import { motion } from 'framer-motion';
import { Grid3x3, FolderPlus } from 'lucide-react';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  horizontalListSortingStrategy,
  useSortable,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useWindowStore, useDockStore, isDockFolder } from '@/store';
import { registeredApps, getAppById } from '@/apps';
import DockFolder from './DockFolder';

interface DockIconProps {
  id: string;
  app: typeof registeredApps[0];
  isRunning: boolean;
  onClick: () => void;
  onRemove: () => void;
}

function DockIcon({ id, app, isRunning, onClick, onRemove }: DockIconProps) {
  const [isHovered, setIsHovered] = useState(false);
  const [showContextMenu, setShowContextMenu] = useState(false);
  const [contextMenuPos, setContextMenuPos] = useState({ x: 0, y: 0 });

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    setContextMenuPos({ x: e.clientX, y: e.clientY });
    setShowContextMenu(true);
  };

  return (
    <>
      <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
        <motion.div
          className="relative flex flex-col items-center group"
          onHoverStart={() => setIsHovered(true)}
          onHoverEnd={() => setIsHovered(false)}
          whileHover={{ scale: 1.4, y: -10 }}
          whileTap={{ scale: 0.9 }}
          transition={{ type: 'spring', stiffness: 300, damping: 20 }}
        >
          <button
            onClick={onClick}
            onContextMenu={handleContextMenu}
            className="relative w-12 h-12 rounded-xl bg-white dark:bg-macos-dark-200 shadow-lg flex items-center justify-center text-2xl hover:shadow-xl transition-shadow cursor-grab active:cursor-grabbing"
          >
            {app.icon}
          </button>

        {/* Running Indicator */}
        {isRunning && (
          <div className="absolute -bottom-1 w-1 h-1 rounded-full bg-gray-800 dark:bg-white" />
        )}

        {/* Tooltip */}
        {isHovered && !showContextMenu && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="absolute -top-10 px-2 py-1 bg-gray-900/90 dark:bg-gray-100/90 text-white dark:text-gray-900 text-xs rounded whitespace-nowrap"
          >
            {app.name}
          </motion.div>
        )}
        </motion.div>
      </div>

      {/* Context Menu */}
      {showContextMenu && (
        <>
          <div
            className="fixed inset-0 z-50"
            onClick={() => setShowContextMenu(false)}
          />
          <div
            className="fixed z-50 bg-white dark:bg-macos-dark-200 rounded-lg shadow-2xl border border-gray-200 dark:border-gray-700 py-1 min-w-[160px]"
            style={{
              left: `${contextMenuPos.x}px`,
              top: `${contextMenuPos.y}px`,
            }}
          >
            <button
              onClick={() => {
                onRemove();
                setShowContextMenu(false);
              }}
              className="w-full px-4 py-2 text-left text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-macos-dark-300 transition-colors"
            >
              Remove from Dock
            </button>
          </div>
        </>
      )}
    </>
  );
}

interface SortableFolderProps {
  id: string;
  folder: any;
  onRemoveApp: (appId: string) => void;
  onDeleteFolder: () => void;
  onRenameFolder: () => void;
}

function SortableFolder({ id, folder, onRemoveApp, onDeleteFolder, onRenameFolder }: SortableFolderProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <DockFolder
        folder={folder}
        onRemoveApp={onRemoveApp}
        onDeleteFolder={onDeleteFolder}
        onRenameFolder={onRenameFolder}
      />
    </div>
  );
}

export default function Dock() {
  const [showCreateFolder, setShowCreateFolder] = useState(false);
  const [folderName, setFolderName] = useState('');
  const [selectedApps, setSelectedApps] = useState<string[]>([]);

  const windows = useWindowStore((state) => state.windows);
  const openWindow = useWindowStore((state) => state.openWindow);
  const dockItems = useDockStore((state) => state.dockItems);
  const removeFromDock = useDockStore((state) => state.removeFromDock);
  const reorderDock = useDockStore((state) => state.reorderDock);
  const createFolder = useDockStore((state) => state.createFolder);
  const deleteFolder = useDockStore((state) => state.deleteFolder);
  const removeAppFromFolder = useDockStore((state) => state.removeAppFromFolder);
  const renameFolder = useDockStore((state) => state.renameFolder);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      const oldIndex = dockItems.findIndex((item) => {
        if (typeof item === 'string') {
          return item === active.id;
        }
        return item.id === active.id;
      });
      const newIndex = dockItems.findIndex((item) => {
        if (typeof item === 'string') {
          return item === over.id;
        }
        return item.id === over.id;
      });

      if (oldIndex !== -1 && newIndex !== -1) {
        reorderDock(oldIndex, newIndex);
      }
    }
  };

  const handleAppClick = (app: typeof registeredApps[0]) => {
    openWindow(app);
  };

  const handleRemoveFromDock = (appId: string) => {
    removeFromDock(appId);
  };

  const handleCreateFolder = () => {
    if (folderName && selectedApps.length > 0) {
      createFolder(folderName, '📁', selectedApps);
      setShowCreateFolder(false);
      setFolderName('');
      setSelectedApps([]);
    }
  };

  const handleDeleteFolder = (folderId: string) => {
    if (confirm('Delete folder? Apps will be moved back to the dock.')) {
      deleteFolder(folderId);
    }
  };

  const handleRenameFolder = (folderId: string) => {
    const folder = dockItems.find((item) => isDockFolder(item) && item.id === folderId);
    if (!folder || !isDockFolder(folder)) return;

    const newName = prompt('Enter new folder name:', folder.name);
    if (newName && newName !== folder.name) {
      renameFolder(folderId, newName);
    }
  };

  // Get items with IDs for sortable
  const dockItemsWithIds = dockItems.map((item) => {
    if (typeof item === 'string') {
      return { id: item, type: 'app' as const, item };
    }
    return { id: item.id, type: 'folder' as const, item };
  });

  // Get all app IDs in dock (not in folders) for folder creation
  const availableApps = dockItems.filter((item) => typeof item === 'string') as string[];

  return (
    <>
      <motion.div
        initial={{ y: 100 }}
        animate={{ y: 0 }}
        transition={{ type: 'spring', stiffness: 100, damping: 20 }}
        className="fixed bottom-2 left-1/2 transform -translate-x-1/2 z-50"
      >
        <div className="px-4 py-2 glass-light dark:glass-dark rounded-2xl border border-gray-200/20 dark:border-gray-700/20 shadow-macos-xl">
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragEnd={handleDragEnd}
          >
            <SortableContext
              items={dockItemsWithIds.map((item) => item.id)}
              strategy={horizontalListSortingStrategy}
            >
              <div className="flex items-end space-x-3">
                {dockItemsWithIds.map((dockItem) => {
                  if (dockItem.type === 'app') {
                    const app = getAppById(dockItem.item);
                    if (!app) return null;
                    return (
                      <DockIcon
                        key={dockItem.id}
                        id={dockItem.id}
                        app={app}
                        isRunning={windows.some((w) => w.appId === app.id)}
                        onClick={() => handleAppClick(app)}
                        onRemove={() => handleRemoveFromDock(app.id)}
                      />
                    );
                  } else {
                    return (
                      <SortableFolder
                        key={dockItem.id}
                        id={dockItem.id}
                        folder={dockItem.item}
                        onRemoveApp={(appId) => removeAppFromFolder(dockItem.item.id, appId)}
                        onDeleteFolder={() => handleDeleteFolder(dockItem.item.id)}
                        onRenameFolder={() => handleRenameFolder(dockItem.item.id)}
                      />
                    );
                  }
                })}

                {/* Separator */}
                <div className="w-px h-12 bg-gray-300/50 dark:bg-gray-600/50 mx-1" />

                {/* Create Folder Button */}
                <motion.div
                  className="relative flex flex-col items-center group"
                  whileHover={{ scale: 1.4, y: -10 }}
                  whileTap={{ scale: 0.9 }}
                  transition={{ type: 'spring', stiffness: 300, damping: 20 }}
                >
                  <button
                    onClick={() => setShowCreateFolder(true)}
                    className="relative w-12 h-12 rounded-xl bg-gradient-to-br from-green-500 to-teal-600 shadow-lg flex items-center justify-center hover:shadow-xl transition-shadow"
                    title="Create Folder"
                  >
                    <FolderPlus className="w-6 h-6 text-white" />
                  </button>

                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    whileHover={{ opacity: 1, y: 0 }}
                    className="absolute -top-10 px-2 py-1 bg-gray-900/90 dark:bg-gray-100/90 text-white dark:text-gray-900 text-xs rounded whitespace-nowrap opacity-0 group-hover:opacity-100 pointer-events-none"
                  >
                    Create Folder
                  </motion.div>
                </motion.div>

                {/* App Gallery Launcher */}
                <motion.div
                  className="relative flex flex-col items-center group"
                  whileHover={{ scale: 1.4, y: -10 }}
                  whileTap={{ scale: 0.9 }}
                  transition={{ type: 'spring', stiffness: 300, damping: 20 }}
                >
                  <button
                    onClick={() => {
                      if ((window as any).openAppGallery) {
                        (window as any).openAppGallery();
                      }
                    }}
                    className="relative w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 shadow-lg flex items-center justify-center hover:shadow-xl transition-shadow"
                    title="App Gallery"
                  >
                    <Grid3x3 className="w-6 h-6 text-white" />
                  </button>

                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    whileHover={{ opacity: 1, y: 0 }}
                    className="absolute -top-10 px-2 py-1 bg-gray-900/90 dark:bg-gray-100/90 text-white dark:text-gray-900 text-xs rounded whitespace-nowrap opacity-0 group-hover:opacity-100 pointer-events-none"
                  >
                    App Gallery
                  </motion.div>
                </motion.div>
              </div>
            </SortableContext>
          </DndContext>
        </div>
      </motion.div>

      {/* Create Folder Modal */}
      {showCreateFolder && (
        <>
          <div
            className="fixed inset-0 bg-black/30 backdrop-blur-sm z-[100]"
            onClick={() => setShowCreateFolder(false)}
          />
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.9 }}
            className="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-[101] w-full max-w-md p-6 glass-light dark:glass-dark rounded-2xl border border-gray-200/20 dark:border-gray-700/20 shadow-macos-xl"
          >
            <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Create Folder
            </h2>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Folder Name
              </label>
              <input
                type="text"
                value={folderName}
                onChange={(e) => setFolderName(e.target.value)}
                placeholder="My Folder"
                className="w-full px-3 py-2 bg-white dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 dark:text-gray-100"
              />
            </div>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Select Apps (minimum 1)
              </label>
              <div className="max-h-64 overflow-y-auto border border-gray-300 dark:border-gray-600 rounded-lg p-2">
                {availableApps.map((appId) => {
                  const app = getAppById(appId);
                  if (!app) return null;
                  const isSelected = selectedApps.includes(appId);
                  return (
                    <label
                      key={appId}
                      className="flex items-center gap-2 p-2 hover:bg-gray-100 dark:hover:bg-macos-dark-300 rounded cursor-pointer"
                    >
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedApps([...selectedApps, appId]);
                          } else {
                            setSelectedApps(selectedApps.filter((id) => id !== appId));
                          }
                        }}
                        className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                      />
                      <span className="text-xl">{app.icon}</span>
                      <span className="text-sm text-gray-900 dark:text-gray-100">{app.name}</span>
                    </label>
                  );
                })}
              </div>
            </div>

            <div className="flex gap-2">
              <button
                onClick={handleCreateFolder}
                disabled={!folderName || selectedApps.length === 0}
                className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
              >
                Create
              </button>
              <button
                onClick={() => {
                  setShowCreateFolder(false);
                  setFolderName('');
                  setSelectedApps([]);
                }}
                className="flex-1 px-4 py-2 bg-gray-200 dark:bg-macos-dark-300 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-300 dark:hover:bg-macos-dark-400 transition-colors"
              >
                Cancel
              </button>
            </div>
          </motion.div>
        </>
      )}
    </>
  );
}
