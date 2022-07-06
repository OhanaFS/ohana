//import AppBase from "./AppBase";
import { Helmet } from 'react-helmet';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AppBase from "./components/AppBase"; //switched to responsive base
import { VFSBrowser } from './components/userFiles';

import Admin_ssogroups from './Admin_ssogroups';

import LoginPage from './LoginPage';
import Admin_configuration from './Admin_configuration';
import Admin_create_key from './Admin_create_key';
import Admin_create_sso_key from './Admin_create_sso_key';
import Admin_key_management from './Admin_key_management';
import Admin_maintenancelogs from './Admin_maintenancelogs';
import Admin_maintenanceresults from './Admin_maintenanceresults';
import Admin_maintenancesettings from './Admin_maintenancesettings';
import Admin_nodes from './Admin_nodes';
import Admin_performmaintenance from './Admin_performmaintenance';
import Admin_runmaintenance from './Admin_runmaintenance';
import Admin_settings from './Admin_settings';
import Admin_ssogroups_inside from './Admin_ssogroups_inside';
import { Admin_statistics } from './Admin_statistics';


export default function Demo() {



  return (
    <Router>
      <Helmet>
        <title>Ohana</title>
      </Helmet>
      <Routes>
        <Route path='/files' element={<AppBase name='Alex' username='@alex' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'><></></AppBase>} />
        <Route path='/admin' element={<VFSBrowser />} />
        <Route path='/' element={<LoginPage />} />
        <Route path="/Admin_create_sso_key" element={<Admin_create_sso_key />} />
        <Route path="/Admin_ssogroups_inside" element={<Admin_ssogroups_inside />} />
        <Route path="/Admin_key_management" element={<Admin_key_management />} />
        <Route path="/Admin_maintenanceresults" element={<Admin_maintenanceresults />} />
        <Route path="/Admin_performmaintenance" element={<Admin_performmaintenance />} />
        <Route path="/Admin_maintenancesettings" element={<Admin_maintenancesettings />} />
        <Route path="/Admin_runmaintenance" element={<Admin_runmaintenance />} />
        <Route path="/Admin_statistics" element={<Admin_statistics />} />
        <Route path="/Admin_ssogroups" element={<Admin_ssogroups />} />
        <Route path="/Admin_nodes" element={<Admin_nodes />} />
        <Route path="/Admin_maintenancelogs" element={<Admin_maintenancelogs />} />
        <Route path="/Admin_settings" element={<Admin_settings />} />
        <Route path="/Admin_configuration" element={<Admin_configuration />} />
        <Route path="/Admin_create_key" element={<Admin_create_key />} />

      </Routes>
    </Router>
  );
}