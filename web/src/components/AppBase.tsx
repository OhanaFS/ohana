import React, { useState } from 'react';
import {
  IconHome2,
  IconStar,
  IconShare,
  IconServer2,
  IconSettings,
  IconEdit,
  IconTool,
  IconRotate2,
  IconUserPlus,
  IconAdjustments,
  IconLogout,
} from '@tabler/icons';
import backgroundimage from '../images/2.webp';

import {
  AppShell,
  Navbar,
  Header,
  Title,
  MediaQuery,
  Burger,
  useMantineTheme,
  Menu,
  createStyles,
} from '@mantine/core';
import { MainLinks } from './MainLinks';
import { UserButton } from './UserButton';
import { useQueryUser } from '../api/auth';

const useStyles = createStyles((theme) => ({
  navbar: {
    backgroundColor:
      theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.white,
    paddingBottom: 0,
  },

  header: {
    padding: theme.spacing.md,
    paddingTop: 0,
    marginLeft: -theme.spacing.md,
    marginRight: -theme.spacing.md,
    color: theme.colorScheme === 'dark' ? theme.white : theme.black,
    borderBottom: `1px solid ${
      theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[3]
    }`,
  },

  links: {
    marginLeft: -theme.spacing.md,
    marginRight: -theme.spacing.md,
  },

  linksInner: {
    paddingTop: theme.spacing.xl,
    paddingBottom: theme.spacing.xl,
  },

  footer: {
    marginLeft: -theme.spacing.md,
    marginRight: -theme.spacing.md,
    borderTop: `1px solid ${
      theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[3]
    }`,
  },
}));

type AppBaseProps = {
  children: React.ReactNode;
  userType: string;
};

export default function AppBase(props: AppBaseProps) {
  const { classes } = useStyles();
  const theme = useMantineTheme();
  const user = useQueryUser();
  const [opened, setOpened] = useState(false);
  const data_user = [
    { icon: <IconHome2 size={16} />, color: 'blue', label: 'Home' },
    { icon: <IconStar size={16} />, color: 'teal', label: 'Favourites' },
    { icon: <IconShare size={16} />, color: 'violet', label: 'Shared' },
  ];
  const data_admin = [
    { icon: <IconHome2 size={16} />, color: 'blue', label: 'Dashboard' },
    { icon: <IconUserPlus size={16} />, color: 'blue', label: 'SSO' },
    { icon: <IconServer2 size={16} />, color: 'blue', label: 'Nodes' },
    { icon: <IconTool size={16} />, color: 'blue', label: 'Maintenance' },
    { icon: <IconSettings size={16} />, color: 'blue', label: 'Settings' },
    { icon: <IconRotate2 size={16} />, color: 'blue', label: 'Rotate Key' },
    { icon: <IconEdit size={16} />, color: 'blue', label: 'Key Management' },
  ];

  return (
    <AppShell
      styles={{
        main: {
          background:
            theme.colorScheme === 'dark'
              ? theme.colors.dark[8]
              : theme.colors.gray[0],
          ...(props.userType === 'user'
            ? {}
            : {
                backgroundImage: `url(${backgroundimage})`,
                backgroundPosition: 'center',
                backgroundSize: 'cover',
                backgroundRepeat: 'no-repeat',
                width: '100vw',
              }),
        },
      }}
      navbarOffsetBreakpoint="sm"
      asideOffsetBreakpoint="sm"
      fixed
      navbar={
        <Navbar
          p="md"
          hiddenBreakpoint="sm"
          hidden={!opened}
          width={{ sm: 250 }}
          className={classes.navbar}
        >
          <Navbar.Section grow className={classes.links}>
            <MainLinks
              links={props.userType === 'user' ? data_user : data_admin}
            />
          </Navbar.Section>
          <Navbar.Section className={classes.footer}>
            <Menu
              style={{ width: '100%' }}
              position="right"
              placement="end"
              control={
                <UserButton
                  name={user.data?.name || ''}
                  email={user.data?.email || ''}
                />
              }
            >
              {
                // @ts-ignore
                props.userType === 'admin' && (
                  <>
                    <Menu.Label>Admin</Menu.Label>
                    <Menu.Item
                      icon={<IconAdjustments size={14} />}
                      onClick={() => {}}
                    >
                      Dashboard
                    </Menu.Item>
                  </>
                )
              }

              <Menu.Label>User</Menu.Label>
              <Menu.Item<'a'>
                color="red"
                icon={<IconLogout size={14} />}
                component="a"
                href="/auth/logout"
              >
                Log out
              </Menu.Item>
            </Menu>
          </Navbar.Section>
        </Navbar>
      }
      header={
        <Header height={70} p="md">
          <div
            style={{ display: 'flex', alignItems: 'center', height: '100%' }}
          >
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
