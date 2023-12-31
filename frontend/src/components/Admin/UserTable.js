import { useMemo, useState } from 'react';
import {
    MaterialReactTable,
    useMaterialReactTable,
} from 'material-react-table';
import { Box, IconButton, Tooltip } from '@mui/material';
import {
    QueryClient,
    QueryClientProvider,
    useMutation,
    useQuery,
    useQueryClient,
} from '@tanstack/react-query';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import LockPersonIcon from '@mui/icons-material/LockPerson';
import LockOpenIcon from '@mui/icons-material/LockOpen';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const UserTableTemplate = () => {
    const [validationErrors, setValidationErrors] = useState({});

    const roles = useMemo(() => [{
        label: "Admin",
        value: {
            id: 0,
            name: "Admin",
        }
    },
    {
        label: "Free Tier",
        value: {
            id: 1,
            name: "Free Tier",
        }
    }], [])

    const columns = useMemo(
        () => [
            {
                accessorKey: 'id',
                header: 'Id',
                enableEditing: false,
                size: 0,
            },
            {
                accessorKey: 'first_name',
                header: 'First Name',
                muiEditTextFieldProps: {
                    type: 'email',
                    required: true,
                    error: !!validationErrors?.firstName,
                    helperText: validationErrors?.firstName,
                    //remove any previous validation errors when user focuses on the input
                    onFocus: () =>
                        setValidationErrors({
                            ...validationErrors,
                            firstName: undefined,
                        }),
                    //optionally add validation checking for onBlur or onChange
                },
            },
            {
                accessorKey: 'last_name',
                header: 'Last Name',
                muiEditTextFieldProps: {
                    type: 'email',
                    required: true,
                    error: !!validationErrors?.lastName,
                    helperText: validationErrors?.lastName,
                    //remove any previous validation errors when user focuses on the input
                    onFocus: () =>
                        setValidationErrors({
                            ...validationErrors,
                            lastName: undefined,
                        }),
                },
            },
            {
                accessorKey: 'email',
                header: 'Email',
                enableEditing: false,
                muiEditTextFieldProps: {
                    type: 'email',
                    required: true,
                    error: !!validationErrors?.email,
                    helperText: validationErrors?.email,
                    //remove any previous validation errors when user focuses on the input
                    onFocus: () =>
                        setValidationErrors({
                            ...validationErrors,
                            email: undefined,
                        }),
                },
            },
            {
                accessorKey: 'created',
                header: 'Created',
                enableEditing: false,
                sortingFn: 'datetime',
                Cell: ({ cell }) => new Date(cell.getValue() * 1000).toLocaleDateString(),  //render Date as a string
            },
            {
                accessorKey: 'role',
                header: 'Role',
                editVariant: 'select',
                editSelectOptions: roles,
                enableEditing: false, // when backend is fixed, it can be enabled
                muiEditTextFieldProps: {
                    select: true
                },
                Cell: ({ cell }) => cell.getValue().name  //render Role name only 
            },
        ],
        [validationErrors, roles]
    );

    const {
        data: users = [],
        isError: loadError,
        isFetching,
        isLoading,
    } = usePopulate();
    const { mutateAsync: updateUser, isPending: isUpdating } = useUpdate();
    const { mutateAsync: deleteUser, isPending: isDeleting } = useDelete();
    const { mutateAsync: toggleEnabledUser, isPending: isTogglingEnableUser } = useToggleEnabled();


    const handleUpdate = async ({ values, table }) => {
        const newValidationErrors = validateUser(values);
        if (Object.values(newValidationErrors).some((error) => error)) {
            setValidationErrors(newValidationErrors);
            return;
        }
        setValidationErrors({});
        await updateUser(values);
        table.setEditingRow(null);
    };

    const handleDelete = (row) => {
        if (window.confirm('Are you sure you want to delete this user?')) {
            deleteUser(row.original.id);
        }
    };


    const handleToggleEnabled = (row) => {
        if (row.original.enabled === true) {
            if (window.confirm('Are you sure you want to disable this user?')) {
                toggleEnabledUser(row.original);
            }
        } else {
            if (window.confirm('Are you sure you want to enable this user?')) {
                toggleEnabledUser(row.original);
            }
        }
    };

    const table = useMaterialReactTable({
        columns,
        data: users,
        createDisplayMode: 'row',
        editDisplayMode: 'row',
        enableEditing: true,
        enableFullScreenToggle: false,
        enableDensityToggle: false,
        getRowId: (row) => row.id,
        muiToolbarAlertBannerProps: loadError
            ? {
                color: 'error',
                children: 'Error loading data',
            }
            : undefined,
        muiTableContainerProps: {
            sx: {
                minHeight: '300px'
            },
        },
        onCreatingRowCancel: () => setValidationErrors({}),
        onEditingRowCancel: () => setValidationErrors({}),
        onEditingRowSave: handleUpdate,
        renderRowActions: ({ row, table }) => (
            <Box sx={{ display: 'flex', gap: '1rem' }}>
                <Tooltip title="Edit">
                    <IconButton onClick={() => table.setEditingRow(row)}>
                        <EditIcon />
                    </IconButton>
                </Tooltip>
                {row.original.enabled === true ?
                    <Tooltip title="Disable">
                        <IconButton onClick={() => handleToggleEnabled(row)}>
                            <LockPersonIcon />
                        </IconButton>
                    </Tooltip> :
                    <Tooltip title="Enable">
                        <IconButton color="error" onClick={() => handleToggleEnabled(row)}>
                            <LockOpenIcon />
                        </IconButton>
                    </Tooltip>
                }
                <Tooltip title="Delete">
                    <IconButton onClick={() => handleDelete(row)}>
                        <DeleteIcon />
                    </IconButton>
                </Tooltip>
            </Box>
        ),
        initialState: {
            density: 'compact'
        },
        state: {
            isLoading: isLoading,
            isSaving: isUpdating || isDeleting || isTogglingEnableUser,
            showAlertBanner: loadError,
            showProgressBars: isFetching,
        },
    });

    return <MaterialReactTable table={table} />;
};

function usePopulate() {
    return useQuery({
        queryKey: ['users'],
        queryFn: async () => {
            //send api request here
            return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/users/', {
                method: "GET",
                mode: "cors",
                credentials: "include"
            })
                .then(response => {
                    if (!response.ok) {
                        response.json().then(content => {
                            reject(content.message)
                        });
                    } else {
                        response.json().then(content => {
                            resolve(content)
                        })
                    }
                })
                .catch(error => {
                    reject(error.message)
                }));
        },
        refetchOnWindowFocus: false,
    });
}

function useUpdate() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (user) => {
            return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/users/' + user.id, {
                method: "PATCH",
                mode: "cors",
                credentials: "include",
                body: JSON.stringify({
                    id: user.id,
                    first_name: user.first_name,
                    last_name: user.last_name,
                    role: user.role
                })
            })
                .then(response => {
                    if (!response.ok) {
                        response.json().then(content => {
                            reject(content.message)
                        });
                    } else {
                        response.json().then(content => {
                            resolve(content)
                        })
                    }
                })
                .catch(error => {
                    reject(error.message)
                }));
        },
        //client side optimistic update
        onMutate: (newUserInfo) => {
            queryClient.setQueryData(['users'], (prevUsers) =>
                prevUsers?.map((prevUser) =>
                    prevUser.id === newUserInfo.id ? newUserInfo : prevUser,
                ),
            );
        }
    });
}

function useToggleEnabled() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (user) => {
            user.enabled = user.enabled === true ? false : true
            return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/users/' + user.id + '/enabled', {
                method: "PUT",
                mode: "cors",
                credentials: "include",
                body: JSON.stringify({
                    id: user.id,
                    enabled: user.enabled
                })
            })
                .then(response => {
                    if (!response.ok) {
                        response.json().then(content => {
                            reject(content.message)
                        });
                    } else {
                        response.json().then(content => {
                            resolve(content)
                        })
                    }
                })
                .catch(error => {
                    reject(error.message)
                }));
        },
        //client side optimistic update
        onMutate: (newUserInfo) => {
            queryClient.setQueryData(['users'], (prevUsers) =>
                prevUsers?.map((prevUser) =>
                    prevUser.id === newUserInfo.id ? newUserInfo : prevUser,
                ),
            );
        }
    });
}

function useDelete() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (userId) => {
            return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/users/' + userId, {
                method: "DELETE",
                mode: "cors",
                credentials: "include"
            })
                .then(response => {
                    if (!response.ok) {
                        response.json().then(content => {
                            reject(content.message)
                        });
                    } else {
                        response.json().then(content => {
                            resolve(content)
                        })
                    }
                })
                .catch(error => {
                    reject(error.message)
                }));
        },
        //client side optimistic update
        onMutate: (userId) => {
            queryClient.setQueryData(['users'], (prevUsers) =>
                prevUsers?.filter((user) => user.id !== userId),
            );
        }
    });
}

const queryClient = new QueryClient();

const UserTable = () => (
    //Put this with your other react-query providers near root of your app
    <QueryClientProvider client={queryClient}>
        <UserTableTemplate />
    </QueryClientProvider>
);

export default UserTable;

const validateRequired = (value) => !!value.length;
const validateEmail = (email) => {
    return !!email.length && email.toLowerCase().match(
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
    );
}

function validateUser(user) {
    return {
        firstName: !validateRequired(user.first_name)
            ? 'First Name is Required'
            : '',
        lastName: !validateRequired(user.last_name) ? 'Last Name is Required' : '',
        email: !validateEmail(user.email) ? 'Incorrect Email Format' : '',
    };
}
