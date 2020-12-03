import React from "react";
import { useAuth0 } from '@auth0/auth0-react';

const LogoutButton = () => {
  const { isAuthenticated, logout } = useAuth0();
  const doLogout = () => {
    logout({
        returnTo: window.location.origin,
    })
  }
  return <button disabled={!isAuthenticated} onClick={doLogout}>Log Out</button>
};

export default LogoutButton;