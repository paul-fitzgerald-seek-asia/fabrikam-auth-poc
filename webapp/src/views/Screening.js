import React, { useState, useEffect } from "react";
import { Alert, Button, Input, InputGroup } from "reactstrap";
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react";
import DataTable from 'react-data-table-component';
import DateTimePicker from 'react-datetime-picker';
import Loading from "../components/Loading";

const apiOrigin = "http://localhost:8080/v1/screening";

const columns = [
  {
    name: "Time",
    selector: function(datum) { return datum.startTime; },
  },
  {
    name: "Duration",
    selector: (datum) => datum.duration,
  },
  {
    name: "Title",
    selector: (datum) => datum.title,
  },
];

const listScreenings = async (state, setState, getAccessTokenSilently) => {
  try {
    const token = await getAccessTokenSilently();
    const response = await fetch(`${apiOrigin}/screening`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const responseData = await response.json();
    setState({
      ...state,
      show: true,
      data: responseData,
    });
  } catch (error) {
    setState({
      ...state,
      error: error.error,
    });
  }
};

export const Screening = () => {
  const [state, setState] = useState({
    startTime: null,
    initialAPIcall: false,
    show: false,
    data: [],
    error: null,
  });

  const {
    getAccessTokenSilently,
    loginWithPopup,
    getAccessTokenWithPopup,
  } = useAuth0();

  const handleConsent = async () => {
    try {
      await getAccessTokenWithPopup();
      setState({
        ...state,
        error: null,
      });
    } catch (error) {
      setState({
        ...state,
        error: error.error,
      });
    }
    await listScreenings(state, setState, getAccessTokenSilently);
  };

  const handleLoginAgain = async () => {
    try {
      await loginWithPopup();
      setState({
        ...state,
        error: null,
      });
    } catch (error) {
      setState({
        ...state,
        error: error.error,
      });
    }
    await listScreenings(state, setState, getAccessTokenSilently);
  };

  const makeCreateCall = async () => {
    try {
      const token = await getAccessTokenSilently();
      const response = await fetch(`${apiOrigin}/screening`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          "startTime": state.startTime,
          "duration": parseInt(document.getElementById("durationInput").value),
          "title": document.getElementById("titleInput").value,
        })
      });
      const screening = await response.json();
      console.log(`Created New Screening with ID = ${screening.screeningID}`);
      listScreenings(state, setState, getAccessTokenSilently);
    } catch (error) {
      setState({
        ...state,
        error: error.error,
      });
    }
  }

  useEffect(() => {
    if (!state.initialAPIcall) {
      listScreenings(state, setState, getAccessTokenSilently)
      state.initialAPIcall = true;
    }
  }, [getAccessTokenSilently, state, state.initialAPIcall])

  const handle = (e, fn) => {
    e.preventDefault();
    fn();
  };

  return (
    <>
      <div className="mb-5">
        {state.error === "consent_required" && (
          <Alert color="warning">
            You need to{" "}
            <a
              href="#/"
              className="alert-link"
              onClick={(e) => handle(e, handleConsent)}
            >
              consent to get access to screenings API!
            </a>
          </Alert>
        )}

        {state.error === "login_required" && (
          <Alert color="warning">
            You need to{" "}
            <a
              href="#/"
              className="alert-link"
              onClick={(e) => handle(e, handleLoginAgain)}
            >
              login again!
            </a>
          </Alert>
        )}

        <h1>Scheduled Screening Sessions</h1>
        <DataTable title="Screening Sessions" persistTableHead={true} columns={columns} data={state.data}></DataTable>
        <hr />
        <div className="screening-form">
          <h3>Book New Screening</h3>
          <InputGroup>
          <DateTimePicker onChange={(t) => {state.startTime = t;}} value={state.startTime} />
            <Input id="durationInput" placeholder="Duration (minutes)" min={0} type="number" step="10" size={250} />
            <Input id="titleInput" placeholder="Title of screening" size={32} />
          </InputGroup>
          <Button color="primary" className="mt-5" onClick={makeCreateCall}>
            Schedule New Screening
          </Button>
        </div>
      </div>
    </>
  );
};

export default withAuthenticationRequired(Screening, {
  onRedirecting: () => <Loading />,
});
