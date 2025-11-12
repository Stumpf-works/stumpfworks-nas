import React, { useState } from 'react';
import { FileInfo, getFileIcon } from '@/api/files';

interface FileThumbnailProps {
  file: FileInfo;
  size?: 'small' | 'medium' | 'large';
}

const isImageFile = (file: FileInfo): boolean => {
  if (!file.mimeType) return false;
  return file.mimeType.startsWith('image/');
};

const FileThumbnail: React.FC<FileThumbnailProps> = ({ file, size = 'medium' }) => {
  const [imageError, setImageError] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);

  // Size classes
  const sizeClasses = {
    small: 'w-12 h-12',
    medium: 'w-20 h-20',
    large: 'w-32 h-32',
  };

  const iconSizes = {
    small: 'text-3xl',
    medium: 'text-5xl',
    large: 'text-7xl',
  };

  // If not an image or image failed to load, show icon
  if (!isImageFile(file) || imageError) {
    return <div className={iconSizes[size]}>{getFileIcon(file)}</div>;
  }

  // Construct image URL
  const token = localStorage.getItem('token');
  const imageUrl = `/api/files/download?path=${encodeURIComponent(file.path)}&token=${token}`;

  return (
    <div className={`${sizeClasses[size]} relative flex items-center justify-center`}>
      {/* Show icon while loading */}
      {!imageLoaded && (
        <div className={`absolute inset-0 flex items-center justify-center ${iconSizes[size]}`}>
          {getFileIcon(file)}
        </div>
      )}

      {/* Thumbnail image */}
      <img
        src={imageUrl}
        alt={file.name}
        loading="lazy"
        className={`${sizeClasses[size]} object-cover rounded-md transition-opacity duration-300 ${
          imageLoaded ? 'opacity-100' : 'opacity-0'
        }`}
        onLoad={() => setImageLoaded(true)}
        onError={() => setImageError(true)}
      />
    </div>
  );
};

export default FileThumbnail;
