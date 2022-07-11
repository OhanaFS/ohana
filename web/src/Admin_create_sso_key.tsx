import {
  Grid,
  Textarea,
  Button,
  Text,
  useMantineTheme,
  Center,
} from '@mantine/core';
import Admin_navigation from './Admin_navigation';

import { BrowserRouter as Router, Link } from 'react-router-dom';
import AppBase from './components/AppBase';

function Admin_create_sso_key() {
  const theme = useMantineTheme();
  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
        <Center>
          <Grid style={{ width: '100vh' }}>
            <Grid.Col
              span={12}
              style={{ marginLeft: '2%', marginTop: '2%', border: '1px solid' }}
            >
              <Text
                underline
                weight={700}
                style={{ marginLeft: '1%', marginTop: '3%' }}
              >
                {' '}
                <h2>Create SSO </h2>{' '}
              </Text>

              <Grid.Col span={12}></Grid.Col>
              <Grid.Col span={12}>
                <Textarea label="SSO Group Name:" radius="xs" size="md" />
              </Grid.Col>

              <div style={{ display: 'flex' }}>
                <Button
                  style={{ marginLeft: 'auto', marginTop: '3%' }}
                  component={Link}
                  to="/sso"
                >
                  Create
                </Button>
              </div>
            </Grid.Col>
          </Grid>
        </Center>
      </AppBase>
    </>
  );
}

export default Admin_create_sso_key;
