import { AppShell, Aside, Burger, Button, Footer, Header, MediaQuery, Navbar,Text, Title, useMantineTheme,ScrollArea } from "@mantine/core";
import { useState  } from "react";

import Admin_statistics from "./Admin_statistics";
import Admin_ssogroups from "./Admin_ssogroups";
import Admin_create_key from "./Admin_create_key";
import Admin_maintenancelogs from "./Admin_maintenancelogs";
import Admin_nodes from "./Admin_nodes";
import Admin_settings from "./Admin_settings";
import Admin_configuration from "./Admin_configuration";
import Admin_runmaintenance from "./Admin_runmaintenance";


import{
    BrowserRouter as Router,
    Link,
    Route,
    Routes


} from "react-router-dom";

import Admin_maintenancesettings from "./Admin_maintenancesettings";
import Admin_maintenanceresults from "./Admin_maintenanceresults";
import Admin_performmaintenance from "./Admin_performmaintenance";
import Admin_key_management from "./Admin_key_management";
import Admin_create_sso_key from "./Admin_create_sso_key";
import Admin_ssogroups_inside from "./Admin_ssogroups_inside";
import LoginPage from "./LoginPage";


function Admin_navigation() {
    const theme = useMantineTheme();
    const [opened, setOpened] = useState(false);
    return (
        <Router>
        <AppShell   
        styles={{
          main: {
            background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
          },
        }}
        navbarOffsetBreakpoint="sm"
        asideOffsetBreakpoint="sm"
        fixed
        navbar={
          <Navbar p="md" hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 200, lg: 200 }}>
            <Navbar.Section>{}</Navbar.Section>
            
            <Navbar.Section grow component={ScrollArea} mx="-xs" px="xs">
        { <div style={{ display:"flex",flexDirection:"column"}}>
            
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/Admin_statistics" >
                View Statistic 
                </Text>
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/Admin_ssogroups" >
                View SSO Groups
            </Text>
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/Admin_nodes" >
                Manage Nodes
            </Text>
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/Admin_maintenancelogs" >
                Run Maintenance
            </Text>
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/Admin_settings" >
              Settings
            </Text>
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}} component={Link} variant="link" to="/Admin_configuration" >
              Key Configuration
            </Text>
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/Admin_key_management" >
              API Key
            </Text>
     



            </div>}
      </Navbar.Section>

     {/*this is the footer for navbar*/}
      <Navbar.Section>{}</Navbar.Section> 
          </Navbar>
        }
     
        footer={
          <Footer height={60} p="md">
            <Text style={{height:"50px",border: '1px solid ', backgroundColor:theme.colors.gray[4]}}component={Link} variant="link" to="/LoginPage" >
             Logout
            </Text>
          </Footer>
        }
        header={
          <Header height={70} p="md">
            <div style={{ display: 'flex', alignItems: 'center', height: '100%' }}>
              <MediaQuery largerThan="sm" styles={{ display: 'none' }}>
                <Burger
                  opened={opened}
                  onClick={() => setOpened((o) => !o)}
                  size="sm"
                  color={theme.colors.gray[6]}
                  mr="xl"
                />
              </MediaQuery>
  
              <Title order={2}>Ohana </Title>
            </div>
          </Header>
   


        }
      >
       <Routes>
       
       <Route path="/Admin_create_sso_key" element={<Admin_create_sso_key/>} />
       <Route path="/Admin_ssogroups_inside" element={<Admin_ssogroups_inside/>} />
       <Route path="/Admin_key_management" element={<Admin_key_management/>} />
       <Route path="/Admin_maintenanceresults" element={<Admin_maintenanceresults/>} />
       <Route path="/Admin_performmaintenance" element={<Admin_performmaintenance/>} />
       <Route path="/Admin_maintenancesettings" element={<Admin_maintenancesettings/>} />
       <Route path="/Admin_runmaintenance" element={<Admin_runmaintenance/>} />
       <Route path="/Admin_statistics" element={<Admin_statistics/>} />
       <Route path="/Admin_ssogroups" element={<Admin_ssogroups/>} />
       <Route path="/Admin_nodes" element={<Admin_nodes/>} />
       <Route path="/Admin_maintenancelogs" element={<Admin_maintenancelogs/>} />
       <Route path="/Admin_settings" element={<Admin_settings/>} />
       <Route path="/Admin_configuration" element={<Admin_configuration/>} />
       <Route path="/Admin_create_key" element={<Admin_create_key/>} />
       <Route path="/LoginPage" element={<LoginPage/>} />

       
       </Routes>
      </AppShell>
      </Router>
    );
       }
   
       
   
   
   
   export default Admin_navigation;