import React, { useState } from 'react';
import { IconHome2, IconStar, IconShare, IconLogin, IconServer2, IconSettings, IconKey } from '@tabler/icons';
import backgroundimage from '../images/2.webp'

import {
  AppShell,
  Navbar,
  Header,
  Title,
  MediaQuery,
  Burger,
  useMantineTheme,
} from '@mantine/core';
import { MainLinks } from './MainLinks';
import { User } from './user';

type AppBaseProps = {
  children: React.ReactNode;
  image: string
  name: string;
  username: string;
  userType: string;
}

export default function AppBase(props: AppBaseProps) {
  const theme = useMantineTheme();
  const [opened, setOpened] = useState(false);
  const data_user = [
    { icon: <IconHome2 size={16} />, color: 'blue', label: 'Home' },
    { icon: <IconStar size={16} />, color: 'teal', label: 'Favourites' },
    { icon: <IconShare size={16} />, color: 'violet', label: 'Shared' },
  ];
  const data_admin = [
    { icon: <IconHome2 size={16} />, color: 'blue', label: 'Dashboard' },
    { icon: <IconLogin size={16} />, color: 'blue', label: 'SSO' },
    { icon: <IconServer2 size={16} />, color: 'blue', label: 'Nodes' },
    { icon: <IconServer2 size={16} />, color: 'blue', label: 'Maintenance' },
    { icon: <IconSettings size={16} />, color: 'blue', label: 'Settings' },
    { icon: <IconServer2 size={16} />, color: 'blue', label: 'Rotate Key' },
    { icon: <IconKey size={16} />, color: 'blue', label: 'Key Management' },
  ]

  return (
    <AppShell
      styles={{
        main: {
          background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
          backgroundImage: props.userType === 'user' ? "":`url(${backgroundimage})`,
          backgroundPosition: props.userType === 'user' ? "":'center',
          backgroundSize: props.userType === 'user' ? "":'cover',
          backgroundRepeat: props.userType === 'user' ? "":'no-repeat',
          width: props.userType === 'user' ? "": '100vw',
        },
      }}
      navbarOffsetBreakpoint="sm"
      asideOffsetBreakpoint="sm"
      fixed
      navbar={
        <Navbar p="md" hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 250 }}>
          <Navbar.Section grow mt="md">
            <MainLinks links={(props.userType === 'user' ? data_user : data_admin)} />
          </Navbar.Section>
          <Navbar.Section>
            <User name={props.name} username={props.username} image={props.image} />
          </Navbar.Section>
        </Navbar>}
      header={
        <Header height={70} p="md">
          <div style={{ display: 'flex', alignItems: 'center', height: '100%' }}>
            <MediaQuery largerThan="sm" styles={{ display: 'none' }}>
              <Burger
                opened={opened}
                onClick={() => setOpened((o) => !o)}
                size="sm"
                color={theme.colors.gray[6]}
                mr="xl"
              />
            </MediaQuery>

            <Title order={2}>Ohana</Title>
          </div>
        </Header>
      }
    >
      {props.children}
    </AppShell>
  );
}
