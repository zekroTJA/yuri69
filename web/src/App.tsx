import { useState } from "react";
import { createGlobalStyle, ThemeProvider } from "styled-components";
import { DarkTheme } from "./theme/theme";

const GlobalStyle = createGlobalStyle`
  @import url("https://fonts.googleapis.com/css2?family=Rubik:wght@300;400;500&display=swap");

  body {
    font-family: 'Rubik', sans-serif;
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
  }

  * {
    box-sizing: border-box;
  }
`;

const App: React.FC = () => {
  return (
    <ThemeProvider theme={DarkTheme}>
      <div>
        <p>pog</p>
      </div>
      <GlobalStyle />
    </ThemeProvider>
  );
};

export default App;
