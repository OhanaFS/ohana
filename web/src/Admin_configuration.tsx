import {
  Grid,
  Button,
  useMantineTheme,
  Textarea,
  Checkbox,
  Table,
  Center,
  Card,
} from '@mantine/core';

import Admin_navigation from './Admin_navigation';
import AppBase from './components/AppBase';

function Admin_configuration() {
  const theme = useMantineTheme();
  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
        <Center style={{}}>
          <Grid style={{ width: '60vh' }}>
            <Card
              style={{
                marginLeft: '0%',
                height: '45vh',
                border: '1px solid ',
                marginTop: '8%',
                width: '160%',
                background:
                  theme.colorScheme === 'dark'
                    ? theme.colors.dark[8]
                    : theme.white,
              }}
              shadow="sm"
              p="xl"
            >
              <Table captionSide="top" verticalSpacing="md">
                <caption
                  style={{
                    textAlign: 'left',
                    fontWeight: 600,
                    fontSize: '24px',
                    color: 'black',
                    marginLeft: '1%',
                  }}
                >
                  Rotate Key
                </caption>
                <tbody>
                  <tr>
                    <td
                      style={{
                        textAlign: 'left',
                        fontWeight: 400,
                        fontSize: '16px',
                        color: 'black',
                        border: 'none',
                      }}
                      width="100%"
                    >
                      {' '}
                      Specify the file/directory location and the system will
                      auto rotate the key
                    </td>
                  </tr>
                  <tr>
                    <td style={{ border: 'none' }}>
                      <Textarea
                        style={{}}
                        label="File location"
                        radius="xs"
                        size="md"
                      />
                    </td>
                  </tr>
                  <tr>
                    <td
                      style={{
                        display: 'flex',
                        textAlign: 'left',
                        fontWeight: 400,
                        fontSize: '16px',
                        color: 'black',
                      }}
                    >
                      Master Key :{' '}
                      <Checkbox style={{ marginLeft: '2%' }}> </Checkbox>
                    </td>
                  </tr>
                </tbody>
              </Table>

              <div style={{ display: 'flex' }}>
                <Button
                  variant="default"
                  color="dark"
                  size="md"
                  style={{ marginLeft: 'auto', marginTop: '3%' }}
                >
                  Rotate Key
                </Button>
              </div>
            </Card>
          </Grid>
        </Center>
      </AppBase>
    </>
  );
}

export default Admin_configuration;
