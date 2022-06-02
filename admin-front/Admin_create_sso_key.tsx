




import { Grid, Textarea, Table, Checkbox, Button,Text, TextInput,Paper, useMantineTheme } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import{
    BrowserRouter as Router,
    Link,
    Route,
    Routes


} from "react-router-dom";
import { useScrollIntoView } from "@mantine/hooks";

function Admin_create_sso_key() {

   const CurrentSSOGroups = [
      "Finance",
      "HR",
      "IT",
      "memem"
   ]
  
   const theme = useMantineTheme();
 return (
     
       <>

<Grid>
<Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   




         
      
      


<Text underline weight={700} style={{marginLeft:"1%",marginTop:"3%"}}> <h2>Create SSO </h2>   </Text>

<Grid.Col span={12}>


</Grid.Col>
<Grid.Col span={12}>


<Textarea
      label="SSO Group Name:"
      radius="xs"
      size="md"
    />

</Grid.Col>


   
      
      
      
      
     


<div style={{ display: "flex" }}>

<Button  style={{ marginLeft: "auto", marginTop:"3%" }}  component={Link} to="/Admin_ssogroups" >Create</Button>

</div>




      
      
      
      
      </Grid.Col>



      </Grid>

</>
   
 );
    }

    



export default Admin_create_sso_key;

