



import { TextInput, Button, Title,useMantineTheme,BackgroundImage, Box, Center ,Text, Group, MediaQuery, CSSObject   } from '@mantine/core';


import img from '../src/images/1.png';

import{
  BrowserRouter as Router,
  Link,
  Route,
  Routes


} from "react-router-dom";


import { useMediaQuery } from '@mantine/hooks';


import {BrowserView, MobileView} from 'react-device-detect';


function LoginPage({setIsToggled}:any) {
  const theme = useMantineTheme();





  const backgroundimage = require('../src/images/5.webp'); 
  
  const highlight: CSSObject = {

    border: '1px solid ',  marginTop: "50%",textAlign: "center", width: "80%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white 
  };
  const mobile = useMediaQuery('(min-width:600px)');
  const desktop = useMediaQuery('(min-width:1400px)');
    return (
 <>   
  <div style={{
    backgroundImage:
      `url(${backgroundimage})`,
    backgroundPosition: "center",
    backgroundSize: "cover",
    backgroundRepeat: "no-repeat",
    width: "100vw",
    height: "100vh",
  }}>
            <BrowserView>
          
          
            {desktop==true ?   
              <Box  sx={{}} mx="auto">
     
              <Center >
             
              <div style={{ border: '1px solid ',  marginTop: "15%",textAlign: "center", width: "20%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
              <Title order={2}>Ohana </Title>
              <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
                                display: 'block',
                                textAlign: "left",
                                width: "90%",
                                height: "10vh",
                                padding: theme.spacing.xs,
                                borderRadius: theme.radius.sm,
                                color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
                            })} />
         
                <Button<typeof Link> style={{marginBottom:"2%"}}variant="default" color="dark" radius="xs" size="md"    component={Link} to="/Admin_statistics"    > 
                                Login Using SSO
                            </Button>
                           
                   </div>
                         
                 
                        
              </Center>
            
              </Box>
            
            
            :   
            
            <Box  sx={{}} mx="auto">
     
            <Center >
           
            <div style={{ border: '1px solid ',  marginTop: "30%",textAlign: "center", width: "30%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
            <Title order={2}>Ohana </Title>
            <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
                              display: 'block',
                              textAlign: "left",
                              width: "90%",
                              height: "10vh",
                              padding: theme.spacing.xs,
                              borderRadius: theme.radius.sm,
                              color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
                          })} />
       
              <Button<typeof Link> style={{marginBottom:"2%"}}variant="default" color="dark" radius="xs" size="md"    component={Link} to="/Admin_statistics"    > 
                              Login Using SSO
                          </Button>
                         
                 </div>
                       
               
                      
            </Center>
          
            </Box>
            
            
            }
    
   
            </BrowserView>


        
            <MobileView>

           


{mobile==true ?         

<Box sx={{}} mx="auto">
     
     <Center >

  
 
       <div style={{  border: '1px solid ',  marginTop: "15%",textAlign: "center", width: "40%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
     <Title order={2}>Ohana </Title>
     <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
                       display: 'block',
                       textAlign: "left",
                       width: "90%",
                       height: "15vh",
                       padding: theme.spacing.xs,
                       borderRadius: theme.radius.sm,
                       color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
                   })} />

       <Button<typeof Link> style={{marginBottom:"2%"}}variant="default" color="dark" radius="xs" size="md"    component={Link} to="/Admin_statistics"    > 
                       Login Using SSO
                   </Button>
                  
          
                   </div>         
     
              
            
     </Center>
   
     </Box>
     
     
     
     : 

<Box sx={{}} mx="auto">
     
     <Center >

  
 
       <div style={{  border: '1px solid ',  marginTop: "50%",textAlign: "center", width: "80%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
     <Title order={2}>Ohana </Title>
     <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
                       display: 'block',
                       textAlign: "left",
                       width: "90%",
                       height: "15vh",
                       padding: theme.spacing.xs,
                       borderRadius: theme.radius.sm,
                       color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
                   })} />

       <Button<typeof Link> style={{marginBottom:"2%"}}variant="default" color="dark" radius="xs" size="md"    component={Link} to="/Admin_statistics"    > 
                       Login Using SSO
                   </Button>
                  
          
                   </div>         
     
              
            
     </Center>
   
     </Box>
     
                  }



            </MobileView>


            </div>
    


 



  


       

      </>


    );
    
}


export default LoginPage;
