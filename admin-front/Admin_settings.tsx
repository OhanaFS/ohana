

import { Grid, Checkbox, Button,Text, Center, Table, Card, useMantineTheme } from "@mantine/core";
import { Console } from "console";
import { useState } from "react";
import { Link } from "react-router-dom";
import Admin_navigation from "./Admin_navigation";

import { useForceUpdate, randomId } from '@mantine/hooks';

function Admin_settings () {
  const theme = useMantineTheme();
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


    const [disable, setDisable] = useState(true);
    let [count,setCount]=useState(0);

 
   

 return (
     
       
    <>

<Center>


  

<Card  style={{ marginLeft: "0%", height: '65vh', border: '1px solid ', marginTop: "3%", width: "60%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white[0] }}
                 shadow="sm"
                 p="xl"            >
<Table striped  verticalSpacing="md"  >
<caption style={{fontWeight:"600",fontSize:"22px",color:"black"}}> <span style={{textAlign: "center"}}>Notification Settings</span> 



</caption>

<thead> 


</thead>
  <tbody style={{}}>


<tr>
<td width="50%" style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}}>  <Text style={{ }}> Allow Cluster health alerts   </Text></td>
<td width="50%" > <Center style={{marginLeft:"80%"}}>  <Checkbox size="md" style={{}} checked={checked0} onChange={(event) => [setChecked(event.currentTarget.checked),setDisable(event.currentTarget.checked)]}/> </Center>   </td>
</tr>

<tr>
<td width="50%" style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}}>  <Text style={{ }}>    Allow server offline alerts </Text></td>
<td width="50%"> <Center style={{marginLeft:"80%"}}> <Checkbox size="md" id="1" style={{}} checked={checked1} onChange={(event) => [setChecked1(event.currentTarget.checked),setDisable(event.currentTarget.checked)]}/> </Center></td>
</tr>

<tr>
<td width="50%" style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}}>  <Text style={{ }}>     Allow supicious action alerts </Text> </td>
<td width="50%"> <Center style={{marginLeft:"80%"}}> <Checkbox size="md" id="1" style={{}}  checked={checked2} onChange={(event) => [setChecked2(event.currentTarget.checked),setDisable(event.currentTarget.checked)]}/> </Center> </td>
</tr>



<tr>
<td width="50%" style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}}>  <span style={{ }}>     Allow server full alert </span> </td>
<td width="50%"> <Center style={{marginLeft:"80%"}}>  <Checkbox size="md" id="1" style={{}} checked={checked3} onChange={(event) => [setChecked3(event.currentTarget.checked),setDisable(event.currentTarget.checked)]}/> </Center></td>
</tr>

<tr >
<td width="50%" style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}} >  <span style={{ }}>      Allow supicious file alerts </span></td>
<td width="50%" style={{}} ><Center style={{marginLeft:"80%"}}> <Checkbox size="md"id="1" style={{}}  checked={checked4} onChange={(event) => [setChecked4(event.currentTarget.checked),setDisable(event.currentTarget.checked)]}/></Center> </td>
</tr>

<tr> 
  <td style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}}>  Backup encryption key    <Text  weight={700}>Current Location:{currentLocation} </Text>   </td>
  <td>  <Button style={{float:"right"}}variant="default" color="dark" size="md"> Backup</Button></td>
</tr>

<tr>
  <td style={{textAlign: "left",fontWeight:"400",fontSize:"18px",color:"black"}}>  Change the redundancy level of the files <Text weight={700}>Current redundancy level:{redundancy} </Text>  </td>
  <td >  <Button style={{float:"right"}}variant="default" color="dark" size="md"> Change</Button></td>
</tr>




</tbody>
<tfoot>
  <td></td>
<td>  
<Button disabled={disable} style={{marginTop:"2%",float:"right"}}variant="default" color="dark" size="md"  >
  Save Pending Changes
</Button> 


</td>
</tfoot>

</Table>



  
  

</Card>


</Center>
    
    </>

   
 );
    }




export default Admin_settings;