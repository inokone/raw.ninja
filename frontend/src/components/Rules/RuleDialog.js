import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, InputLabel, Select, MenuItem, FormLabel } from "@mui/material";
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

const timingUnits = [
    {
        label: "Day",
        value: 1,
    },
    {
        label: "Week",
        value: 7,
    },
    {
        label: "Month (30 days)",
        value: 30,
    },
    {
        label: "Year",
        value: 365,
    }
]

const timingFor = (duration) => {
    if (!duration) {
        return ""
    }
    let unit = timingUnits.slice().reverse().filter((unit) => duration % unit.value === 0)[0]
    return duration / unit.value
}

const timingUnitFor = (duration) => {
    if (!duration) {
        return ""
    }
    return timingUnits.slice().reverse().filter((unit) => duration % unit.value === 0)[0].value
}

function RuleDialog({ classes, onSave, onCancel, open, input, constants }) {
    const [name, setName] = React.useState(input ? input.name : "")
    const [nameError, setNameError] = React.useState(null)
    const [timing, setTiming] = React.useState(input ? timingFor(input.timing) : 1)
    const [timingError, setTimingError] = React.useState(null)
    const [timingUnit, setTimingUnit] = React.useState(input ? timingUnitFor(input.timing) : 1)
    const [actionID, setActionID] = React.useState(input && input.action ? input.action.id : 1)
    const [actionTargetID, setActionTargetID] = React.useState(input && input.target ? input.target.id : 0)
    const [isTargeted, setIsTargeted] = React.useState(input && input.action && input.action.targeted)

    const createRule = () => {
        let timingValue = timing * timingUnit
        let action = constants.actions.filter((act) => act.id === actionID)[0]
        let actionTarget = actionTargetID && constants.targets.filter((tgt) => tgt.id === actionTargetID)[0]
        let desc = "After " + timingValue + " days " + action.name.toLowerCase()
        if (isTargeted) {
            desc += (" " + actionTarget.name.toLowerCase())
        }
        desc += "."
        return {
            id: input ? input.id : null,
            name: name,
            description: desc,
            timing: timingValue,
            action: {
                id: actionID
            },
            target: actionTargetID ? {
                id: actionTargetID
            }: null
        }
    }

    const handleActionChange = React.useCallback((id) => {
        setIsTargeted(id && constants.actions.filter((act) => act.id === id)[0].targeted)
        setActionID(id)
    }, [constants])

    return (
        <Dialog
            PaperProps={{ style: { overflowY: 'visible', zIndex: 1 } }}
            open={open}
            fullWidth
            scroll="paper"
            onClose={onCancel}
        >
            <DialogTitle>{input ? "Change" : "Create"} rule</DialogTitle>
            <DialogContent style={{ overflowY: 'visible' }}>
                <TextField
                    type="text"
                    name="rule"
                    autoComplete="rule"
                    variant='outlined'
                    color='primary'
                    label="Name"
                    value={name}
                    error={nameError}
                    onChange={e => {
                        setNameError(e.target.value === '')
                        setName(e.target.value)
                    }}
                    fullWidth
                    required
                    sx={{ mb: 4, mt: 1, backgroundColor: "#fff", borderRadius: 1 }}
                />
                <TextField
                    fullWidth
                    type="number"
                    name="timing"
                    autoComplete="timing"
                    variant='outlined'
                    color='primary'
                    label="Timing"
                    value={timing}
                    error={timingError}
                    InputProps={{ inputProps: { min: 1, max: 31 } }}
                    onChange={e => {
                        setTimingError(e.target.value === '')
                        setTiming(e.target.value)
                    }}
                    required
                    sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                />
                <FormLabel id="timing-unit-helper-label">Timing Unit</FormLabel>
                <Select
                    fullWidth
                    labelId="timing-unit-helper-label"
                    id="timing-unit-select"
                    value={timingUnit}
                    onChange={e => {
                        setTimingUnit(e.target.value)
                    }}
                >
                    {timingUnits.map(({ value, label }, index) => (<MenuItem key={index} value={value}>{label}</MenuItem>))}
                </Select>
                <InputLabel id="action-helper-label">Action</InputLabel>
                <Select
                    fullWidth
                    labelId="action-helper-label"
                    id="action-select"
                    value={actionID}
                    onChange={e => {
                        handleActionChange(e.target.value)
                    }}
                >
                    {constants && constants.actions && constants.actions.map(({ id, name }, index) => (<MenuItem key={index} value={id}>{name}</MenuItem>))}
                </Select>
                {isTargeted &&
                <>
                    <InputLabel id="action-target-helper-label">Target</InputLabel>
                    <Select
                        fullWidth
                        labelId="action-target-helper-label"
                        id="action-target-select"
                        value={actionTargetID}
                        onChange={e => {
                            setActionTargetID(e.target.value)
                        }}
                    >
                        {constants.targets.map(({ id, name }, index) => (<MenuItem key={index} value={id}>{name}</MenuItem>))}
                    </Select>
                </>
                }
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
                    onClick={() => onSave(createRule())}
                    variant="contained"
                    color="primary"
                >
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
}

RuleDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onSave: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired,
    constants: PropTypes.object.isRequired,
    input: PropTypes.object
};

export default withStyles(styles, { withTheme: true })(RuleDialog);