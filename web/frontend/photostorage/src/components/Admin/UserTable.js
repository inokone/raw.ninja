import React from 'react';
import { styled } from '@mui/material/styles';

import { Alert, Paper, TableContainer, TableBody, Table, TableHead, TableRow, Box } from "@mui/material";
import TableCell, { tableCellClasses } from '@mui/material/TableCell';

import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env;

const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
        border: 0,
    },
}));

const UserTable = () => {
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(true)
    const [users, setUsers] = React.useState(null)

    const asDate = (unixTimestamp) => {
        return new Date(unixTimestamp * 1000).toLocaleDateString()
    }

    React.useEffect(() => {
        const loadUsers = () => {
            fetch(REACT_APP_API_PREFIX + '/api/v1/users/', {
                method: "GET",
                mode: "cors",
                credentials: "include"
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(response.status + ": " + response.statusText);
                    } else {
                        response.json().then(content => {
                            setLoading(false)
                            setUsers(content)
                        })
                    }
                })
                .catch(error => {
                    setError(error)
                    setLoading(false)
                });
        }

        if (!users) {
            loadUsers()
        }
    },)

    return (
        <>
            {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
            {loading ? <ProgressDisplay /> :
                <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
                    <TableContainer component={Paper} style={{ width: 1200 }}>
                        <Table style={{ width: 1200 }}>
                            <TableHead>
                                <TableRow>
                                    <StyledTableCell>E-mail</StyledTableCell>
                                    <StyledTableCell>Phone</StyledTableCell>
                                    <StyledTableCell>Role</StyledTableCell>
                                    <StyledTableCell>Registered</StyledTableCell>
                                    <StyledTableCell>Last Updated</StyledTableCell>
                                    <StyledTableCell>Deleted</StyledTableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {users.map((user) => {
                                    return (
                                        <StyledTableRow>
                                            <StyledTableCell>{user.email}</StyledTableCell>
                                            <StyledTableCell>{user.phone}</StyledTableCell>
                                            <StyledTableCell>{user.role.name}</StyledTableCell>
                                            <StyledTableCell>{asDate(user.created)}</StyledTableCell>
                                            <StyledTableCell>{asDate(user.updated)}</StyledTableCell>
                                            <StyledTableCell>{user.deleted}</StyledTableCell>
                                        </StyledTableRow>
                                    )
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Box>
            }
        </>
    );
}

export default UserTable