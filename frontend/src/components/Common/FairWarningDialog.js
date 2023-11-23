import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Typography } from "@mui/material";
import { Link } from "react-router-dom"
import withStyles from '@mui/styles/withStyles';
import ColoredButton from "./ColoredButton";

const styles = theme => ({
    dialogActions: {
        paddingTop: theme.spacing(2),
        paddingBottom: theme.spacing(2),
        paddingRight: theme.spacing(2)
    },
    dialog: {
        zIndex: 1400
    }
});

function FairWarningDialog(props) {
    const { classes, onClose, open, theme } = props;
    return (
        <Dialog
            open={open}
            scroll="paper"
            onClose={onClose}
            className={classes.dialog}
        >
            <DialogTitle>Fair Warning</DialogTitle>
            <DialogContent>
                <Typography variant="h6" color="primary" paragraph>
                   State of this application
                </Typography>
                <Typography paragraph>
                    This current version of the application is used for development and 
                     <strong> BETA testing</strong>. We are on a journey for creating new features and 
                    identifying existing issues in the application - including potential 
                    security related issues.
                </Typography>
                <Typography variant="h6" color="primary" paragraph>
                    What you can expect
                </Typography>
                <Typography paragraph>
                    We can <strong>wipe all data</strong> from the application at any given time without
                    any prior notice. We can not take any responsibility for your loss
                    as stated in the <Link to="/terms">terms of use</Link>.
                    Also, the current environment has rather strict limitations on 
                    computation and reliability, so the application can be slow or
                    unavailable for various reasons.
                    Only free tier is available at this point.
                </Typography>
                <Typography variant="h6" color="primary" paragraph>
                    Why is it public, then?
                </Typography>
                <Typography paragraph>
                    First, we really value your opinion, so please tell us at <strong>support@raw.ninja</strong>.
                    Second, it is easier to test on various devices on a life-like
                    environment.
                </Typography>
            </DialogContent>
            <DialogActions className={classes.dialogActions}>
                <ColoredButton
                    onClick={onClose}
                    variant="contained"
                    color={theme.palette.common.black}
                >                   
                Close
                </ColoredButton>
            </DialogActions>
        </Dialog>
    );
}

FairWarningDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onClose: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired,
    theme: PropTypes.object.isRequired
};

export default withStyles(styles, { withTheme: true })(FairWarningDialog);