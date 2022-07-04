


import { Grid, Table, Button, Text, Center, Checkbox, useMantineTheme, Card, ScrollArea,Dialog } from "@mantine/core";
import Admin_navigation from "./Admin_navigation";

import {
   Link,


} from "react-router-dom";
import { useState } from "react";

export interface ConsoleDetails {
    groupList: Array<any>;
    addObjectLabel:string;
    deleteObjectLabel:string;
    tableHeader: Array<string>;
    tableBody : Array<string>;
    caption: string;
    pointerEvents : boolean;
    consoleWidth : number;
    consoleHeight : number;
  }   
 

function Admin_console( props:ConsoleDetails ) {

    // take from database.

    const theme = useMantineTheme();
    let [CurrentSSOGroups, setValue] = useState(props.groupList);
    
    function add(){
        setValue (prevValue => CurrentSSOGroups.concat("asd"));
       
       }
       
      
   

   
  
      const [checkedOne, setCheckedOne] = useState([""]);
     
      function deleteGroup(){
        

        checkedOne.forEach(element =>{
        
          setCheckedOne(checkedOne.filter(item=>item!==element));
            setValue(CurrentSSOGroups.filter(item=>item!==element));

            console.log(CurrentSSOGroups);
           
          
        
     
      })
        
     /*   checkedOne.forEach(element => {
          console.log("Element is " +element);


          CurrentSSOGroups.forEach(item=>{  
            console.log("item " + item);
            setValue(CurrentSSOGroups.filter(item=>item!==element));
            setCheckedOne(checkedOne.filter(item=>item!==element));
          
          });
        

        
          });*/
       
        }
      function update(index:string){
      
        setCheckedOne(prevValue => checkedOne.concat(index));
      } 
      function remove(index:string){
     
       setCheckedOne(checkedOne.filter(item=>item!==index));

      } 
     
    const ths = (
     
        <tr >
          
          <th style={{ width: "80%", textAlign: "left", fontWeight: "700", fontSize: "16px", color: "black" }}>{props.tableHeader}</th>
         
        </tr>
      );
      const rows = CurrentSSOGroups.map((items,index) => (
        <tr  >
          <td width="80%" style={{ textAlign: "left", fontWeight: "400", fontSize: "16px", color: "black" }}>

         <Text color="dark" style={{ marginBottom: "20%", height: "50px",pointerEvents:props.pointerEvents?  'auto' : 'none'  }} component={Link} to="/Admin_ssogroups_inside" variant="link"   >
                  {items}
                </Text>

          </td>
          <td >

         
               <Checkbox onChange={(event) => event.currentTarget.checked? update(items) : remove(items)}  ></Checkbox> 
           </td>
        </tr>
      ));

 
    return (
 
       <>

     
             <Center style={{ marginLeft: "15%" }}>
 
 <Grid style={{ width: props.consoleWidth+"vh" }}>
 
   <Grid.Col span={12}>

     <Card style={{ marginLeft: "0%", height: props.consoleHeight+"vh", border: '1px solid ', marginTop: "4%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}
       shadow="sm"
       p="xl"
     >
      
       <Card.Section style={{ textAlign: 'left', marginLeft: "0%" }}>
 
 
       </Card.Section>
 
 
 
 
 
       <ScrollArea style={{ height: "90%", width: "100%", marginTop: "1%" }}>
   
         <Table captionSide="top" verticalSpacing="sm" >
           <caption style={{ textAlign: "center", fontWeight: "600", fontSize: "24px", color: "black",marginTop:"2%" }}>{props.caption} </caption>
           <thead>{ths}</thead>
   
           <tbody>{rows}</tbody>
  
         </Table>
       
       </ScrollArea>
             <tr>
      <td width={"80%"}>   <Button variant="default" color="dark" size="md" style={{ marginLeft: "auto", marginTop: "3%" } } onClick={()=>add()}   >Add {props.addObjectLabel}</Button></td>
       
      <td>   <Button variant="default" color="dark" size="md" style={{ marginLeft: "auto", marginTop: "3%" }} onClick={()=>deleteGroup() }  >Delete {props.deleteObjectLabel}</Button></td>
      
      </tr>
     </Card>
 
 
 
 
 
 
 
 
 
   </Grid.Col>
 
 </Grid>
 
 </Center>
     
       </>
 
    );
 }
 
 
 
 
 
 export default Admin_console;
 