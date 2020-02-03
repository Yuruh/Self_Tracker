import React from 'react';
import './App.css';
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";

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
      </div>
  );
};

export default App;