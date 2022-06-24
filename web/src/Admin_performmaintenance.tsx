

import { Grid, Button, Text, Card, useMantineTheme, Image, Center } from "@mantine/core";
import { Link } from "react-router-dom";
import Admin_navigation from "./Admin_navigation";

import img2 from '../src/images/3.png';

function Admin_performmaintenance() {
   const logs = [
      "Turning server one offline",
      "Cleaning server one",
      "Turning server one online",
      "Server one is back online",
      "Turning server two offline"
   ]
   const theme = useMantineTheme();
   return (


      <>
         <Admin_navigation>
            <Center>
               <Grid style={{ width: "100vh" }}>
                  <Grid.Col span={6} style={{ height: "600px", marginTop: "1%" }}>

                     <Image style={{ textAlign: 'left', marginLeft: "1%" }}

                        src={img2} />

                     <Button color="blue" style={{ marginLeft: '20%', width: '20%' }}>
                        Pause
                     </Button>

                     <Button color="blue" component={Link} to="/Admin_maintenanceresults" style={{ marginLeft: '25%', width: '20%', }}>
                        Stop
                     </Button>


                  </Grid.Col>






                  <Grid.Col span={6} style={{ marginTop: "1%" }}>
                     <Card style={{ marginLeft: "10%", height: '500px', border: '1px solid ', marginTop: "5%", width: "60%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                        shadow="sm"
                        p="xl"

                     >
                        <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                           <Text underline weight={700} style={{ marginLeft: "3%", marginTop: "3%" }}> <h2>Run Scheduled Maintenance</h2>   </Text>

                        </Card.Section>
                        <div style={{ marginLeft: "3%" }}>

                           {logs.map(logs => <p>- {logs}</p>)}

                        </div>

                     </Card>

                  </Grid.Col>
               </Grid>
            </Center>
         </Admin_navigation>
      </>


   );
}





export default Admin_performmaintenance;


