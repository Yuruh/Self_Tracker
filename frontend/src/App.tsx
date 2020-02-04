import React from 'react';
import './App.css';
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import Api from "./Api";

const App = () => {
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
        <Button onClick={async () => {
            try {
                const res = await Api.getSpotifyUrl();
//                console.log(res.data)
                window.location.replace(res.data);
            } catch (e) {
                console.log(e.message)
            }
        }}>Authorize Spotify</Button>
      </div>
  );
};

export default App;