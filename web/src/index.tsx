import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { setChonkyDefaults } from 'chonky';
import { ChonkyIconFA } from 'chonky-icon-fontawesome';

import { MantineProvider } from '@mantine/core';
import { NotificationsProvider } from '@mantine/notifications';

setChonkyDefaults({ iconComponent: ChonkyIconFA });
const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <MantineProvider withNormalizeCSS withGlobalStyles>
    <NotificationsProvider autoClose={4000}>
      <App />
    </NotificationsProvider>
  </MantineProvider>
);
