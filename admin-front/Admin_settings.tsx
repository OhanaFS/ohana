

import { Grid, Checkbox, Button,Text } from "@mantine/core";
import { Link } from "react-router-dom";
import Admin_navigation from "./Admin_navigation";



function Admin_settings() {

 return (
     
       
    <>
       
    
       <Grid>
  
   

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"1%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>  Allow Cluster health alerts   </Text></Grid.Col>
      <Grid.Col span={2}  style={{ marginTop:"1%"}}>
         
         
         
      <Checkbox
     
    />
         
         
         
         </Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>  Allow server offline alerts </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  /></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Allow supicious action alerts </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Allow server full alert </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>
     
      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>     Allow supicious file alerts  </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>

      <Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <div style={{ marginLeft:"1%"}}>    Backup encryption key 
      
      <Button style={{marginLeft:"70%",height:"30px"}}> Backup</Button>
      </div>
       </Grid.Col>
      
      
    


      <Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>    <div style={{ marginLeft:"1%"}}>  Change the redundancy level of the files
      
      <Button style={{marginLeft:"55%",height:"30px"}}> Edit</Button>
      </div>
      
      
      
      
      </Grid.Col>
     

      <Grid.Col span={12} style={{textAlign:'right' ,marginLeft:"2%",marginTop:"2%",maxWidth:"50%"}}>     <Button radius="md" size="xs"  >
      Save Pending Changes
        
    </Button > </Grid.Col> 

    </Grid>
    </>

   
 );
    }

    



export default Admin_settings;