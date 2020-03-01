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
    Link, Redirect, useLocation
} from "react-router-dom";
import HomePage from "./HomePage";

// @ts-ignore
const PrivateRoute = ({ component: Component, ...rest }) => (
    <Route {...rest} render={props => {
        const token = Api.token;
        if (!token) {
            // not logged in so redirect to login page with the return url
            return <Redirect to={{ pathname: '/login', state: { from: props.location } }} />
        } else {
            console.log("Attempt to access a private route worked")
        }

        // authorised so return component
        return <Component {...props} />
    }} />
);

const App = () => {
  return (
      <Router>
          <div>
              {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
              <Switch>
                  <Route path="/login">
                      <LoginPage/>
                  </Route>
                  <Route path="/register">
                      <RegisterPage/>
                  </Route>
                  <Route path="/spotify-auth">
                      <AuthSpotifyPage/>
                  </Route>
                  <PrivateRoute component={HomePage} path="/"/>
              </Switch>
          </div>
      </Router>
  );
};

function LoginPage() {
    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [redirect, setRedirect] = React.useState<boolean>(false);

    function onChangeEmail(event: any) {
        setEmail(event.target.value)
    }

    function onChangePassword(event: any) {
        setPassword(event.target.value)
    }

    async function login() {
        await Api.login(email, password);

        setRedirect(true)
    }
    if (redirect) {
        return <Redirect to={{ pathname: "/" }} />
    }

    return (
        <div style={{
            width: "100%",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flexDirection: "column",
        }}>
            <Typography component={"h1"}>Log in</Typography>
            <TextField required label="Email" placeholder="awesome@mail.com" value={email} onChange={onChangeEmail}/>
            <TextField value={password} onChange={onChangePassword} required label="Password" type="password"/>
            <Button onClick={login}>Login</Button>
            <Link to={"/register"}>Register here</Link>
        </div>
    )
}

function AuthSpotifyPage() {
    function useQuery() {
        return new URLSearchParams(useLocation().search);
    }

    const [redirectTo, setRedirectTo] = React.useState("");

    const query = useQuery();

    const code = query.get("code");
    const state = query.get("state");

    if (!code || ! state) {
        return <p>Invalid params</p>
    }
    if (redirectTo === "") {
        Api.registerSpotify(code, state)
            .then(() => setRedirectTo("/"))
            .catch();
    }
    if (redirectTo !== "") {
        return <Redirect to={{ pathname: "/" }} />
    }

    return (
        <div>Ongoing validation</div>
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

export default App;