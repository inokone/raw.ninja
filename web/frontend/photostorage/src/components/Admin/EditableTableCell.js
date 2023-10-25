import React from 'react';
import { styled } from '@mui/material/styles';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import { CircularProgress, TextField, Tooltip } from "@mui/material";
import ReportProblemOutlinedIcon from '@mui/icons-material/ReportProblemOutlined';


const EditableCell = (props) => {
    const [error, setError] = React.useState(null)
    const [working, setWorking] = React.useState(false)
    const [editing, setEditing] = React.useState(false)
    const { onCellEdit, value, formatter } = props

    const handleChange = (event) => {
        setWorking(true)
        onCellEdit(event.target.value)
        .then(result => {
            setWorking(false)
        }).catch(error => {
            setWorking(false)
            setError(error.message)
        })
    };

    return (
        <TableCell onClick={() => !editing && !working ? setEditing(true) : setError(null)} 
            sx={{ ...(error !== null) && {border: 1, borderColor: 'red'} }}>
            {working && (
                <CircularProgress sx={{ mr: 2 }} size={20} thickness={2} />
            )}
            {error && (
                <Tooltip title={error}>
                    <ReportProblemOutlinedIcon sx={{ color: 'red' }}/>
                </Tooltip>
            )}
            {editing ? (
                <TextField
                    value={value}
                    onChange={handleChange}
                    onBlur={() => setEditing(false)}
                />
            ) : (
                formatter(value)
            )}
        </TableCell>
    )
}

const EditableTableCell = styled(EditableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

export default EditableTableCell