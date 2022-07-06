



import { TextInput, Button, Title, useMantineTheme, BackgroundImage, Box, Center, Text, Group, MediaQuery, CSSObject, MantineProvider } from '@mantine/core';




import {
  BrowserRouter as Router,
  Link,



} from "react-router-dom";


import { useMediaQuery, useViewportSize } from '@mantine/hooks';


import { BrowserView, MobileView } from 'react-device-detect';


function LoginPage({ setIsToggled }: any) {
  const theme = useMantineTheme();





  const backgroundimage = require('../src/images/5.webp');

  const highlight: CSSObject = {
    backgroundImage:
      `url(${backgroundimage})`,
    backgroundPosition: "center",
    backgroundSize: "cover",
    backgroundRepeat: "no-repeat",
    width: "100vw",
    height: "100vh",

  };
  const mobile = useMediaQuery('(min-width:600px)');



  const { height, width } = useViewportSize();

  return (
    <>
  
      <MantineProvider
        theme={{
          breakpoints: {
            xs: 576,
            sm: 768,
            md: 992,
            lg: 1200,
            xl: 1400,
          },
        }}
      >


        <BrowserView>

    
          {width > 1400 &&
            <MediaQuery largerThan="lg" styles={highlight}>
           
              <Box sx={(theme) => ({
      
      })} mx="auto"> 
           
              <div>Width: {width}px, height: {height}px</div>
                screen is xl
           
                <Center >
               
                  <div style={{ border: '1px solid ', marginTop: "15%", textAlign: "center", width: "20%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>

                  </div>
            
                </Center>
              
              </Box>
             
            </MediaQuery>
          }

          {width > 1200 && width < 1400 &&

            <MediaQuery largerThan="lg" smallerThan="xl" styles={highlight}>
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is lg

                <Center >
                  <div style={{ border: '1px solid ', marginTop: "20%", textAlign: "center", width: "25%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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
                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>
                  </div>
                </Center>

              </Box>
            </MediaQuery>
          }


          {width > 992 && width < 1200 &&

            <MediaQuery largerThan="md" smallerThan="lg" styles={highlight}>
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is md
                <Center >

                  <div style={{ border: '1px solid ', marginTop: "25%", textAlign: "center", width: "30%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>

                  </div>
                </Center>
              </Box>
            </MediaQuery>
          }

          {width > 768 && width < 992 &&

            <MediaQuery largerThan="sm" smallerThan="md" styles={highlight}>
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is sm
                <Center >

                  <div style={{ border: '1px solid ', marginTop: "25%", textAlign: "center", width: "35%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>

                  </div>
                </Center>
              </Box>
            </MediaQuery>
          }

{width < 768 && width >576 &&

<MediaQuery largerThan="xs" smallerThan="sm" styles={highlight}>
  <Box sx={{}} mx="auto">
    screen is xs
    <Center >

      <div style={{ border: '1px solid ', marginTop: "30%", textAlign: "center", width: "50%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

        <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
          Login Using SSO
        </Button>

      </div>
    </Center>
  </Box>
</MediaQuery>
}
{width < 576  &&

<MediaQuery  smallerThan="xs" styles={highlight}>
<Box sx={{}} mx="auto">
  screen is less than xs 
  <Center >

    <div style={{ border: '1px solid ', marginTop: "50%", textAlign: "center", width: "60%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

      <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
        Login Using SSO
      </Button>

    </div>
  </Center>
</Box>
</MediaQuery>
}














        </BrowserView>



        <MobileView>


        {width > 1400 &&
            <MediaQuery largerThan="lg" styles={highlight}>
              
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is xl
                <Center >

                  <div style={{ border: '1px solid ', marginTop: "15%", textAlign: "center", width: "20%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>

                  </div>
                </Center>
              </Box>
            </MediaQuery>
          }

          {width > 1200 && width < 1400 &&

            <MediaQuery largerThan="lg" smallerThan="xl" styles={highlight}>
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is lg

                <Center >
                  <div style={{ border: '1px solid ', marginTop: "20%", textAlign: "center", width: "25%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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
                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>
                  </div>
                </Center>

              </Box>
            </MediaQuery>
          }


          {width > 992 && width < 1200 &&

            <MediaQuery largerThan="md" smallerThan="lg" styles={highlight}>
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is md
                <Center >

                  <div style={{ border: '1px solid ', marginTop: "25%", textAlign: "center", width: "30%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>

                  </div>
                </Center>
              </Box>
            </MediaQuery>
          }

          {width > 768 && width < 992 &&

            <MediaQuery largerThan="sm" smallerThan="md" styles={highlight}>
              <Box sx={{}} mx="auto">
              <div>Width: {width}px, height: {height}px</div>
                screen is sm
                <Center >

                  <div style={{ border: '1px solid ', marginTop: "25%", textAlign: "center", width: "35%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

                    <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
                      Login Using SSO
                    </Button>

                  </div>
                </Center>
              </Box>
            </MediaQuery>
          }

{width < 768 && width >576 &&

<MediaQuery largerThan="xs" smallerThan="sm" styles={highlight}>
  <Box sx={{}} mx="auto">
    screen is xs
    <Center >

      <div style={{ border: '1px solid ', marginTop: "10%", textAlign: "center", width: "50%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
        <Title order={2}>Ohana </Title>
        <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
          display: 'block',
          textAlign: "left",
          width: "90%",
          height: "20vh",
          padding: theme.spacing.xs,
          borderRadius: theme.radius.sm,
          color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
        })} />

        <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
          Login Using SSO
        </Button>

      </div>
    </Center>
  </Box>
</MediaQuery>
}
{width < 576  &&

<MediaQuery  smallerThan="xs" styles={highlight}>
<Box sx={{}} mx="auto">
  screen is less than xs 
  <Center >

    <div style={{ border: '1px solid ', marginTop: "50%", textAlign: "center", width: "60%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.white }}>
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

      <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
        Login Using SSO
      </Button>

    </div>
  </Center>
</Box>
</MediaQuery>
}














        </MobileView>













      </MantineProvider>

    </>


  );

}


export default LoginPage;
