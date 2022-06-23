

import { Grid, Textarea, Table, Checkbox, Button,Text, TextInput, Center } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import{
    BrowserRouter as Router,
    Link,
    Route,
    Routes


} from "react-router-dom";


function Admin_ssogroups_inside() {

    const data = [["Tom"],["Peter"],["Raymond"]]


 return (
     
       <>

<Center>
<Grid style={{width:"80vh"}}>
<Grid.Col span={12} style={{height:"500px", marginLeft:"2%",marginTop:"2%", border: '1px solid'}}>    




         
      
      <Table>
<tr>

<td>
<Text underline weight={700} style={{marginLeft:"1%",marginTop:"1%"}}> <h2>SSO Group: </h2>   </Text>

</td>




</tr>



</Table>
<Grid.Col span={12}>
<Table style={{height:"300px"}}>


<tr > 
    <td>
    <Text underline weight={700} style={{marginLeft:"1%",marginTop:"1%"}}> List of Users inside this group </Text>
    </td>


</tr>
 

{data.map((userlist, index) => {
        return (
          
          <tr>
            {userlist.map((user, sIndex) => {
              return  <td> {user} </td>;
            })}

         
          </tr>



        );
      })}

<tr>

    <td>

  

<Button  style={{ marginLeft: "", marginTop:"1%" }} component={Link} to="/Admin_create_key"  >Add User</Button>


    </td>
    <td>
    <Button  style={{ marginLeft: "", marginTop:"3%" }}  >Delete User</Button>

    </td>

</tr>
    


    

   

    </Table>
</Grid.Col>




      
      
      
      
     







      
      
      
      
      </Grid.Col>



      </Grid>
      </Center>
</>
   
 );
    }

    



export default Admin_ssogroups_inside;