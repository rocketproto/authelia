import React from "react";

import { Theme, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { useTranslation } from "react-i18next";

import FailureIcon from "@root/components/FailureIcon";

const AccessDenied = function () {
    const styles = useStyles();
    const { t: translate } = useTranslation();
    return (
        <div id="access-denied-stage">
            <div className={styles.iconContainer}>
                <FailureIcon />
            </div>
            <Typography>{translate("Access Denied")}</Typography>
        </div>
    );
};

export default AccessDenied;

const useStyles = makeStyles((theme: Theme) => ({
    iconContainer: {
        marginBottom: theme.spacing(2),
        flex: "0 0 100%",
    },
}));
