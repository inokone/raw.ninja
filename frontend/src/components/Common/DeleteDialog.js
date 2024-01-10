import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from "@mui/material";
import withStyles from '@mui/styles/withStyles';

const styles = theme => ({
    dialogActions: {
        justifyContent: "flex-end",
        paddingTop: theme.spacing(2),
        paddingBottom: theme.spacing(2),
        paddingRight: theme.spacing(2)
    },
    dialog: {
        zIndex: 1400
    },
    backIcon: {
        marginRight: theme.spacing(1)
    }
});

function DeleteDialog({ classes, onDelete, onCancel, open, name }) {

    return (
        <Dialog
            PaperProps={{ style: { overflowY: 'visible', zIndex: 1 } }}
            open={open}
            fullWidth
            scroll="paper"
            onClose={onCancel}
        >
            <DialogTitle></DialogTitle>
            <DialogContent style={{ overflowY: 'visible' }}>
                <Typography align='center'>Are you sure you want to delete {name}?</Typography>
            </DialogContent>
            <DialogActions className={classes.dialogActions}>
                <Button
                    onClick={onCancel}
                    variant="contained"
                    color="secondary"
                >
                    Cancel
                </Button>               
                <Button
                    onClick={onDelete}
                    variant="contained"
                    color="primary"
                >
                    Delete
                </Button>
            </DialogActions>
        </Dialog>
    );
}

DeleteDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onDelete: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired,
    name: PropTypes.string.isRequired
};

export default withStyles(styles, { withTheme: true })(DeleteDialog);