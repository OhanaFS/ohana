




import { Grid, Table, Button, Text, Center } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import {
   Link,


} from "react-router-dom";
import { useState } from "react";


function Admin_ssogroups() {

   // take from database.
   let [CurrentSSOGroups, setValue] = useState([
      "Finance",
      "HR",
      "IT",

   ]);

   function add() {

      setValue(prevValue => prevValue.concat('sds')
      );

   }

   return (

      <>
         <Admin_navigation>
            <Center>
               <Grid style={{ width: "120vh" }} >






                  <Grid.Col >

                     <Table >

                        <tr>
                           <td>
                              <Text underline weight={700} style={{ marginLeft: "2%", marginTop: "30%" }}> <h2>Current SSO Groups</h2>   </Text>
                           </td>
                           <td>
                              <Button style={{ marginLeft: "Auto", marginTop: "2%" }} component={Link} to="/Admin_create_sso_key"  >Create SSO Groups</Button>
                           </td>
                        </tr>


                     </Table >

                     <div style={{}}>

                        {CurrentSSOGroups.map(CurrentSSOGroups =>

                           <Grid.Col span={10} style={{}}>   <Button style={{ width: '100%', marginLeft: "2%", marginTop: "1%", maxWidth: "100%", border: '1px solid' }} component={Link} to="/Admin_ssogroups_inside"> {CurrentSSOGroups} </Button></Grid.Col>



                        )}

                     </div>




                  </Grid.Col>







                  <Grid.Col span={12}>








                     <Button onClick={add}>ADD static SSO group</Button>








                  </Grid.Col>



               </Grid>
            </Center>
         </Admin_navigation>
      </>

   );
}





export default Admin_ssogroups;

