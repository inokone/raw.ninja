import React from "react";
import PropTypes from "prop-types";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, FormControl, InputLabel, Select, FormHelperText, MenuItem } from "@mui/material";
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
        return null
    }
    timingUnits.reverse().forEach((unit) => {
        if (duration % unit.value === 0) {
            return duration / unit.value
        }
    })
}

const timingUnitFor = (duration) => {
    if (!duration) {
        return null
    }
    timingUnits.reverse().forEach((unit) => {
        if (duration % unit.value === 0) {
            return unit
        }
    })
}

function RuleDialog({ classes, onSave, onCancel, open, input, constants }) {
    const [name, setName] = React.useState(input ? input.name : null)
    const [nameError, setNameError] = React.useState(null)
    const [timing, setTiming] = React.useState(input ? timingFor(input.timing) : null)
    const [timingError, setTimingError] = React.useState(null)
    const [timingUnit, setTimingUnit] = React.useState(input ? timingUnitFor(input.timing) : null)
    const [action, setAction] = React.useState(input ? input.action : null)
    const [actionTarget, setActionTarget] = React.useState(input ? input.actionTarget : null)

    const options = constants || {
        actions: [
            {
                value: 1,
                label: "Delete",
                targeted: false
            },
            {
                value: 2,
                label: "Move to",
                targeted: true
            },
        ],
        targets: [
            {
                value: 1,
                label: "Standard storage"
            },
            {
                value: 2,
                label: "Frozen storage"
            },
            {
                value: 3,
                label: "Bin"
            }
        ]
    }

    const createRule = () => {
        let timingValue = timing * timingUnit
        return {
            name: name,
            description: "After " + timingValue + " days " + action.label.toLower() + (actionTarget && " " + actionTarget.label.toLower()),
            timing: timingValue,
            action: action.value,
            target: actionTarget.value
        }
    }

    return (
        <Dialog
            open={open}
            scroll="paper"
            onClose={onCancel}
            className={classes.dialog}
        >
            <DialogTitle>{input ? "Change" : "Create"} rule</DialogTitle>
            <DialogContent>
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
                        setNameError(name === '')
                        setName(e.target.value)
                    }}
                    fullWidth
                    required
                    sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                />
                <FormControl sx={{ m: 1, minWidth: 120 }}>
                    <TextField
                        type="text"
                        name="timing"
                        autoComplete="timing"
                        variant='outlined'
                        color='primary'
                        label="Timing"
                        value={timing}
                        error={timingError}
                        onChange={e => {
                            setTimingError(timing === '')
                            setTiming(e.target.value)
                        }}
                        fullWidth
                        required
                        sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                    />
                    <FormHelperText>When the action should happen</FormHelperText>
                </FormControl>
                <FormControl sx={{ m: 1, minWidth: 120 }}>
                    <InputLabel id="timing-unit-helper-label">Timing</InputLabel>
                    <Select
                        labelId="timing-unit-helper-label"
                        id="timing-unit-select"
                        value={timingUnit}
                        label="Timing Unit"
                        onChange={e => {
                            setTimingUnit(e.target.value)
                        }}
                    >
                        {timingUnits.map(unit => {
                            return (<MenuItem key={unit.value + unit.label} value={unit.value}>
                                <em>{unit.label}</em>
                            </MenuItem>)
                        })}
                    </Select>
                </FormControl>
                <FormControl sx={{ m: 1, minWidth: 120 }}>
                    <InputLabel id="action-helper-label">Action</InputLabel>
                    <Select
                        labelId="action-helper-label"
                        id="action-select"
                        value={action}
                        label="Action"
                        onChange={e => {
                            setAction(e.target.value)
                        }}
                    >
                        {options.actions.map(action => {
                            return (<MenuItem key={action.value + action.label} value={action}>
                                <em>{action.label}</em>
                            </MenuItem>)
                        })}
                    </Select>
                    <FormHelperText>What action should happen</FormHelperText>
                </FormControl>
                {action && action.targeted &&
                    <FormControl sx={{ m: 1, minWidth: 120 }}>
                        <InputLabel id="action-target-helper-label">Target</InputLabel>
                        <Select
                            labelId="action-target-helper-label"
                            id="action-target-select"
                            value={actionTarget}
                            label="Target"
                            onChange={e => {
                                setActionTarget(e.target.value)
                            }}
                        >
                            {options.actions.map(target => {
                                return (<MenuItem key={target.value + target.label} value={target}>
                                    <em>{target.label}</em>
                                </MenuItem>)
                            })}
                        </Select>
                        <FormHelperText>What is the target of the action</FormHelperText>
                    </FormControl>
                }
            </DialogContent>
            <DialogActions className={classes.dialogActions}>
                <Button
                    onClick={() => onSave(createRule())}
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

RuleDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    onSave: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired
};

export default withStyles(styles, { withTheme: true })(RuleDialog);