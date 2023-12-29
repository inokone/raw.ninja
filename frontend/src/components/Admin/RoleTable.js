import { useMemo, useState } from 'react';
import {
    MaterialReactTable,
    useMaterialReactTable,
} from 'material-react-table';
import {
    useQuery,
    QueryClient,
    QueryClientProvider,
    useQueryClient,
    useMutation
} from '@tanstack/react-query';
import { Box, IconButton, Tooltip } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const RoleTableTemplate = () => {
    const [validationErrors, setValidationErrors] = useState({});

    const formatBytes = (bytes, decimals = 2) => {
        if (!+bytes) return '0 Bytes'

        const k = 1024
        const dm = decimals < 0 ? 0 : decimals
        const sizes = ['Bytes', 'KBytes', 'MBytes', 'GBytes', 'TBytes', 'PBytes']

        const i = Math.floor(Math.log(bytes) / Math.log(k))

        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    }

    const usePopulate = () => {
        return useQuery({
            queryKey: ['roles'],
            queryFn: async () => {
                return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/roles/', {
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
            mutationFn: async (role) => {
                return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/roles/' + role.id, {
                    method: "PUT",
                    mode: "cors",
                    credentials: "include",
                    body: JSON.stringify({
                        id: role.id,
                        name: role.name,
                        quota: parseInt(role.quota)
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
            onMutate: (newRoleInfo) => {
                queryClient.setQueryData(['roles'], (prevRoles) =>
                    prevRoles?.map((prevRole) =>
                        prevRole.id === newRoleInfo.id ? newRoleInfo : prevRole,
                    ),
                );
            },
        });
    }

    function useDelete() {
        const queryClient = useQueryClient();
        return useMutation({
            mutationFn: async (roleId) => {
                return await new Promise((resolve, reject) => fetch(REACT_APP_API_PREFIX + '/api/v1/roles/' + roleId, {
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
            onMutate: (roleId) => {
                queryClient.setQueryData(['roles'], (prevRoles) =>
                    prevRoles?.filter((role) => role.id !== roleId),
                );
            },
        });
    }

    const handleUpdate = async ({ values, table }) => {
        const newValidationErrors = validateRole(values);
        if (Object.values(newValidationErrors).some((error) => error)) {
            setValidationErrors(newValidationErrors);
            return;
        }
        setValidationErrors({});
        await updateRole(values);
        table.setEditingRow(null);
    };

    const handleDelete = (row) => {
        if (window.confirm('Are you sure you want to delete this role?')) {
            deleteRole(row.original.id);
        }
    };

    const columns = useMemo(
        () => [
            {
                accessorKey: 'id',
                header: 'ID',
                size: 50,
                enableEditing: false
            },
            {
                accessorKey: 'name',
                header: 'Name',
                size: 200,
                error: !!validationErrors?.name,
                helperText: validationErrors?.name,
                onFocus: () =>
                    setValidationErrors({
                        ...validationErrors,
                        name: undefined,
                    }),
                enableEditing: false
            },
            {
                accessorKey: 'quota',
                header: 'Quota',
                Cell: ({ cell }) => cell.getValue() <= 0 ? 'Unlimited' : formatBytes(cell.getValue()),
                size: 100,
                error: !!validationErrors?.quota,
                helperText: validationErrors?.quota,
                onFocus: () =>
                    setValidationErrors({
                        ...validationErrors,
                        quota: undefined,
                    }),
            },
        ],
        [validationErrors],
    );

    const {
        data: roles = [],
        isError: loadError,
        isFetching,
        isLoading,
    } = usePopulate();
    const { mutateAsync: updateRole, isPending: isUpdating } = useUpdate();
    const { mutateAsync: deleteRole, isPending: isDeleting } = useDelete();

    const table = useMaterialReactTable({
        columns,
        data: roles,
        enableFullScreenToggle: false,
        enableDensityToggle: false,
        createDisplayMode: 'row',
        editDisplayMode: 'row',
        enableEditing: true,
        muiToolbarAlertBannerProps: loadError
            ? {
                color: 'error',
                children: 'Error loading data',
            }
            : undefined,
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
                <Tooltip title="Delete">
                    <IconButton color="error" onClick={() => handleDelete(row)}>
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
            isSaving: isUpdating || isDeleting,
            showAlertBanner: loadError,
            showProgressBars: isFetching,
        },
    });

    return <MaterialReactTable table={table} />;
};

const queryClient = new QueryClient();

const RoleTable = () => (
    //Put this with your other react-query providers near root of your app
    <QueryClientProvider client={queryClient}>
        <RoleTableTemplate />
    </QueryClientProvider>
);

export default RoleTable;

const validateRequired = (value) => !!value.length;

function validateRole(role) {
    return {
        name: !validateRequired(role.name) ? 'Role Name is Required' : '',
        quota: !validateRequired(role.quota) ? 'Quota is Required' : '',
    };
}
