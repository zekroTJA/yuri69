import styled, { createGlobalStyle, ThemeProvider } from 'styled-components';
import { DarkTheme } from './theme/theme';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { MainRoute } from './routes/Main';
import { LoginRoute } from './routes/Login';
import { SoundsRoute } from './routes/Sounds';
import { useWSHooks } from './hooks/useWSHooks';
import './fonts.css';
import { SettingsRoute } from './routes/Settings';
import { UploadRoute } from './routes/Upload';
import 'react-contexify/dist/ReactContexify.css';
import { SnackBar } from './components/SnackBar';
import { EditRoute } from './routes/Edit';
import { WsDisconnectScreen } from './components/WsDisconnectScreen';
import { StatsRoute } from './routes/Stats';
import { useStore } from './store';
import { useEffect } from 'react';
import { useApi } from './hooks/useApi';
import { AdminRoute } from './routes/Admin';

const GlobalStyle = createGlobalStyle`
  body {
    font-family: 'Rubik', sans-serif;
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
    padding: 0;
    margin: 0;
  }

  * {
    box-sizing: border-box;
  }

  h1, h2, h3, h4, h5, h6 {
    margin-top: 0;
  }
`;

const Outlet = styled.div`
  width: 100vw;
  height: 100vh;
`;

const App: React.FC = () => {
  useWSHooks();

  return (
    <ThemeProvider theme={DarkTheme}>
      <Outlet>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<MainRoute />}>
              <Route index element={<SoundsRoute />} />
              <Route path="sounds/:uid" element={<EditRoute />} />
              <Route path="upload" element={<UploadRoute />} />
              <Route path="settings" element={<SettingsRoute />} />
              <Route path="stats" element={<StatsRoute />} />
              <Route path="admin" element={<AdminRoute />} />
            </Route>
            <Route path="/login" element={<LoginRoute />} />
            <Route path="*" element={<Navigate to="/" />} />
          </Routes>
        </BrowserRouter>
        <WsDisconnectScreen />
        <SnackBar />
      </Outlet>
      <GlobalStyle />
    </ThemeProvider>
  );
};

export default App;
