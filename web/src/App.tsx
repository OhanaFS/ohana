//import AppBase from "./AppBase";
import { Helmet } from 'react-helmet';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AppBase from './components/AppBase'; //switched to responsive base
import { VFSBrowser } from './components/userFiles';

import { AdminDashboard } from './AdminDashboard';
import {LoginPage} from './LoginPage';
import { AdminConfiguration } from './AdminConfiguration';
import { AdminCreateKey } from './AdminCreateKey';
import { AdminCreateSsoKey } from './AdminCreateSsoKey';
import { AdminKeyManagement } from './AdminKeyManagement';
import { AdminMaintenanceLogs } from './AdminMaintenanceLogs';
import { AdminMaintenanceResults } from './AdminMaintenanceResults';
import { AdminMaintenanceSettings } from './AdminMaintenanceSettings';
import { AdminNodes } from './AdminNodes';
import { AdminPerformMaintenance } from './AdminPerformMaintenance';
import { AdminRunMaintenance } from './AdminRunMaintenance';
import { AdminSettings } from './AdminSettings';
import { AdminSsoGroups } from './AdminSsoGroups';
import { AdminSsoGroupsInside } from './AdminSsoGroupsInside';

export default function Demo() {
  return (
    <Router>
      <Helmet>
        <title>Ohana</title>
      </Helmet>
      <Routes>
        <Route
          path="/files"
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
        <Route 
          path="/home" 
          element={<VFSBrowser />} 
          />
        <Route 
          path="/" 
          element={<LoginPage />} 
          />
        <Route 
          path="/create_sso_key"
          element={<AdminCreateSsoKey/>} 
          />
        <Route 
          path="/ssogroups_inside" 
          element={<AdminSsoGroupsInside />} 
          />
        <Route 
          path="/key%20management" 
          element={<AdminKeyManagement />} 
          />
        <Route
          path="/maintenanceresults"
          element={<AdminMaintenanceResults />}
        />
        <Route
          path="/performmaintenance"
          element={<AdminPerformMaintenance />}
        />
        <Route 
          path="/maintenance" 
          element={<AdminMaintenanceSettings />} 
          />
        <Route 
          path="/runmaintenance" 
          element={<AdminRunMaintenance />} 
          />
        <Route 
          path="/dashboard"
          element={<AdminDashboard />} 
          />
        <Route 
          path="/sso" 
          element={<AdminSsoGroups />} 
          />
        <Route 
          path="/nodes" 
          element={<AdminNodes />} 
          />
        <Route 
          path="/maintenancelogs" 
          element={<AdminMaintenanceLogs />} 
          />
        <Route 
          path="/settings" 
          element={<AdminSettings />} 
          />
        <Route 
          path="/rotate%20key" 
          element={<AdminConfiguration />} 
          />
        <Route 
        path="/Admin_create_key" 
        element={<AdminCreateKey />} 
        />
      </Routes>
    </Router>
  );
}
