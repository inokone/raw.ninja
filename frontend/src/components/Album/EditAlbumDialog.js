import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Autocomplete, Chip } from "@mui/material";
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

function EditAlbumDialog({ classes, onSave, onCancel, open, input }) {
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
                    sx={{ mb: 4, mt: 1, backgroundColor: "#fff", borderRadius: 1 }}
                />
                <Autocomplete
                    multiple
                    id="tags"
                    options={[].map(a => a)}
                    value={tags}
                    onChange={(event, newValue) => {
                        setTags(newValue);
                    }}
                    variant='outlined'
                    color='primary'
                    freeSolo
                    renderTags={(value, getTagProps) =>
                        value.map((option, index) => (
                            <Chip variant="outlined" label={option} {...getTagProps({ index })} />
                        ))
                    }
                    renderInput={(params) => (
                        <TextField
                            {...params}
                            variant='outlined'
                            color='primary'
                            label="Tags"
                            placeholder="Tags"
                        />
                    )}
                />
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
                    onClick={() => onSave(name, tags)}
                    variant="contained"
                    color="primary"
                >
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
}

EditAlbumDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onSave: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired,
    input: PropTypes.object
};

export default withStyles(styles, { withTheme: true })(EditAlbumDialog);