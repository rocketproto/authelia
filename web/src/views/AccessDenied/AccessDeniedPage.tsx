import React from "react";

import { Button, Grid, Theme } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";

import { LogoutRoute as SignOutRoute } from "@constants/Routes";
import LoginLayout from "@layouts/LoginLayout";
import AccessDenied from "@views/AccessDenied/AccessDenied";

export interface Props {
    name: string;
}

// Originally "AuthenticatedView"
const AccessDeniedView = function () {
    const styles = useStyles();
    const navigate = useNavigate();
    const { t: translate } = useTranslation();

    const props: Props = {
        name: "Rocket Proto" // Need to call API
    }

    const handleLogoutClick = () => {
        navigate(SignOutRoute);
    };

    return (
        <LoginLayout id="authenticated-stage" title={`${translate("Hi")} ${props.name}`} showBrand>
            <Grid container>
                <Grid item xs={12}>
                    <Button color="secondary" onClick={handleLogoutClick} id="logout-button">
                        {translate("Logout")}
                    </Button>
                </Grid>
                <Grid item xs={12} className={styles.mainContainer}>
                    <AccessDenied />
                </Grid>
            </Grid>
        </LoginLayout>
    );
};

export default AccessDeniedView;

const useStyles = makeStyles((theme: Theme) => ({
    mainContainer: {
        border: "1px solid #d6d6d6",
        borderRadius: "10px",
        padding: theme.spacing(4),
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
}));
