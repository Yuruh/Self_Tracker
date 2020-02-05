import React from 'react';
import './App.css';
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import Api from "./Api";
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from "react-router-dom";

const App = () => {
  return (
      <Router>
          <div>
              <nav>
                  <ul>
                      <li>
                          <Link to="/login">Login</Link>
                      </li>
                      <li>
                          <Link to="/home">Home</Link>
                      </li>
                      <li>
                          <Link to="/register">Register</Link>
                      </li>
                  </ul>
              </nav>

              {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
              <Switch>
                  <Route path="/login">
                      <LoginPage/>
                  </Route>
                  <Route path="/home">
                      <HomePage/>
                  </Route>
                  <Route path="/register">
                      <RegisterPage/>
                  </Route>
              </Switch>
          </div>
      </Router>
  );
};

function LoginPage() {
    return (
        <div style={{
            width: "100%",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flexDirection: "column",
        }}>
            <Typography component={"h1"}>Log in</Typography>
            <TextField required label="Email" placeholder="awesome@mail.com" />
            <TextField required label="Password" type="password"/>
            <Link to={"/register"}>Register here</Link>
        </div>
    )
}

function RegisterPage() {
    return (
        <div style={{
            width: "100%",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flexDirection: "column",
        }}>
            <Typography component={"h1"}>Register</Typography>
            <TextField required label="Email" placeholder="awesome@mail.com" />
            <TextField required label="Password" type="password"/>
            <Link to={"/login"}>Log in here</Link>
        </div>
    )
}

function HomePage() {
    return <Button onClick={async () => {
        try {
            const res = await Api.getSpotifyUrl();
//                console.log(res.data)
            window.location.replace(res.data);
        } catch (e) {
            console.log(e.message)
        }
    }}>Authorize Spotify</Button>
}

export default App;