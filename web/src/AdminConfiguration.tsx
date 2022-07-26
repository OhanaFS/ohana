import { Button, Textarea, Checkbox, Text } from '@mantine/core';

import AppBase from './components/AppBase';

export function AdminConfiguration() {
  return (
    <AppBase
      userType="admin"
      name="Alex Simmons"
      username="@alex"
      image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'flex-start',
          height: '100%',
        }}
      >
        <div
          style={{
            display: 'flex',
            border: '1px solid #ccc',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'flex-start',
            width: '90%',
            backgroundColor: 'white',
            borderRadius: '10px',
            padding: '20px',
            maxWidth: '500px',
          }}
        >
          <caption
            style={{
              textAlign: 'center',
              fontWeight: 600,
              fontSize: '24px',
              color: 'black',
              marginBottom: '20px',
              alignSelf: 'center',
            }}
          >
            Rotate Key
          </caption>
          <Textarea
            label="Specify the file/directory location and the system will
            auto rotate the key"
            radius="md"
            size="lg"
          />
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              margin: '20px 0',
            }}
          >
            <Text>Master Key :</Text>
            <Checkbox style={{ marginLeft: '10px' }}> </Checkbox>
          </div>
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{ alignSelf: 'flex-end' }}
          >
            Rotate Key
          </Button>
        </div>
      </div>
    </AppBase>
  );
}
