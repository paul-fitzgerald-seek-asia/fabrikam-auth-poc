import React from "react";
import { Router, Route, Switch } from "react-router-dom";
import { Container } from "reactstrap";
import Loading from "./components/Loading";
import Footer from "./components/Footer";
import Header from "./components/Header";
import Home from "./views/Home";
import Screening from "./views/Screening";
import { useAuth0 } from "@auth0/auth0-react";
import history from "./utils/history";

// styles
import "./App.css";

const App = () => {
  const auth = useAuth0();

  if (auth.error) {
    return <div>Oops... {auth.error.message}</div>;
  }

  if (auth.isLoading) {
    return <Loading />;
  }

  return (
    <Router history={history}>
      <div id="app" className="d-flex flex-column h-100">
        <Header auth={auth} />
        <Container className="flex-grow-1 mt-5">
          <Switch>
            <Route path="/" exact component={Home} />
            <Route path="/screening" component={Screening} />
          </Switch>
        </Container>
        <Footer />
      </div>
    </Router>
  );
};

export default App;
