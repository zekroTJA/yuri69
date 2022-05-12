import { format } from 'date-fns';

export const formatDate = (date: string | Date) => {
  if (!date) return 'n/a';
  const _date = typeof date === 'string' ? new Date(date) : date;
  return format(_date, 'dd/LL/yyyy HH:mm:ss O');
};
