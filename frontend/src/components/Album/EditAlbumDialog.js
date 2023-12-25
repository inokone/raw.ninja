import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField } from "@mui/material";
import withStyles from '@mui/styles/withStyles';

const styles = theme => ({
    dialogActions: {
        justifyContent: "flex-start",
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

function EditAlbumDialog(props) {
    const { classes, onSave, onCancel, open, input } = props;
    const [name, setName] = React.useState(input.name)
    const [nameError, setNameError] = React.useState(false)
    const [tags, setTags] = React.useState(input.tags)

    return (
        <Dialog
            open={open}
            scroll="paper"
            onClose={onCancel}
            className={classes.dialog}
        >
            <DialogTitle>Edit album</DialogTitle>
            <DialogContent>
                <TextField
                    type="text"
                    name="album"
                    autoComplete="album"
                    variant='outlined'
                    color='primary'
                    label="Name"
                    value={name}
                    error={nameError}
                    onChange={e => {
                        setNameError(name === '')
                        setName(e.target.value)
                    }}
                    fullWidth
                    required
                    sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                />
                <TextField
                    type="text"
                    name="tags"
                    variant='outlined'
                    color='primary'
                    label="Tags"
                    value={tags}
                    onChange={e => {
                        setTags(e.target.value)
                    }}
                    fullWidth
                    sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                />
            </DialogContent>
            <DialogActions className={classes.dialogActions}>
                <Button
                    onClick={() => onSave(name, tags)}
                    variant="contained"
                    color="primary"
                >
                    Save
                </Button>
                <Button
                    onClick={onCancel}
                    variant="contained"
                    color="secondary"
                >
                    Cancel
                </Button>
            </DialogActions>
        </Dialog>
    );
}

EditAlbumDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onSave: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired
};

export default withStyles(styles, { withTheme: true })(EditAlbumDialog);