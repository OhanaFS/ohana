

import { Grid, Button, Paper, useMantineTheme,Text, GroupedTransition, Textarea, Checkbox, Table, Center } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";



function Admin_configuration() {
    const theme = useMantineTheme();
 return (
     
       <>
       <Center>
<Grid style={{width:"100vh"}}>
<Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"80%", border: '1px solid'}}>   




         
      
      


<Text underline weight={700} style={{marginLeft:"1%",marginTop:"3%"}}> <h2>Rotate Key</h2>   </Text>

<Grid.Col span={12}>

<Text style={{marginTop:"3%"}}>  Specify the file/directory location and the system will auto rotate the key  </Text>

</Grid.Col>
<Grid.Col span={12}>


<Textarea
      label="File location"
      radius="xs"
      size="md"
    />

</Grid.Col>

<Grid.Col span={12} >  
<Table>
<tr>
    <td width="15%">Master Key :</td>
    <td width="90%">   <Checkbox style={{}}> Edit</Checkbox></td>
    </tr>
</Table>

  
      </Grid.Col>
   
      
      
      
      
     


<div style={{ display: "flex" }}>

<Button  style={{ marginLeft: "auto", marginTop:"3%" }}  >Rotate Key</Button>

</div>




      
      
      
      
      </Grid.Col>



      </Grid>
      </Center>
</>
   
 );
    }

    



export default Admin_configuration;