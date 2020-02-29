import {List} from "@material-ui/core";
import ListItem from "@material-ui/core/ListItem";
import ListItemAvatar from "@material-ui/core/ListItemAvatar";
import Avatar from "@material-ui/core/Avatar";
import ListItemText from "@material-ui/core/ListItemText";
import React from "react"
import ListSubheader from "@material-ui/core/ListSubheader";
import ListItemSecondaryAction from "@material-ui/core/ListItemSecondaryAction";
import Button from "@material-ui/core/Button";
import Api from "./Api";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CancelIcon from '@material-ui/icons/Cancel';

async function connectSpotify() {
    try {
        const res = await Api.getSpotifyUrl();
        window.location.replace(res.data.url);
    } catch (e) {
        console.log(e.message)
    }
}

export default function ConnectorList() {
    return <List subheader={<ListSubheader>Connectors</ListSubheader>}>
        <Connector avatarSrc={"https://cdn0.capterra-static.com/logos/150/2137143-1574690999.png"}
                   isConnected={true}
                   title={"Affect-tag"}
                   onConnect={() => console.log("connect aftg")}
        />
        <Connector title={"Spotify"}
                   onConnect={connectSpotify}
                   isConnected={false}
                   avatarSrc={"https://upload.wikimedia.org/wikipedia/fr/6/60/Spotify_logo_sans_texte.png"}/>
    </List>
}

function Connector(props: {
    avatarSrc: string,
    title: string,
    isConnected: boolean,
    onConnect: () => void,
}) {
    return <ListItem>
        <ListItemIcon>
            {props.isConnected ? <CheckCircleIcon style={{color: "green"}}/> : <CancelIcon style={{color: "red"}}/>}
        </ListItemIcon>
        <ListItemAvatar>
            <Avatar src={props.avatarSrc}/>
        </ListItemAvatar>
        <ListItemText primary={props.title}/>
        <ListItemSecondaryAction>
            {/*props.isConnected && <Switch color="primary"/>*/}
            <Button onClick={props.onConnect}>{props.isConnected ? "Reconnect" : "Connect"}</Button>
        </ListItemSecondaryAction>
    </ListItem>

}