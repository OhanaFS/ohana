

import { Grid, Textarea, Table, Checkbox, Button,Text, TextInput } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import{
    BrowserRouter as Router,
    Link,
    Route,
    Routes


} from "react-router-dom";

function Admin_create_key() {

 return (
     
       <>

<Grid>
<Grid.Col span={12} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>    




         
      
      


<Text underline weight={700} style={{marginLeft:"1%",marginTop:"3%"}}> <h2>Create API Key</h2>   </Text>


<Grid.Col span={12}>
<Table>
    <tr>
<td>
<TextInput
      label="API Key"
      radius="xs"
      size="md"
      required
    />

    </td>
    <td>

    <Button >Generate</Button>
    </td>

    </tr>

    </Table>
</Grid.Col>

<Grid.Col span={12}>

<TextInput
      label="Description"
      radius="xs"
      size="md"
      required
    />

</Grid.Col>


      
      
      
      
     


<div style={{ display: "flex" }}>

<Button  style={{ marginLeft: "auto", marginTop:"3%" }} component={Link} to="/Admin_key_management"  >Create Key</Button>

</div>




      
      
      
      
      </Grid.Col>



      </Grid>

</>
   
 );
    }

    



export default Admin_create_key;