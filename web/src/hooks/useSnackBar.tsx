import { useRef } from 'react';
import { SnackBarType } from '../components/SnackBar';
import { useStore } from '../store';

export const useSnackBar = () => {
  const [snackBar, setSnackBar] = useStore((s) => [s.snackBar, s.setSnackBar]);
  const timerRef = useRef<ReturnType<typeof setTimeout>>();

  const hide = () => {
    setSnackBar({ show: false });
  };

  const show = (
    content: string | JSX.Element,
    type: SnackBarType = 'info',
    duration: number = 4000,
  ) => {
    if (timerRef.current) clearTimeout(timerRef.current);
    const _content = typeof content === 'string' ? <span>{content}</span> : content;
    setSnackBar({ show: true, content: _content, type });
    timerRef.current = setTimeout(() => hide(), duration);
  };

  return { show, hide };
};
