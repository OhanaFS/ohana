




import { Grid, Textarea, Table, Checkbox, Button,Text, TextInput,Paper, useMantineTheme } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import{
    BrowserRouter as Router,
    Link,
    Route,
    Routes


} from "react-router-dom";
import { useScrollIntoView } from "@mantine/hooks";

function Admin_ssogroups() {
   // take from database.
   const CurrentSSOGroups = [
      "Finance",
      "HR",
      "IT",
      
   ]
  
   const theme = useMantineTheme();
 return (
     
       <>

<Grid >




         
      
   <Grid.Col > 
  
<Table >

<tr>
   <td>
<Text underline weight={700} style={{marginLeft:"2%",marginTop:"0%"}}> <h2>Current SSO Groups</h2>   </Text>
</td>
<td>
<Button  style={{ marginLeft: "Auto", marginTop:"2%" }} component={Link} to="/Admin_create_sso_key"  >Create SSO Groups</Button>
</td>
</tr>


</Table >

<div style={{  }}>

{CurrentSSOGroups.map(CurrentSSOGroups => 

 <Grid.Col span={10} style={{ }}>   <Button style={{ width:'100%',marginLeft:"2%",marginTop:"1%",maxWidth:"100%", border: '1px solid'}}  component={Link} to="/Admin_ssogroups_inside"> {CurrentSSOGroups} </Button></Grid.Col>
   
   
   
   )}

</div>
          
    


</Grid.Col>







<Grid.Col span={12}>





      
      
      
      
     



      
      
      
      
      </Grid.Col>



      </Grid>

</>
   
 );
    }

    



export default Admin_ssogroups;

