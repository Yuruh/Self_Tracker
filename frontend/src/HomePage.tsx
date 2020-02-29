import React, {ChangeEvent} from "react";
import ConnectorList from "./ConnectorList";
import Header from "./Header";
import Switch from "@material-ui/core/Switch";
import {createStyles, FormControlLabel, Theme} from "@material-ui/core";
import Container from "@material-ui/core/Container";
import Tooltip from "@material-ui/core/Tooltip";
import makeStyles from "@material-ui/core/styles/makeStyles";
import Api from "./Api";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        tooltip: {
            fontSize: 16,
        },
    }),
);

async function changeRecording(elem: ChangeEvent, checked: boolean) {
    try {
        await Api.recordActivty(checked);
    } catch (e) {
        console.log(e);
    }
}

export default function HomePage() {
    const classes = useStyles({});

    return <div>
        <Header/>
        <Container style={{width: "100%", marginTop: 30}}>
            <Tooltip classes={{tooltip: classes.tooltip}} arrow title={"When you are recording, contextual data from each connector is retrieved and sent to Affect-tag RX as an emotionally analyzable tag"}>
            <FormControlLabel
                control={<Switch color="primary" onChange={changeRecording}/>}
                label="Recording Activity"
            />
            </Tooltip>
            <div style={{maxWidth: 500}}>
                <ConnectorList/>
            </div>
        </Container>
    </div>
}
