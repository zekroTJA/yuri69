export const sanitizeUid = (uid: string) => uid.replaceAll(/[^\w-]/g, '_').toLowerCase();

