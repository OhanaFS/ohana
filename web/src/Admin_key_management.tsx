

import { Grid, Table, Button, Text, Center, Card, Checkbox, useMantineTheme, ScrollArea, Modal, TextInput } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import {
  BrowserRouter as Router,
  Link,


} from "react-router-dom";
import { Settings } from "tabler-icons-react";
import { useState } from "react";
import Admin_console from "./Admin_console";
import AppBase from "./components/AppBase";


function Admin_key_management() {

  let [data, setValue] = useState([
    "128c1d5d-2359-4ba1-8739-2cd30d694d67",
    "128c1d5d-2359-4ba1-8739-2cd30d69sds67",
  
    
 ]);
  
  
 
  return (
    <>
    
    <AppBase userType="admin" name='Alex Simmons' username='@alex' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'>
       

   

    
<Admin_console consoleWidth={80} consoleHeight={80} groupList={data} addObjectLabel = "Key" deleteObjectLabel="Key" tableHeader={["Key ID"]} tableBody={[]} caption="API Key Management Console" pointerEvents={false} ></Admin_console>

</AppBase>

    </>
  );
}





export default Admin_key_management;