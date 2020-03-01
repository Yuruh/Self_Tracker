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
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import TextField from "@material-ui/core/TextField";
import DialogActions from "@material-ui/core/DialogActions";

async function connectSpotify() {
    try {
        const res = await Api.getSpotifyUrl();
        window.location.replace(res.data.url);
    } catch (e) {
        console.log(e.message)
    }
}

async function connectAftg(key: string) {
    try {
        const res = await Api.registerAftg(key);
    } catch (e) {
        console.log(e.message)
    }
}

export default function ConnectorList() {
    const [open, setOpen] = React.useState(false);
    const [key, setKey] = React.useState("");

    const handleClickOpen = () => {
        setOpen(true);
    };

    const onAftgIntegrate = () => {
        connectAftg(key).then(() => handleClose())
            .catch((e) => console.log(e));
    };

    const handleClose = () => {
        setOpen(false);
    };

    return <div>
        <List subheader={<ListSubheader>Connectors</ListSubheader>}>
            <Connector avatarSrc={"https://cdn0.capterra-static.com/logos/150/2137143-1574690999.png"}
                       isConnected={true}
                       title={"Affect-tag"}
                       onConnect={handleClickOpen}
            />
            <Connector title={"Spotify"}
                       onConnect={connectSpotify}
                       isConnected={false}
                       avatarSrc={"https://upload.wikimedia.org/wikipedia/fr/6/60/Spotify_logo_sans_texte.png"}/>
        </List>
        <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
            <DialogTitle id="form-dialog-title">Affect Tag Integration</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    To integrate with your affect-tag accout, please enter your API Key provided in your "Account"
                    Tab on the Affect-tag cloud platform.
                </DialogContentText>
                <TextField
                    autoFocus
                    margin="dense"
                    id="api key"
                    label="API Key"
                    type="text"
                    value={key}
                    onChange={(e) => setKey(e.target.value)}
                    fullWidth
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} color="primary">
                    Cancel
                </Button>
                <Button onClick={onAftgIntegrate} color="primary">
                    Integrate
                </Button>
            </DialogActions>
        </Dialog>
    </div>
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