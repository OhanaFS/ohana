import React, { Suspense } from 'react';
import { Helmet } from 'react-helmet';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AppBase from './components/AppBase'; //switched to responsive base

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { AdminDashboard } from './AdminDashboard';
import { AdminConfiguration } from './AdminConfiguration';
import { AdminNodes } from './AdminNodes';
import { AdminPerformMaintenance } from './AdminPerformMaintenance';
import { AdminRunMaintenance } from './AdminRunMaintenance';
import AdminMaintenanceLogs from './AdminMaintenanceLogs';

const LoginPage = React.lazy(() => import('./LoginPage'));
const VFSBrowser = React.lazy(
  () => import('./components/FileBrowser/FileBrowser')
);
const SharingPage = React.lazy(
  () => import('./components/Sharing/SharingPage')
);
const SharedFavoritesList = React.lazy(
  () => import('./components/FileBrowser/SharedFavoritesList')
);

const queryClient = new QueryClient({
  defaultOptions: { queries: { keepPreviousData: true } },
});

export default function Demo() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Helmet>
          <title>Ohana</title>
        </Helmet>
        <Routes>
          <Route
            path="/blank"
            element={
              <AppBase userType="user">
                <></>
              </AppBase>
            }
          />
          <Route
            path="/home"
            element={
              <Suspense>
                <VFSBrowser />
              </Suspense>
            }
          />
          <Route
            path="/home/:id"
            element={
              <Suspense>
                <VFSBrowser />
              </Suspense>
            }
          />
          <Route
            path="/share/:id"
            element={
              <Suspense>
                <SharingPage />
              </Suspense>
            }
          />
          <Route
            path="/favorites"
            element={
              <Suspense>
                <SharedFavoritesList list="favorites" />
              </Suspense>
            }
          />
          <Route
            path="/shared"
            element={
              <Suspense>
                <SharedFavoritesList list="shared" />
              </Suspense>
            }
          />
          <Route
            path="/"
            element={
              <Suspense>
                <LoginPage />
              </Suspense>
            }
          />
          <Route
            path="/performmaintenance"
            element={<AdminPerformMaintenance />}
          />
          <Route path="/runmaintenance" element={<AdminRunMaintenance />} />
          <Route path="dashboard" element={<AdminDashboard />} />
          <Route path="/nodes" element={<AdminNodes />} />
          <Route path="/maintenance" element={<AdminMaintenanceLogs />} />
          <Route path="/settings" element={<AdminConfiguration />} />
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}
