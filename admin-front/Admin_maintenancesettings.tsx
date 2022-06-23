import { Grid, Button, Checkbox,Text, Center } from "@mantine/core";
import { useListState } from "@mantine/hooks";

import React, { useState } from "react";

import{
   BrowserRouter as Router,
   Link,
   Route,
   Routes
 
 
 } from "react-router-dom";

function Admin_maintenancesettings() {
   
   



    
 return (
     
       
    <>  
        <Center>
        <Grid style={{width:"100vh"}}> 
      <Grid.Col span={12} style={{ marginLeft:"2%" }}>  <Text underline weight={700} >  Maintenance Settings </Text></Grid.Col>
   

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"0%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>  Crawl the list of files to remove permissions from expired users   </Text></Grid.Col>
      <Grid.Col span={2} >
         
         
         
      <Checkbox
     
    />
         
         
         
         </Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Purging orphaned files and shards </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"  /></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>    Purge a user and their files </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>    Crawl all of the files to make sure it has full replicas</Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>
     
      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>     Quick File Check (Only checks current versions of files to see if it’s fine and is not corrupted) </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>     Full File Check (Checks all fragments to ensure that it’s not corrupted) </Text> </Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1" /></Grid.Col>
      
      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   DB integrity Check </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}><Checkbox size="md" id="1"/></Grid.Col>

      <Grid.Col span={10} style={{ marginLeft:"2%",marginTop:"2%",maxWidth:"50%", border: '1px solid'}}>   <Text style={{ marginLeft:"1%"}}>   Use Default Settings </Text></Grid.Col>
      <Grid.Col span={2} style={{ marginTop:"2%"}}>
         
         
         
         <Checkbox 
     
      
        
      
       
         /></Grid.Col>

      <Grid.Col span={12} style={{textAlign:'right' ,marginLeft:"2%",marginTop:"2%",maxWidth:"50%"}}>     <Button radius="md" size="xs"  component={Link} to="/Admin_runmaintenance">
      Apply Settings
        
    </Button > </Grid.Col> 

    </Grid>

   
    
    
</Center>
    </>

   
 );
    }

    



export default Admin_maintenancesettings;

