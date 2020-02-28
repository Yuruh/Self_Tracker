import {withStyles} from "@material-ui/core";
import Box from "@material-ui/core/Box";

const BoxCenter = withStyles(theme => ({
    root: {
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
    },
}))(Box);

export {BoxCenter};
