



import { TextInput, Button, Title,useMantineTheme,BackgroundImage, Box, Center    } from '@mantine/core';


import img from '../src/images/1.png';

import{
  BrowserRouter as Router,
  Link,
  Route,
  Routes


} from "react-router-dom";



function LoginPage() {

    const theme = useMantineTheme();
  
    return (
   
        <Box sx={{}} mx="auto">
       <BackgroundImage 
           
          src={img}
    
        >
      <Center >
     
      <div style={{ border: '1px solid ',  marginTop: "5%",textAlign: "center", width: "20%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}>
      <Title order={2}>Ohana </Title>
      <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
                        display: 'block',
                        textAlign: "left",
                        width: "95%",
                        padding: theme.spacing.xs,
                        borderRadius: theme.radius.sm,
                        color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
                    })} />

        <Button variant="default" color="dark" radius="xs" size="md"     >
                        Login Using SSO
                    </Button>
                      </div>
      </Center>
      </BackgroundImage>
      </Box>

      
 
    
    );
    
}


export default LoginPage;
