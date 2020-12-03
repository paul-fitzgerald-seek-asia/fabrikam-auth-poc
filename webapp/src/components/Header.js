import React from "react";
import Logo from './Logo'
import LoginButton from './LoginButton'
import LogoutButton from './LogoutButton'
import { useAuth0 } from '@auth0/auth0-react';

const SessionControl = () => {
  let authResult, button;
  const { isAuthenticated, user } = useAuth0();
  if (isAuthenticated) {
    const { email } = user
    authResult = `Logged in as: ${email}`
    button = (
      <LogoutButton />
    );
  } else {
    authResult = "You are not logged in"
    button = (
      <LoginButton />
    );
  }
  return (
    <div className="vert-bump horz-bump fg-contrast float-right">
      <span>{authResult}</span> &nbsp; {button}
    </div>
  );
}

const Header = () => (
  <header className="bg-contrast h-100">
    <Logo />
    <SessionControl />
  </header>
);

export default Header;
