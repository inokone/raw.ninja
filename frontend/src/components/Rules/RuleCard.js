import * as React from 'react';
import PropTypes from "prop-types";
import { Card, Typography, CardContent, CardActions, Button } from "@mui/material";
import RuleDialog from './RuleDialog';
import { useTheme } from '@mui/styles';

const RuleCard = ({ rule, setRule, constants, onDelete }) => {
    const theme = useTheme();
    const [isEditDialogOpen, setIsEditDialogOpen] = React.useState(false);

    const handleEditDialogOpen = React.useCallback(() => {
        setIsEditDialogOpen(true);
    }, [setIsEditDialogOpen]);

    const handleEditDialogClose = React.useCallback(() => {
        setIsEditDialogOpen(false);
    }, [setIsEditDialogOpen]);

    const handleEditDialogSave = React.useCallback((rule) => {
        setIsEditDialogOpen(false);
        setRule(rule)
    }, [setIsEditDialogOpen, setRule]);

    return (
        <>
        <RuleDialog constants={constants} open={isEditDialogOpen} onCancel={handleEditDialogClose} onSave={handleEditDialogSave} input={rule} />
            <Card style={{
                height: "100%",
                display: "flex",
                flexDirection: "column",
            }} sx={{ width: '220px', height: '100px', position: 'relative', cursor: "pointer", margin: 1, bgcolor: theme.palette.common.white }}>
                <Typography sx={{ fontSize: 17, margin: 1 }} color="text.secondary" onClick={handleEditDialogOpen}>{rule.name}</Typography>
                <CardContent onClick={handleEditDialogOpen}>
                    <Typography variant="body2">{rule.description}</Typography>
                </CardContent>
                <CardActions disableSpacing sx={{ mt: "auto" }}>
                    <Button color="secondary" onClick={onDelete}>
                        Delete
                    </Button>
                </CardActions>
        </Card>
        </>
    );
}

RuleCard.propTypes = {
    rule: PropTypes.object.isRequired,
    setRule: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
    constants: PropTypes.object.isRequired
};

export default RuleCard;