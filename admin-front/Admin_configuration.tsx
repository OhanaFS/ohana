

import { Grid, Button, Paper, useMantineTheme,Text, GroupedTransition, Textarea, Checkbox, Table, Center, Card, ScrollArea } from "@mantine/core";
import { NONAME } from "dns";
import { Link } from "react-router-dom";
import Admin_navigation from "./Admin_navigation";



function Admin_configuration() {
    const theme = useMantineTheme();
 return (
     
       <>
  



      <Center style={{}}>
  
  <Grid style={{width:"60vh"}}> 
  
        
  
        <Card  style={{ marginLeft: "0%", height: '45vh', border: '1px solid ', marginTop: "10%", width: "160%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}
                       shadow="sm"
                       p="xl"
  
                    >
    
                     
                    
   
      
      
       
        <Table captionSide="top" verticalSpacing="md" >
        <caption style={{textAlign: "left",fontWeight:"600",fontSize:"24px",color:"black",marginLeft:"1%"}}>Rotate Key</caption>
        <tbody >
        <tr >
      <td style={{textAlign: "left",fontWeight:"400",fontSize:"16px",color:"black",border:"none"}} width="100%" > Specify the file/directory location and the system will auto rotate the key</td>

        </tr>
     <tr >
      <td style={{ border:"none"}}>
     <Textarea style={{}}
      label="File location"
      radius="xs"
      size="md"
    />
</td>

     </tr>
      <tr>
    <td style={{display:"flex", textAlign: "left",fontWeight:"400",fontSize:"16px",color:"black"}}>Master Key :   <Checkbox style={{marginLeft:"2%"}} > </Checkbox></td>
      
    </tr>
    </tbody>
      </Table>
  
    
       <div style={{ display: "flex" }}>
       <Button  variant="default" color="dark" size="md" style={{ marginLeft: "auto", marginTop:"3%" }}  >Rotate Key</Button>
  
       </div>
    
                    </Card>   
        
        
         
          
  
       
  
  
  
  
      </Grid>
  
      </Center>

</>
   
 );
    }

    



export default Admin_configuration;