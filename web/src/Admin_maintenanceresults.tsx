import { useMantineTheme, Grid, Button, Card,Text, Paper, Center, Table } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import AppBase from "./components/AppBase";



function Admin_maintenanceresults() {
    const logs = [
        "Turning server one offline",
        "Cleaning server one",
        "Turning server one online",
        "Server one is back online",
        "Turning server two offline",
        "Turning server two online",
        "Maintenance is not completed.",
     ]

     const serverOnline = 641;
     const nodeOnline = 120;
     const theme = useMantineTheme();
 
    const rows = logs.map((items) => (
      <tr >
        <td width="80%" style={{ textAlign: "left", fontWeight: "400", fontSize: "16px", color: "black",border:"none" }}>{items}</td>
        
      </tr>
    ));
 return (
     
       
    <>
          <AppBase userType="admin" name='Alex Simmons' username='@alex' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'>
       
       
    
     <Center>
       <Grid style={{width:"100vh"}}> 
          

                   
     

         

 


       






            <Grid.Col span={12} style={{ marginTop: "" }}>
               <Card style={{ marginLeft: "10%", height: '500px', border: '1px solid ', marginTop: "5%", width: "60%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}
                  shadow="sm"
                  p="xl"

               >
                 
                  <Table captionSide="top"  >
          <caption style={{ textAlign: "center", fontWeight: "600", fontSize: "20px", color: "black" }}>  Maintenance Logs</caption>
       
          <tbody>{rows}</tbody>

        </Table>
                  <Button variant="default" color="dark" style={{ textAlign: "right", marginLeft: '70%' } }>
                     Export logs
                  </Button>
               </Card>

            </Grid.Col>
         </Grid>
         </Center>
       </AppBase>
    </>

   
 );
    }

    



export default Admin_maintenanceresults;