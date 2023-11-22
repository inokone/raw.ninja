import React from 'react';
import { styled } from '@mui/material/styles';

import { Alert, Paper, TableContainer, TableBody, Table, TableHead, TableRow, Box } from "@mui/material";
import TableCell, { tableCellClasses } from '@mui/material/TableCell';

import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

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
    const [loading, setLoading] = React.useState(false)
    const [users, setUsers] = React.useState(null)

    const asDate = (unixTimestamp) => {
        return new Date(unixTimestamp * 1000).toLocaleDateString()
    }

    const loadUsers = () => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/users/', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
                    response.json().then(content => {
                        setUsers(content)
                        setLoading(false)
                    })
                }
            })
            .catch(error => {
                setError(error.message)
                setLoading(false)
            });
    }

    React.useEffect(() => {
        if (!loading && !users && !error) {
            loadUsers()
        }
    },[users, loading, error])

    return (
        <React.Fragment>
            {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {users &&
                <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
                    <TableContainer component={Paper} style={{ flex: 0.85 }}>
                        <Table style={{ flex: 0.85 }}>
                            <TableHead>
                                <TableRow>
                                    <StyledTableCell>E-mail</StyledTableCell>
                                    <StyledTableCell>Name</StyledTableCell>
                                    <StyledTableCell>Role</StyledTableCell>
                                    <StyledTableCell>Registered</StyledTableCell>
                                    <StyledTableCell>Last Updated</StyledTableCell>
                                    <StyledTableCell>Deleted</StyledTableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {users.map((user) => {
                                    return (
                                        <StyledTableRow key={user.id}>
                                            <StyledTableCell>{user.email}</StyledTableCell>
                                            <StyledTableCell>{user.first_name + " " + user.last_name}</StyledTableCell>
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
                </Box>}
        </React.Fragment>
    );
}

export default UserTable