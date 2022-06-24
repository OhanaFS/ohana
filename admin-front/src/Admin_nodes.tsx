


import { Grid, Textarea, Table, Checkbox, Button,Text, TextInput, Center } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import{
    BrowserRouter as Router,
    Link,
    Route,
    Routes


} from "react-router-dom";


function Admin_nodes() {

    const data = [["192.168.1.1"],["192.168.1.2"]]


 return (
     
       <>
       <Admin_navigation>
<Center>
<Grid style={{width:"100vh"}}>
<Grid.Col span={12} style={{height:"500px", marginLeft:"2%",marginTop:"2%", border: '1px solid'}}>    




         
      
      <Table>
<tr>

<td>
<Text underline weight={700} style={{marginLeft:"1%",marginTop:"1%"}}> <h2>Node Mangement</h2>   </Text>

</td>




</tr>



</Table>
<Grid.Col span={12}>
<Table style={{height:"300px"}}>


<tr > 
    <td>
    <Text underline weight={700} style={{marginLeft:"1%",marginTop:"1%"}}>IP address of the node </Text>
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

  

<Button  style={{ marginLeft: "", marginTop:"1%" }} component={Link} to="/"  >Add Node</Button>


    </td>
    <td>
    <Button  style={{ marginLeft: "", marginTop:"3%" }}  >Disconnect Node</Button>

    </td>

</tr>
    


    

   

    </Table>
</Grid.Col>




      
      
      
      
     







      
      
      
      
      </Grid.Col>



      </Grid>
      </Center>
      </Admin_navigation>
</>
   
 );
    }

    



export default Admin_nodes;