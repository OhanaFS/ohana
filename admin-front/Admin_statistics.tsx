



import { Box, Button, Card, Center, Container, Table, Text, TextInput, Title, useMantineTheme, Image, Grid } from '@mantine/core';

import img from '../src/images/2.png';
import img2 from '../src/images/3.png';
function Admin_statistics() {
   const theme = useMantineTheme();
   //static data
   const noOfFilesWithRepl = 46;
   const noOfFilesWithoutRepl = 545;
   const diskUsage = 46;
   const logs = ["User a have store file A", "User b have store file B", "User c have store file C"];

   return (


      <>
         <Table>
            <tr>
               <td style={{ width: "35%" }}>
                  <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "1%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                     shadow="sm"
                     p="xl"      >         
                           
                     <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                     <Text weight={700}>  Number of files stored: {noOfFilesWithRepl}  </Text>
                     </Card.Section>

                     <Image style={{ textAlign: 'left', marginLeft: "1%", marginTop: "10%" }}
                        radius="md"
                        src={img}/>

                  </Card>

               </td>
               <td style={{ width: "35%" }}>
                  <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "1%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                     shadow="sm"
                     p="xl"

                  >
                     <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "0%" }}>
                        <Text weight={700}> Total size of files stored: {noOfFilesWithoutRepl}<br></br> (not incl. replicas) </Text>

                     </Card.Section>

                     <Image style={{ textAlign: 'left', marginLeft: "1%", marginTop: "10%" }}
                        radius="md"
                        src={img}

                     />

                  </Card>

               </td>
               <td rowSpan={2} style={{ width: "35%" }}>
                  <Card style={{ marginLeft: "3%", height: '700px', border: '1px solid ', marginTop: "1%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                     shadow="sm"
                     p="xl"

                  >
                     <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                         <h2>Logs</h2> 

                     </Card.Section>

                     {logs.map(logs => <p>{logs}</p>)}
                     <Button variant="light" color="blue" style={{ textAlign: "right", marginLeft: '70%' }}>
                        Export logs
                     </Button>
                  </Card>

               </td>
            </tr>
            <tr>
               <td>
                  <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "3%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                     shadow="sm"
                     p="xl"

                  >
                     <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                        <Text weight={700}>  Total Disk usage: {diskUsage}gb  </Text>

                     </Card.Section>

                     <Image height={300} style={{ textAlign: 'left', marginLeft: "1%", marginTop: "2%" }}

                        src={img2}

                     />

                  </Card>

               </td>
               <td>
                  <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "3%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[3] }}
                     shadow="sm"
                     p="xl"

                  >
                     <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                        <Text weight={700}> Cluster Health :  </Text>

                     </Card.Section>

                     <Image style={{ textAlign: 'left', marginLeft: "1%", marginTop: "2%" }}
                        radius="md"
                        src={img2}

                     />

                  </Card>

               </td>


            </tr>

         </Table>





      </>


   );
}





export default Admin_statistics;
