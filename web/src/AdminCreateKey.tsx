import {
  Grid,
  Table,
  Button,
  Text,
  TextInput,
  Center,
} from '@mantine/core';

import { BrowserRouter as Router, Link } from 'react-router-dom';
import { useState } from 'react';
import AppBase from './components/AppBase';

function generateRandomString() {
  var result = '';
  var characters =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  var charactersLength = characters.length;
  for (var i = 0; i < 16; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }

  return result;
}

export function AdminCreateKey() {
  let [value, setValue] = useState('');
  function generateKeys() {
    setValue((prevValue) => generateRandomString());
  }

  return (
    <AppBase
      userType="admin"
      name="Alex Simmons"
      username="@alex"
      image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
    >
      <Center>
        <Grid style={{ width: '80vh' }}>
          <Grid.Col
            span={12}
            style={{
              marginLeft: '2%',
              marginTop: '2%',
              maxWidth: '50%',
              border: '1px solid',
            }}
          >
            <Text
              underline
              weight={700}
              style={{ marginLeft: '1%', marginTop: '3%' }}
            >
              {' '}
              <h2>Create API Key</h2>{' '}
            </Text>

            <Grid.Col span={12}>
              <Table>
                <tr>
                  <td>
                    <TextInput
                      label="API Key"
                      radius="xs"
                      size="md"
                      required
                      value={value}
                      onChange={(event) => setValue(event.currentTarget.value)}
                    />
                  </td>
                  <td>
                    <Button onClick={generateKeys}>Generate</Button>
                  </td>
                </tr>
              </Table>
            </Grid.Col>

            <Grid.Col span={12}>
              <TextInput label="Description" radius="xs" size="md" required />
            </Grid.Col>

            <div style={{ display: 'flex' }}>
              <Button
                style={{ marginLeft: 'auto', marginTop: '3%' }}
                component={Link}
                to="/Admin_key_management"
              >
                Create Key
              </Button>
            </div>
          </Grid.Col>
        </Grid>
      </Center>
    </AppBase>
  );
}