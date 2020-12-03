import React, { Fragment } from "react";
import { Button } from 'reactstrap';
import { Link } from 'react-router-dom'
import Hero from "../components/Hero";

const Home = () => {
  return (
    <Fragment>
      <Hero>
        <Link to="/screening"><Button className="navButton">View and/or Schedule Screening Sessions</Button></Link>
      </Hero>
    </Fragment>
  )
};

export default Home;
