import {
  Grid,
  Table,
  Button,
  Center,
  ScrollArea,
  Card,
  Checkbox,
  useMantineTheme,
} from '@mantine/core';
import AppBase from './components/AppBase';

export function AdminSsoGroupsInside() {
  const data = [['Tom'], ['Peter'], ['Raymond']];

  const ths = (
    <tr>
      <th
        style={{
          width: '80%',
          textAlign: 'left',
          fontWeight: '700',
          fontSize: '16px',
          color: 'black',
        }}
      >
        List of Users inside this group
      </th>
    </tr>
  );
  const rows = data.map((items) => (
    <tr>
      <td
        width="80%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items}
      </td>
      <td>
        <Checkbox></Checkbox>{' '}
      </td>
    </tr>
  ));

  const theme = useMantineTheme();

  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
        <Center style={{ marginRight: '' }}>
          <Grid style={{ width: '50vh' }}>
            <Grid.Col span={12}>
              <Card
                style={{
                  marginLeft: '0%',
                  height: '50vh',
                  border: '1px solid ',
                  marginTop: '3%',
                  width: '160%',
                  background:
                    theme.colorScheme === 'dark'
                      ? theme.colors.dark[8]
                      : theme.white,
                }}
                shadow="sm"
                p="xl"
              >
                <Card.Section
                  style={{ textAlign: 'left', marginLeft: '0%' }}
                ></Card.Section>

                <ScrollArea
                  style={{ height: '90%', width: '100%', marginTop: '1%' }}
                >
                  <Table
                    captionSide="top"
                    striped
                    highlightOnHover
                    verticalSpacing="sm"
                  >
                    <caption
                      style={{
                        textAlign: 'center',
                        fontWeight: '600',
                        fontSize: '24px',
                        color: 'black',
                      }}
                    >
                      User Management Console
                    </caption>
                    <thead>{ths}</thead>
                    <tbody>{rows}</tbody>
                  </Table>
                </ScrollArea>
                <tr>
                  <td width={'80%'}>
                    {' '}
                    <Button
                      variant="default"
                      color="dark"
                      size="md"
                      style={{ marginLeft: 'auto', marginTop: '3%' }}
                    >
                      Add User
                    </Button>
                  </td>

                  <td>
                    {' '}
                    <Button
                      variant="default"
                      color="dark"
                      size="md"
                      style={{ marginLeft: 'auto', marginTop: '3%' }}
                    >
                      Delete User
                    </Button>
                  </td>
                </tr>
              </Card>
            </Grid.Col>
          </Grid>
        </Center>
      </AppBase>
    </>
  );
}


