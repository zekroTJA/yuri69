export function buildQueryParams(params: { [key: string]: any }): string {
  if (Object.keys(params).length === 0) return '';
  return (
    '?' +
    Object.keys(params)
      .filter((k) => !!params[k])
      .map((k) => `${k}=${params[k]}`)
      .join('&')
  );
}
