

import { Grid, Checkbox, Button,Text, Center } from "@mantine/core";
import { Console } from "console";
import { useState } from "react";
import { Link } from "react-router-dom";
import Admin_navigation from "./Admin_navigation";

import { useForceUpdate, randomId } from '@mantine/hooks';

function Admin_settings () {

   //retrieve from database
   let ConfigurationSettings = [
      {name: 'clusterAlerts', setting: "true"},
      {name: 'sActionAlerts', setting: "true"},
      {name: 'supiciousAlerts', setting: "false"},
      {name: 'serverAlerts', setting: "true"},
      {name: 'sFileAlerts', setting: "true"},
      {name: 'BackupLocation', setting: "C:\\\Users\\\admin"},
      {name: 'redundancy', setting: "Low"}
    ];


    const [checked0, setChecked] = useState(() => {
      if (ConfigurationSettings[0].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked1, setChecked1] = useState(() => {
      if (ConfigurationSettings[1].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked2, setChecked2] = useState(() => {
      if (ConfigurationSettings[2].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked3, setChecked3] = useState(() => {
      if (ConfigurationSettings[3].setting=="true") {
        return true;
      }
  
      return false;
    });
    const [checked4, setChecked4] = useState(() => {
      if (ConfigurationSettings[4].setting=="true") {
        return true;
      }
  
      return false;
    });

    let currentLocation = ConfigurationSettings[5].setting;
    let redundancy= ConfigurationSettings[6].setting;



    


   

 return (
     
       
    <>
      <Center >
       <Grid  style={{width:"100vh"}}>
  
     

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"80%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>  Allow Cluster health alerts   </Text></Grid.Col>
      <Grid.Col span={2}  style={{ marginTop:"5%"}}>
         
         
         
      <Checkbox size="md" id="clusterAlerts"    checked={checked0} onChange={(event) => setChecked(event.currentTarget.checked)}/>
         
         
         
         </Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"80%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>  Allow server offline alerts </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"   checked={checked1} onChange={(event) => setChecked1(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"80%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Allow supicious action alerts </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  checked={checked2} onChange={(event) => setChecked2(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"80%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Allow server full alert </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"   checked={checked3} onChange={(event) => setChecked3(event.currentTarget.checked)} /></Grid.Col>
     
      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"80%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>     Allow supicious file alerts  </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"   checked={checked4} onChange={(event) => setChecked4(event.currentTarget.checked)}/></Grid.Col>

      <Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"90%", border: '1px solid'}} >   <div style={{ marginLeft:"1%"}}>    Backup encryption key  
      
      <Button style={{marginLeft:"70%",height:"30px"}}> Backup</Button>
      <Text weight={700}>Current Location:{currentLocation} </Text>
      </div>
       </Grid.Col>
      
      
    


      <Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"90%", border: '1px solid'}}>    <div style={{ marginLeft:"1%"}}>  Change the redundancy level of the files 
     
      <Button style={{marginLeft:"54%",height:"30px"}}> Change</Button>
      <Text weight={700}>Current redundancy level:{redundancy} </Text>
      </div>
      
      
      
      
      </Grid.Col>
     

      <Grid.Col span={12} style={{textAlign:'right' ,marginLeft:"2%",marginTop:"2%",maxWidth:"90%"}}>     <Button radius="md" size="xs"  >
      Save Pending Changes
        
    </Button > </Grid.Col> 

    </Grid>
    </Center>
    </>

   
 );
    }




export default Admin_settings;