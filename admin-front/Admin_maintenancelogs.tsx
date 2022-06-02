

import { Card, Grid, Paper, Slider, useMantineTheme,Text, Button } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import { useScrollIntoView } from '@mantine/hooks';
import React from "react";

import{
  BrowserRouter as Router,
  Link,
  Route,
  Routes


} from "react-router-dom";
function Admin_maintenancelogs() {
    const theme = useMantineTheme();
    const {  scrollableRef } = useScrollIntoView();

    const maintenanceLogs = [[ 
    "Maintenance date: 19/2/22 23:59",
    "Total servers online : 3000",
    "Total nodes online : 15200",
    "Total time estimated time : 2:23:59"],
    [ 
        "Maintenance date: 19/2/22 23:59",
        "Total servers online : 3000",
        "Total nodes online : 15200",
        "Total time estimated time : 2:23:59"],
        [ 
            "Maintenance date: 19/2/22 23:59",
            "Total servers online : 3000",
            "Total nodes online : 15200",
            "Total time estimated time : 2:23:59"],
            [ 
              "Maintenance date: 19/2/22 23:59",
              "Total servers online : 3000",
              "Total nodes online : 15200",
              "Total time estimated time : 2:23:59"],
              [ 
                "Maintenance date: 19/2/22 23:59",
                "Total servers online : 3000",
                "Total nodes online : 15200",
                "Total time estimated time : 2:23:59"]
        ]

 return (
        
       <>

<Grid>
      <Grid.Col span={4}>

         
      
      
       <Paper ref={scrollableRef} style={{ marginLeft: "3%", height: '500px',  marginTop: "1%", width: 
       "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[1] ,overflowY: 'scroll',  flex: 1 }}>
       

       <Text underline weight={700} style={{marginLeft:"3%",marginTop:"3%"}}>  Run Scheduled Maintenance  </Text>

        <br></br>
       <Text underline weight={700} style={{marginLeft:"3%",marginTop:"3%"}}>  Past Maintenance:  </Text>
       {maintenanceLogs.map((items, index) => {
        return (
          <ol>
            {items.map((subItems, sIndex) => {
              return <ol> {subItems} </ol>;
            })}
          </ol>
        );
      })}
       


     

       </Paper>

       <div style={{ display: "flex" }}>
     
      <Button  style={{ marginLeft: "auto", marginTop:"3%" }} component={Link} to="/Admin_runmaintenance"  >Perform Maintainance</Button>

      </div>



      </Grid.Col>
    
    </Grid>
</>
   
 );
    }

    



export default Admin_maintenancelogs;