export enum AppTheme {
  DARK = 0,
  LIGHT = 1,
}

export const DarkTheme = {
  background: '#101018',
  background2: '#171721',
  background3: '#1d1d26',

  text: '#f4f4f5',

  accent: '#137ed6',

  white: '#f4f4f5',
  whiteDarker: '#dddddd',
  blurple: '#5865f2',
  blurpleDarker: '#4450d6',
  gray: '#455A64',
  darkGray: '#1e1e1e',
  red: '#ed4245',
  orange: '#f57c00',
  yellow: '#fbc02d',
  green: '#43a047',
  lime: '#57f287',
  cyan: '#03a9f4',
  pink: '#eb459e',
};

export const LightTheme: Theme = {
  ...DarkTheme,

  background: '#fffffe',
  background2: '#dddddd',
  background3: '#cecece',

  text: '#212121',
};

export const DefaultTheme = DarkTheme;
export type Theme = typeof DefaultTheme;
