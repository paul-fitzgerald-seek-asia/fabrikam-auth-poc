import React from "react";
import { useAuth0 } from '@auth0/auth0-react';

const LoginButton = () => {
  const { loginWithRedirect } = useAuth0();
  const doLogin = () => {
    loginWithRedirect({
      redirect_uri: window.location.origin,
    })
  }
  return <button onClick={doLogin}>Log In</button>
};

export default LoginButton;