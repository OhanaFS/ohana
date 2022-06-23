



import { TextInput, Button, Title,useMantineTheme,BackgroundImage, Box, Center ,Text   } from '@mantine/core';


import img from '../src/images/1.png';

import{
  BrowserRouter as Router,
  Link,
  Route,
  Routes


} from "react-router-dom";




function LoginPage({setIsToggled}:any) {

  
    const theme = useMantineTheme();
    const backgroundimage = require('../src/images/5.webp'); 
    return (
 <>   
     <div
      style={{
        backgroundImage:
          `url(${backgroundimage})`,
        backgroundPosition: "center",
        backgroundSize: "cover",
        backgroundRepeat: "no-repeat",
        width: "100vw",
        height: "100vh",
      }}


      
    >

<Box sx={{}} mx="auto">
     
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

        <Button variant="default" color="dark" radius="xs" size="md"  onClick={()=>setIsToggled(true)}      > 
                        Login Using SSO
                    </Button>
      
           
                 
         
                      </div>
      </Center>
    
      </Box>




    </div>
       
      
      </>


    );
    
}


export default LoginPage;
