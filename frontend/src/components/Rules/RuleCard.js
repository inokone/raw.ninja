import * as React from 'react';
import { Card, Typography, CardContent } from "@mui/material";
import RuleDialog from './RuleDialog';

const RuleCard = ({ rule, setRule, constants }) => {
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
        <Card style={{ flex: 1 }} sx={{ position: 'relative', cursor: "pointer", margin: 1, bgcolor: "lightgrey" }} onClick={handleEditDialogOpen}>
            <CardContent>
                <Typography sx={{ fontSize: 14 }} color="text.secondary">{rule.name}</Typography>
                <Typography variant="body2">{rule.description}</Typography>
            </CardContent>
        </Card>
        </>
    );
}

export default RuleCard;