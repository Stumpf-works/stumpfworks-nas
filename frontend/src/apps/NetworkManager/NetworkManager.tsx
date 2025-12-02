import NetworkConfig from './components/NetworkConfig';

export function NetworkManager() {
  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      <NetworkConfig />
    </div>
  );
}
