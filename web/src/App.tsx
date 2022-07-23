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
import { AdminMaintenanceResults } from './AdminMaintenanceResults';
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
              <AppBase
                userType="user"
                name="Alex"
                username="@alex"
                image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
              >
                <></>
              </AppBase>
            }
          />
          <Route path="/home" element={<VFSBrowser />} />
          <Route path="/" element={<LoginPage />} />
          <Route path="/insidessogroup" element={<AdminSsoGroupsInside />} />
          <Route path="/key_management" element={<AdminKeyManagement />} />
          <Route
            path="/maintenanceresults"
            element={<AdminMaintenanceResults />}
          />
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
