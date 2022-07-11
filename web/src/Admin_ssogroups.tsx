import {
  Grid,
  Table,
  Button,
  Text,
  Center,
  Checkbox,
  useMantineTheme,
  Card,
  ScrollArea,
} from '@mantine/core';
import Admin_navigation from './Admin_navigation';

import { Link } from 'react-router-dom';
import { useState } from 'react';
import { ResponsiveContainer } from 'recharts';
import Admin_console from './Admin_console';
import AppBase from './components/AppBase';

function Admin_ssogroups() {
  /*
    groupList: Array<string>;
    addObjectLabel:string;
    deleteObjectLabel:string;
    tableHeader: string;
    caption: string;
    pointerEvents : boolean; 
    conso
  */

  const SSOGroupList = ['Hr', 'asd'];

  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
        <Admin_console
          consoleWidth={80}
          consoleHeight={60}
          groupList={SSOGroupList}
          addObjectLabel="Group"
          deleteObjectLabel="Group"
          tableHeader={['Current SSO Groups']}
          tableBody={[]}
          caption="SSO Group Management Console"
          pointerEvents={true}
        ></Admin_console>
      </AppBase>
    </>
  );
}

export default Admin_ssogroups;
