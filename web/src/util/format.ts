export const possessive = (name: string) => {
  return name.endsWith('s') ? name + "'" : name + "'s";
};
