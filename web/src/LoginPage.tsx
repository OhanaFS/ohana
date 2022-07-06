import React from 'react';
import CSS from 'csstype';
import { MediaQuery, Title, TextInput, Button, useMantineTheme, CSSObject } from "@mantine/core";
import { Link } from "react-router-dom";
import './assets/styles.css'

function LoginPage() {
  const theme = useMantineTheme();

  return (
    <div style={{
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      height: "100vh",
    }}>
      <div className='loginBox'>
        <Title style={{
          marginBottom: '15 px'
        }} order={2}>Ohana </Title>
        <TextInput required label="Email" placeholder="Email" sx={(theme) => ({
          display: 'block',
          textAlign: "left",
          width: "90%",
          height: "10vh",
          paddingLeft: theme.spacing.md,
          borderRadius: theme.radius.sm,
          color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
        })} />

        <Button<typeof Link> style={{ marginBottom: "2%" }} variant="default" color="dark" radius="xs" size="md" component={Link} to="/Admin_statistics"    >
          Login Using SSO
        </Button>
      </div>
    </div>
  )
}

export default LoginPage