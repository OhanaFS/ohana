import { useMantineTheme, Grid, Button, Card,Text, Paper, Center } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";



function Admin_maintenanceresults() {
    const logs = [
        "Turning server one offline",
        "Cleaning server one",
        "Turning server one online",
        "Server one is back online",
        "Turning server two offline",
        "Turning server two online",
        "Maintenance is completed.",
     ]

     const serverOnline = 641;
     const nodeOnline = 120;
     const theme = useMantineTheme();
 return (
     
       
    <>
       
    
     <Center>
       <Grid style={{width:"100vh"}}> 
            <Grid.Col span={6} style={{ height: "600px", marginTop: "1%" }}>

                   
       <Paper  style={{ marginLeft: "3%", height: '500px',  marginTop: "1%", width: 
       "95%",border: '1px solid', background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}>
       

       <Text underline weight={700} style={{marginLeft:"3%",marginTop:"3%"}}>    <h2>Maintenance </h2>  </Text>
       <Text  style={{marginLeft:"3%",marginTop:"3%"}}>        Maintenance is completed. {serverOnline}  </Text>
       <Text  style={{marginLeft:"3%",marginTop:"3%"}}>     - Total server online : {serverOnline}  </Text>
       <Text  style={{marginLeft:"3%",marginTop:"3%"}}>    - Total nodes online : {nodeOnline} </Text>
      

   

     

       </Paper>

         

 


            </Grid.Col>






            <Grid.Col span={6} style={{ marginTop: "" }}>
               <Card style={{ marginLeft: "10%", height: '500px', border: '1px solid ', marginTop: "5%", width: "90%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                  shadow="sm"
                  p="xl"

               >
                  <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                     <Text underline weight={700} style={{ marginLeft: "3%", marginTop: "3%" }}> <h2>Scheduled Maintenance Logs</h2>   </Text>

                  </Card.Section>
                  <div style={{ marginLeft: "3%" }}>

                     {logs.map(logs => <p>- {logs}</p>)}

                  </div>
                  <Button variant="light" color="blue" style={{ textAlign: "right", marginLeft: '70%' } }>
                     Export logs
                  </Button>
               </Card>

            </Grid.Col>
         </Grid>
         </Center>
    </>

   
 );
    }

    



export default Admin_maintenanceresults;