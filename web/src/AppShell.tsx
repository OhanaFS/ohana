import React, { useState } from 'react';
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
}

export default function AppBase(props: AppBaseProps) {
  const theme = useMantineTheme();
  const [opened, setOpened] = useState(false);
  return (
    <AppShell
      styles={{
        main: {
          background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
        },
      }}
      navbarOffsetBreakpoint="sm"
      asideOffsetBreakpoint="sm"
      fixed
      navbar={
        <Navbar p="md" hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 200 }}>
          <Navbar.Section grow mt="md">
          <MainLinks />
        </Navbar.Section>
        <Navbar.Section>
          <User />
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