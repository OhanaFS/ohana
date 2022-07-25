import {
  useMantineTheme,
  Grid,
  Button,
  Card,
  Center,
  Table,
} from '@mantine/core';
import AppBase from './components/AppBase';

export function AdminMaintenanceResults() {
  const logs = [
    'Turning server one offline',
    'Cleaning server one',
    'Turning server one online',
    'Server one is back online',
    'Turning server two offline',
    'Turning server two online',
    'Maintenance is not completed.',
  ];

  const theme = useMantineTheme();

  const rows = logs.map((items) => (
    <tr>
      <td
        width="80%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
          border: 'none',
        }}
      >
        {items}
      </td>
    </tr>
  ));
  return (
    <>
      <AppBase>
        <Center>
          <Grid style={{ width: '100vh' }}>
            <Grid.Col span={12} style={{ marginTop: '' }}>
              <Card
                style={{
                  marginLeft: '10%',
                  height: '500px',
                  border: '1px solid ',
                  marginTop: '5%',
                  width: '60%',
                  background:
                    theme.colorScheme === 'dark'
                      ? theme.colors.dark[8]
                      : theme.white,
                }}
                shadow="sm"
                p="xl"
              >
                <Table captionSide="top">
                  <caption
                    style={{
                      textAlign: 'center',
                      fontWeight: '600',
                      fontSize: '20px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    Maintenance Logs
                  </caption>

                  <tbody>{rows}</tbody>
                </Table>
                <Button
                  variant="default"
                  color="dark"
                  style={{ textAlign: 'right', marginLeft: '70%' }}
                >
                  Export logs
                </Button>
              </Card>
            </Grid.Col>
          </Grid>
        </Center>
      </AppBase>
    </>
  );
}
