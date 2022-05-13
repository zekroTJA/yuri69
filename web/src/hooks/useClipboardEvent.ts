import { useEffect } from 'react';

export const useClipboardEvent = (handler: (ev: ClipboardEvent) => void) => {
  useEffect(() => {
    document.addEventListener('paste', handler);
    return () => document.removeEventListener('paste', handler);
  }, []);
};
