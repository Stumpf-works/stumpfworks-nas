import { useWindowStore } from '@/store';
import Window from '@/components/Window';

export default function WindowManager() {
  const windows = useWindowStore((state) => state.windows);

  return (
    <div className="relative w-full h-full">
      {windows
        .filter((w) => w.state !== 'minimized')
        .map((window) => (
          <Window key={window.id} window={window} />
        ))}
    </div>
  );
}
