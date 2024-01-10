import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField} from "@mui/material";
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

function RuleSetDialog({ classes, onSave, onCancel, open, input }) {
    const [name, setName] = React.useState(input ? input.name : null)
    const [nameError, setNameError] = React.useState(null)
    const [description, setDescription] = React.useState(input ? input.description : null)

    return (
        <Dialog
            open={open}
            scroll="paper"
            onClose={onCancel}
            className={classes.dialog}
        >
            <DialogTitle>{input ? "Edit" : "Create"} rule set</DialogTitle>
            <DialogContent>
                <TextField
                    type="text"
                    name="ruleset"
                    autoComplete="ruleset"
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
                <TextField
                    type="text"
                    name="rulesetdescription"
                    autoComplete="rulesetdescription"
                    variant='outlined'
                    color='primary'
                    label="Description"
                    value={description}
                    onChange={e => {
                        setDescription(e.target.value)
                    }}
                    fullWidth
                    required
                    sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
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
                    onClick={() => onSave({
                        name: name,
                        description: description
                    })}
                    variant="contained"
                    color="primary"
                >
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    )
}

RuleSetDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onSave: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired
};

export default withStyles(styles, { withTheme: true })(RuleSetDialog);