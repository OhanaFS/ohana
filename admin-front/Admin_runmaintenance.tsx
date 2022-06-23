

import {  Grid,  useMantineTheme,Text, SimpleGrid, Checkbox, Button, Center } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import{
  BrowserRouter as Router,
  Link,
  Route,
  Routes


} from "react-router-dom";
import { useState } from "react";

function Admin_runmaintenance() {
    const theme = useMantineTheme();
    let MaintenanceSettings = [
      {name: 'CrawlPermissions', setting: "true"},
      {name: 'PurgeOrphanedFile', setting: "true"},
      {name: 'PurgeUser', setting: "false"},
      {name: 'CrawlReplicas', setting: "true"},
      {name: 'QuickCheck', setting: "true"},
      {name: 'FullCheck', setting: "false"},
      {name: 'DBCheck', setting: "true"}
    ];
    const [checked0, setChecked] = useState(() => {
      if (MaintenanceSettings[0].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked1, setChecked1] = useState(() => {
      if (MaintenanceSettings[1].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked2, setChecked2] = useState(() => {
      if (MaintenanceSettings[2].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked3, setChecked3] = useState(() => {
      if (MaintenanceSettings[3].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked4, setChecked4] = useState(() => {
      if (MaintenanceSettings[4].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked5, setChecked5] = useState(() => {
      if (MaintenanceSettings[5].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked6, setChecked6] = useState(() => {
      if (MaintenanceSettings[6].setting=="true") {
        return true;
      }
  
      return false;
    });
   
    
 return (
     
       
    <>
  <Center>
    <Grid style={{width:"100vh"}}> 
      <Grid.Col span={4} style={{ marginLeft:"2%" }}>  <Text underline weight={700} >  Run Scheduled Maintenance  </Text></Grid.Col>
      <Grid.Col span={2} style={{textAlign:'right'}} >   <Button radius="md" size="xs"  component={Link} to="/Admin_maintenancesettings">
      Settings
    </Button></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"0%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>  Crawl the list of files to remove permissions from expired users   </Text></Grid.Col>
      <Grid.Col span={2} ><Checkbox size="md" checked={checked0} onChange={(event) => setChecked(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Purging orphaned files and shards </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  checked={checked1} onChange={(event) => setChecked1(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>    Purge a user and their files </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  checked={checked2} onChange={(event) => setChecked2(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>    Crawl all of the files to make sure it has full replicas</Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  checked={checked3} onChange={(event) => setChecked3(event.currentTarget.checked)}/></Grid.Col>
     
      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>     Quick File Check (Only checks current versions of files to see if it’s fine and is not corrupted) </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md"id="1"  checked={checked4} onChange={(event) => setChecked4(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>     Full File Check (Checks all fragments to ensure that it’s not corrupted) </Text> </Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md"id="1"  checked={checked5} onChange={(event) => setChecked5(event.currentTarget.checked)}/></Grid.Col>
      
      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   DB integrity Check </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  checked={checked6} onChange={(event) => setChecked6(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={12} style={{textAlign:'right' ,marginLeft:"2%",marginTop:"2%",maxWidth:"50%"}}>    
       <Button radius="md" size="xs"  component={Link} to="/Admin_performmaintenance">
      Run Maintenance
    </Button> </Grid.Col> 
    
    </Grid>

    </Center>
 
 
   
</>

   
 );
    }

    



export default Admin_runmaintenance;