import { Title, Button } from '@mantine/core';
import './assets/styles.css';
import backgroundimage from './images/2.webp';

function LoginPage() {
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

        <Button<'a'>
          style={{ marginBottom: '2%' }}
          variant="default"
          color="dark"
          radius="xs"
          size="md"
          component={'a'}
          href="/auth/login"
        >
          Login Using SSO
        </Button>
      </div>
    </div>
  );
}

export default LoginPage;
