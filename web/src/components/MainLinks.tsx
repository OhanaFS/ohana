import React from 'react';
import { Link } from "react-router-dom";
import { ThemeIcon, UnstyledButton, Group, Text } from '@mantine/core';

interface MainLinkProps {
  icon: React.ReactNode;
  color: string;
  label: string;
}

function MainLink({ icon, color, label }: MainLinkProps) {
  return (
    <Link style={{
      textDecoration: 'none',
    }} to={label.toLowerCase()}>
      <UnstyledButton
        sx={(theme) => ({
          display: 'block',
          width: '100%',
          padding: theme.spacing.xs,
          borderRadius: theme.radius.sm,
          color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,

          '&:hover': {
            backgroundColor:
              theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
          },
        })}
      >
        <Group>
          <ThemeIcon color={color} variant="light">
            {icon}
          </ThemeIcon>

          <Text size="md">{label}</Text>
        </Group>
      </UnstyledButton>
    </Link>
  );
}

type MainLinksProps = {
  links: MainLinkProps[];
}

export function MainLinks(props: MainLinksProps) {
  const links = props.links.map((link) => <MainLink {...link} key={link.label} />);
  return <div>{links}</div>;
}
