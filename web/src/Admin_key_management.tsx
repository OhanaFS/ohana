

import { Grid, Table, Button, Text, Center, Card, Checkbox, useMantineTheme, ScrollArea, Modal, TextInput } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import {
  BrowserRouter as Router,
  Link,


} from "react-router-dom";
import { Settings } from "tabler-icons-react";
import { useState } from "react";
import Admin_console from "./Admin_console";


function Admin_key_management() {

  let [data, setValue] = useState([
    "128c1d5d-2359-4ba1-8739-2cd30d694d67",
    "128c1d5d-2359-4ba1-8739-2cd30d69sds67",
  
    
 ]);
  
  
 
  return (
    <>
    
    <Admin_navigation>

   

    
<Admin_console consoleWidth={80} consoleHeight={80} groupList={data} addObjectLabel = "Key" deleteObjectLabel="Key" tableHeader={["Key ID"]} tableBody={[]} caption="API Key Management Console" pointerEvents={false} ></Admin_console>

    </Admin_navigation>

    </>
  );
}





export default Admin_key_management;