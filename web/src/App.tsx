//import AppBase from "./AppBase";
import { Helmet } from 'react-helmet';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AppBase from './components/AppBase'; //switched to responsive base
import { VFSBrowser } from './components/userFiles';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { AdminDashboard } from './AdminDashboard';
import { LoginPage } from './LoginPage';
import { AdminConfiguration } from './AdminConfiguration';
import { AdminKeyManagement } from './AdminKeyManagement';
import { AdminNodes } from './AdminNodes';
import { AdminPerformMaintenance } from './AdminPerformMaintenance';
import { AdminRunMaintenance } from './AdminRunMaintenance';
import { AdminSettings } from './AdminSettings';
import { AdminSsoGroups } from './AdminSsoGroups';
import { AdminSsoGroupsInside } from './AdminSsoGroupsInside';
import AdminMaintenanceLogs from './AdminMaintenanceLogs';

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
              <AppBase>
                <></>
              </AppBase>
            }
          />
          <Route path="/home" element={<VFSBrowser />} />
          <Route path="/" element={<LoginPage />} />
          <Route path="/insidessogroup" element={<AdminSsoGroupsInside />} />
          <Route path="/key_management" element={<AdminKeyManagement />} />
          <Route
            path="/performmaintenance"
            element={<AdminPerformMaintenance />}
          />
          <Route path="/runmaintenance" element={<AdminRunMaintenance />} />
          <Route path="dashboard" element={<AdminDashboard />} />
          <Route path="/sso" element={<AdminSsoGroups />} />
          <Route path="/nodes" element={<AdminNodes />} />
          <Route path="/maintenance" element={<AdminMaintenanceLogs />} />
          <Route path="/settings" element={<AdminSettings />} />
          <Route path="/rotate_key" element={<AdminConfiguration />} />
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}
