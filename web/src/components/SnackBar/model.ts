export type SnackBarType = 'error' | 'warning' | 'info' | 'success';

export type SnackBarModel = {
  show: boolean;
  type: SnackBarType;
  content: JSX.Element;
};
