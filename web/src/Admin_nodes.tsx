


import { Grid, Table, Button, Text, Center, Checkbox, useMantineTheme, Card, ScrollArea } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";
import {
  Link,


} from "react-router-dom";
import Admin_console from "./Admin_console";


function Admin_nodes() {

  const data = ["192.168.1.1", "192.168.1.2"]
  

  return (

    <>
      <Admin_navigation>
     
        


      <Admin_console consoleWidth={80} consoleHeight={80} groupList={data} addObjectLabel = "Node" deleteObjectLabel="Node"  tableBody={[]} tableHeader={["IP address of the node"]} caption="Node Management Console" pointerEvents={false} ></Admin_console>
      </Admin_navigation>
    </>

  );
}





export default Admin_nodes;