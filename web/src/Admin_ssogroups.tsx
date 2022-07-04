



import { Grid, Table, Button, Text, Center, Checkbox, useMantineTheme, Card, ScrollArea } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import {
   Link,


} from "react-router-dom";
import { useState } from "react";
import { ResponsiveContainer } from "recharts";
import Admin_console from "./Admin_console";


function Admin_ssogroups() {



  
  /*
    groupList: Array<string>;
    addObjectLabel:string;
    deleteObjectLabel:string;
    tableHeader: string;
    caption: string;
    pointerEvents : boolean; 
    conso
  */
 
  const SSOGroupList = ["Hr","asd"];
 
 

   return (
 
      <>
     
         <Admin_navigation>
       
      <Admin_console consoleWidth={80} consoleHeight={60} groupList={SSOGroupList} addObjectLabel = "Group" deleteObjectLabel="Group" tableHeader={["Current SSO Groups"]} tableBody={[]} caption="SSO Group Management Console" pointerEvents={true} ></Admin_console>
      
         </Admin_navigation>
      </>

   );
}





export default Admin_ssogroups;


