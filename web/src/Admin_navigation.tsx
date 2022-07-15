import {
  AppShell,
  Burger,
  Header,
  MediaQuery,
  Navbar,
  Text,
  Title,
  useMantineTheme,
  createStyles,
} from '@mantine/core';
import { useState } from 'react';

import { BrowserRouter as Router, Link } from 'react-router-dom';

import Admin_maintenancesettings from './Admin_maintenancesettings';
import Admin_maintenanceresults from './Admin_maintenanceresults';
import Admin_performmaintenance from './Admin_performmaintenance';
import Admin_key_management from './Admin_key_management';
import Admin_create_sso_key from './Admin_create_sso_key';
import Admin_ssogroups_inside from './Admin_ssogroups_inside';
import { BrowserView, MobileView } from 'react-device-detect';

import LoginPage from './LoginPage';
import {
  Icon as TablerIcon,
  User,
  Settings,
  Dashboard,
  Key,
  Rotate,
  Server,
  Edit,
} from 'tabler-icons-react';

type AppBaseProps = {
  children: React.ReactNode;
};
function Admin_navigation(props: AppBaseProps) {
  const theme = useMantineTheme();

  const useStyles = createStyles((theme, _params, getRef) => {
    const icon = getRef('icon');

    return {
      navbar: {
        backgroundColor: theme.colors.gray[0],
      },

      version: {
        backgroundColor: theme.colors[theme.primaryColor][2],
        color: theme.white,
        fontWeight: 700,
      },

      header: {
        paddingBottom: theme.spacing.md,
        marginBottom: theme.spacing.md * 1.5,
        borderBottom: `1px solid ${theme.colors[theme.primaryColor][7]}`,
      },

      footer: {
        paddingTop: theme.spacing.md,
        marginTop: theme.spacing.md,
        borderTop: `1px solid ${theme.colors.gray[7]}`,
      },

      link: {
        ...theme.fn.focusStyles(),
        display: 'flex',
        alignItems: 'center',
        textDecoration: 'none',
        fontSize: theme.fontSizes.sm,

        color: theme.white,
        padding: `${theme.spacing.xs}px ${theme.spacing.sm}px`,
        borderRadius: theme.radius.sm,
        fontWeight: 500,

        '&:hover': {
          backgroundColor: theme.colors.gray[4],
        },
      },

      linkIcon: {
        ref: icon,
        color: theme.black,
        opacity: 0.75,
        marginRight: theme.spacing.sm,
      },

      linkActive: {
        '&, &:hover': {
          backgroundColor: theme.colors.gray[4],
          [`& .${icon}`]: {
            opacity: 0.9,
          },
        },
      },
    };
  });
  interface NavbarLinkProps {
    icon: TablerIcon;
    label: string;
    active?: boolean;
    onClick?(): void;
  }

  const data = [
    { link: '', label: 'Dashboard', icon: Dashboard, to: '/Admin_statistics' },
    { link: '', label: 'SSO', icon: User, to: '/Admin_ssogroups' },
    { link: '', label: 'Nodes', icon: Server, to: '/Admin_nodes' },
    { link: '', label: 'Maintenance', icon: Key, to: '/Admin_maintenancelogs' },
    { link: '', label: 'Settings', icon: Settings, to: '/Admin_settings' },
    { link: '', label: 'Rotate Key', icon: Rotate, to: '/Admin_configuration' },
    {
      link: '',
      label: 'Key Management',
      icon: Edit,
      to: '/Admin_key_management',
    },
  ];
  const { classes, cx } = useStyles();
  const backgroundimage = require('../src/images/2.webp');
  const [opened, setOpened] = useState(false);
  const [active, setActive] = useState('');
  const links = data.map((item) => (
    <a
      className={cx(classes.link, {
        [classes.linkActive]: item.label === active,
      })}
      href={item.link}
      key={item.label}
      onClick={(event) => {
        event.preventDefault();
        setActive(item.label);
      }}
    >
      <item.icon className={classes.linkIcon} />

      <span>
        {' '}
        <Text color="dark" component={Link} variant="link" to={item.to}>
          {item.label}
        </Text>
      </span>
    </a>
  ));
  return (
    <>
      <BrowserView>
        <AppShell
          styles={{
            main: {
              background:
                theme.colorScheme === 'dark'
                  ? theme.colors.dark[8]
                  : theme.colors.gray[1],
              backgroundImage: `url(${backgroundimage})`,
              backgroundPosition: 'center',
              backgroundSize: 'cover',
              backgroundRepeat: 'no-repeat',
              width: '100vw',
              height: '100vh',
            },
          }}
          navbarOffsetBreakpoint="sm"
          asideOffsetBreakpoint="sm"
          fixed
          navbar={
            <Navbar
              height={'100vh'}
              width={{ sm: '10vw' }}
              p="md"
              className={classes.navbar}
            >
              <Navbar.Section grow>{links}</Navbar.Section>

              <Navbar.Section className={classes.footer}>
                <a href="#" className={classes.link}>
                  <Text
                    color="dark"
                    style={{ marginBottom: '20%', height: '50px' }}
                    component={Link}
                    to="/"
                    variant="link"
                  >
                    Logout
                  </Text>
                </a>
              </Navbar.Section>
            </Navbar>
          }
          header={
            <Header height={70} p="md">
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  height: '100%',
                }}
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

                <Title style={{ marginLeft: '1%' }} order={2}>
                  Ohana{' '}
                </Title>
              </div>
            </Header>
          }
        >
          {props.children}
        </AppShell>
      </BrowserView>

      <MobileView>
        <AppShell
          styles={{
            main: {
              background:
                theme.colorScheme === 'dark'
                  ? theme.colors.dark[8]
                  : theme.colors.gray[1],
              backgroundImage: `url(${backgroundimage})`,
              backgroundPosition: 'center',
              backgroundSize: 'cover',
              backgroundRepeat: 'no-repeat',
              width: '100vw',
              height: '100vh',
            },
          }}
          navbarOffsetBreakpoint="sm"
          asideOffsetBreakpoint="sm"
          fixed
          navbar={
            <Navbar
              height={'100vh'}
              width={{ sm: '1' }}
              p="md"
              className={classes.navbar}
            >
              <Navbar.Section grow>{links}</Navbar.Section>

              <Navbar.Section className={classes.footer}>
                <a href="#" className={classes.link}>
                  <Text
                    color="dark"
                    style={{ marginBottom: '20%', height: '50px' }}
                    component={Link}
                    to="/"
                    variant="link"
                  >
                    Logout
                  </Text>
                </a>
              </Navbar.Section>
            </Navbar>
          }
          header={
            <Header height={70} p="md">
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  height: '100%',
                }}
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

                <Title style={{ marginLeft: '1%' }} order={2}>
                  Ohana{' '}
                </Title>
              </div>
            </Header>
          }
        >
          {props.children}
        </AppShell>
      </MobileView>
    </>
  );
}

export default Admin_navigation;
