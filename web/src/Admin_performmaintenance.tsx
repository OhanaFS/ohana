

import { Grid, Button, Text, Card, useMantineTheme, Image, Center } from "@mantine/core";
import { Link } from "react-router-dom";
import Admin_navigation from "./Admin_navigation";

import img2 from '../src/images/3.png';
import { ResponsiveContainer, PieChart, Pie, Legend, Tooltip } from "recharts";
import { Cell } from "tabler-icons-react";

function Admin_performmaintenance() {
   const logs = [
      "Turning server one offline",
      "Cleaning server one",
      "Turning server one online",
      "Server one is back online",
      "Turning server two offline"
   ]

   const ClusterHealthChartData = [

      {
         "name": "No of Healthy Nodes",
         "value": 2000000
       },
       {
         "name": "No of Unhealthy Nodes",
         "value": 0
       },
   ]
   const RADIAN = Math.PI / 180;



const barColors = ["#1f77b4", "#ff0000"]
   const theme = useMantineTheme();
   return (


      <>
         <Admin_navigation>
            <Center>
               <Grid style={{ width: "40vh" }}>
                  <Grid.Col span={12} style={{  }}>
                     <Card  style={{ marginLeft: "0%", height: '60vh', border: '1px solid ', marginTop: "3%", width: "160%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}
      shadow="sm"
      p="xl">
                  <Text style={{textAlign:"center"}}><h2>Maintenance Progress</h2></Text>
                  <ResponsiveContainer width={500} height={300} >
                            <PieChart width={500} height={300}>
                           
                                <Pie
                                    data={ClusterHealthChartData}
                                    cx="50%"
                                    cy="50%"
                                    labelLine={false}
                                 
                                    outerRadius={100}
                                    fill="#8884d8"
                                    dataKey="value"
                                >
                                    {ClusterHealthChartData.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={barColors[index % barColors.length]} />
                                    ))}
                                </Pie>
                            
                            </PieChart>
                         
                        </ResponsiveContainer>
                     <Button variant="default" color="dark"  style={{ marginLeft: '20%', width: '20%' }}>
                        Pause
                     </Button>
                
                     <Button  variant="default" color="dark"  component={Link} to="/Admin_maintenanceresults" style={{ marginLeft: '25%', width: '20%', }}>
                        Stop
                     </Button>

                     </Card>
                  </Grid.Col>






               
               </Grid>
            </Center>
         </Admin_navigation>
      </>


   );
}





export default Admin_performmaintenance;


