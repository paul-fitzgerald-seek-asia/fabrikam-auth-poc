import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import appconfig from './appconfig'
import { Auth0Provider } from "@auth0/auth0-react";

ReactDOM.render(
  <Auth0Provider
    domain={appconfig.auth.domain}
    clientId={appconfig.auth.clientId}
    redirectUri={window.location.origin}
    audience={appconfig.auth.apiAudience}
  >
    <App />
  </Auth0Provider>,
  document.getElementById("root")
);
