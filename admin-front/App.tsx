

import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Admin_maintenance from './Admin_maintenancelogs';
import Admin_maintenanceSettings from './Admin_runmaintenance';

import img from '../src/images/1.png';

import Admin_navigation from "./Admin_navigation";
import Admin_ssogroups from './Admin_ssogroups';
import Admin_statistics from './Admin_statistics';
import LoginPage from './LoginPage';
import Admin_configuration from './Admin_configuration';
import { useState } from 'react';
import { Box, BackgroundImage,useMantineTheme, Center, Title, TextInput, Button } from '@mantine/core';
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


export default function Demo() {

  const [Login,setIsToggled]= useState(false);
  const theme = useMantineTheme();
  return (
    <>
    <div> 

    {Login  && <Admin_navigation  />  
    
    }


    {!Login && <LoginPage setIsToggled={setIsToggled} />

   
    }
    
   
</div>
 
    
    
 </>

   



    
  





  );


}