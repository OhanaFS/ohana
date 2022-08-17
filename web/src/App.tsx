//import AppBase from "./AppBase";
import { Helmet } from 'react-helmet';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AppBase from './components/AppBase'; //switched to responsive base
import { VFSBrowser } from './components/FileBrowser/FileBrowser';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { AdminDashboard } from './AdminDashboard';
import { LoginPage } from './LoginPage';
import { AdminConfiguration } from './AdminConfiguration';
import { AdminNodes } from './AdminNodes';
import { AdminPerformMaintenance } from './AdminPerformMaintenance';
import { AdminRunMaintenance } from './AdminRunMaintenance';
import AdminMaintenanceLogs from './AdminMaintenanceLogs';
import SharingPage from './components/Sharing/SharingPage';

const queryClient = new QueryClient();

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
          <Route path="/home" element={<VFSBrowser />} />
          <Route path="/home/:id" element={<VFSBrowser />} />
          <Route path="/share/:id" element={<SharingPage />} />
          <Route path="/" element={<LoginPage />} />
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
