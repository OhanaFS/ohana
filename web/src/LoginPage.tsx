import { Title, TextInput, Button, useMantineTheme } from '@mantine/core';
import { Link } from 'react-router-dom';
import './assets/styles.css';
import backgroundimage from './images/2.webp';

export function LoginPage() {
  const theme = useMantineTheme();
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        backgroundImage: `url(${backgroundimage})`,
        backgroundPosition: 'center',
        backgroundSize: 'cover',
        backgroundRepeat: 'no-repeat',
        width: '100vw',
      }}
    >
      <div className="loginBox">
        <Title
          style={{
            marginBottom: '15 px',
          }}
          order={2}
        >
          Ohana{' '}
        </Title>
        <TextInput
          required
          label="Email"
          placeholder="Email"
          sx={(theme) => ({
            display: 'block',
            textAlign: 'left',
            width: '90%',
            height: '10vh',
            paddingLeft: theme.spacing.md,
            borderRadius: theme.radius.sm,
            color:
              theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
          })}
        />

        <Button<typeof Link>
          style={{ marginBottom: '2%' }}
          variant="default"
          color="dark"
          radius="xs"
          size="md"
          component={Link}
          to="/dashboard"
        >
          Login Using SSO
        </Button>
      </div>
    </div>
  );
}
