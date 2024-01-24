import React, { useEffect } from "react";

import { Button, Grid, Theme } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { useTranslation } from "react-i18next";
import { Navigate, To, useNavigate, useSearchParams } from "react-router-dom";

import { IndexRoute, LogoutRoute } from "@constants/Routes";
import LoginLayout from "@layouts/LoginLayout";
import GenericError from "@views/Error/GenericError";
import { Errors } from "@models/Errors";
import { ErrorCode, RedirectionURL } from "@constants/SearchParams";
import ForbiddenError from "@views/Error/NamedErrors/ForbiddenError";
import { useUserInfoPOST } from "@hooks/UserInfo";
import { useAutheliaState } from "@hooks/State";
import { ComponentOrLoading } from "@views/Generic/ComponentOrLoading";
import { AuthenticationLevel } from "@services/State";
import { useConfiguration } from "@hooks/Configuration";

const BaseErrorView = function () {
    const styles = useStyles();
    const navigate = useNavigate();
    const { t: translate } = useTranslation();
    const NavigateTo: To = { pathname: LogoutRoute };
    const [state, fetchState, , fetchStateError] = useAutheliaState();
    const [searchParams] = useSearchParams();
    let searchParamsOverride = new URLSearchParams();
    const [configuration, fetchConfiguration, , fetchConfigurationError] = useConfiguration();
    const [userInfo, fetchUserInfo, , fetchUserInfoError] = useUserInfoPOST();

    // Fetch the state when portal is mounted.
    useEffect(() => {
        fetchState();
    }, [fetchState]);

    // Fetch preferences and configuration when user is authenticated.
    useEffect(() => {
        if (state && state.authentication_level >= AuthenticationLevel.OneFactor) {
            fetchUserInfo();
            fetchConfiguration();
        }
        console.log("STATE:", state);
    }, [state, fetchUserInfo, fetchConfiguration]);

    useEffect(() => {
        switch (searchParams.get(ErrorCode)) {
            case Errors.forbidden: {
                console.log("Forbidden Handling:")
                if (searchParams.has(RedirectionURL)) {
                    searchParamsOverride.set(RedirectionURL, searchParams.get(RedirectionURL) as string);
                    NavigateTo.search = searchParamsOverride.toString();
                }
                if (state?.authentication_level === AuthenticationLevel.Unauthenticated) {
                    NavigateTo.pathname = IndexRoute;
                    // navigate(NavigateTo);
                    console.warn("Navigating to: ", NavigateTo);
                }
                break;
            }
            default: {
                break;
            }
        }
    }, [
        state,
        userInfo,
        fetchUserInfoError,
        searchParams,
    ]);

    const handleErrorCodeJSX = () => {
        switch (searchParams.get(ErrorCode)) {
            case Errors.forbidden: {
                return <ForbiddenError />;
            }
            default: {
                return <GenericError />;
            }
        }
    };

    const handleLogoutClick = () => {
        // NOTE Not working with RD
        navigate(NavigateTo);
    };

    const infoLoaded =
        userInfo !== undefined;

    return (
        <ComponentOrLoading ready={infoLoaded}>
            <LoginLayout
                id="base-error-stage"
                title={`${translate("Hi")} ${userInfo?.display_name || ""}`}
                showBrand
            >
                <Grid container>
                    <Grid item xs={12}>
                        <Button color="secondary" onClick={handleLogoutClick} id="logout-button">
                            {translate("Logout")}
                        </Button>
                    </Grid>
                    <Grid item xs={12} className={styles.mainContainer}>
                        {handleErrorCodeJSX()}
                    </Grid>
                </Grid>
            </LoginLayout>
        </ComponentOrLoading>
    );
};

export default BaseErrorView;

const useStyles = makeStyles((theme: Theme) => ({
    mainContainer: {
        border: "1px solid #d6d6d6",
        borderRadius: "10px",
        padding: theme.spacing(4),
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
}));
